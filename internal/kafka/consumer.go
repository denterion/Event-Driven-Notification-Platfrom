package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokerAddress, topic, groupID string){

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic: topic, 
		GroupID: groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	log.Println("Kafka consumer started on topic:", topic)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil{
			log.Println("error reading message:", err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("Message received: key=%s values=%s offset=%d partion=%d\n",
		string(m.Key), string(m.Value), m.Offset, m.Partition)
	}
}