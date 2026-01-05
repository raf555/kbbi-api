package loggerfx

import (
	"context"
	"log/slog"

	contexts "github.com/raf555/kbbi-api/internal/context"
	"github.com/raf555/kbbi-api/internal/context/contextfx"
	"github.com/raf555/kbbi-api/internal/logger"
	"go.uber.org/fx"
)

var Provider = fx.Module(
	"logger",
	fx.Provide(logger.New),
	contextfx.Provider(func(log *slog.Logger) contexts.ContextDecoratorFn {
		return func(ctx context.Context) context.Context {
			return logger.WithContext(ctx, log)
		}
	}),
	fx.Invoke(func(l *slog.Logger) {
		slog.SetDefault(l)
	}),
)
