package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "service-url-shortener/pkg/grpc/shortener"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type shortenerServer struct {
	pb.UnimplementedShortenerServer
}

// Create returns short URL from given original.
func (s *shortenerServer) Create(ctx context.Context, data *pb.ShortenerData) (*pb.ShortenerData, error) {
	return &pb.ShortenerData{URL: fmt.Sprint("SHORT_URL_FROM_SERVER:", data.GetURL())}, nil
}

// Get returns original URL from given short.
func (s *shortenerServer) Get(ctx context.Context, data *pb.ShortenerData) (*pb.ShortenerData, error) {
	return &pb.ShortenerData{URL: fmt.Sprint("ORIGINAL_URL_FROM_SERVER:", data.GetURL())}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShortenerServer(grpcServer, &shortenerServer{})
	grpcServer.Serve(lis)
}
