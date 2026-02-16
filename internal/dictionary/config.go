package dictionary

import "github.com/raf555/kbbi-api/internal/encoding"

type Configuration struct {
	WOTD       AssetConfig `env:",prefix=ASSETS_WOTD_"`
	Dictionary AssetConfig `env:",prefix=ASSETS_DICTIONARY_"`

	AssetsEncryptionKey encoding.HexString `env:"ASSETS_ENCRYPTION_KEY, required"`
	AssetsEncryptionIV  encoding.HexString `env:"ASSETS_ENCRYPTION_IV, required"`
	AssetsDirectory     string             `env:"ASSETS_DIRECTORY, default=./assets/"`
}

type AssetConfig struct {
	DownloadURL string `env:"DOWNLOAD_URL"`
}
