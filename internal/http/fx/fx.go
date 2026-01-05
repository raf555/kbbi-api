package httpfx

import (
	"context"
	"fmt"
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
		temp := s.BaseContext
		s.BaseContext = func(nl net.Listener) context.Context {
			if temp != nil {
				ctx = temp(nl)
			}
			// using context.WithoutCancel because the global context will be canceled upon shutdown.
			// we want the server to shutdown gracefully without canceling all serving requests immediately.
			return context.WithoutCancel(ctx)
		}
		return s
	}),
	fx.Invoke(func(ctx context.Context, conf httpsrv.Config, log *slog.Logger, hs *http.Server, shutdowner fx.Shutdowner) error {
		srv, err := server.New(conf.Port)
		if err != nil {
			return fmt.Errorf("server.New: %w", err)
		}

		go func() {
			err := srv.ServeHTTP(hs)
			if err != nil {
				log.ErrorContext(ctx, "http server exited unexpectedly", logger.Error(err))
				_ = shutdowner.Shutdown(fx.ExitCode(1))
			}
		}()

		log.InfoContext(ctx, "listening http...", slog.Int("port", conf.Port))

		return nil
	}),
)
