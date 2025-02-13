package entry

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/repositories/dict"
)

type (
	Handler struct {
		logger *slog.Logger
		dict   *dict.Dictionary
	}
)

func New(logger *slog.Logger, dict *dict.Dictionary) *Handler {
	return &Handler{logger, dict}
}

func (h *Handler) RedirectToLowercase(ctx *gin.Context) {
	param := ctx.Param("entry")

	if lowered := strings.ToLower(param); lowered != param {
		ctx.Redirect(http.StatusMovedPermanently, lowered)
		ctx.Abort()
	}
}
