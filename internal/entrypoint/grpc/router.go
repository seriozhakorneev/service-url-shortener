package grpc

import (
	"google.golang.org/grpc"

	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

// NewRouter -.
func NewRouter(s grpc.ServiceRegistrar, l logger.Interface, t usecase.Shortener) {
	{
		newShortenerRoutes(s, t, l)
	}
}
