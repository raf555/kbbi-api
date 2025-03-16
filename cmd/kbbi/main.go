package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raf555/kbbi-api/internal/config"
)

type (
	MainDependency struct {
		HTTPServer *http.Server
		Logger     *slog.Logger
	}
)

func main() {
	dep, err := InitializeDependency()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := dep.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			dep.Logger.Error("Server returned a non server closed error", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	dep.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := dep.HTTPServer.Shutdown(ctx); err != nil {
		dep.Logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		panic(err)
	}

	dep.Logger.Info("Server exited")
}

func NewHTTPServer(config *config.Configuration, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(config.Port)),
		Handler:      router,
		ReadTimeout:  config.HTTPServerReadTimeout,
		WriteTimeout: config.HTTPServerWriteTimeout,
	}
}
