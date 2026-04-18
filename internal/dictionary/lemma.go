package dictionary

import (
	"github.com/raf555/kbbi-api/pkg/kbbi"
)

// internal representation of [kbbi.Lemma]
type (
	Lemma struct {
		Lemma   string  `json:"lemma"`
		Entries []Entry `json:"entries"`
	}

	Entry struct {
		Entry            string            `json:"entry"`
		BaseWord         string            `json:"baseWord"`
		EntryVariants    []string          `json:"entryVariants"`
		Pronunciation    string            `json:"pronunciation"`
		Definitions      []EntryDefinition `json:"definitions"`
		NonStandardWords []string          `json:"nonStandardWords"`
		WordVariants     []string          `json:"variants"`
		CompoundWords    []string          `json:"compoundWords"`
		DerivedWords     []string          `json:"derivedWords"`
		Proverbs         []string          `json:"proverbs"`
		Metaphors        []string          `json:"metaphors"`
		IsPrecategorical bool              `json:"isPrecategorical"`
	}

	EntryDefinition struct {
		Definition       string       `json:"definition"`
		RawDefinition    string       `json:"rawDefinition"`
		ReferencedLemma  string       `json:"referencedLemma"`
		Labels           []EntryLabel `json:"labels"`
		UsageExamples    []string     `json:"usageExamples"`
		RawUsageExamples []string     `json:"rawUsageExamples"`
	}

	EntryLabel struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
)

func (l Lemma) ToKBBI() kbbi.Lemma {
	entries := make([]kbbi.Entry, len(l.Entries))
	for i := range l.Entries {
		entries[i] = l.Entries[i].toKBBI()
	}

	return kbbi.Lemma{
		Lemma:   l.Lemma,
		Entries: entries,
	}
}

func (e Entry) toKBBI() kbbi.Entry {
	definitions := make([]kbbi.EntryDefinition, len(e.Definitions))
	for i := range e.Definitions {
		definitions[i] = e.Definitions[i].toKBBI()
	}

	return kbbi.Entry{
		Entry:            e.Entry,
		BaseWord:         e.BaseWord,
		EntryVariants:    e.EntryVariants,
		Pronunciation:    e.Pronunciation,
		Definitions:      definitions,
		NonStandardWords: e.NonStandardWords,
		WordVariants:     e.WordVariants,
		CompoundWords:    e.CompoundWords,
		DerivedWords:     e.DerivedWords,
		Proverbs:         e.Proverbs,
		Metaphors:        e.Metaphors,
		IsPrecategorical: e.IsPrecategorical,
	}
}

func (d EntryDefinition) toKBBI() kbbi.EntryDefinition {
	labels := make([]kbbi.EntryLabel, len(d.Labels))
	for i := range d.Labels {
		labels[i] = d.Labels[i].toKBBI()
	}

	return kbbi.EntryDefinition{
		Definition:      d.Definition,
		ReferencedLemma: d.ReferencedLemma,
		Labels:          labels,
		UsageExamples:   d.UsageExamples,
	}
}

func (l EntryLabel) toKBBI() kbbi.EntryLabel {
	return kbbi.EntryLabel{
		Code: l.Code,
		Name: l.Name,
		Kind: l.Kind,
	}
}
