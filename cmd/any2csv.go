package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		fmt.Fprintf(os.Stderr, "any2csv - Convert any format to CSV\n")
		fmt.Fprintf(os.Stderr, "=====================================\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  File mode:  %s <input-file> [output-file|-]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  SQL mode:   %s --dsn=\"user:pass@host\" --sql=\"SELECT * FROM table\" [output-file|-]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "              Use '-' for stdout, omit for auto-generated filename\n\n")

		fmt.Fprintf(os.Stderr, "What is CSV (.csv)?\n")
		fmt.Fprintf(os.Stderr, "  Comma-Separated Values: Tabular data in plain text format\n")
		fmt.Fprintf(os.Stderr, "  Compatible with Excel, Google Sheets, and all spreadsheet tools\n")
		fmt.Fprintf(os.Stderr, "  Best for: Data analysis, spreadsheets, reporting, data exchange\n\n")

		fmt.Fprintf(os.Stderr, "=== SQL DATABASE SUPPORT ===\n\n")
		fmt.Fprintf(os.Stderr, "Export directly from MySQL or PostgreSQL to CSV format:\n\n")

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
		fmt.Fprintf(os.Stderr, "  %s --dsn=\"root:pass@localhost\" --sql=\"SELECT * FROM users\" users.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --dsn=\"root:pass@localhost\" --table=\"mydb.users\" - | head\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --driver=postgre --dsn=\"user:pass@pghost\" --table=\"public.orders\" orders.csv\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "=== FILE CONVERSION ===\n\n")

		fmt.Fprintf(os.Stderr, "Supported input formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .parquet → Columnar binary format\n")
		fmt.Fprintf(os.Stderr, "  .jsonl   → JSON Lines format\n")
		fmt.Fprintf(os.Stderr, "  .msgpack → Binary serialization\n\n")

		fmt.Fprintf(os.Stderr, "Input compression formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .gz  → Gzip (standard compression, widely supported, slow)\n")
		fmt.Fprintf(os.Stderr, "  .zst → Zstandard (RECOMMENDED: best balance of speed & compression)\n")
		fmt.Fprintf(os.Stderr, "  .lz4 → LZ4 (fastest compression, moderate compression ratio)\n")
		fmt.Fprintf(os.Stderr, "  .br  → Brotli (best compression, but very slow)\n")
		fmt.Fprintf(os.Stderr, "  .xz  → XZ/LZMA (excellent compression, extremely slow - avoid)\n\n")

		fmt.Fprintf(os.Stderr, "Output compression (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .csv     → Plain text (no compression)\n")
		fmt.Fprintf(os.Stderr, "  .csv.gz  → Gzip compression\n")
		fmt.Fprintf(os.Stderr, "  .csv.zst → Zstandard (RECOMMENDED: best balance)\n")
		fmt.Fprintf(os.Stderr, "  .csv.lz4 → LZ4 (fastest)\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s data.jsonl                      → data.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.parquet -                  → stdout\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.jsonl output.csv           → output.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.parquet data.csv.gz        → data.csv.gz (with Gzip)\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Schema Support:\n")
		fmt.Fprintf(os.Stderr, "  ✅ Automatically handles ANY structure - no schema required!\n")
		fmt.Fprintf(os.Stderr, "  ✅ Column order is sorted alphabetically for consistency.\n\n")

		fmt.Fprintf(os.Stderr, "See also: ./any2jsonl, ./any2parquet\n")
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
		for _, ext := range []string{".parquet", ".jsonl", ".msgpack", ".csv"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		outputFile += ".csv"
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
		if err := writeCSVToWriter(os.Stdout, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Converting %s -> %s\n", inputFile, outputFile)
		// Write to CSV file (compression auto-detected from filename)
		if err := fileiterator.WriteOutput(outputFile, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
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
		if err := writeCSVToWriter(os.Stdout, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Write to file
		if err := fileiterator.WriteOutput(outputFile, records); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
		stat, _ := os.Stat(outputFile)
		fmt.Fprintf(os.Stderr, "Written %s (%d bytes, %.2f MB)\n", outputFile, stat.Size(), float64(stat.Size())/1024/1024)
	}
}

func writeCSVToWriter(w *os.File, records []map[string]any) error {
	if len(records) == 0 {
		return nil
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Get all unique column names and sort them for consistency
	columnSet := make(map[string]bool)
	for _, record := range records {
		for key := range record {
			columnSet[key] = true
		}
	}

	columns := make([]string, 0, len(columnSet))
	for col := range columnSet {
		columns = append(columns, col)
	}
	sort.Strings(columns)

	// Write header
	if err := csvWriter.Write(columns); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		row := make([]string, len(columns))
		for i, col := range columns {
			if val, ok := record[col]; ok && val != nil {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
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
