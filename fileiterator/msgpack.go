package fileiterator

import (
	"bufio"
	"fmt"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

// SaveMsgPack saves data to a MessagePack file with buffered IO
// Uses 4MB buffer for maximum performance
func SaveMsgPack(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Use buffered writer (4MB buffer)
	writer := bufio.NewWriterSize(file, bufferSize)
	defer writer.Flush()

	encoder := msgpack.NewEncoder(writer)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode msgpack: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("MessagePack saved: %s\n", filename)
	return nil
}

// LoadMsgPack loads data from a MessagePack file with buffered IO
// Uses 4MB buffer for maximum performance
func LoadMsgPack(filename string, dest any) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Use buffered reader (4MB buffer)
	reader := bufio.NewReaderSize(file, bufferSize)

	decoder := msgpack.NewDecoder(reader)
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode msgpack: %w", err)
	}

	fmt.Printf("MessagePack loaded: %s\n", filename)
	return nil
}

// SaveMsgPackCompressed saves data to a compressed MessagePack file
// Supports all compression formats via file extension:
// .gz (gzip), .zst (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
// Common usage: filename.msgpack.zst (MessagePack + Zstandard)
func SaveMsgPackCompressed(filename string, data any) error {
	writer := FUCreate(filename) // Auto-detects compression from extension
	defer writer.Close()

	// Use buffered writer for max speed
	bufWriter := bufio.NewWriterSize(writer, bufferSize)
	defer bufWriter.Flush()

	encoder := msgpack.NewEncoder(bufWriter)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode msgpack: %w", err)
	}

	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("MessagePack saved (compressed): %s\n", filename)
	return nil
}

// LoadMsgPackCompressed loads data from a compressed MessagePack file
// Supports all compression formats via FUOpen auto-detection
func LoadMsgPackCompressed(filename string, dest any) error {
	reader := FUOpen(filename) // Auto-detects compression
	defer reader.Close()

	// Use buffered reader for max speed
	bufReader := bufio.NewReaderSize(reader, bufferSize)

	decoder := msgpack.NewDecoder(bufReader)
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode msgpack: %w", err)
	}

	fmt.Printf("MessagePack loaded (compressed): %s\n", filename)
	return nil
}

// IterateMsgPack iterates over a MessagePack stream file
// Each record in the file is decoded and passed to the processor
// Supports compressed files via FUOpen auto-detection
func IterateMsgPack(filename string, processor func(any) error) error {
	reader := FUOpen(filename) // Auto-detects compression
	defer reader.Close()

	// Use buffered reader for max speed
	bufReader := bufio.NewReaderSize(reader, bufferSize)

	decoder := msgpack.NewDecoder(bufReader)
	count := 0

	for {
		var record any
		err := decoder.Decode(&record)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("line %d: decode error: %w", count+1, err)
		}

		count++
		if err := processor(record); err != nil {
			return fmt.Errorf("line %d: processor error: %w", count, err)
		}
	}

	fmt.Printf("File %s. Records processed: %d\n", filename, count)
	return nil
}

// IterateMsgPackTyped iterates over a MessagePack stream file with type safety
// Each record is decoded into type T and passed to the processor
// Supports compressed files via FUOpen auto-detection
func IterateMsgPackTyped[T any](filename string, processor func(T) error) error {
	reader := FUOpen(filename) // Auto-detects compression
	defer reader.Close()

	// Use buffered reader for max speed
	bufReader := bufio.NewReaderSize(reader, bufferSize)

	decoder := msgpack.NewDecoder(bufReader)
	count := 0

	for {
		var record T
		err := decoder.Decode(&record)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("line %d: decode error: %w", count+1, err)
		}

		count++
		if err := processor(record); err != nil {
			return fmt.Errorf("line %d: processor error: %w", count, err)
		}
	}

	fmt.Printf("File %s. Records processed: %d\n", filename, count)
	return nil
}

// SaveMsgPackMap saves a map to MessagePack file with buffered IO
func SaveMsgPackMap(filename string, data map[string]any) error {
	return SaveMsgPack(filename, data)
}

// LoadMsgPackMap loads a map from MessagePack file with buffered IO
func LoadMsgPackMap(filename string) (map[string]any, error) {
	var data map[string]any
	err := LoadMsgPack(filename, &data)
	return data, err
}

// SaveMsgPackMapCompressed saves a map to compressed MessagePack file
// Common usage: data.msgpack.zst (MessagePack + Zstandard)
func SaveMsgPackMapCompressed(filename string, data map[string]any) error {
	return SaveMsgPackCompressed(filename, data)
}

// LoadMsgPackMapCompressed loads a map from compressed MessagePack file
func LoadMsgPackMapCompressed(filename string) (map[string]any, error) {
	var data map[string]any
	err := LoadMsgPackCompressed(filename, &data)
	return data, err
}
