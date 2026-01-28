package fileiterator_test

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestFUOpenPlainFile(t *testing.T) {
	// Create temp plain file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testData := []byte("Hello, World!")

	err := ioutil.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test FUOpen
	r := fileiterator.FUOpen(testFile)
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(data, testData) {
		t.Errorf("Expected %s, got %s", testData, data)
	}
}

func TestFUOpenGzipFile(t *testing.T) {
	// Create temp gzipped file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.gz")
	testData := []byte("Hello, Gzip!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	gz := gzip.NewWriter(f)
	gz.Write(testData)
	gz.Close()
	f.Close()

	// Test FUOpen with automatic decompression
	r := fileiterator.FUOpen(testFile)
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(data, testData) {
		t.Errorf("Expected %s, got %s", testData, data)
	}
}

func TestFUOpenZstdFile(t *testing.T) {
	// Create temp zstd file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.zst")
	testData := []byte("Hello, Zstd!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	zw, _ := zstd.NewWriter(f)
	zw.Write(testData)
	zw.Close()
	f.Close()

	// Test FUOpen with automatic decompression
	r := fileiterator.FUOpen(testFile)
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(data, testData) {
		t.Errorf("Expected %s, got %s", testData, data)
	}
}

func TestLoadBinFile(t *testing.T) {
	tmpDir := t.TempDir()
	testData := []byte("Test data for LoadBinFile")

	// Test plain file
	plainFile := filepath.Join(tmpDir, "plain.txt")
	ioutil.WriteFile(plainFile, testData, 0644)
	var result1 []byte
	fileiterator.LoadBinFile(plainFile, &result1)
	if !bytes.Equal(result1, testData) {
		t.Errorf("Plain file: Expected %s, got %s", testData, result1)
	}

	// Test gzipped file
	gzFile := filepath.Join(tmpDir, "test.txt.gz")
	f, _ := os.Create(gzFile)
	gz := gzip.NewWriter(f)
	gz.Write(testData)
	gz.Close()
	f.Close()
	var result2 []byte
	fileiterator.LoadBinFile(gzFile, &result2)
	if !bytes.Equal(result2, testData) {
		t.Errorf("Gzip file: Expected %s, got %s", testData, result2)
	}

	// Test zstd file
	zstFile := filepath.Join(tmpDir, "test.txt.zst")
	f2, _ := os.Create(zstFile)
	zw, _ := zstd.NewWriter(f2)
	zw.Write(testData)
	zw.Close()
	f2.Close()
	var result3 []byte
	fileiterator.LoadBinFile(zstFile, &result3)
	if !bytes.Equal(result3, testData) {
		t.Errorf("Zstd file: Expected %s, got %s", testData, result3)
	}
}

func TestIterateLines(t *testing.T) {
	tmpDir := t.TempDir()
	testLines := []string{"Line 1", "Line 2", "Line 3"}

	// Test plain file
	plainFile := filepath.Join(tmpDir, "plain.txt")
	ioutil.WriteFile(plainFile, []byte("Line 1\nLine 2\nLine 3\n"), 0644)
	var result1 []string
	fileiterator.IterateLines(plainFile, func(line string) {
		result1 = append(result1, line)
	})
	if len(result1) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result1))
	}

	// Test gzipped file
	gzFile := filepath.Join(tmpDir, "test.txt.gz")
	f, _ := os.Create(gzFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	gz.Close()
	f.Close()
	var result2 []string
	fileiterator.IterateLines(gzFile, func(line string) {
		result2 = append(result2, line)
	})
	if len(result2) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result2))
	}

	// Test zstd file
	zstFile := filepath.Join(tmpDir, "test.txt.zst")
	f2, _ := os.Create(zstFile)
	zw, _ := zstd.NewWriter(f2)
	zw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	zw.Close()
	f2.Close()
	var result3 []string
	fileiterator.IterateLines(zstFile, func(line string) {
		result3 = append(result3, line)
	})
	if len(result3) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result3))
	}
}

func TestLoadBinGzFile(t *testing.T) {
	tmpDir := t.TempDir()
	testData := []byte("Test data for LoadBinGzFile")

	// Create gzipped file
	gzFile := filepath.Join(tmpDir, "test.gz")
	f, _ := os.Create(gzFile)
	gz := gzip.NewWriter(f)
	gz.Write(testData)
	gz.Close()
	f.Close()

	var result []byte
	fileiterator.LoadBinGzFile(gzFile, &result)
	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestLoadBinZstdFile(t *testing.T) {
	tmpDir := t.TempDir()
	testData := []byte("Test data for LoadBinZstdFile")

	// Create zstd file
	zstFile := filepath.Join(tmpDir, "test.zst")
	f, _ := os.Create(zstFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write(testData)
	zw.Close()
	f.Close()

	var result []byte
	fileiterator.LoadBinZstdFile(zstFile, &result)
	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestIterateLinesGz(t *testing.T) {
	tmpDir := t.TempDir()
	testLines := []string{"Line 1", "Line 2", "Line 3"}

	// Create gzipped file
	gzFile := filepath.Join(tmpDir, "test.gz")
	f, _ := os.Create(gzFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	gz.Close()
	f.Close()

	var result []string
	fileiterator.IterateLinesGz(gzFile, func(line string) {
		result = append(result, line)
	})
	if len(result) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result))
	}
}

func TestIterateLinesZstd(t *testing.T) {
	tmpDir := t.TempDir()
	testLines := []string{"Line 1", "Line 2", "Line 3"}

	// Create zstd file
	zstFile := filepath.Join(tmpDir, "test.zst")
	f, _ := os.Create(zstFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	zw.Close()
	f.Close()

	var result []string
	fileiterator.IterateLinesZstd(zstFile, func(line string) {
		result = append(result, line)
	})
	if len(result) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result))
	}
}

func TestLoadIDTabGzFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create gzipped tab-separated file with hex IDs
	gzFile := filepath.Join(tmpDir, "test.gz")
	f, _ := os.Create(gzFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("A\tFIRST\nB\tSECOND\n14\tTWENTY\n"))
	gz.Close()
	f.Close()

	var ids []int32
	var names []string
	fileiterator.LoadIDTabGzFile(gzFile, func(id int32, name string) {
		ids = append(ids, id)
		names = append(names, name)
	})

	expectedIDs := []int32{10, 11, 20}
	expectedNames := []string{"first", "second", "twenty"}

	if len(ids) != len(expectedIDs) {
		t.Fatalf("Expected %d entries, got %d", len(expectedIDs), len(ids))
	}

	for i := range ids {
		if ids[i] != expectedIDs[i] {
			t.Errorf("ID[%d]: Expected %d, got %d", i, expectedIDs[i], ids[i])
		}
		if names[i] != expectedNames[i] {
			t.Errorf("Name[%d]: Expected %s, got %s", i, expectedNames[i], names[i])
		}
	}
}
