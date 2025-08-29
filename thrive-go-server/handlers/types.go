package handlers

type UserMessage struct {
	Message  string `json:"message" binding:"required"`
	UserTime string `json:"userTime"`
}
