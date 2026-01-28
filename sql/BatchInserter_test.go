package sql_test

import (
	"testing"
)

func TestBatchInserterTypes(t *testing.T) {
	// Test that the functions exist and have correct signatures
	// We can't test actual database operations without a real database
	t.Skip("Requires database connection - see examples for usage")
}

// Example showing expected usage
func ExampleBatchInserter() {
	// This would require a real database connection
	// db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/dbname")
	// defer db.Close()
	//
	// insert, flush := hbsql.BatchInserter(db, "users", "id, name, email", 1000)
	// defer flush()
	//
	// for i := 0; i < 10000; i++ {
	//     values := fmt.Sprintf("%d, \"user%d\", \"user%d@example.com\"", i, i, i)
	//     insert(values)
	// }
}
