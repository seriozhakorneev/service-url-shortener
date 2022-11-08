package grpc

import (
	"google.golang.org/grpc"

	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

// NewRouter -.
func NewRouter(s grpc.ServiceRegistrar, t usecase.Shortener, l logger.Interface) {
	{
		newShortenerRoutes(s, t, l)
	}
}
