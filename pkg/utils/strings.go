package utils

import (
	"strings"
	"unicode"
)

func BreakLines(line string, lineLength int) []string {
	start := 0
	lastWhitespaceIdx := 0
	lines := []string{}
	for i, r := range line {
		if unicode.IsSpace(r) {
			if lastWhitespaceIdx == i-1 && start == lastWhitespaceIdx {
				start = i
			}
			lastWhitespaceIdx = i
		}
		if i-start > lineLength && lastWhitespaceIdx > start {
			lines = append(lines, strings.TrimSpace(line[start:lastWhitespaceIdx]))
			start = lastWhitespaceIdx
		}
	}
	lines = append(lines, strings.TrimSpace(line[start:]))
	return lines
}
