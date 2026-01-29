package fileiterator_test

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/parf/homebase-go-lib/fileiterator"
	"github.com/vmihailenco/msgpack/v5"
)

func TestSaveMsgPack(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack")

	testData := map[string]any{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}

	err := fileiterator.SaveMsgPack(testFile, testData)
	if err != nil {
		t.Fatalf("Failed to save MessagePack: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}

func TestLoadMsgPack(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack")

	testData := map[string]any{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}

	err := fileiterator.SaveMsgPack(testFile, testData)
	if err != nil {
		t.Fatalf("Failed to save MessagePack: %v", err)
	}

	var loadedData map[string]any
	err = fileiterator.LoadMsgPack(testFile, &loadedData)
	if err != nil {
		t.Fatalf("Failed to load MessagePack: %v", err)
	}

	if loadedData["name"] != testData["name"] {
		t.Errorf("Name mismatch: got %v, want %v", loadedData["name"], testData["name"])
	}
}

func TestSaveLoadMsgPackMap(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack")

	testData := map[string]any{
		"name":  "Jane Doe",
		"age":   25,
		"city":  "New York",
		"score": 95.5,
	}

	err := fileiterator.SaveMsgPackMap(testFile, testData)
	if err != nil {
		t.Fatalf("Failed to save MessagePack map: %v", err)
	}

	loadedData, err := fileiterator.LoadMsgPackMap(testFile)
	if err != nil {
		t.Fatalf("Failed to load MessagePack map: %v", err)
	}

	if loadedData["name"] != testData["name"] {
		t.Errorf("Name mismatch")
	}
	if loadedData["city"] != testData["city"] {
		t.Errorf("City mismatch")
	}
}

func TestLoadMsgPackCompressed(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack.gz")

	testData := map[string]any{
		"name":    "Compressed User",
		"version": 1,
	}

	// Manually create compressed MessagePack file
	f, _ := os.Create(testFile)
	gz := gzip.NewWriter(f)
	encoder := msgpack.NewEncoder(gz)
	encoder.Encode(testData)
	gz.Close()
	f.Close()

	// Load compressed
	var loadedData map[string]any
	err := fileiterator.LoadMsgPackCompressed(testFile, &loadedData)
	if err != nil {
		t.Fatalf("Failed to load compressed MessagePack: %v", err)
	}

	if loadedData["name"] != testData["name"] {
		t.Errorf("Name mismatch in compressed data")
	}
}

func TestIterateMsgPack(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack")

	// Create a stream of MessagePack records
	f, _ := os.Create(testFile)
	encoder := msgpack.NewEncoder(f)
	records := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Charlie"},
	}
	for _, record := range records {
		encoder.Encode(record)
	}
	f.Close()

	// Iterate
	count := 0
	err := fileiterator.IterateMsgPack(testFile, func(record any) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to iterate MessagePack: %v", err)
	}

	if count != len(records) {
		t.Errorf("Expected %d records, got %d", len(records), count)
	}
}

func TestIterateMsgPackTyped(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack")

	type User struct {
		ID   int    `msgpack:"id"`
		Name string `msgpack:"name"`
	}

	// Create a stream of MessagePack records
	f, _ := os.Create(testFile)
	encoder := msgpack.NewEncoder(f)
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	for _, user := range users {
		encoder.Encode(user)
	}
	f.Close()

	// Iterate with type safety
	count := 0
	err := fileiterator.IterateMsgPackTyped(testFile, func(user User) error {
		count++
		if user.ID == 0 || user.Name == "" {
			t.Errorf("Invalid user data: %+v", user)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to iterate typed MessagePack: %v", err)
	}

	if count != len(users) {
		t.Errorf("Expected %d records, got %d", len(users), count)
	}
}

func TestSaveMsgPackCompressed(t *testing.T) {
	tmpDir := t.TempDir()

	testFormats := []struct {
		name string
		ext  string
	}{
		{"Gzip", ".msgpack.gz"},
		{"Zstd", ".msgpack.zst"},
		{"Zlib", ".msgpack.zlib"},
		{"LZ4", ".msgpack.lz4"},
		{"Brotli", ".msgpack.br"},
		{"XZ", ".msgpack.xz"},
	}

	testData := map[string]any{
		"name":    "Compressed Test",
		"version": 2,
		"active":  true,
	}

	for _, tf := range testFormats {
		t.Run(tf.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test"+tf.ext)

			// Save compressed
			err := fileiterator.SaveMsgPackCompressed(testFile, testData)
			if err != nil {
				t.Fatalf("Failed to save compressed MessagePack: %v", err)
			}

			// Load compressed
			var loadedData map[string]any
			err = fileiterator.LoadMsgPackCompressed(testFile, &loadedData)
			if err != nil {
				t.Fatalf("Failed to load compressed MessagePack: %v", err)
			}

			// Verify data matches
			if loadedData["name"] != testData["name"] {
				t.Errorf("Name mismatch in compressed data: got %v, want %v", loadedData["name"], testData["name"])
			}
			// MessagePack encodes integers as int64/int8 depending on value
			// Compare as numbers, not strict type equality
			version, ok := loadedData["version"].(int8)
			if !ok {
				t.Errorf("Version type mismatch: got %T", loadedData["version"])
			} else if version != 2 {
				t.Errorf("Version mismatch in compressed data: got %v, want 2", version)
			}
		})
	}
}

func TestIterateMsgPackCompressed(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.msgpack.zst")

	// Create compressed MessagePack stream
	f, _ := os.Create(testFile)
	zw, _ := zstd.NewWriter(f)
	encoder := msgpack.NewEncoder(zw)
	records := []map[string]any{
		{"id": 1, "value": "First"},
		{"id": 2, "value": "Second"},
	}
	for _, record := range records {
		encoder.Encode(record)
	}
	zw.Close()
	f.Close()

	// Iterate compressed stream
	count := 0
	err := fileiterator.IterateMsgPack(testFile, func(record any) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to iterate compressed MessagePack: %v", err)
	}

	if count != len(records) {
		t.Errorf("Expected %d records, got %d", len(records), count)
	}
}
