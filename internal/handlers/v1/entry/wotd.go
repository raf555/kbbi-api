package entry

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/dto/response"
)

// Entry godoc
// @Summary      Get Lemma of The Day
// @Description  Redirect to the lemma of the day
// @Tags         entry
// @Success      200      {object}  kbbi.Lemma
// @Success      302      {object}  kbbi.Lemma
// @Failure      500      {object}  response.Error
// @Router       /api/v1/entry/_wotd [get]
func (h *Handler) WOTD(ctx *gin.Context) {
	wotd, err := h.dict.LemmaOfTheDay()
	if err != nil {
		h.logger.ErrorContext(ctx, "Unexpected WOTD failure", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, url.PathEscape(wotd.Lemma))
}
