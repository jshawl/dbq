package ui

import (
	"context"
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/db"
)

type Model struct {
	Results     db.QueryResult
	Err         error
	DB          *db.DB
	ResultsPane ResultsPaneModel
	QueryPane   QueryPaneModel
}

type DBMsg struct {
	DB *db.DB
}

type QueryMsg struct {
	Duration time.Duration
	Err      error
	Results  db.QueryResult
	Query    string
}

func Run() {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func InitialModel() Model {
	return Model{
		DB:      nil,
		Err:     nil,
		Results: db.QueryResult{},
		//nolint:exhaustruct
		ResultsPane: ResultsPaneModel{},
		//nolint:exhaustruct
		QueryPane: QueryPaneModel{}.New(),
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

func query(sql string, db *db.DB) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		results := db.Query(ctx, sql)

		return QueryMsg{
			Err:      results.Err,
			Results:  results.Results,
			Duration: results.Duration,
			Query:    sql,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.ResultsPane, cmd = m.ResultsPane.Update(msg)
	cmds = append(cmds, cmd)
	m.QueryPane, cmd = m.QueryPane.Update(msg)
	cmds = append(cmds, cmd)

	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			return m.cycleFocus(), nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.ResultsPane = m.ResultsPane.Resize(msg.Width, msg.Height, lipgloss.Height(m.QueryPane.View()))
	case QueryExecMsg:
		return m, query(msg.Value, m.DB)
	case QueryMsg:
		var cmd tea.Cmd

		if msg.Err == nil {
			m.QueryPane = m.QueryPane.Blur()
			m.ResultsPane = m.ResultsPane.Focus()
		}

		return m, cmd
	case DBMsg:
		m.DB = msg.DB

		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf(
			"%s\n%s",
			m.QueryPane.View(),
			m.Err.Error(),
		)
	}

	return fmt.Sprintf(
		"%s\n%s",
		withFocusView(m.QueryPane.View(), m.QueryPane.Focused()),
		withFocusView(m.ResultsPane.View(), m.ResultsPane.Focused()),
	)
}

func (m Model) cycleFocus() Model {
	if m.QueryPane.Focused() {
		m.QueryPane = m.QueryPane.Blur()
		m.ResultsPane = m.ResultsPane.Focus()
	} else {
		m.QueryPane = m.QueryPane.Focus()
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
