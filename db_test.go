package main

import (
	"context"
	"fmt"
	"testing"
)

type mockStore struct {
	queryFunc func(ctx context.Context, sql string) (QueryResult, error)
	closeFunc func(ctx context.Context) error
}

func (m *mockStore) Query(ctx context.Context, sql string) (QueryResult, error) {
	return m.queryFunc(ctx, sql)
}
func (m *mockStore) Close(ctx context.Context) error {
	if m.closeFunc != nil {
		return m.closeFunc(ctx)
	}
	return nil
}

func TestQueryFormatsResults(t *testing.T) {
	m := &mockStore{
		queryFunc: func(ctx context.Context, sql string) (QueryResult, error) {
			return QueryResult{
				{"id": int32(1), "name": "Alice"},
				{"id": int32(2), "name": "Bob"},
			}, nil
		},
	}

	ctx := context.Background()
	got, err := m.Query(ctx, "SELECT * FROM users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := QueryResult{
		{"id": int32(1), "name": "Alice"},
		{"id": int32(2), "name": "Bob"},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d rows, got %d", len(want), len(got))
	}
	for i := range want {
		for k, v := range want[i] {
			if fmt.Sprintf("%v", got[i][k]) != fmt.Sprintf("%v", v) {
				t.Errorf("row %d column %s: want %v, got %v", i, k, v, got[i][k])
			}
		}
	}
}
