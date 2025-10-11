package searchableviewport

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/search"
)

//nolint:recvcheck // to match bubbletea interface
type Model struct {
	Height int
	Width  int
	Search search.Model

	content          string
	highlightContent string
	currentMatch     int
	matches          []search.SearchMatch
	ready            bool
	viewport         viewport.Model
}

type WindowSizeMsg struct {
	Height int
	Width  int
}

type SearchDirection int

const (
	SearchDirectionDown SearchDirection = iota
	SearchDirectionUp
)

func NewSearchableViewportModel() Model {
	return Model{
		Height: 0,
		Width:  0,
		Search: search.NewSearchModel(),

		content:          "",
		highlightContent: "",
		currentMatch:     -1,
		matches:          nil,
		ready:            false,
		viewport:         viewport.New(0, 0),
	}
}

func (model *Model) SetContent(str string) {
	model.content = str
	model.Search = search.NewSearchModel()
	model.viewport.SetContent(str)
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	const footerHeight = 1

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if model.Search.Focused() {
			updatedSearchModel, cmd := model.Search.Update(msg)
			model.Search = updatedSearchModel

			return model, cmd
		}
		//nolint:gocritic
		switch msg.String() {
		case "n":
			var cmd tea.Cmd

			model.currentMatch = (model.currentMatch + 1) % len(model.matches)
			model.highlightContent = search.Highlight(model.content, model.matches, model.currentMatch)
			model.viewport.SetContent(model.highlightContent)
			match := model.matches[model.currentMatch]
			model.viewport.YOffset = GetYOffset(
				match.ScreenYPosition,
				model.viewport.YOffset,
				model.viewport.TotalLineCount(),
				model.viewport.Height,
				SearchDirectionDown,
			)

			model.viewport, cmd = model.viewport.Update(msg)

			return model, cmd
		case "N":
			var cmd tea.Cmd

			model.currentMatch = (model.currentMatch - 1) % len(model.matches)
			if model.currentMatch < 0 {
				model.currentMatch = len(model.matches) - 1
			}
			model.highlightContent = search.Highlight(model.content, model.matches, model.currentMatch)
			model.viewport.SetContent(model.highlightContent)
			match := model.matches[model.currentMatch]
			model.viewport.YOffset = GetYOffset(
				match.ScreenYPosition,
				model.viewport.YOffset,
				model.viewport.TotalLineCount(),
				model.viewport.Height,
				SearchDirectionUp,
			)

			model.viewport, cmd = model.viewport.Update(msg)

			return model, cmd
		}
	case WindowSizeMsg:
		height := msg.Height - footerHeight
		if !model.ready {
			model.viewport = viewport.New(msg.Width, height)
			model.ready = true
		} else {
			model.viewport.Width = msg.Width
			model.viewport.Height = height
		}

		return model, nil
	case search.SearchMsg:
		value := msg.Value
		model.matches = search.Search(model.content, value)
		model.currentMatch = 0
		model.highlightContent = search.Highlight(model.content, model.matches, model.currentMatch)
		model.viewport.SetContent(model.highlightContent)

		return model, nil
	case search.SearchClearMsg:
		model.highlightContent = ""
		model.viewport.SetContent(model.content)

		return model, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	model.Search, cmd = model.Search.Update(msg)
	cmds = append(cmds, cmd)
	model.viewport, cmd = model.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func GetYOffset(
	screenYPosition int,
	viewportYOffset int,
	viewportTotalLineCount int,
	viewportHeight int,
	searchDirection SearchDirection,
) int {
	log.Printf("before change syp: %d, vo: %d, vtlc: %d vh: %d", screenYPosition, viewportYOffset, viewportTotalLineCount, viewportHeight)
	// below current viewport
	if screenYPosition > viewportYOffset+viewportHeight {
		log.Printf("syp is below: %d", screenYPosition)
		if searchDirection == SearchDirectionDown {
			return screenYPosition
		} else {
			return viewportTotalLineCount - viewportHeight
		}
	}

	// above current viewport
	if screenYPosition < viewportYOffset {
		log.Printf("syp is above: %d", screenYPosition)
		if searchDirection == SearchDirectionDown {
			return screenYPosition
		} else {
			maybeYOffset := screenYPosition - viewportHeight + 1
			log.Printf("maybeYOffset: %d", maybeYOffset)
			if maybeYOffset < 0 {
				return 0
			}
			return maybeYOffset
		}
	}

	// top offset is too far
	maxYOffset := viewportTotalLineCount - viewportHeight
	if screenYPosition > maxYOffset {
		log.Printf("too far: %d", maxYOffset)
		return maxYOffset
	}

	log.Printf("already visible: %d", viewportYOffset)
	// already visible
	return viewportYOffset
}

func (model Model) View() string {
	return model.viewport.View()
}

func (model Model) FooterView() string {
	if model.Search.Focused() || model.Search.Value != "" {
		return model.Search.View()
	}

	return ""
}
