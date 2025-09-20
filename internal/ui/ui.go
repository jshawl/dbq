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
	db "github.com/jshawl/dbq/internal/db"
)

type Model struct {
	TextInput textinput.Model
	Query     string
	Results   db.DBQueryResult
	Err       error
	DB        *db.DB
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
	//nolint:exhaustive
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Query = m.TextInput.Value()

			msg := query(m.Query, m.DB)

			return m, msg
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			var cmd tea.Cmd

			m.TextInput, cmd = m.TextInput.Update(msg)

			return m, cmd
		}

	case QueryMsg:
		m.Results = msg.Results
		m.Err = msg.Err

		return m, nil
	case DBMsg:
		m.DB = msg.DB

		return m, nil
	default:
		return m, nil
	}
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.promptView(), m.resultsView())
}

func (m Model) promptView() string {
	return m.TextInput.View() + "\n"
}

func (m Model) resultsView() string {
	durationStr := ""
	if m.Results.Duration.Seconds() > 0 {
		durationStr = fmt.Sprintf("%.3fs\n", m.Results.Duration.Seconds())
	}

	if m.Err != nil {
		return fmt.Sprintf("%s\n%s\n%s", m.Query, durationStr, m.Err.Error())
	}

	jsonStr := ""

	if len(m.Results.Results) > 0 {
		jsonData, err := json.MarshalIndent(m.Results.Results, "", "  ")
		if err != nil {
			panic(err)
		}

		jsonStr = string(jsonData)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.Query, durationStr, jsonStr)
}
