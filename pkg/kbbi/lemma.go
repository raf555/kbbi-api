// Package kbbi contains public structs that can be used by API clients.
// This package is directly used by the API server.
package kbbi

type (
	// Lemma contains the information of one lemma.
	// One lemma can hold multiple entries.
	// E.g. lemma `apel` consists entries of `apel (1)`, `apel (2)`, etc.
	Lemma struct {
		// Lemma is a single dictionary entry. E.g. `apel`.
		Lemma string `json:"lemma"`

		// Entries holds all entries information for this lemma.
		Entries []Entry `json:"entries"`
	}

	// Entry contains all informations related to the entry.
	// All fields will always be a non-nil value.
	Entry struct {
		// Entry is the entry word. E.g. `apel (1)`.
		Entry string `json:"entry"`

		// BaseWord is the base word for a given entry (if any).
		// I.e. `kata dasar`.
		// E.g. `menyukai` has a base word of `suka`.
		BaseWord string `json:"baseWord"`

		// EntryVariants contains variants of how the entry can be referred (if any).
		// E.g. `terselip` can be alternatively referred as `terselip ke luar`.
		//
		// It is possible that the variant does not have any entries in the dictionary.
		EntryVariants []string `json:"entryVariants"`

		// Pronunciation describes the way in which a word is prononunced (if any).
		// E.g. `apel` can be prononunced as apêl.
		Pronunciation string `json:"pronunciation"`

		// Definitions contains the meaning of the entry.
		// A single entry can have multiple meanings or definitions.
		// E.g. `suka` has multiple meanings depending on the context.
		//
		// Each definition has its own information, such as labels and usage examples.
		//
		// The definitions can be empty depending on the entry.
		// If it is empty, usually it can be referred from the information of the other fields (e.g. BaseWord).
		Definitions []EntryDefinition `json:"definitions"`

		// NonStandardWords contains the non-standard forms of the entry (if any).
		// I.e. `bentuk tidak baku`.
		// E.g. `apotek` has a non-standard form of `apotik`.
		NonStandardWords []string `json:"nonStandardWords"`

		// WordVariants contains the alternative words of the entry (if any).
		// I.e. `varian`.
		// E.g. `ude` has a alternative word of `udeh`.
		//
		// The difference between WordVariants and `EntryVariants` is that
		// WordVariants guaranteed to have at least 1 entry in the dictionary.
		WordVariants []string `json:"variants"`

		// CompoundWords contains the compound words of the entry (if any).
		// I.e. `gabungan kata`.
		// E.g. `kacang` has a compound word of `kacang atom`.
		CompoundWords []string `json:"compoundWords"`

		// DerivedWords contains the derived words of the entry (if any).
		// I.e. `kata turunan`.
		// E.g. `suka` has a derived word of `menyukai`.
		DerivedWords []string `json:"derivedWords"`

		// Proverbs contains the proverbs of the entry (if any).
		// I.e. `peribahasa`.
		// E.g. `kacang` is used in `kacang lupa akan kulitnya` proverb.
		Proverbs []string `json:"proverbs"`

		// Metaphors contains metaphors of this entry (if any).
		// I.e. `kiasan`.
		// E.g. `leher` is used in `leher terasa panjang` metaphor.
		Metaphors []string `json:"metaphors"`
	}

	// EntryDefinition contains the detail of the entry's definition.
	EntryDefinition struct {
		// Definition contains the meaning of the entry.
		Definition string `json:"definition"`

		// ReferencedLemma contains referenced lemma in the definition if present.
		//
		// Some entries have no direct meaning, so instead it refers the other lemma as the definition.
		// Usually it has the definition of `lihat [lemma]`.
		//
		// In other case, the entry is usually a non-standard form of the other lemma.
		// Usually it has the definition of `bentuk tidak baku dari [lemma]`.
		ReferencedLemma string `json:"referencedLemma"`

		// Labels contains the label of this definition if present.
		// In the dictionary, they are usually placed at the front of the meaning.
		// E.g. `su.ka a cak mudah sekali ...; kerap kali ...`
		Labels []EntryLabel `json:"labels"`

		// UsageExamples contains usage example of the entry for this meaning if any.
		// In the dictionary, they are usually placed at the end of the meaning.
		// E.g. `su.ka a cak mudah sekali ...; kerap kali ...: memang dia -- lupa; pensil semacam ini -- patah`
		UsageExamples []string `json:"usageExamples"`
	}

	// EntryLabel contains the label information of the entry for a definition.
	EntryLabel struct {
		// Code is the label short form.
		// E.g. `n`, `Huk`, `cak`, etc.
		Code string `json:"code"`

		// Code is the label actual name.
		// E.g. `nomina`, `Hukum`, `cakapan`, etc.
		Name string `json:"name"`

		// Kind is the label kind.
		// E.g. `Kelas Kata`, `Bidang`, `Ragam`, etc.
		Kind string `json:"kind"`
	}
)
