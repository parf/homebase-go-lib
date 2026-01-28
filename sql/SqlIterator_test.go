package sql_test

import (
	"testing"
)

func TestSqlIteratorTypes(t *testing.T) {
	// Test that the functions exist and have correct signatures
	// We can't test actual database operations without a real database
	t.Skip("Requires database connection - see examples for usage")
}

// Example showing expected usage
func ExampleSqlIterator() {
	// This would require a real database connection
	// hbsql.SqlIterator("user:pass@tcp(host:3306)/db", "SELECT * FROM users LIMIT 10", func(row *sql.Rows) {
	//     var id int
	//     var name string
	//     row.Scan(&id, &name)
	//     fmt.Printf("User: %d - %s\n", id, name)
	// })
}
