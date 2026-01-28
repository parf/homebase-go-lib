package fileiterator

import (
	"encoding/csv"
	"fmt"
	"io"

	hb "github.com/parf/homebase-go-lib"
)

// CSVOptions configures CSV parsing behavior
type CSVOptions struct {
	Comma            rune // Field delimiter (default: ',')
	Comment          rune // Comment character (default: 0, disabled)
	SkipHeader       bool // Skip first row as header
	TrimLeadingSpace bool // Trim leading space in fields
}

// DefaultCSVOptions returns default CSV parsing options
func DefaultCSVOptions() CSVOptions {
	return CSVOptions{
		Comma:            ',',
		Comment:          0,
		SkipHeader:       false,
		TrimLeadingSpace: false,
	}
}

// IterateCSV processes a CSV file row by row.
// Supports compression auto-detection by extension (.gz, .zst).
// Each row is passed to the processor function as a slice of strings.
//
// filename - "filename" or "http://url" (with optional .gz or .zst extension)
// opts - CSV parsing options (use DefaultCSVOptions() for defaults)
// processor - function that receives each CSV row
//
// Example:
//
//	fileiterator.IterateCSV("data.csv.gz", fileiterator.DefaultCSVOptions(), func(row []string) error {
//	    fmt.Printf("Row: %v\n", row)
//	    return nil
//	})
func IterateCSV(filename string, opts CSVOptions, processor func([]string) error) error {
	fi := hb.FUOpen(filename) // Auto-detects compression
	defer fi.Close()

	reader := csv.NewReader(fi)
	reader.Comma = opts.Comma
	reader.Comment = opts.Comment
	reader.TrimLeadingSpace = opts.TrimLeadingSpace

	rowNum := 0

	// Skip header if requested
	if opts.SkipHeader {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				fmt.Printf("File %s. Rows processed: 0 (empty file)\n", filename)
				return nil
			}
			return fmt.Errorf("failed to read header: %w", err)
		}
		rowNum++ // Count header
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("row %d: CSV parse error: %w", rowNum+1, err)
		}

		rowNum++

		if err := processor(record); err != nil {
			return fmt.Errorf("row %d: processor error: %w", rowNum, err)
		}
	}

	fmt.Printf("File %s. Rows processed: %d\n", filename, rowNum)
	return nil
}

// IterateCSVMap processes a CSV file with header row, returning each row as a map.
// Supports compression auto-detection by extension (.gz, .zst).
// First row is used as headers (keys), subsequent rows are returned as map[header]value.
//
// filename - "filename" or "http://url" (with optional .gz or .zst extension)
// opts - CSV parsing options (SkipHeader is ignored, first row is always header)
// processor - function that receives each CSV row as a map
//
// Example:
//
//	fileiterator.IterateCSVMap("users.csv", fileiterator.DefaultCSVOptions(), func(row map[string]string) error {
//	    fmt.Printf("Name: %s, Email: %s\n", row["name"], row["email"])
//	    return nil
//	})
func IterateCSVMap(filename string, opts CSVOptions, processor func(map[string]string) error) error {
	fi := hb.FUOpen(filename) // Auto-detects compression
	defer fi.Close()

	reader := csv.NewReader(fi)
	reader.Comma = opts.Comma
	reader.Comment = opts.Comment
	reader.TrimLeadingSpace = opts.TrimLeadingSpace

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("empty CSV file (no header)")
		}
		return fmt.Errorf("failed to read header: %w", err)
	}

	rowNum := 1 // Header is row 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("row %d: CSV parse error: %w", rowNum+1, err)
		}

		rowNum++

		// Create map from headers and values
		rowMap := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				rowMap[header] = record[i]
			} else {
				rowMap[header] = ""
			}
		}

		if err := processor(rowMap); err != nil {
			return fmt.Errorf("row %d: processor error: %w", rowNum, err)
		}
	}

	fmt.Printf("File %s. Rows processed: %d (excluding header)\n", filename, rowNum-1)
	return nil
}
