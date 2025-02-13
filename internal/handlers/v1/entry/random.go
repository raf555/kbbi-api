package entry

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/dto/response"
)

// Entry godoc
// @Summary      Get Random Lemma
// @Description  Redirect to the random lemma
// @Tags         entry
// @Success      302      {object}  kbbi.Lemma
// @Failure      500      {object}  response.Error
// @Router       /api/v1/entry/_random [get]
func (h *Handler) Random(ctx *gin.Context) {
	lemma, err := h.dict.RandomLemma()
	if err != nil {
		h.logger.ErrorContext(ctx, "Unexpected Random Lemma failure", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, response.Error{Message: response.ErrInternalServerError})
		return
	}

	ctx.Redirect(http.StatusFound, lemma.Lemma)
}
