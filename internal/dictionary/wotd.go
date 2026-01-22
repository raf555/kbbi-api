package dictionary

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type WOTD struct {
	lemmaIndexes []int
	epoch        int64
}

func NewWOTD(env Configuration) (*WOTD, error) {
	var lemmaIndexes []int
	if err := ReadAsset("wotd.db", env.AssetsDirectory, env.AssetsEncryptionKey, env.AssetsEncryptionIV).To(&lemmaIndexes); err != nil {
		return nil, fmt.Errorf("ReadAsset: %w", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("time.LoadLocation: %w", err)
	}

	repo := &WOTD{
		lemmaIndexes: lemmaIndexes,
		epoch:        time.Date(2022, time.October, 30, 0, 0, 0, 0, loc).UnixMilli(),
	}

	return repo, nil
}

func (w *WOTD) RandomLemmaIndex() int {
	return rand.IntN(len(w.lemmaIndexes))
}

func (w *WOTD) TodayLemmaIndex() int {
	currentTime := time.Now().UnixMilli()
	daysSinceEpoch := (currentTime - w.epoch) / (1000 * 60 * 60 * 24)
	j := daysSinceEpoch % int64(len(w.lemmaIndexes))

	return w.lemmaIndexes[j]
}
