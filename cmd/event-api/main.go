package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "ProjectNotification/api/proto"
	"ProjectNotification/internal/config"
	"ProjectNotification/internal/kafka"
	"ProjectNotification/internal/models"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEventServiceServer
	producer *kafka.Producer
}

func (s *server) SendEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {

	log.Println("event received: ", req.EventType, req.UserId)

	event := models.Event{
		EventType: req.EventType,
		UserID:    req.UserId,
		Timestamp: time.Now().Unix(),
		Payload: models.UserRegisteredPayload{
			Email: req.Email,
		},
	}

	err := s.producer.Send(ctx, event)
	if err != nil {
		return nil, err
	}

	return &pb.EventResponse{
		Status: "event accepted",
	}, nil
}

func main() {
	config.OverloadDotEnv()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "notifications.events")

	producer := kafka.NewProducer(broker, topic)
	defer producer.Close()

	grpcServer := grpc.NewServer()

	pb.RegisterEventServiceServer(grpcServer, &server{producer: producer})

	log.Println("gRPC server started on 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
