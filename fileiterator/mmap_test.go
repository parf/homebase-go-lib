package fileiterator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestMmapOpen(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin")

	// Create test file
	testData := []byte("Hello Memory-Mapped File!")
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Open with mmap
	mmapFile, err := fileiterator.MmapOpen(testFile)
	if err != nil {
		t.Fatalf("Failed to mmap file: %v", err)
	}
	defer mmapFile.Close()

	// Verify data
	if len(mmapFile.Data) != len(testData) {
		t.Errorf("Data length mismatch: got %d, want %d", len(mmapFile.Data), len(testData))
	}

	if string(mmapFile.Data) != string(testData) {
		t.Errorf("Data mismatch: got %s, want %s", string(mmapFile.Data), string(testData))
	}
}

func TestLoadMmap(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin")

	// Create test file with larger data
	testData := make([]byte, 1024*1024) // 1MB
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load with mmap
	data, err := fileiterator.LoadMmap(testFile)
	if err != nil {
		t.Fatalf("Failed to load mmap: %v", err)
	}

	// Verify size
	if len(data) != len(testData) {
		t.Errorf("Data length mismatch: got %d, want %d", len(data), len(testData))
	}

	// Verify content (spot check)
	for i := 0; i < 100; i++ {
		idx := i * 1000
		if data[idx] != testData[idx] {
			t.Errorf("Data mismatch at index %d: got %d, want %d", idx, data[idx], testData[idx])
		}
	}
}

func TestMmapOpenNonExistent(t *testing.T) {
	_, err := fileiterator.MmapOpen("/nonexistent/file.bin")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestMmapOpenEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.bin")

	// Create empty file
	err := os.WriteFile(testFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Mmap empty file - should fail (can't mmap 0 bytes)
	_, err = fileiterator.MmapOpen(testFile)
	if err == nil {
		t.Error("Expected error when mmapping empty file")
	}
}

func TestMmapRandomAccess(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "random.bin")

	// Create file with known pattern
	size := 10 * 1024 * 1024 // 10MB
	testData := make([]byte, size)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load with mmap for random access
	data, err := fileiterator.LoadMmap(testFile)
	if err != nil {
		t.Fatalf("Failed to load mmap: %v", err)
	}

	// Random access test - much faster with mmap than reading file
	testIndices := []int{0, 1000, 50000, 100000, 500000, 1000000, 5000000, 9999999}
	for _, idx := range testIndices {
		expected := byte(idx % 256)
		if data[idx] != expected {
			t.Errorf("Random access failed at index %d: got %d, want %d", idx, data[idx], expected)
		}
	}
}
