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
		wotd                   WOTDRepo
		stats                  Stats
		longestLemmaLength     int
		inverseIndex           map[string]*lemmaIndex
		inverseNormalizedIndex map[string]*lemmaIndex
		lemmas                 []kbbi.Lemma
	}
)

func NewDictionary(cfg Configuration, logger *slog.Logger, wotd WOTDRepo) (*Dictionary, error) {
	start := time.Now()
	logger.Info("Started reading dictionary asset")
	var assetData AssetData
	if err := ReadAsset("dict.db", cfg.AssetsDirectory, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV).To(&assetData); err != nil {
		return nil, fmt.Errorf("ReadAsset: %w", err)
	}
	logger.Info("Finished reading dictionary asset", slog.String("elapsed", time.Since(start).String()))

	longestLemmaLength := 0
	inverseIdx := make(map[string]*lemmaIndex, len(assetData.Lemmas))
	inverseNormalizedIndex := make(map[string]*lemmaIndex)

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
		if normalized := Normalize(lemma.Lemma, false); normalized != lemma.Lemma { // lemma has normalized form
			// p.s. not removing punctuation here to make exact match.
			// Don't want `s.t` to have the result of `st.` or other similar case since it's probably wrong.
			// So for now only care for removing diacritics.
			//
			// If the inverseNormalizedIndex of the normalized lemma is already occupied, ignore (only use the first one).
			if _, ok := inverseNormalizedIndex[normalized]; !ok {
				inverseNormalizedIndex[normalized] = idx
			}
		}

		lemmaLength := len(lemma.Lemma)
		longestLemmaLength = max(longestLemmaLength, lemmaLength)
	}

	return &Dictionary{
		wotd:                   wotd,
		stats:                  assetData.Stats,
		longestLemmaLength:     longestLemmaLength,
		inverseIndex:           inverseIdx,
		inverseNormalizedIndex: inverseNormalizedIndex,
		lemmas:                 assetData.Lemmas,
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

	index := d.lookupInverseIndex(lemma)
	if index == nil {
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

func (d *Dictionary) lookupInverseIndex(lemma string) *lemmaIndex {
	// lookup on exact index first
	index, ok := d.inverseIndex[lemma]
	if ok {
		return index
	}

	// if not found, normalize the lemma, and check on the normalized index
	normalized := Normalize(lemma, false)
	index, ok = d.inverseNormalizedIndex[normalized]
	if ok {
		return index
	}

	// otherwise not found
	return nil
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
