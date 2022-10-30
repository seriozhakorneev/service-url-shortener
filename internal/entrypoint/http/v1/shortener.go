package v1

import (
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

//TODO docs

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *shortenerRoutes) create(c *gin.Context) {

	var request long
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

	c.JSON(http.StatusOK, short{shortURL})
}

//TODO docs

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *shortenerRoutes) get(c *gin.Context) {
	var request short
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	URI, err := r.t.Lengthen(c.Request.Context(), request.ShortURI)
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusInternalServerError, "shortener service problems")

		return
	}

	c.JSON(http.StatusOK, long{URI})
}

type long struct {
	URI string `json:"URI"`
}

type short struct {
	ShortURI string `json:"short_URI"`
}
