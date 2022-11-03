package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"service-url-shortener/internal/entity"
	internal "service-url-shortener/internal/errors"
	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

type shortenerRoutes struct {
	t usecase.Shortener
	l logger.Interface
}

func newShortenerRoutes(handler *gin.RouterGroup, t usecase.Shortener, l logger.Interface) {
	r := &shortenerRoutes{t, l}

	h := handler.Group("/shortener")
	{
		h.POST("/create", r.create)
		h.GET("/get", r.get)
	}
}

// TODO docs
func (r *shortenerRoutes) create(c *gin.Context) {
	var request entity.ShortenerData
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	shortURL, err := r.t.Shorten(c.Request.Context(), request.URL)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.JSON(http.StatusOK, entity.ShortenerData{URL: shortURL})
}

// TODO docs
func (r *shortenerRoutes) get(c *gin.Context) {
	var request entity.ShortenerData
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	URL, err := r.t.Lengthen(c.Request.Context(), request.URL)
	if err != nil {
		if errors.Is(err, internal.ErrNotFoundURL) {
			errorResponse(c, http.StatusNotFound, err.Error())

			return
		}
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.JSON(http.StatusOK, entity.ShortenerData{URL: URL})
}
