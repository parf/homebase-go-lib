# API Documentation

## Package: hb (github.com/parf/homebase-go-lib)

### General Utilities

#### Scale(nn uint32) byte
Logarithmic base-4 scaling function that maps uint32 values to 0-9 scale.

```go
scale := hb.Scale(1024)  // Returns 6
```

Scale ranges:
- 0 → 0
- 1 → 1
- 4 → 2
- 16 → 3
- 64 → 4
- 256 → 5
- 1024 → 6
- 4096 → 7
- 16384 → 8
- 16385+ → 9

#### Any2uint32(iii any) (uint32, error)
Convert various integer types to uint32.

```go
val, err := hb.Any2uint32(someInt)
```

#### DumpSortedMap(m map[string]any)
Print a map in sorted key order.

```go
hb.DumpSortedMap(myMap)
```

### Performance Monitoring

#### Clistat (package: clistat)
CLI statistics tracker for monitoring operations.

```go
import "github.com/parf/homebase-go-lib/clistat"

stat := clistat.New(5) // 5 second reporting interval
for i := 0; i < 1000000; i++ {
    // do work
    stat.Hit()
}
stat.Finish()
```

#### ParallelRunner / SequentialRunner
Task runners with performance tracking.

```go
runner := hb.NewParallelRunner()
runner.Run("task1", func() {
    // do work
})
runner.Run("task2", func() {
    // do work
})
runner.Finish()
```

#### MemReport(event string)
Report memory allocation statistics.

```go
hb.MemReport("After loading data")
```

### Scheduling

#### JobScheduler
Periodic job scheduler.

```go
scheduler := hb.NewJobScheduler(60, func() {
    // job to run every 60 seconds
})
scheduler.Start()
defer scheduler.Stop()
```

### Database

#### BatchInserter(db *sql.DB, table, fields string, bufferSize int)
Batch insert for SQL operations.

**WARNING**: This is UNSAFE - you must escape values yourself!

```go
insert, flush := hb.BatchInserter(db, "users", "id, name, email", 1000)
defer flush()

for _, user := range users {
    values := fmt.Sprintf("%d, \"%s\", \"%s\"", user.ID, user.Name, user.Email)
    insert(values)
}
```

#### SqlIterator(connection, sql string, processor SqlRowProcessor)
Iterate over SQL query results with statistics.

```go
hb.SqlIterator("user:pass@tcp(host:3306)/db", "SELECT * FROM users", func(row *sql.Rows) {
    var id int
    var name string
    row.Scan(&id, &name)
    // process row
})
```

### File Processing

#### FUOpen(file_or_url string) io.ReadCloser
Open file or HTTP URL.

```go
r := hb.FUOpen("http://example.com/data.txt")
defer r.Close()
```

#### LoadLinesGzFile(filename string, processor func(string))
Process lines in gzipped file.

```go
hb.LoadLinesGzFile("data.txt.gz", func(line string) {
    // process line
})
```

#### LoadBinGzFile(filename string, dest *[]byte)
Load gzipped file into byte buffer.

```go
var data []byte
hb.LoadBinGzFile("data.bin.gz", &data)
```

#### ZlibFileIterator(filename string, recordSize int, processor func([]byte))
Iterate over zlib-compressed binary records.

```go
hb.ZlibFileIterator("data.zlib", 10, func(record []byte) {
    // process 10-byte record
})
```

### Debugging

#### Debug(prefix string, level int) DebugFunction
Create debug function with output to stderr.

```go
debug := hb.Debug("MyApp", 3)
debug(0, "Error: %s", err)     // Always shown
debug(1, "Warning: %s", msg)   // Shown if level >= 1
debug(3, "Debug: %v", data)    // Shown if level >= 3
debug(4, "Trace: %v", detail)  // Not shown (level 4 > 3)
```

Debug levels:
- -1: FATAL (panic after message)
- 0: ERROR (always shown)
- 1: WARNING
- 2: NOTICE
- 3: INFO
- 4+: DEBUG

#### DebugLog(prefix string, level int) DebugFunction
Same as Debug but with timestamp prefix (uses log package).

```go
debug := hb.DebugLog("MyApp", 2)
debug(1, "Starting process")
```

### System Logging

#### SysLogNotice(message string)
Write notice to syslog.

```go
hb.SysLogNotice("Application started")
```

#### SysLogError(message string)
Write error to syslog.

```go
hb.SysLogError("Failed to connect: " + err.Error())
```

## Package: sql (github.com/parf/homebase-go-lib/sql)

### SqlExtra - Dynamic SQL Query Results

#### Types

```go
type SqlRow  map[string]any        // Single row as map[column]value
type SqlRows []SqlRow               // Multiple rows
```

#### WildSqlQuery(db *sql.DB, query string) (SqlRows, error)

Execute SQL query and return results as maps. Perfect for dynamic queries where the schema is not known at compile time.

```go
import hbsql "github.com/parf/homebase-go-lib/sql"

db, _ := sql.Open("mysql", "user:pass@tcp(host:3306)/db")
rows, err := hbsql.WildSqlQuery(db, "SELECT * FROM users LIMIT 10")
if err != nil {
    log.Fatal(err)
}

for _, row := range rows {
    fmt.Printf("ID: %s, Name: %s\n", row["id"], row["name"])
}
```

**Features:**
- Works with any SELECT query
- Returns column names dynamically
- All values as strings (easily convertible)
- NULL values handled as nil

**Note:** All non-NULL values are returned as strings. For type-specific handling, use standard database/sql with explicit scanning.
