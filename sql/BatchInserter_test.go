package sql_test

import (
	"testing"

	"github.com/parf/homebase-go-lib/sql"
)

func TestEscapeValue(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "NULL"},
		{"string simple", "hello", "'hello'"},
		{"string with single quote", "John's Pizza", "'John''s Pizza'"},
		{"string with backslash", "C:\\path\\to\\file", "'C:\\\\path\\\\to\\\\file'"},
		{"string with both", "It's\\cool", "'It''s\\\\cool'"},
		{"int", 42, "42"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},
		{"uint", uint(42), "42"},
		{"float64", 3.14, "3.140000"},
		{"bool true", true, "1"},
		{"bool false", false, "0"},
		{"empty string", "", "''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sql.EscapeValue(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBatchInserterWithSlices(t *testing.T) {
	// Test that the functions exist and have correct signatures
	// We can't test actual database operations without a real database
	t.Skip("Requires database connection - see examples for usage")
}

// Example showing expected usage with manual escaping (old style)
func ExampleBatchInserter_manual() {
	// This would require a real database connection
	// db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/dbname")
	// defer db.Close()
	//
	// insert, flush := hbsql.BatchInserter(db, "users", "id, name, email", 1000)
	// defer flush()
	//
	// for i := 0; i < 10000; i++ {
	//     // Manual escaping - UNSAFE if data comes from user input
	//     values := fmt.Sprintf("%d, \"user%d\", \"user%d@example.com\"", i, i, i)
	//     insert(values)
	// }
}

// Example showing expected usage with auto-escaping (new style - recommended)
func ExampleBatchInserter_autoEscape() {
	// This would require a real database connection
	// db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/dbname")
	// defer db.Close()
	//
	// insert, flush := hbsql.BatchInserter(db, "users", "id, name, email", 1000)
	// defer flush()
	//
	// for i := 0; i < 10000; i++ {
	//     // Auto-escaping - SAFE, handles quotes and special characters
	//     insert([]any{i, fmt.Sprintf("user%d", i), fmt.Sprintf("user%d@example.com", i)})
	// }
	//
	// // Also works with typed slices:
	// insert([]string{"value1", "value2", "value3"})
	// insert([]int{1, 2, 3})
	// insert([]any{1, "John's Pizza", 99.95, true, nil})
}
