package main

import (
	"ProjectNotification/internal/kafka"
)

func main(){
	broker := "localhost:9092"
	topic := "notifications.events"
	groupID := "test-consumer-group"

	kafka.StartConsumer(broker, topic, groupID)
}