package searchableviewport_test

import (
	"strings"
	"testing"

	"github.com/jshawl/dbq/internal/searchableviewport"
)

func initializeViewport(model searchableviewport.Model, t *testing.T) searchableviewport.Model {
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
}

func TestView(t *testing.T) {
	t.Parallel()

	model := initializeViewport(searchableviewport.NewSearchableViewportModel(), t)
	model.SetContent("success!")
	if !strings.Contains(model.View(), "success!") {
		t.Fatal("expected view to have updated content")
	}
}
