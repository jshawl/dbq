package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewportModel struct {
	Height    int
	Width     int
	YPosition int

	ready    bool
	viewport viewport.Model
	focused  bool
}

func (model ViewportModel) Update(msg tea.Msg) (ViewportModel, tea.Cmd) {
	//nolint:gocritic
	switch msg.(type) {
	case tea.KeyMsg:
		if !model.focused {
			return model, nil
		}
	}

	m, cmd := model.viewport.Update(msg)
	model.viewport = m

	return model, cmd
}

func (model ViewportModel) Resize(width int, height int, yposition int) ViewportModel {
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

func (model ViewportModel) Focus() ViewportModel {
	model.focused = true

	return model
}

func (model ViewportModel) Focused() bool {
	return model.focused
}

func (model ViewportModel) Blur() ViewportModel {
	model.focused = false

	return model
}

func (model ViewportModel) New(width int, height int) viewport.Model {
	return viewport.New(width, height)
}

func (model ViewportModel) SetContent(content string) ViewportModel {
	model.viewport.SetContent(content)

	return model
}

func (model ViewportModel) View() string {
	return model.viewport.View()
}
