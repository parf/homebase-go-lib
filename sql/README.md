# SQL Package

The `sql` package provides utilities for working with SQL databases.

## SqlExtra

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
