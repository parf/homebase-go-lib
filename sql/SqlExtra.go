package sql

/*
	Execute SQL Select-or-Alike Statement return list of map{ field: (string)value rows }
 */

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type (
	SqlRow  map[string]any
	SqlRows []SqlRow
)

// WildSqlQuery executes a SQL query and returns results as a slice of maps.
// Each row is represented as a map with column names as keys and string values.
func WildSqlQuery(db *sql.DB, query string) (rz SqlRows, err error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]any, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		kv := SqlRow{}
		for i, col := range values {
			if col == nil {
				kv[columns[i]] = col
				continue
			}
			kv[columns[i]] = string(col)
		}
		rz = append(rz, kv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rz, nil
}
