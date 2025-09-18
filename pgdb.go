package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PGDB struct {
	conn *pgx.Conn
}

var (
	ErrConnect = errors.New("failed to call Connect")
	ErrQuery   = errors.New("failed to call Query")
	ErrValues  = errors.New("failed to call Values")
	ErrRows    = errors.New("failed to call Rows")
	ErrClose   = errors.New("failed to call Close")
)

func NewPostgresDB(ctx context.Context, dsn string) (PGDB, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return PGDB{}, fmt.Errorf("%w: %w", ErrConnect, err)
	}

	db := PGDB{conn: conn}

	return db, nil
}

func (db PGDB) Query(ctx context.Context, sql string) (QueryResult, error) {
	rows, err := db.conn.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQuery, err)
	}
	defer rows.Close()

	var results QueryResult

	fieldDescriptions := rows.FieldDescriptions()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrValues, err)
		}

		row := make(map[string]interface{})
		for i, field := range fieldDescriptions {
			row[field.Name] = values[i]
		}

		results = append(results, row)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrRows, rowsErr)
	}

	return results, nil
}

func (db PGDB) Close(ctx context.Context) error {
	err := db.conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrClose, err)
	}

	return nil
}
