// package main implements a query tool for databases
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pgdb, err := NewPostgresDB(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
		log.Printf("Failed to format results: %v", err)
	}

	log.Println(string(jsonData))
}
