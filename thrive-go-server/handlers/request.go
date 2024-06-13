package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"thrive/server/auth"
	"thrive/server/chatgpt"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)

func NewChatGPTRequest(messages []chatgpt.Message) *chatgpt.ChatGPTRequest {
	return &chatgpt.ChatGPTRequest{
		Model:    "gpt-4o",
		Messages: messages,
		Stream:   false,
	}
}

func CallChatGPT(c *gin.Context, messages []chatgpt.Message) (*chatgpt.Message, error) {
	// c.Header("Content-Type", "text/event-stream")
	// c.Header("Cache-Control", "no-cache")
	// c.Header("Connection", "keep-alive")
	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))

	// if the first message is not the system message, add a system message
	if len(messages) == 0 || messages[0].Role != chatgpt.SystemRole {
		messages = append([]chatgpt.Message{{Role: chatgpt.SystemRole, Content: chatgpt.SystemMessage}, {
			Role:    chatgpt.AssistantRole,
			Content: `Hi! I'm Diego, your assistant.\nI'm going to ask you a few questions to check your Spanish level.\nCan you introduce yourself in Spanish?\n\n¿Cómo te llamas? ¿De dónde eres? ¿Qué haces en tu tiempo libre?`,
		}}, messages...)
	}
	jsonData, err := json.Marshal(NewChatGPTRequest(messages))

	if err != nil {
		return nil, errors.New("failed to marshal request")
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to make request")
	}
	defer resp.Body.Close()

	var chatGPTResponse chatgpt.ChatGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatGPTResponse); err != nil {
		return nil, errors.New("failed to parse response")
	}

	// scanner := bufio.NewScanner(resp.Body)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	if strings.HasPrefix(line, "data: ") {
	// 		data := strings.TrimPrefix(line, "data: ")
	// 		if data != "[DONE]" {
	// 			var responseData chatgpt.StreamingResponse
	// 			if err := json.Unmarshal([]byte(data), &responseData); err != nil {
	// 				c.SSEvent("error", gin.H{"error": "Failed to unmarshal response data"})
	// 				return nil, errors.New("failed to unmarshal response data")
	// 			}
	// 			if responseData.Choices[0].Delta.Content != nil {
	// 				responseMessage.Content += *responseData.Choices[0].Delta.Content
	// 				c.SSEvent("message", gin.H{"content": responseMessage.Content})
	// 			}
	// 		}
	// 	}
	// }

	return &chatGPTResponse.Choices[0].Message, nil

}

func PostMessageHandler(c *gin.Context) {
	var request UserMessage
	println("PostMessageHandler")
	memberProfile := auth.ValidateWixUser(c)

	if memberProfile == nil || memberProfile.Member.ID == "" {
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

	existingMessages, err = dbClient.GetChat(c, memberProfile.Member.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages := append(*existingMessages, chatgpt.Message{Role: chatgpt.UserRole, Content: request.Message})

	chatGPTResponseMessage, err := CallChatGPT(c, messages)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, *chatGPTResponseMessage)

	if err := dbClient.UpdateChat(c, memberProfile.Member.ID, db.SavedChatRecord{
		Messages: messages,
		MemberId: memberProfile.Member.ID,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Marshal the chatGPTResponse struct back to JSON
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal ChatGPT response"})
		return
	}

	// Print the full JSON response
	fmt.Println(string(jsonResponse))

	c.JSON(200, messages)

}

func GetChatHandler(c *gin.Context) {
	memberProfile := auth.ValidateWixUser(c)
	if memberProfile == nil || memberProfile.Member.ID == "" {
		c.JSON(400, gin.H{"error": "No user instance"})
		return
	}

	dbClient, err := db.NewClient(c, "thrive-chat")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages, err := dbClient.GetChat(c, memberProfile.Member.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, messages)
}
