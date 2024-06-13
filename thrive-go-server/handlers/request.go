package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"thrive/server/auth"
	"thrive/server/chatgpt"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)

func NewChatGPTRequest(messages []chatgpt.Message) *chatgpt.ChatGPTRequest {
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
	}
	return &chatgpt.ChatGPTRequest{
		Model:      "gpt-4o",
		Messages:   messages,
		Stream:     false,
		Tools:      toolsArray,
		ToolChoice: "required",
	}
}

func handleToolCalls(toolCalls []chatgpt.ToolCall) {
	for _, toolCall := range toolCalls {
		switch toolCall.Function.Name {
		case "estimateUserLevel":
			// Arguments should be a JSON string like this {\"estimatedLevel\":\"C1\"}
			println("Estimate user level", toolCall.Function.Arguments)
			arguments := map[string]string{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			fmt.Println(arguments)
			fmt.Println("Estimated level:", arguments["estimatedLevel"])
		}
	}

}

func CallChatGPT(c *gin.Context, messages []chatgpt.Message) (*chatgpt.Message, error) {
	// c.Header("Content-Type", "text/event-stream")
	// c.Header("Cache-Control", "no-cache")
	// c.Header("Connection", "keep-alive")
	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))

	// if the first message is not the system message, add a system message
	if len(messages) == 0 || messages[0].Role != chatgpt.SystemRole {
		messages = append(chatgpt.InitialMessages, messages...)
	}
	jsonData, err := json.Marshal(NewChatGPTRequest(messages))

	// print the request JSON
	fmt.Println(string(jsonData))

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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}

	// Create a new reader with the response body
	respBody := bytes.NewReader(body)

	// Pretty-print the JSON response
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return nil, errors.New("failed to pretty-print JSON")
	}

	// Print the pretty-printed JSON to the console
	fmt.Println(prettyJSON.String())

	// Create a new decoder with the response body reader
	var chatGPTResponse chatgpt.ChatGPTResponse
	if err := json.NewDecoder(respBody).Decode(&chatGPTResponse); err != nil {
		return nil, errors.New("failed to parse response")
	}

	fmt.Println("decoded response")
	choice := chatGPTResponse.Choices[0]
	if choice.Message.ToolCalls != nil {
		go handleToolCalls(choice.Message.ToolCalls)
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

	return choice.Message.ToMessage(), nil

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

	if len(*existingMessages) == 0 {
		existingMessages = &chatgpt.InitialMessages
	}

	messages := append(*existingMessages, chatgpt.Message{Role: chatgpt.UserRole, Content: request.Message})

	chatGPTResponseMessage, err := CallChatGPT(c, messages)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, *chatGPTResponseMessage)

	if err := dbClient.UpdateChat(c, memberProfile.Member.ID, db.SavedChatRecord{
		Messages:   messages,
		MemberId:   memberProfile.Member.ID,
		MemberName: memberProfile.Member.Profile.Nickname,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	println("PostMessageHandler finished")

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

	if len(*messages) == 0 {
		newMessages := chatgpt.InitialMessages[1:]
		messages = &newMessages
	}

	c.JSON(200, messages)
}
