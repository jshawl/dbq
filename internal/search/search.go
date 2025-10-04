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

type SearchClearMsg struct{}

type Model struct {
	Value     string
	focused   bool
	textInput textinput.Model
}

func NewSearchModel() Model {
	textInput := textinput.New()
	textInput.Prompt = "/"
	textInput.Focus()
	textInput.Cursor.SetMode(1)

	return Model{
		Value:     textInput.Value(),
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
			return model.Focus(), nil
		case "esc":
			model.focused = false
			model.textInput.SetValue("")
			model.Value = ""

			return model, func() tea.Msg {
				return SearchClearMsg{}
			}
		case "enter":
			value := model.textInput.Value()
			model.Value = value
			model.textInput.Blur()
			model.focused = false

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
	model.textInput.Focus()

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

	for matchIndex, match := range matches {
		builder.WriteString(str[start:match.BufferStart])

		if matchIndex == 0 {
			firstCharacter := WithBlackBackground(str[match.BufferStart : match.BufferStart+1])
			rest := WithYellowBackground(str[match.BufferStart+1 : match.BufferEnd])
			builder.WriteString(firstCharacter + rest)
		} else {
			builder.WriteString(WithYellowBackground(str[match.BufferStart:match.BufferEnd]))
		}

		start = match.BufferEnd
	}

	builder.WriteString(str[start:])
	result := builder.String()

	return result
}

func WithYellowBackground(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("11")).Render(str)
}

func WithBlackBackground(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#000000")).Render(str)
}
