package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type TestResultEntry struct {
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Action  string  `json:"Action"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
	Time    string  `json:"Time"`
}

func main() {
	// Parse the test results from stdin
	testResults := make([]TestResultEntry, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var entry TestResultEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip malformed JSON entries
		}
		testResults = append(testResults, entry)
	}

	csvWriter := csv.NewWriter(os.Stdout)

	// Write CSV header
	header := []string{"PACKAGE", "TEST", "ACTION", "ELAPSED"}
	if err := csvWriter.Write(header); err != nil {
		fmt.Printf("failed to write CSV header: %v", err)
		os.Exit(1)
	}

	// Write test results
	for _, result := range testResults {
		if result.Action == "output" {
			log.Print(result.Output)
		}
		if result.Test != "" && (result.Action == "pass" || result.Action == "fail") {
			record := []string{
				result.Package,
				result.Test,
				result.Action,
				fmt.Sprintf("%f", result.Elapsed),
			}
			if err := csvWriter.Write(record); err != nil {
				fmt.Printf("failed to write record to CSV: %v", err)
				os.Exit(1)
			}
		}
	}

	// Flush the writer (write whole output to stdout) and check for errors
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		fmt.Printf("error flushing the CSV writer: %v", err)
		os.Exit(1)
	}
}
