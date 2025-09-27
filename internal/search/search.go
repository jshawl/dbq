package search

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SearchMatch struct {
	Start int
	End   int
}

func Search(str string, substring string) []SearchMatch {
	estimatedMatches := len(str) / len(substring)
	if estimatedMatches < 1 {
		estimatedMatches = 1
	}

	matches := make([]SearchMatch, 0, estimatedMatches)

	var builder strings.Builder

	runes := []rune(substring)
	for i, r := range runes {
		builder.WriteString(regexp.QuoteMeta(string(r)))

		if i != len(runes)-1 {
			// optional line breaks in the substring
			builder.WriteString("(?:\r\n|\r|\n)?")
		}
	}

	re := regexp.MustCompile(builder.String())

	for _, match := range re.FindAllStringIndex(str, -1) {
		matches = append(matches, SearchMatch{
			Start: match[0],
			End:   match[1],
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
		builder.WriteString(str[start:match.Start])
		builder.WriteString(style.Render(str[match.Start:match.End]))
		start = match.End
	}

	builder.WriteString(str[start:])
	result := builder.String()

	return result
}
