package httpsrv

import "time"

type Config struct {
	Port                   int           `env:"PORT, default=8888"`
	HTTPServerReadTimeout  time.Duration `env:"HTTP_SERVER_READ_TIMEOUT, default=1s"`
	HTTPServerWriteTimeout time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT, default=1s"`
}
