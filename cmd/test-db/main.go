package main

import (
	"log"
	"time"

	"ProjectNotification/internal/models"
	"ProjectNotification/internal/notification"
	"ProjectNotification/internal/postgres"
)

func main() {
	db := postgres.NewDB()
	repo := notification.NewRepository(db)

	notification := models.Notification{
		ID:        "test-1",
		UserID:    "user-123",
		Type:      "email",
		Payload:   map[string]interface{}{"message": "hello"},
		Status:    "pending",
		CreatedAt: time.Now().Unix(),
	}

	err := repo.Save(notification)
	if err != nil {
		log.Fatal(err)
	}
}
