package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/apache/arrow/go/v14/parquet/file"
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

	// Read parquet records and convert to FlatBuffer
	records := make([][]byte, 0)

	// Get reader for all row groups
	reader := file.NewParquetReader(pf)
	numRows := reader.NumRows()

	fmt.Printf("Reading %d rows from parquet file...\n", numRows)

	// Read all rows
	for i := int64(0); i < numRows; i++ {
		// Read row as map
		row, err := reader.ReadRow()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading row %d: %v\n", i, err)
			os.Exit(1)
		}

		// Convert row to JSON
		jsonData, err := json.Marshal(row)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling row %d: %v\n", i, err)
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
