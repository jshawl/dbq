package history

import (
	"context"
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	Value  string
	cursor int64
	db     *sql.DB
}

func Init(path string) Model {
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
		cursor: 0,
		db:     database,
		Value:  "",
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
	Entry string
}

type PushedMsg struct {
	id int64
}

type TravelMsg struct {
	Direction string
}

type TraveledMsg struct {
	cursor int64
	value  string
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	//nolint:exhaustive
	switch msg := msg.(type) {
	case PushMsg:
		entry := msg.Entry

		return model, model.push(entry)
	case PushedMsg:
		cursor := msg.id
		model.cursor = cursor

		return model, nil
	case TravelMsg:
		return model, model.travel(msg.Direction)

	case TraveledMsg:
		model.cursor = msg.cursor
		model.Value = msg.value

		return model, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			return model, func() tea.Msg { return TravelMsg{Direction: "previous"} }
		case tea.KeyDown:
			return model, func() tea.Msg { return TravelMsg{Direction: "next"} }
		}
	}

	return model, nil
}

func (model Model) push(entry string) tea.Cmd {
	return func() tea.Msg {
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

		result, _ := stmt.ExecContext(context.Background(), entry)

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

		return PushedMsg{
			id: lastInsertId,
		}
	}
}

func (model Model) travel(direction string) tea.Cmd {
	return func() tea.Msg {
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
			if model.cursor == 0 {
				sql = "select id, query from history where 0 = (?) order by id desc limit 1;"
			} else {
				sql = "select id, query from history where id < (?) order by id desc limit 1;"
			}
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
			return TraveledMsg{
				cursor: model.cursor,
				value:  model.Value,
			}
		}

		return TraveledMsg{
			cursor: id,
			value:  query,
		}
	}
}
