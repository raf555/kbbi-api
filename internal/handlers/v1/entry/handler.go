package entry

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
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

var (
	errUnexpectedEntryNo = errors.New("unexpected `entryNo` query")
)

func New(logger *slog.Logger, dict *dict.Dictionary) *Handler {
	return &Handler{logger, dict}
}

func (h *Handler) RedirectToLowercase(ctx *gin.Context) {
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
		ctx.Request.URL.Path = ctx.Request.URL.RawPath

		ctx.Redirect(http.StatusMovedPermanently, url.PathEscape(lowered))
		ctx.Abort()
	}
}
