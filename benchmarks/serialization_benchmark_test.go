package benchmarks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/parf/homebase-go-lib/fileiterator"
	"github.com/vmihailenco/msgpack/v5"
)

// TestRecord represents a typical data record
type TestRecord struct {
	ID        int64   `json:"id" msgpack:"id"`
	Name      string  `json:"name" msgpack:"name"`
	Email     string  `json:"email" msgpack:"email"`
	Age       int     `json:"age" msgpack:"age"`
	Score     float64 `json:"score" msgpack:"score"`
	Active    bool    `json:"active" msgpack:"active"`
	Category  string  `json:"category" msgpack:"category"`
	Timestamp int64   `json:"timestamp" msgpack:"timestamp"`
}

const (
	numRecords = 1_000_000 // 1 million records
)

var testDataset []TestRecord

// Generate test dataset once
func init() {
	testDataset = make([]TestRecord, numRecords)
	for i := 0; i < numRecords; i++ {
		testDataset[i] = TestRecord{
			ID:        int64(i),
			Name:      "User Name " + string(rune(i%100)),
			Email:     "user" + string(rune(i%100)) + "@example.com",
			Age:       20 + (i % 60),
			Score:     float64(i%100) + 0.5,
			Active:    i%2 == 0,
			Category:  "Category" + string(rune(i%10)),
			Timestamp: 1640000000 + int64(i),
		}
	}
}

// ============================================================================
// WRITE BENCHMARKS
// ============================================================================

func BenchmarkWrite_JSONL_Plain(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.jsonl")
		f, _ := os.Create(file)
		encoder := json.NewEncoder(f)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		f.Close()
	}
}

func BenchmarkWrite_JSONL_Gzip(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.jsonl.gz")
		w := fileiterator.FUCreate(file)
		encoder := json.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_JSONL_Zstd(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.jsonl.zst")
		w := fileiterator.FUCreate(file)
		encoder := json.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_JSONL_LZ4(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.jsonl.lz4")
		w := fileiterator.FUCreate(file)
		encoder := json.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_MsgPack_Plain(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.msgpack")
		f, _ := os.Create(file)
		encoder := msgpack.NewEncoder(f)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		f.Close()
	}
}

func BenchmarkWrite_MsgPack_Gzip(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.msgpack.gz")
		w := fileiterator.FUCreate(file)
		encoder := msgpack.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_MsgPack_Zstd(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.msgpack.zst")
		w := fileiterator.FUCreate(file)
		encoder := msgpack.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_MsgPack_LZ4(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.msgpack.lz4")
		w := fileiterator.FUCreate(file)
		encoder := msgpack.NewEncoder(w)
		for _, record := range testDataset {
			encoder.Encode(record)
		}
		w.Close()
	}
}

func BenchmarkWrite_FlatBuffer_Plain(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.fb")
		// Create a single FlatBuffer with all records
		builder := flatbuffers.NewBuilder(1024 * 1024 * 100) // 100MB initial

		// For simplicity, just store as byte vectors (in real use, define proper schema)
		offsets := make([]flatbuffers.UOffsetT, numRecords)
		for j, record := range testDataset {
			data, _ := json.Marshal(record)
			offsets[j] = builder.CreateByteVector(data)
		}

		builder.Finish(offsets[0])
		fileiterator.SaveFlatBuffer(file, builder)
	}
}

func BenchmarkWrite_FlatBuffer_Zstd(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, "test.fb.zst")
		builder := flatbuffers.NewBuilder(1024 * 1024 * 100)

		offsets := make([]flatbuffers.UOffsetT, numRecords)
		for j, record := range testDataset {
			data, _ := json.Marshal(record)
			offsets[j] = builder.CreateByteVector(data)
		}

		builder.Finish(offsets[0])
		fileiterator.SaveFlatBufferCompressed(file, builder)
	}
}

