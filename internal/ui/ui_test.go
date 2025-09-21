package ui_test

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/history"
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

	model := ui.InitialModel()
	model.TextInput.SetValue("SELECT * FROM users LIMIT 1;")

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

func makeResults(duration time.Duration, userID int, userIDs ...int) db.DBQueryResult {
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

	return db.DBQueryResult{
		Duration: duration,
		Results:  rows,
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

	msg := testutil.AssertMsgType[ui.DBMsg](t, cmd)
	if msg.DB == nil {
		t.Fatal("expected DBmsg to contain db")
	}
}

//nolint:cyclop
func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("keys - enter", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyEnter))

		queryMsg := testutil.AssertMsgType[ui.QueryMsg](t, cmd)
		if len(queryMsg.Results.Results) == 0 {
			t.Fatal("expected results")
		}
	})

	t.Run("keys - ctrl-c", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyCtrlC))

		testutil.AssertMsgType[tea.QuitMsg](t, cmd)
	})

	t.Run("keys - tab", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		if !model.TextInput.Focused() {
			t.Fatal("expected text input to be focused")
		}

		updatedModel, _ := model.Update(testutil.MakeKeyMsg(tea.KeyTab))
		model = assertModelType[ui.Model](t, updatedModel)

		if model.TextInput.Focused() {
			t.Fatal("expected text input not to be focused")
		}

		updatedModel, _ = model.Update(testutil.MakeKeyMsg(tea.KeyTab))
		model = assertModelType[ui.Model](t, updatedModel)

		if !model.TextInput.Focused() {
			t.Fatal("expected text input to be focused")
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

	// t.Run("QueryMsg", func(t *testing.T) {
	// 	t.Parallel()

	// 	userID := 789
	// 	model := setupDatabaseModel(t)
	// 	updatedModel, _ := model.Update(ui.QueryMsg{
	// 		Err:     nil,
	// 		Results: makeResults(0, userID),
	// 	})

	// 	typedModel := assertModelType[ui.Model](t, updatedModel)

	// 	got := typedModel.Results.Results[0]["id"]
	// 	if got != userID {
	// 		t.Fatalf("expected first result to have id %d got %d", userID, got)
	// 	}

	// 	if !typedModel.ResultsPane.Focused() {
	// 		t.Fatal("expected requery results to focus on viewport")
	// 	}
	// })

	// t.Run("QueryMsg - err", func(t *testing.T) {
	// 	t.Parallel()

	// 	model := setupDatabaseModel(t)
	// 	updatedModel, _ := model.Update(ui.QueryMsg{
	// 		Err: errSQL,
	// 		Results: db.DBQueryResult{
	// 			Results:  db.QueryResult{},
	// 			Duration: 0,
	// 		},
	// 	})

	// 	typedModel := assertModelType[ui.Model](t, updatedModel)

	// 	if typedModel.ResultsPane.Focused() {
	// 		t.Fatal("expected requery results not to focus on viewport")
	// 	}
	// })

	t.Run("history.TraveledMsg", func(t *testing.T) {
		t.Parallel()

		model := setupDatabaseModel(t)
		_, cmd := model.Update(history.TraveledMsg{})

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
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

	// t.Run("duration with 1 row", func(t *testing.T) {
	// 	t.Parallel()

	// 	model := setupDatabaseModel(t)
	// 	model.Results = makeResults(time.Millisecond*2345, 123)

	// 	view := model.View()
	// 	if !strings.Contains(view, "(1 row in 2.345s)") {
	// 		t.Fatalf("expected model error to be visible\n %s", view)
	// 	}
	// })

	// t.Run("duration with 2 rows", func(t *testing.T) {
	// 	t.Parallel()

	// 	model := setupDatabaseModel(t)
	// 	model.Results = makeResults(time.Millisecond*2345, 123, 456)

	// 	view := model.View()
	// 	if !strings.Contains(view, "(2 rows in 2.345s)") {
	// 		t.Fatalf("expected duration to be visible\n %s", view)
	// 	}
	// })

	// t.Run("errors", func(t *testing.T) {
	// 	t.Parallel()

	// 	model := setupDatabaseModel(t)
	// 	model.Err = errSQL

	// 	view := model.View()
	// 	if !strings.Contains(view, "sql error") {
	// 		t.Fatal("expected model error to be visible")
	// 	}
	// })

	// t.Run("results", func(t *testing.T) {
	// 	t.Parallel()

	// 	model := setupDatabaseModel(t)
	// 	updatedModel, _ := model.Update(ui.QueryMsg{
	// 		Err:     nil,
	// 		Results: makeResults(0, 666),
	// 	})

	// 	view := updatedModel.View()

	// 	matched, _ := regexp.MatchString(
	// 		`---\s+\ncreated_at: 2025-09-21T15:41:22\s+\nid: 666`,
	// 		view,
	// 	)
	// 	if !matched {
	// 		t.Fatalf("expected results to be visible, got \n %s", view)
	// 	}
	// })
}
