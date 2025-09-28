package search

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SearchMatch struct {
	BufferStart     int
	BufferEnd       int
	ScreenYPosition int
}

func Search(str string, substring string) []SearchMatch {
	estimatedMatches := len(str) / len(substring)
	if estimatedMatches < 1 {
		estimatedMatches = 1
	}

	matches := make([]SearchMatch, 0, estimatedMatches)

	re := regexp.MustCompile(substring)

	newLineRe := regexp.MustCompile("\n")

	for _, match := range re.FindAllStringIndex(str, -1) {
		matches = append(matches, SearchMatch{
			BufferStart:     match[0],
			BufferEnd:       match[1],
			ScreenYPosition: len(newLineRe.FindAllStringIndex(str, match[0])),
		})
	}

	return matches
}

func Highlight(str string, matches []SearchMatch) string {
	var builder strings.Builder

	start := 0

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("11"))

	for _, match := range matches {
		builder.WriteString(str[start:match.BufferStart])
		builder.WriteString(style.Render(str[match.BufferStart:match.BufferEnd]))
		start = match.BufferEnd
	}

	builder.WriteString(str[start:])
	result := builder.String()

	return result
}
