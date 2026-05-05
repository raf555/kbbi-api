package loggerfx

import (
	"context"
	"log/slog"

	contexts "github.com/raf555/kbbi-api/internal/context"
	"github.com/raf555/kbbi-api/internal/context/contextfx"
	"github.com/raf555/kbbi-api/internal/logger"
	salomeconfig "github.com/raf555/salome/config/v1"
	"go.uber.org/fx"
)

var Provider = fx.Module(
	"logger",
	fx.Provide(salomeconfig.LoadConfigTo[logger.Config], fx.Private),
	fx.Provide(logger.New),
	contextfx.Provider(func(log *slog.Logger) contexts.ContextDecoratorFn {
		return func(ctx context.Context) context.Context {
			return logger.WithContext(ctx, log)
		}
	}),
)
