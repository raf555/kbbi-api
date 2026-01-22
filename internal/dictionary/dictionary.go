package dictionary

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/raf555/kbbi-api/pkg/kbbi"
)

type (
	lemmaIndex struct {
		idx        int         // index in lemmas.
		entryNoMap map[int]int // key is entry number, starts from 1. value is the actual index in the entries.
	}

	Dictionary struct {
		wotd               WOTDRepo
		stats              Stats
		longestLemmaLength int
		inverseIndex       map[string]*lemmaIndex
		lemmas             []kbbi.Lemma
	}
)

func NewDictionary(cfg Configuration, logger *slog.Logger, wotd WOTDRepo) (*Dictionary, error) {
	start := time.Now()
	logger.Info("Started reading dictionary asset")
	var assetData AssetData
	if err := ReadAsset("dict.db", cfg.AssetsDirectory, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV).To(&assetData); err != nil {
		return nil, fmt.Errorf("assets.Read: %w", err)
	}
	logger.Info("Finished reading dictionary asset", slog.String("elapsed", time.Since(start).String()))

	longestLemmaLength := 0
	inverseIdx := make(map[string]*lemmaIndex, len(assetData.Lemmas))

	for i, lemma := range assetData.Lemmas {
		idx := &lemmaIndex{
			idx:        i,
			entryNoMap: map[int]int{},
		}

		// lookup and map entry index if any
		for j, def := range lemma.Entries {
			_, entryNo, ok := FindEntryNoFromLemma(def.Entry)
			if !ok {
				continue
			}

			idx.entryNoMap[entryNo] = j
		}

		inverseIdx[lemma.Lemma] = idx

		lemmaLength := len(lemma.Lemma)
		longestLemmaLength = max(longestLemmaLength, lemmaLength)
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

func (d *Dictionary) Stats() Stats {
	return d.stats
}

func (d *Dictionary) Lemma(lemma string, entryNo int) (kbbi.Lemma, error) {
	if lemma == "" {
		return kbbi.Lemma{}, ErrUnexpectedEmptyLemma
	}

	if len(lemma) > d.longestLemmaLength {
		return kbbi.Lemma{}, ErrLemmaTooLong
	}

	index, ok := d.inverseIndex[lemma]
	if !ok {
		return kbbi.Lemma{}, ErrLemmaNotFound
	}

	lemmaData := d.lemmas[index.idx]

	if entryNo < 0 {
		return kbbi.Lemma{}, ErrUnexpectedEntryNumber
	}

	if entryNo > 0 {
		if entryNo > len(lemmaData.Entries) {
			return kbbi.Lemma{}, ErrEntryNotFound
		}

		// first, lookup for entry index in the map.
		// if found, use that.
		// otherwise, fallback to the index in the entries list.
		entryIdx, ok := index.entryNoMap[entryNo]
		if !ok {
			entryIdx = entryNo - 1
		}

		lemmaData.Entries = lemmaData.Entries[entryIdx : entryIdx+1]
	}

	return lemmaData, nil
}

func (d *Dictionary) RandomLemma() kbbi.Lemma {
	randomIdx := rand.IntN(len(d.lemmas))
	return d.lemmas[randomIdx]
}

func (d *Dictionary) LemmaOfTheDay() (kbbi.Lemma, error) {
	wotdIdx := d.wotd.TodayLemmaIndex()
	idx := wotdIdx - 1

	if !d.indexInDictRange(idx) {
		return kbbi.Lemma{}, fmt.Errorf("%w: %d", ErrUnexpectedWotdIndex, idx)
	}

	return d.lemmas[idx], nil
}
