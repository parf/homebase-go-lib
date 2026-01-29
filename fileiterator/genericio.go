package fileiterator

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

// ReadInput reads any supported format and returns generic records
func ReadInput(filename string) ([]map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	// Remove compression extension if present
	if ext == ".gz" || ext == ".zst" || ext == ".lz4" || ext == ".br" || ext == ".xz" {
		base := strings.TrimSuffix(filename, ext)
		ext = strings.ToLower(filepath.Ext(base))
	}

	switch ext {
	case ".jsonl", ".ndjson":
		return readJSONL(filename)
	case ".msgpack", ".mp":
		return readMsgPack(filename)
	case ".csv":
		return readCSV(filename)
	case ".parquet":
		return readParquetGeneric(filename)
	default:
		return nil, fmt.Errorf("unsupported input format: %s (supported: .jsonl, .csv, .msgpack, .parquet)", ext)
	}
}

func readJSONL(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := IterateJSONL(filename, func(record map[string]any) error {
		records = append(records, record)
		return nil
	})
	return records, err
}

func readMsgPack(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := IterateMsgPack(filename, func(data any) error {
		if m, ok := data.(map[string]any); ok {
			records = append(records, m)
		}
		return nil
	})
	return records, err
}

func readCSV(filename string) ([]map[string]any, error) {
	r := FUOpen(filename)
	defer r.Close()

	csvReader := csv.NewReader(r)

	// Read header
	headerRow, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Copy header to avoid issues with ReuseRecord
	header := make([]string, len(headerRow))
	copy(header, headerRow)

	var records []map[string]any
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		record := make(map[string]any)
		for i, value := range row {
			if i >= len(header) {
				continue
			}
			record[header[i]] = inferValue(value)
		}
		records = append(records, record)
	}

	return records, nil
}

func inferValue(s string) any {
	// Try bool
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Try int
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	// Try float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// Default to string
	return s
}

func readParquetGeneric(filename string) ([]map[string]any, error) {
	var records []map[string]any
	err := IterateParquetAny(filename, func(record map[string]any) error {
		records = append(records, record)
		return nil
	})
	return records, err
}

// WriteOutput writes records to any supported format
func WriteOutput(filename string, records []map[string]any) error {
	ext := strings.ToLower(filepath.Ext(filename))
	// Remove compression extension if present
	baseFilename := filename
	if ext == ".gz" || ext == ".zst" || ext == ".lz4" || ext == ".br" || ext == ".xz" {
		baseFilename = strings.TrimSuffix(filename, ext)
		ext = strings.ToLower(filepath.Ext(baseFilename))
	}

	switch ext {
	case ".jsonl", ".ndjson":
		return writeJSONL(filename, records)
	case ".msgpack", ".mp":
		return writeMsgPack(filename, records)
	case ".csv":
		return writeCSV(filename, records)
	case ".parquet":
		return WriteParquetAny(filename, records)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: .jsonl, .csv, .msgpack, .parquet)", ext)
	}
}

func writeJSONL(filename string, records []map[string]any) error {
	w := FUCreate(filename)
	defer w.Close()

	encoder := json.NewEncoder(w)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}
	return nil
}

func writeMsgPack(filename string, records []map[string]any) error {
	w := FUCreate(filename)
	defer w.Close()

	encoder := msgpack.NewEncoder(w)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}
	return nil
}

func writeCSV(filename string, records []map[string]any) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to write")
	}

	w := FUCreate(filename)
	defer w.Close()

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Get headers from first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}

	// Write header
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		row := make([]string, len(headers))
		for i, key := range headers {
			row[i] = fmt.Sprintf("%v", record[key])
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}
