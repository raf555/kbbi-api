package util

import (
	"strconv"
	"unicode"
)

// FindEntryNoFromLemma will return the cleaned lemma from the entry number
// and will return the corresponding entry number if any.
// ok will be true for above case.
//
// e.g. Apel (2) will return (Apel, 2, true)
func FindEntryNoFromLemma(lemma string) (string, int, bool) {
	length := len(lemma)

	if length < 5 || lemma[length-1] != ')' { // The lemma must contain at least `()`, a digit, a space, and a letter
		return "", 0, false
	}

	digitStartIdx := -1
	digitEndIdx := -1
	openingParenthesisIdx := -1

	// Scan from the end to find `(digits)`
	for i := length - 2; i >= 0; i-- {
		char := rune(lemma[i])

		if unicode.IsDigit(char) {
			if digitEndIdx == -1 {
				digitEndIdx = i + 1
			}
			digitStartIdx = i
		} else if char == '(' {
			if digitStartIdx == -1 {
				return "", 0, false // No digits inside `()`
			}
			openingParenthesisIdx = i
			break
		} else {
			return "", 0, false
		}
	}

	whiteSpaceLoc := openingParenthesisIdx - 1 // before `(`
	if lemma[whiteSpaceLoc] != ' ' {
		return "", 0, false
	}

	digitStr := lemma[digitStartIdx:digitEndIdx]
	entryNo, err := strconv.Atoi(digitStr)
	if err != nil {
		return "", 0, false
	}

	lemmaText := lemma[:whiteSpaceLoc]

	return lemmaText, entryNo, true
}
