package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"service-url-shortener/internal/usecase"
	"service-url-shortener/pkg/logger"
)

type EchoRoutes struct {
	t usecase.Echo
	l logger.Interface
}

func newEchoRoutes(handler *gin.RouterGroup, t usecase.Echo, l logger.Interface) {
	r := &EchoRoutes{t, l}

	h := handler.Group("/echo")
	{
		h.POST("/reflect", r.reflect)
	}
}

// @Summary     Reflect
// @Description Reflects the received data, overwrites the data if overwrite rules are active
// @ID          reflect
// @Tags  	    echo
// @Accept      json
// @Produce     json
// @Param       request body entity.JSONObject true "Set up any json object"
// @Success     200 {object} entity.JSONObject
// @Failure     400 {object} response
// @Router      /echo/reflect [post]
func (r *EchoRoutes) reflect(c *gin.Context) {
	var data map[string]any

	err := c.BindJSON(&data)
	if err != nil {
		r.l.Error(err, "http - v1 - reflect")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	r.t.Rewrite(data)

	c.JSON(http.StatusOK, data)
}
