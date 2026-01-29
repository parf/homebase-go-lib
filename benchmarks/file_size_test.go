package benchmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/parf/homebase-go-lib/fileiterator"
	"github.com/vmihailenco/msgpack/v5"
)

func TestGenerateFileSizes(t *testing.T) {
	tmpDir := t.TempDir()

	formats := []struct {
		name string
		ext  string
	}{
		{"JSONL Plain", ".jsonl"},
		{"JSONL Gzip", ".jsonl.gz"},
		{"JSONL Zstd", ".jsonl.zst"},
		{"JSONL LZ4", ".jsonl.lz4"},
		{"JSONL Brotli", ".jsonl.br"},
		{"JSONL XZ", ".jsonl.xz"},
		{"MsgPack Plain", ".msgpack"},
		{"MsgPack Gzip", ".msgpack.gz"},
		{"MsgPack Zstd", ".msgpack.zst"},
		{"MsgPack LZ4", ".msgpack.lz4"},
		{"MsgPack Brotli", ".msgpack.br"},
		{"MsgPack XZ", ".msgpack.xz"},
	}

	fmt.Println("\n=== File Size Comparison (1M records) ===")
	fmt.Printf("%-20s %15s %10s\n", "Format", "Size (bytes)", "Size (MB)")
	fmt.Println("-----------------------------------------------------------")

	// JSONL formats
	for _, format := range formats {
		if format.ext == ".jsonl" || (len(format.ext) > 6 && format.ext[0:6] == ".jsonl") {
			file := filepath.Join(tmpDir, "test"+format.ext)

			var w interface{ Write([]byte) (int, error); Close() error }
			if format.ext == ".jsonl" {
				f, _ := os.Create(file)
				w = f
			} else {
				w = fileiterator.FUCreate(file)
			}

			encoder := json.NewEncoder(w)
			for _, record := range testDataset {
				encoder.Encode(record)
			}
			w.Close()

			stat, _ := os.Stat(file)
			sizeMB := float64(stat.Size()) / 1024 / 1024
			fmt.Printf("%-20s %15d %10.2f\n", format.name, stat.Size(), sizeMB)
		}
	}

	// MsgPack formats
	for _, format := range formats {
		if format.ext == ".msgpack" || (len(format.ext) > 8 && format.ext[0:8] == ".msgpack") {
			file := filepath.Join(tmpDir, "test"+format.ext)

			var w interface{ Write([]byte) (int, error); Close() error }
			if format.ext == ".msgpack" {
				f, _ := os.Create(file)
				w = f
			} else {
				w = fileiterator.FUCreate(file)
			}

			encoder := msgpack.NewEncoder(w)
			for _, record := range testDataset {
				encoder.Encode(record)
			}
			w.Close()

			stat, _ := os.Stat(file)
			sizeMB := float64(stat.Size()) / 1024 / 1024
			fmt.Printf("%-20s %15d %10.2f\n", format.name, stat.Size(), sizeMB)
		}
	}

	// FlatBuffer formats
	fbFormats := []struct {
		name string
		ext  string
	}{
		{"FlatBuffer Plain", ".fb"},
		{"FlatBuffer Zstd", ".fb.zst"},
	}

	for _, format := range fbFormats {
		file := filepath.Join(tmpDir, "test"+format.ext)
		builder := flatbuffers.NewBuilder(1024 * 1024 * 100)

		offsets := make([]flatbuffers.UOffsetT, numRecords)
		for j, record := range testDataset {
			data, _ := json.Marshal(record)
			offsets[j] = builder.CreateByteVector(data)
		}

		builder.Finish(offsets[0])

		if format.ext == ".fb" {
			fileiterator.SaveFlatBuffer(file, builder)
		} else {
			fileiterator.SaveFlatBufferCompressed(file, builder)
		}

		stat, _ := os.Stat(file)
		sizeMB := float64(stat.Size()) / 1024 / 1024
		fmt.Printf("%-20s %15d %10.2f\n", format.name, stat.Size(), sizeMB)
	}

	fmt.Println()
}
