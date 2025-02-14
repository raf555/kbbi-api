package wotd

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/raf555/kbbi-api/internal/assets"
	"github.com/raf555/kbbi-api/internal/config"
)

type (
	Repository struct {
		lemmaIndexes []int
		epoch        int64
	}
)

func New(env *config.Configuration) (*Repository, error) {
	var lemmaIndexes []int
	if err := assets.Read("wotd.db", env.AssetsEncryptionKey, env.AssetsEncryptionIV).To(&lemmaIndexes); err != nil {
		return nil, fmt.Errorf("error reading the WOTD db: %w", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("failed loading Asia/Jakarta time location: %w", err)
	}

	repo := &Repository{
		lemmaIndexes: lemmaIndexes,
		epoch:        time.Date(2022, 10, 30, 0, 0, 0, 0, loc).UnixMilli(),
	}

	return repo, nil
}

func (w *Repository) RandomLemmaIndex() int {
	return rand.IntN(len(w.lemmaIndexes))
}

func (w *Repository) TodayLemmaIndex() int {
	currentTime := time.Now().UnixMilli()
	daysSinceEpoch := (currentTime - w.epoch) / (1000 * 60 * 60 * 24)
	j := daysSinceEpoch % int64(len(w.lemmaIndexes))

	return w.lemmaIndexes[j]
}
