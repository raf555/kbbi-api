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
	Lemma   string `uri:"entry" validate:"required"`
	EntryNo int    `form:"entryNo" validate:"gte=0"`
}

// transform will look for entry number in the lemma itself.
// if it's present, it'll overwrite the lemma without the number and as well as modify the entryNo.
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
