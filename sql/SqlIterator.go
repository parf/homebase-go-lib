package sql

import (
	"database/sql"
	"fmt"
	"log"
	"log/syslog"
	"path/filepath"
	"runtime"

	"github.com/parf/homebase-go-lib/clistat"
	_ "github.com/go-sql-driver/mysql"
)

// sysLogError writes error to syslog
func sysLogError(message string) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Error in syslog, ", message)
		return
	}
	logwriter, _ := syslog.New(syslog.LOG_NOTICE, filepath.Base(filename))
	if logwriter != nil {
		logwriter.Err("go-api " + message)
	}
}

// SqlRowProcessor is a function type for processing SQL rows
type SqlRowProcessor func(row *sql.Rows)

// SqlIterator iterates over SQL query and shows statistics
//
// sample connection: "parf:mv700@tcp(hdb2:3306)/visits_log"
// sample sql:        "SELECT FL, T, C, B, G, V, Blocked, L FROM flTCBGVL limit 10"
func SqlIterator(connection string, sql_ string, processor SqlRowProcessor) {
	fmt.Println("Iterating SQL: " + sql_)
	db, e := sql.Open("mysql", connection)
	if e != nil {
		log.Println("SqlIterator Connection Error:", e, connection)
		sysLogError("SqlIterator Connection Error: " + e.Error() + " " + connection)
		return
	}
	defer db.Close()
	results, e := db.Query(sql_)
	if e != nil {
		log.Println("SqlIterator db Query Error:", e, sql_)
		sysLogError("SqlIterator db Query Error: " + e.Error() + " " + sql_)
		return
	}
	stat := clistat.New(10)
	defer stat.Finish()
	for results.Next() {
		processor(results)
		stat.Hit()
	}
}
