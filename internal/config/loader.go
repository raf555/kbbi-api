package config

import (
	"github.com/raf555/salome/config/v1"
)

// Deprecated: use [github.com/raf555/salome/config/v1.LoadConfigTo]
func EnvConfigProvider[T any](provider config.Provider) (T, error) {
	return config.LoadConfigTo[T](provider)
}
