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
	data := make(map[string]any)
	ok := cache.Map(filename, data)
	if ok {
		t.Fatal("expected cache miss (false), got hit")
	}
	if len(data) != 0 {
		t.Fatalf("expected empty map on miss, got %v", data)
	}

	// Write
	input := map[string]any{"name": "Alice", "age": int64(30)}
	cache.WriteMap(filename, input)

	// Hit: file now exists
	data = make(map[string]any)
	ok = cache.Map(filename, data)
	if !ok {
		t.Fatal("expected cache hit (true), got miss")
	}
	if data["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", data["name"])
	}
	if data["age"] != int64(30) {
		t.Errorf("expected age=30 (int64), got %v (%T)", data["age"], data["age"])
	}
}

func TestStringKeyedMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "strmap.parquet")

	input := map[string]any{
		"city":   "NYC",
		"score":  float64(99.5),
		"active": true,
		"count":  int32(42),
	}

	cache.WriteMap(filename, input)

	m := make(map[string]any)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}
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

	m := make(map[string]any)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}

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

	m := make(map[string]any)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}

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

	// Read back into same type — no type assertion needed!
	m := make(map[uint32]string)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}
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

	m := make(map[uint32]uint64)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}
	if m[10] != 1000 {
		t.Errorf("key 10: expected 1000, got %v", m[10])
	}
	if m[20] != 2000 {
		t.Errorf("key 20: expected 2000, got %v", m[20])
	}
}

func TestIntAnyMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "intany.parquet")

	input := map[int]any{
		1: float64(1.1),
		2: float64(2.2),
		3: float64(3.3),
	}

	cache.WriteMap(filename, input)

	// int keys promoted to int64 in storage, so read with int64 keys
	m := make(map[int64]any)
	if !cache.Map(filename, m) {
		t.Fatal("expected hit")
	}
	if m[1] != float64(1.1) {
		t.Errorf("key 1: expected 1.1, got %v", m[1])
	}
	if m[2] != float64(2.2) {
		t.Errorf("key 2: expected 2.2, got %v", m[2])
	}
}

func TestCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "corrupt.parquet")

	os.WriteFile(filename, []byte("not a parquet file"), 0644)

	m := make(map[string]any)
	if cache.Map(filename, m) {
		t.Fatal("expected miss for corrupted file")
	}
}

func TestEmptyMapIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "empty.parquet")

	cache.WriteMap(filename, map[string]any{})

	m := make(map[string]any)
	if cache.Map(filename, m) {
		t.Fatal("expected miss for empty map write")
	}
}

func TestBatchIterateKV(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "batch_kv.parquet")

	input := map[uint32]string{
		1: "one", 2: "two", 3: "three", 4: "four", 5: "five",
	}
	cache.WriteMap(filename, input)

	var totalKeys int
	collected := make(map[uint32]string)

	m := make(map[uint32]string)
	ok := cache.MapBatchIterate(filename, m, func() error {
		totalKeys += len(m)
		for k, v := range m {
			collected[k] = v
		}
		return nil
	})
	if !ok {
		t.Fatal("expected successful batch iteration")
	}
	if totalKeys != 5 {
		t.Errorf("expected 5 total keys, got %d", totalKeys)
	}
	for k, v := range input {
		if collected[k] != v {
			t.Errorf("key %d: expected %q, got %q", k, v, collected[k])
		}
	}
}

func TestBatchIterateStringMap(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "batch_map.parquet")

	input := map[string]any{"x": int32(1), "y": int32(2)}
	cache.WriteMap(filename, input)

	calls := 0
	m := make(map[string]any)
	ok := cache.MapBatchIterate(filename, m, func() error {
		calls++
		if m["x"] != int32(1) || m["y"] != int32(2) {
			t.Errorf("unexpected values: %v", m)
		}
		return nil
	})
	if !ok {
		t.Fatal("expected successful iteration")
	}
	if calls != 1 {
		t.Errorf("expected 1 batch call for single-row map, got %d", calls)
	}
}

func TestBatchIterateMissing(t *testing.T) {
	m := make(map[string]any)
	ok := cache.MapBatchIterate("/nonexistent.parquet", m, func() error {
		t.Fatal("should not be called")
		return nil
	})
	if ok {
		t.Fatal("expected false for missing file")
	}
}
