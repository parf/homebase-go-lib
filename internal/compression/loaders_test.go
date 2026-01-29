package compression_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/parf/homebase-go-lib/internal/compression"
	"github.com/ulikunitz/xz"
)

func TestLoadBinGzFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.gz")
	testData := []byte("Test gzip data")

	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write(testData)
	gz.Close()
	f.Close()

	var result []byte
	compression.LoadBinGzFile(testFile, &result)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestLoadBinZstdFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.zst")
	testData := []byte("Test zstd data")

	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write(testData)
	zw.Close()
	f.Close()

	var result []byte
	compression.LoadBinZstdFile(testFile, &result)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestIterateLinesGz(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.gz")

	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	gz.Close()
	f.Close()

	var lines []string
	compression.IterateLinesGz(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestIterateLinesZstd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.zst")

	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	zw.Close()
	f.Close()

	var lines []string
	compression.IterateLinesZstd(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestLoadBinLz4File(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.lz4")
	testData := []byte("Test LZ4 data")

	f, _ := os.Create(testFile)
	lzw := lz4.NewWriter(f)
	lzw.Write(testData)
	lzw.Close()
	f.Close()

	var result []byte
	compression.LoadBinLz4File(testFile, &result)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestIterateLinesLz4(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.lz4")

	f, _ := os.Create(testFile)
	lzw := lz4.NewWriter(f)
	lzw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	lzw.Close()
	f.Close()

	var lines []string
	compression.IterateLinesLz4(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestLoadBinBrotliFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.br")
	testData := []byte("Test Brotli data")

	f, _ := os.Create(testFile)
	brw := brotli.NewWriter(f)
	brw.Write(testData)
	brw.Close()
	f.Close()

	var result []byte
	compression.LoadBinBrotliFile(testFile, &result)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestIterateLinesBrotli(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.br")

	f, _ := os.Create(testFile)
	brw := brotli.NewWriter(f)
	brw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	brw.Close()
	f.Close()

	var lines []string
	compression.IterateLinesBrotli(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestLoadBinXzFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin.xz")
	testData := []byte("Test XZ data")

	f, _ := os.Create(testFile)
	xzw, _ := xz.NewWriter(f)
	xzw.Write(testData)
	xzw.Close()
	f.Close()

	var result []byte
	compression.LoadBinXzFile(testFile, &result)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}

func TestIterateLinesXz(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.xz")

	f, _ := os.Create(testFile)
	xzw, _ := xz.NewWriter(f)
	io.WriteString(xzw, "Line 1\nLine 2\nLine 3\n")
	xzw.Close()
	f.Close()

	var lines []string
	compression.IterateLinesXz(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestLoadIDTabGzFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tab.gz")

	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("A\tAlpha\nB\tBeta\nC\tGamma\n"))
	gz.Close()
	f.Close()

	count := 0
	compression.LoadIDTabGzFile(testFile, func(id int32, name string) {
		count++
	})

	if count != 3 {
		t.Errorf("Expected 3 entries, got %d", count)
	}
}
