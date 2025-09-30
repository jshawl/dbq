package ui_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/db"
	"github.com/jshawl/dbq/internal/history"
	"github.com/jshawl/dbq/internal/testutil"
	"github.com/jshawl/dbq/internal/ui"
)

func TestQueryPane_Update(t *testing.T) {
	t.Parallel()

	t.Run("keys - unfocused", func(t *testing.T) {
		t.Parallel()

		model := ui.NewQueryPaneModel()
		model = model.Blur()
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyEnter))

		if cmd != nil {
			t.Fatal("expected query pane to ignore keys while not focused")
		}
	})

	t.Run("keys - enter", func(t *testing.T) {
		t.Parallel()

		model := ui.NewQueryPaneModel()
		want := "select * from posts limit 1;"
		model.TextInput.SetValue(want)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyEnter))

		queryMsg := testutil.AssertMsgType[ui.QueryExecMsg](t, cmd)

		if queryMsg.Value != want {
			t.Fatalf("expected QueryMsg.Value to be set, got %s", queryMsg.Value)
		}
	})
	t.Run("history.SetInputValueMsg", func(t *testing.T) {
		t.Parallel()

		model := ui.NewQueryPaneModel()
		_, cmd := model.Update(history.SetInputValueMsg{
			Value: "select * from posts limit 1;",
		})

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})

	t.Run("QueryResponseReceivedMsg - Err", func(t *testing.T) {
		t.Parallel()

		model := ui.NewQueryPaneModel()
		_, cmd := model.Update(ui.QueryResponseReceivedMsg{
			QueryMsg: ui.QueryMsg{
				Duration: 0,
				Err:      errSQL,
				Query:    "not sql",
				Results:  db.QueryResult{},
			},
		})

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})

	t.Run("QueryResponseReceivedMsg - Results", func(t *testing.T) {
		t.Parallel()

		model := ui.NewQueryPaneModel()
		_, cmd := model.Update(ui.QueryResponseReceivedMsg{
			QueryMsg: ui.QueryMsg{
				Duration: time.Millisecond * 2345,
				Err:      nil,
				Results:  makeResults(456),
				Query:    "select * from foo;",
			},
		})

		msg := testutil.AssertMsgType[history.PushMsg](t, cmd)

		if msg.Query != "select * from foo;" {
			t.Fatal("expected history push msg")
		}
	})
}

func TestQueryPane_View(t *testing.T) {
	t.Parallel()

	model := ui.NewQueryPaneModel()

	view := model.View()
	if !strings.Contains(view, "> SELECT") {
		t.Fatalf("expected view to have placeholder, got %s", view)
	}
}
