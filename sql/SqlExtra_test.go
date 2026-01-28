package sql_test

import (
	"testing"

	hbsql "github.com/parf/homebase-go-lib/sql"
	_ "github.com/go-sql-driver/mysql"
)

func TestSqlRowTypes(t *testing.T) {
	// Test that the types are defined correctly
	var row hbsql.SqlRow
	row = make(map[string]any)
	row["test"] = "value"

	if row["test"] != "value" {
		t.Error("SqlRow map should work correctly")
	}

	var rows hbsql.SqlRows
	rows = append(rows, row)

	if len(rows) != 1 {
		t.Error("SqlRows should be a slice of SqlRow")
	}
}

// TestWildSqlQueryWithMockDB would require a real database or mock
// This is a placeholder for when database testing is set up
func TestWildSqlQueryStructure(t *testing.T) {
	// Test is skipped without a real database connection
	t.Skip("Requires database connection - see examples for usage")
}

// Example test showing the expected structure
func ExampleWildSqlQuery() {
	// This would require a real database connection
	// db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/dbname")
	// defer db.Close()
	//
	// rows, err := hbsql.WildSqlQuery(db, "SELECT id, name FROM users LIMIT 10")
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// for _, row := range rows {
	//     fmt.Printf("ID: %s, Name: %s\n", row["id"], row["name"])
	// }
}
