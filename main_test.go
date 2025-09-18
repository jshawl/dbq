package main

import "testing"

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		_ = run(DSN)
	})
	t.Run("error - empty dsn", func(t *testing.T) {
		t.Parallel()

		err := run("")
		if err == nil {
			t.Fatalf("expected error for empty dsn")
		}
	})

	t.Run("error - invalid dsn", func(t *testing.T) {
		t.Parallel()

		err := run("localhost:1234")
		if err == nil {
			t.Fatalf("expected error for invalid dsn")
		}
	})
}
