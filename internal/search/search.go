package search

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchMatch struct {
	BufferStart int
	BufferEnd   int
}

type SearchMsg struct {
	Value string
}

type Model struct {
	focused   bool
	textInput textinput.Model
}

func NewSearchModel() Model {
	textInput := textinput.New()
	textInput.Prompt = "/"
	textInput.Focus()
	textInput.Cursor.SetMode(1)

	return Model{
		focused:   false,
		textInput: textInput,
	}
}

func (model Model) View() string {
	return model.textInput.View()
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
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
		case "enter":
			value := model.textInput.Value()
			model.textInput.Blur()
			// TODO focused == false allows scrolling the viewport,
			// but then the search term is no longer visible
			// model.focused = false
			return model, func() tea.Msg {
				return SearchMsg{
					Value: value,
				}
			}
		}
	}

	model.textInput, cmd = model.textInput.Update(msg)

	return model, cmd
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
