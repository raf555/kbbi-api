package httpfx

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/raf555/kbbi-api/internal/config"
	"github.com/raf555/kbbi-api/internal/http/httpsrv"
	"github.com/raf555/kbbi-api/internal/logger"
	"github.com/raf555/kbbi-api/internal/server"
	"go.uber.org/fx"
)

var ServerInvoker = fx.Module(
	"http_server",
	fx.Provide(config.EnvConfigProvider[httpsrv.Config], fx.Private),
	fx.Provide(
		fx.Annotate(
			httpsrv.NewServer,
			fx.OnStop(func(ctx context.Context, s *http.Server) error {
				return s.Shutdown(ctx)
			}),
		),
		fx.Private,
	),
	fx.Decorate(func(ctx context.Context, s *http.Server) *http.Server {
		s.BaseContext = func(_ net.Listener) context.Context {
			return context.WithoutCancel(ctx)
		}
		return s
	}),
	fx.Provide(
		func(c httpsrv.Config) (*server.Server, error) {
			srv, err := server.New(c.Port)
			return srv, err
		},
		fx.Private,
	),
	fx.Invoke(func(ctx context.Context, conf httpsrv.Config, log *slog.Logger, s *server.Server, hs *http.Server, shutdowner fx.Shutdowner) {
		go func() {
			err := s.ServeHTTP(hs)
			if err != nil {
				log.ErrorContext(ctx, "http server exited unexpectedly", logger.Error(err))
				shutdowner.Shutdown(fx.ExitCode(1))
			}
		}()

		log.InfoContext(ctx, "listening http...", slog.Int("port", conf.Port))
	}),
)

func HandlerProvider(constructor any, opts ...fx.Annotation) fx.Option {
	opt := []fx.Annotation{
		fx.As(new(httpsrv.RouterRegistrar)),
		fx.ResultTags(`group:"http.controllers"`),
	}
	opt = append(opt, opts...)

	return fx.Provide(
		fx.Annotate(
			constructor,
			opt...,
		),
	)
}
