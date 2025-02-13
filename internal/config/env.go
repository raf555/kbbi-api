package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raf555/kbbi-api/internal/models/key"
	"github.com/sethvargo/go-envconfig"
)

type (
	Configuration struct {
		Port int `env:"PORT, default=8888"`

		HTTPServerReadTimeout  time.Duration `env:"HTTP_SERVER_READ_TIMEOUT, default=1s"`
		HTTPServerWriteTimeout time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT, default=1s"`

		AssetsEncryptionKey key.HexString `env:"ASSETS_ENCRYPTION_KEY, required"`
		AssetsEncryptionIV  key.HexString `env:"ASSETS_ENCRYPTION_IV, required"`
	}
)

func init() {
	if err := readAndLoadEnvFile(); err != nil {
		panic(fmt.Errorf("failed loading .env file: %w", err))
	}
}

func ReadConfig() (*Configuration, error) {
	ctx := context.Background()

	var c Configuration
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("failed loading environment variables: %w", err)
	}

	return &c, nil
}

func readAndLoadEnvFile() error {
	_, err := os.Stat(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
