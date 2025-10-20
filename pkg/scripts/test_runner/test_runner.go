package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"sync"
)

type TestType string

const (
	TestTypeUnit                  TestType = "unit"
	TestTypeIntegration           TestType = "integration"
	TestTypeAcceptance            TestType = "acceptance"
	TestTypeAccountLevel          TestType = "account_level"
	TestTypeFunctional            TestType = "functional"
	TestTypeArchitecture          TestType = "architecture"
	TestTypeMainTerraformUseCases TestType = "main_terraform_use_cases"
)

var AllTestTypes = []TestType{
	TestTypeUnit,
	TestTypeIntegration,
	TestTypeAcceptance,
	TestTypeAccountLevel,
	TestTypeFunctional,
	TestTypeArchitecture,
	TestTypeMainTerraformUseCases,
}

var TestMakeCommand = map[TestType]string{
	TestTypeUnit:                  "test-unit",
	TestTypeIntegration:           "test-integration",
	TestTypeAcceptance:            "test-acceptance",
	TestTypeAccountLevel:          "test-account-level-features",
	TestTypeFunctional:            "test-functional",
	TestTypeArchitecture:          "test-architecture",
	TestTypeMainTerraformUseCases: "test-main-terraform-use-cases",
}

func ToTestType(s string) (TestType, error) {
	if slices.Contains(AllTestTypes, TestType(s)) {
		return TestType(s), nil
	}
	return "", fmt.Errorf("unknown test type: %s", s)
}

func main() {
	testWorkflowId, ok := os.LookupEnv("TEST_SF_TF_TEST_WORKFLOW_ID")
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_WORKFLOW_ID is not set")
	}
	log.Println("Processing with the following workflow id: ", testWorkflowId)

	testType, ok := os.LookupEnv("TEST_SF_TF_TEST_TYPE")
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_TYPE is not set")
	}

	outputPath, ok := os.LookupEnv("TEST_SF_TF_TEST_RESULTS_OUTPUT_PATH")
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_TYPE is not set")
	}

	mappedTestType, err := ToTestType(testType)
	if err != nil {
		log.Fatal("Failed to parse TEST_SF_TF_TEST_TYPE:", err)
	}

	testResults := runTest(mappedTestType)
	if err := processAndSaveTestResults(testWorkflowId, testResults, outputPath); err != nil {
		log.Fatalf("Failed to processed the test results, err = %s", err)
	}

	log.Println("Successfully processed test results")
}

type TestResultEntry struct {
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Action  string  `json:"Action"`
	Elapsed float64 `json:"Elapsed"`
	Time    string  `json:"Time"`
}

func processAndSaveTestResults(id string, results *bytes.Buffer, path string) error {
	testResults := make([]TestResultEntry, 0)
	scanner := bufio.NewScanner(results)
	for scanner.Scan() {
		var entry TestResultEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal test result entry: %w", err)
		}
		testResults = append(testResults, entry)
	}

	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Write CSV header
	header := []string{"TEST_RUN_ID", "PACKAGE", "TEST", "ACTION", "ELAPSED"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write test results
	for _, result := range testResults {
		if result.Test != "" && (result.Action == "pass" || result.Action == "fail") {
			record := []string{
				id,
				result.Package,
				result.Test,
				result.Action,
				fmt.Sprintf("%f", result.Elapsed),
			}
			if err := csvWriter.Write(record); err != nil {
				return fmt.Errorf("failed to write record to CSV: %w", err)
			}
		}
	}

	return nil
}

func runTest(testType TestType) *bytes.Buffer {
	cmd := exec.Command("make", TestMakeCommand[testType])

	buf := NewSyncBuffer()

	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Env = cmd.Environ()

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		if err := cmd.Run(); err != nil {
			log.Printf("Running %s tests run ended with errors: %s", testType, err)
		}
	}()

	resultingBuffer := new(bytes.Buffer)
	reader := bufio.NewReader(buf)

loop:
	for {
		select {
		case <-doneChan:
			for reader.Buffered() != 0 {
				text, _ := reader.ReadBytes('\n')

				var entry map[string]string
				_ = json.Unmarshal(text, &entry)

				if v, ok := entry["Output"]; ok && v != "" {
					fmt.Print(v)
				}

				if v, ok := entry["Action"]; ok && v != "" {
					resultingBuffer.Write(text)
				}
			}
			break loop
		default:
			text, _ := reader.ReadBytes('\n')

			var entry map[string]string
			_ = json.Unmarshal(text, &entry)

			if v, ok := entry["Output"]; ok && v != "" {
				fmt.Print(v)
			}

			if v, ok := entry["Action"]; ok && v != "" {
				resultingBuffer.Write(text)
			}
		}
	}

	return resultingBuffer
}

type SyncBuffer struct {
	buf    *bytes.Buffer
	reader *bufio.Reader
	mutex  *sync.Mutex
}

func NewSyncBuffer() *SyncBuffer {
	buf := new(bytes.Buffer)
	return &SyncBuffer{
		buf:    buf,
		reader: bufio.NewReader(buf),
		mutex:  new(sync.Mutex),
	}
}

func (b *SyncBuffer) Read(p []byte) (n int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buf.Read(p)
}

func (b *SyncBuffer) Write(p []byte) (n int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buf.Write(p)
}
