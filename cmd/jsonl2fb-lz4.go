package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	// Parse command line flags
	inputFile := flag.String("input", "", "Input JSONL file (required)")
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

	// Read JSONL file and convert to FlatBuffer records
	records := make([][]byte, 0)
	err := fileiterator.IterateJSONL(*inputFile, func(obj map[string]any) error {
		// Marshal object back to JSON (to store in FlatBuffer)
		jsonData, err := json.Marshal(obj)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		// Create FlatBuffer with JSON data as byte vector
		builder := flatbuffers.NewBuilder(1024)
		dataOffset := builder.CreateByteVector(jsonData)
		builder.Finish(dataOffset)

		records = append(records, builder.FinishedBytes())
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSONL: %v\n", err)
		os.Exit(1)
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
		fmt.Printf("  Input:  %.2f MB\n", float64(inStat.Size())/1024/1024)
		fmt.Printf("  Output: %.2f MB (%.1f%% of original)\n",
			float64(outStat.Size())/1024/1024, ratio)
	}
}
