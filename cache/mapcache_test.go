package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/parf/homebase-go-lib/cache"
)

func TestCacheMissAndHit(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.parquet")

	// Miss: file doesn't exist
	data, ok := cache.Map(filename)
	if ok {
		t.Fatal("expected cache miss (false), got hit")
	}
	if data != nil {
		t.Fatalf("expected nil on miss, got %v", data)
	}

	// Write
	input := map[string]any{"name": "Alice", "age": int64(30)}
	cache.WriteMap(filename, input)

	// Hit: file now exists
	data, ok = cache.Map(filename)
	if !ok {
		t.Fatal("expected cache hit (true), got miss")
	}
	m := data.(map[string]any)
	if m["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", m["name"])
	}
	if m["age"] != int64(30) {
		t.Errorf("expected age=30 (int64), got %v (%T)", m["age"], m["age"])
	}
}

func TestStringKeyedMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "strmap.parquet")

	input := map[string]any{
		"city":    "NYC",
		"score":   float64(99.5),
		"active":  true,
		"count":   int32(42),
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	m := data.(map[string]any)
	if m["city"] != "NYC" {
		t.Errorf("city: expected NYC, got %v", m["city"])
	}
	if m["score"] != float64(99.5) {
		t.Errorf("score: expected 99.5, got %v (%T)", m["score"], m["score"])
	}
	if m["active"] != true {
		t.Errorf("active: expected true, got %v", m["active"])
	}
	if m["count"] != int32(42) {
		t.Errorf("count: expected int32(42), got %v (%T)", m["count"], m["count"])
	}
}

func TestRoundTripAllTypes(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "alltypes.parquet")

	input := map[string]any{
		"v_string":  "hello",
		"v_bool":    true,
		"v_int8":    int8(8),
		"v_int16":   int16(16),
		"v_int32":   int32(32),
		"v_int64":   int64(64),
		"v_uint8":   uint8(80),
		"v_uint16":  uint16(160),
		"v_uint32":  uint32(320),
		"v_uint64":  uint64(640),
		"v_float32": float32(3.14),
		"v_float64": float64(6.28),
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	m := data.(map[string]any)

	tests := []struct {
		key      string
		expected any
	}{
		{"v_string", "hello"},
		{"v_bool", true},
		{"v_int8", int8(8)},
		{"v_int16", int16(16)},
		{"v_int32", int32(32)},
		{"v_int64", int64(64)},
		{"v_uint8", uint8(80)},
		{"v_uint16", uint16(160)},
		{"v_uint32", uint32(320)},
		{"v_uint64", uint64(640)},
		{"v_float32", float32(3.14)},
		{"v_float64", float64(6.28)},
	}

	for _, tt := range tests {
		got := m[tt.key]
		if got != tt.expected {
			t.Errorf("key %s: expected %v (%T), got %v (%T)",
				tt.key, tt.expected, tt.expected, got, got)
		}
	}
}

func TestIntAndUintPromotion(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "promote.parquet")

	input := map[string]any{
		"v_int":  int(999),
		"v_uint": uint(888),
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	m := data.(map[string]any)

	// int → int64 on read
	if v, ok := m["v_int"].(int64); !ok || v != 999 {
		t.Errorf("v_int: expected int64(999), got %v (%T)", m["v_int"], m["v_int"])
	}
	// uint → uint64 on read
	if v, ok := m["v_uint"].(uint64); !ok || v != 888 {
		t.Errorf("v_uint: expected uint64(888), got %v (%T)", m["v_uint"], m["v_uint"])
	}
}

func TestUint32StringMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "u32str.parquet")

	input := map[uint32]string{
		1:   "one",
		2:   "two",
		100: "hundred",
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	m := data.(map[uint32]any)
	if m[1] != "one" {
		t.Errorf("key 1: expected 'one', got %v", m[1])
	}
	if m[2] != "two" {
		t.Errorf("key 2: expected 'two', got %v", m[2])
	}
	if m[100] != "hundred" {
		t.Errorf("key 100: expected 'hundred', got %v", m[100])
	}
}

func TestUint32Uint64Map(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "u32u64.parquet")

	input := map[uint32]uint64{
		10: 1000,
		20: 2000,
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	m := data.(map[uint32]any)
	if m[uint32(10)] != uint64(1000) {
		t.Errorf("key 10: expected uint64(1000), got %v (%T)", m[uint32(10)], m[uint32(10)])
	}
	if m[uint32(20)] != uint64(2000) {
		t.Errorf("key 20: expected uint64(2000), got %v (%T)", m[uint32(20)], m[uint32(20)])
	}
}

func TestIntAnyMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "intany.parquet")

	// All values same type → preserved
	input := map[int]any{
		1: float64(1.1),
		2: float64(2.2),
		3: float64(3.3),
	}

	cache.WriteMap(filename, input)

	data, ok := cache.Map(filename)
	if !ok {
		t.Fatal("expected hit")
	}

	// int keys are promoted to int64 in storage
	m := data.(map[int64]any)
	if m[int64(1)] != float64(1.1) {
		t.Errorf("key 1: expected 1.1, got %v", m[int64(1)])
	}
	if m[int64(2)] != float64(2.2) {
		t.Errorf("key 2: expected 2.2, got %v", m[int64(2)])
	}
}

func TestCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "corrupt.parquet")

	os.WriteFile(filename, []byte("not a parquet file"), 0644)

	data, ok := cache.Map(filename)
	if ok {
		t.Fatal("expected miss for corrupted file")
	}
	if data != nil {
		t.Fatalf("expected nil, got %v", data)
	}
}

func TestEmptyMapIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "empty.parquet")

	// Writing empty map should not create file
	cache.WriteMap(filename, map[string]any{})

	_, ok := cache.Map(filename)
	if ok {
		t.Fatal("expected miss for empty map write")
	}
}
