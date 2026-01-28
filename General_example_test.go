package hb_test

import (
	"fmt"

	hb "github.com/parf/homebase-go-lib"
)

func ExampleScale() {
	// Scale small values
	fmt.Println(hb.Scale(0))
	fmt.Println(hb.Scale(1))
	fmt.Println(hb.Scale(4))
	fmt.Println(hb.Scale(16))
	fmt.Println(hb.Scale(64))
	fmt.Println(hb.Scale(256))
	fmt.Println(hb.Scale(1024))
	fmt.Println(hb.Scale(4096))
	fmt.Println(hb.Scale(16384))
	fmt.Println(hb.Scale(100000))

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleScale_ranges() {
	// Demonstrate ranges
	values := []uint32{0, 1, 5, 17, 65, 257, 1025, 5000, 20000}
	for _, v := range values {
		fmt.Printf("%d -> scale %d\n", v, hb.Scale(v))
	}

	// Output:
	// 0 -> scale 0
	// 1 -> scale 1
	// 5 -> scale 3
	// 17 -> scale 4
	// 65 -> scale 5
	// 257 -> scale 6
	// 1025 -> scale 7
	// 5000 -> scale 8
	// 20000 -> scale 9
}
