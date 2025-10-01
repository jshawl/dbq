package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/search"
)

type ResultsPaneModel struct {
	Height    int
	Width     int
	YPosition int
	Duration  time.Duration
	Results   db.QueryResult
	Err       error
	Search    search.Model

	ready    bool
	viewport viewport.Model
	focused  bool
}

type WindowSizeMsg struct {
	Height    int
	Width     int
	YPosition int
}

func NewResultsPaneModel() ResultsPaneModel {
	return ResultsPaneModel{
		Height:    0,
		Width:     0,
		YPosition: 0,
		Duration:  0,
		Results:   db.QueryResult{},
		Err:       nil,
		Search:    search.NewSearchModel(),

		ready:    false,
		focused:  false,
		viewport: NewViewportModel(0, 0),
	}
}

func (model ResultsPaneModel) Update(msg tea.Msg) (ResultsPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, nil
		}

		if model.Search.Focused() {
			updatedSearchModel, cmd := model.Search.Update(msg)
			model.Search = updatedSearchModel

			return model, cmd
		}
	case QueryResponseReceivedMsg:
		model.Duration = msg.Duration
		model.Err = msg.Err
		model.Results = msg.Results
		model.viewport.SetContent(model.ResultsView())
		model.viewport.YPosition = 0

		return model, nil
	case WindowSizeMsg:
		model = model.Resize(msg.Width, msg.Height, msg.YPosition)

		return model, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	model.Search, cmd = model.Search.Update(msg)
	cmds = append(cmds, cmd)
	model.viewport, cmd = model.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model ResultsPaneModel) Resize(width int, height int, yposition int) ResultsPaneModel {
	footerHeight := lipgloss.Height(model.footerView())

	height = height - footerHeight - yposition
	if !model.ready {
		model.viewport = NewViewportModel(width, height)
		model.viewport.YPosition = 0
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

func NewViewportModel(width int, height int) viewport.Model {
	return viewport.New(width, height)
}

func (model ResultsPaneModel) View() string {
	return fmt.Sprintf("%s\n%s", model.viewport.View(), model.footerView())
}

func (model ResultsPaneModel) ResultsView() string {
	if model.Err != nil {
		return model.Err.Error()
	}

	var builder strings.Builder

	for row := range model.Results {
		builder.WriteString("---\n")

		keys := make([]string, 0, len(model.Results[row]))
		for key := range model.Results[row] {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			builder.WriteString(fmt.Sprintf("%s: %v\n", key, model.Results[row][key]))
		}
	}

	return builder.String()
}

func (model ResultsPaneModel) footerView() string {
	if model.Search.Focused() {
		return model.Search.View()
	}

	if model.Duration.Seconds() == 0 {
		return ""
	}

	numStr := "1 row"

	numResults := len(model.Results)
	if numResults != 1 {
		numStr = fmt.Sprintf("%d rows", numResults)
	}

	return fmt.Sprintf("(%s in %.3fs)", numStr, model.Duration.Seconds())
}
