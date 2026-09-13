package ui

import (
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct{}

func New() (*HTTPHandler, error) {
	return &HTTPHandler{}, nil
}

func (h *HTTPHandler) MustRegisterRoutes(g *gin.Engine) {
	g.Static("/ui", "assets/view")
}
