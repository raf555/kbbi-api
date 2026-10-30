package dictionary

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Dictionary struct {
	wotd               WOTDRepo
	stats              Stats
	longestLemmaLength int
	// lemmas keeps the asset order, which is sorted by NormalizedForm. Search depends on that.
	lemmas []wrappedLemma
	// index maps an exact lemma to its index in lemmas.
	index map[string]int32
	// normalizedIndex maps a lemma without diacritics to its index in lemmas.
	// It only holds lemmas whose normalized form differs from the lemma (a few hundred).
	normalizedIndex map[string]int32
}

type wrappedLemma struct {
	Lemma

	NormalizedForm string
}

func NewDictionary(cfg Configuration, logger *slog.Logger, wotd WOTDRepo) (*Dictionary, error) {
	start := time.Now()
	logger.Info("Started reading dictionary asset")

	var reader *reader
	if url := cfg.Dictionary.DownloadURL; url != "" {
		logger.Info("reading dictionary asset from URL", slog.String("url", url))
		reader = ReadAssetFromURL(url, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV)
	} else {
		reader = ReadAsset("dict.db", cfg.AssetsDirectory, cfg.AssetsEncryptionKey, cfg.AssetsEncryptionIV)
	}

	var dict *Dictionary
	err := reader.Read(func(plain io.Reader) error {
		var err error
		dict, err = loadDictionary(plain, wotd)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("ReadAsset: %w", err)
	}

	logger.Info("Finished reading dictionary asset", slog.String("elapsed", time.Since(start).String()))

	return dict, nil
}

// loadDictionary builds a Dictionary from the plain JSON asset.
// Lemmas are converted to their compact form one at a time, so the decoded asset is never fully in memory.
func loadDictionary(plain io.Reader, wotd WOTDRepo) (*Dictionary, error) {
	d := &Dictionary{
		wotd:            wotd,
		normalizedIndex: map[string]int32{},
	}
	labels := labelInterner{}

	add := func(l lemmaJSON) {
		if d.lemmas == nil {
			// "stats" comes before "lemmas" in the asset, so the count is known by now.
			d.lemmas = make([]wrappedLemma, 0, d.stats.LemmaCount)
			d.index = make(map[string]int32, d.stats.LemmaCount)
		}

		idx := int32(len(d.lemmas))
		d.index[l.Lemma] = idx
		if normalized := Normalize(l.Lemma, false); normalized != l.Lemma {
			// Only diacritics are removed, not punctuation, to keep the match exact.
			// `s.t` should not resolve to `st.` or similar, since that is probably wrong.
			// The first lemma with a given normalized form wins.
			if _, ok := d.normalizedIndex[normalized]; !ok {
				d.normalizedIndex[normalized] = idx
			}
		}

		normalizedForm := Normalize(l.Lemma, true)
		if normalizedForm == l.Lemma { // reuse lemma string when possible.
			normalizedForm = l.Lemma
		}

		d.longestLemmaLength = max(d.longestLemmaLength, len(l.Lemma))
		d.lemmas = append(d.lemmas, wrappedLemma{
			Lemma:          l.compact(labels),
			NormalizedForm: normalizedForm,
		})
	}

	if err := decodeDictionaryAsset(plain, &d.stats, add); err != nil {
		return nil, err
	}

	return d, nil
}

// decodeDictionaryAsset reads `{"stats": ..., "lemmas": [...]}` and calls add for each lemma as it is decoded.
func decodeDictionaryAsset(plain io.Reader, stats *Stats, add func(lemmaJSON)) error {
	dec := jsontext.NewDecoder(plain)

	if err := readDelim(dec, '{'); err != nil {
		return err
	}

	for dec.PeekKind() != '}' {
		name, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("read field name: %w", err)
		}

		switch name.String() {
		case "stats":
			if err := json.UnmarshalDecode(dec, stats); err != nil {
				return fmt.Errorf("decode stats: %w", err)
			}
		case "lemmas":
			if err := readDelim(dec, '['); err != nil {
				return err
			}
			for dec.PeekKind() != ']' {
				var l lemmaJSON
				if err := json.UnmarshalDecode(dec, &l); err != nil {
					return fmt.Errorf("decode lemma: %w", err)
				}
				add(l)
			}
			if err := readDelim(dec, ']'); err != nil {
				return err
			}
		default:
			if err := dec.SkipValue(); err != nil {
				return fmt.Errorf("skip %q: %w", name.String(), err)
			}
		}
	}

	return readDelim(dec, '}')
}

func readDelim(dec *jsontext.Decoder, want jsontext.Kind) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return fmt.Errorf("read %v: %w", want, err)
	}
	if tok.Kind() != want {
		return fmt.Errorf("want %v, got %v", want, tok.Kind())
	}
	return nil
}

func (d *Dictionary) indexInDictRange(idx int) bool {
	return 0 <= idx && idx < len(d.lemmas)
}

func (d *Dictionary) Stats() Stats {
	return d.stats
}

func (d *Dictionary) Lemma(lemma string, entryNo int) (Lemma, error) {
	if lemma == "" {
		return Lemma{}, ErrUnexpectedEmptyLemma
	}

	if len(lemma) > d.longestLemmaLength {
		return Lemma{}, ErrLemmaTooLong
	}

	idx, ok := d.lookup(lemma)
	if !ok {
		return Lemma{}, ErrLemmaNotFound
	}

	lemmaData := d.lemmas[idx].Lemma

	if entryNo < 0 {
		return Lemma{}, ErrUnexpectedEntryNumber
	}

	if entryNo > 0 {
		// There can be multiple entries with the same number, e.g. ketak (4).
		// It could be a KBBI input mistake, but the website returns all of them, so do the same.
		var entries []Entry
		for _, entry := range lemmaData.Entries {
			if _, no, ok := FindEntryNoFromLemma(entry.Entry); ok && no == entryNo {
				entries = append(entries, entry)
			}
		}
		if len(entries) == 0 {
			return Lemma{}, ErrEntryNotFound
		}

		lemmaData.Entries = entries
	}

	return lemmaData, nil
}

// lookup returns the index in d.lemmas of an exact match, falling back to the lemma without diacritics.
func (d *Dictionary) lookup(lemma string) (int, bool) {
	if idx, ok := d.index[lemma]; ok {
		return int(idx), true
	}

	idx, ok := d.normalizedIndex[Normalize(lemma, false)]
	return int(idx), ok
}

func (d *Dictionary) RandomLemma() Lemma {
	randomIdx := rand.IntN(len(d.lemmas))
	return d.lemmas[randomIdx].Lemma
}

func (d *Dictionary) LemmaOfTheDay() (Lemma, error) {
	wotdIdx := d.wotd.TodayLemmaIndex()
	idx := wotdIdx - 1

	if !d.indexInDictRange(idx) {
		return Lemma{}, fmt.Errorf("%w: %d", ErrUnexpectedWotdIndex, idx)
	}

	return d.lemmas[idx].Lemma, nil
}

// Search provides a list of lemmas based on prefix, number of result depends on limit.
// Search behaves similarly with search feature on the KBBI application.
//
// If prefix is empty, Search returns top limit lemmas.
func (d *Dictionary) Search(prefix string, limit uint) []Lemma {
	if prefix == "" {
		return lo.Map(d.lemmas[:min(limit, uint(len(d.lemmas)))], func(lemma wrappedLemma, _ int) Lemma { return lemma.Lemma })
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

	return lo.Map(d.lemmas[leftIdx:rightIdx][:min(limit, uint(rightIdx-leftIdx))], func(lemma wrappedLemma, _ int) Lemma { return lemma.Lemma })
}
