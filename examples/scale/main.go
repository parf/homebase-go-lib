package main

import (
	"fmt"

	hb "github.com/parf/homebase-go-lib"
)

func main() {
	fmt.Println("Scale Function Examples")
	fmt.Println("=======================")
	fmt.Println()

	// Test various values
	testValues := []uint32{
		0, 1, 2, 3, 4, 5,
		10, 15, 16, 17,
		50, 64, 65,
		100, 256, 257,
		500, 1024, 1025,
		2000, 4096, 4097,
		10000, 16384, 16385,
		50000, 100000,
	}

	fmt.Println("Value    -> Scale")
	fmt.Println("------------------")
	for _, val := range testValues {
		scale := hb.Scale(val)
		fmt.Printf("%-8d -> %d\n", val, scale)
	}

	fmt.Println()
	fmt.Println("Boundaries:")
	fmt.Println("-----------")
	boundaries := []uint32{0, 1, 4, 16, 64, 256, 1024, 4096, 16384}
	for i, boundary := range boundaries {
		fmt.Printf("Scale %d starts at: %d\n", i, boundary)
	}
}
