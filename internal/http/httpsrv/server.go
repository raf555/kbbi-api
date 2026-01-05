package httpsrv

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/http/httpres"
	sloggin "github.com/samber/slog-gin"
	"go.uber.org/fx"
)

type RouterRegistrar interface {
	MustRegisterRoutes(g *gin.Engine)
}

type ServerParam struct {
	fx.In

	Conf      Config
	Slog      *slog.Logger
	Routerers []RouterRegistrar `group:"http.controllers"`
}

func NewServer(param ServerParam) *http.Server {
	g := gin.New()
	g.ContextWithFallback = true
	g.UseRawPath = true

	logger := param.Slog.With(slog.String("label", "http_server"))

	registerMiddlewares(g, logger)

	for _, routerer := range param.Routerers {
		routerer.MustRegisterRoutes(g)
	}

	return &http.Server{
		Handler:      g,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadTimeout:  param.Conf.HTTPServerReadTimeout,
		WriteTimeout: param.Conf.HTTPServerWriteTimeout,
	}
}

func registerMiddlewares(router *gin.Engine, logger *slog.Logger) {
	router.Use(sloggin.NewWithConfig(logger, sloggin.Config{
		Filters: []sloggin.Filter{
			sloggin.IgnorePathContains("/healthzzz"),
		},
	}))

	router.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {
		logger.ErrorContext(ctx, "Panic occurred", slog.Any("panic", err))

		ctx.JSON(http.StatusInternalServerError, httpres.Error{Message: http.StatusText(http.StatusInternalServerError)})
	}))

	router.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, httpres.Error{Message: http.StatusText(http.StatusMethodNotAllowed)})
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, httpres.Error{Message: http.StatusText(http.StatusNotFound)})
	})

	// TODO: maybe some metrics endpoint such as prometheus.
}
