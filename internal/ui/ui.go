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

type QueryResponseReceivedMsg struct {
	QueryMsg
}

func Run() {
	f, _ := tea.LogToFile("debug.log", "debug")

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	p := tea.NewProgram(NewUIModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	log.Println("ui.Run()")

	_, err := p.Run()
	if err != nil {
		defer func() {
			log.Fatal(err)
		}()
	}
}

func NewUIModel() Model {
	return Model{
		DB:      nil,
		Err:     nil,
		Results: db.QueryResult{},
		//nolint:exhaustruct
		ResultsPane: ResultsPaneModel{},
		QueryPane:   NewQueryPaneModel(),
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

func dispatch(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			return m.cycleFocus(), nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.QueryPane.History.Cleanup()

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m, dispatch(WindowSizeMsg{
			Width:     msg.Width,
			Height:    msg.Height,
			YPosition: lipgloss.Height(m.QueryPane.View()),
		})
	case QueryExecMsg:
		return m, query(msg.Value, m.DB)
	case QueryMsg:
		if msg.Err == nil {
			m.QueryPane = m.QueryPane.Blur()
			m.ResultsPane = m.ResultsPane.Focus()
		}

		return m, dispatch(QueryResponseReceivedMsg{
			QueryMsg: msg,
		})
	case DBMsg:
		m.DB = msg.DB

		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.ResultsPane, cmd = m.ResultsPane.Update(msg)
	cmds = append(cmds, cmd)

	m.QueryPane, cmd = m.QueryPane.Update(msg)
	cmds = append(cmds, cmd)

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
