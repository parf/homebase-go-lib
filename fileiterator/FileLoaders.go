package fileiterator

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// FUOpen opens a file or URL and returns an io.ReadCloser.
// Automatically detects and decompresses .gz and .zst files based on extension.
func FUOpen(file_or_url string) io.ReadCloser {
	var base io.ReadCloser

	if strings.HasPrefix(file_or_url, "http") {
		resp, err := http.Get(file_or_url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			err := fmt.Errorf("Url: %s - Unexpected HTTP Code %d", file_or_url, resp.StatusCode)
			panic(err)
		}
		base = resp.Body
	} else {
		r, err := os.Open(file_or_url)
		if err != nil {
			panic(err)
		}
		base = r
	}

	// Auto-detect compression by extension
	if strings.HasSuffix(file_or_url, ".gz") {
		gz, err := gzip.NewReader(base)
		if err != nil {
			base.Close()
			panic(err)
		}
		return &combinedCloser{Reader: gz, closers: []io.Closer{gz, base}}
	} else if strings.HasSuffix(file_or_url, ".zst") {
		zr, err := zstd.NewReader(base)
		if err != nil {
			base.Close()
			panic(err)
		}
		return &zstdReadCloser{decoder: zr, base: base}
	}

	return base
}

// combinedCloser closes multiple closers
type combinedCloser struct {
	io.Reader
	closers []io.Closer
}

func (c *combinedCloser) Close() error {
	var firstErr error
	for _, closer := range c.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// zstdReadCloser wraps zstd decoder to implement io.ReadCloser
type zstdReadCloser struct {
	decoder *zstd.Decoder
	base    io.ReadCloser
}

func (z *zstdReadCloser) Read(p []byte) (int, error) {
	return z.decoder.Read(p)
}

func (z *zstdReadCloser) Close() error {
	z.decoder.Close()
	return z.base.Close()
}

// LoadBinFile loads a file (with automatic decompression if .gz or .zst) into a byte buffer
func LoadBinFile(filename string, dest *[]byte) {
	fi := FUOpen(filename) // FUOpen handles compression automatically
	defer fi.Close()
	var err error
	*dest, err = ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// LoadLinesFile processes lines in a file (with automatic decompression if .gz or .zst)
func LoadLinesFile(filename string, processor func(string)) {
	fi := FUOpen(filename) // FUOpen handles compression automatically
	defer fi.Close()
	scanner := bufio.NewScanner(fi)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// openFileOrURL opens a file or HTTP URL and returns an io.ReadCloser (without decompression)
func openFileOrURL(file_or_url string) io.ReadCloser {
	if strings.HasPrefix(file_or_url, "http") {
		resp, err := http.Get(file_or_url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			err := fmt.Errorf("Url: %s - Unexpected HTTP Code %d", file_or_url, resp.StatusCode)
			panic(err)
		}
		return resp.Body
	}
	r, err := os.Open(file_or_url)
	if err != nil {
		panic(err)
	}
	return r
}

// LoadBinGzFile loads a gzipped file into a byte buffer (explicit gzip)
func LoadBinGzFile(filename string, dest *[]byte) {
	fi := openFileOrURL(filename)
	defer fi.Close()
	fz, err := gzip.NewReader(fi)
	if err != nil {
		panic(err)
	}
	defer fz.Close()
	*dest, err = ioutil.ReadAll(fz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// LoadBinZstdFile loads a zstd-compressed file into a byte buffer (explicit zstd)
func LoadBinZstdFile(filename string, dest *[]byte) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	fz, err := zstd.NewReader(fi)
	if err != nil {
		panic(err)
	}
	defer fz.Close()
	*dest, err = ioutil.ReadAll(fz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// LoadLinesGzFile processes lines in a gzipped file (explicit gzip)
func LoadLinesGzFile(filename string, processor func(string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		fmt.Printf("File %s\n", filename)
		panic(err)
	}
	defer fz.Close()
	scanner := bufio.NewScanner(fz)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// LoadLinesZstdFile processes lines in a zstd-compressed file (explicit zstd)
func LoadLinesZstdFile(filename string, processor func(string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	fz, err := zstd.NewReader(fi)
	if err != nil {
		fmt.Printf("File %s\n", filename)
		panic(err)
	}
	defer fz.Close()
	scanner := bufio.NewScanner(fz)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// LoadIDTabGzFile iterates over TAB separated (ID <tab> NAME) GZIP file
// ID is parsed as hexadecimal int32, NAME is converted to lowercase
func LoadIDTabGzFile(filename string, processor func(int32, string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()
	fz, err := gzip.NewReader(fi)
	if err != nil {
		panic(err)
	}
	defer fz.Close()
	scanner := bufio.NewScanner(fz)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		lc := strings.Split(line, "\t")
		id, err := strconv.ParseInt(lc[0], 16, 32)
		if err != nil {
			panic(err)
		}
		processor(int32(id), strings.ToLower(lc[1]))
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}
