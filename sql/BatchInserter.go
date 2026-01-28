package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

/**
 * UNSAFE SQL Batch Inserter
 * Avoid Injections !!! YOU HAVE TO ESCAPE Values YOURSELF !!
 *   maybe use https://godoc.org/github.com/golang-sql/sqlexp
 *
 * table  - table name
 * fields - comma delimited field list
 *
 * Usage:
 *
 * insert, flush := BatchInserter(sql.Open(...), "finance.invoice", "id, name, amount", 1000)
 * defer flush()
 * for ..... {
 * 		insert("1, \"Red Sox\", 99.95")   // values(...) sql-row as string - YOU have to escape it
 * }
 *
 *
 * TODO:
 *   implement this via channels and parallel go thread
 *   https://goinbigdata.com/golang-wait-for-all-goroutines-to-finish/
 *
 */

// BatchDBInserter creates a batch inserter with a new database connection.
//
// Usage:
//
//	connect_string := "parf:passwd(rxdb:3306)/visits_log"
//	insert, flush := hb.BatchDBInserter(connect_string, "tableName", "field1, field2, ...", 10000)
//	defer flush()
//
//	for .... {
//		values := fmt.Sprintf("%d, %d, %d, %d, %d", ...)
//		insert(values)
//	}
func BatchDBInserter(db, table, fields string, bufferSize int) (insert func(string), flushClose func()) {
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

// BatchInserter creates a batch inserter for an existing database connection.
// It accumulates insert values and flushes them in batches when bufferSize is reached.
func BatchInserter(db *sql.DB, table string, fields string, bufferSize int) (insert func(string), flush func()) {
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
	insert = func(row string) {
		buffer = append(buffer, row)
		cnt++
		if cnt > bufferSize {
			flush()
		}
	}
	return
}
