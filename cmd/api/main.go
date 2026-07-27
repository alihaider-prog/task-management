package main

import (
	// "database/sql"
	"fmt"
	"log"

	"task-managment/internal/database"

	_ "github.com/lib/pq" // To register the driver.
)

func main() {
	// Initialize DB connection
	database.Init(".env")

	// Use db.DB for database operations
	rows, err := database.DB.Query("SELECT NOW()")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var now string
	for rows.Next() {
		if err := rows.Scan(&now); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		fmt.Println("Current time from DB:", now)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("rows failed: %v", err)
	}

	fmt.Println("connected")
}
