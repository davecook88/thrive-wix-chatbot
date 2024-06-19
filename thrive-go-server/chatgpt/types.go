package chatgpt

type ChatGPTRequest struct {
	Model      string    `json:"model" binding:"required"`
	Messages   []Message `json:"messages" binding:"required,dive"`
	Stream     bool      `json:"stream" binding:"default:false"`
	Tools      []Tools   `json:"tools"`
	ToolChoice *string   `json:"tool_choice"`
}

type Role string

const (
	UserRole      Role = "user"
	AssistantRole Role = "assistant"
	SystemRole    Role = "system"
)

type Message struct {
	Role       Role    `json:"role" binding:"required,oneof=user assistant"`
	Content    string  `json:"content" binding:"required"`
	ToolCallId *string `json:"tool_call_id"`
	Name       *string `json:"name"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Id       string           `json:"id"`
	Function ToolCallFunction `json:"function"`
}

type ResponseMessage struct {
	Role      Role       `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

func (m *ResponseMessage) ToMessage() *Message {
	return &Message{
		Role:    m.Role,
		Content: m.Content,
	}
}

type ToolType string

const (
	FunctionType ToolType = "function"
)

type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description" binding:"required"`
	Parameters  interface{} `json:"parameters"`
}

type Tools struct {
	ToolType ToolType     `json:"type"`
	Funtion  ToolFunction `json:"function"`
}

func NewToolFunction(name, description string, parameters interface{}) *Tools {
	newFunc := ToolFunction{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}
	return &Tools{ToolType: FunctionType, Funtion: newFunc}
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
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	LogProbs     interface{}     `json:"logprobs"`
	FinishReason string          `json:"finish_reason"`
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
