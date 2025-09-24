package ui_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/history"
	"github.com/jshawl/dbq/internal/testutil"
	"github.com/jshawl/dbq/internal/ui"
)

func TestQueryPane_Update(t *testing.T) {
	t.Parallel()

	t.Run("keys - enter", func(t *testing.T) {
		t.Parallel()

		//nolint:exhaustruct
		model := ui.QueryPaneModel{}.New()
		want := "select * from posts limit 1;"
		model.TextInput.SetValue(want)
		_, cmd := model.Update(testutil.MakeKeyMsg(tea.KeyEnter))

		queryMsg := testutil.AssertMsgType[ui.QueryExecMsg](t, cmd)

		if queryMsg.Value != want {
			t.Fatalf("expected QueryMsg.Value to be set, got %s", queryMsg.Value)
		}
	})
	t.Run("history.TraveledMsg", func(t *testing.T) {
		t.Parallel()

		//nolint:exhaustruct
		model := ui.QueryPaneModel{}.New()
		_, cmd := model.Update(history.TraveledMsg{
			Value: "select * from posts limit 1;",
		})

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})
}

func TestQueryPane_View(t *testing.T) {
	t.Parallel()

	//nolint:exhaustruct
	model := ui.QueryPaneModel{}.New()

	view := model.View()
	if !strings.Contains(view, "> SELECT") {
		t.Fatalf("expected view to have placeholder, got %s", view)
	}
}
