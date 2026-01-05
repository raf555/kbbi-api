package homefx

import (
	"github.com/raf555/kbbi-api/internal/dictionary"
	"github.com/raf555/kbbi-api/internal/home"
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"home",

	fx.Provide(
		func(dict *dictionary.Dictionary) home.AssetStatsFetcher {
			return dict
		},
	),

	httpfx.HandlerProvider(
		home.NewHTTPHandler,
	),
)
