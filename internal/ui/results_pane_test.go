package ui_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/ui"
)

func TestResultsPane_Update(t *testing.T) {
	t.Parallel()

	t.Run("QueryResponseReceivedMsg", func(t *testing.T) {
		t.Parallel()

		userID := 789
		model := ui.NewResultsPaneModel()
		updatedModel, _ := model.Update(ui.QueryResponseReceivedMsg{
			QueryMsg: ui.QueryMsg{
				Duration: 0,
				Err:      nil,
				Results:  makeResults(userID),
				Query:    "select * from posts",
			},
		})

		got := updatedModel.Results[0]["id"]
		if got != userID {
			t.Fatalf("expected first result to have id %d got %d", userID, got)
		}
	})

	t.Run("QueryResponseReceivedMsg - err", func(t *testing.T) {
		t.Parallel()

		model := ui.NewResultsPaneModel()
		updatedModel, _ := model.Update(ui.QueryResponseReceivedMsg{
			QueryMsg: ui.QueryMsg{
				Duration: 0,
				Err:      errSQL,
				Results:  db.QueryResult{},
				Query:    "not sql",
			},
		})

		if !errors.Is(updatedModel.Err, errSQL) {
			t.Fatal("expected query msg err to update model")
		}
	})
}

func TestResultsPane_View(t *testing.T) {
	t.Parallel()

	t.Run("duration with 1 row", func(t *testing.T) {
		t.Parallel()

		model := ui.NewResultsPaneModel()
		model.Duration = time.Millisecond * 2345
		model.Results = makeResults(123)

		view := model.View()
		if !strings.Contains(view, "(1 row in 2.345s)") {
			t.Fatalf("expected model duration to be visible\n %s", view)
		}
	})

	t.Run("duration with 2 rows", func(t *testing.T) {
		t.Parallel()

		model := ui.NewResultsPaneModel()
		model.Duration = time.Millisecond * 2345
		model.Results = makeResults(123, 456)

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

		model := ui.NewResultsPaneModel()
		model.Results = makeResults(666)

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

		model := ui.NewResultsPaneModel()
		model.Err = errSQL

		view := model.ResultsView()
		if !strings.Contains(view, "sql error") {
			t.Fatal("expected model error to be visible")
		}
	})
}
