package ui

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/history"
)

type Model struct {
	TextInput textinput.Model
	Query     string
	Results   db.DBQueryResult
	Err       error
	DB        *db.DB
	History   history.Model
	viewport  viewport.Model
	ready     bool
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
		DB:    nil,
		Err:   nil,
		Query: "",
		Results: db.DBQueryResult{
			Results:  db.QueryResult{},
			Duration: time.Duration(0),
		},
		TextInput: input,
		History:   history.Init("/tmp/.dbqhistory"),
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
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)
	m.History, cmd = m.History.Update(msg)
	cmds = append(cmds, cmd)

	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Query = m.TextInput.Value()

			return m, query(m.Query, m.DB)
		case tea.KeyCtrlC, tea.KeyEsc:
			m.History.Cleanup()

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.TextInput.View())
		footerHeight := lipgloss.Height("after")
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case QueryMsg:
		var cmd tea.Cmd

		m.Results = msg.Results
		m.Err = msg.Err

		if msg.Err == nil {
			m.History, cmd = m.History.Update(history.PushMsg{Entry: m.Query})
			m.viewport.SetContent(m.resultsView())
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

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			m.TextInput.View(),
			m.Query,
			m.durationView(),
			m.Err.Error(),
		)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.TextInput.View(), m.viewport.View(), "after viewport")
}

func (m Model) durationView() string {
	if m.Results.Duration.Seconds() == 0 {
		return ""
	}

	numStr := "1 row"

	numResults := len(m.Results.Results)
	if numResults != 1 {
		numStr = fmt.Sprintf("%d rows", numResults)
	}

	return fmt.Sprintf("(%s in %.3fs)\n", numStr, m.Results.Duration.Seconds())
}

func (m Model) resultsView() string {
	jsonStr := ""

	if len(m.Results.Results) > 0 {
		jsonData, err := json.MarshalIndent(m.Results.Results, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		jsonStr = string(jsonData)
	}

	return fmt.Sprintf("%s\n%s", m.durationView(), jsonStr)
}
