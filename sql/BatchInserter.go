package sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

/**
 * SQL Batch Inserter with Auto-Escaping Support
 *
 * table  - table name
 * fields - comma delimited field list
 *
 * Usage (SAFE - with auto-escaping):
 *
 * insert, flush := BatchInserter(sql.Open(...), "finance.invoice", "id, name, amount", 1000)
 * defer flush()
 * for ..... {
 * 		insert([]any{1, "Red Sox", 99.95})  // Values are automatically escaped
 * }
 *
 * Usage (UNSAFE - manual escaping for backward compatibility):
 *
 * for ..... {
 * 		insert("1, 'Red Sox', 99.95")  // String format - YOU must escape values yourself
 * }
 *
 * TODO:
 *   implement this via channels and parallel go thread
 *   https://goinbigdata.com/golang-wait-for-all-goroutines-to-finish/
 *
 */

// BatchDBInserter creates a batch inserter with a new database connection.
//
// Usage with auto-escaping (recommended):
//
//	connect_string := "parf:passwd(rxdb:3306)/visits_log"
//	insert, flush := sql.BatchDBInserter(connect_string, "tableName", "field1, field2, ...", 10000)
//	defer flush()
//
//	for .... {
//		insert([]any{val1, val2, val3, val4, val5})  // Auto-escaped
//	}
//
// Usage with manual string (legacy):
//
//	for .... {
//		values := fmt.Sprintf("%d, %d, %d, %d, %d", ...)  // You must escape
//		insert(values)
//	}
func BatchDBInserter(db, table, fields string, bufferSize int) (insert func(any), flushClose func()) {
	_db, err := sql.Open("mysql", db)
	if err != nil {
		panic(err)
	}
	insert, flush := BatchInserter(_db, table, fields, bufferSize)
	flushClose = func() {
		flush()
		_db.Close()
	}
	return
}

// EscapeValue converts a value to a SQL-safe string representation.
// Handles NULL, strings (with escaping), numbers, booleans, and other types.
// This function is exported so users can use it for manual escaping if needed.
func EscapeValue(val any) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case string:
		// Escape single quotes by doubling them and backslashes
		escaped := strings.ReplaceAll(v, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "''")
		return "'" + escaped + "'"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		// For other types, convert to string and escape
		str := fmt.Sprintf("%v", v)
		escaped := strings.ReplaceAll(str, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "''")
		return "'" + escaped + "'"
	}
}

// valuesToString converts values (string or slice) to a SQL VALUES string.
// If values is a string, returns it as-is (unsafe, user must escape).
// If values is a slice/array, escapes each element and joins with commas (safe).
func valuesToString(values any) string {
	// If it's a string, use as-is (backward compatible)
	if str, ok := values.(string); ok {
		return str
	}

	// Use reflection to handle any slice type
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(fmt.Sprintf("BatchInserter: values must be string, slice, or array (got %T)", values))
	}

	// Extract slice elements and escape them
	length := v.Len()
	if length == 0 {
		return ""
	}

	escaped := make([]string, length)
	for i := 0; i < length; i++ {
		escaped[i] = EscapeValue(v.Index(i).Interface())
	}

	return strings.Join(escaped, ", ")
}

// BatchInserter creates a batch inserter for an existing database connection.
// It accumulates insert values and flushes them in batches when bufferSize is reached.
//
// The insert function accepts either:
//   - string: values as SQL string (UNSAFE - you must escape yourself)
//   - slice/array: values as slice (SAFE - automatically escaped)
//
// Example with auto-escaping:
//
//	insert([]any{1, "John's Pizza", 99.95})
//	insert([]string{"value1", "value2", "value3"})
//	insert([]int{1, 2, 3})
//
// Example with manual escaping (backward compatible):
//
//	insert("1, 'John''s Pizza', 99.95")
func BatchInserter(db *sql.DB, table string, fields string, bufferSize int) (insert func(any), flush func()) {
	buffer := []string{}
	sql_prefix := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", table, fields)
	cnt := 0
	flush = func() {
		bl := len(buffer)
		if bl == 0 {
			return
		}
		sq := "" // sql text
		if bl == 1 {
			sq = "(" + buffer[0] + ")"
		} else {
			last, buffer := buffer[bl-1], buffer[:bl-1]
			sq = "(" + strings.Join(buffer, "),(") + "),(" + last + ")"
		}
		sqlt := sql_prefix + sq
		// fmt.Println("SQL: ", sqlt)
		_, err := db.Exec(sqlt)
		if err != nil {
			fmt.Println("Insert Error. Table: " + table)
			panic(err)
		}
		buffer = []string{}
		cnt = 0
	}
	insert = func(values any) {
		row := valuesToString(values)
		buffer = append(buffer, row)
		cnt++
		if cnt > bufferSize {
			flush()
		}
	}
	return
}
