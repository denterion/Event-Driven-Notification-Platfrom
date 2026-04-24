package main

import (
	"ProjectNotification/internal/config"
	"ProjectNotification/internal/delivery/telegram"
	"ProjectNotification/internal/kafka"
	"ProjectNotification/internal/notification"
	"ProjectNotification/internal/postgres"
	"ProjectNotification/internal/redis"
	"context"
	"os"
	"time"
)

func main() {
	config.OverloadDotEnv()

	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "notifications.events")
	groupID := getEnv("KAFKA_GROUP_ID", "test-consumer-group")

	db := postgres.NewDB()
	defer db.Close()

	redisClient := redis.NewFromEnv()
	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = redisClient.Ping(ctx)
	}

	repo := notification.NewRepository(db)
	handler := notification.NewHandler(repo, redisClient, telegram.NewFromEnv())

	kafka.StartConsumer(broker, topic, groupID, handler)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
