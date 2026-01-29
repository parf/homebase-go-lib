package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "any2jsonl - Convert any format to JSONL\n")
		fmt.Fprintf(os.Stderr, "=======================================\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [output-file]\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Full Benchmark Results:\n")
		fmt.Fprintf(os.Stderr, "  https://github.com/parf/homebase-go-lib/blob/main/benchmarks/serialization-benchmark-result.md\n\n")

		fmt.Fprintf(os.Stderr, "What is JSONL (.jsonl)?\n")
		fmt.Fprintf(os.Stderr, "  JSON Lines: One JSON object per line (human-readable text)\n")
		fmt.Fprintf(os.Stderr, "  Easy to read, edit, and debug with any text editor\n")
		fmt.Fprintf(os.Stderr, "  Best for: Debugging, data inspection, human-readable exports\n\n")

		fmt.Fprintf(os.Stderr, "Supported input formats:\n")
		fmt.Fprintf(os.Stderr, "  - Parquet: Columnar binary format (.parquet)\n")
		fmt.Fprintf(os.Stderr, "  - FlatBuffer: Zero-copy binary format (.fb, .fb.lz4)\n")
		fmt.Fprintf(os.Stderr, "  - MsgPack: Binary serialization (.msgpack, .msgpack.gz, .msgpack.zst, etc.)\n")
		fmt.Fprintf(os.Stderr, "  - CSV: Comma-separated values (.csv, .csv.gz, .csv.zst, etc.)\n\n")

		fmt.Fprintf(os.Stderr, "Output compression (auto-detected from extension):\n")
		fmt.Fprintf(os.Stderr, "  .jsonl      - Plain text (no compression)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.gz   - Gzip compression\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.zst  - Zstandard (RECOMMENDED: best balance)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.lz4  - LZ4 (fastest)\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.br   - Brotli\n")
		fmt.Fprintf(os.Stderr, "  .jsonl.xz   - XZ\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s data.csv                        → data.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.parquet output.jsonl       → output.jsonl\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.csv data.jsonl.zst         → data.jsonl.zst (with Zstd)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.msgpack data.jsonl.gz      → data.jsonl.gz (with Gzip)\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Output: JSONL (.jsonl)\n")
		fmt.Fprintf(os.Stderr, "  Human-readable text, one JSON object per line\n")
		fmt.Fprintf(os.Stderr, "  Larger and slower than binary formats, but easy to inspect\n\n")

		fmt.Fprintf(os.Stderr, "Performance (1M records):\n")
		fmt.Fprintf(os.Stderr, "  Plain JSONL: 156MB, 1.93s read, 1.38s write\n")
		fmt.Fprintf(os.Stderr, "  JSONL+Zstd:   43MB, 1.91s read, 0.84s write (RECOMMENDED)\n")
		fmt.Fprintf(os.Stderr, "  JSONL+LZ4:    64MB, 1.97s read, 0.88s write (fastest)\n\n")

		fmt.Fprintf(os.Stderr, "Note: Assumes record structure with fields:\n")
		fmt.Fprintf(os.Stderr, "  id, name, email, age, score, active, category, timestamp\n")
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

	// Read all records from input
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

type Record struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Age       int64   `json:"age"`
	Score     float64 `json:"score"`
	Active    bool    `json:"active"`
	Category  string  `json:"category"`
	Timestamp int64   `json:"timestamp"`
}

func readInput(filename string) ([]Record, error) {
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

func readParquet(filename string) ([]Record, error) {
	var records []Record

	pf, err := file.OpenParquetFile(filename, false)
	if err != nil {
		return nil, err
	}
	defer pf.Close()

	reader, err := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
	if err != nil {
		return nil, err
	}

	tbl, err := reader.ReadTable(context.Background())
	if err != nil {
		return nil, err
	}
	defer tbl.Release()

	numRows := int(tbl.NumRows())

	idCol := tbl.Column(0).Data().Chunk(0).(*array.Int64)
	nameCol := tbl.Column(1).Data().Chunk(0).(*array.String)
	emailCol := tbl.Column(2).Data().Chunk(0).(*array.String)
	ageCol := tbl.Column(3).Data().Chunk(0).(*array.Int64)
	scoreCol := tbl.Column(4).Data().Chunk(0).(*array.Float64)
	activeCol := tbl.Column(5).Data().Chunk(0).(*array.Boolean)
	categoryCol := tbl.Column(6).Data().Chunk(0).(*array.String)
	timestampCol := tbl.Column(7).Data().Chunk(0).(*array.Int64)

	for i := 0; i < numRows; i++ {
		rec := Record{
			ID:        idCol.Value(i),
			Name:      nameCol.Value(i),
			Email:     emailCol.Value(i),
			Age:       ageCol.Value(i),
			Score:     scoreCol.Value(i),
			Active:    activeCol.Value(i),
			Category:  categoryCol.Value(i),
			Timestamp: timestampCol.Value(i),
		}
		records = append(records, rec)
	}

	return records, nil
}

func readFlatBuffer(filename string) ([]Record, error) {
	var records []Record
	err := fileiterator.IterateFlatBufferList(filename, func(data []byte) error {
		// Note: Simplified - would need full FlatBuffer parsing with generated code
		return fmt.Errorf("FlatBuffer reading not fully implemented - use Parquet or MsgPack instead")
	})
	return records, err
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

func writeJSONL(filename string, records []Record) error {
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
