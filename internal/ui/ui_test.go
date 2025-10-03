package ui_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/testutil"
	ui "github.com/jshawl/dbq/internal/ui"
)

var errSQL = errors.New("sql error")

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

	model := ui.NewUIModel("postgres://admin:password@localhost:5432/dbq_test")
	// model.TextInput.SetValue("SELECT * FROM users LIMIT 1;")

	cmd := model.Init()
	msg := testutil.AssertMsgType[ui.DBMsg](t, cmd)
	model = assertModelType[ui.Model](t, model)
	updatedModel, _ := model.Update(msg)
	model = assertModelType[ui.Model](t, updatedModel)

	windowSizeMsg := tea.WindowSizeMsg{Width: 80, Height: 20}
	updatedModel, _ = model.Update(windowSizeMsg)
	model = assertModelType[ui.Model](t, updatedModel)

	t.Cleanup(func() {
		err := model.DB.Close(t.Context())
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	return model
}

func makeResults(userID int, userIDs ...int) db.QueryResult {
	rows := []map[string]interface{}{
		{
			"id":         userID,
			"created_at": "2025-09-21T15:41:22",
		},
	}
	if len(userIDs) > 0 {
		rows = append(rows, map[string]interface{}{
			"id":         userIDs[0],
			"created_at": "2025-09-21T15:41:22",
		})
	}

	return rows
}

func TestInit(t *testing.T) {
	t.Parallel()

	model := ui.NewUIModel("postgres://admin:password@localhost:5432/dbq_test")
	cmd := model.Init()

	msg := testutil.AssertMsgType[ui.DBMsg](t, cmd)
	if msg.DB == nil {
		t.Fatal("expected DBmsg to contain db")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	t.Run("keys - ctrl-c", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyCtrlC))

		testutil.AssertMsgType[tea.QuitMsg](t, cmd)
	})

	t.Run("keys - tab", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		if !model.QueryPane.Focused() {
			t.Fatal("expected query pane to be focused")
		}

		updatedModel, _ := model.Update(testutil.MakeKeyMsg(tea.KeyTab))
		model = assertModelType[ui.Model](t, updatedModel)

		if model.QueryPane.Focused() {
			t.Fatal("expected query pane not to be focused")
		}

		updatedModel, _ = model.Update(testutil.MakeKeyMsg(tea.KeyTab))
		model = assertModelType[ui.Model](t, updatedModel)

		if !model.QueryPane.Focused() {
			t.Fatal("expected query pane to be focused")
		}
	})

	t.Run("keys - other", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeySpace))

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})

	t.Run("QueryExecMsg", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(ui.QueryExecMsg{
			Value: "select * from users where id = 1",
		})

		msg := testutil.AssertMsgType[ui.QueryMsg](t, cmd)
		if len(msg.Results) != 1 {
			t.Fatal("expected QueryExecMsg to query actual db")
		}
	})

	t.Run("QueryMsg", func(t *testing.T) {
		t.Parallel()

		userID := 789
		model := setupDatabaseModel(t)
		updatedModel, _ := model.Update(ui.QueryMsg{
			Duration: 0,
			Err:      nil,
			Results:  makeResults(0, userID),
			Query:    "select * from users where userID = 789",
		})

		typedModel := assertModelType[ui.Model](t, updatedModel)

		if !typedModel.ResultsPane.Focused() {
			t.Fatal("expected requery results to focus on viewport")
		}
	})

	t.Run("QueryMsg - err", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		updatedModel, _ := model.Update(ui.QueryMsg{
			Duration: 0,
			Err:      errSQL,
			Results:  db.QueryResult{},
			Query:    "not sql",
		})

		typedModel := assertModelType[ui.Model](t, updatedModel)

		if typedModel.ResultsPane.Focused() {
			t.Fatal("expected requery results not to focus on viewport")
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

	t.Run("QueryPane input visible", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)

		view := model.View()
		if !strings.Contains(view, "> SELECT") {
			t.Fatalf("expected view to contain a text input:\n%s", view)
		}
	})

	t.Run("ResultsPane Err", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		model.Err = errSQL

		view := model.View()
		if !strings.Contains(view, errSQL.Error()) {
			t.Fatalf("expected view to contain a text input:\n%s", view)
		}
	})
}
