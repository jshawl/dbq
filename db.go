package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type DB struct {
	inner PGDBI // my interface, not the struct
}

type DBQueryResult struct {
	Results  PostgresQueryResult
	Duration time.Duration
}

var (
	ErrDBQuery = errors.New("failed to call Query")
	ErrDBClose = errors.New("failed to call Close")
)

func NewDB(inner PGDBI) *DB {
	return &DB{inner: inner}
}

func (db *DB) Query(ctx context.Context, query string) (DBQueryResult, error) {
	start := time.Now()

	postgresQueryResults, err := db.inner.Query(ctx, query)
	if err != nil {
		return DBQueryResult{}, fmt.Errorf("%w: %w", ErrDBQuery, err)
	}

	return DBQueryResult{
		Results:  postgresQueryResults,
		Duration: time.Since(start),
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	err := db.inner.Close(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDBClose, err)
	}

	return nil
}
