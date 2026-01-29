package fileiterator_test

import (
	"bytes"
	"path/filepath"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestSaveLoadFlatBufferList(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fb")

	// Create multiple FlatBuffer records
	records := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		builder := flatbuffers.NewBuilder(1024)
		testData := []byte("FlatBuffer Record " + string(rune('A'+i)))
		dataOffset := builder.CreateByteVector(testData)
		builder.Finish(dataOffset)
		records[i] = builder.FinishedBytes()
	}

	// Save list
	err := fileiterator.SaveFlatBufferList(testFile, records)
	if err != nil {
		t.Fatalf("Failed to save FlatBuffer list: %v", err)
	}

	// Load list
	loadedRecords := make([][]byte, 0)
	err = fileiterator.IterateFlatBufferList(testFile, func(data []byte) error {
		loadedRecords = append(loadedRecords, data)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to iterate FlatBuffer list: %v", err)
	}

	// Verify
	if len(loadedRecords) != len(records) {
		t.Errorf("Record count mismatch: got %d, want %d", len(loadedRecords), len(records))
	}

	for i, loaded := range loadedRecords {
		if !bytes.Equal(loaded, records[i]) {
			t.Errorf("Record %d mismatch", i)
		}
	}
}

func TestSaveLoadFlatBufferListCompressed(t *testing.T) {
	tmpDir := t.TempDir()

	testFormats := []struct {
		name string
		ext  string
	}{
		{"Gzip", ".fb.gz"},
		{"Zstd", ".fb.zst"},
		{"Zstd-1", ".fb.zst1"},
		{"Zstd-2", ".fb.zst2"},
		{"Zlib", ".fb.zlib"},
		{"LZ4", ".fb.lz4"},
		{"Brotli", ".fb.br"},
		{"XZ", ".fb.xz"},
	}

	// Create test records
	records := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		builder := flatbuffers.NewBuilder(1024)
		testData := []byte("Compressed FlatBuffer List Record #" + string(rune('0'+i)))
		dataOffset := builder.CreateByteVector(testData)
		builder.Finish(dataOffset)
		records[i] = builder.FinishedBytes()
	}

	for _, tf := range testFormats {
		t.Run(tf.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test"+tf.ext)

			// Save compressed list
			err := fileiterator.SaveFlatBufferList(testFile, records)
			if err != nil {
				t.Fatalf("Failed to save compressed FlatBuffer list: %v", err)
			}

			// Load compressed list
			loadedRecords := make([][]byte, 0)
			err = fileiterator.IterateFlatBufferList(testFile, func(data []byte) error {
				loadedRecords = append(loadedRecords, data)
				return nil
			})
			if err != nil {
				t.Fatalf("Failed to iterate compressed FlatBuffer list: %v", err)
			}

			// Verify
			if len(loadedRecords) != len(records) {
				t.Errorf("Record count mismatch: got %d, want %d", len(loadedRecords), len(records))
			}

			for i, loaded := range loadedRecords {
				if !bytes.Equal(loaded, records[i]) {
					t.Errorf("Record %d mismatch", i)
				}
			}
		})
	}
}

func TestIterateFlatBufferListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.fb")

	// Save empty list
	err := fileiterator.SaveFlatBufferList(testFile, [][]byte{})
	if err != nil {
		t.Fatalf("Failed to save empty list: %v", err)
	}

	// Iterate empty list
	count := 0
	err = fileiterator.IterateFlatBufferList(testFile, func(data []byte) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to iterate empty list: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 records, got %d", count)
	}
}

func TestIterateFlatBufferListLZ4(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fb.lz4")

	// Create large dataset to test LZ4 compression
	records := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		builder := flatbuffers.NewBuilder(1024)
		// Create larger test data to make compression effective
		testData := make([]byte, 1000)
		for j := range testData {
			testData[j] = byte(i % 256)
		}
		dataOffset := builder.CreateByteVector(testData)
		builder.Finish(dataOffset)
		records[i] = builder.FinishedBytes()
	}

	// Save with LZ4 compression
	err := fileiterator.SaveFlatBufferList(testFile, records)
	if err != nil {
		t.Fatalf("Failed to save LZ4 FlatBuffer list: %v", err)
	}

	// Load with LZ4 decompression
	loadedCount := 0
	err = fileiterator.IterateFlatBufferList(testFile, func(data []byte) error {
		loadedCount++
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to iterate LZ4 FlatBuffer list: %v", err)
	}

	if loadedCount != len(records) {
		t.Errorf("Record count mismatch: got %d, want %d", loadedCount, len(records))
	}

	t.Logf("Successfully saved and loaded 100 FlatBuffer records with LZ4 compression")
}
