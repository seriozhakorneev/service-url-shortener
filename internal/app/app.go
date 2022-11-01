package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"service-url-shortener/config"
	v1 "service-url-shortener/internal/entrypoint/http/v1"
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
		l.Fatal(fmt.Errorf("app - Run - digitiser.New: %w", err))
	}

	shortenerUseCase := usecase.New(
		repo.New(pg),
		&d,
		fmt.Sprintf("%s:%s/", cfg.URL.Blank, cfg.HTTP.Port),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, shortenerUseCase)

	log.Printf("swagger docs on  http://localhost:%v/swagger/index.html", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
