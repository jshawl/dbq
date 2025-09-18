// package main implements a query tool for databases
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	err := run(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("%v", err)
	}
}

var (
	ErrDatabaseURLUnset = errors.New("DATABASE_URL is not set")
	ErrNewPostgresDB    = errors.New("failed to connect to database")
	ErrFormatResults    = errors.New("failed to format results")
)

func run(dsn string) error {
	ctx := context.Background()

	if dsn == "" {
		return fmt.Errorf("%w", ErrDatabaseURLUnset)
	}

	pgdb, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewPostgresDB, err)
	}

	database := NewDB(pgdb)

	defer func() {
		err := database.Close(ctx)
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	query, err := database.Query(ctx, "select * from users limit 2;")
	if err != nil {
		log.Println(err)
	}

	log.Printf("success in %s\n", query.Duration)

	jsonData, err := json.MarshalIndent(query.Results, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFormatResults, err)
	}

	log.Println(string(jsonData))

	return nil
}
