package cmdfx

import (
	"context"

	"github.com/raf555/kbbi-api/internal/config"
	contextx "github.com/raf555/kbbi-api/internal/context"
	"github.com/raf555/kbbi-api/internal/context/contextfx"
	"github.com/raf555/salome/melt/metric"
	"github.com/raf555/salome/melt/otel"
	"github.com/raf555/salome/melt/trace"
	otelmetric "go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var OtelModule = fx.Module(
	"otel",

	fx.Provide(
		fx.Annotate(
			func(ctx context.Context, cfg config.ServerConfig) (otel.OpenTelemetry, error) {
				return otel.NewOrNoop(ctx, cfg.ServiceName)
			},
			fx.OnStop(func(ctx context.Context, ot otel.OpenTelemetry) error {
				return ot.Shutdown(ctx)
			}),
		),
		func(ot otel.OpenTelemetry) (otelmetric.MeterProvider, oteltrace.TracerProvider) {
			return ot.MeterProvider(), ot.TracerProvider()
		},
	),

	fx.Provide(
		func(cfg config.ServerConfig, mp otelmetric.MeterProvider) (*metric.RecorderProvider, error) {
			return metric.New(cfg.ServiceName, mp)
		},
		func(rp *metric.RecorderProvider) metric.Recorder {
			return rp.DefaultRecorder()
		},

		func(cfg config.ServerConfig, tp oteltrace.TracerProvider) *trace.TracerProvider {
			return trace.New(cfg.ServiceName, tp)
		},
		func(tp *trace.TracerProvider) trace.Tracer {
			return tp.Tracer()
		},
	),

	contextfx.Provider(func(tracer trace.Tracer) contextx.ContextDecoratorFn {
		return func(ctx context.Context) context.Context {
			return trace.WithContext(ctx, tracer)
		}
	}),

	contextfx.Provider(func(mr metric.Recorder) contextx.ContextDecoratorFn {
		return func(ctx context.Context) context.Context {
			return metric.WithContext(ctx, mr)
		}
	}),
)
