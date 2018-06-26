package main

import (
	"context"
	"log"
	"math"
	"net"
	"strings"

	pb "github.com/backendservice/big-go/big"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

var users []*pb.UserRequest

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func (s *server) RegistUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	users = append(users, in)
	return &pb.UserResponse{Message: "User (" + in.Name + ") register success"}, nil
}

func (s *server) FindUser(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	filtereduser, count := Filter(in)
	return &pb.FindResponse{Count: count, Result: filtereduser, Message: "Hello " + in.Religion}, nil
}

func Filter(condition *pb.FindRequest) ([]*pb.UserRequest, int32) {
	var tmpuser []*pb.UserRequest
	var count int32
	for i := 0; i < len(users); i++ {
		if condition.GetStartage() < users[i].Age && condition.GetEndage() > users[i].Age &&
			CheckDistance(condition.GetDistance(), condition.GetLatitude(), users[i].Latitude, condition.GetLongitude(), condition.GetDistance()) &&
			strings.Compare(condition.GetGender(), users[i].Gender) == 0 &&
			strings.Compare(condition.GetReligion(), users[i].Religion) == 0 &&
			strings.Compare(condition.GetNationality(), users[i].Nationality) == 0 {
			tmpuser = append(tmpuser, users[i])
			count++
		}
	}
	return tmpuser, count
}

func CheckDistance(distance int32, latitude1 float64, latitude2 float64, longitude1 float64, longitude2 float64) bool {

	var dis = math.Acos(math.Sin(latitude1)*math.Sin(latitude2) + math.Cos(latitude1)*math.Cos(latitude2)*math.Cos(longitude1-longitude2))

	if dis > float64(distance) {
		return false
	}
	return true
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
