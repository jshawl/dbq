package searchableviewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

//nolint:recvcheck // to match bubbletea interface
type Model struct {
	Height int
	Width  int

	ready    bool
	viewport viewport.Model
}

type WindowSizeMsg struct {
	Height int
	Width  int
}

func NewSearchableViewportModel() Model {
	return Model{
		Height: 0,
		Width:  0,

		ready:    false,
		viewport: viewport.New(0, 0),
	}
}

func (model *Model) SetContent(str string) {
	model.viewport.SetContent(str)
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	const footerHeight = 1
	//nolint:gocritic
	switch msg := msg.(type) {
	case WindowSizeMsg:
		height := msg.Height - footerHeight
		if !model.ready {
			model.viewport = viewport.New(msg.Width, height)
			model.ready = true
		} else {
			model.viewport.Width = msg.Width
			model.viewport.Height = height
		}

		return model, nil
	}

	var cmd tea.Cmd

	model.viewport, cmd = model.viewport.Update(msg)

	return model, cmd
}

func (model Model) View() string {
	return model.viewport.View()
}
