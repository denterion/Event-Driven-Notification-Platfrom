package kafka

import (
	"context"
	"encoding/json"
	"log"

	"ProjectNotification/internal/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokerAddress, topic string) *Producer {
	return &Producer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
		}),
	}
}

func (p *Producer) Send(ctx context.Context, event models.Event) error {

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(event.UserID),
		Value: data,
	}
	err = p.Writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Println("failed to write message:", err)
		return err
	}
	log.Println("event sent to Kafka:", string(data))
	return nil
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}
