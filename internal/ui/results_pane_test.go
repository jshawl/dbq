package ui_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/ui"
)

func TestResultsPane_Update(t *testing.T) {
	t.Parallel()

	t.Run("QueryMsg", func(t *testing.T) {
		t.Parallel()

		userID := 789
		model := ui.ResultsPaneModel{}
		updatedModel, _ := model.Update(ui.QueryMsg{
			Err:     nil,
			Results: makeResults(0, userID),
		})

		typedModel := assertModelType[ui.ResultsPaneModel](t, updatedModel)

		got := typedModel.Results.Results[0]["id"]
		if got != userID {
			t.Fatalf("expected first result to have id %d got %d", userID, got)
		}
	})

	t.Run("QueryMsg - err", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{}
		updatedModel, _ := model.Update(ui.QueryMsg{
			Err:     errSQL,
			Results: db.DBQueryResult{},
		})

		typedModel := assertModelType[ui.ResultsPaneModel](t, updatedModel)

		if typedModel.Err != errSQL {
			t.Fatal("expected query msg err to update model")
		}
	})
}

func TestResultsPane_View(t *testing.T) {
	t.Parallel()

	t.Run("duration with 1 row", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Results: makeResults(time.Millisecond*2345, 123),
		}

		view := model.View()
		if !strings.Contains(view, "(1 row in 2.345s)") {
			t.Fatalf("expected model error to be visible\n %s", view)
		}
	})

	t.Run("duration with 2 rows", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Results: makeResults(time.Millisecond*2345, 123, 456),
		}

		view := model.View()
		if !strings.Contains(view, "(2 rows in 2.345s)") {
			t.Fatalf("expected duration to be visible\n %s", view)
		}
	})
}

func TestResultsPane_ResultsView(t *testing.T) {
	t.Parallel()

	t.Run("results", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Err:     nil,
			Results: makeResults(0, 666),
		}

		view := model.ResultsView()

		matched, _ := regexp.MatchString(
			`---\ncreated_at: 2025-09-21T15:41:22\nid: 666`,
			view,
		)
		if !matched {
			t.Fatalf("expected results to be visible, got \n %s", view)
		}
	})

	t.Run("errors", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Err:     errSQL,
			Results: db.DBQueryResult{},
		}

		view := model.ResultsView()
		if !strings.Contains(view, "sql error") {
			t.Fatal("expected model error to be visible")
		}
	})
}
