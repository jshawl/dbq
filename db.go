package main

import (
	"context"
	"fmt"
	"time"
)

type DB struct {
	inner PGDB
}

type DBQueryResult struct {
	Results  PostgresQueryResult
	Duration time.Duration
}

func NewDB(inner PGDB) *DB {
	return &DB{inner: inner}
}

func (db *DB) Query(ctx context.Context, query string) (DBQueryResult, error) {
	start := time.Now()
	postgresQueryResults, err := db.inner.Query(ctx, query)
	if err != nil {
		return DBQueryResult{}, fmt.Errorf("query %q failed: %w", query, err)
	}
	return DBQueryResult{
		Results:  postgresQueryResults,
		Duration: time.Since(start),
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.inner.Close(ctx)
}
