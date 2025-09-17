package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	db := NewDB(pgdb)

	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	query, err := db.Query(ctx, "select * from users limit 2;")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("success in %s\n", query.Duration)
	jsonData, err := json.MarshalIndent(query.Results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format results: %v", err)
	}

	fmt.Println(string(jsonData))
}
