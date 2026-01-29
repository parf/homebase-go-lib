package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "any2fb - Convert to FlatBuffer (FIXED SCHEMA ONLY)\n")
		fmt.Fprintf(os.Stderr, "====================================================\n\n")
		fmt.Fprintf(os.Stderr, "⚠️  LIMITATION: FlatBuffer requires generated code for each schema.\n")
		fmt.Fprintf(os.Stderr, "   This tool only supports a FIXED schema (id, name, email, age, score, active, category, timestamp).\n")
		fmt.Fprintf(os.Stderr, "   For ANY schema support, use ./any2parquet instead!\n\n")

		fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [output-file]\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "What is FlatBuffer (.fb)?\n")
		fmt.Fprintf(os.Stderr, "  FlatBuffer is a binary format with zero-copy deserialization.\n")
		fmt.Fprintf(os.Stderr, "  Fast reads (0.06s for 1M records), but large files (160MB uncompressed).\n")
		fmt.Fprintf(os.Stderr, "  Best for: Hot data paths where read speed is critical.\n\n")

		fmt.Fprintf(os.Stderr, "Full Benchmark Results:\n")
		fmt.Fprintf(os.Stderr, "  https://github.com/parf/homebase-go-lib/blob/main/benchmarks/serialization-benchmark-result.md\n\n")

		fmt.Fprintf(os.Stderr, "Supported input formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  - JSONL: JSON Lines, one JSON object per line → .jsonl\n")
		fmt.Fprintf(os.Stderr, "  - CSV: Comma-separated values with header row → .csv, .tsv, .psv\n")
		fmt.Fprintf(os.Stderr, "  - MsgPack: Binary serialization format → .msgpack\n")
		fmt.Fprintf(os.Stderr, "  - Parquet: Columnar storage format (RECOMMENDED) → .parquet\n\n")

		fmt.Fprintf(os.Stderr, "Input compression formats (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .gz  → Gzip (standard compression, slow)\n")
		fmt.Fprintf(os.Stderr, "  .zst → Zstandard (best balance of speed and compression)\n")
		fmt.Fprintf(os.Stderr, "  .lz4 → LZ4 (fastest compression, larger files)\n")
		fmt.Fprintf(os.Stderr, "  .br  → Brotli (best compression, very slow)\n")
		fmt.Fprintf(os.Stderr, "  .xz  → XZ/LZMA (excellent compression, extremely slow - avoid)\n\n")

		fmt.Fprintf(os.Stderr, "Output compression (recognized extension → format):\n")
		fmt.Fprintf(os.Stderr, "  .fb     → FlatBuffer (uncompressed, 160MB for 1M records)\n")
		fmt.Fprintf(os.Stderr, "  .fb.lz4 → FlatBuffer + LZ4 (RECOMMENDED: 66MB, fastest)\n")
		fmt.Fprintf(os.Stderr, "  .fb.zst → FlatBuffer + Zstandard (better compression)\n")
		fmt.Fprintf(os.Stderr, "  .fb.gz  → FlatBuffer + Gzip\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s data.jsonl.gz                  → data.fb\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.csv.zst output.fb         → output.fb\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.parquet data.fb.lz4       → data.fb.lz4 (with LZ4 - RECOMMENDED)\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Fixed Schema (required fields):\n")
		fmt.Fprintf(os.Stderr, "  id (int64), name (string), email (string), age (int64)\n")
		fmt.Fprintf(os.Stderr, "  score (float64), active (bool), category (string), timestamp (int64)\n\n")

		fmt.Fprintf(os.Stderr, "⚠️  For generic schema support, use ./any2parquet instead!\n")
		fmt.Fprintf(os.Stderr, "See also: ./any2parquet (supports ANY schema with auto-inference)\n")
		os.Exit(1)
	}

	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	if outputFile == "" {
		outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		// Handle double extensions like .jsonl.gz
		for _, ext := range []string{".gz", ".zst", ".lz4", ".br", ".xz"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		for _, ext := range []string{".jsonl", ".csv", ".msgpack", ".parquet"} {
			outputFile = strings.TrimSuffix(outputFile, ext)
		}
		outputFile += ".fb"
	}

	fmt.Printf("Converting %s -> %s\n", inputFile, outputFile)

	// Read all records from input (using fixed schema)
	records, err := readInput(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d records\n", len(records))

	// Write to FlatBuffer List (compression auto-detected from filename)
	if err := writeFlatBufferList(outputFile, records); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing FlatBuffer: %v\n", err)
		os.Exit(1)
	}

	stat, _ := os.Stat(outputFile)
	fmt.Printf("Written %s (%d bytes, %.2f MB)\n", outputFile, stat.Size(), float64(stat.Size())/1024/1024)
}

