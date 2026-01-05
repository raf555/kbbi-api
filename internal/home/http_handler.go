package home

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/http/httphandler"
)

type HTTPHandler struct {
	statsFetcher AssetStatsFetcher
	baseResponse *HomeResponse
}

func NewHTTPHandler(statsFetcher AssetStatsFetcher) *HTTPHandler {
	return &HTTPHandler{
		statsFetcher: statsFetcher,

		baseResponse: &HomeResponse{
			Message:       "Welcome to the (unofficial) KBBI API",
			Stats:         statsFetcher.Stats(),
			Documentation: "https://kbbi.raf555.dev/swagger/index.html",
			Issues:        "https://github.com/raf555/kbbi-api/issues",
		},
	}
}

func (h *HTTPHandler) MustRegisterRoutes(g *gin.Engine) {
	g.GET("/",
		httphandler.MakeSimpleHandler(h.Home),
	)

	g.GET("/healthzzz",
		httphandler.MakeSimpleHandler(h.Health),
	)
}

func (h *HTTPHandler) Home(ctx context.Context) (*HomeResponse, error) {
	return h.baseResponse, nil
}

func (h *HTTPHandler) Health(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{
		Message: "OK",
	}, nil
}
