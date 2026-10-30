package dictionary

import (
	"strings"

	"github.com/raf555/kbbi-api/pkg/kbbi"
	"github.com/samber/lo"
)

// Internal representation of [kbbi.Lemma].
// The whole dictionary lives in memory, so these types are laid out to keep it small.
// Convert with [Lemma.ToKBBI] before serving.
type (
	Lemma struct {
		Lemma   string
		Entries []Entry
	}

	Entry struct {
		Entry         string
		BaseWord      string
		Pronunciation string
		Definitions   []EntryDefinition
		// RelatedEntry is nil when all of its lists are empty, which is true for ~90% of entries.
		RelatedEntry     *EntryRelated
		IsPrecategorical bool
	}

	EntryRelated struct {
		EntryVariants    []string
		NonStandardWords []string
		WordVariants     []string
		CompoundWords    []string
		DerivedWords     []string
		Proverbs         []string
		Metaphors        []string
	}

	EntryDefinition struct {
		Definition      string
		ReferencedLemma string
		// Labels may share its backing array with other definitions; never modify it.
		Labels []EntryLabel
		// Extra is nil when the definition has no usage examples and its raw text equals Definition,
		// which is true for ~78% of definitions.
		Extra *EntryDefinitionExtra
	}

	EntryDefinitionExtra struct {
		RawDefinition    string
		UsageExamples    []string
		RawUsageExamples []string
	}

	EntryLabel struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
)

// Asset representation of [Lemma] and [Entry].
type (
	lemmaJSON struct {
		Lemma   string      `json:"lemma"`
		Entries []entryJSON `json:"entries"`
	}

	entryJSON struct {
		Entry            string           `json:"entry"`
		BaseWord         string           `json:"baseWord"`
		EntryVariants    []string         `json:"entryVariants"`
		Pronunciation    string           `json:"pronunciation"`
		Definitions      []definitionJSON `json:"definitions"`
		NonStandardWords []string         `json:"nonStandardWords"`
		WordVariants     []string         `json:"variants"`
		CompoundWords    []string         `json:"compoundWords"`
		DerivedWords     []string         `json:"derivedWords"`
		Proverbs         []string         `json:"proverbs"`
		Metaphors        []string         `json:"metaphors"`
		IsPrecategorical bool             `json:"isPrecategorical"`
	}

	definitionJSON struct {
		Definition       string       `json:"definition"`
		RawDefinition    string       `json:"rawDefinition"`
		ReferencedLemma  string       `json:"referencedLemma"`
		Labels           []EntryLabel `json:"labels"`
		UsageExamples    []string     `json:"usageExamples"`
		RawUsageExamples []string     `json:"rawUsageExamples"`
	}
)

// noRelated stands in for a nil [Entry.RelatedEntry] so responses still get empty lists instead of null.
var noRelated = EntryRelated{
	EntryVariants:    []string{},
	NonStandardWords: []string{},
	WordVariants:     []string{},
	CompoundWords:    []string{},
	DerivedWords:     []string{},
	Proverbs:         []string{},
	Metaphors:        []string{},
}

func (l lemmaJSON) compact(labels labelInterner) Lemma {
	entries := make([]Entry, len(l.Entries))
	for i, e := range l.Entries {
		entries[i] = e.compact(labels)
	}

	return Lemma{
		Lemma:   l.Lemma,
		Entries: entries,
	}
}

func (e entryJSON) compact(labels labelInterner) Entry {
	definitions := make([]EntryDefinition, len(e.Definitions))
	for i, d := range e.Definitions {
		definitions[i] = d.compact(labels)
	}

	var related *EntryRelated
	hasRelated := len(e.EntryVariants) > 0 || len(e.NonStandardWords) > 0 || len(e.WordVariants) > 0 ||
		len(e.CompoundWords) > 0 || len(e.DerivedWords) > 0 || len(e.Proverbs) > 0 || len(e.Metaphors) > 0
	if hasRelated {
		related = &EntryRelated{
			EntryVariants:    e.EntryVariants,
			NonStandardWords: e.NonStandardWords,
			WordVariants:     e.WordVariants,
			CompoundWords:    e.CompoundWords,
			DerivedWords:     e.DerivedWords,
			Proverbs:         e.Proverbs,
			Metaphors:        e.Metaphors,
		}
	}

	return Entry{
		Entry:            e.Entry,
		BaseWord:         e.BaseWord,
		Pronunciation:    e.Pronunciation,
		Definitions:      definitions,
		RelatedEntry:     related,
		IsPrecategorical: e.IsPrecategorical,
	}
}

