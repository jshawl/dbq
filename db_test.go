package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type mockPGDB struct {
	queryCalled bool
	closeCalled bool
	results     PostgresQueryResult
	queryErr    error
	closeErr    error
}

func (m *mockPGDB) Query(_ context.Context, _ string) (PostgresQueryResult, error) {
	m.queryCalled = true

	return m.results, m.queryErr
}

func (m *mockPGDB) Close(_ context.Context) error {
	m.closeCalled = true

	return m.closeErr
}

func TestQueryReturnsResultsAndDuration(t *testing.T) {
	t.Parallel()

	want := PostgresQueryResult{
		{"id": int32(1), "name": "Alice"},
		{"id": int32(2), "name": "Bob"},
	}
	mock := &mockPGDB{
		closeCalled: false,
		queryCalled: false,
		results:     want,
		queryErr:    nil,
		closeErr:    nil,
	}

	db := NewDB(mock)

	ctx := context.Background()
	got, _ := db.Query(ctx, "SELECT * FROM users")

	if !mock.queryCalled {
		t.Error("expected Query to call inner PGDB.Query")
	}

	for i := range want {
		for k, v := range want[i] {
			if got.Results[i][k] != v {
				t.Errorf("row %d column %s: want %v, got %v", i, k, v, got.Results[i][k])
			}
		}
	}

	if got.Duration <= 0 {
		t.Error("expected positive Duration")
	}
}

var ErrTestQuery = errors.New("boom")

func TestQueryPropagatesError(t *testing.T) {
	t.Parallel()

	mock := &mockPGDB{
		closeCalled: false,
		closeErr:    nil,
		results:     nil,
		queryCalled: false,
		queryErr:    fmt.Errorf("%w", ErrTestQuery),
	}
	db := NewDB(mock)

	_, err := db.Query(context.Background(), "bad sql")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestCloseCallsInnerClose(t *testing.T) {
	t.Parallel()

	mock := &mockPGDB{
		closeCalled: false,
		closeErr:    nil,
		results:     nil,
		queryCalled: false,
		queryErr:    fmt.Errorf("%w", ErrTestQuery),
	}
	db := NewDB(mock)

	err := db.Close(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.closeCalled {
		t.Error("expected Close to call inner PGDB.Close")
	}
}
