package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

var globalDB *DB

func init() {
	var err error
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	globalDB = &DB{conn: conn}
}

type QueryResult []map[string]interface{}

func (db *DB) query(sql string) (QueryResult, error) {
	ctx := context.Background()
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

func (db *DB) close() {
	db.conn.Close(context.Background())
}
