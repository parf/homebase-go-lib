package fileiterator

import (
	"bufio"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// FUOpen opens a file or URL and returns an io.ReadCloser.
// Automatically detects and decompresses files based on extension:
// .gz (gzip), .zst/.zst1/.zst2 (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
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
	} else if strings.HasSuffix(file_or_url, ".zst") || strings.HasSuffix(file_or_url, ".zst1") || strings.HasSuffix(file_or_url, ".zst2") {
		// All zstd levels use the same decoder
		zr, err := zstd.NewReader(base)
		if err != nil {
			base.Close()
			panic(err)
		}
		return &zstdReadCloser{decoder: zr, base: base}
	} else if strings.HasSuffix(file_or_url, ".zlib") || strings.HasSuffix(file_or_url, ".zz") {
		zr, err := zlib.NewReader(base)
		if err != nil {
			base.Close()
			panic(err)
		}
		return &combinedCloser{Reader: zr, closers: []io.Closer{zr, base}}
	} else if strings.HasSuffix(file_or_url, ".lz4") {
		lzr := lz4.NewReader(base)
		return &simpleReadCloser{Reader: lzr, base: base}
	} else if strings.HasSuffix(file_or_url, ".br") {
		brr := brotli.NewReader(base)
		return &simpleReadCloser{Reader: brr, base: base}
	} else if strings.HasSuffix(file_or_url, ".xz") {
		xzr, err := xz.NewReader(base)
		if err != nil {
			base.Close()
			panic(err)
		}
		return &simpleReadCloser{Reader: xzr, base: base}
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

// simpleReadCloser wraps an io.Reader with a base closer
type simpleReadCloser struct {
	io.Reader
	base io.ReadCloser
}

func (s *simpleReadCloser) Close() error {
	return s.base.Close()
}

// FUCreate creates a file and returns an io.WriteCloser.
// Automatically compresses based on file extension:
// .gz (gzip), .zst (zstd default), .zst1 (zstd level 1), .zst2 (zstd level 2),
// .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
func FUCreate(filename string) io.WriteCloser {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Auto-detect compression by extension
	if strings.HasSuffix(filename, ".gz") {
		gzw := gzip.NewWriter(file)
		return &combinedWriteCloser{Writer: gzw, closers: []io.Closer{gzw, file}}
	} else if strings.HasSuffix(filename, ".zst1") {
		// Zstd level 1 - fastest compression
		zw, err := zstd.NewWriter(file, zstd.WithEncoderLevel(zstd.SpeedFastest))
		if err != nil {
			file.Close()
			panic(err)
		}
		return &zstdWriteCloser{encoder: zw, base: file}
	} else if strings.HasSuffix(filename, ".zst2") {
		// Zstd level 2 - fast compression
		zw, err := zstd.NewWriter(file, zstd.WithEncoderLevel(zstd.SpeedDefault))
		if err != nil {
			file.Close()
			panic(err)
		}
		return &zstdWriteCloser{encoder: zw, base: file}
	} else if strings.HasSuffix(filename, ".zst") {
		// Zstd default level (3)
		zw, err := zstd.NewWriter(file)
		if err != nil {
			file.Close()
			panic(err)
		}
		return &zstdWriteCloser{encoder: zw, base: file}
	} else if strings.HasSuffix(filename, ".zlib") || strings.HasSuffix(filename, ".zz") {
		zlibw := zlib.NewWriter(file)
		return &combinedWriteCloser{Writer: zlibw, closers: []io.Closer{zlibw, file}}
	} else if strings.HasSuffix(filename, ".lz4") {
		lzw := lz4.NewWriter(file)
		return &simpleWriteCloser{Writer: lzw, base: file}
	} else if strings.HasSuffix(filename, ".br") {
		brw := brotli.NewWriter(file)
		return &simpleWriteCloser{Writer: brw, base: file}
	} else if strings.HasSuffix(filename, ".xz") {
		xzw, err := xz.NewWriter(file)
		if err != nil {
			file.Close()
			panic(err)
		}
		return &simpleWriteCloser{Writer: xzw, base: file}
	}

	return file
}

// combinedWriteCloser closes multiple closers in sequence
type combinedWriteCloser struct {
	io.Writer
	closers []io.Closer
}

func (c *combinedWriteCloser) Close() error {
	var firstErr error
	for _, closer := range c.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// zstdWriteCloser wraps zstd encoder to implement io.WriteCloser
type zstdWriteCloser struct {
	encoder *zstd.Encoder
	base    io.WriteCloser
}

func (z *zstdWriteCloser) Write(p []byte) (int, error) {
	return z.encoder.Write(p)
}

func (z *zstdWriteCloser) Close() error {
	if err := z.encoder.Close(); err != nil {
		z.base.Close()
		return err
	}
	return z.base.Close()
}

// simpleWriteCloser wraps an io.Writer with a base closer
type simpleWriteCloser struct {
	io.Writer
	base io.WriteCloser
}

func (s *simpleWriteCloser) Close() error {
	// Note: Some writers need to be closed themselves, but we can't
	// determine that from the interface. This works for most cases.
	// For writers that need explicit Close (like lz4.Writer), the
	// combinedWriteCloser should be used instead.
	if closer, ok := s.Writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			s.base.Close()
			return err
		}
	}
	return s.base.Close()
}

// LoadBinFile loads a file with automatic decompression into a byte buffer
// Supported: .gz (gzip), .zst (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
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

// IterateLines processes lines in a file with automatic decompression
// Supported: .gz (gzip), .zst (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
func IterateLines(filename string, processor func(string)) {
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

// IterateIDTabFile iterates over TAB separated (ID <tab> NAME) file with automatic decompression
// ID is parsed as hexadecimal int32, NAME is converted to lowercase
// Supported: .gz (gzip), .zst (zstd), .zlib/.zz (zlib), .lz4 (lz4), .br (brotli), .xz (xz)
func IterateIDTabFile(filename string, processor func(int32, string)) {
	fi := FUOpen(filename) // FUOpen handles compression automatically
	defer fi.Close()
	scanner := bufio.NewScanner(fi)
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

