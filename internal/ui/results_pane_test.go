package ui_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/testutil"
	"github.com/jshawl/dbq/internal/ui"
)

func TestResultsPane_Update(t *testing.T) {
	t.Parallel()

	t.Run("WindowSizeMsg", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      0,
			IsSearching: false,
			Results:     db.QueryResult{},
			Width:       0,
			YPosition:   0,
		}
		updatedModel, _ := model.Update(ui.WindowSizeMsg{
			Height:    42,
			Width:     37,
			YPosition: 2,
		})

		// height - yposition (2) - footer (1)
		if updatedModel.Height != 39 {
			t.Fatalf("expected WindowSizeMsg.Height to be 40, got %d", updatedModel.Height)
		}

		if updatedModel.Width != 37 {
			t.Fatalf("expected WindowSizeMsg.Width to be 37, got %d", updatedModel.Width)
		}
	})

	t.Run("QueryResponseReceivedMsg", func(t *testing.T) {
		t.Parallel()

		userID := 789
		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     db.QueryResult{},
			Width:       80,
			YPosition:   0,
		}
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

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     db.QueryResult{},
			Width:       80,
			YPosition:   0,
		}
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

	t.Run("slash IsSearching", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     db.QueryResult{},
			Width:       80,
			YPosition:   0,
		}
		model = model.Focus()

		updatedModel, _ := model.Update(tea.KeyMsg{
			Alt:   false,
			Paste: false,
			Type:  tea.KeyRunes,
			Runes: []rune{'/'},
		})

		if !updatedModel.IsSearching {
			t.Fatal("expected IsSearching to be true")
		}
	})

	t.Run("esc IsSearching", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: true,
			Results:     db.QueryResult{},
			Width:       80,
			YPosition:   0,
		}
		model = model.Focus()

		updatedModel, _ := model.Update(testutil.MakeKeyMsg(tea.KeyEscape))

		if updatedModel.IsSearching {
			t.Fatal("expected IsSearching to be false")
		}
	})
}

func TestResultsPane_View(t *testing.T) {
	t.Parallel()

	t.Run("duration with 1 row", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    time.Millisecond * 2345,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     makeResults(123),
			Width:       80,
			YPosition:   0,
		}

		view := model.View()
		if !strings.Contains(view, "(1 row in 2.345s)") {
			t.Fatalf("expected model duration to be visible\n %s", view)
		}
	})

	t.Run("duration with 2 rows", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    time.Millisecond * 2345,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     makeResults(123, 456),
			Width:       80,
			YPosition:   0,
		}

		view := model.View()
		if !strings.Contains(view, "(2 rows in 2.345s)") {
			t.Fatalf("expected duration to be visible\n %s", view)
		}
	})

	t.Run("IsSearching", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: true,
			Results:     makeResults(123, 456),
			Width:       80,
			YPosition:   0,
		}

		view := model.View()
		if !strings.Contains(view, "/") {
			t.Fatalf("expected searching to show slash \n %s", view)
		}
	})
}

func TestResultsPane_ResultsView(t *testing.T) {
	t.Parallel()

	t.Run("results", func(t *testing.T) {
		t.Parallel()

		model := ui.ResultsPaneModel{
			Duration:    0,
			Err:         nil,
			Height:      80,
			IsSearching: false,
			Results:     makeResults(0, 666),
			Width:       80,
			YPosition:   0,
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
			Duration:    0,
			Err:         errSQL,
			Height:      80,
			IsSearching: false,
			Results:     db.QueryResult{},
			Width:       80,
			YPosition:   0,
		}

		view := model.ResultsView()
		if !strings.Contains(view, "sql error") {
			t.Fatal("expected model error to be visible")
		}
	})
}

func TestResultsPane_Resize(t *testing.T) {
	t.Parallel()

	model := ui.ResultsPaneModel{
		Duration:    0,
		Err:         errSQL,
		Height:      80,
		IsSearching: false,
		Results:     db.QueryResult{},
		Width:       80,
		YPosition:   0,
	}

	model = model.Resize(20, 30, 1)
	if model.Width != 20 {
		t.Fatal("expected resize to set width")
	}

	model = model.Resize(40, 50, 1)
	// subtracts the header and footer
	if model.Height != 48 {
		t.Fatalf("expected resize to set height, got %d", model.Height)
	}
}
