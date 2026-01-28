package fileiterator_test

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestIterateCSVPlain(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.csv")

	// Create CSV file
	data := `name,age,city
Alice,30,NYC
Bob,25,LA
`
	os.WriteFile(testFile, []byte(data), 0644)

	// Iterate with header skip
	opts := fileiterator.DefaultCSVOptions()
	opts.SkipHeader = true

	var names []string
	err := fileiterator.IterateCSV(testFile, opts, func(row []string) error {
		names = append(names, row[0])
		return nil
	})

	if err != nil {
		t.Fatalf("IterateCSV failed: %v", err)
	}

	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}
}

func TestIterateCSVGzip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.csv.gz")

	// Create gzipped CSV file
	data := `a,b,c
1,2,3
4,5,6
`
	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write([]byte(data))
	gz.Close()
	f.Close()

	// Iterate without skipping header
	count := 0
	err := fileiterator.IterateCSV(testFile, fileiterator.DefaultCSVOptions(), func(row []string) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateCSV failed: %v", err)
	}

	if count != 3 { // Header + 2 data rows
		t.Errorf("Expected 3 rows, got %d", count)
	}
}

func TestIterateCSVZstd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.csv.zst")

	// Create zstd CSV file
	data := `x,y
10,20
30,40
`
	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	zw.Write([]byte(data))
	zw.Close()
	f.Close()

	// Iterate with header skip
	opts := fileiterator.DefaultCSVOptions()
	opts.SkipHeader = true

	count := 0
	err := fileiterator.IterateCSV(testFile, opts, func(row []string) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateCSV failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 data rows, got %d", count)
	}
}

func TestIterateCSVMap(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.csv")

	// Create CSV file with header
	data := `name,age,city
Alice,30,NYC
Bob,25,LA
`
	os.WriteFile(testFile, []byte(data), 0644)

	// Iterate as map
	var results []map[string]string
	err := fileiterator.IterateCSVMap(testFile, fileiterator.DefaultCSVOptions(), func(row map[string]string) error {
		results = append(results, row)
		return nil
	})

	if err != nil {
		t.Fatalf("IterateCSVMap failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(results))
	}

	if results[0]["name"] != "Alice" {
		t.Errorf("Expected first name 'Alice', got '%s'", results[0]["name"])
	}

	if results[1]["age"] != "25" {
		t.Errorf("Expected second age '25', got '%s'", results[1]["age"])
	}
}

func TestIterateCSVCustomDelimiter(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsv")

	// Create TSV file (tab-separated)
	data := "a\tb\tc\n1\t2\t3\n"
	os.WriteFile(testFile, []byte(data), 0644)

	// Iterate with tab delimiter
	opts := fileiterator.DefaultCSVOptions()
	opts.Comma = '\t'

	count := 0
	err := fileiterator.IterateCSV(testFile, opts, func(row []string) error {
		if len(row) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(row))
		}
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateCSV failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}
