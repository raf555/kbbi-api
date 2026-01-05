package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPHandler struct{}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) MustRegisterRoutes(g *gin.Engine) {
	SwaggerInfo.Title = "KBBI API"
	SwaggerInfo.Description = "Probably the most complete KBBI API you will ever find."

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
