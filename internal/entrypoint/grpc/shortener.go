package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "service-url-shortener/internal/entrypoint/grpc/shortener_proto"
	internal "service-url-shortener/internal/errors"
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

// Create returns existed or make new short URL from given original.
func (s *shortenerRoutes) Create(ctx context.Context, request *pb.ShortenerData) (*pb.ShortenerData, error) {
	err := validateURL(request.URL)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	short, err := s.t.Shorten(ctx, request.URL)
	if err != nil {
		s.l.Error(err, "grpc - Shortener - Create")
		return nil, status.Error(codes.Internal, "shortener service problems")
	}

	return &pb.ShortenerData{URL: short}, nil
}

// Get returns original URL from given short if exists.
func (s *shortenerRoutes) Get(ctx context.Context, request *pb.ShortenerData) (*pb.ShortenerData, error) {
	err := validateURL(request.URL)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	original, err := s.t.Lengthen(ctx, request.URL)
	if err != nil {

		if errors.Is(err, internal.ErrLengthTooHigh) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, internal.ErrNotFoundURL) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.l.Error(err, "grpc - Shortener - Get")
		return nil, status.Error(codes.Internal, "shortener service problems")
	}

	return &pb.ShortenerData{URL: original}, nil
}
