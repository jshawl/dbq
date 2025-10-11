package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/searchableviewport"
)

type ResultsPaneModel struct {
	Duration           time.Duration
	Results            db.QueryResult
	Err                error
	SearchableViewport searchableviewport.Model

	focused bool
}

func NewResultsPaneModel() ResultsPaneModel {
	return ResultsPaneModel{
		Duration:           0,
		Results:            db.QueryResult{},
		Err:                nil,
		SearchableViewport: searchableviewport.NewSearchableViewportModel(),

		focused: false,
	}
}

func (model ResultsPaneModel) Update(msg tea.Msg) (ResultsPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, nil
		}
	case QueryResponseReceivedMsg:
		model.Duration = msg.Duration
		model.Err = msg.Err
		model.Results = msg.Results
		model.SearchableViewport.SetContent(model.ResultsView())

		return model, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	model.SearchableViewport, cmd = model.SearchableViewport.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
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

func (model ResultsPaneModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		model.SearchableViewport.View(),
		model.footerView(),
	)
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
	if model.SearchableViewport.FooterView() != "" && model.focused {
		return model.SearchableViewport.FooterView()
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
