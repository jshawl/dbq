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

	model.History, cmd = model.History.Update(msg)
	cmds = append(cmds, cmd)
	model.TextInput, cmd = model.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, tea.Batch(cmds...)
		}
		switch msg.Type {
		case tea.KeyEnter:
			return model, func() tea.Msg {
				return QueryExecMsg{
					Value: model.TextInput.Value(),
				}
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			model.History.Cleanup()
		}
	case history.TraveledMsg:
		model.TextInput.SetValue(model.History.Value)
		model.TextInput.SetCursor(len(model.History.Value))

		return model, tea.Batch(cmds...)
	case QueryMsg:
		if msg.Err == nil {
			model.History, cmd = model.History.Update(history.PushMsg{Entry: msg.Query})
			cmds = append(cmds, cmd)
		}
	}

	return model, tea.Batch(cmds...)
}

func (model QueryPaneModel) Focused() bool {
	return model.focused
}

func (model QueryPaneModel) Focus() QueryPaneModel {
	model.focused = true

	return model
}

func (model QueryPaneModel) Blur() QueryPaneModel {
	model.focused = false

	return model
}

func (model QueryPaneModel) View() string {
	return model.TextInput.View()
}
