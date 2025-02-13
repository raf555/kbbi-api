package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/models/assets"
	"github.com/raf555/kbbi-api/internal/repositories/dict"
)

type (
	Handler struct {
		dictRepo *dict.Dictionary

		baseResponse *HomeResponse
	}

	HomeResponse struct {
		Message       string       `json:"message"`
		Stats         assets.Stats `json:"stats"`
		Documentation string       `json:"documentation"`
		Issues        string       `json:"issues"`
	}
)

func New(dictRepo *dict.Dictionary) *Handler {
	return &Handler{
		dictRepo: dictRepo,

		baseResponse: &HomeResponse{
			Message:       "Welcome to the (unofficial) KBBI API",
			Stats:         dictRepo.Stats(),
			Documentation: "https://kbbi.raf555.dev/swagger/index.html",
			Issues:        "https://github.com/raf555/kbbi-api/issues",
		},
	}
}

func (c *Handler) Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.baseResponse)
}

func (c *Handler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
