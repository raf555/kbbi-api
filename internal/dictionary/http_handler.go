package dictionary

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/http/httperr"
	"github.com/raf555/kbbi-api/internal/http/httphandler"
)

type HTTPHandler struct {
	dict DictionaryRepo
}

func NewHTTPHandler(dict DictionaryRepo) *HTTPHandler {
	return &HTTPHandler{
		dict: dict,
	}
}

func (h *HTTPHandler) MustRegisterRoutes(g *gin.Engine) {
	entryGroupV1 := g.Group("/api/v1/entry")

	entryGroupV1.GET("/_random",
		httphandler.MakeSimpleRedirectHandler(
			h.Random,
		),
	)

	entryGroupV1.GET("/:entry",
		h.redirectToLowercase,
		httphandler.MakeHandler(
			h.Entry,
			httphandler.DefaultRequestBinder,
			httphandler.WithPureJSONSerializer(),
		),
	)

	entryGroupV1.GET("/_wotd",
		httphandler.MakeSimpleRedirectHandler(
			h.WOTD,
		),
	)
}

func (*HTTPHandler) redirectToLowercase(ctx *gin.Context) {
	param := ctx.Param("entry")

	if lowered := strings.ToLower(param); lowered != param {
		// TODO: fix this hack maybe?
		//
		// for some reason, the ctx.Request.URL.Path is having the unescaped path,
		// which caused the underlying http.Redirect to not form the new path properly if the entry has a slash.
		// 	e.g. /api/v1/entry/termometer%20maks%2fmin%20Fahrenheit -> /api/v1/entry/termometer maks/min Fahrenheit
		// the underlying http.Redirect will wrongly redirect the path which will become like this.
		//	e.g. new path (lowered): termometer maks/min fahrenheit
		//	i.e. /api/v1/entry/termometer maks/min Fahrenheit -> /api/v1/entry/termometer maks/termometer maks/min fahrenheit
		//
		// this is a hack to make sure it properly redirects to the correct lemma path.
		if ctx.Request.URL.RawPath != "" {
			ctx.Request.URL.Path = ctx.Request.URL.RawPath
		}

		path := url.PathEscape(lowered)
		if query := ctx.Request.URL.RawQuery; query != "" {
			path += fmt.Sprintf("?%s", ctx.Request.URL.RawQuery)
		}

		ctx.Redirect(http.StatusMovedPermanently, path)
		ctx.Abort()
	}
}

// Entry godoc
// @Summary      Show Lemma Information
// @Description  Show the information of provided lemma
// @Tags         entry
// @Accept       json
// @Produce      json
// @Param        entry    path      string  true  "Lemma. E.g. apel, aku (2), etc."
// @Param        entryNo  query     int	  	false "Lemma's entry number (optional). Start from 1. Will be skipped if there's entry number in the lemma." minimum(1)
// @Success      200   	  {object}  kbbi.Lemma
// @Failure      400      {object}  httpres.Error
// @Failure      404      {object}  httpres.Error
// @Failure      414      {object}  httpres.Error
// @Failure      500      {object}  httpres.Error
// @Router       /api/v1/entry/{entry} [get]
func (h *HTTPHandler) Entry(ctx context.Context, req *EntryRequest) (*EntryResponse, error) {
	req.transform()

	data, err := h.dict.Lemma(req.Lemma, req.EntryNo)
	if err != nil {
		wrappedErr := fmt.Errorf("h.dict.Lemma: %w", err)
		switch {
		case errors.Is(err, ErrUnexpectedEmptyLemma):
			return nil, httperr.Wrap(wrappedErr, http.StatusBadRequest, "empty lemma")
		case errors.Is(err, ErrUnexpectedEntryNumber):
			return nil, httperr.Wrapf(wrappedErr, http.StatusBadRequest, "invalid entry number: %d", req.EntryNo)
		case errors.Is(err, ErrLemmaNotFound):
			return nil, httperr.Wrap(wrappedErr, http.StatusNotFound, "lemma not found")
		case errors.Is(err, ErrEntryNotFound):
			return nil, httperr.Wrap(wrappedErr, http.StatusNotFound, "lemma's entry not found")
		case errors.Is(err, ErrLemmaTooLong):
			return nil, httperr.Wrap(wrappedErr, http.StatusRequestURITooLong, "lemma is too long")
		default:
			return nil, wrappedErr
		}
	}

	return &EntryResponse{data}, nil
}

// Entry godoc
// @Summary      Get Random Lemma
// @Description  Redirect to the random lemma
// @Tags         entry
// @Success      200      {object}  kbbi.Lemma
// @Success      302      {object}  kbbi.Lemma
// @Failure      500      {object}  httpres.Error
// @Router       /api/v1/entry/_random [get]
func (h *HTTPHandler) Random(ctx context.Context) (httphandler.RedirectResult, error) {
	lemma := h.dict.RandomLemma()

	return httphandler.RedirectResult{
		Code: http.StatusFound,
		Path: url.PathEscape(lemma.Lemma),
	}, nil
}

// Entry godoc
// @Summary      Get Lemma of The Day
// @Description  Redirect to the lemma of the day
// @Tags         entry
// @Success      200      {object}  kbbi.Lemma
// @Success      302      {object}  kbbi.Lemma
// @Failure      500      {object}  httpres.Error
// @Router       /api/v1/entry/_wotd [get]
func (h *HTTPHandler) WOTD(ctx context.Context) (httphandler.RedirectResult, error) {
	wotd, err := h.dict.LemmaOfTheDay()
	if err != nil {
		return httphandler.RedirectResult{}, fmt.Errorf("h.dict.LemmaOfTheDay: %w", err)
	}

	return httphandler.RedirectResult{
		Code: http.StatusFound,
		Path: url.PathEscape(wotd.Lemma),
	}, nil
}
