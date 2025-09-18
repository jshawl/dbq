package main

import (
	"context"
	"fmt"
	"testing"
)

const DSN = "postgres://admin:password@localhost:5432/dbq_test"

func setupDatabase(t *testing.T, dsn string) PGDB {
	t.Helper()

	ctx := context.Background()

	database, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Cleanup(func() {
		err := database.Close(ctx)
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})

	return database
}

func TestNewPostgresDB(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		database := setupDatabase(t, DSN)
		_ = database
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		_, err := NewPostgresDB(t.Context(), "")
		if err == nil {
			t.Fatal("expected error for pgdb Connect")
		}
	})
}

func TestPGDB_Query(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		database := setupDatabase(t, DSN)

		want := []map[string]interface{}{
			{
				"id":         1,
				"first_name": "John",
			},
			{
				"id":         2,
				"first_name": "Jane",
			},
		}

		have, err := database.Query(context.Background(), "SELECT * FROM users")
		if err != nil {
			t.Fatalf("%v", err)
		}

		for i := range want {
			for k, v := range want[i] {
				if fmt.Sprintf("%v", have[i][k]) != fmt.Sprintf("%v", v) {
					t.Errorf("row %d column %s: want %v, got %v", i, k, v, have[i][k])
				}
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		database := setupDatabase(t, DSN)

		_, err := database.Query(context.Background(), "! not sql !")
		if err == nil {
			t.Fatalf("expected error for pgdb Query")
		}
	})
}
