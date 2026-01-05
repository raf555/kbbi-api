package httpfx

import (
	"github.com/raf555/kbbi-api/internal/http/httpsrv"
	"go.uber.org/fx"
)

func HandlerProvider(constructor any, opts ...fx.Annotation) fx.Option {
	opt := []fx.Annotation{
		fx.As(new(httpsrv.RouterRegistrar)),
		fx.ResultTags(`group:"http.controllers"`),
	}
	opt = append(opt, opts...)

	return fx.Provide(
		fx.Annotate(
			constructor,
			opt...,
		),
	)
}
