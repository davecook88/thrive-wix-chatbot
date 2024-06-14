package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"thrive/server/auth"
	"thrive/server/chatgpt"
	"thrive/server/db"
	"thrive/server/wix"

	"github.com/gin-gonic/gin"
)

func NewChatGPTRequest(messages *[]chatgpt.Message) *chatgpt.ChatGPTRequest {
	toolsArray := []chatgpt.Tools{
		*chatgpt.NewToolFunction(
			"estimateUserLevel",
			`Estimate the user's Spanish level based on their responses. 
			This should be included with every response with the current best estimate of the user's level.`,
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"estimatedLevel": map[string]interface{}{
						"type": "string",
						"enum": []string{"A1", "A2", "B1", "B2", "C1"},
					},
				},
			}),
		*chatgpt.NewToolFunction(
			"addNotes",
			`Add notes about the user for the CRM.
				Summarize interests, goals, and any other relevant information.
				If the user makes consistent mistakes, note them here.
				Keep notes to a few words at most.`,
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"notes": map[string]interface{}{
						"type": "string",
					},
				},
			}),
	}
	toolChoice := "required"
	return &chatgpt.ChatGPTRequest{
		Model:      "gpt-4o",
		Messages:   *messages,
		Stream:     false,
		Tools:      toolsArray,
		ToolChoice: &toolChoice,
	}
}

func handleLevelEstimation(level string, wixContact *wix.Contact) {
	existingLabels := wixContact.Info.LabelKeys
	// remove any existing level labels, they begin with custom.level-
	filteredLabels := []string{}
	for _, label := range existingLabels.Items {
		if !strings.HasPrefix(label, "custom.level-") {
			filteredLabels = append(filteredLabels, label)
		}
	}
	// append new level label
	filteredLabels = append(filteredLabels, "custom.level-"+strings.ToLower(level))
	wixContact.Info.LabelKeys.Items = filteredLabels
}

func handleToolCalls(toolCalls *[]chatgpt.ToolCall, wixMember *wix.WixMember) {
	wixClient := wix.NewWixClient()
	contact, err := wixClient.GetContact(wixMember.ContactId)
	if err != nil {
		fmt.Println("failed to get contact", err)
		return
	}
	for _, toolCall := range *toolCalls {
		switch toolCall.Function.Name {
		case "estimateUserLevel":
			// Arguments should be a JSON string like this {\"estimatedLevel\":\"C1\"}
			arguments := map[string]string{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			handleLevelEstimation(arguments["estimatedLevel"], contact)

		case "addNotes":
			// Arguments should be a JSON string like this {\"notes\":\"This user is interested in learning Spanish for travel.\"}
			arguments := map[string]string{}
			fmt.Println("addNotes", toolCall.Function.Arguments)
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			// Add notes to the contact
			existingNotes := contact.Info.ExtendedFields.Items["custom.notes"]
			newNotes := existingNotes + "\n*" + arguments["notes"] + "*"
			contact.Info.ExtendedFields.Items["custom.notes"] = newNotes

		default:
			continue
		}
	}
	wixClient.UpdateContact(contact.ID, contact.Revision, contact.Info)

}

func makeToolCall(client *chatgpt.ChatGPTClient, chatGPTRequest *chatgpt.ChatGPTRequest, wixMember *wix.WixMember) {
	chatGPTResponse, err := client.MakeRequest(chatGPTRequest)
	if err != nil {
		return
	}
	choice := chatGPTResponse.Choices[0]
	if choice.Message.ToolCalls != nil {
		go handleToolCalls(&choice.Message.ToolCalls, wixMember)
	}
}

func CallChatGPT(c *gin.Context, messages *[]chatgpt.Message, wixMember *wix.WixMember) (*chatgpt.Message, error) {

	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))
	_messages := *messages

	// if the first message is not the system message, add a system message
	if len(_messages) == 0 || _messages[0].Role != chatgpt.SystemRole {
		_messages = append(chatgpt.InitialMessages, _messages...)
	}
	newChatGPTRequest := NewChatGPTRequest(&_messages)
	userMessageCount := 0
	for _, message := range _messages {
		if message.Role == chatgpt.UserRole {
			userMessageCount++
		}
	}

	// make the tool call on every 3rd message
	if userMessageCount > 3 && userMessageCount%3 == 0 {
		go makeToolCall(client, newChatGPTRequest, wixMember)
	}
	// copy newChatGPTRequest
	newChatGPTRequest = NewChatGPTRequest(&_messages)
	// if no content is returned, make another request with no function call
	newChatGPTRequest.Tools = nil
	newChatGPTRequest.ToolChoice = nil
	chatGPTResponse, err := client.MakeRequest(newChatGPTRequest)
	if err != nil {
		return nil, errors.New("failed to make request")
	}
	choice := chatGPTResponse.Choices[0]

	fmt.Println("returning message", choice)

	return choice.Message.ToMessage(), nil

}

func PostMessageHandler(c *gin.Context) {
	var request UserMessage
	println("PostMessageHandler")
	memberProfile := auth.ValidateWixUser(c)

	if memberProfile == nil || memberProfile.ID == "" {
		return
	}
	fmt.Println(memberProfile)

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

	messages = append(messages, *chatGPTResponseMessage)

	if err := dbClient.UpdateChat(c, memberProfile.ID, db.SavedChatRecord{
		Messages:   messages,
		MemberId:   memberProfile.ID,
		MemberName: memberProfile.Profile.Nickname,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	println("PostMessageHandler finished")

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
