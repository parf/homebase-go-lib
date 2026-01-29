package fileiterator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestFUCreateZstd1(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.zst1")
	testData := []byte("Test Zstd Level 1 Compression")

	// Write with Zstd-1
	w := fileiterator.FUCreate(testFile)
	_, err := w.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	w.Close()

	// Read back with FUOpen
	r := fileiterator.FUOpen(testFile)
	defer r.Close()

	readData := make([]byte, len(testData))
	n, err := r.Read(readData)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if n != len(testData) || string(readData) != string(testData) {
		t.Errorf("Data mismatch: got %s, want %s", string(readData), string(testData))
	}
}

func TestFUCreateZstd2(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.zst2")
	testData := []byte("Test Zstd Level 2 Compression")

	// Write with Zstd-2
	w := fileiterator.FUCreate(testFile)
	_, err := w.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	w.Close()

	// Read back with FUOpen
	r := fileiterator.FUOpen(testFile)
	defer r.Close()

	readData := make([]byte, len(testData))
	n, err := r.Read(readData)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if n != len(testData) || string(readData) != string(testData) {
		t.Errorf("Data mismatch: got %s, want %s", string(readData), string(testData))
	}
}

func TestZstdLevelComparison(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test data - repetitive for better compression
	testData := make([]byte, 100000)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	formats := []string{".zst1", ".zst2", ".zst"}
	sizes := make(map[string]int64)

	// Write with each format
	for _, ext := range formats {
		file := filepath.Join(tmpDir, "test"+ext)
		w := fileiterator.FUCreate(file)
		w.Write(testData)
		w.Close()

		stat, _ := os.Stat(file)
		sizes[ext] = stat.Size()
	}

	t.Logf("Compression sizes for 100KB data:")
	t.Logf("  Zstd-1: %d bytes", sizes[".zst1"])
	t.Logf("  Zstd-2: %d bytes", sizes[".zst2"])
	t.Logf("  Zstd-3: %d bytes", sizes[".zst"])

	// Verify all produce compressed output
	for ext, size := range sizes {
		if size >= int64(len(testData)) {
			t.Errorf("%s produced larger file than input: %d >= %d", ext, size, len(testData))
		}
	}
}
