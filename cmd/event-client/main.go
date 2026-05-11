package main

import (
	"context"
	"log"
	"time"

	pb "ProjectNotification/api/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("connection error:", err)
	}
	defer conn.Close()

	client := pb.NewEventServiceClient(conn)

	events := []*pb.EventRequest{
		{
			EventType: "user_registered",
			UserId:    "123",
			Email:     "test@mail.com",
			Payload:   "new user signup",
		},
		{
			EventType: "order_created",
			UserId:    "123",
			Email:     "test@mail.com",
			Payload:   "ORD-1001",
		},
		{
			EventType: "payment_succeeded",
			UserId:    "123",
			Email:     "test@mail.com",
			Payload:   "ORD-1001",
		},
		{
			EventType: "payment_failed",
			UserId:    "456",
			Email:     "buyer@mail.com",
			Payload:   "ORD-1002",
		},
		{
			EventType: "password_reset_requested",
			UserId:    "789",
			Email:     "reset@mail.com",
			Payload:   "password reset",
		},
	}

	for _, event := range events {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		response, err := client.SendEvent(ctx, event)
		cancel()
		if err != nil {
			log.Fatal("SendEvent error", err)
		}
		log.Println(event.EventType, ":", response.Status)
	}
}
