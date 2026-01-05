package httphandler

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/raf555/kbbi-api/internal/http/httperr"
)

type RequestBinder[req any] = func(GinMinimalContext) (*req, error)

// DefaultRequestBinder will try to bind from header, query, and uri into reqT.
func DefaultRequestBinder[reqT any](ctx GinMinimalContext) (*reqT, error) {
	var req reqT
	if err := commonBinder(ctx, &req); err != nil {
		return nil, err
	}

	if err := validateStruct(ctx, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func commonBinder[reqT any](ctx GinMinimalContext, req *reqT) error {
	if err := ctx.ShouldBindHeader(req); err != nil {
		return httperr.Wrap(err, http.StatusBadRequest, "failed to bind header from request")
	}

	if err := ctx.ShouldBindQuery(req); err != nil {
		return httperr.Wrap(err, http.StatusBadRequest, "failed to bind query from request")
	}

	if err := ctx.ShouldBindUri(req); err != nil {
		return httperr.Wrap(err, http.StatusBadRequest, "failed to bind uri params from request")
	}

	return nil
}

func validateStruct[reqT any](ctx context.Context, req *reqT) error {
	if err := validate.StructCtx(ctx, req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return httperr.Wrap(err, http.StatusBadRequest, "failed to validate request")
		}
		return httperr.Wrap(err, http.StatusBadRequest, errs[0].Translate(validateTranslator))
	}
	return nil
}

func NoopRequestBinder[reqT any](ctx GinMinimalContext) (*reqT, error) {
	var zero reqT
	return &zero, nil
}
