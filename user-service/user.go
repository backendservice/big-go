package main

import (
	"context"
	"log"
	"net"

	pb "github.com/backendservice/big-go/big"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func (s *server) RegistUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{Message: "Hello " + in.Name}, nil
}

func (s *server) FindUser(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	return &pb.FindResponse{Message: "Hello " + in.Religion}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBigServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
