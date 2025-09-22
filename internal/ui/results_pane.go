package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/db"
)

type ResultsPaneModel struct {
	Height    int
	Width     int
	YPosition int
	Results   db.DBQueryResult
	Err       error

	ready    bool
	viewport viewport.Model
	focused  bool
}

func (model ResultsPaneModel) Init() tea.Cmd {
	return nil
}

func (model ResultsPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, nil
		}
	case QueryMsg:
		model.Results = msg.Results
		model.Err = msg.Err
		model.viewport.SetContent(model.ResultsView())
	}

	m, cmd := model.viewport.Update(msg)
	model.viewport = m

	return model, cmd
}

func (model ResultsPaneModel) Resize(width int, height int, yposition int) ResultsPaneModel {
	footerHeight := lipgloss.Height(model.footerView())

	height = height - footerHeight - yposition
	if !model.ready {
		model.viewport = model.New(width, height)
		model.viewport.YPosition = yposition
		model.ready = true
	} else {
		model.viewport.Width = width
		model.viewport.Height = height
	}

	model.Width = model.viewport.Width
	model.Height = model.viewport.Height

	return model
}

func (model ResultsPaneModel) Focus() ResultsPaneModel {
	model.focused = true

	return model
}

func (model ResultsPaneModel) Focused() bool {
	return model.focused
}

func (model ResultsPaneModel) Blur() ResultsPaneModel {
	model.focused = false

	return model
}

func (model ResultsPaneModel) New(width int, height int) viewport.Model {
	return viewport.New(width, height)
}

func (model ResultsPaneModel) SetContent(content string) ResultsPaneModel {
	model.viewport.SetContent(content)

	return model
}

func (model ResultsPaneModel) View() string {
	return fmt.Sprintf("%s\n%s", model.viewport.View(), model.footerView())
}

func (model ResultsPaneModel) ResultsView() string {
	if model.Err != nil {
		return model.Err.Error()
	}

	var builder strings.Builder

	for row := range model.Results.Results {
		builder.WriteString("---\n")

		keys := make([]string, 0, len(model.Results.Results[row]))
		for key := range model.Results.Results[row] {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			builder.WriteString(fmt.Sprintf("%s: %v\n", key, model.Results.Results[row][key]))
		}
	}

	return builder.String()
}

func (model ResultsPaneModel) footerView() string {
	if model.Results.Duration.Seconds() == 0 {
		return ""
	}

	numStr := "1 row"

	numResults := len(model.Results.Results)
	if numResults != 1 {
		numStr = fmt.Sprintf("%d rows", numResults)
	}

	return fmt.Sprintf("(%s in %.3fs)", numStr, model.Results.Duration.Seconds())
}
