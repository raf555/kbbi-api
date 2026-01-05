package cmdfx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/raf555/kbbi-api/internal/config/load" // load dotenv
	"github.com/raf555/kbbi-api/internal/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var errUnexpectedExit = errors.New("cmdfx: unexpected exit")

var (
	startTimeout = 30 * time.Second
	stopTimeout  = 30 * time.Second
)

// Run runs application provided fx options. It provides common things such as global context, dependency (from [KitchenSink])
// and graceful shutdown.
func Run(ctx context.Context, options ...fx.Option) error {
	app, runErr := runContainer(ctx, true, options...)
	if runErr != nil && !errors.Is(runErr, errUnexpectedExit) {
		return fmt.Errorf("run application: %w", runErr)
	}

	if err := stopContainer(ctx, app); err != nil {
		logger.FromContext(ctx).WarnContext(ctx, "failed to stop application", logger.Error(err))
	}

	// must be unexpected exit
	if runErr != nil {
		return runErr
	}

	logger.FromContext(ctx).InfoContext(ctx, "shutdown successful")

	return nil
}

func runContainer(ctx context.Context, waitSig bool, options ...fx.Option) (*fx.App, error) {
	containerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	options = append([]fx.Option{
		KitchenSink,
		fx.WithLogger(func(_ *slog.Logger) fxevent.Logger {
			// return &fxevent.SlogLogger{
			// 	Logger: l,
			// }
			return fxevent.NopLogger
		}),
		fx.Supply(
			fx.Annotate(
				containerCtx,
				fx.As(new(context.Context)),
				fx.OnStop(cancel),
			),
		),
		fx.StartTimeout(startTimeout),
		fx.StopTimeout(stopTimeout),
	}, options...)

	app := fx.New(options...)
	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("building app: %w", err)
	}

	if err := app.Start(ctx); err != nil {
		return nil, fmt.Errorf("starting app: %w", err)
	}

	if waitSig {
		sig := <-app.Wait()
		logger.FromContext(ctx).InfoContext(ctx, "application is shutting down from signal", slog.Any("sig", sig))
		if sig.ExitCode == 1 {
			return app, errUnexpectedExit
		}
	} else {
		logger.FromContext(ctx).InfoContext(ctx, "application is shutting down")
	}

	return app, nil
}

func stopContainer(ctx context.Context, app *fx.App) error {
	ctx, done := context.WithTimeout(context.WithoutCancel(ctx), stopTimeout)
	defer done()

	if err := app.Stop(ctx); err != nil {
		return fmt.Errorf("stopping app: %w", err)
	}

	return nil
}
