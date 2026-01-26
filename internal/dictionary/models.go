package dictionary

import (
	"github.com/raf555/kbbi-api/pkg/kbbi"
)

type AssetData struct {
	Stats  Stats        `json:"stats"`
	Lemmas []kbbi.Lemma `json:"lemmas"`
}

type Stats struct {
	Edition    string `json:"edition"`
	EntryCount int    `json:"entryCount"`
	LemmaCount int    `json:"lemmaCount"`
}

type EntryRequest struct {
	Lemma string `uri:"entry" validate:"required"`
	// EntryNo is optional; value 0 means "no specific entry number requested".
	EntryNo int `form:"entryNo" validate:"gte=0"`
}

// transform mutates the EntryRequest in place by looking for an entry number in the lemma string.
// If an entry number is present, it updates Lemma to exclude the number and sets EntryNo accordingly.
func (e *EntryRequest) transform() {
	// override if there's any entry number in the lemma
	if newLemma, entryNo, ok := FindEntryNoFromLemma(e.Lemma); ok {
		e.Lemma = newLemma
		e.EntryNo = entryNo
	}
}

type EntryResponse struct {
	kbbi.Lemma
}

type SearchRequest struct {
	Lemma string `form:"entry"`
	Limit uint   `form:"limit" validate:"max=100"`
}

type SearchResponse struct {
	Entries []string `json:"entries"`
}
