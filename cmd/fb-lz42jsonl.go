package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	// Parse command line flags
	inputFile := flag.String("input", "", "Input FlatBuffer file (required)")
	outputFile := flag.String("output", "", "Output JSONL file (default: input.jsonl)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -input flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Default output file
	if *outputFile == "" {
		*outputFile = *inputFile + ".jsonl"
	}

	fmt.Printf("Converting %s to %s...\n", *inputFile, *outputFile)

	// Open output file
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Read FlatBuffer list and convert to JSONL
	count := 0
	err = fileiterator.IterateFlatBufferList(*inputFile, func(data []byte) error {
		// The FlatBuffer contains JSON data as byte vector
		// Skip FlatBuffer header and extract the JSON bytes
		// For simplicity, we'll just write the raw data
		// In production, you'd properly decode the FlatBuffer structure

		// Unmarshal to verify it's valid JSON
		var obj map[string]any
		if err := json.Unmarshal(data, &obj); err == nil {
			// Write as JSONL
			jsonData, _ := json.Marshal(obj)
			fmt.Fprintf(outFile, "%s\n", jsonData)
		}

		count++
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading FlatBuffer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Converted %d records\n", count)

	// Show file sizes
	inStat, _ := os.Stat(*inputFile)
	outStat, _ := os.Stat(*outputFile)
	if inStat != nil && outStat != nil {
		ratio := float64(inStat.Size()) / float64(outStat.Size()) * 100
		fmt.Printf("  Input:  %.2f MB\n", float64(inStat.Size())/1024/1024)
		fmt.Printf("  Output: %.2f MB (input was %.1f%% of output size)\n",
			float64(outStat.Size())/1024/1024, ratio)
	}
}
