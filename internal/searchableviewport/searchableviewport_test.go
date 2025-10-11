package searchableviewport_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/dbq/internal/search"
	"github.com/jshawl/dbq/internal/searchableviewport"
	"github.com/muesli/termenv"
)

func initializeViewport(t *testing.T, model searchableviewport.Model) searchableviewport.Model {
	t.Helper()

	msg := searchableviewport.WindowSizeMsg{
		Height: 10,
		Width:  10,
	}
	updatedModel, _ := model.Update(msg)

	return updatedModel
}

func TestNewSearchableViewportModel(t *testing.T) {
	t.Parallel()

	model := searchableviewport.NewSearchableViewportModel()

	if model.Height != 0 || model.Width != 0 {
		t.Fatal("expected model viewport defaults")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("WindowSizeMsg", func(t *testing.T) {
		t.Parallel()

		model := searchableviewport.NewSearchableViewportModel()
		msg := searchableviewport.WindowSizeMsg{
			Height: 10,
			Width:  10,
		}
		updatedModel, cmd := model.Update(msg)

		if cmd != nil {
			t.Fatal("expected no cmd from window resize")
		}

		_, cmd = updatedModel.Update(msg)
		if cmd != nil {
			t.Fatal("expected no cmd from second window resize")
		}
	})

	t.Run("keys ignored if search not focused", func(t *testing.T) {
		t.Parallel()

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		msg := tea.KeyMsg{
			Alt:   false,
			Paste: false,
			Type:  tea.KeyRunes,
			Runes: []rune{'j'},
		}

		updatedModel, _ := model.Update(msg)
		if strings.Contains(updatedModel.Search.Value, "j") {
			t.Fatal("expected keys not to update search value while unfocused")
		}
	})

	t.Run("keys update search model if focused", func(t *testing.T) {
		t.Parallel()

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		model.Search = model.Search.Focus()
		msg := tea.KeyMsg{
			Alt:   false,
			Paste: false,
			Type:  tea.KeyRunes,
			Runes: []rune{'j'},
		}

		updatedModel, _ := model.Update(msg)
		if !strings.Contains(updatedModel.Search.Value, "j") {
			t.Fatal("expected keys to update search value while focused")
		}
	})

	t.Run("search.SearchMsg", func(t *testing.T) {
		t.Parallel()
		lipgloss.SetColorProfile(termenv.TrueColor)

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		model.SetContent("abcd")

		msg := search.SearchMsg{
			Value: "abc",
		}

		updatedModel, _ := model.Update(msg)
		if updatedModel.View() == model.View() {
			t.Fatal("expected view to have highlighted content")
		}
	})

	t.Run("search.SearchClearMsg", func(t *testing.T) {
		t.Parallel()
		lipgloss.SetColorProfile(termenv.TrueColor)

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		model.SetContent("abcd")

		msg := search.SearchMsg{
			Value: "abc",
		}

		updatedModel, _ := model.Update(msg)
		if updatedModel.View() == model.View() {
			t.Fatal("expected view to have highlighted content")
		}

		updatedModel, _ = model.Update(search.SearchClearMsg{})
		if updatedModel.View() != model.View() {
			t.Fatal("expected highlighted view to be unhighlighted")
		}
	})

	t.Run("n navigates search results", func(t *testing.T) {
		t.Parallel()

		lipgloss.SetColorProfile(termenv.TrueColor)

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		model.SetContent("abcd abef")

		msg := search.SearchMsg{
			Value: "ab",
		}

		updatedModel, _ := model.Update(msg)

		navigatedModel, _ := updatedModel.Update(tea.KeyMsg{
			Alt:   false,
			Paste: false,
			Type:  tea.KeyRunes,
			Runes: []rune{'n'},
		})
		if updatedModel.View() == navigatedModel.View() {
			t.Fatal("expected highlighted view to update mark")
		}
	})
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
	model.Search.Value = "prev search"
	model.SetContent("success!")

	if model.Search.Value == "prev search" {
		t.Fatal("expected SetContent to reset search model")
	}
}

func TestGetYOffset(t *testing.T) {
	t.Parallel()

	t.Run("down - match below current view", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(11, 0, 30, 10, searchableviewport.SearchDirectionDown)
		if offset != 11 {
			t.Fatalf("expected offset to be next screen position, got %d", offset)
		}
	})

	t.Run("down - match above current view", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(1, 10, 30, 10, searchableviewport.SearchDirectionDown)
		if offset != 1 {
			t.Fatalf("expected offset to be next screen position, got %d", offset)
		}
	})

	t.Run("down - match after current view but near end", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(29, 20, 30, 10, searchableviewport.SearchDirectionDown)
		if offset != 20 {
			t.Fatalf("expected offset to be top of end, got %d", offset)
		}
	})

	t.Run("up - match below current view", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(11, 0, 30, 10, searchableviewport.SearchDirectionUp)
		if offset != 20 {
			t.Fatalf("expected offset to page up from end, got %d", offset)
		}
	})

	t.Run("up - match above current view", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(1, 10, 30, 10, searchableviewport.SearchDirectionUp)
		if offset != 0 {
			t.Fatalf("expected offset to be 0, got %d", offset)
		}
	})

	t.Run("up - match above current view close to top", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(9, 10, 30, 10, searchableviewport.SearchDirectionUp)
		if offset != 0 {
			t.Fatalf("expected offset to be 0, got %d", offset)
		}
	})

	t.Run("up - match after current view but near end", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(29, 20, 30, 10, searchableviewport.SearchDirectionUp)
		if offset != 20 {
			t.Fatalf("expected offset to be top of end, got %d", offset)
		}
	})

	t.Run("currently visible", func(t *testing.T) {
		t.Parallel()

		offset := searchableviewport.GetYOffset(2, 0, 30, 10, searchableviewport.SearchDirectionDown)
		if offset != 0 {
			t.Fatal("expected viewport offset not to change")
		}
	})
}

func TestView(t *testing.T) {
	t.Parallel()

	model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
	model.SetContent("success!")

	if !strings.Contains(model.View(), "success!") {
		t.Fatal("expected view to have updated content")
	}
}

func TestFooterView(t *testing.T) {
	t.Parallel()

	t.Run("focused input", func(t *testing.T) {
		t.Parallel()
		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())

		model.Search = model.Search.Focus()
		if model.FooterView() == "" {
			t.Fatal("expected footer view to render text input")
		}
	})

	t.Run("unfocused input", func(t *testing.T) {
		t.Parallel()

		model := initializeViewport(t, searchableviewport.NewSearchableViewportModel())
		if model.FooterView() != "" {
			t.Fatal("expected footer view to render empty string")
		}
	})
}
