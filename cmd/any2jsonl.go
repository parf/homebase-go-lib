package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/parf/homebase-go-lib/fileiterator"
)

var (
	sqlFlag    = flag.String("sql", "", "SQL query to execute")
	tableFlag  = flag.String("table", "", "Table name (generates SELECT * FROM table)")
	driverFlag = flag.String("driver", "mysql", "Database driver: mysql or postgre")
	dsnFlag    = flag.String("dsn", "", "Database connection string")
)

func main() {
	flag.Parse()

	// Detect SQL mode
	isSQLMode := *sqlFlag != "" || *tableFlag != "" || *dsnFlag != "" || *driverFlag != "mysql"

	if isSQLMode {
		handleSQLMode()
	} else {
		handleFileMode()
	}
}

func handleFileMode() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "any2jsonl - Convert any format to JSONL\n")
		fmt.Fprintf(os.Stderr, "=======================================\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  File mode:  %s <input-file> [output-file|-]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  SQL mode:   %s --dsn=\"user:pass@host\" --sql=\"SELECT * FROM table\" [output-file|-]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "              Use '-' for stdout, omit for auto-generated filename\n\n")

		fmt.Fprintf(os.Stderr, "What is JSONL (.jsonl)?\n")
		fmt.Fprintf(os.Stderr, "  JSON Lines: One JSON object per line (human-readable text)\n")
		fmt.Fprintf(os.Stderr, "  Easy to read, edit, and debug with any text editor\n")
		fmt.Fprintf(os.Stderr, "  Best for: Debugging, data inspection, human-readable exports\n\n")

		fmt.Fprintf(os.Stderr, "=== SQL DATABASE SUPPORT ===\n\n")
		fmt.Fprintf(os.Stderr, "Export directly from MySQL or PostgreSQL to JSONL format:\n\n")

		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  --sql=\"SELECT * FROM table\"  SQL query to execute\n")
		fmt.Fprintf(os.Stderr, "  --table=\"schema.table\"       Table name (alternative to --sql)\n")
		fmt.Fprintf(os.Stderr, "  --driver=mysql               Database driver: mysql or postgre (default: mysql)\n")
		fmt.Fprintf(os.Stderr, "  --dsn=\"connection-string\"    Database connection string\n\n")

		fmt.Fprintf(os.Stderr, "DSN Format:\n")
		fmt.Fprintf(os.Stderr, "  MySQL:      user:password@tcp(host:3306)/database\n")
		fmt.Fprintf(os.Stderr, "  PostgreSQL: host=localhost port=5432 user=myuser password=mypass sslmode=disable\n")
		fmt.Fprintf(os.Stderr, "  Simplified: user:password@host (auto-expanded)\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s --dsn=\"root:pass@localhost\" --sql=\"SELECT * FROM users\" users.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --dsn=\"root:pass@localhost\" --table=\"mydb.users\" - | jq\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --driver=postgre --dsn=\"user:pass@pghost\" --table=\"public.orders\" orders.jsonl.zst\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "=== FILE CONVERSION ===\n\n")

		fmt.Fprintf(os.Stderr, "Full Benchmark Results:\n")
		fmt.Fprintf(os.Stderr, "  https://github.com/parf/homebase-go-lib/blob/main/benchmarks/serialization-benchmark-result.md\n\n")

		fmt.Fprintf(os.Stderr, "Supported input formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .parquet → Columnar binary format\n")
		fmt.Fprintf(os.Stderr, "  .msgpack → Binary serialization\n")
		fmt.Fprintf(os.Stderr, "  .csv     → Comma-separated values\n\n")

		fmt.Fprintf(os.Stderr, "Input compression formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .gz  → Gzip (standard compression, widely supported, slow)\n")
		fmt.Fprintf(os.Stderr, "  .zst → Zstandard (RECOMMENDED: best balance of speed & compression)\n")
		fmt.Fprintf(os.Stderr, "  .lz4 → LZ4 (fastest compression, moderate compression ratio)\n")
		fmt.Fprintf(os.Stderr, "  .br  → Brotli (best compression, but very slow)\n")
		fmt.Fprintf(os.Stderr, "  .xz  → XZ/LZMA (excellent compression, extremely slow - avoid)\n\n")

		fmt.Fprintf(os.Stderr, "Output compression (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .jsonl     → Plain text (no compression)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.gz  → Gzip compression\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.zst → Zstandard (RECOMMENDED: best balance)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.lz4 → LZ4 (fastest)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.br  → Brotli\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.xz  → XZ\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s data.csv                        → data.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.csv -                      → stdout\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.parquet output.jsonl       → output.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.csv data.jsonl.zst         → data.jsonl.zst (with Zstd)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.msgpack data.jsonl.gz      → data.jsonl.gz (with Gzip)\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Performance (1M records):\n")
		fmt.Fprintf(os.Stderr, "  Plain JSONL: 156MB, 1.93s read, 1.38s write\n")
		fmt.Fprintf(os.Stderr, "  JSONL+Zstd:   43MB, 1.91s read, 0.84s write (RECOMMENDED)\n")
		fmt.Fprintf(os.Stderr, "  JSONL+LZ4:    64MB, 1.97s read, 0.88s write (fastest)\n\n")

		fmt.Fprintf(os.Stderr, "Schema Support:\n")
		fmt.Fprintf(os.Stderr, "  ✅ Automatically handles ANY structure - no schema required!\n")
		fmt.Fprintf(os.Stderr, "  ✅ Preserves all fields and types from input data.\n\n")

		fmt.Fprintf(os.Stderr, "See also: ./any2parquet (convert to efficient Parquet format)\n")
		os.Exit(1)
	}

	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	if outputFile == "" {
		outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		// Handle double extensions
		for _, ext := range []string{".gz", ".zst", ".lz4", ".br", ".xz"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		for _, ext := range []string{".parquet", ".msgpack", ".csv"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		outputFile += ".jsonl"
	}

	// Read all records from input (supports ANY schema)
	records, err := fileiterator.ReadInput(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Read %d records\n", len(records))

	// If output is "-", write to stdout. Otherwise, write to file.
	if outputFile == "-" {
		// Write to stdout
		encoder := json.NewEncoder(os.Stdout)
		for _, record := range records {
			if err := encoder.Encode(record); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Converting %s -> %s\n", inputFile, outputFile)
		// Write to JSONL file (compression auto-detected from filename)
		if err := fileiterator.WriteOutput(outputFile, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing JSONL: %v\n", err)
			os.Exit(1)
		}
		stat, _ := os.Stat(outputFile)
		fmt.Fprintf(os.Stderr, "Written %s (%d bytes, %.2f MB)\n", outputFile, stat.Size(), float64(stat.Size())/1024/1024)
	}
}

func handleSQLMode() {
	// Validate required flags
	if *dsnFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --dsn is required for SQL queries\n")
		os.Exit(1)
	}

	if *sqlFlag == "" && *tableFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: Either --sql or --table is required\n")
		os.Exit(1)
	}

	if *sqlFlag != "" && *tableFlag != "" {
		fmt.Fprintf(os.Stderr, "Error: Cannot use both --sql and --table\n")
		os.Exit(1)
	}

	// Generate SQL query from table if needed
	sqlQuery := *sqlFlag
	if sqlQuery == "" && *tableFlag != "" {
		sqlQuery = fmt.Sprintf("SELECT * FROM %s", *tableFlag)
	}

	// Get optional output file from positional argument
	outputFile := ""
	if flag.NArg() > 0 {
		outputFile = flag.Arg(0)
	}

	// Normalize DSN
	dsn := normalizeDSN(*driverFlag, *dsnFlag)

	fmt.Fprintf(os.Stderr, "Executing SQL query: %s\n", sqlQuery)

	// Read from SQL
	records, err := fileiterator.ReadSQLInput(*driverFlag, dsn, sqlQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing SQL query: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Read %d records\n", len(records))

	// If output file is "-" or empty, write to stdout. Otherwise, write to file.
	if outputFile == "" || outputFile == "-" {
		// Write to stdout
		encoder := json.NewEncoder(os.Stdout)
		for _, record := range records {
			if err := encoder.Encode(record); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		// Write to file
		if err := fileiterator.WriteOutput(outputFile, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing JSONL: %v\n", err)
			os.Exit(1)
		}
		stat, _ := os.Stat(outputFile)
		fmt.Fprintf(os.Stderr, "Written %s (%d bytes, %.2f MB)\n", outputFile, stat.Size(), float64(stat.Size())/1024/1024)
	}
}

func normalizeDSN(driver, dsn string) string {
	// If DSN contains "/" or "=" or "sslmode=", it's already in proper format
	if strings.Contains(dsn, "/") || strings.Contains(dsn, "=") || strings.Contains(dsn, "sslmode=") {
		return dsn
	}

	// Simplified format: "user:password@host"
	if !strings.Contains(dsn, "@") {
		return dsn // Return as-is if not in simplified format
	}

	parts := strings.Split(dsn, "@")
	if len(parts) != 2 {
		return dsn
	}

	userPass := parts[0]
	host := parts[1]

	switch driver {
	case "mysql":
		// MySQL: user:password@tcp(host:3306)/
		if !strings.Contains(host, ":") {
			host = host + ":3306"
		}
		return fmt.Sprintf("%s@tcp(%s)/", userPass, host)

	case "postgre", "postgres", "postgresql":
		// PostgreSQL: host=localhost port=5432 user=myuser password=mypass sslmode=disable
		userParts := strings.Split(userPass, ":")
		user := userParts[0]
		pass := ""
		if len(userParts) > 1 {
			pass = userParts[1]
		}

		hostPort := strings.Split(host, ":")
		hostName := hostPort[0]
		port := "5432"
		if len(hostPort) > 1 {
			port = hostPort[1]
		}

		return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
			hostName, port, user, pass)
	}

	return dsn
}
