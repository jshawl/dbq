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

		if result[0].Start != 2 {
			t.Fatalf("expected ow to Start at 2, got %d", result[0].Start)
		}

		if result[0].End != 4 {
			t.Fatalf("expected ow to End at 4, got %d", result[0].End)
		}

		if result[1].Start != 8 {
			t.Fatalf("expected ow to Start at 8, got %d", result[1].Start)
		}

		if result[1].End != 10 {
			t.Fatalf("expected ow to End at 10, got %d", result[1].End)
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		t.Parallel()

		result := search.Search("bro\nwn", "ow")

		if result[0].Start != 2 {
			t.Fatalf("expected ow to Start at 2, got %d", result[0].Start)
		}

		if result[0].End != 5 {
			t.Fatalf("expected ow to End at 5, got %d", result[0].End)
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

	t.Run("two lines, ignores line breaks", func(t *testing.T) {
		t.Parallel()

		result := search.Search("bro\nwn", "ow")

		highlighted := search.Highlight("bro\nwn", result)

		if highlighted == "brown clown" {
			t.Fatalf("failed to highlight, got %s", highlighted)
		}
	})
}
