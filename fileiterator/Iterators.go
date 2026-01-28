package fileiterator

import (
	"bufio"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/clistat"
)

// IterateBinaryRecords iterates over a file of binary records of fixed recordSize.
// Automatically detects compression format by extension: .gz (gzip), .zst (zstd), .zlib/.zz (zlib)
// If no compression extension is detected, processes file as plain binary.
// Calls processor function on every record.
//
// filename - "filename" or "http://url"
// recordSize - size of each binary record in bytes
func IterateBinaryRecords(filename string, recordSize int, processor func([]byte)) {
	var b io.Reader // base reader
	var closer io.Closer

	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			panic(err)
		}
		b = resp.Body
		closer = resp.Body
	} else {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		b = bufio.NewReader(f)
		closer = f
	}
	defer closer.Close()

	// Detect compression by extension
	var r io.Reader
	var err error

	if strings.HasSuffix(filename, ".gz") {
		r, err = gzip.NewReader(b)
		if err != nil {
			panic(err)
		}
	} else if strings.HasSuffix(filename, ".zst") {
		zr, err := zstd.NewReader(b)
		if err != nil {
			panic(err)
		}
		defer zr.Close()
		r = zr
	} else if strings.HasSuffix(filename, ".zlib") || strings.HasSuffix(filename, ".zz") {
		r, err = zlib.NewReader(b)
		if err != nil {
			panic(err)
		}
	} else {
		// No compression - plain binary file
		r = b
	}

	buf := make([]byte, recordSize)
	stat := clistat.New(10)
	fmt.Printf("Loading: %v\n", filename)
	for {
		n, err := io.ReadFull(r, buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("cnt: %d read:%d\n", stat.Cnt, n)
				fmt.Println(err)
			}
			break
		}
		stat.Hit()
		processor(buf)
	}
	stat.Finish()
}

// IterateZlibRecords iterates over zlib-compressed file of binary records (explicit zlib)
// This is the original function for backward compatibility
func IterateZlibRecords(filename string, recordSize int, processor func([]byte)) {
	var b io.Reader
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			panic(err)
		}
		b = resp.Body
	} else {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b = bufio.NewReader(f)
	}

	r, err := zlib.NewReader(b) // RFC 1950
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf := make([]byte, recordSize)
	stat := clistat.New(10)
	fmt.Printf("Loading: %v\n", filename)
	for {
		n, err := io.ReadFull(r, buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("cnt: %d read:%d\n", stat.Cnt, n)
				fmt.Println(err)
			}
			break
		}
		stat.Hit()
		processor(buf)
	}
	stat.Finish()
}

// IterateGzipRecords iterates over gzip-compressed file of binary records (explicit gzip)
func IterateGzipRecords(filename string, recordSize int, processor func([]byte)) {
	var b io.Reader
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			panic(err)
		}
		b = resp.Body
	} else {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b = bufio.NewReader(f)
	}

	r, err := gzip.NewReader(b)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf := make([]byte, recordSize)
	stat := clistat.New(10)
	fmt.Printf("Loading: %v\n", filename)
	for {
		n, err := io.ReadFull(r, buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("cnt: %d read:%d\n", stat.Cnt, n)
				fmt.Println(err)
			}
			break
		}
		stat.Hit()
		processor(buf)
	}
	stat.Finish()
}

// IterateZstdRecords iterates over zstd-compressed file of binary records (explicit zstd)
func IterateZstdRecords(filename string, recordSize int, processor func([]byte)) {
	var b io.Reader
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			panic(err)
		}
		b = resp.Body
	} else {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b = bufio.NewReader(f)
	}

	r, err := zstd.NewReader(b)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf := make([]byte, recordSize)
	stat := clistat.New(10)
	fmt.Printf("Loading: %v\n", filename)
	for {
		n, err := io.ReadFull(r, buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("cnt: %d read:%d\n", stat.Cnt, n)
				fmt.Println(err)
			}
			break
		}
		stat.Hit()
		processor(buf)
	}
	stat.Finish()
}
