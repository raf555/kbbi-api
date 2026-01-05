package dictionary

import "errors"

var (
	ErrLemmaNotFound         = errors.New("dictionary: lemma not found")
	ErrLemmaTooLong          = errors.New("dictionary: lemma length too long")
	ErrEntryNotFound         = errors.New("dictionary: entry not found")
	ErrUnexpectedEmptyLemma  = errors.New("dictionary: unexpected empty lemma")
	ErrUnexpectedEntryNumber = errors.New("dictionary: unexpected entry number")
	ErrUnexpectedWotdIndex   = errors.New("dictionary: unexpected wotd lemma index")
)
