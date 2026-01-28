package hb

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
)

// FUOpen opens a file or URL and returns an io.ReadCloser
func FUOpen(file_or_url string) io.ReadCloser {
	if strings.HasPrefix(file_or_url, "http") {
		resp, err := http.Get(file_or_url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			err := fmt.Errorf("Url: %s - Unexpected HTTP Code %d", file_or_url, resp.StatusCode)
			panic(err)
		}
		return resp.Body // Reader
	}
	r, err := os.Open(file_or_url)
	if err != nil {
		panic(err)
	}
	return r
}

// LoadLinesGzFile processes lines in a gzipped file
func LoadLinesGzFile(filename string, processor func(string)) {
	fi := FUOpen(filename)
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

// LoadBinGzFile loads a gzipped file into a byte buffer
func LoadBinGzFile(filename string, dest *[]byte) {
	fi := FUOpen(filename)
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
	fmt.Printf("File %s loaded. %d bytes\n", filename, len(*dest)) // len(*dest) == size
}

// LoadIDTabGzFile iterates over TAB separated (ID <tab> NAME) GZIP file
func LoadIDTabGzFile(filename string, processor func(int32, string)) {
	fi := FUOpen(filename)
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
