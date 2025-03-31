package entry

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/dto/response"
	"github.com/raf555/kbbi-api/internal/repositories/dict"
	"github.com/raf555/kbbi-api/internal/util"
)

// Entry godoc
// @Summary      Show Lemma Information
// @Description  Show the information of provided lemma
// @Tags         entry
// @Accept       json
// @Produce      json
// @Param        entry    path      string  true  "Lemma. E.g. apel, aku (2), etc."
// @Param        entryNo  query     int	  	false "Lemma's entry number (optional). Start from 1. Will be skipped if there's entry number in the lemma." minimum(1)
// @Success      200   	  {object}  kbbi.Lemma
// @Failure      400      {object}  response.Error
// @Failure      404      {object}  response.Error
// @Failure      414      {object}  response.Error
// @Failure      500      {object}  response.Error
// @Router       /api/v1/entry/{entry} [get]
func (h *Handler) Entry(ctx *gin.Context) {
	lemma := ctx.Param("entry")

	entryNo, ok := h.getEntryNo(ctx, &lemma)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorOf(errUnexpectedEntryNo))
		return
	}

	data, err := h.dict.Lemma(lemma, entryNo)
	if err != nil {
		switch {
		case errors.Is(err, dict.ErrUnexpectedEmptyLemma) || errors.Is(err, dict.ErrUnexpectedEntryNumber):
			ctx.JSON(http.StatusBadRequest, response.ErrorOf(err))
		case errors.Is(err, dict.ErrLemmaNotFound) || errors.Is(err, dict.ErrEntryNotFound):
			ctx.JSON(http.StatusNotFound, response.ErrorOf(err))
		case errors.Is(err, dict.ErrLemmaTooLong):
			ctx.JSON(http.StatusRequestURITooLong, response.ErrorOf(err))
		default:
			h.logger.ErrorContext(ctx, "Unexpected lemma retrieval failure", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError)
		}

		return
	}

	ctx.PureJSON(http.StatusOK, data)
}

// getEntryNo will initially look for entry number in the lemma itself.
// if it's present, it'll overwrite the lemma without the number and return the number.
// entryNo in the query param will be skipped.
//
// otherwise, getEntryNo will return false if the entryNo is present in the query param and is an invalid number.
// on other case, it'll always return true (indicates no entry number to query).
func (h *Handler) getEntryNo(ctx *gin.Context, lemma *string) (*int, bool) {
	// override if there's any entry number in the lemma
	if newLemma, entryNo, ok := util.FindEntryNoFromLemma(*lemma); ok {
		*lemma = newLemma
		return &entryNo, true
	}

	entryQuery := ctx.Query("entryNo")
	if entryQuery == "" {
		return nil, true
	}

	entryNo, err := strconv.Atoi(entryQuery)
	if err != nil {
		return nil, false
	}

	return &entryNo, true
}
