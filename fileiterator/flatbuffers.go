package fileiterator

import (
	"bufio"
	"fmt"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
)

const bufferSize = 4 * 1024 * 1024 // 4MB buffer for max speed

// SaveFlatBuffer saves a FlatBuffer to a file with buffered IO
// Uses 4MB buffer for maximum performance
func SaveFlatBuffer(filename string, builder *flatbuffers.Builder) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Use buffered writer for max speed
	writer := bufio.NewWriterSize(file, bufferSize)
	defer writer.Flush()

	_, err = writer.Write(builder.FinishedBytes())
	if err != nil {
		return fmt.Errorf("failed to write buffer: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("FlatBuffer saved: %s (%d bytes)\n", filename, len(builder.FinishedBytes()))
	return nil
}

// LoadFlatBuffer loads a FlatBuffer from a file with buffered IO
// Returns the raw bytes for zero-copy access
// Uses 4MB buffer for maximum performance
func LoadFlatBuffer(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size for pre-allocation
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Pre-allocate buffer
	data := make([]byte, stat.Size())

	// Use buffered reader for max speed
	reader := bufio.NewReaderSize(file, bufferSize)
	_, err = reader.Read(data)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	fmt.Printf("FlatBuffer loaded: %s (%d bytes)\n", filename, len(data))
	return data, nil
}

// SaveFlatBufferCompressed saves a FlatBuffer to a compressed file
// Supports all compression formats via file extension:
// .gz (gzip), .zst (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
func SaveFlatBufferCompressed(filename string, builder *flatbuffers.Builder) error {
	writer := FUCreate(filename) // Auto-detects compression from extension
	defer writer.Close()

	// Use buffered writer for max speed
	bufWriter := bufio.NewWriterSize(writer, bufferSize)
	defer bufWriter.Flush()

	_, err := bufWriter.Write(builder.FinishedBytes())
	if err != nil {
		return fmt.Errorf("failed to write buffer: %w", err)
	}

	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("FlatBuffer saved (compressed): %s (%d bytes)\n", filename, len(builder.FinishedBytes()))
	return nil
}

// LoadFlatBufferCompressed loads a FlatBuffer from a compressed file
// Supports all compression formats via FUOpen auto-detection
func LoadFlatBufferCompressed(filename string) ([]byte, error) {
	reader := FUOpen(filename) // Auto-detects compression
	defer reader.Close()

	// Use buffered reader for max speed
	bufReader := bufio.NewReaderSize(reader, bufferSize)

	// Read all data
	data := make([]byte, 0, bufferSize)
	buf := make([]byte, bufferSize)
	for {
		n, err := bufReader.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read buffer: %w", err)
		}
	}

	fmt.Printf("FlatBuffer loaded (compressed): %s (%d bytes)\n", filename, len(data))
	return data, nil
}
