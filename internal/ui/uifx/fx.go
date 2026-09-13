package uifx

import (
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"github.com/raf555/kbbi-api/internal/ui"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"ui",
	httpfx.HandlerProvider(ui.New),
)
