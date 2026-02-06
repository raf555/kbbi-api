package dictionaryfx

import (
	"github.com/raf555/kbbi-api/internal/dictionary"
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"github.com/raf555/salome/config/v1"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"dictionary",

	fx.Provide(config.LoadConfigTo[dictionary.Configuration], fx.Private),

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
