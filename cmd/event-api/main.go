package main

import (
	"context"
	"log"
	"net"

	pb "ProjectNotification/api/proto"
	"ProjectNotification/internal/kafka"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEventServiceServer
	producer *kafka.Producer
}

func (s *server) SendEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {

	log.Println("event received: ", req.EventType, req.UserId)

	err := s.producer.Send(ctx, req.UserId, req.Payload)
	if err != nil{
		return nil, err
	}

	return &pb.EventResponse{
		Status: "event accepted",
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	producer := kafka.NewProducer("localhost:9092", "notifications.events")
	defer producer.Close()

	grpcServer := grpc.NewServer()

	pb.RegisterEventServiceServer(grpcServer, &server{producer: producer})

	log.Println("gRPC server started on 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
