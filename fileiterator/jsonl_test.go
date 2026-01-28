package fileiterator_test

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestIterateJSONLPlain(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl")

	// Create JSONL file
	data := `{"id":1,"name":"Alice"}
{"id":2,"name":"Bob"}
{"id":3,"name":"Charlie"}
`
	os.WriteFile(testFile, []byte(data), 0644)

	// Iterate and collect results
	var names []string
	err := fileiterator.IterateJSONL(testFile, func(obj map[string]any) error {
		names = append(names, obj["name"].(string))
		return nil
	})

	if err != nil {
		t.Fatalf("IterateJSONL failed: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}
}

func TestIterateJSONLGzip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl.gz")

	// Create gzipped JSONL file
	data := `{"id":1,"value":"test1"}
{"id":2,"value":"test2"}
`
	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte(data))
	gz.Close()
	f.Close()

	// Iterate and count
	count := 0
	err := fileiterator.IterateJSONL(testFile, func(obj map[string]any) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateJSONL failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 records, got %d", count)
	}
}

func TestIterateJSONLZstd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl.zst")

	// Create zstd JSONL file
	data := `{"x":1}
{"x":2}
{"x":3}
`
	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write([]byte(data))
	zw.Close()
	f.Close()

	// Iterate and sum
	sum := 0
	err := fileiterator.IterateJSONL(testFile, func(obj map[string]any) error {
		sum += int(obj["x"].(float64))
		return nil
	})

	if err != nil {
		t.Fatalf("IterateJSONL failed: %v", err)
	}

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func TestIterateJSONLTyped(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "users.jsonl")

	// Create JSONL file
	data := `{"id":1,"name":"Alice"}
{"id":2,"name":"Bob"}
`
	os.WriteFile(testFile, []byte(data), 0644)

	// Iterate with typed struct
	var users []User
	err := fileiterator.IterateJSONLTyped(testFile, func(user User) error {
		users = append(users, user)
		return nil
	})

	if err != nil {
		t.Fatalf("IterateJSONLTyped failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	if users[0].Name != "Alice" {
		t.Errorf("Expected first user name 'Alice', got '%s'", users[0].Name)
	}
}
