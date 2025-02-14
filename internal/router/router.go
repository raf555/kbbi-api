package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/docs"
	"github.com/raf555/kbbi-api/internal/dto/response"
	"github.com/raf555/kbbi-api/internal/handlers/v1/entry"
	"github.com/raf555/kbbi-api/internal/handlers/v1/home"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(
	logger *slog.Logger,

	entryHandler *entry.Handler,
	homeHandler *home.Handler,
) *gin.Engine {

	router := gin.New()
	router.UseRawPath = true

	{
		router.Use(sloggin.New(logger))

		router.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {
			logger.ErrorContext(ctx, "Panic occured", slog.Any("panic", err))
			ctx.JSON(http.StatusInternalServerError, response.Error{Message: response.ErrInternalServerError})
		}))

		router.NoMethod(func(ctx *gin.Context) {
			ctx.JSON(http.StatusMethodNotAllowed, response.Error{Message: response.ErrMethodNotAllowed})
		})

		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, response.Error{Message: response.ErrNotFound})
		})
	}

	{
		router.GET("/", homeHandler.Home)
		router.GET("/healthzzz", homeHandler.Health)

		docs.SwaggerInfo.Title = "KBBI API"
		docs.SwaggerInfo.Description = "Probably the most complete KBBI API you will ever find."
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	apiGroupV1 := router.Group("/api/v1")
	{
		entryGroup := apiGroupV1.Group("/entry")
		{
			entryGroup.GET("/_wotd", entryHandler.WOTD)
			entryGroup.GET("/_random", entryHandler.Random)
			entryGroup.GET("/:entry", entryHandler.RedirectToLowercase, entryHandler.Entry)
		}
	}

	return router
}
