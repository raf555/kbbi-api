package swaggerfx

import (
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"github.com/raf555/kbbi-api/internal/swagger"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"swagger",

	httpfx.HandlerProvider(
		swagger.NewHTTPHandler,
	),
)
