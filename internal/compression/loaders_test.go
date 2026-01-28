package compression_test

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/internal/compression"
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

func TestLoadLinesGzFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.gz")

	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	gz.Close()
	f.Close()

	var lines []string
	compression.LoadLinesGzFile(testFile, func(line string) {
		lines = append(lines, line)
	})

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestLoadLinesZstdFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt.zst")

	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write([]byte("Line 1\nLine 2\nLine 3\n"))
	zw.Close()
	f.Close()

	var lines []string
	compression.LoadLinesZstdFile(testFile, func(line string) {
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
