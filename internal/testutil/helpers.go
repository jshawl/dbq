package testutil

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

//nolint:ireturn
func AssertMsgType[T interface{}](t *testing.T, cmd tea.Cmd) T {
	t.Helper()

	if cmd == nil {
		t.Fatalf("%T cmd is nil ", new(T))
	}

	msg := cmd()
	typed, ok := msg.(T)

	if !ok {
		t.Fatalf("Expected msg to be of type %T, got %T", *new(T), msg)
	}

	return typed
}

func MakeKeyMsg(key tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{
		Alt:   false,
		Paste: false,
		Runes: nil,
		Type:  key,
	}
}
