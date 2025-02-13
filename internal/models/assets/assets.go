package assets

import "github.com/raf555/kbbi-api/pkg/kbbi"

type (
	AssetData struct {
		Stats  Stats        `json:"stats"`
		Lemmas []kbbi.Lemma `json:"lemmas"`
	}

	Stats struct {
		Edition    string `json:"edition"`
		EntryCount int    `json:"entryCount"`
		LemmaCount int    `json:"lemmaCount"`
	}
)
