package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	defer globalDB.close()
	results, err := globalDB.query("select * from users limit 2;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
