package client

import "time"

//Client interface for chat clients
type Client interface {
	Send(m *BroadcastMessage)
	GetID() string
	Close()
}

// BroadcastMessage message from a client targeted to brodcasting to clients
type BroadcastMessage struct {
	MessageID int
	//Message user message
	Message,
	//UserID user ID
	UserID string
	//Timestamp
	Timestamp string
}

func NewBroadcastMessage(userID, message string) *BroadcastMessage {
	return &BroadcastMessage{
		UserID:    userID,
		Message:   message,
		Timestamp: time.Now().Format("3:04PM"),
	}
}
