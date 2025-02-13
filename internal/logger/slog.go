package logger

import (
	"log/slog"
	"os"

	slogformatter "github.com/samber/slog-formatter"
)

func NewLogger() *slog.Logger {
	timeDurationFormatter := slogformatter.FormatByKind(slog.KindDuration, func(v slog.Value) slog.Value {
		return slog.StringValue(v.Duration().String())
	})
	return slog.New(
		slogformatter.NewFormatterHandler(timeDurationFormatter)(
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)
}
