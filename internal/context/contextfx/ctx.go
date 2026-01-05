package contextfx

import (
	"context"

	contextx "github.com/raf555/kbbi-api/internal/context"
	"go.uber.org/fx"
)

type contextDecoratorParams struct {
	fx.In

	Decorators []contextx.ContextDecoratorFn `group:"ctx.decorators"`
}

var ContextDecoratorOption = fx.Options(
	fx.Decorate(func(
		ctx context.Context,
		params contextDecoratorParams,
	) context.Context {
		for _, decoratorFn := range params.Decorators {
			ctx = decoratorFn(ctx)
		}
		return ctx
	}),
)

func Provider(constructor any, opts ...fx.Annotation) fx.Option {
	opt := []fx.Annotation{
		fx.ResultTags(`group:"ctx.decorators"`),
	}
	opt = append(opt, opts...)

	return fx.Provide(
		fx.Annotate(
			constructor,
			opt...,
		),
	)
}
