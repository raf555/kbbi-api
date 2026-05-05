package logger

import (
	"log/slog"
	"os"

	"github.com/raf555/salome/melt/log"
	slogformatter "github.com/samber/slog-formatter"
)

func New(cfg Config) *slog.Logger {
	var handler slog.Handler

	handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Level})
	handler = log.WithOtelHandler(handler)
	handler = log.WithFormatter(handler, slogformatter.FormatByKind(slog.KindDuration, func(v slog.Value) slog.Value {
		return slog.StringValue(v.Duration().String())
	}))

	return log.New(handler)
}
