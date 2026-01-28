# Changelog

## 2026-01-28 - Initial Release

### Added from RDUtil

All utilities from `/rd/go/src/RDUtil/` have been integrated into the `hb` package:

#### Core Files Added

1. **BatchInserter.go** - Batch SQL insert utilities
   - `BatchInserter()` - Batch inserter with existing DB connection
   - `BatchDBInserter()` - Batch inserter with new DB connection

2. **Debug.go** - Configurable debugging utilities
   - `Debug()` - Debug output to STDERR
   - `DebugLog()` - Debug output with timestamps

3. **General.go** - General utilities
   - `Any2uint32()` - Type conversion utility
   - `DumpSortedMap()` - Sorted map printing
   - `Scale()` - Logarithmic base-4 scaling (moved from Scale.go)

4. **GZLoaders.go** - Gzip file processing
   - `FUOpen()` - Universal file/URL opener
   - `LoadLinesGzFile()` - Process gzipped text files line by line
   - `LoadBinGzFile()` - Load gzipped binary files
   - `LoadIDTabGzFile()` - Process tab-separated gzipped files

5. **JobScheduler.go** - Periodic job scheduling
   - `NewJobScheduler()` - Create job scheduler
   - `Start()` / `Stop()` - Control scheduler
   - `IsRunning()` - Check scheduler status

6. **Runner.go** - Task runners with performance tracking
   - `NewParallelRunner()` - Parallel task execution
   - `NewSequentialRunner()` - Sequential task execution
   - `MemReport()` - Memory usage reporting

7. **SqlIterator.go** - SQL query iteration
   - `SqlIterator()` - Iterate over SQL queries with statistics

8. **Syslog.go** - System logging utilities
   - `SysLogNotice()` - Write notice to syslog
   - `SysLogError()` - Write error to syslog

9. **Zlib.go** - Zlib compression utilities
   - `ZlibFileIterator()` - Iterate over zlib-compressed binary files

#### Packages Added

- **clistat/** - CLI statistics tracking package
  - Moved from `/rd/go/src/Clistat/`
  - `New()` - Create new statistics tracker
  - `Hit()` - Record hit
  - `Finish()` - Print final statistics

- **sql/** - SQL utilities package
  - Moved from `/rd/go/src/RDUtil/sql/`
  - `WildSqlQuery()` - Execute SQL and return results as maps
  - `SqlRow` / `SqlRows` types for dynamic query results

#### Tests Added

- `Debug_test.go` - Debug function tests
- `General_test.go` - General utilities tests (including Scale)
- `General_example_test.go` - Example tests for General functions
- `JobScheduler_test.go` - Job scheduler tests
- `Runner_test.go` - Task runner tests
- `clistat/clistat_test.go` - Clistat package tests

#### Examples Added

- `examples/clistat/` - Clistat usage example
- `examples/scale/` - Scale function demonstration
- `examples/sql/` - SQL WildSqlQuery example

#### Documentation Added

- `docs/API.md` - Comprehensive API documentation
- `docs/CHANGELOG.md` - This file
- `clistat/README.md` - Clistat package documentation
- `sql/README.md` - SQL package documentation

### Changes

- Package name changed from `RDUtil` to `hb`
- Updated imports to use `github.com/parf/homebase-go-lib`
- Updated `interface{}` to `any` (Go 1.18+)
- Updated Clistat imports to use local package
- Added proper Go module support (go 1.25.6)

### Dependencies

- `github.com/go-sql-driver/mysql` v1.9.3 - MySQL driver for database utilities
- `filippo.io/edwards25519` v1.1.0 - Indirect dependency

### File Organization

```
.
├── BatchInserter.go
├── Debug.go
├── General.go
├── GZLoaders.go
├── JobScheduler.go
├── Runner.go
├── SqlIterator.go
├── Syslog.go
├── Zlib.go
├── homebase.go
├── clistat/
│   ├── clistat.go
│   ├── clistat_test.go
│   └── README.md
├── sql/
│   ├── SqlExtra.go
│   ├── SqlExtra_test.go
│   └── README.md
├── examples/
│   ├── clistat/main.go
│   ├── scale/main.go
│   └── sql/main.go
├── docs/
│   ├── API.md
│   ├── CHANGELOG.md
│   └── README.md
└── testdata/
```

### Test Coverage

All main packages have comprehensive test coverage:
- Core package: 13 tests
- Clistat package: 4 tests
- SQL package: 2 tests
- Total: 19 tests, all passing

### Breaking Changes

None - this is the initial release.

### Migration from RDUtil

If migrating from the original RDUtil package:

```go
// Old
import "difive.com/lib/RDUtil"
RDUtil.Scale(100)

// New
import hb "github.com/parf/homebase-go-lib"
hb.Scale(100)
```

Clistat package:
```go
// Old
import "difive.com/lib/Clistat"
stat := Clistat.New(10)

// New
import "github.com/parf/homebase-go-lib/clistat"
stat := clistat.New(10)
```
