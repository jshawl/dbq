package history_test

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/history"
	"github.com/jshawl/dbq/internal/testutil"
)

func setupHistoryStore(t *testing.T) string {
	t.Helper()

	return t.TempDir()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("creates a sqlite db", func(t *testing.T) {
		t.Parallel()

		dir := setupHistoryStore(t)

		path := dir + "/foo.db"
		if fileExists(path) {
			t.Fatal("expected db not to be created")
		}

		h := history.Init(path)

		h.Cleanup()

		if !fileExists(path) {
			t.Fatal("expected db to be created")
		}
	})

	t.Run("does not overwrite existing db", func(t *testing.T) {
		t.Parallel()

		dir := setupHistoryStore(t)
		path := dir + "/foo.db"
		h1 := history.Init(path)
		h2 := history.Init(path)

		h1.Cleanup()
		h2.Cleanup()
	})
}

func setupHistoryModel(t *testing.T) history.Model {
	t.Helper()

	dir := setupHistoryStore(t)
	path := dir + "/foo.db"

	model := history.Init(path)

	t.Cleanup(func() {
		defer model.Cleanup()
	})

	return model
}

//nolint:cyclop
func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("PushMsg", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)
		_, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 1;"})
		testutil.AssertMsgType[history.PushedMsg](t, cmd)
	})

	t.Run("UnknownMsg", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		_, cmd = hist.Update(nil)
		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})

	t.Run("TravelMsg(previous) does not update the Value if no results", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.TravelMsg{Direction: "previous"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "" {
			t.Fatalf("expected empty string, got %s", hist.Value)
		}
	})

	t.Run("TravelMsg(previous) does not update the Value if one result", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 1;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.TravelMsg{Direction: "previous"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "" {
			t.Fatalf("expected empty string, got %s", hist.Value)
		}
	})

	t.Run("TravelMsg(previous) returns the last row if the cursor is 0", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 1;"})
		testutil.AssertMsgType[history.PushedMsg](t, cmd)

		// reset
		hist, _ = hist.Update(history.PushedMsg{})

		hist, cmd = hist.Update(history.TravelMsg{Direction: "previous"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "select * from users limit 1;" {
			t.Fatalf("expected current, got %s", hist.Value)
		}
	})

	t.Run("TravelMsg(previous) sets the Value to the penultimate row", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 1;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 2;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 3;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.TravelMsg{Direction: "previous"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "select * from users limit 2;" {
			t.Fatalf("expected current, got %s", hist.Value)
		}
	})

	t.Run("TravelMsg(next) sets the Value to the next row", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 1;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.PushMsg{Entry: "select * from users limit 2;"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.PushedMsg](t, cmd))

		hist, cmd = hist.Update(history.TravelMsg{Direction: "previous"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "select * from users limit 1;" {
			t.Fatalf("expected current, got %s", hist.Value)
		}

		hist, cmd = hist.Update(history.TravelMsg{Direction: "next"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "select * from users limit 2;" {
			t.Fatalf("expected current, got %s", hist.Value)
		}
	})

	t.Run("TravelMsg(next) does not update the Value if no results", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist, cmd = hist.Update(history.TravelMsg{Direction: "next"})
		hist, _ = hist.Update(testutil.AssertMsgType[history.TraveledMsg](t, cmd))

		if hist.Value != "" {
			t.Fatalf("expected empty string, got %s", hist.Value)
		}
	})

	t.Run("tea.KeyMsg(up)", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		_, cmd = hist.Update(testutil.MakeKeyMsg(tea.KeyUp))

		msg := testutil.AssertMsgType[history.TravelMsg](t, cmd)
		if msg.Direction != "previous" {
			t.Fatal("expected msg.Direction to be 'previous'")
		}
	})

	t.Run("tea.KeyMsg(down)", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		_, cmd = hist.Update(testutil.MakeKeyMsg(tea.KeyDown))

		msg := testutil.AssertMsgType[history.TravelMsg](t, cmd)
		if msg.Direction != "next" {
			t.Fatal("expected msg.Direction to be 'previous'")
		}
	})
}
