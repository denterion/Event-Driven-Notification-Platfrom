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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	response, err := client.SendEvent(ctx, &pb.EventRequest{
		EventType: "user_registered",
		UserId:    "123",
		Email:     "test@mail.com",
		Payload:   "new user signup",
	})

	if err != nil {
		log.Fatal("SendEvent error", err)
	}

	log.Println("Server response:", response.Status)
}
