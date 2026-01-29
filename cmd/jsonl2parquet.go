package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	inputFile := flag.String("input", "", "Input JSONL file (supports .gz, .zst, .lz4, .br, .xz)")
	outputFile := flag.String("output", "", "Output Parquet file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: jsonl2parquet -input <jsonl_file> -output <parquet_file>")
		fmt.Println("\nExample:")
		fmt.Println("  jsonl2parquet -input data.jsonl.gz -output data.parquet")
		os.Exit(1)
	}

	// Read first record to infer schema
	var firstRecord map[string]interface{}
	reader := fileiterator.FUOpen(*inputFile)
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Error: Empty JSONL file\n")
		os.Exit(1)
	}

	if err := json.Unmarshal(scanner.Bytes(), &firstRecord); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing first JSONL record: %v\n", err)
		os.Exit(1)
	}

	// Infer Arrow schema from first record
	fields := make([]arrow.Field, 0, len(firstRecord))
	fieldNames := make([]string, 0, len(firstRecord))

	for key, value := range firstRecord {
		fieldNames = append(fieldNames, key)
		var arrowType arrow.DataType

		switch v := value.(type) {
		case float64:
			// JSON numbers are always float64
			if v == float64(int64(v)) {
				arrowType = arrow.PrimitiveTypes.Int64
			} else {
				arrowType = arrow.PrimitiveTypes.Float64
			}
		case bool:
			arrowType = arrow.FixedWidthTypes.Boolean
		case string:
			arrowType = arrow.BinaryTypes.String
		default:
			arrowType = arrow.BinaryTypes.String // Default to string
		}

		fields = append(fields, arrow.Field{Name: key, Type: arrowType})
	}

	schema := arrow.NewSchema(fields, nil)

	// Create output file
	f, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Create Parquet writer with Snappy compression
	writerProps := parquet.NewWriterProperties(parquet.WithCompression(compress.Codecs.Snappy))
	arrowProps := pqarrow.DefaultWriterProps()
	writer, err := pqarrow.NewFileWriter(schema, f, writerProps, arrowProps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Parquet writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	// Build Arrow record batch
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	// Process all records
	count := 0

	// Process first record
	appendRecord(builder, firstRecord, fieldNames)
	count++

	// Process remaining records
	reader2 := fileiterator.FUOpen(*inputFile)
	defer reader2.Close()
	scanner2 := bufio.NewScanner(reader2)
	scanner2.Scan() // Skip first record (already processed)

	for scanner2.Scan() {
		var record map[string]interface{}
		if err := json.Unmarshal(scanner2.Bytes(), &record); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Skipping invalid JSON at line %d: %v\n", count+1, err)
			continue
		}

		appendRecord(builder, record, fieldNames)
		count++

		if count%100000 == 0 {
			fmt.Printf("Processed %d records...\n", count)
		}
	}

	// Write record batch
	rec := builder.NewRecord()
	defer rec.Release()

	if err := writer.Write(rec); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing Parquet data: %v\n", err)
		os.Exit(1)
	}

	writer.Close()
	f.Close()

	// Show file sizes
	inputStat, _ := os.Stat(*inputFile)
	outputStat, _ := os.Stat(*outputFile)

	fmt.Printf("\nConversion complete!\n")
	fmt.Printf("  Records: %d\n", count)
	fmt.Printf("  Input:   %.2f MB (%s)\n", float64(inputStat.Size())/1024/1024, *inputFile)
	fmt.Printf("  Output:  %.2f MB (%s)\n", float64(outputStat.Size())/1024/1024, *outputFile)
	fmt.Printf("  Ratio:   %.1fx\n", float64(inputStat.Size())/float64(outputStat.Size()))
}

func appendRecord(builder *array.RecordBuilder, record map[string]interface{}, fieldNames []string) {
	for i, fieldName := range fieldNames {
		value, exists := record[fieldName]

		switch b := builder.Field(i).(type) {
		case *array.Int64Builder:
			if !exists || value == nil {
				b.AppendNull()
			} else if v, ok := value.(float64); ok {
				b.Append(int64(v))
			} else {
				b.AppendNull()
			}
		case *array.Float64Builder:
			if !exists || value == nil {
				b.AppendNull()
			} else if v, ok := value.(float64); ok {
				b.Append(v)
			} else {
				b.AppendNull()
			}
		case *array.BooleanBuilder:
			if !exists || value == nil {
				b.AppendNull()
			} else if v, ok := value.(bool); ok {
				b.Append(v)
			} else {
				b.AppendNull()
			}
		case *array.StringBuilder:
			if !exists || value == nil {
				b.AppendNull()
			} else {
				b.Append(fmt.Sprintf("%v", value))
			}
		default:
			fmt.Fprintf(os.Stderr, "Warning: Unknown field type %v\n", reflect.TypeOf(b))
		}
	}
}
