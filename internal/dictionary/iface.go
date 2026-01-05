package dictionary

import "github.com/raf555/kbbi-api/pkg/kbbi"

type WOTDRepo interface {
	RandomLemmaIndex() int
	TodayLemmaIndex() int
}

type DictionaryRepo interface {
	Lemma(lemma string, entryNo int) (kbbi.Lemma, error)
	RandomLemma() kbbi.Lemma
	LemmaOfTheDay() (kbbi.Lemma, error)
}
