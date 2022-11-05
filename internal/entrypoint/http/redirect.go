package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	internal "service-url-shortener/internal/errors"
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

func (r *redirectRoutes) get(c *gin.Context) {
	URL, err := r.t.Lengthen(c.Request.Context(), c.Request.RequestURI)
	if err != nil {
		if errors.Is(err, internal.ErrNotFoundURL) {
			errorResponse(c, http.StatusNotFound, err.Error())

			return
		}
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.Redirect(http.StatusSeeOther, URL)
}
