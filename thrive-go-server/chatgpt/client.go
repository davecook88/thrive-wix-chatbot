package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ChatGPTClient struct {
	ApiKey string
	*http.Client
}

func NewChatGPTClient(apiKey string) *ChatGPTClient {
	return &ChatGPTClient{
		ApiKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *ChatGPTClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.Client.Do(req)
}

func (client *ChatGPTClient) MakeRequest(req *ChatGPTRequest) (*ChatGPTResponse, error) {
	jsonData, err := json.Marshal(req)

	// print the request JSON
	fmt.Println(string(jsonData))

	if err != nil {
		return nil, errors.New("failed to marshal request")
	}

	r, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request")
	}
	resp, err := client.Do(r)
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
	var chatGPTResponse ChatGPTResponse
	if err := json.NewDecoder(respBody).Decode(&chatGPTResponse); err != nil {
		return nil, errors.New("failed to parse response")
	}
	return &chatGPTResponse, nil

}
