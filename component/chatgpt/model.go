package chatgpt

import "time"

type Conversation struct {
	ID            string
	Nickname      string // client request id
	LastMessageID string
	LastAskTime   time.Time
	Server        *Server
}

type Response struct {
	Message        string       `json:"message"`
	ConversationID string       `json:"conversation_id"`
	MessageID      string       `json:"message_id"`
	FunctionCall   FunctionCall `json:"function_call"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string      `json:"type"`
		Properties interface{} `json:"properties"`
		Required   []string    `json:"required"`
	} `json:"parameters"`
}
