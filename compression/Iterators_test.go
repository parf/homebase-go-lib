package compression_test

import (
	"compress/gzip"
	"compress/zlib"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/compression"
)

func TestIterateBinaryRecordsPlain(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin")

	// Create plain binary file with 3 records of 10 bytes each
	data := []byte("0123456789ABCDEFGHIJKLMNOPQRST")
	err := os.WriteFile(testFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Iterate over records
	count := 0
	compression.IterateBinaryRecords(testFile, 10, func(record []byte) {
		count++
		if len(record) != 10 {
			t.Errorf("Expected record size 10, got %d", len(record))
		}
	})

	if count != 3 {
		t.Errorf("Expected 3 records, got %d", count)
	}
}

func TestIterateBinaryRecordsGzip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.gz")

	// Create gzipped binary file with 2 records of 5 bytes each
	data := []byte("1234567890")
	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write(data)
	gz.Close()
	f.Close()

	// Iterate over records
	count := 0
	compression.IterateBinaryRecords(testFile, 5, func(record []byte) {
		count++
		if len(record) != 5 {
			t.Errorf("Expected record size 5, got %d", len(record))
		}
	})

	if count != 2 {
		t.Errorf("Expected 2 records, got %d", count)
	}
}

func TestIterateBinaryRecordsZstd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.zst")

	// Create zstd binary file with 4 records of 8 bytes each
	data := []byte("AAAABBBBCCCCDDDD")
	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write(data)
	zw.Close()
	f.Close()

	// Iterate over records
	count := 0
	compression.IterateBinaryRecords(testFile, 8, func(record []byte) {
		count++
		if len(record) != 8 {
			t.Errorf("Expected record size 8, got %d", len(record))
		}
	})

	if count != 2 {
		t.Errorf("Expected 2 records, got %d", count)
	}
}

func TestIterateBinaryRecordsZlib(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.zlib")

	// Create zlib binary file with 3 records of 6 bytes each
	data := []byte("111111222222333333")
	f, _ := os.Create(testFile)
	zw := zlib.NewWriter(f)
	zw.Write(data)
	zw.Close()
	f.Close()

	// Iterate over records
	count := 0
	compression.IterateBinaryRecords(testFile, 6, func(record []byte) {
		count++
		if len(record) != 6 {
			t.Errorf("Expected record size 6, got %d", len(record))
		}
	})

	if count != 3 {
		t.Errorf("Expected 3 records, got %d", count)
	}
}
