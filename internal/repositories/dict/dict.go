package dict

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/raf555/kbbi-api/internal/assets"
	"github.com/raf555/kbbi-api/internal/config"
	assets_model "github.com/raf555/kbbi-api/internal/models/assets"
	"github.com/raf555/kbbi-api/internal/repositories/wotd"
	"github.com/raf555/kbbi-api/pkg/kbbi"
)

type (
	Dictionary struct {
		wotd *wotd.Repository

		stats assets_model.Stats

		longestLemmaLength int

		inverseIndex map[string]int
		lemmas       []kbbi.Lemma
	}
)

var (
	ErrLemmaNotFound         = errors.New("lemma not found")
	ErrLemmaTooLong          = errors.New("lemma length too long")
	ErrEntryNotFound         = errors.New("entry not found")
	ErrUnexpectedEmptyLemma  = errors.New("unexpected empty lemma")
	ErrUnexpectedEntryNumber = errors.New("unexpected entry number")
	ErrUnexpectedRandomIndex = errors.New("unexpected random lemma index")
	ErrUnexpectedWotdIndex   = errors.New("unexpected wotd lemma index")
)

func New(env *config.Configuration, logger *slog.Logger, wotd *wotd.Repository) (*Dictionary, error) {
	start := time.Now()
	logger.Info("Started reading dictionary asset")
	var assetData assets_model.AssetData
	if err := assets.Read("dict.db", env.AssetsDirectory, env.AssetsEncryptionKey, env.AssetsEncryptionIV).To(&assetData); err != nil {
		return nil, fmt.Errorf("error reading the dictionary: %w", err)
	}
	logger.Info("Finished reading dictionary asset", slog.String("elapsed", time.Since(start).String()))

	longestLemmaLength := 0
	inverseIdx := make(map[string]int, len(assetData.Lemmas))

	for i, lemma := range assetData.Lemmas {
		inverseIdx[lemma.Lemma] = i

		lemmaLength := len(lemma.Lemma)
		if lemmaLength > longestLemmaLength {
			longestLemmaLength = lemmaLength
		}
	}

	return &Dictionary{
		wotd:               wotd,
		stats:              assetData.Stats,
		longestLemmaLength: longestLemmaLength,
		inverseIndex:       inverseIdx,
		lemmas:             assetData.Lemmas,
	}, nil
}

func (d *Dictionary) indexInDictRange(idx int) bool {
	return 0 <= idx && idx < len(d.lemmas)
}

func (d *Dictionary) Stats() assets_model.Stats {
	return d.stats
}

func (d *Dictionary) Lemma(lemma string, entryNoPtr *int) (kbbi.Lemma, error) {
	if lemma == "" {
		return kbbi.Lemma{}, ErrUnexpectedEmptyLemma
	}

	if len(lemma) > d.longestLemmaLength {
		return kbbi.Lemma{}, ErrLemmaTooLong
	}

	idx, ok := d.inverseIndex[lemma]
	if !ok {
		return kbbi.Lemma{}, ErrLemmaNotFound
	}

	lemmaData := d.lemmas[idx]

	if entryNoPtr != nil {
		entryNo := *entryNoPtr
		if entryNo <= 0 {
			return kbbi.Lemma{}, ErrUnexpectedEntryNumber
		}

		if entryNo > len(lemmaData.Entries) {
			return kbbi.Lemma{}, ErrEntryNotFound
		}

		entryIdx := entryNo - 1
		lemmaData.Entries = lemmaData.Entries[entryIdx:entryNo]
	}

	return lemmaData, nil
}

func (d *Dictionary) RandomLemma() (kbbi.Lemma, error) {
	randomIdx := rand.IntN(len(d.lemmas))

	return d.lemmas[randomIdx], nil
}

func (d *Dictionary) LemmaOfTheDay() (kbbi.Lemma, error) {
	wotdIdx := d.wotd.TodayLemmaIndex()
	idx := wotdIdx - 1

	if !d.indexInDictRange(idx) {
		return kbbi.Lemma{}, fmt.Errorf("%w: %d", ErrUnexpectedWotdIndex, idx)
	}

	return d.lemmas[idx], nil
}
