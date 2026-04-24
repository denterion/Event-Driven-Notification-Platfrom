package kafka

import (
	"ProjectNotification/internal/models"
	"ProjectNotification/internal/notification"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokerAddress, topic, groupID string, handler *notification.Handler) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	log.Println("Kafka consumer started on topic:", topic)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("error reading message:", err)
			time.Sleep(time.Second)
			continue
		}
		var event models.Event

		err = json.Unmarshal(m.Value, &event)
		if err != nil {
			log.Println("failed to parse event:", err)
			continue
		}

		log.Println("EVENT RECEIVED:")
		log.Println("type:", event.EventType)
		log.Println("user:", event.UserID)
		log.Println("payload:", event.Payload)

		if handler != nil {
			if err := handler.Handle(context.Background(), event); err != nil {
				if err == notification.ErrUnsupportedEventType {
					log.Println("event ignored (unsupported type):", event.EventType)
					continue
				}
				log.Println("failed to handle event:", err)
			}
		}
	}
}
