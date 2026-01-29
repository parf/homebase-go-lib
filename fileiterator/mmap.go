package fileiterator

import (
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

// MmapFile represents a memory-mapped file with lifecycle management
type MmapFile struct {
	Data []byte
	mmap mmap.MMap
	file *os.File
}

// Close unmaps the memory and closes the file
func (m *MmapFile) Close() error {
	var errs []error

	if m.mmap != nil {
		if err := m.mmap.Unmap(); err != nil {
			errs = append(errs, fmt.Errorf("failed to unmap: %w", err))
		}
	}

	if m.file != nil {
		if err := m.file.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close file: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}

	return nil
}

// MmapOpen opens a file with memory mapping for ultra-fast read access
// Returns MmapFile which must be closed by the caller
// Best for: Large read-heavy files where random access is needed
// Not suitable for: Compressed files, streaming, or write operations
func MmapOpen(filename string) (*MmapFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	mmapData, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to map file: %w", err)
	}

	fmt.Printf("Memory-mapped file: %s (%d bytes)\n", filename, len(mmapData))

	return &MmapFile{
		Data: mmapData,
		mmap: mmapData,
		file: file,
	}, nil
}

// LoadMmap loads a file using memory mapping for ultra-fast access
// Returns the data as a byte slice. The underlying memory mapping is managed
// internally and will be cleaned up when appropriate.
//
// WARNING: The returned slice is backed by mmap. Do not use it after the
// program terminates or if you need guaranteed persistence across process restarts.
// For long-lived data that outlives the file access, copy the bytes.
//
// Best for: Large read-heavy files where random access is needed
// Not suitable for: Compressed files, streaming, or write operations
func LoadMmap(filename string) ([]byte, error) {
	mmapFile, err := MmapOpen(filename)
	if err != nil {
		return nil, err
	}

	// Note: We return the data but keep the mmap active.
	// The OS will handle cleanup, but for explicit control use MmapOpen/Close
	return mmapFile.Data, nil
}
