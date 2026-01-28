package fileiterator

import (
	"bufio"
	"encoding/json"
	"fmt"

	hb "github.com/parf/homebase-go-lib"
)

// IterateJSONL processes a JSONL (JSON Lines) file line by line.
// Supports compression auto-detection by extension (.gz, .zst).
// Each line is parsed as JSON and passed to the processor function.
//
// filename - "filename" or "http://url" (with optional .gz or .zst extension)
// processor - function that receives each parsed JSON object
//
// Example:
//
//	fileiterator.IterateJSONL("data.jsonl.gz", func(obj map[string]any) error {
//	    fmt.Printf("ID: %v, Name: %v\n", obj["id"], obj["name"])
//	    return nil
//	})
func IterateJSONL(filename string, processor func(map[string]any) error) error {
	fi := hb.FUOpen(filename) // Auto-detects compression
	defer fi.Close()

	scanner := bufio.NewScanner(fi)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			return fmt.Errorf("line %d: JSON parse error: %w", lineNum, err)
		}

		if err := processor(obj); err != nil {
			return fmt.Errorf("line %d: processor error: %w", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fmt.Printf("File %s. Lines processed: %d\n", filename, lineNum)
	return nil
}

// IterateJSONLTyped processes a JSONL file with a typed struct.
// Supports compression auto-detection by extension (.gz, .zst).
// Each line is parsed into the provided type T.
//
// filename - "filename" or "http://url" (with optional .gz or .zst extension)
// processor - function that receives each parsed object of type T
//
// Example:
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//	fileiterator.IterateJSONLTyped("users.jsonl", func(user User) error {
//	    fmt.Printf("User: %s\n", user.Name)
//	    return nil
//	})
func IterateJSONLTyped[T any](filename string, processor func(T) error) error {
	fi := hb.FUOpen(filename) // Auto-detects compression
	defer fi.Close()

	scanner := bufio.NewScanner(fi)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		var obj T
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			return fmt.Errorf("line %d: JSON parse error: %w", lineNum, err)
		}

		if err := processor(obj); err != nil {
			return fmt.Errorf("line %d: processor error: %w", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fmt.Printf("File %s. Lines processed: %d\n", filename, lineNum)
	return nil
}
