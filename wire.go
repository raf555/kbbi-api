//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/raf555/kbbi-api/internal/config"
	"github.com/raf555/kbbi-api/internal/handlers/v1/entry"
	"github.com/raf555/kbbi-api/internal/handlers/v1/home"
	"github.com/raf555/kbbi-api/internal/logger"
	"github.com/raf555/kbbi-api/internal/repositories/dict"
	"github.com/raf555/kbbi-api/internal/repositories/wotd"
	"github.com/raf555/kbbi-api/internal/router"
)

var (
	noDepSet = wire.NewSet(
		config.ReadConfig,
	)

	mainDepSet = wire.NewSet(
		logger.NewLogger,
	)

	repoSet = wire.NewSet(
		wotd.New,
		dict.New,
	)

	handlerSet = wire.NewSet(
		home.New,
		entry.New,
	)

	httpServerSet = wire.NewSet(
		router.New,
		NewHTTPServer,
		wire.Struct(
			new(MainDependency),
			"HTTPServer",
			"Logger",
		),
	)
)

func InitializeDependency() (*MainDependency, error) {
	wire.Build(
		noDepSet,
		mainDepSet,
		repoSet,
		handlerSet,
		httpServerSet,
	)
	return nil, nil
}
