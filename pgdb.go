package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PGDB struct {
	conn *pgx.Conn
}

func NewPostgresDB(ctx context.Context, dsn string) (PGDB, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return PGDB{}, err
	}
	db := PGDB{conn: conn}
	return db, nil
}

type PostgresQueryResult []map[string]interface{}

func (db *PGDB) Query(ctx context.Context, sql string) (PostgresQueryResult, error) {
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

func (db *PGDB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
