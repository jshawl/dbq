package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/history"
)

type Model struct {
	TextInput   textinput.Model
	Query       string
	Results     db.DBQueryResult
	Err         error
	DB          *db.DB
	History     history.Model
	ResultsPane ResultsPaneModel
}

type DBMsg struct {
	DB *db.DB
}

type QueryMsg struct {
	Err     error
	Results db.DBQueryResult
}

func Run() {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func InitialModel() Model {
	input := textinput.New()
	input.Placeholder = "SELECT * FROM users LIMIT 1;"
	input.Focus()
	input.CharLimit = 256
	input.Width = 80

	return Model{
		DB:        nil,
		Err:       nil,
		Query:     "",
		TextInput: input,
		History:   history.Init("/tmp/.dbqhistory"),
		Results: db.DBQueryResult{
			Duration: 0,
			Results:  db.QueryResult{},
		},
		//nolint:exhaustruct
		ResultsPane: ResultsPaneModel{},
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		pgdb, _ := db.NewPostgresDB(ctx, "postgres://admin:password@localhost:5432/dbq_test")
		db := db.NewDB(pgdb)

		return DBMsg{
			DB: db,
		}
	}
}

func query(q string, db *db.DB) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		results, err := db.Query(ctx, q)

		return QueryMsg{
			Err:     err,
			Results: results,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd         tea.Cmd
		cmds        []tea.Cmd
		resultsPane tea.Model
	)

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)
	m.History, cmd = m.History.Update(msg)
	cmds = append(cmds, cmd)
	resultsPane, cmd = m.ResultsPane.Update(msg)
	m.ResultsPane = resultsPane.(ResultsPaneModel)

	cmds = append(cmds, cmd)

	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Query = m.TextInput.Value()

			return m, query(m.Query, m.DB)
		case tea.KeyTab:
			return m.cycleFocus(), nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.History.Cleanup()

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.ResultsPane = m.ResultsPane.Resize(msg.Width, msg.Height, lipgloss.Height(m.TextInput.View()))

	case QueryMsg:
		var cmd tea.Cmd

		if msg.Err == nil {
			m.History, cmd = m.History.Update(history.PushMsg{Entry: m.Query})
			m.TextInput.Blur()
			m.ResultsPane = m.ResultsPane.Focus()
		}

		return m, cmd
	case DBMsg:
		m.DB = msg.DB

		return m, nil
	case history.TraveledMsg:
		m.TextInput.SetValue(m.History.Value)
		m.TextInput.SetCursor(len(m.History.Value))

		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf(
			"%s\n%s\n%s",
			m.TextInput.View(),
			m.Query,
			m.Err.Error(),
		)
	}

	return fmt.Sprintf(
		"%s\n%s",
		withFocusView(m.TextInput.View(), m.TextInput.Focused()),
		withFocusView(m.ResultsPane.View(), m.ResultsPane.Focused()),
	)
}

func (m Model) cycleFocus() Model {
	if m.TextInput.Focused() {
		m.TextInput.Blur()
		m.ResultsPane = m.ResultsPane.Focus()
	} else {
		m.TextInput.Focus()
		m.ResultsPane = m.ResultsPane.Blur()
	}

	return m
}

func withFocusView(view string, focused bool) string {
	if !focused {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#bbb"))

		return style.Render(view)
	}

	return view
}