func (d definitionJSON) compact(labels labelInterner) EntryDefinition {
	// Only drop the extras when toKBBI can rebuild the same output without them.
	// The usage lists must be non-nil: an asset without raw fields (e.g. 2023-10) serves them as null.
	var extra *EntryDefinitionExtra
	noExtra := d.RawDefinition == d.Definition &&
		d.UsageExamples != nil && len(d.UsageExamples) == 0 &&
		d.RawUsageExamples != nil && len(d.RawUsageExamples) == 0
	if !noExtra {
		extra = &EntryDefinitionExtra{
			RawDefinition:    d.RawDefinition,
			UsageExamples:    d.UsageExamples,
			RawUsageExamples: d.RawUsageExamples,
		}
	}

	return EntryDefinition{
		Definition:      d.Definition,
		ReferencedLemma: d.ReferencedLemma,
		Labels:          labels.intern(d.Labels),
		Extra:           extra,
	}
}

// labelInterner makes definitions with the same labels share one slice.
// The dictionary has ~260k label lists but fewer than 2k distinct ones.
type labelInterner map[string][]EntryLabel

func (in labelInterner) intern(labels []EntryLabel) []EntryLabel {
	if len(labels) == 0 {
		return labels
	}

	var key strings.Builder
	for _, l := range labels {
		key.WriteString(l.Code)
		key.WriteByte(0)
		key.WriteString(l.Name)
		key.WriteByte(0)
		key.WriteString(l.Kind)
		key.WriteByte(0)
	}

	if shared, ok := in[key.String()]; ok {
		return shared
	}
	in[key.String()] = labels
	return labels
}

func (l Lemma) ToKBBI(raw bool) kbbi.Lemma {
	entries := make([]kbbi.Entry, len(l.Entries))
	for i := range l.Entries {
		entries[i] = l.Entries[i].toKBBI(raw)
	}

	return kbbi.Lemma{
		Lemma:   l.Lemma,
		Entries: entries,
	}
}

func (e Entry) toKBBI(raw bool) kbbi.Entry {
	definitions := make([]kbbi.EntryDefinition, len(e.Definitions))
	for i := range e.Definitions {
		definitions[i] = e.Definitions[i].toKBBI(raw)
	}

	related := e.RelatedEntry
	if related == nil {
		related = &noRelated
	}

	return kbbi.Entry{
		Entry:            e.Entry,
		BaseWord:         e.BaseWord,
		EntryVariants:    related.EntryVariants,
		Pronunciation:    e.Pronunciation,
		Definitions:      definitions,
		NonStandardWords: related.NonStandardWords,
		WordVariants:     related.WordVariants,
		CompoundWords:    related.CompoundWords,
		DerivedWords:     related.DerivedWords,
		Proverbs:         related.Proverbs,
		Metaphors:        related.Metaphors,
		IsPrecategorical: e.IsPrecategorical,
	}
}

func (d EntryDefinition) toKBBI(raw bool) kbbi.EntryDefinition {
	labels := make([]kbbi.EntryLabel, len(d.Labels))
	for i := range d.Labels {
		labels[i] = d.Labels[i].toKBBI()
	}

	extra := d.Extra
	if extra == nil {
		extra = &EntryDefinitionExtra{
			RawDefinition:    d.Definition,
			UsageExamples:    []string{},
			RawUsageExamples: []string{},
		}
	}

	return kbbi.EntryDefinition{
		Definition:      lo.Ternary(raw, extra.RawDefinition, d.Definition),
		ReferencedLemma: d.ReferencedLemma,
		Labels:          labels,
		UsageExamples:   lo.Ternary(raw, extra.RawUsageExamples, extra.UsageExamples),
	}
}

func (l EntryLabel) toKBBI() kbbi.EntryLabel {
	return kbbi.EntryLabel{
		Code: l.Code,
		Name: l.Name,
		Kind: l.Kind,
	}
}
