package handlers

import "thrive/server/chatgpt"

func NewChatGPTRequestCheckLevel(messages *[]chatgpt.Message) *chatgpt.ChatGPTRequest {
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

func NewChatGPTRequestConversation(messages *[]chatgpt.Message) *chatgpt.ChatGPTRequest {
	toolsArray := []chatgpt.Tools{
		*chatgpt.NewToolFunction(
			"getServices",
			`Get all available services to suggest the best option to the user.`,
			map[string]interface{}{}),
	}
	return &chatgpt.ChatGPTRequest{
		Model:    "gpt-4o",
		Messages: *messages,
		Stream:   false,
		Tools:    toolsArray,
	}
}
