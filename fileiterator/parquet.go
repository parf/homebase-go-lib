package fileiterator

import (
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
)

// ParquetRecord represents a generic record for Parquet operations
// Standard schema: id, name, email, age, score, active, category, timestamp
type ParquetRecord struct {
	ID        int64
	Name      string
	Email     string
	Age       int64
	Score     float64
	Active    bool
	Category  string
	Timestamp int64
}

// IterateParquet reads Parquet file and calls processor for each record
// Automatically handles compression detection via FUOpen
// Schema: id, name, email, age, score, active, category, timestamp
func IterateParquet(filename string, processor func(ParquetRecord) error) error {
	pf, err := file.OpenParquetFile(filename, false)
	if err != nil {
		return err
	}
	defer pf.Close()

	reader, err := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
	if err != nil {
		return err
	}

	tbl, err := reader.ReadTable(nil)
	if err != nil {
		return err
	}
	defer tbl.Release()

	numRows := int(tbl.NumRows())

	// Extract columns
	idCol := tbl.Column(0).Data().Chunk(0).(*array.Int64)
	nameCol := tbl.Column(1).Data().Chunk(0).(*array.String)
	emailCol := tbl.Column(2).Data().Chunk(0).(*array.String)
	ageCol := tbl.Column(3).Data().Chunk(0).(*array.Int64)
	scoreCol := tbl.Column(4).Data().Chunk(0).(*array.Float64)
	activeCol := tbl.Column(5).Data().Chunk(0).(*array.Boolean)
	categoryCol := tbl.Column(6).Data().Chunk(0).(*array.String)
	timestampCol := tbl.Column(7).Data().Chunk(0).(*array.Int64)

	for i := 0; i < numRows; i++ {
		rec := ParquetRecord{
			ID:        idCol.Value(i),
			Name:      nameCol.Value(i),
			Email:     emailCol.Value(i),
			Age:       ageCol.Value(i),
			Score:     scoreCol.Value(i),
			Active:    activeCol.Value(i),
			Category:  categoryCol.Value(i),
			Timestamp: timestampCol.Value(i),
		}
		if err := processor(rec); err != nil {
			return err
		}
	}

	return nil
}

// WriteParquet writes records to Parquet file with Snappy compression
// Automatically handles additional compression via FUCreate (if filename has .gz/.zst/.lz4 extension)
// Schema: id, name, email, age, score, active, category, timestamp
func WriteParquet(filename string, records []ParquetRecord) error {
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

	// Create output file with auto-compression detection
	f := FUCreate(filename)
	defer f.Close()

	// Create Parquet writer with Snappy compression
	writerProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
	arrowProps := pqarrow.DefaultWriterProps()
	writer, err := pqarrow.NewFileWriter(schema, f, writerProps, arrowProps)
	if err != nil {
		return err
	}
	defer writer.Close()

	// Build Arrow record batch
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	for _, record := range records {
		builder.Field(0).(*array.Int64Builder).Append(record.ID)
		builder.Field(1).(*array.StringBuilder).Append(record.Name)
		builder.Field(2).(*array.StringBuilder).Append(record.Email)
		builder.Field(3).(*array.Int64Builder).Append(record.Age)
		builder.Field(4).(*array.Float64Builder).Append(record.Score)
		builder.Field(5).(*array.BooleanBuilder).Append(record.Active)
		builder.Field(6).(*array.StringBuilder).Append(record.Category)
		builder.Field(7).(*array.Int64Builder).Append(record.Timestamp)
	}

	rec := builder.NewRecord()
	defer rec.Release()

	// Write record batch
	if err := writer.Write(rec); err != nil {
		return err
	}

	return nil
}
