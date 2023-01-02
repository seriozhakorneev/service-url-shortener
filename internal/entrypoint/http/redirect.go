package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	service "service-url-shortener/internal/errors"
	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

type redirectRoutes struct {
	t usecase.Shortener
	l logger.Interface
}

func newRedirectRoute(handler *gin.Engine, t usecase.Shortener, l logger.Interface) {
	r := &redirectRoutes{t, l}
	handler.GET("/*any", r.get)
}

// get redirects to original url hidden under /<url string identifier>
func (r *redirectRoutes) get(c *gin.Context) {
	URL, err := r.t.Lengthen(c.Request.Context(), c.Request.RequestURI)
	if err != nil {
		switch {
		// caching error - just prints log with error, no error response
		case errors.Is(err, service.ErrCaching):
			r.l.Warn(err.Error(), "http - v1 - get")
		case errors.Is(err, service.ErrImpossibleShortURL):
			errorResponse(c, http.StatusBadRequest, service.ErrImpossibleShortURL.Error())
		case errors.Is(err, service.ErrNotFoundURL):
			errorResponse(c, http.StatusNotFound, service.ErrNotFoundURL.Error())
		default:
			r.l.Error(err, "http - v1 - get")
			errorResponse(c, http.StatusInternalServerError, "shortener service problems")
		}
	}

	c.Redirect(http.StatusSeeOther, URL)
}
