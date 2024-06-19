package handlers

import (
	"errors"
	"fmt"
	"os"
	"thrive/server/auth"
	"thrive/server/chatgpt"
	"thrive/server/db"
	"thrive/server/wix"
	"time"

	"github.com/gin-gonic/gin"
)

var MAX_USER_MESSAGES = 20

func makeToolCall(client *chatgpt.ChatGPTClient, chatGPTRequest *chatgpt.ChatGPTRequest, wixMember *wix.WixMember) {
	println("making tool call")
	chatGPTResponse, err := client.MakeRequest(chatGPTRequest)
	if err != nil {
		return
	}
	choice := chatGPTResponse.Choices[0]
	if choice.Message.ToolCalls != nil {
		println("handling tool calls")
		go handleToolCalls(&choice.Message.ToolCalls, wixMember)
	}
}

func CallChatGPT(c *gin.Context, messages *[]chatgpt.Message, wixMember *wix.WixMember) (*[]chatgpt.Message, error) {

	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))
	_messages := *messages

	// if the first message is not the system message, add a system message
	if len(_messages) == 0 || _messages[0].Role != chatgpt.SystemRole {
		_messages = append(chatgpt.InitialMessages, _messages...)
	}
	newChatGPTRequest := NewChatGPTRequestConversation(&_messages)
	userMessageCount := 0
	for _, message := range _messages {
		if message.Role == chatgpt.UserRole {
			userMessageCount++
		}
	}

	// if user message count is greater than MAX_MESSAGE_COUNT
	// return a standard message telling the user that they have reached the limit
	if userMessageCount > MAX_USER_MESSAGES {
		return &[]chatgpt.Message{{
			Role:    chatgpt.AssistantRole,
			Content: "Your placement test is complete. One of our teachers will reach out shortly via email.",
		}}, nil
	}

	// make the tool call on every 3rd message
	println("userMessageCount", userMessageCount)
	if userMessageCount == 1 || userMessageCount%3 == 0 {
		go makeToolCall(client, NewChatGPTRequestCheckLevel(messages), wixMember)
	}
	// copy newChatGPTRequest
	responseMessages := []chatgpt.Message{}

	// if no content is returned, make another request with no function call
	chatGPTResponse, err := client.MakeRequest(newChatGPTRequest)
	if err != nil {
		return nil, errors.New("failed to make request")
	}
	responseMessages = append(responseMessages, *chatGPTResponse.Choices[0].Message.ToMessage())
	choice := chatGPTResponse.Choices[0]
	if choice.Message.ToolCalls != nil {
		fmt.Println("handling tool calls", choice.Message.ToolCalls)
		funcResponseStr := handleToolCallsWithResponse(&choice.Message.ToolCalls)
		responseMsg := chatgpt.Message{
			Role:       chatgpt.AssistantRole,
			Content:    *funcResponseStr,
			ToolCallId: &choice.Message.ToolCalls[0].Id,
			Name:       &choice.Message.ToolCalls[0].Function.Name,
		}
		responseMessages = append(responseMessages, responseMsg)
		// append the response message to the messages and call chatgpt again
		_messages = append(_messages, responseMsg)
		newChatGPTRequest = NewChatGPTRequestConversation(&_messages)
		chatGPTResponse, err = client.MakeRequest(newChatGPTRequest)
		if err != nil {
			return nil, errors.New("failed to make request")
		}
		responseMessages = append(responseMessages, *chatGPTResponse.Choices[0].Message.ToMessage())
	}

	return &responseMessages, nil
}

func PostMessageHandler(c *gin.Context) {
	var request UserMessage
	memberProfile := auth.ValidateWixUser(c)

	if memberProfile == nil || memberProfile.ID == "" {
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dbClient, err := db.NewClient(c, "thrive-chat")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var existingMessages *[]chatgpt.Message

	existingMessages, err = dbClient.GetChat(c, memberProfile.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(*existingMessages) == 0 {
		existingMessages = &chatgpt.InitialMessages
	}

	messages := append(*existingMessages, chatgpt.Message{Role: chatgpt.UserRole, Content: request.Message})

	chatGPTResponseMessage, err := CallChatGPT(c, &messages, memberProfile)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, *chatGPTResponseMessage...)

	if err := dbClient.UpdateChat(c, memberProfile.ID, db.SavedChatRecord{
		Messages:    messages,
		MemberId:    memberProfile.ID,
		MemberName:  memberProfile.Profile.Nickname,
		LastUpdated: time.Now().Format(time.RFC3339),
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, messages)

}

func GetChatHandler(c *gin.Context) {
	memberProfile := auth.ValidateWixUser(c)
	if memberProfile == nil || memberProfile.ID == "" {
		c.JSON(400, gin.H{"error": "No user instance"})
		return
	}

	dbClient, err := db.NewClient(c, "thrive-chat")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages, err := dbClient.GetChat(c, memberProfile.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(*messages) == 0 {
		newMessages := chatgpt.InitialMessages[1:]
		messages = &newMessages
	}

	c.JSON(200, messages)
}
