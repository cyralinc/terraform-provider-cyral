package cyral

import (
	"testing"
)

func TestElementsMatch(t *testing.T) {
	testCases := []struct {
		desc        string
		this, other []string
		expectMatch bool
	}{
		{
			desc:        "empty lists",
			this:        []string{},
			other:       []string{},
			expectMatch: true,
		},
		{
			desc:        "lists with different size",
			this:        []string{"1"},
			other:       []string{"2", "3"},
			expectMatch: false,
		},
		{
			desc:        "lists with the same size but different",
			this:        []string{"1", "2"},
			other:       []string{"2", "3"},
			expectMatch: false,
		},
		{
			desc:        "equal lists with shuffled elements",
			this:        []string{"1", "3", "2"},
			other:       []string{"3", "1", "2"},
			expectMatch: true,
		},
	}

	for _, testCase := range testCases {
		match := elementsMatch(testCase.this, testCase.other)
		if match != testCase.expectMatch {
			t.Errorf("For test %q, expected match=%t got match=%t",
				testCase.desc, testCase.expectMatch, match)
		}
	}

}
