package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBreakLines(tt *testing.T) {
	for _, test := range []struct {
		name     string
		line     string
		length   int
		expected []string
	}{
		{
			name:   "short",
			line:   "This is a line",
			length: 20,
			expected: []string{
				"This is a line",
			},
		},
		{
			name:   "long",
			line:   "This is a long line",
			length: 10,
			expected: []string{
				"This is a",
				"long line",
			},
		},
		{
			name:   "word size exceeds line length",
			line:   "This line has a long word: abcdefghijklmnopqrstuvwxyz",
			length: 20,
			expected: []string{
				"This line has a long",
				"word:",
				"abcdefghijklmnopqrstuvwxyz",
			},
		},
		{
			name:   "word size exceeds line length but in the middle",
			line:   "This is a line. abcdefghijklmnopqrstuvwxyz is a really long word.",
			length: 20,
			expected: []string{
				"This is a line.",
				"abcdefghijklmnopqrstuvwxyz",
				"is a really long",
				"word.",
			},
		},
		{
			name:   "lots o spaces",
			line:   "This  line    has a          lot  of   extra    whitespace.",
			length: 20,
			expected: []string{
				"This  line    has a",
				"lot  of   extra",
				"whitespace.",
			},
		},
		{
			name:     "start with a one-character word",
			line:     "A line with a one-character word",
			length:   100,
			expected: []string{"A line with a one-character word"},
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			assert := require.New(t)
			lines := BreakLines(test.line, test.length)
			assert.Equal(test.expected, lines)
		})
	}
}
