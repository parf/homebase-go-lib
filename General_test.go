package hb_test

import (
	"testing"

	hb "github.com/parf/homebase-go-lib"
)

func TestScale(t *testing.T) {
	tests := []struct {
		input    uint32
		expected byte
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 2},
		{5, 3},
		{16, 3},
		{17, 4},
		{64, 4},
		{65, 5},
		{256, 5},
		{257, 6},
		{1024, 6},
		{1025, 7},
		{4096, 7},
		{4097, 8},
		{16384, 8},
		{16385, 9},
		{100000, 9},
		{999999, 9},
	}

	for _, tt := range tests {
		result := hb.Scale(tt.input)
		if result != tt.expected {
			t.Errorf("Scale(%d) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestScaleBoundaries(t *testing.T) {
	// Test exact boundary values
	boundaries := map[uint32]byte{
		0:     0,
		1:     1,
		4:     2,
		16:    3,
		64:    4,
		256:   5,
		1024:  6,
		4096:  7,
		16384: 8,
	}

	for input, expected := range boundaries {
		result := hb.Scale(input)
		if result != expected {
			t.Errorf("Scale(%d) = %d, expected %d", input, result, expected)
		}
	}
}

func TestScaleMaxValue(t *testing.T) {
	// Test that very large values return 9
	largeValues := []uint32{
		20000,
		50000,
		100000,
		1000000,
		^uint32(0), // Maximum uint32 value
	}

	for _, val := range largeValues {
		result := hb.Scale(val)
		if result != 9 {
			t.Errorf("Scale(%d) = %d, expected 9", val, result)
		}
	}
}
