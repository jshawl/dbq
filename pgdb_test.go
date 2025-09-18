package main

import (
	"fmt"
	"testing"
)

func TestPGDBConnect(t *testing.T) {
	t.Parallel()

	dsn := "postgres://admin:password@localhost:5432/dbq_test"

	ctx := t.Context()

	db, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close(ctx)

	// Test basic query
	results, _ := db.Query(ctx, "SELECT 1 as test_col")
	fmt.Printf("%v", results)
}
