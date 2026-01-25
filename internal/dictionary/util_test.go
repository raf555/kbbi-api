package dictionary_test

import (
	"fmt"
	"testing"

	"github.com/raf555/kbbi-api/internal/dictionary"
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
		t.Run(fmt.Sprintf("input=%s", tc.in), func(t *testing.T) {
			lemma, entry, ok := dictionary.FindEntryNoFromLemma(tc.in)

			if !tc.expectedOk {
				assert.False(t, ok)
				assert.Zero(t, entry)
				assert.Zero(t, lemma)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tc.expectedLemma, lemma)
				assert.Equal(t, tc.expentedEntry, entry)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tcs := []struct {
		in                 string
		removePunctuations bool
		expected           string
	}{
		{
			in:                 "",
			removePunctuations: true,
			expected:           "",
		},
		{
			in:                 "hello",
			removePunctuations: true,
			expected:           "hello",
		},
		{
			in:                 "HELLO",
			removePunctuations: true,
			expected:           "HELLO",
		},
		{
			in:                 "hello world",
			removePunctuations: true,
			expected:           "hello world",
		},
		{
			in:                 "hello123",
			removePunctuations: true,
			expected:           "hello123",
		},
		{
			in:                 "hello-world",
			removePunctuations: true,
			expected:           "helloworld",
		},
		{
			in:                 "hello_world",
			removePunctuations: true,
			expected:           "helloworld",
		},
		{
			in:                 "hello!@#$%^&*()world",
			removePunctuations: true,
			expected:           "helloworld",
		},
		{
			in:                 "café",
			removePunctuations: true,
			expected:           "cafe",
		},
		{
			in:                 "naïve",
			removePunctuations: true,
			expected:           "naive",
		},
		{
			in:                 "façade",
			removePunctuations: true,
			expected:           "facade",
		},
		{
			in:                 "Apél",
			removePunctuations: true,
			expected:           "Apel",
		},
		{
			in:                 "  multiple   spaces  ",
			removePunctuations: true,
			expected:           "  multiple   spaces  ",
		},
		{
			in:                 "hello...world",
			removePunctuations: true,
			expected:           "helloworld",
		},
		{
			in:                 "test@123#abc",
			removePunctuations: true,
			expected:           "test123abc",
		},
		{
			in:                 "ñoño",
			removePunctuations: true,
			expected:           "nono",
		},
		{
			in:                 "中文",
			removePunctuations: true,
			expected:           "",
		},
		{
			in:                 "",
			removePunctuations: false,
			expected:           "",
		},
		{
			in:                 "hello",
			removePunctuations: false,
			expected:           "hello",
		},
		{
			in:                 "hello-world",
			removePunctuations: false,
			expected:           "hello-world",
		},
		{
			in:                 "hello_world",
			removePunctuations: false,
			expected:           "hello_world",
		},
		{
			in:                 "hello!@#$%^&*()world",
			removePunctuations: false,
			expected:           "hello!@#$%^&*()world",
		},
		{
			in:                 "café",
			removePunctuations: false,
			expected:           "cafe",
		},
		{
			in:                 "naïve",
			removePunctuations: false,
			expected:           "naive",
		},
		{
			in:                 "façade",
			removePunctuations: false,
			expected:           "facade",
		},
		{
			in:                 "Apél",
			removePunctuations: false,
			expected:           "Apel",
		},
		{
			in:                 "hello...world",
			removePunctuations: false,
			expected:           "hello...world",
		},
		{
			in:                 "test@123#abc",
			removePunctuations: false,
			expected:           "test@123#abc",
		},
		{
			in:                 "ñoño",
			removePunctuations: false,
			expected:           "nono",
		},
		{
			in:                 "中文",
			removePunctuations: false,
			expected:           "中文",
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("input=%s (removePunctuation=%t)", tc.in, tc.removePunctuations), func(t *testing.T) {
			result := dictionary.Normalize(tc.in, tc.removePunctuations)
			assert.Equal(t, tc.expected, result)
		})
	}
}
