package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/parf/homebase-go-lib/fileiterator"
)

func main() {
	inputFile := flag.String("input", "", "Input CSV file (supports .gz, .zst, .lz4, .br, .xz)")
	outputFile := flag.String("output", "", "Output Parquet file")
	hasHeader := flag.Bool("header", true, "CSV has header row")
	delimiter := flag.String("delimiter", ",", "CSV delimiter (comma, tab, pipe, etc.)")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: csv2parquet -input <csv_file> -output <parquet_file> [-header=true] [-delimiter=,]")
		fmt.Println("\nExample:")
		fmt.Println("  csv2parquet -input data.csv.gz -output data.parquet")
		fmt.Println("  csv2parquet -input data.tsv -output data.parquet -delimiter=tab")
		os.Exit(1)
	}

	// Parse delimiter
	var delimRune rune
	switch *delimiter {
	case "comma", ",":
		delimRune = ','
	case "tab", "\\t":
		delimRune = '\t'
	case "pipe", "|":
		delimRune = '|'
	case "semicolon", ";":
		delimRune = ';'
	default:
		if len(*delimiter) == 1 {
			delimRune = rune((*delimiter)[0])
		} else {
			fmt.Fprintf(os.Stderr, "Error: Invalid delimiter '%s'\n", *delimiter)
			os.Exit(1)
		}
	}

	// Open CSV file
	reader := fileiterator.FUOpen(*inputFile)
	defer reader.Close()

	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimRune
	csvReader.TrimLeadingSpace = true

	// Read header or infer column names
	var headers []string
	firstRow, err := csvReader.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	if *hasHeader {
		headers = firstRow
	} else {
		// Generate column names: col0, col1, col2, ...
		headers = make([]string, len(firstRow))
		for i := range headers {
			headers[i] = fmt.Sprintf("col%d", i)
		}
	}

	// Infer schema by reading first few rows
	var sampleRows [][]string
	if !*hasHeader {
		sampleRows = append(sampleRows, firstRow)
	}

	for i := 0; i < 100; i++ { // Sample first 100 rows
		row, err := csvReader.Read()
		if err != nil {
			break
		}
		sampleRows = append(sampleRows, row)
	}

	// Infer types for each column
	columnTypes := make([]arrow.DataType, len(headers))
	for i := range headers {
		columnTypes[i] = inferColumnType(sampleRows, i)
	}

	// Create Arrow schema
	fields := make([]arrow.Field, len(headers))
	for i, header := range headers {
		fields[i] = arrow.Field{Name: header, Type: columnTypes[i]}
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

	// Process sampled rows
	for _, row := range sampleRows {
		appendCSVRow(builder, row, columnTypes)
	}
	count := len(sampleRows)

	// Process remaining rows
	reader2 := fileiterator.FUOpen(*inputFile)
	defer reader2.Close()
	csvReader2 := csv.NewReader(reader2)
	csvReader2.Comma = delimRune
	csvReader2.TrimLeadingSpace = true

	if *hasHeader {
		csvReader2.Read() // Skip header
	}

	// Skip already processed rows
	for i := 0; i < len(sampleRows); i++ {
		csvReader2.Read()
	}

	for {
		row, err := csvReader2.Read()
		if err != nil {
			break
		}

		appendCSVRow(builder, row, columnTypes)
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
	fmt.Printf("  Columns: %d\n", len(headers))
	fmt.Printf("  Input:   %.2f MB (%s)\n", float64(inputStat.Size())/1024/1024, *inputFile)
	fmt.Printf("  Output:  %.2f MB (%s)\n", float64(outputStat.Size())/1024/1024, *outputFile)
	fmt.Printf("  Ratio:   %.1fx\n", float64(inputStat.Size())/float64(outputStat.Size()))
}

func inferColumnType(rows [][]string, colIndex int) arrow.DataType {
	hasFloat := false
	hasInt := false
	hasBool := false

	for _, row := range rows {
		if colIndex >= len(row) || row[colIndex] == "" {
			continue
		}

		value := row[colIndex]

		// Check if boolean
		if value == "true" || value == "false" || value == "TRUE" || value == "FALSE" {
			hasBool = true
			continue
		}

		// Check if integer
		if _, err := strconv.ParseInt(value, 10, 64); err == nil {
			hasInt = true
			continue
		}

		// Check if float
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			hasFloat = true
			continue
		}

		// If none of above, it's a string
		return arrow.BinaryTypes.String
	}

	if hasBool && !hasInt && !hasFloat {
		return arrow.FixedWidthTypes.Boolean
	}
	if hasFloat {
		return arrow.PrimitiveTypes.Float64
	}
	if hasInt {
		return arrow.PrimitiveTypes.Int64
	}

	return arrow.BinaryTypes.String
}

func appendCSVRow(builder *array.RecordBuilder, row []string, columnTypes []arrow.DataType) {
	for i, value := range row {
		if i >= len(columnTypes) {
			break
		}

		switch b := builder.Field(i).(type) {
		case *array.Int64Builder:
			if value == "" {
				b.AppendNull()
			} else if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				b.Append(v)
			} else {
				b.AppendNull()
			}
		case *array.Float64Builder:
			if value == "" {
				b.AppendNull()
			} else if v, err := strconv.ParseFloat(value, 64); err == nil {
				b.Append(v)
			} else {
				b.AppendNull()
			}
		case *array.BooleanBuilder:
			if value == "" {
				b.AppendNull()
			} else if value == "true" || value == "TRUE" || value == "1" {
				b.Append(true)
			} else if value == "false" || value == "FALSE" || value == "0" {
				b.Append(false)
			} else {
				b.AppendNull()
			}
		case *array.StringBuilder:
			if value == "" {
				b.AppendNull()
			} else {
				b.Append(value)
			}
		}
	}
}
