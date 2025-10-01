package search

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchMatch struct {
	BufferStart int
	BufferEnd   int
}

type Model struct {
	focused bool
}

func NewSearchModel() Model {
	return Model{
		focused: false,
	}
}

func (model Model) View() string {
	return "/" // + input view
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	//nolint:gocritic
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			model.focused = true

			return model, nil
		case "esc":
			model.focused = false

			return model, nil
		}
	}

	return model, nil
}

func (model Model) Focus() Model {
	model.focused = true

	return model
}

func (model Model) Focused() bool {
	return model.focused
}

func Search(str string, substring string) []SearchMatch {
	estimatedMatches := len(str) / len(substring)
	if estimatedMatches < 1 {
		estimatedMatches = 1
	}

	matches := make([]SearchMatch, 0, estimatedMatches)

	re := regexp.MustCompile(substring)

	for _, match := range re.FindAllStringIndex(str, -1) {
		matches = append(matches, SearchMatch{
			BufferStart: match[0],
			BufferEnd:   match[1],
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
