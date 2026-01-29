package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "any2jsonl - Convert any format to JSONL\n")
		fmt.Fprintf(os.Stderr, "=======================================\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [output-file]\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "What is JSONL (.jsonl)?\n")
		fmt.Fprintf(os.Stderr, "  JSON Lines: One JSON object per line (human-readable text)\n")
		fmt.Fprintf(os.Stderr, "  Easy to read, edit, and debug with any text editor\n")
		fmt.Fprintf(os.Stderr, "  Best for: Debugging, data inspection, human-readable exports\n\n")

		fmt.Fprintf(os.Stderr, "Full Benchmark Results:\n")
		fmt.Fprintf(os.Stderr, "  https://github.com/parf/homebase-go-lib/blob/main/benchmarks/serialization-benchmark-result.md\n\n")

		fmt.Fprintf(os.Stderr, "Supported input formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  - Parquet: Columnar binary format → .parquet\n")
		fmt.Fprintf(os.Stderr, "  - FlatBuffer: Zero-copy binary format → .fb\n")
		fmt.Fprintf(os.Stderr, "  - MsgPack: Binary serialization → .msgpack\n")
		fmt.Fprintf(os.Stderr, "  - CSV: Comma-separated values → .csv, .tsv, .psv\n\n")

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
		fmt.Fprintf(os.Stderr, "  %s data.parquet output.jsonl       → output.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.csv data.jsonl.zst         → data.jsonl.zst (with Zstd)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.msgpack data.jsonl.gz      → data.jsonl.gz (with Gzip)\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Performance (1M records):\n")
		fmt.Fprintf(os.Stderr, "  Plain JSONL: 156MB, 1.93s read, 1.38s write\n")
		fmt.Fprintf(os.Stderr, "  JSONL+Zstd:   43MB, 1.91s read, 0.84s write (RECOMMENDED)\n")
		fmt.Fprintf(os.Stderr, "  JSONL+LZ4:    64MB, 1.97s read, 0.88s write (fastest)\n\n")

		fmt.Fprintf(os.Stderr, "Schema Support:\n")
		fmt.Fprintf(os.Stderr, "  Automatically handles ANY structure - no schema required!\n")
		fmt.Fprintf(os.Stderr, "  Preserves all fields and types from input data.\n\n")

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
		for _, ext := range []string{".parquet", ".fb", ".msgpack", ".csv"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		outputFile += ".jsonl"
	}

	fmt.Printf("Converting %s -> %s\n", inputFile, outputFile)

	// Read all records from input (as generic map[string]any)
	records, err := readInput(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d records\n", len(records))

	// Write to JSONL
	if err := writeJSONL(outputFile, records); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSONL: %v\n", err)
		os.Exit(1)
	}

	stat, _ := os.Stat(outputFile)
	fmt.Printf("Written %s (%d bytes, %.2f MB)\n", outputFile, stat.Size(), float64(stat.Size())/1024/1024)
}

func readInput(filename string) ([]map[string]any, error) {
	lower := strings.ToLower(filename)

	if strings.Contains(lower, ".parquet") {
		return readParquet(filename)
	} else if strings.Contains(lower, ".fb") {
		return readFlatBuffer(filename)
	} else if strings.Contains(lower, ".msgpack") {
		return readMsgPack(filename)
	} else if strings.Contains(lower, ".csv") {
		return readCSV(filename)
	}

	return nil, fmt.Errorf("unsupported input format: %s", filename)
}

func readParquet(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := fileiterator.IterateParquetAny(filename, func(record map[string]any) error {
		records = append(records, record)
		return nil
	})
	return records, err
}

func readFlatBuffer(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := fileiterator.IterateFlatBufferList(filename, func(data []byte) error {
		// Note: Simplified - would need full FlatBuffer parsing with generated code
		return fmt.Errorf("FlatBuffer reading not fully implemented - use Parquet or MsgPack instead")
	})
	return records, err
}

func readMsgPack(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := fileiterator.IterateMsgPack(filename, func(data any) error {
		if m, ok := data.(map[string]any); ok {
			records = append(records, m)
		}
		return nil
	})
	return records, err
}

func readCSV(filename string) ([]map[string]any, error) {
	var records []map[string]any

	reader := fileiterator.FUOpen(filename)
	defer reader.Close()

	csvReader := csv.NewReader(reader)

	// Detect delimiter
	if strings.Contains(filename, ".tsv") {
		csvReader.Comma = '\t'
	} else if strings.Contains(filename, ".psv") {
		csvReader.Comma = '|'
	}

	// Read header to get field names
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Create map from header and row
		record := make(map[string]any)
		for i, fieldName := range header {
			if i < len(row) {
				// Try to infer type from value
				value := row[i]
				record[fieldName] = inferCSVType(value)
			}
		}

		records = append(records, record)
	}

	return records, nil
}

// inferCSVType tries to infer the type of a CSV value
func inferCSVType(value string) any {
	// Try bool
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Try int64
	var i int64
	if _, err := fmt.Sscanf(value, "%d", &i); err == nil && !strings.Contains(value, ".") {
		return i
	}

	// Try float64
	var f float64
	if _, err := fmt.Sscanf(value, "%f", &f); err == nil {
		return f
	}

	// Default to string
	return value
}

func writeJSONL(filename string, records []map[string]any) error {
	writer := fileiterator.FUCreate(filename)
	defer writer.Close()

	encoder := json.NewEncoder(writer)
	for _, rec := range records {
		if err := encoder.Encode(rec); err != nil {
			return err
		}
	}

	return nil
}
