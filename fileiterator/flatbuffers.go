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
// Compression formats: .gz, .zst, .zst1, .zst2, .zlib, .zz, .lz4, .br, .xz
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

// IterateFlatBufferList iterates over a FlatBuffer file containing multiple records
// Each record should be prefixed with a 4-byte length (uint32, little-endian)
// Supports compressed files via FUOpen auto-detection (.fb, .fb.gz, .fb.zst, .fb.lz4, etc.)
//
// File format: [length1:uint32][record1:bytes][length2:uint32][record2:bytes]...
func IterateFlatBufferList(filename string, processor func([]byte) error) error {
	reader := FUOpen(filename) // Auto-detects compression
	defer reader.Close()

	// Use buffered reader for max speed
	bufReader := bufio.NewReaderSize(reader, bufferSize)

	count := 0
	lengthBuf := make([]byte, 4)

	for {
		// Read length prefix (4 bytes, little-endian uint32)
		_, err := bufReader.Read(lengthBuf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("record %d: failed to read length: %w", count+1, err)
		}

		// Decode length (little-endian)
		length := uint32(lengthBuf[0]) | uint32(lengthBuf[1])<<8 | uint32(lengthBuf[2])<<16 | uint32(lengthBuf[3])<<24

		// Read record data
		recordData := make([]byte, length)
		n, err := bufReader.Read(recordData)
		if err != nil {
			return fmt.Errorf("record %d: failed to read data: %w", count+1, err)
		}
		if uint32(n) != length {
			return fmt.Errorf("record %d: incomplete read: got %d bytes, expected %d", count+1, n, length)
		}

		count++
		if err := processor(recordData); err != nil {
			return fmt.Errorf("record %d: processor error: %w", count, err)
		}
	}

	fmt.Printf("FlatBuffer list: %s. Records processed: %d\n", filename, count)
	return nil
}

// SaveFlatBufferList saves multiple FlatBuffer records to a file with length prefixes
// Each record is prefixed with a 4-byte length (uint32, little-endian)
// Supports compression via file extension (.fb.gz, .fb.zst, .fb.lz4, etc.)
func SaveFlatBufferList(filename string, records [][]byte) error {
	var writer interface {
		Write([]byte) (int, error)
		Close() error
	}

	// Determine if compression is needed based on extension
	if filename[len(filename)-3:] == ".fb" {
		// Plain file
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		writer = file
	} else {
		// Compressed file
		writer = FUCreate(filename)
	}
	defer writer.Close()

	// Use buffered writer for max speed
	bufWriter := bufio.NewWriterSize(writer, bufferSize)
	defer bufWriter.Flush()

	lengthBuf := make([]byte, 4)
	totalBytes := 0

	for i, record := range records {
		// Write length prefix (little-endian uint32)
		length := uint32(len(record))
		lengthBuf[0] = byte(length)
		lengthBuf[1] = byte(length >> 8)
		lengthBuf[2] = byte(length >> 16)
		lengthBuf[3] = byte(length >> 24)

		_, err := bufWriter.Write(lengthBuf)
		if err != nil {
			return fmt.Errorf("record %d: failed to write length: %w", i+1, err)
		}

		// Write record data
		_, err = bufWriter.Write(record)
		if err != nil {
			return fmt.Errorf("record %d: failed to write data: %w", i+1, err)
		}

		totalBytes += 4 + len(record)
	}

	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("FlatBuffer list saved: %s (%d records, %d bytes)\n", filename, len(records), totalBytes)
	return nil
}
