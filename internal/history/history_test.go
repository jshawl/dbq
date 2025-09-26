package history_test

import (
	"math"
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

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("PushMsg", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		_, cmd = hist.Update(history.PushMsg{Query: "select * from users limit 1;"})
		if cmd == nil {
			t.Fatal("expected PushMsg to return a msg")
		}

		_, _ = hist.Update(cmd())
	})

	t.Run("tea.KeyMsg(up)", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist.Push("select * from users limit 1;")
		hist.Push("select * from users limit 2;")
		hist = hist.SetCursor(2)

		// key event
		_, cmd = hist.Update(testutil.MakeKeyMsg(tea.KeyUp))
		// query
		_, cmd = hist.Update(cmd())
		// response
		msg := testutil.AssertMsgType[history.SetInputValueMsg](t, cmd)
		if msg.Value != "select * from users limit 1;" {
			t.Fatalf("expected msg.Direction to be 'previous', got %s", msg.Value)
		}
	})

	t.Run("tea.KeyMsg(down)", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)

		hist.Push("select * from users limit 1;")
		hist.Push("select * from users limit 4;")
		hist = hist.SetCursor(1)

		// key event
		_, cmd = hist.Update(testutil.MakeKeyMsg(tea.KeyDown))
		// query
		_, cmd = hist.Update(cmd())
		// response
		msg := testutil.AssertMsgType[history.SetInputValueMsg](t, cmd)
		if msg.Value != "select * from users limit 4;" {
			t.Fatal("wrong value")
		}
	})

	t.Run("unknownMsg", func(t *testing.T) {
		t.Parallel()

		var cmd tea.Cmd

		hist := setupHistoryModel(t)
		_, cmd = hist.Update(testutil.MakeKeyMsg(tea.KeyLeft))

		if cmd != nil {
			t.Fatal("expected cmd to be nil")
		}
	})
}

func TestPrevious(t *testing.T) {
	t.Parallel()

	t.Run("no entries - no cursor movement", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		cursor, query := hist.Previous()

		if cursor != 0 {
			t.Fatal("expected cursor to be 0")
		}

		if query != "" {
			t.Fatal("expected query to be ''")
		}
	})

	t.Run("one entry - no cursor movement", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		hist.Push("select * from users limit 1;")
		hist.SetCursor(1)
		cursor, query := hist.Previous()

		if cursor != 1 {
			t.Fatalf("expected cursor to be 1, got %d", cursor)
		}

		if query != "select * from users limit 1;" {
			t.Fatal("expected query to be 'select * from users limit 1;'")
		}
	})

	t.Run("two entries - cursor decrements", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		hist.Push("select * from users limit 1;")
		hist.Push("select * from users limit 2;")
		hist.SetCursor(3)
		cursor, query := hist.Previous()

		if cursor != 2 {
			t.Fatalf("expected cursor to be 2, got %d", cursor)
		}

		if query != "select * from users limit 2;" {
			t.Fatalf("expected query to be 'select * from users limit 2;', got %s", query)
		}
	})
}

func TestNext(t *testing.T) {
	t.Parallel()

	t.Run("no entries - no cursor movement", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		cursor, query := hist.Next()

		if cursor != math.MaxInt32 {
			t.Fatalf("expected cursor to be math.MaxInt32, got %d", cursor)
		}

		if query != "" {
			t.Fatal("expected query to be ''")
		}
	})

	t.Run("one entry - no cursor movement", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		hist.Push("select * from users limit 1;")
		hist = hist.SetCursor(1)
		cursor, query := hist.Next()

		if cursor != math.MaxInt32 {
			t.Fatalf("expected cursor to be math.MaxInt32, got %d", cursor)
		}

		if query != "" {
			t.Fatalf("expected query to be '', got %s", query)
		}
	})

	t.Run("two entries - cursor increments", func(t *testing.T) {
		t.Parallel()

		hist := setupHistoryModel(t)
		hist.Push("select * from users limit 1;")

		hist.Push("select * from users limit 2;")
		hist.Push("select * from users limit 3;")
		hist = hist.SetCursor(1)
		cursor, query := hist.Next()

		if cursor != 2 {
			t.Fatalf("expected cursor to be 2, got %d", cursor)
		}

		if query != "select * from users limit 2;" {
			t.Fatalf("expected query to be 'select * from users limit 2;', got %s", query)
		}
	})
}
