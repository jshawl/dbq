package search_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/search"
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

func TestHighlight(t *testing.T) {
	t.Parallel()

	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Run("one line, two matches", func(t *testing.T) {
		t.Parallel()

		result := search.Search("green queen", "ee")

		highlighted := search.Highlight("green queen", result)

		if highlighted == "green queen" {
			t.Fatalf("failed to highlight, got %s", highlighted)
		}
	})

	t.Run("two lines, two matches", func(t *testing.T) {
		t.Parallel()

		result := search.Search("brown\nclown", "ow")

		highlighted := search.Highlight("brown\nclown", result)

		if highlighted == "brown clown" {
			t.Fatalf("failed to highlight, got %s", highlighted)
		}
	})
}
