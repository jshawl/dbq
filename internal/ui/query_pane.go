package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/history"
)

type QueryPaneModel struct {
	History   history.Model
	TextInput textinput.Model

	focused bool
}

type QueryExecMsg struct {
	Value string
}

func (model QueryPaneModel) New() QueryPaneModel {
	input := textinput.New()
	input.Placeholder = "SELECT * FROM users LIMIT 1;"
	input.Focus()
	input.CharLimit = 256
	input.Width = 80
	input.Cursor.SetMode(1)
	model.TextInput = input
	model.History = history.Init("/tmp/.dbqhistory")
	model.focused = true

	return model
}

func (model QueryPaneModel) Update(msg tea.Msg) (QueryPaneModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, nil
		}

		if msg.Type == tea.KeyEnter {
			return model, dispatch(QueryExecMsg{
				Value: model.TextInput.Value(),
			})
		}
	case history.SetInputValueMsg:
		model.TextInput.SetValue(msg.Value)
		model.TextInput.SetCursor(len(msg.Value))

		return model, nil
	case QueryResponseReceivedMsg:
		if msg.Err != nil {
			return model, nil
		}

		return model, dispatch(history.PushMsg{Query: msg.Query})
	}

	model.History, cmd = model.History.Update(msg)
	cmds = append(cmds, cmd)

	model.TextInput, cmd = model.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model QueryPaneModel) Focused() bool {
	return model.focused
}

func (model QueryPaneModel) Focus() QueryPaneModel {
	model.focused = true
	model.TextInput.Focus()

	return model
}

func (model QueryPaneModel) Blur() QueryPaneModel {
	model.focused = false
	model.TextInput.Blur()

	return model
}

func (model QueryPaneModel) View() string {
	return model.TextInput.View()
}
