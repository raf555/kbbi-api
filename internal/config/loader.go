package config

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func Load(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("os.Stat: %w", err)
	}

	err = godotenv.Load(file)
	if err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}

	return nil
}

var (
	validatorOnce sync.Once
	vl            *validator.Validate
)

func LoadEnvTo[T any](dst *T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := envconfig.Process(ctx, dst); err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}

	validatorOnce.Do(func() {
		vl = validator.New()
	})

	if err := vl.StructCtx(ctx, dst); err != nil {
		return fmt.Errorf("vl.StructCtx: %w", err)
	}

	return nil
}

func EnvConfigProvider[T any]() (T, error) {
	var dst T
	if err := LoadEnvTo(&dst); err != nil {
		var zero T
		return zero, fmt.Errorf("LoadEnvTo: %w", err)
	}
	return dst, nil
}
