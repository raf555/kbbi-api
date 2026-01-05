package dictionary

import "github.com/raf555/kbbi-api/internal/encoding"

type Configuration struct {
	AssetsEncryptionKey encoding.HexString `env:"ASSETS_ENCRYPTION_KEY, required"`
	AssetsEncryptionIV  encoding.HexString `env:"ASSETS_ENCRYPTION_IV, required"`
	AssetsDirectory     string             `env:"ASSETS_DIRECTORY, default=./assets/"`
}
