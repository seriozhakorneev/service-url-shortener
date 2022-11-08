package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"service-url-shortener/config"
	rgrpc "service-url-shortener/internal/entrypoint/grpc"
	"service-url-shortener/internal/entrypoint/http"
	"service-url-shortener/internal/usecase"
	"service-url-shortener/internal/usecase/digitiser"
	"service-url-shortener/internal/usecase/repo"
	"service-url-shortener/pkg/httpserver"
	"service-url-shortener/pkg/logger"
	"service-url-shortener/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use Case
	d, err := digitiser.New(cfg.Digitiser.Base, cfg.Digitiser.Length)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - digitiser.NewShortener: %w", err))
	}

	shortenerUseCase := usecase.NewShortener(
		repo.New(pg),
		&d,
		fmt.Sprintf("%s:%s/", cfg.URL.Blank, cfg.HTTP.Port),
	)

	//TODO !!!!

	lis, err := net.Listen(cfg.GRPC.Network, fmt.Sprintf("localhost:%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// GRPC Server
	grpcServer := grpc.NewServer()
	rgrpc.NewRouter(grpcServer, shortenerUseCase, l)

	//TODO IT BLOCKS ALL NEXT CODE
	err = grpcServer.Serve(lis)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpcServer - grpcServer.Serve: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, l, shortenerUseCase)

	log.Printf("swagger docs on  http://localhost:%v/swagger/index.html", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// TODO GRPC SERVER HERE
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	l.Debug("Server shutdown")

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
