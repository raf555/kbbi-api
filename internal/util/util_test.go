package util

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindEntryNoFromLemma(t *testing.T) {
	tcs := []struct {
		in string

		expectedOk    bool
		expectedLemma string
		expentedEntry int
	}{
		{
			in:         "",
			expectedOk: false,
		},
		{
			in:         "apel (",
			expectedOk: false,
		},
		{
			in:         "apel )",
			expectedOk: false,
		},
		{
			in:         "apel ()",
			expectedOk: false,
		},
		{
			in:         "apel (a)",
			expectedOk: false,
		},
		{
			in:         "apel 1)",
			expectedOk: false,
		},
		{
			in:         "111111)",
			expectedOk: false,
		},
		{
			in:         "(111111)",
			expectedOk: false,
		},
		{
			in:         "apel(1)",
			expectedOk: false,
		},
		{
			// still a valid case since it ends with entry number format.
			// however, the lemma will be empty string.
			in:            " (111111)",
			expectedOk:    true,
			expectedLemma: "",
			expentedEntry: 111111,
		},
		{
			in:            "apel (1)",
			expectedOk:    true,
			expectedLemma: "apel",
			expentedEntry: 1,
		},
		{
			in:            "some lemma (1696969)",
			expectedOk:    true,
			expectedLemma: "some lemma",
			expentedEntry: 1696969,
		},
		{
			in:            "another-one  (1123123)",
			expectedOk:    true,
			expectedLemma: "another-one ",
			expentedEntry: 1123123,
		},
	}

	for _, tc := range tcs {
		log.Println("testing case:", tc.in)

		lemma, entry, ok := FindEntryNoFromLemma(tc.in)

		if !tc.expectedOk {
			assert.False(t, ok)
			assert.Zero(t, entry)
			assert.Zero(t, lemma)
		} else {
			assert.True(t, ok)
			assert.Equal(t, tc.expectedLemma, lemma)
			assert.Equal(t, tc.expentedEntry, entry)
		}
	}
}
