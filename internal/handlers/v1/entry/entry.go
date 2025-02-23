package entry

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/dto/response"
	"github.com/raf555/kbbi-api/internal/repositories/dict"
)

// Entry godoc
// @Summary      Show Lemma Information
// @Description  Show the information of provided lemma
// @Tags         entry
// @Accept       json
// @Produce      json
// @Param        entry    path      string  true  "Lemma"
// @Param        entryNo  query     int	  	false "Lemma's entry number (optional). Start from 1." minimum(1)
// @Success      200   	  {object}  kbbi.Lemma
// @Failure      400      {object}  response.Error
// @Failure      404      {object}  response.Error
// @Failure      414      {object}  response.Error
// @Failure      500      {object}  response.Error
// @Router       /api/v1/entry/{entry} [get]
func (h *Handler) Entry(ctx *gin.Context) {
	entryNo, ok := h.getEntryNo(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.Error{Message: errUnexpectedEntryNo})
		return
	}

	lemma := ctx.Param("entry")

	data, err := h.dict.Lemma(lemma, entryNo)
	if err != nil {
		switch {
		case errors.Is(err, dict.ErrUnexpectedEmptyLemma) || errors.Is(err, dict.ErrUnexpectedEntryNumber):
			ctx.JSON(http.StatusBadRequest, response.Error{Message: err})
		case errors.Is(err, dict.ErrLemmaNotFound) || errors.Is(err, dict.ErrEntryNotFound):
			ctx.JSON(http.StatusNotFound, response.Error{Message: err})
		case errors.Is(err, dict.ErrLemmaTooLong):
			ctx.JSON(http.StatusRequestURITooLong, response.Error{Message: err})
		default:
			h.logger.ErrorContext(ctx, "Unexpected lemma retrieval failure", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, response.Error{Message: response.ErrInternalServerError})
		}

		return
	}

	ctx.PureJSON(http.StatusOK, data)
}

func (h *Handler) getEntryNo(ctx *gin.Context) (*int, bool) {
	entryQuery := ctx.Query("entryNo")
	if entryQuery == "" {
		return nil, true
	}

	entryNoInt64, err := strconv.ParseInt(entryQuery, 10, 32)
	if err != nil {
		return nil, false
	}

	entryNo := int(entryNoInt64)
	return &entryNo, true
}
