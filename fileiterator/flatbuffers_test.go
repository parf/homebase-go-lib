package fileiterator_test

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func TestSaveFlatBuffer(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fb")

	// Create a simple FlatBuffer
	builder := flatbuffers.NewBuilder(1024)
	testData := []byte("Hello FlatBuffers!")
	dataOffset := builder.CreateByteVector(testData)
	builder.Finish(dataOffset)

	// Save
	err := fileiterator.SaveFlatBuffer(testFile, builder)
	if err != nil {
		t.Fatalf("Failed to save FlatBuffer: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}

func TestLoadFlatBuffer(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fb")

	// Create and save a FlatBuffer
	builder := flatbuffers.NewBuilder(1024)
	testData := []byte("Hello FlatBuffers!")
	dataOffset := builder.CreateByteVector(testData)
	builder.Finish(dataOffset)

	err := fileiterator.SaveFlatBuffer(testFile, builder)
	if err != nil {
		t.Fatalf("Failed to save FlatBuffer: %v", err)
	}

	// Load
	data, err := fileiterator.LoadFlatBuffer(testFile)
	if err != nil {
		t.Fatalf("Failed to load FlatBuffer: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Loaded data is empty")
	}

	// Verify data matches
	originalBytes := builder.FinishedBytes()
	if !bytes.Equal(data, originalBytes) {
		t.Errorf("Loaded data doesn't match original")
	}
}

func TestSaveFlatBufferCompressed(t *testing.T) {
	tmpDir := t.TempDir()

	testFormats := []struct {
		name string
		ext  string
	}{
		{"Gzip", ".fb.gz"},
		{"Zstd", ".fb.zst"},
		{"Zlib", ".fb.zlib"},
		{"LZ4", ".fb.lz4"},
		{"Brotli", ".fb.br"},
		{"XZ", ".fb.xz"},
	}

	for _, tf := range testFormats {
		t.Run(tf.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test"+tf.ext)

			// Create a FlatBuffer
			builder := flatbuffers.NewBuilder(1024)
			testData := []byte("Compressed FlatBuffer Test!")
			dataOffset := builder.CreateByteVector(testData)
			builder.Finish(dataOffset)
			originalBytes := builder.FinishedBytes()

			// Save compressed
			err := fileiterator.SaveFlatBufferCompressed(testFile, builder)
			if err != nil {
				t.Fatalf("Failed to save compressed FlatBuffer: %v", err)
			}

			// Load compressed
			data, err := fileiterator.LoadFlatBufferCompressed(testFile)
			if err != nil {
				t.Fatalf("Failed to load compressed FlatBuffer: %v", err)
			}

			// Verify data matches
			if !bytes.Equal(data, originalBytes) {
				t.Errorf("Loaded data doesn't match original")
			}
		})
	}
}

func TestLoadFlatBufferGz(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fb.gz")

	// Create a FlatBuffer
	builder := flatbuffers.NewBuilder(1024)
	testData := []byte("Hello Compressed FlatBuffers!")
	dataOffset := builder.CreateByteVector(testData)
	builder.Finish(dataOffset)
	originalBytes := builder.FinishedBytes()

	// Manually create compressed file
	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	gz.Write(originalBytes)
	gz.Close()
	f.Close()

	// Load compressed
	data, err := fileiterator.LoadFlatBufferCompressed(testFile)
	if err != nil {
		t.Fatalf("Failed to load compressed FlatBuffer: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Loaded data is empty")
	}

	// Verify data matches
	if !bytes.Equal(data, originalBytes) {
		t.Errorf("Loaded data doesn't match original")
	}
}
