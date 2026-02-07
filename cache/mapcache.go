package cache

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/metadata"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
)

// Map reads a cached map from a parquet file into dest.
// dest must be a pre-allocated map (e.g., make(map[uint32]string)).
// Returns true on cache hit (dest populated), false on miss.
//
// Usage:
//
//	data := make(map[uint32]string)
//	ok := cache.Map("lookup.parquet", data)
func Map(filename string, dest any) bool {
	return readCacheInto(filename, dest) == nil
}

// WriteMap writes any map to a parquet cache file.
// Supports map[K]V where K and V are scalar types.
// Errors are silently ignored (cache is best-effort).
func WriteMap(filename string, data any) {
	_ = writeCache(filename, data)
}

// --- type inference ---

// inferPreciseType maps a Go value to its exact Arrow DataType.
func inferPreciseType(value any) arrow.DataType {
	switch value.(type) {
	case string:
		return arrow.BinaryTypes.String
	case bool:
		return arrow.FixedWidthTypes.Boolean
	case int8:
		return arrow.PrimitiveTypes.Int8
	case int16:
		return arrow.PrimitiveTypes.Int16
	case int32:
		return arrow.PrimitiveTypes.Int32
	case int64:
		return arrow.PrimitiveTypes.Int64
	case int:
		return arrow.PrimitiveTypes.Int64
	case uint8:
		return arrow.PrimitiveTypes.Uint8
	case uint16:
		return arrow.PrimitiveTypes.Uint16
	case uint32:
		return arrow.PrimitiveTypes.Uint32
	case uint64:
		return arrow.PrimitiveTypes.Uint64
	case uint:
		return arrow.PrimitiveTypes.Uint64
	case float32:
		return arrow.PrimitiveTypes.Float32
	case float64:
		return arrow.PrimitiveTypes.Float64
	default:
		return arrow.BinaryTypes.String
	}
}

// arrowTypeForKind returns the Arrow type for a reflect.Kind.
func arrowTypeForKind(k reflect.Kind) arrow.DataType {
	switch k {
	case reflect.String:
		return arrow.BinaryTypes.String
	case reflect.Bool:
		return arrow.FixedWidthTypes.Boolean
	case reflect.Int8:
		return arrow.PrimitiveTypes.Int8
	case reflect.Int16:
		return arrow.PrimitiveTypes.Int16
	case reflect.Int32:
		return arrow.PrimitiveTypes.Int32
	case reflect.Int64, reflect.Int:
		return arrow.PrimitiveTypes.Int64
	case reflect.Uint8:
		return arrow.PrimitiveTypes.Uint8
	case reflect.Uint16:
		return arrow.PrimitiveTypes.Uint16
	case reflect.Uint32:
		return arrow.PrimitiveTypes.Uint32
	case reflect.Uint64, reflect.Uint:
		return arrow.PrimitiveTypes.Uint64
	case reflect.Float32:
		return arrow.PrimitiveTypes.Float32
	case reflect.Float64:
		return arrow.PrimitiveTypes.Float64
	default:
		return arrow.BinaryTypes.String
	}
}

// --- value appending ---

// appendPreciseValue appends a Go value to the matching Arrow builder.
func appendPreciseValue(builder array.Builder, value any) error {
	if value == nil {
		builder.AppendNull()
		return nil
	}

	switch b := builder.(type) {
	case *array.StringBuilder:
		if s, ok := value.(string); ok {
			b.Append(s)
		} else {
			b.Append(fmt.Sprintf("%v", value))
		}
	case *array.BooleanBuilder:
		b.Append(value.(bool))
	case *array.Int8Builder:
		b.Append(value.(int8))
	case *array.Int16Builder:
		b.Append(value.(int16))
	case *array.Int32Builder:
		b.Append(value.(int32))
	case *array.Int64Builder:
		switch v := value.(type) {
		case int64:
			b.Append(v)
		case int:
			b.Append(int64(v))
		default:
			return fmt.Errorf("cannot append %T to Int64Builder", value)
		}
	case *array.Uint8Builder:
		b.Append(value.(uint8))
	case *array.Uint16Builder:
		b.Append(value.(uint16))
	case *array.Uint32Builder:
		b.Append(value.(uint32))
	case *array.Uint64Builder:
		switch v := value.(type) {
		case uint64:
			b.Append(v)
		case uint:
			b.Append(uint64(v))
		default:
			return fmt.Errorf("cannot append %T to Uint64Builder", value)
		}
	case *array.Float32Builder:
		b.Append(value.(float32))
	case *array.Float64Builder:
		b.Append(value.(float64))
	default:
		return fmt.Errorf("unsupported builder type: %T", builder)
	}
	return nil
}

