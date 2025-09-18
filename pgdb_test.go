package main

import (
	"fmt"
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	t.Parallel()

	dsn := "postgres://admin:password@localhost:5432/dbq_test"

	ctx := t.Context()

	database, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		t.Fatalf("%v", err)
	}

	defer func() {
		err := database.Close(ctx)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()
}

func TestPGDBQuery(t *testing.T) {
	t.Parallel()

	dsn := "postgres://admin:password@localhost:5432/dbq_test"

	ctx := t.Context()

	database, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		t.Fatalf("%v", err)
	}

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

	have, err := database.Query(t.Context(), "SELECT * FROM users")
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

	defer func() {
		err := database.Close(ctx)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()
}
