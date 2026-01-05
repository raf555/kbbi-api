package dictionaryfx

import (
	"github.com/raf555/kbbi-api/internal/config"
	"github.com/raf555/kbbi-api/internal/dictionary"
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"dictionary",

	fx.Provide(config.EnvConfigProvider[dictionary.Configuration], fx.Private),

	fx.Provide(
		fx.Annotate(
			dictionary.NewWOTD,
			fx.As(new(dictionary.WOTDRepo)),
		),
		fx.Private,
	),

	fx.Provide(
		fx.Annotate(
			dictionary.NewDictionary,
			fx.As(new(dictionary.DictionaryRepo)),
			fx.As(fx.Self()),
		),
	),

	httpfx.HandlerProvider(
		dictionary.NewHTTPHandler,
	),
)