// appendReflectValue appends a reflect.Value to an Arrow builder using the appropriate type conversion.
func appendReflectValue(builder array.Builder, rv reflect.Value) error {
	switch b := builder.(type) {
	case *array.StringBuilder:
		b.Append(rv.String())
	case *array.BooleanBuilder:
		b.Append(rv.Bool())
	case *array.Int8Builder:
		b.Append(int8(rv.Int()))
	case *array.Int16Builder:
		b.Append(int16(rv.Int()))
	case *array.Int32Builder:
		b.Append(int32(rv.Int()))
	case *array.Int64Builder:
		b.Append(rv.Int())
	case *array.Uint8Builder:
		b.Append(uint8(rv.Uint()))
	case *array.Uint16Builder:
		b.Append(uint16(rv.Uint()))
	case *array.Uint32Builder:
		b.Append(uint32(rv.Uint()))
	case *array.Uint64Builder:
		b.Append(rv.Uint())
	case *array.Float32Builder:
		b.Append(float32(rv.Float()))
	case *array.Float64Builder:
		b.Append(rv.Float())
	default:
		return fmt.Errorf("unsupported builder type: %T", builder)
	}
	return nil
}

// --- value reading ---

// getPreciseValue extracts the native Go value from an Arrow array.
func getPreciseValue(arr arrow.Array, index int) (any, error) {
	if arr.IsNull(index) {
		return nil, nil
	}

	switch a := arr.(type) {
	case *array.String:
		return a.Value(index), nil
	case *array.Boolean:
		return a.Value(index), nil
	case *array.Int8:
		return a.Value(index), nil
	case *array.Int16:
		return a.Value(index), nil
	case *array.Int32:
		return a.Value(index), nil
	case *array.Int64:
		return a.Value(index), nil
	case *array.Uint8:
		return a.Value(index), nil
	case *array.Uint16:
		return a.Value(index), nil
	case *array.Uint32:
		return a.Value(index), nil
	case *array.Uint64:
		return a.Value(index), nil
	case *array.Float32:
		return a.Value(index), nil
	case *array.Float64:
		return a.Value(index), nil
	default:
		return nil, fmt.Errorf("unsupported array type: %T", arr)
	}
}

// --- write path ---

func writeCache(filename string, data any) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("expected map, got %T", data)
	}
	if rv.Len() == 0 {
		return fmt.Errorf("empty map")
	}

	keyKind := rv.Type().Key().Kind()

	if keyKind == reflect.String {
		return writeMapColumns(filename, rv)
	}
	return writeKVColumns(filename, rv)
}

// writeMapColumns writes map[string]any as column-per-key, single row.
func writeMapColumns(filename string, rv reflect.Value) error {
	// Extract keys and values
	keys := make([]string, 0, rv.Len())
	values := make(map[string]any, rv.Len())
	for _, k := range rv.MapKeys() {
		key := k.String()
		keys = append(keys, key)
		val := rv.MapIndex(k).Interface()
		values[key] = val
	}
	sort.Strings(keys)

	// Build schema with precise types
	fields := make([]arrow.Field, len(keys))
	for i, k := range keys {
		fields[i] = arrow.Field{
			Name:     k,
			Type:     inferPreciseType(values[k]),
			Nullable: true,
		}
	}
	md := arrow.MetadataFrom(map[string]string{
		"_cache_format":   "map",
		"_cache_key_type": "string",
	})
	schema := arrow.NewSchema(fields, &md)

	// Write parquet
	return writeParquet(filename, schema, func(builder *array.RecordBuilder) error {
		for i, k := range keys {
			if err := appendPreciseValue(builder.Field(i), values[k]); err != nil {
				return fmt.Errorf("field %s: %w", k, err)
			}
		}
		return nil
	}, 1)
}

// writeKVColumns writes map[K]V as key-value columns, multi-row.
func writeKVColumns(filename string, rv reflect.Value) error {
	mapType := rv.Type()
	keyKind := mapType.Key().Kind()
	keyArrowType := arrowTypeForKind(keyKind)

	// Determine value type: if all values same type, use it; otherwise string
	valType := mapType.Elem()
	isAnyValue := valType.Kind() == reflect.Interface

	var valArrowType arrow.DataType
	if isAnyValue {
		// Detect from actual values
		var firstType arrow.DataType
		allSame := true
		for _, k := range rv.MapKeys() {
			v := rv.MapIndex(k).Elem()
			t := inferPreciseType(v.Interface())
			if firstType == nil {
				firstType = t
			} else if firstType != t {
				allSame = false
				break
			}
		}
		if allSame && firstType != nil {
			valArrowType = firstType
		} else {
			valArrowType = arrow.BinaryTypes.String
		}
	} else {
		valArrowType = arrowTypeForKind(valType.Kind())
	}

	fields := []arrow.Field{
		{Name: "_key", Type: keyArrowType},
		{Name: "_val", Type: valArrowType},
	}

	md := arrow.MetadataFrom(map[string]string{
		"_cache_format":   "kv",
		"_cache_key_type": keyKind.String(),
	})
	schema := arrow.NewSchema(fields, &md)

	numRows := rv.Len()
	// Sort keys for deterministic output
	mapKeys := rv.MapKeys()
	sortReflectKeys(mapKeys)

	return writeParquet(filename, schema, func(builder *array.RecordBuilder) error {
		for _, k := range mapKeys {
			// Append key
			if err := appendReflectValue(builder.Field(0), k); err != nil {
				return fmt.Errorf("key: %w", err)
			}
			// Append value
			v := rv.MapIndex(k)
			if isAnyValue {
				v = v.Elem()
			}
			if valArrowType == arrow.BinaryTypes.String && isAnyValue {
				// Mixed types promoted to string
				builder.Field(1).(*array.StringBuilder).Append(fmt.Sprintf("%v", v.Interface()))
			} else {
				if err := appendReflectValue(builder.Field(1), v); err != nil {
					return fmt.Errorf("value: %w", err)
				}
			}
		}
		return nil
	}, numRows)
}

