package models

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Payload   interface{}
	Status    string
	CreatedAt int64
}
