package http

import (
	"github.com/gin-gonic/gin"

	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

//  TODO DOCS :

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Shortener) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	newRedirectRoute(handler, t, l)
}
