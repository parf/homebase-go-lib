package main

import (
	"fmt"
	"log"
	"os"

	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	fmt.Println("Compression Support Example")
	fmt.Println("===========================")
	fmt.Println()

	// Example 1: FUOpen with auto-decompression
	fmt.Println("Example 1: FUOpen with automatic decompression")
	fmt.Println("-----------------------------------------------")
	fmt.Println("FUOpen automatically detects compression by file extension:")
	fmt.Println("  - .gz files are decompressed with gzip")
	fmt.Println("  - .zst files are decompressed with zstd")
	fmt.Println("  - Other files are opened as-is")
	fmt.Println()

	// Create test file (uncompressed)
	testFile := "/tmp/test.txt"
	err := os.WriteFile(testFile, []byte("Hello, World!\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(testFile)

	// Open plain file
	r := fileiterator.FUOpen(testFile)
	fmt.Printf("Opened: %s\n", testFile)
	r.Close()
	fmt.Println()

	// Example 2: LoadBinFile with auto-decompression
	fmt.Println("Example 2: LoadBinFile with automatic decompression")
	fmt.Println("----------------------------------------------------")
	var data []byte
	fileiterator.LoadBinFile(testFile, &data)
	fmt.Printf("Loaded %d bytes: %s\n", len(data), string(data))
	fmt.Println()

	// Example 3: IterateLines with auto-decompression
	fmt.Println("Example 3: IterateLines with automatic decompression")
	fmt.Println("-----------------------------------------------------")
	lineCount := 0
	fileiterator.IterateLines(testFile, func(line string) {
		lineCount++
		fmt.Printf("Line %d: %s\n", lineCount, line)
	})
	fmt.Println()

	fmt.Println("Supported compression formats:")
	fmt.Println("  - Gzip (.gz)")
	fmt.Println("  - Zstd (.zst)")
	fmt.Println()
	fmt.Println("All loader functions now support automatic decompression!")
}
