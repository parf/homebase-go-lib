# SQL Package

The `sql` package provides utilities for working with SQL databases.

## Modules

### BatchInserter - Batch SQL Inserts

**WARNING**: UNSAFE - You must escape values yourself to prevent SQL injection!

Efficiently insert large batches of data into SQL databases.

```go
import hbsql "github.com/parf/homebase-go-lib/sql"

db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/dbname")
defer db.Close()

insert, flush := hbsql.BatchInserter(db, "users", "id, name, email", 1000)
defer flush()

for i := 0; i < 10000; i++ {
    // WARNING: You must escape values yourself!
    values := fmt.Sprintf("%d, \"%s\", \"%s\"", i, escapeString(name), escapeString(email))
    insert(values)
}
```

**Also available:** `BatchDBInserter` - Opens database connection for you.

### SqlIterator - Query Iteration with Statistics

Iterate over SQL query results with automatic progress tracking.

```go
import hbsql "github.com/parf/homebase-go-lib/sql"

hbsql.SqlIterator("user:pass@tcp(host:3306)/db", "SELECT id, name FROM users", func(row *sql.Rows) {
    var id int
    var name string
    row.Scan(&id, &name)
    fmt.Printf("User: %d - %s\n", id, name)
})
```

### SqlExtra - Dynamic Query Results

Execute SQL SELECT statements and get results as maps.

### Types

```go
type SqlRow  map[string]any        // Single row as map[column]value
type SqlRows []SqlRow               // Multiple rows
```

### Functions

#### WildSqlQuery

```go
func WildSqlQuery(db *sql.DB, query string) (SqlRows, error)
```

Executes a SQL query and returns results as a slice of maps. Each row is represented as a map with column names as keys and string values.

**Note**: All values are returned as strings. NULL values are returned as nil.

## Usage Example

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    hbsql "github.com/parf/homebase-go-lib/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Execute query and get results as maps
    rows, err := hbsql.WildSqlQuery(db, "SELECT id, name, email FROM users LIMIT 10")
    if err != nil {
        log.Fatal(err)
    }

    // Process results
    for _, row := range rows {
        fmt.Printf("ID: %s, Name: %s, Email: %s\n",
            row["id"], row["name"], row["email"])
    }

    // Access specific row
    if len(rows) > 0 {
        firstRow := rows[0]
        fmt.Printf("First user ID: %s\n", firstRow["id"])
    }
}
```

## Benefits

- **Dynamic queries**: Works with any SELECT query without knowing the schema in advance
- **Simple interface**: Results as maps are easy to work with
- **Type flexibility**: All values returned as strings, easily convertible
- **NULL handling**: NULL values are properly handled

## Limitations

- All non-NULL values are converted to strings
- For type-specific handling, use standard `database/sql` with explicit type scanning
- Not optimized for very large result sets (loads all rows into memory)

## Use Cases

- Dynamic query results where schema is not known at compile time
- Quick prototyping and debugging
- Administration tools
- Configuration queries
- Metadata queries
