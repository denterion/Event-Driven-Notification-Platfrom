package models

type Event struct {
	EventType string      `json:"event_type"`
	UserID    string      `json:"user_id"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type UserRegisteredPayload struct {
	Email string `json:"email"`
}
