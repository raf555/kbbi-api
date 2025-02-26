package entry

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Entry godoc
// @Summary      Get Random Lemma
// @Description  Redirect to the random lemma
// @Tags         entry
// @Success      302      {object}  kbbi.Lemma
// @Failure      500      {object}  response.Error
// @Router       /api/v1/entry/_random [get]
func (h *Handler) Random(ctx *gin.Context) {
	lemma := h.dict.RandomLemma()
	ctx.Redirect(http.StatusFound, url.PathEscape(lemma.Lemma))
}
