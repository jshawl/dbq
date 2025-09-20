package ui_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/jshawl/dbq/internal/db"
	ui "github.com/jshawl/dbq/internal/ui"
)

var errSQL = errors.New("sql error")

//nolint:ireturn
func assertMsgType[T interface{}](t *testing.T, cmd tea.Cmd) T {
	t.Helper()

	msg := cmd()
	typed, ok := msg.(T)

	if !ok {
		t.Fatalf("Expected msg to be of type %T, got %T", *new(T), msg)
	}

	return typed
}

//nolint:ireturn
func assertModelType[T tea.Model](t *testing.T, model tea.Model) T {
	t.Helper()

	typed, ok := model.(T)
	if !ok {
		t.Fatalf("Expected msg to be of type %T, got %T", *new(T), model)
	}

	return typed
}

func setupDatabaseModel(t *testing.T) ui.Model {
	t.Helper()

	model := ui.InitialModel()
	model.TextInput.SetValue("SELECT * FROM users LIMIT 1;")
	cmd := model.Init()
	msg := assertMsgType[ui.DBMsg](t, cmd)
	updatedModel, _ := model.Update(msg)

	typedModel, ok := updatedModel.(ui.Model)
	if !ok {
		t.Fatal("expected updated model of type Model")
	}

	t.Cleanup(func() {
		err := typedModel.DB.Close(t.Context())
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	return typedModel
}

func makeKeyMsg(key tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{
		Alt:   false,
		Paste: false,
		Runes: nil,
		Type:  key,
	}
}

func makeResults(duration time.Duration, userID int) db.DBQueryResult {
	return db.DBQueryResult{
		Duration: duration,
		Results: db.QueryResult{
			map[string]interface{}{
				"id": userID,
			},
		},
	}
}

func TestInitialModel(t *testing.T) {
	t.Parallel()

	model := ui.InitialModel()

	if model.TextInput.Placeholder != "SELECT * FROM users LIMIT 1;" {
		t.Fatal("expected placeholder to be a select statement")
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	model := ui.InitialModel()
	cmd := model.Init()

	msg := assertMsgType[ui.DBMsg](t, cmd)
	if msg.DB == nil {
		t.Fatal("expected DBmsg to contain db")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("keys - enter", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(makeKeyMsg(tea.KeyEnter))

		queryMsg := assertMsgType[ui.QueryMsg](t, cmd)
		if len(queryMsg.Results.Results) == 0 {
			t.Fatal("expected results")
		}
	})

	t.Run("keys - ctrl-c", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(makeKeyMsg(tea.KeyCtrlC))

		assertMsgType[tea.QuitMsg](t, cmd)
	})

	t.Run("keys - other", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(makeKeyMsg(tea.KeySpace))

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})

	t.Run("QueryMsg", func(t *testing.T) {
		t.Parallel()

		userID := 789
		model := setupDatabaseModel(t)
		updatedModel, cmd := model.Update(ui.QueryMsg{
			Err:     nil,
			Results: makeResults(0, userID),
		})

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}

		typedModel := assertModelType[ui.Model](t, updatedModel)

		got := typedModel.Results.Results[0]["id"]
		if got != userID {
			t.Fatalf("expected first result to have id %d got %d", userID, got)
		}
	})

	t.Run("unknown msg", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)

		_, cmd := model.Update(nil)
		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})
}

func TestView(t *testing.T) {
	t.Parallel()

	t.Run("duration", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		model.Results = makeResults(time.Millisecond*2345, 123)

		view := model.View()
		if !strings.Contains(view, "2.345s") {
			t.Fatalf("expected model error to be visible\n %s", view)
		}
	})

	t.Run("errors", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		model.Err = errSQL

		view := model.View()
		if !strings.Contains(view, "sql error") {
			t.Fatal("expected model error to be visible")
		}
	})

	t.Run("results", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		model.Results = makeResults(0, 123)

		view := model.View()
		if !strings.Contains(view, "\"id\": 123") {
			t.Fatalf("expected results to be visible, got \n %s", view)
		}
	})
}
