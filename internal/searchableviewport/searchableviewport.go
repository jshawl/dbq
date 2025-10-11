package searchableviewport

import tea "github.com/charmbracelet/bubbletea"

type Model struct{}

func NewSearchableViewportModel() Model {
	return Model{}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return model, nil
}

func (model Model) View() string {
	return ""
}
