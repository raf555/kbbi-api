package logger

import "log/slog"

// Error is an alias for `slog.String("error", err.Error())`
func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
