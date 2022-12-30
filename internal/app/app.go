package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"service-url-shortener/config"
	grpcroutes "service-url-shortener/internal/entrypoint/grpc"
	"service-url-shortener/internal/entrypoint/http"
	"service-url-shortener/internal/usecase"
	"service-url-shortener/internal/usecase/digitiser"
	"service-url-shortener/internal/usecase/repo"
	grpcI "service-url-shortener/pkg/grpc/interceptors"
	grpcS "service-url-shortener/pkg/grpc/server"
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
	d, err := digitiser.New(
		cfg.Digitiser.Base,
		cfg.Digitiser.Length,
		cfg.Digitiser.MaxRepoInt,
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - digitiser.NewShortener: %w", err))
	}

	shortenerUseCase := usecase.NewShortener(
		repo.New(pg),
		&d,
		fmt.Sprintf("%s:%s/", cfg.URL.Blank, cfg.HTTP.Port),
	)

	// GRPC Server
	grpcSer := grpc.NewServer(grpc.UnaryInterceptor(grpcI.Logs))
	grpcroutes.NewRouter(grpcSer, l, shortenerUseCase)

	grpcServer, err := grpcS.New(grpcSer, cfg.GRPC.Network, cfg.GRPC.Port, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpcServer - server.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, l, shortenerUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

httpContinue:
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))

		// GRPC shutdown
		err = grpcServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
		}
		goto httpContinue
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}
}
