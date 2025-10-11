package search_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/search"
	"github.com/jshawl/dbq/internal/testutil"
	"github.com/muesli/termenv"
)

func TestSearch(t *testing.T) {
	t.Parallel()

	t.Run("one line", func(t *testing.T) {
		t.Parallel()

		result := search.Search("brown clown", "ow")

		if result[0].BufferStart != 2 {
			t.Fatalf("expected ow to Start at 2, got %d", result[0].BufferStart)
		}

		if result[0].BufferEnd != 4 {
			t.Fatalf("expected ow to End at 4, got %d", result[0].BufferEnd)
		}

		if result[1].BufferStart != 8 {
			t.Fatalf("expected ow to Start at 8, got %d", result[1].BufferStart)
		}

		if result[1].BufferEnd != 10 {
			t.Fatalf("expected ow to End at 10, got %d", result[1].BufferEnd)
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		t.Parallel()

		result := search.Search("br\nown", "ow")

		if result[0].BufferStart != 3 {
			t.Fatalf("expected ow to Start at 3, got %d", result[0].BufferStart)
		}

		if result[0].BufferEnd != 5 {
			t.Fatalf("expected ow to End at 5, got %d", result[0].BufferEnd)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		t.Parallel()

		result := search.Search("a", "abc")
		if len(result) != 0 {
			t.Fatalf("expected no matches, got %d", len(result))
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("slash IsSearching", func(t *testing.T) {
		t.Parallel()

		model := search.NewSearchModel()
		model = model.Focus()

		updatedModel, _ := model.Update(tea.KeyMsg{
			Alt:   false,
			Paste: false,
			Type:  tea.KeyRunes,
			Runes: []rune{'/'},
		})

		if !updatedModel.Focused() {
			t.Fatal("expected IsSearching to be true")
		}
	})

	t.Run("esc IsSearching", func(t *testing.T) {
		t.Parallel()

		model := search.NewSearchModel()
		model = model.Focus()

		updatedModel, _ := model.Update(testutil.MakeKeyMsg(tea.KeyEscape))

		if updatedModel.Focused() {
			t.Fatal("expected IsSearching to be false")
		}
	})
}

func TestHighlight(t *testing.T) {
	t.Parallel()

	lipgloss.SetColorProfile(termenv.TrueColor)

	t.Run("marks the first match", func(t *testing.T) {
		t.Parallel()

		result := search.Search("abcd", "bc")
		highlighted := search.Highlight("abcd", result, 0)

		if highlighted != "a"+search.WithBlackBackground("b")+search.WithYellowBackground("c")+"d" {
			t.Fatalf("expected first character to be marked, got %s", highlighted)
		}
	})

	t.Run("two matches only marks the first highlight", func(t *testing.T) {
		t.Parallel()

		result := search.Search("brown clown", "ow")

		highlighted := search.Highlight("brown clown", result, 0)

		expected := "br" +
			search.WithBlackBackground("o") +
			search.WithYellowBackground("w") +
			"n cl" +
			search.WithYellowBackground("ow") +
			"n"
		if highlighted != expected {
			t.Fatalf("failed to highlight, got %s", highlighted)
		}
	})
}
