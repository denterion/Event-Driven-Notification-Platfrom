package kafka

import (
	"context"
	"log"

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

func (p *Producer) Send(ctx context.Context, key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	err := p.Writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Println("failed to write message:", err)
		return err
	}
	log.Println("event sent to Kafka:", string(value))
	return nil
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}
