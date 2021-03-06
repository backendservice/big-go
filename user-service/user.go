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
	return &pb.UserResponse{Message: "User (" + in.Name + ") register success", Code: 200}, nil
}

func (s *server) FindUser(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	filtereduser, count := filter(in)
	return &pb.FindResponse{Count: count, Result: filtereduser, Message: "Hello " + in.Religion, Code: 200}, nil
}

func filter(condition *pb.FindRequest) ([]*pb.UserRequest, int32) {
	var tmpuser []*pb.UserRequest
	var count int32
	for i := 0; i < len(users); i++ {
		if (strings.Compare(condition.GetAgeType(), "") == 0 || strings.Compare(condition.GetAgeType(), "E") == 0 && condition.GetAge() == users[i].Age ||
			strings.Compare(condition.GetAgeType(), "U") == 0 && condition.GetAge() < users[i].Age ||
			strings.Compare(condition.GetAgeType(), "D") == 0 && condition.GetAge() > users[i].Age) &&
			(condition.GetDistance() == 0 || checkDistance(condition.GetDistance(), float64(condition.GetLatitude()), float64(users[i].Latitude), float64(condition.GetLongitude()), float64(users[i].Longitude))) &&
			(strings.Compare(condition.GetGender(), "") == 0 || strings.Compare(condition.GetGender(), users[i].Gender) == 0) &&
			(strings.Compare(condition.GetReligion(), "") == 0 || strings.Compare(condition.GetReligion(), users[i].Religion) == 0) &&
			(strings.Compare(condition.GetNationality(), "") == 0 || strings.Compare(condition.GetNationality(), users[i].Nationality) == 0) {
			tmpuser = append(tmpuser, users[i])
			count++
		}
	}
	return tmpuser, count
}

func checkDistance(distance int32, latitude1 float64, latitude2 float64, longitude1 float64, longitude2 float64) bool {

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
