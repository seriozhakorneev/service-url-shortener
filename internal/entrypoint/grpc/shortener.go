package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "service-url-shortener/internal/entrypoint/grpc/shortener_proto"
	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

type shortenerRoutes struct {
	pb.UnimplementedShortenerServer
	t usecase.Shortener
	l logger.Interface
}

func newShortenerRoutes(s grpc.ServiceRegistrar, t usecase.Shortener, l logger.Interface) {
	r := &shortenerRoutes{t: t, l: l}

	pb.RegisterShortenerServer(s, r)
}

// Create returns short URL from given original.
func (s *shortenerRoutes) Create(ctx context.Context, data *pb.ShortenerData) (*pb.ShortenerData, error) {
	//s.t.Shorten()
	return &pb.ShortenerData{URL: fmt.Sprint("SHORT_URL_FROM_SERVER:", data.GetURL())}, nil
}

// Get returns original URL from given short.
func (s *shortenerRoutes) Get(ctx context.Context, data *pb.ShortenerData) (*pb.ShortenerData, error) {
	//s.t.Lengthen()
	return &pb.ShortenerData{URL: fmt.Sprint("ORIGINAL_URL_FROM_SERVER:", data.GetURL())}, nil
}
