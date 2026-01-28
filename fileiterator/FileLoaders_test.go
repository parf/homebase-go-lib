package fileiterator_test

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/parf/homebase-go-lib/fileiterator"
	"github.com/ulikunitz/xz"
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

func TestFUOpenZlibFile(t *testing.T) {
	// Create temp zlib file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.zlib")
	testData := []byte("Hello, Zlib!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	zw := zlib.NewWriter(f)
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

func TestFUOpenZzFile(t *testing.T) {
	// Create temp .zz file (zlib with .zz extension)
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.zz")
	testData := []byte("Hello, Zlib (.zz)!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	zw := zlib.NewWriter(f)
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

func TestFUOpenLz4File(t *testing.T) {
	// Create temp lz4 file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.lz4")
	testData := []byte("Hello, LZ4!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	lzw := lz4.NewWriter(f)
	lzw.Write(testData)
	lzw.Close()
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

func TestFUOpenBrotliFile(t *testing.T) {
	// Create temp brotli file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.br")
	testData := []byte("Hello, Brotli!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	brw := brotli.NewWriter(f)
	brw.Write(testData)
	brw.Close()
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

func TestFUOpenXzFile(t *testing.T) {
	// Create temp xz file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.xz")
	testData := []byte("Hello, XZ!")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	xzw, err := xz.NewWriter(f)
	if err != nil {
		t.Fatalf("Failed to create xz writer: %v", err)
	}
	xzw.Write(testData)
	xzw.Close()
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

	// Test zlib file
	zlibFile := filepath.Join(tmpDir, "test.txt.zlib")
	f3, _ := os.Create(zlibFile)
	zlibw := zlib.NewWriter(f3)
	zlibw.Write(testData)
	zlibw.Close()
	f3.Close()
	var result4 []byte
	fileiterator.LoadBinFile(zlibFile, &result4)
	if !bytes.Equal(result4, testData) {
		t.Errorf("Zlib file: Expected %s, got %s", testData, result4)
	}

	// Test lz4 file
	lz4File := filepath.Join(tmpDir, "test.txt.lz4")
	f4, _ := os.Create(lz4File)
	lz4w := lz4.NewWriter(f4)
	lz4w.Write(testData)
	lz4w.Close()
	f4.Close()
	var result5 []byte
	fileiterator.LoadBinFile(lz4File, &result5)
	if !bytes.Equal(result5, testData) {
		t.Errorf("LZ4 file: Expected %s, got %s", testData, result5)
	}

	// Test brotli file
	brFile := filepath.Join(tmpDir, "test.txt.br")
	f5, _ := os.Create(brFile)
	brw := brotli.NewWriter(f5)
	brw.Write(testData)
	brw.Close()
	f5.Close()
	var result6 []byte
	fileiterator.LoadBinFile(brFile, &result6)
	if !bytes.Equal(result6, testData) {
		t.Errorf("Brotli file: Expected %s, got %s", testData, result6)
	}

	// Test xz file
	xzFile := filepath.Join(tmpDir, "test.txt.xz")
	f6, _ := os.Create(xzFile)
	xzw, _ := xz.NewWriter(f6)
	xzw.Write(testData)
	xzw.Close()
	f6.Close()
	var result7 []byte
	fileiterator.LoadBinFile(xzFile, &result7)
	if !bytes.Equal(result7, testData) {
		t.Errorf("XZ file: Expected %s, got %s", testData, result7)
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

	// Test zlib file
	zlibFile := filepath.Join(tmpDir, "test.txt.zlib")
	f4, _ := os.Create(zlibFile)
	zlibw := zlib.NewWriter(f4)
	zlibw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	zlibw.Close()
	f4.Close()
	var result4 []string
	fileiterator.IterateLines(zlibFile, func(line string) {
		result4 = append(result4, line)
	})
	if len(result4) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result4))
	}

	// Test lz4 file
	lz4File := filepath.Join(tmpDir, "test.txt.lz4")
	f5, _ := os.Create(lz4File)
	lz4w := lz4.NewWriter(f5)
	lz4w.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	lz4w.Close()
	f5.Close()
	var result5 []string
	fileiterator.IterateLines(lz4File, func(line string) {
		result5 = append(result5, line)
	})
	if len(result5) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result5))
	}

	// Test brotli file
	brFile := filepath.Join(tmpDir, "test.txt.br")
	f6, _ := os.Create(brFile)
	brw := brotli.NewWriter(f6)
	brw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	brw.Close()
	f6.Close()
	var result6 []string
	fileiterator.IterateLines(brFile, func(line string) {
		result6 = append(result6, line)
	})
	if len(result6) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result6))
	}

	// Test xz file
	xzFile := filepath.Join(tmpDir, "test.txt.xz")
	f7, _ := os.Create(xzFile)
	xzw, _ := xz.NewWriter(f7)
	io.WriteString(xzw, "Line 1\nLine 2\nLine 3\n")
	xzw.Close()
	f7.Close()
	var result7 []string
	fileiterator.IterateLines(xzFile, func(line string) {
		result7 = append(result7, line)
	})
	if len(result7) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(result7))
	}
}

func TestIterateIDTabFile(t *testing.T) {
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
	fileiterator.IterateIDTabFile(gzFile, func(id int32, name string) {
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
