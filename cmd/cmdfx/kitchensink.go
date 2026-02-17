package cmdfx

import (
	"github.com/raf555/kbbi-api/internal/config/configfx"
	"github.com/raf555/kbbi-api/internal/context/contextfx"
	"github.com/raf555/kbbi-api/internal/logger/loggerfx"
	"go.uber.org/fx"
)

// KitchenSink holds a common dependency for the application.
var KitchenSink = fx.Options(
	configfx.Module,
	loggerfx.Provider,
	OtelModule,
	contextfx.ContextDecoratorOption,
)
