package hb

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/parf/homebase-go-lib/clistat"
)

// ZlibFileIterator iterates over Zlib compressed file of binary records of "recordSize"
// and calls recordProcessor on every item
//
// filename - "filename" or "http://url"
func ZlibFileIterator(filename string, recordSize int, processor func([]byte)) {
	var b io.Reader // base reader
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			panic(err)
		}
		b = resp.Body // Reader
	} else {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		b = bufio.NewReader(f) // bufio.NewReader(os.Stdin)
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
		//n, err := r.Read(buf)  -- BAD !!! == https://golang.org/pkg/bufio/#Reader.Read == FUCK
		n, err := io.ReadFull(r, buf) // read recordSize BYTES !!!
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