// Record with FIXED schema for FlatBuffer
type Record struct {
	ID        int64
	Name      string
	Email     string
	Age       int64
	Score     float64
	Active    bool
	Category  string
	Timestamp int64
}

func readInput(filename string) ([]Record, error) {
	lower := strings.ToLower(filename)

	if strings.Contains(lower, ".parquet") {
		return readParquet(filename)
	} else if strings.Contains(lower, ".msgpack") {
		return readMsgPack(filename)
	} else if strings.Contains(lower, ".csv") {
		return readCSV(filename)
	} else if strings.Contains(lower, ".jsonl") {
		return readJSONL(filename)
	}

	return nil, fmt.Errorf("unsupported input format: %s", filename)
}

func readJSONL(filename string) ([]Record, error) {
	var records []Record
	err := fileiterator.IterateJSONL(filename, func(line map[string]any) error {
		rec := Record{
			ID:        int64(line["id"].(float64)),
			Name:      line["name"].(string),
			Email:     line["email"].(string),
			Age:       int64(line["age"].(float64)),
			Score:     line["score"].(float64),
			Active:    line["active"].(bool),
			Category:  line["category"].(string),
			Timestamp: int64(line["timestamp"].(float64)),
		}
		records = append(records, rec)
		return nil
	})
	return records, err
}

func readCSV(filename string) ([]Record, error) {
	var records []Record

	reader := fileiterator.FUOpen(filename)
	defer reader.Close()

	csvReader := csv.NewReader(reader)

	// Detect delimiter
	if strings.Contains(filename, ".tsv") {
		csvReader.Comma = '\t'
	} else if strings.Contains(filename, ".psv") {
		csvReader.Comma = '|'
	}

	// Skip header
	_, err := csvReader.Read()
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

		var rec Record
		fmt.Sscanf(row[0], "%d", &rec.ID)
		rec.Name = row[1]
		rec.Email = row[2]
		fmt.Sscanf(row[3], "%d", &rec.Age)
		fmt.Sscanf(row[4], "%f", &rec.Score)
		rec.Active = row[5] == "true"
		rec.Category = row[6]
		fmt.Sscanf(row[7], "%d", &rec.Timestamp)

		records = append(records, rec)
	}

	return records, nil
}

func readMsgPack(filename string) ([]Record, error) {
	var records []Record
	err := fileiterator.IterateMsgPack(filename, func(data any) error {
		m := data.(map[string]any)
		rec := Record{
			ID:        m["id"].(int64),
			Name:      m["name"].(string),
			Email:     m["email"].(string),
			Age:       m["age"].(int64),
			Score:     m["score"].(float64),
			Active:    m["active"].(bool),
			Category:  m["category"].(string),
			Timestamp: m["timestamp"].(int64),
		}
		records = append(records, rec)
		return nil
	})
	return records, err
}

func readParquet(filename string) ([]Record, error) {
	var records []Record
	err := fileiterator.IterateParquetAny(filename, func(data map[string]any) error {
		rec := Record{
			ID:        data["id"].(int64),
			Name:      data["name"].(string),
			Email:     data["email"].(string),
			Age:       data["age"].(int64),
			Score:     data["score"].(float64),
			Active:    data["active"].(bool),
			Category:  data["category"].(string),
			Timestamp: data["timestamp"].(int64),
		}
		records = append(records, rec)
		return nil
	})
	return records, err
}

func writeFlatBufferList(filename string, records []Record) error {
	// Build individual record FlatBuffers
	var recordBytes [][]byte

	for _, rec := range records {
		builder := flatbuffers.NewBuilder(0)

		// Create strings
		nameOffset := builder.CreateString(rec.Name)
		emailOffset := builder.CreateString(rec.Email)
		categoryOffset := builder.CreateString(rec.Category)

		// Build TestRecord (fixed schema)
		builder.StartObject(8)
		builder.PrependInt64Slot(0, rec.ID, 0)
		builder.PrependUOffsetTSlot(1, nameOffset, 0)
		builder.PrependUOffsetTSlot(2, emailOffset, 0)
		builder.PrependInt32Slot(3, int32(rec.Age), 0)
		builder.PrependFloat64Slot(4, rec.Score, 0)
		builder.PrependBoolSlot(5, rec.Active, false)
		builder.PrependUOffsetTSlot(6, categoryOffset, 0)
		builder.PrependInt64Slot(7, rec.Timestamp, 0)
		builder.Finish(builder.EndObject())

		recordBytes = append(recordBytes, builder.FinishedBytes())
	}

	// Save as FlatBuffer list (compression auto-detected from filename)
	return fileiterator.SaveFlatBufferList(filename, recordBytes)
}
