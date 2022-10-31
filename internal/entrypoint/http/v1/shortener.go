package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

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
	var request shortenerData
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	shortURL, err := r.t.Shorten(c.Request.Context(), request.URI)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.JSON(http.StatusOK, shortenerData{shortURL})
}

// TODO docs
func (r *shortenerRoutes) get(c *gin.Context) {
	var request shortenerData
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	URI, err := r.t.Lengthen(c.Request.Context(), request.URI)
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.JSON(http.StatusOK, shortenerData{URI})
}

type shortenerData struct {
	URI string `json:"URI"`
}

// UnmarshalJSON переопределенный json.unmarshal для shortenerData
func (d *shortenerData) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}

	tmp := struct {
		URI string `json:"URI"`
	}{}

	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}

	if tmp.URI == "" {
		return fmt.Errorf("required field 'URI' is empty")
	}

	*d = shortenerData{URI: tmp.URI}
	return nil
}
