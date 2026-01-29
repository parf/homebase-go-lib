package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v7"
)

type Record struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Age       int     `json:"age"`
	Score     float64 `json:"score"`
	Active    bool    `json:"active"`
	Category  string  `json:"category"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	count := flag.Int("count", 100, "Number of records to generate")
	seed := flag.Int64("seed", 12345, "Random seed for reproducibility")
	flag.Parse()

	gofakeit.Seed(*seed)

	records := make([]Record, *count)
	categories := []string{"Electronics", "Books", "Clothing", "Food", "Sports", "Toys", "Home", "Garden", "Health", "Beauty"}

	for i := 0; i < *count; i++ {
		records[i] = Record{
			ID:        int64(i),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Age:       gofakeit.Number(18, 80),
			Score:     gofakeit.Float64Range(0, 100),
			Active:    gofakeit.Bool(),
			Category:  gofakeit.RandomString(categories),
			Timestamp: gofakeit.Date().Unix(),
		}
	}

	// Write JSONL
	jsonlFile, _ := os.Create("sample-data.jsonl")
	defer jsonlFile.Close()
	jsonEncoder := json.NewEncoder(jsonlFile)
	for _, rec := range records {
		jsonEncoder.Encode(rec)
	}
	fmt.Printf("Generated sample-data.jsonl (%d records)\n", *count)

	// Write CSV
	csvFile, _ := os.Create("sample-data.csv")
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write([]string{"id", "name", "email", "age", "score", "active", "category", "timestamp"})
	for _, rec := range records {
		csvWriter.Write([]string{
			fmt.Sprintf("%d", rec.ID),
			rec.Name,
			rec.Email,
			fmt.Sprintf("%d", rec.Age),
			fmt.Sprintf("%.2f", rec.Score),
			fmt.Sprintf("%t", rec.Active),
			rec.Category,
			fmt.Sprintf("%d", rec.Timestamp),
		})
	}
	csvWriter.Flush()
	fmt.Printf("Generated sample-data.csv (%d records)\n", *count)
}
