package httphandler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/http/httperr"
	"github.com/raf555/kbbi-api/internal/http/httpres"
	"github.com/raf555/kbbi-api/internal/logger"
)

// Handler is a simple HTTP request handler which accepts request and response.
type Handler[req, res any] = func(context.Context, *req) (*res, error)

// MakeHandler wraps http controller with request binder and sends response accordingly.
//
// handler can return a http error from [httperr] package to specify status code for the response.
//
// error handling for requestBinder is the same as handler.
//
// Both implementation of handler is encouraged to use errors from package [httperr] to wrap the error to make the best result.
// Otherwise, this handler will always return 5xx error.
func MakeHandler[reqT, resT any](
	handler Handler[reqT, resT],
	requestBinder RequestBinder[reqT],
	opts ...handlerOption,
) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := &ginCtx{gCtx}

		req, err := requestBinder(ctx)
		if err != nil {
			sendResponse(ctx, (*resT)(nil), err, opts...)
			return
		}

		res, err := handler(ctx, req)
		sendResponse(ctx, res, err, opts...)
	}
}

type SimpleHandler[res any] = func(context.Context) (*res, error)

func MakeSimpleHandler[res any](handler SimpleHandler[res]) gin.HandlerFunc {
	return MakeHandler(
		func(ctx context.Context, _ *struct{}) (*res, error) {
			return handler(ctx)
		},
		NoopRequestBinder,
	)
}

func sendResponse[resT any](ctx *ginCtx, res *resT, err error, opts ...handlerOption) {
	options := resolveOptions(ctx, opts...)

	statusCode := httperr.HTTPStatusCode(err)
	statusCodeMessage := http.StatusText(statusCode)

	if err != nil {
		logger.FromContext(ctx).ErrorContext(ctx, "HTTP error occurred", logger.Error(err), slog.Int("code", statusCode))

		innerErrMsg := statusCodeMessage

		errMsg, ok := httperr.HTTPResponseMessage(err)
		if ok {
			innerErrMsg = errMsg
		}

		options.serializer(statusCode, &httpres.Error{Message: innerErrMsg})
		return
	}

	if res != nil {
		resAsAny := any(res)
		// override if status code present in res

		if resCoder, ok := resAsAny.(interface{ GetCode() (int, bool) }); ok {
			if code, ok := resCoder.GetCode(); ok {
				statusCode = code
			}
		}
	}

	options.serializer(statusCode, res)
}

type (
	RedirectResult struct {
		Code int
		Path string
	}

	// RedirectHandler is a simple HTTP request handler which accepts request and redirects into given path.
	RedirectHandler[req any] = func(context.Context, *req) (RedirectResult, error)

	// SimpleRedirectHandler is a simple HTTP request handler which redirects into given path.
	SimpleRedirectHandler = func(context.Context) (RedirectResult, error)
)

func MakeRedirectHandler[reqT any](
	handler RedirectHandler[reqT],
	requestBinder RequestBinder[reqT],
) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := &ginCtx{gCtx}

		req, err := requestBinder(ctx)
		if err != nil {
			sendResponse(ctx, (*struct{})(nil), err)
			return
		}

		result, err := handler(ctx, req)
		if err != nil {
			sendResponse(ctx, (*struct{})(nil), err)
			return
		}

		code := result.Code
		if code == 0 {
			code = http.StatusTemporaryRedirect
		}

		gCtx.Redirect(code, result.Path)
	}
}

func MakeSimpleRedirectHandler(handler SimpleRedirectHandler) gin.HandlerFunc {
	return MakeRedirectHandler(
		func(ctx context.Context, _ *struct{}) (RedirectResult, error) {
			return handler(ctx)
		},
		NoopRequestBinder,
	)
}
