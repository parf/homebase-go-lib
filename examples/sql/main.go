package main

import (
	"database/sql"
	"fmt"
	"log"

	hbsql "github.com/parf/homebase-go-lib/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Example connection string (update with your credentials)
	// Format: "user:password@tcp(host:port)/database"
	connectionString := "user:pass@tcp(localhost:3306)/testdb"

	fmt.Println("SQL WildQuery Example")
	fmt.Println("====================")
	fmt.Println()

	// Note: This example requires a real database connection
	// Update the connection string above with your database credentials

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		log.Println("Update the connection string in this example to test with a real database")
		return
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Printf("Cannot reach database: %v", err)
		log.Println("Make sure the database is running and credentials are correct")
		return
	}

	// Example query
	query := "SELECT * FROM users LIMIT 5"

	fmt.Printf("Executing query: %s\n\n", query)

	rows, err := hbsql.WildSqlQuery(db, query)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	fmt.Printf("Found %d rows\n\n", len(rows))

	// Print results
	for i, row := range rows {
		fmt.Printf("Row %d:\n", i+1)
		for key, value := range row {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println()
	}

	// Access specific fields
	if len(rows) > 0 {
		firstRow := rows[0]
		fmt.Println("First row fields:")
		for key := range firstRow {
			fmt.Printf("  - %s\n", key)
		}
	}
}
