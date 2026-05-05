package logger

import "log/slog"

type Config struct {
	Level slog.Level `env:"LOG_LEVEL,default=INFO"`
}
