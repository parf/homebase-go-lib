package hb

/*
	general methods missing from common GO libraries :)

*/

import (
	"fmt"
	"sort"
)

// Any2uint32 converts various integer types to uint32
func Any2uint32(iii any) (r uint32, err error) {
	switch v := iii.(type) {
	case uint64:
		r = uint32(v)
	case int64:
		r = uint32(v)
	case int32:
		r = uint32(v)
	case uint32:
		r = v
	case int16:
		r = uint32(v)
	case uint16:
		r = uint32(v)
	case int8:
		r = uint32(v)
	case uint8:
		r = uint32(v)
	default:
		return 0, fmt.Errorf("UINT32 typecast error. type: %T value %v", iii, iii)
	}
	return
}

// DumpSortedMap prints a map in key-sorted order
func DumpSortedMap(m map[string]any) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	// Sort the slice of keys
	sort.Strings(keys)

	// Use sorted keys to access map values in order
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, m[k])
	}
}

// HB::Scale($nn, 4) implementation

// Scale returns a 0..9 scale value for uint32 numbers using logarithmic base-4.
// Returns 0 for nn=0, 1 for nn=1, then logBase4(nn)+1 capped at 9.
// This is a copy of HB::scale(nn, 4) method for int16 numbers range.
func Scale(nn uint32) byte {
	var log4_hits = []uint32{0, 1, 4, 16, 64, 256, 1024, 4096, 16384}
	for r, v := range log4_hits {
		if nn <= v {
			return byte(r)
		}
	}
	return 9
}
