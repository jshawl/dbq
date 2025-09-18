package main

import (
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	t.Parallel()

	dsn := "postgres://admin:password@localhost:5432/dbq_test"

	ctx := t.Context()

	db, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close(ctx)
}
