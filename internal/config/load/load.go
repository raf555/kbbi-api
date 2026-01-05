package load

import (
	"fmt"

	"github.com/raf555/kbbi-api/internal/config"
)

var EnvFilename = ".env"

func init() {
	if err := config.Load(EnvFilename); err != nil {
		panic(fmt.Errorf("config.Load(%s): %w", EnvFilename, err))
	}
}
