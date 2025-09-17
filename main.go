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
	db, err := NewDB(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	results, err := db.Query(ctx, "select * from users limit 2;")
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format results: %v", err)
	}

	fmt.Println(string(jsonData))
}
