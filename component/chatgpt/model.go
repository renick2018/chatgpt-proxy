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
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
}