func BenchmarkWrite_Parquet(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < b.N; i++ {
		filename := filepath.Join(tmpDir, "test.parquet")

		// Create Arrow schema
		schema := arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "name", Type: arrow.BinaryTypes.String},
				{Name: "email", Type: arrow.BinaryTypes.String},
				{Name: "age", Type: arrow.PrimitiveTypes.Int64},
				{Name: "score", Type: arrow.PrimitiveTypes.Float64},
				{Name: "active", Type: arrow.FixedWidthTypes.Boolean},
				{Name: "category", Type: arrow.BinaryTypes.String},
				{Name: "timestamp", Type: arrow.PrimitiveTypes.Int64},
			},
			nil,
		)

		// Create output file
		f, _ := os.Create(filename)
		defer f.Close()

		// Create Parquet writer with Snappy compression
		writerProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
		arrowProps := pqarrow.DefaultWriterProps()
		writer, _ := pqarrow.NewFileWriter(schema, f, writerProps, arrowProps)
		defer writer.Close()

		// Build Arrow record batch
		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, schema)
		defer builder.Release()

		for _, record := range testDataset {
			builder.Field(0).(*array.Int64Builder).Append(record.ID)
			builder.Field(1).(*array.StringBuilder).Append(record.Name)
			builder.Field(2).(*array.StringBuilder).Append(record.Email)
			builder.Field(3).(*array.Int64Builder).Append(int64(record.Age))
			builder.Field(4).(*array.Float64Builder).Append(record.Score)
			builder.Field(5).(*array.BooleanBuilder).Append(record.Active)
			builder.Field(6).(*array.StringBuilder).Append(record.Category)
			builder.Field(7).(*array.Int64Builder).Append(record.Timestamp)
		}

		rec := builder.NewRecord()
		defer rec.Release()

		// Write record batch
		writer.Write(rec)
	}
}

// ============================================================================
// READ BENCHMARKS
// ============================================================================

func BenchmarkRead_JSONL_Plain(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.jsonl")

	// Setup: Write file once
	f, _ := os.Create(file)
	encoder := json.NewEncoder(f)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	f.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateJSONLTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_JSONL_Gzip(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.jsonl.gz")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := json.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateJSONLTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_JSONL_Zstd(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.jsonl.zst")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := json.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateJSONLTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_JSONL_LZ4(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.jsonl.lz4")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := json.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateJSONLTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_MsgPack_Plain(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.msgpack")

	// Setup: Write file once
	f, _ := os.Create(file)
	encoder := msgpack.NewEncoder(f)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	f.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateMsgPackTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_MsgPack_Gzip(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.msgpack.gz")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := msgpack.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateMsgPackTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_MsgPack_Zstd(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.msgpack.zst")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := msgpack.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateMsgPackTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_MsgPack_LZ4(b *testing.B) {
	tmpDir := b.TempDir()
	file := filepath.Join(tmpDir, "test.msgpack.lz4")

	// Setup: Write file once
	w := fileiterator.FUCreate(file)
	encoder := msgpack.NewEncoder(w)
	for _, record := range testDataset {
		encoder.Encode(record)
	}
	w.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		fileiterator.IterateMsgPackTyped(file, func(record TestRecord) error {
			count++
			return nil
		})
	}
}

func BenchmarkRead_Parquet(b *testing.B) {
	tmpDir := b.TempDir()
	filename := filepath.Join(tmpDir, "test.parquet")

	// Setup: Write file once
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "email", Type: arrow.BinaryTypes.String},
			{Name: "age", Type: arrow.PrimitiveTypes.Int64},
			{Name: "score", Type: arrow.PrimitiveTypes.Float64},
			{Name: "active", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "category", Type: arrow.BinaryTypes.String},
			{Name: "timestamp", Type: arrow.PrimitiveTypes.Int64},
		},
		nil,
	)

	f, _ := os.Create(filename)
	writerProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
	arrowProps := pqarrow.DefaultWriterProps()
	writer, _ := pqarrow.NewFileWriter(schema, f, writerProps, arrowProps)

	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)

	for _, record := range testDataset {
		builder.Field(0).(*array.Int64Builder).Append(record.ID)
		builder.Field(1).(*array.StringBuilder).Append(record.Name)
		builder.Field(2).(*array.StringBuilder).Append(record.Email)
		builder.Field(3).(*array.Int64Builder).Append(int64(record.Age))
		builder.Field(4).(*array.Float64Builder).Append(record.Score)
		builder.Field(5).(*array.BooleanBuilder).Append(record.Active)
		builder.Field(6).(*array.StringBuilder).Append(record.Category)
		builder.Field(7).(*array.Int64Builder).Append(record.Timestamp)
	}

	rec := builder.NewRecord()
	writer.Write(rec)
	rec.Release()
	builder.Release()
	writer.Close()
	f.Close()

	// Benchmark reading
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pf, _ := file.OpenParquetFile(filename, false)
		reader, _ := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
		tbl, _ := reader.ReadTable(context.Background())
		tbl.Release()
		pf.Close()
	}
}
