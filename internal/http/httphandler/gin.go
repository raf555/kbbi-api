package httphandler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginCtx struct {
	*gin.Context
}

func (g *ginCtx) Request() *http.Request {
	return g.Context.Request
}

type GinMinimalContext interface {
	context.Context

	Request() *http.Request

	ShouldBindJSON(obj any) error
	ShouldBindQuery(obj any) error
	ShouldBindUri(obj any) error
	ShouldBindHeader(obj any) error
}
