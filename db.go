package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, url string) (*DB, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	db := &DB{conn: conn}
	return db, nil
}

type QueryResult []map[string]interface{}

func (db *DB) Query(ctx context.Context, sql string) (QueryResult, error) {
	rows, err := db.conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results QueryResult
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

func (db *DB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
