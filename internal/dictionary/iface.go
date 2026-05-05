package dictionary

type WOTDRepo interface {
	RandomLemmaIndex() int
	TodayLemmaIndex() int
}

type DictionaryRepo interface {
	Lemma(lemma string, entryNo int) (Lemma, error)
	RandomLemma() Lemma
	LemmaOfTheDay() (Lemma, error)
	Search(prefix string, limit uint) []Lemma
}
