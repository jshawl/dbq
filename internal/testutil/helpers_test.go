package testutil_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jshawl/dbq/internal/testutil"
)

type Msg struct {
	ok bool
}

func TestAssertMsgType(t *testing.T) {
	t.Parallel()

	cmd := func() tea.Msg {
		return Msg{
			ok: true,
		}
	}

	msg := testutil.AssertMsgType[Msg](t, cmd)
	if !msg.ok {
		t.Fatal("expected msg.ok to be true")
	}
}
