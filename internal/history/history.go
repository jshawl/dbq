package history

import (
	"context"
	"database/sql"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	cursor int64
	db     *sql.DB
}

func NewHistoryModel(path string) Model {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
		create table if not exists history (
			id integer not null primary key,
			query text,
			created_at datetime default current_timestamp
		);
	`

	_, err = database.ExecContext(context.Background(), sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return Model{
		cursor: math.MaxInt32,
		db:     database,
	}
}

func (h Model) Cleanup() {
	defer func() {
		err := h.db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

type PushMsg struct {
	Query string
}

type pushedMsg struct {
	id int64
}

type traveledMsg struct {
	cursor int64
	query  string
}

type SetInputValueMsg struct {
	Value string
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	//nolint:exhaustive
	switch msg := msg.(type) {
	case PushMsg:
		query := msg.Query

		return model, model.push(query)
	case pushedMsg:
		cursor := msg.id
		model.cursor = cursor

		return model, nil
	case traveledMsg:
		model.cursor = msg.cursor

		return model, model.dispatch(SetInputValueMsg{Value: msg.query})
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			return model, model.travelCmd("previous")
		case tea.KeyDown:
			return model, model.travelCmd("next")
		}
	}

	return model, nil
}

func (model Model) SetCursor(cursor int64) Model {
	model.cursor = cursor

	return model
}

func (model Model) Previous() (int64, string) {
	return model.travel("previous")
}

func (model Model) Next() (int64, string) {
	return model.travel("next")
}

func (model Model) Push(query string) int64 {
	transaction, err := model.db.BeginTx(
		context.Background(),
		&sql.TxOptions{ReadOnly: false, Isolation: 0},
	)
	if err != nil {
		log.Fatal("db begin err")
	}

	stmt, err := transaction.PrepareContext(
		context.Background(),
		"insert into history (query) values (?)",
	)
	if err != nil {
		log.Fatal("prepare err")
	}

	result, _ := stmt.ExecContext(context.Background(), query)

	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("exec err")
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return lastInsertId
}

func (model Model) dispatch(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (model Model) push(query string) tea.Cmd {
	return func() tea.Msg {
		return pushedMsg{
			id: model.Push(query),
		}
	}
}

func (model Model) travel(direction string) (int64, string) {
	_, err := model.db.BeginTx(
		context.Background(),
		&sql.TxOptions{ReadOnly: true, Isolation: 0},
	)
	if err != nil {
		log.Fatal("db begin err")
	}

	var sql string
	if direction == "next" {
		sql = "select id, query from history where id > (?) order by id asc limit 1;"
	}

	if direction == "previous" {
		sql = "select id, query from history where id < (?) order by id desc limit 1;"
	}

	stmt, _ := model.db.PrepareContext(
		context.Background(),
		sql,
	)

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var query string

	//nolint:varnamelen
	var id int64

	err = stmt.QueryRowContext(context.Background(), model.cursor).Scan(&id, &query)
	if err != nil {
		cursor := model.cursor
		if direction == "next" {
			cursor = math.MaxInt32
		}

		if direction == "previous" {
			cursor = 0
		}

		return cursor, query
	}

	log.Printf("history.Traveled id: %d  query: %s", id, query)

	return id, query
}

func (model Model) travelCmd(direction string) tea.Cmd {
	return func() tea.Msg {
		cursor, query := model.travel(direction)

		return traveledMsg{
			cursor: cursor,
			query:  query,
		}
	}
}
