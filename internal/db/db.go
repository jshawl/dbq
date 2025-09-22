package db

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type QueryResult []map[string]interface{}

type Queryable interface {
	Query(ctx context.Context, sql string) (QueryResult, error)
	Close(ctx context.Context) error
}

type DB struct {
	inner Queryable
}

type DBQueryResult struct {
	Err      error
	Results  QueryResult
	Duration time.Duration
}

var (
	ErrDBQuery = errors.New("failed to call Query")
	ErrDBClose = errors.New("failed to call Close")
)

func NewDB(inner Queryable) *DB {
	return &DB{inner: inner}
}

func (db *DB) Query(ctx context.Context, query string) DBQueryResult {
	start := time.Now()

	postgresQueryResults, err := db.inner.Query(ctx, query)

	return DBQueryResult{
		Err:      err,
		Results:  postgresQueryResults,
		Duration: time.Since(start),
	}
}

func (db *DB) Close(ctx context.Context) error {
	err := db.inner.Close(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDBClose, err)
	}

	return nil
}