// writeParquet creates a parquet file with the given schema and data.
func writeParquet(filename string, schema *arrow.Schema, fill func(*array.RecordBuilder) error, expectedRows int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writerProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
	arrowProps := pqarrow.DefaultWriterProps()
	writer, err := pqarrow.NewFileWriter(schema, f, writerProps, arrowProps)
	if err != nil {
		return err
	}
	defer writer.Close()

	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	if err := fill(builder); err != nil {
		return err
	}

	rec := builder.NewRecord()
	defer rec.Release()

	return writer.Write(rec)
}

// --- read path ---

// readCacheInto reads a parquet cache file and populates the dest map.
// The dest map's key/value types determine how parquet data is converted.
func readCacheInto(filename string, dest any) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("dest must be a map, got %T", dest)
	}

	pf, err := file.OpenParquetFile(filename, false)
	if err != nil {
		return err
	}
	defer pf.Close()

	// Read file-level metadata
	kvMeta := pf.MetaData().KeyValueMetadata()
	cacheFormat := getMetaValue(kvMeta, "_cache_format")

	reader, err := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
	if err != nil {
		return err
	}

	tbl, err := reader.ReadTable(context.Background())
	if err != nil {
		return err
	}
	defer tbl.Release()

	if tbl.NumRows() == 0 {
		return fmt.Errorf("empty parquet file")
	}

	destKeyType := rv.Type().Key()
	destValType := rv.Type().Elem()

	if cacheFormat == "kv" {
		return readKVInto(tbl, rv, destKeyType, destValType)
	}
	return readMapInto(tbl, rv, destValType)
}

// getMetaValue reads a value from parquet key-value metadata.
func getMetaValue(kvMeta metadata.KeyValueMetadata, key string) string {
	for i := 0; i < kvMeta.Len(); i++ {
		k := kvMeta.Keys()[i]
		if k == key {
			return kvMeta.Values()[i]
		}
	}
	return ""
}

// readMapInto reads column-per-key format into a string-keyed map.
func readMapInto(tbl arrow.Table, dest reflect.Value, valType reflect.Type) error {
	numCols := int(tbl.NumCols())
	schema := tbl.Schema()

	for colIdx := 0; colIdx < numCols; colIdx++ {
		fieldName := schema.Field(colIdx).Name
		col := tbl.Column(colIdx).Data().Chunk(0)
		value, err := getPreciseValue(col, 0)
		if err != nil {
			return fmt.Errorf("column %s: %w", fieldName, err)
		}
		dest.SetMapIndex(
			reflect.ValueOf(fieldName),
			toReflectValue(value, valType),
		)
	}
	return nil
}

// readKVInto reads key-value columns into a typed map.
func readKVInto(tbl arrow.Table, dest reflect.Value, keyType, valType reflect.Type) error {
	numRows := int(tbl.NumRows())
	keyCol := tbl.Column(0).Data().Chunk(0)
	valCol := tbl.Column(1).Data().Chunk(0)

	for i := 0; i < numRows; i++ {
		k, err := getPreciseValue(keyCol, i)
		if err != nil {
			return err
		}
		v, err := getPreciseValue(valCol, i)
		if err != nil {
			return err
		}
		dest.SetMapIndex(
			toReflectValue(k, keyType),
			toReflectValue(v, valType),
		)
	}
	return nil
}

// toReflectValue converts a value to the target reflect.Type.
// For interface{}/any targets, preserves the precise type.
// For concrete targets, converts using reflect.Convert or fmt.Sprintf fallback.
func toReflectValue(value any, targetType reflect.Type) reflect.Value {
	if value == nil {
		return reflect.Zero(targetType)
	}
	if targetType.Kind() == reflect.Interface {
		return reflect.ValueOf(value)
	}
	rv := reflect.ValueOf(value)
	if rv.Type().ConvertibleTo(targetType) {
		return rv.Convert(targetType)
	}
	// Fallback: stringify for string targets
	if targetType.Kind() == reflect.String {
		return reflect.ValueOf(fmt.Sprintf("%v", value))
	}
	return reflect.Zero(targetType)
}

// --- helpers ---

// sortReflectKeys sorts reflect.Value keys for deterministic output.
func sortReflectKeys(keys []reflect.Value) {
	if len(keys) == 0 {
		return
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		switch a.Kind() {
		case reflect.String:
			return a.String() < b.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() < b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return a.Uint() < b.Uint()
		case reflect.Float32, reflect.Float64:
			return a.Float() < b.Float()
		default:
			return fmt.Sprintf("%v", a.Interface()) < fmt.Sprintf("%v", b.Interface())
		}
	})
}
