package dictionary

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.com/raf555/kbbi-api/pkg/kbbi"
	"github.com/samber/lo"
)

type Dictionary struct {
	wotd                   WOTDRepo
	stats                  Stats
	longestLemmaLength     int
	inverseIndex           map[string]*lemmaIndex
	inverseNormalizedIndex map[string]*lemmaIndex
	lemmas                 []wrappedLemma
}

type lemmaIndex struct {
	idx        int           // index in lemmas.
	entryNoMap map[int][]int // key is entry number, starts from 1. value is the actual index in the entries.
}

type wrappedLemma struct {
	kbbi.Lemma

	NormalizedForm string
}

func NewDictionary(cfg Configuration, logger *slog.Logger, wotd WOTDRepo) (*Dictionary, error) {
	start := time.Now()
	logger.Info("Started reading dictionary asset")

	var assetData AssetData

	var reader *reader
	if url := cfg.Dictionary.DownloadURL; url != "" {
		logger.Info("reading dictionary asset from URL", slog.String("url", url))
		reader = ReadAssetFromURL(url, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV)
	} else {
		reader = ReadAsset("dict.db", cfg.AssetsDirectory, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV)
	}

	if err := reader.To(&assetData); err != nil {
		return nil, fmt.Errorf("ReadAsset: %w", err)
	}

	logger.Info("Finished reading dictionary asset", slog.String("elapsed", time.Since(start).String()))

	longestLemmaLength := 0
	inverseIdx := make(map[string]*lemmaIndex, len(assetData.Lemmas))
	inverseNormalizedIndex := make(map[string]*lemmaIndex)
	lemmas := make([]wrappedLemma, 0, len(assetData.Lemmas))

	for i, lemma := range assetData.Lemmas {
		idx := &lemmaIndex{
			idx:        i,
			entryNoMap: map[int][]int{},
		}

		// lookup and map entry index if any
		for j, def := range lemma.Entries {
			_, entryNo, ok := FindEntryNoFromLemma(def.Entry)
			if !ok {
				continue
			}

			// there can be multiple entries with same number. E.g. ketak (4)
			// could be misinput from KBBI but for now making the behavior the same as the website.
			idx.entryNoMap[entryNo] = append(idx.entryNoMap[entryNo], j)
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
		lemmas = append(lemmas, wrappedLemma{
			Lemma:          lemma,
			NormalizedForm: Normalize(lemma.Lemma, true),
		})
	}

	return &Dictionary{
		wotd:                   wotd,
		stats:                  assetData.Stats,
		longestLemmaLength:     longestLemmaLength,
		inverseIndex:           inverseIdx,
		inverseNormalizedIndex: inverseNormalizedIndex,
		lemmas:                 lemmas,
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

		entryIndexes, ok := index.entryNoMap[entryNo]
		if !ok {
			return kbbi.Lemma{}, ErrEntryNotFound
		}

		lemmaData.Entries = lo.Map(entryIndexes, func(idx int, _ int) kbbi.Entry {
			return lemmaData.Entries[idx]
		})
	}

	return lemmaData.Lemma, nil
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
	return d.lemmas[randomIdx].Lemma
}

func (d *Dictionary) LemmaOfTheDay() (kbbi.Lemma, error) {
	wotdIdx := d.wotd.TodayLemmaIndex()
	idx := wotdIdx - 1

	if !d.indexInDictRange(idx) {
		return kbbi.Lemma{}, fmt.Errorf("%w: %d", ErrUnexpectedWotdIndex, idx)
	}

	return d.lemmas[idx].Lemma, nil
}

// Search provides a list of lemmas based on prefix, number of result depends on limit.
// Search behaves similarly with search feature on the KBBI application.
//
// If prefix is empty, Search returns top limit lemmas.
func (d *Dictionary) Search(prefix string, limit uint) []kbbi.Lemma {
	if prefix == "" {
		return lo.Map(d.lemmas[:min(int(limit), len(d.lemmas))], func(lemma wrappedLemma, _ int) kbbi.Lemma { return lemma.Lemma })
	}

	prefix = strings.ToLower(Normalize(prefix, true))

	leftIdx, _ := slices.BinarySearchFunc(d.lemmas, prefix, func(curr wrappedLemma, search string) int {
		return strings.Compare(curr.NormalizedForm, search)
	})

	rightIdx, _ := slices.BinarySearchFunc(d.lemmas, prefix+"\uffff", func(curr wrappedLemma, search string) int {
		return strings.Compare(curr.NormalizedForm, search)
	})

	if leftIdx >= len(d.lemmas) || rightIdx < leftIdx || !strings.HasPrefix(d.lemmas[leftIdx].NormalizedForm, prefix) {
		return nil
	}

	return lo.Map(d.lemmas[leftIdx:rightIdx][:min(limit, uint(rightIdx-leftIdx))], func(lemma wrappedLemma, _ int) kbbi.Lemma { return lemma.Lemma })
}
