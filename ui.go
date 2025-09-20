package main

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
)

type Model struct {
	textInput textinput.Model
	query     string
	results   DBQueryResult
	err       error
	db        *DB
}

type DBMsg struct {
	db *DB
}

type QueryMsg struct {
	err     error
	results DBQueryResult
}

func ui() {
	p := tea.NewProgram(initialModel())

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func initialModel() Model {
	input := textinput.New()
	input.Placeholder = "SELECT * FROM users LIMIT 1;"
	input.Focus()
	input.CharLimit = 256
	input.Width = 80

	return Model{
		db:    nil,
		err:   nil,
		query: "",
		results: DBQueryResult{
			Results:  QueryResult{},
			Duration: time.Duration(0),
		},
		textInput: input,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		pgdb, _ := NewPostgresDB(ctx, "postgres://admin:password@localhost:5432/dbq_test")
		db := NewDB(pgdb)

		return DBMsg{
			db: db,
		}
	}
}

func query(q string, db *DB) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		results, err := db.Query(ctx, q)

		return QueryMsg{
			err:     err,
			results: results,
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
			m.query = m.textInput.Value()

			msg := query(m.query, m.db)

			return m, msg
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			var cmd tea.Cmd

			m.textInput, cmd = m.textInput.Update(msg)

			return m, cmd
		}

	case QueryMsg:
		m.results = msg.results
		m.err = msg.err

		return m, nil
	case DBMsg:
		m.db = msg.db

		return m, nil
	default:
		return m, nil
	}
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.promptView(), m.resultsView())
}

func (m Model) promptView() string {
	return m.textInput.View() + "\n"
}

func (m Model) resultsView() string {
	durationStr := ""
	if m.results.Duration.Seconds() > 0 {
		durationStr = fmt.Sprintf("%.3fs\n", m.results.Duration.Seconds())
	}

	if m.err != nil {
		return fmt.Sprintf("%s\n%s\n%s", m.query, durationStr, m.err.Error())
	}

	jsonStr := ""

	if len(m.results.Results) > 0 {
		jsonData, err := json.MarshalIndent(m.results.Results, "", "  ")
		if err != nil {
			panic(err)
		}

		jsonStr = string(jsonData)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.query, durationStr, jsonStr)
}
