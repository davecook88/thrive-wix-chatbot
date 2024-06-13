package chatgpt

type ChatGPTRequest struct {
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required,dive"`
	Stream   bool      `json:"stream" binding:"default:false"`
}

type Role string

const (
	UserRole      Role = "user"
	AssistantRole Role = "assistant"
	SystemRole    Role = "system"
)

type Message struct {
	Role    Role   `json:"role" binding:"required,oneof=user assistant"`
	Content string `json:"content" binding:"required"`
}

type ChatGPTResponse struct {
	Choices           []ResponseChoice `json:"choices"`
	Error             *string          `json:"error"`
	ID                string           `json:"id"`
	Obejct            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	Usage             Usage            `json:"usage"`
	SystemFingerprint string           `json:"system_fingerprint"`
}

type ResponseChoice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	LogProbs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Delta struct {
	Content *string `json:"content"`
}

type StreamingResponseChoice struct {
	Index        int         `json:"index"`
	Delta        Delta       `json:"delta"`
	LogProbs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type StreamingResponse struct {
	Id                string                    `json:"id"`
	Object            string                    `json:"object"`
	Created           int64                     `json:"created"`
	Model             string                    `json:"model"`
	SystemFingerprint string                    `json:"system_fingerprint"`
	Choices           []StreamingResponseChoice `json:"choices"`
}
