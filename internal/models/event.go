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

type OrderCreatedPayload struct {
	OrderID string `json:"order_id"`
	Email   string `json:"email"`
	Details string `json:"details"`
}

type PaymentSucceededPayload struct {
	OrderID string `json:"order_id"`
	Details string `json:"details"`
}

type PaymentFailedPayload struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}

type PasswordResetRequestedPayload struct {
	Email string `json:"email"`
}
