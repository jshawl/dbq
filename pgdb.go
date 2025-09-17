package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PGDBS struct {
	conn *pgx.Conn
}

type PGDBI interface {
	Query(ctx context.Context, sql string) (PostgresQueryResult, error)
	Close(ctx context.Context) error
}

func NewPostgresDB(ctx context.Context, dsn string) (PGDBS, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return PGDBS{}, err
	}
	db := PGDBS{conn: conn}
	return db, nil
}

type PostgresQueryResult []map[string]interface{}

func (db PGDBS) Query(ctx context.Context, sql string) (PostgresQueryResult, error) {
	rows, err := db.conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results PostgresQueryResult
	fieldDescriptions := rows.FieldDescriptions()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, field := range fieldDescriptions {
			row[field.Name] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

func (db PGDBS) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
