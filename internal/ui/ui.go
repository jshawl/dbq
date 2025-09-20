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
	tea "github.com/charmbracelet/bubbletea"
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
}

type DBMsg struct {
	DB *db.DB
}

type QueryMsg struct {
	Err     error
	Results db.DBQueryResult
}

func Run() {
	p := tea.NewProgram(InitialModel())

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
		History:   history.Init("/Users/jesse/.dbqhistory"),
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

//nolint:ireturn
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

	case QueryMsg:
		var cmd tea.Cmd

		m.Results = msg.Results
		m.Err = msg.Err

		if msg.Err == nil {
			m.History, cmd = m.History.Update(history.PushMsg{Entry: m.Query})
		}

		return m, cmd
	case DBMsg:
		m.DB = msg.DB

		return m, nil
	case history.TraveledMsg:
		var cmd tea.Cmd

		m.TextInput.SetValue(m.History.Value)
		m.TextInput.SetCursor(len(m.History.Value))

		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			m.promptView(),
			m.Query,
			m.durationView(),
			m.Err.Error(),
		)
	}

	return fmt.Sprintf("%s\n%s", m.promptView(), m.resultsView())
}

func (m Model) promptView() string {
	return m.TextInput.View() + "\n"
}

func (m Model) durationView() string {
	if m.Results.Duration.Seconds() == 0 {
		return ""
	}

	return fmt.Sprintf("%.3fs\n", m.Results.Duration.Seconds())
}

func (m Model) resultsView() string {
	jsonStr := ""

	if len(m.Results.Results) > 0 {
		jsonData, err := json.MarshalIndent(m.Results.Results, "", "  ")
		if err != nil {
			panic(err)
		}

		jsonStr = string(jsonData)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.Query, m.durationView(), jsonStr)
}
