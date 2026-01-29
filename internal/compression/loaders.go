package compression

import (
	"bufio"
	"bytes"
	"compress/gzip"
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

// openFileOrURL opens a file or HTTP URL and returns an io.ReadCloser
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

// LoadBinGzFile loads a gzipped file into a byte buffer
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

// LoadBinZstdFile loads a zstd-compressed file into a byte buffer
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

// IterateLinesGz processes lines in a gzipped file (explicit gzip)
func IterateLinesGz(filename string, processor func(string)) {
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

// IterateLinesZstd processes lines in a zstd-compressed file
func IterateLinesZstd(filename string, processor func(string)) {
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

// LoadBinLz4File loads an LZ4-compressed file into a byte buffer
func LoadBinLz4File(filename string, dest *[]byte) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	fz := lz4.NewReader(fi)
	var err error
	*dest, err = ioutil.ReadAll(fz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// IterateLinesLz4 processes lines in an LZ4-compressed file
func IterateLinesLz4(filename string, processor func(string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	fz := lz4.NewReader(fi)
	scanner := bufio.NewScanner(fz)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// LoadBinBrotliFile loads a Brotli-compressed file into a byte buffer
func LoadBinBrotliFile(filename string, dest *[]byte) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	br := brotli.NewReader(fi)
	var err error
	*dest, err = ioutil.ReadAll(br)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// IterateLinesBrotli processes lines in a Brotli-compressed file
func IterateLinesBrotli(filename string, processor func(string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	br := brotli.NewReader(fi)
	scanner := bufio.NewScanner(br)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// LoadBinXzFile loads an XZ-compressed file into a byte buffer
func LoadBinXzFile(filename string, dest *[]byte) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	xzr, err := xz.NewReader(fi)
	if err != nil {
		panic(err)
	}
	*dest, err = ioutil.ReadAll(xzr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest))
}

// IterateLinesXz processes lines in an XZ-compressed file
func IterateLinesXz(filename string, processor func(string)) {
	fi := openFileOrURL(filename)
	defer fi.Close()

	xzr, err := xz.NewReader(fi)
	if err != nil {
		fmt.Printf("File %s\n", filename)
		panic(err)
	}
	scanner := bufio.NewScanner(xzr)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}

// LoadIDTabGzFile iterates over TAB separated (ID <tab> NAME) GZIP file
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
		processor(int32(id), fmt.Sprintf("%s", bytes.ToLower([]byte(lc[1]))))
		count++
	}
	fmt.Printf("File %s. Lines processed: %d\n", filename, count)
}
