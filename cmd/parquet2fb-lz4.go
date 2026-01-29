package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	// Parse command line flags
	inputFile := flag.String("input", "", "Input Parquet file (required)")
	outputFile := flag.String("output", "", "Output FlatBuffer file (default: input.fb.lz4)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -input flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Default output file
	if *outputFile == "" {
		*outputFile = *inputFile + ".fb.lz4"
	}

	fmt.Printf("Converting %s to %s...\n", *inputFile, *outputFile)

	// Open parquet file
	pf, err := file.OpenParquetFile(*inputFile, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening parquet file: %v\n", err)
		os.Exit(1)
	}
	defer pf.Close()

	// Create Arrow reader
	reader, err := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Arrow reader: %v\n", err)
		os.Exit(1)
	}

	// Read table
	tbl, err := reader.ReadTable(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading table: %v\n", err)
		os.Exit(1)
	}
	defer tbl.Release()

	numRows := tbl.NumRows()
	fmt.Printf("Reading %d rows from parquet file...\n", numRows)

	// Convert to FlatBuffer records
	records := make([][]byte, 0, numRows)

	// Get column names
	schema := tbl.Schema()
	fieldNames := make([]string, schema.NumFields())
	for i := 0; i < schema.NumFields(); i++ {
		fieldNames[i] = schema.Field(i).Name
	}

	// Process each row
	for rowIdx := int64(0); rowIdx < numRows; rowIdx++ {
		row := make(map[string]interface{})

		// Extract values from each column
		for colIdx := 0; colIdx < int(tbl.NumCols()); colIdx++ {
			col := tbl.Column(colIdx)
			fieldName := fieldNames[colIdx]

			// Get value from the first chunk (assume single chunk for simplicity)
			if col.Len() == 0 {
				continue
			}

			chunk := col.Data().Chunk(0)

			// Handle different Arrow types
			switch arr := chunk.(type) {
			case *array.Int64:
				if !arr.IsNull(int(rowIdx)) {
					row[fieldName] = arr.Value(int(rowIdx))
				}
			case *array.Float64:
				if !arr.IsNull(int(rowIdx)) {
					row[fieldName] = arr.Value(int(rowIdx))
				}
			case *array.String:
				if !arr.IsNull(int(rowIdx)) {
					row[fieldName] = arr.Value(int(rowIdx))
				}
			case *array.Boolean:
				if !arr.IsNull(int(rowIdx)) {
					row[fieldName] = arr.Value(int(rowIdx))
				}
			}
		}

		// Convert row to JSON
		jsonData, err := json.Marshal(row)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling row %d: %v\n", rowIdx, err)
			os.Exit(1)
		}

		// Create FlatBuffer with JSON data
		builder := flatbuffers.NewBuilder(1024)
		dataOffset := builder.CreateByteVector(jsonData)
		builder.Finish(dataOffset)

		records = append(records, builder.FinishedBytes())
	}

	// Save as FlatBuffer list with LZ4 compression
	err = fileiterator.SaveFlatBufferList(*outputFile, records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing FlatBuffer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Converted %d records\n", len(records))

	// Show file sizes
	inStat, _ := os.Stat(*inputFile)
	outStat, _ := os.Stat(*outputFile)
	if inStat != nil && outStat != nil {
		ratio := float64(outStat.Size()) / float64(inStat.Size()) * 100
		fmt.Printf("  Input:  %.2f MB (Parquet)\n", float64(inStat.Size())/1024/1024)
		fmt.Printf("  Output: %.2f MB (FlatBuffer + LZ4, %.1f%% of parquet size)\n",
			float64(outStat.Size())/1024/1024, ratio)
	}
}
