package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockRows struct {
	idx    int
	fields []pgconn.FieldDescription
	values [][]interface{}
}

func (m *mockRows) Close()                                       {}
func (m *mockRows) Err() error                                   { return nil }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return m.fields }
func (m *mockRows) Next() bool {
	if m.idx < len(m.values) {
		m.idx++
		return true
	}
	return false
}
func (m *mockRows) Values() ([]interface{}, error) {
	if m.idx == 0 || m.idx > len(m.values) {
		return nil, errors.New("out of range")
	}
	return m.values[m.idx-1], nil
}

// Unused pgx.Rows methods can be stubbed:
func (m *mockRows) RawValues() [][]byte            { return nil }
func (m *mockRows) Conn() *pgx.Conn                { return nil }
func (m *mockRows) CommandTag() pgconn.CommandTag  { return pgconn.CommandTag{} }
func (m *mockRows) Scan(dest ...interface{}) error { return nil }

type mockConn struct {
	rows pgx.Rows
}

func (m *mockConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return m.rows, nil
}
func (m *mockConn) Close(ctx context.Context) error { return nil }

func TestQueryFormatsResults(t *testing.T) {
	ctx := context.Background()

	fields := []pgconn.FieldDescription{
		{Name: "id"},
		{Name: "name"},
	}
	values := [][]interface{}{
		{1, "Alice"},
		{2, "Bob"},
	}

	fakeRows := &mockRows{fields: fields, values: values}
	db := &DB{conn: &mockConn{rows: fakeRows}}

	got, err := db.Query(ctx, "SELECT * FROM users")
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
