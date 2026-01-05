package home

import "github.com/raf555/kbbi-api/internal/dictionary"

type AssetStatsFetcher interface {
	Stats() dictionary.Stats
}
