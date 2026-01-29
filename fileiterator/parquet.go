package fileiterator

import (
	"fmt"
	"sort"

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

// WriteParquetAny writes generic records to Parquet file with Snappy compression
// Automatically infers schema from data - supports ANY record structure
// Handles compression via FUCreate (if filename has .gz/.zst/.lz4 extension)
// Supported types: int64, float64, string, bool
func WriteParquetAny(filename string, records []map[string]any) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to write")
	}

	// Infer schema from all records
	schema, fieldOrder, err := inferSchema(records)
	if err != nil {
		return err
	}

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

	// Append records
	for _, record := range records {
		for i, fieldName := range fieldOrder {
			value := record[fieldName]
			if err := appendValue(builder.Field(i), value); err != nil {
				return fmt.Errorf("error appending field %s: %w", fieldName, err)
			}
		}
	}

	rec := builder.NewRecord()
	defer rec.Release()

	// Write record batch
	if err := writer.Write(rec); err != nil {
		return err
	}

	return nil
}

// inferSchema infers Arrow schema from records
// Returns schema and field order for consistent field ordering
func inferSchema(records []map[string]any) (*arrow.Schema, []string, error) {
	if len(records) == 0 {
		return nil, nil, fmt.Errorf("cannot infer schema from empty records")
	}

	// Collect all field names and infer types
	fieldTypes := make(map[string]arrow.DataType)

	// Sample first record to get field names
	for key, value := range records[0] {
		fieldTypes[key] = inferType(value)
	}

	// Verify all records have compatible types
	for _, record := range records[1:] {
		for key, value := range record {
			if value == nil {
				continue // Skip nil values
			}
			inferredType := inferType(value)
			if existing, ok := fieldTypes[key]; ok {
				// Check type compatibility
				if existing != inferredType {
					// Promote to string if types don't match
					fieldTypes[key] = arrow.BinaryTypes.String
				}
			} else {
				// New field found in later record
				fieldTypes[key] = inferredType
			}
		}
	}

	// Create sorted field list for deterministic ordering
	fieldNames := make([]string, 0, len(fieldTypes))
	for name := range fieldTypes {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	// Build Arrow fields
	fields := make([]arrow.Field, len(fieldNames))
	for i, name := range fieldNames {
		fields[i] = arrow.Field{
			Name: name,
			Type: fieldTypes[name],
		}
	}

	schema := arrow.NewSchema(fields, nil)
	return schema, fieldNames, nil
}

// inferType infers Arrow type from Go value
func inferType(value any) arrow.DataType {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return arrow.PrimitiveTypes.Int64
	case uint, uint8, uint16, uint32, uint64:
		return arrow.PrimitiveTypes.Int64
	case float32, float64:
		return arrow.PrimitiveTypes.Float64
	case bool:
		return arrow.FixedWidthTypes.Boolean
	case string:
		return arrow.BinaryTypes.String
	default:
		// Default to string for unknown types
		return arrow.BinaryTypes.String
	}
}

// appendValue appends a value to an Arrow array builder
func appendValue(builder array.Builder, value any) error {
	if value == nil {
		builder.AppendNull()
		return nil
	}

	switch b := builder.(type) {
	case *array.Int64Builder:
		switch v := value.(type) {
		case int:
			b.Append(int64(v))
		case int64:
			b.Append(v)
		case int32:
			b.Append(int64(v))
		case int16:
			b.Append(int64(v))
		case int8:
			b.Append(int64(v))
		case uint:
			b.Append(int64(v))
		case uint64:
			b.Append(int64(v))
		case uint32:
			b.Append(int64(v))
		case uint16:
			b.Append(int64(v))
		case uint8:
			b.Append(int64(v))
		case float64:
			b.Append(int64(v))
		default:
			return fmt.Errorf("cannot convert %T to int64", value)
		}
	case *array.Float64Builder:
		switch v := value.(type) {
		case float64:
			b.Append(v)
		case float32:
			b.Append(float64(v))
		case int:
			b.Append(float64(v))
		case int64:
			b.Append(float64(v))
		default:
			return fmt.Errorf("cannot convert %T to float64", value)
		}
	case *array.BooleanBuilder:
		switch v := value.(type) {
		case bool:
			b.Append(v)
		default:
			return fmt.Errorf("cannot convert %T to bool", value)
		}
	case *array.StringBuilder:
		switch v := value.(type) {
		case string:
			b.Append(v)
		default:
			// Convert anything to string
			b.Append(fmt.Sprintf("%v", value))
		}
	default:
		return fmt.Errorf("unsupported builder type: %T", builder)
	}

	return nil
}

// IterateParquetAny reads Parquet file and calls processor for each record as map[string]any
// Automatically handles compression detection via file extension
// Supports ANY Parquet schema
func IterateParquetAny(filename string, processor func(map[string]any) error) error {
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
	numCols := int(tbl.NumCols())
	schema := tbl.Schema()

	// Process each row
	for i := 0; i < numRows; i++ {
		record := make(map[string]any)

		for colIdx := 0; colIdx < numCols; colIdx++ {
			fieldName := schema.Field(colIdx).Name
			col := tbl.Column(colIdx).Data().Chunk(0)

			// Extract value based on column type
			value, err := getValueFromColumn(col, i)
			if err != nil {
				return fmt.Errorf("error reading column %s at row %d: %w", fieldName, i, err)
			}

			record[fieldName] = value
		}

		if err := processor(record); err != nil {
			return err
		}
	}

	return nil
}

// getValueFromColumn extracts value from Arrow array at given index
func getValueFromColumn(arr arrow.Array, index int) (any, error) {
	if arr.IsNull(index) {
		return nil, nil
	}

	switch a := arr.(type) {
	case *array.Int64:
		return a.Value(index), nil
	case *array.Int32:
		return int64(a.Value(index)), nil
	case *array.Float64:
		return a.Value(index), nil
	case *array.Float32:
		return float64(a.Value(index)), nil
	case *array.Boolean:
		return a.Value(index), nil
	case *array.String:
		return a.Value(index), nil
	default:
		return nil, fmt.Errorf("unsupported array type: %T", arr)
	}
}
