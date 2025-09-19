package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"log"
	"os"
	"os/exec"
	"slices"
	"sync"
)

type TestType string

const (
	TestTypeUnit         TestType = "unit"
	TestTypeIntegration  TestType = "integration"
	TestTypeAcceptance   TestType = "acceptance"
	TestTypeAccountLevel TestType = "account_level"
	TestTypeFunctional   TestType = "functional"
	TestTypeArchitecture TestType = "architecture"
)

var AllTestTypes = []TestType{
	TestTypeUnit,
	TestTypeIntegration,
	TestTypeAcceptance,
	TestTypeAccountLevel,
	TestTypeFunctional,
	TestTypeArchitecture,
}

var TestMakeCommand = map[TestType]string{
	TestTypeUnit:         "test-unit",
	TestTypeIntegration:  "test-integration",
	TestTypeAcceptance:   "test-acceptance",
	TestTypeAccountLevel: "test-account-level-features",
	TestTypeFunctional:   "test-functional",
	TestTypeArchitecture: "test-architecture",
}

func ToTestType(s string) (TestType, error) {
	if slices.Contains(AllTestTypes, TestType(s)) {
		return TestType(s), nil
	}
	return "", fmt.Errorf("unknown test type: %s", s)
}

var testResultsTableId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_TABLE")

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

	mappedTestType, err := ToTestType(testType)
	if err != nil {
		log.Fatal("Failed to parse TEST_SF_TF_TEST_TYPE:", err)
	}

	dirName, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	testResultsDirName := dirName + "/test_results"

	clientConfig, err := sdk.ProfileConfig(testprofiles.Fourth)
	if err != nil {
		log.Fatal("Failed to get client config:", err)
	}

	client, err := sdk.NewClient(clientConfig)
	if err != nil {
		log.Fatal("Failed to create a new client:", err)
	}

	// Create a temporary stage for test results
	testResultsStageId := sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_STAGE")
	if err := client.Stages.CreateInternal(
		context.Background(),
		sdk.NewCreateInternalStageRequest(testResultsStageId).
			WithTemporary(sdk.Pointer(true)).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(sdk.Pointer(sdk.FileFormatTypeJSON))),
	); err != nil {
		log.Fatal("Failed to create test results stage:", err)
	}

	log.Printf("Running %s tests", testType)

	testResults := runTest(mappedTestType)
	if err := processTestResults(mappedTestType, testWorkflowId, client, testResultsStageId, testResultsDirName, testResults); err != nil {
		log.Fatalf("Failed to processed the test results, err = %s", err)
	}

	log.Println("Successfully processed test results")
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
				// TODO: Try to unify with code below
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

			// TODO: Remove if (workaround for now)
			if v, ok := entry["Action"]; ok && v != "" {
				resultingBuffer.Write(text)
			}
		}
	}

	return resultingBuffer
}

func processTestResults(testType TestType, testWorkflowId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsDirName string, testResults *bytes.Buffer) error {
	fileLocation, err := stageTestResults(testType, testWorkflowId, client, testResultsStageId, testResultsDirName, testResults)
	if err != nil {
		return fmt.Errorf("failed to stage test results for test type %s, err = %w", testType, err)
	}

	// Create a temporary table for storing and transforming test results before inserting them into the final table
	temporaryTableId := sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", fmt.Sprintf("TEST_RESULTS_TEMP_%s", testType))
	if err := client.Tables.Create(context.Background(), sdk.NewCreateTableRequest(temporaryTableId, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("TEST_RUN_ID", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("TEST_TYPE", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("PACKAGE_NAME", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("TEST_NAME", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("ACTION", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("ELAPSED", sdk.DataTypeFloat),
		*sdk.NewTableColumnRequest("FINISHED_AT", sdk.DataTypeTimestampNTZ),
	}).WithKind(sdk.Pointer(sdk.TemporaryTableKind))); err != nil {
		return fmt.Errorf("failed to create temporary table for test results, err = %w", err)
	}

	// Store the data in the temporary table
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`
	copy into %s(test_run_id, test_type, package_name, test_name, action, elapsed, finished_at)
	from (select '%s', '%s', $1:Package::text, $1:Test::text, $1:Action::text, $1:Elapsed::float, $1:Time::timestamp from %s)
	on_error = 'continue';
	`, temporaryTableId.FullyQualifiedName(), testWorkflowId, testType, fileLocation.ToSql())); err != nil {
		return fmt.Errorf("failed to copy test results file from stage to the target table, err = %w", err)
	}

	// Prepare the data before putting it into the final table.
	if err := transformTemporaryTableData(client, temporaryTableId); err != nil {
		return fmt.Errorf("failed to transform temporary table data, err = %w", err)
	}

	// Move the data from the temporary table to the final table
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`insert into %s select * from %s;`, testResultsTableId.FullyQualifiedName(), temporaryTableId.FullyQualifiedName())); err != nil {
		return fmt.Errorf("failed to copy test results from the temporary table to the target one, err = %w", err)
	}

	return nil
}

func stageTestResults(testType TestType, testWorkflowId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsDirName string, testResults *bytes.Buffer) (sdk.Location, error) {
	fileName := fmt.Sprintf("%s_test_%s_output.json", testWorkflowId, testType)
	testResultsFilePath := testResultsDirName + "/" + fileName

	if err := os.MkdirAll(testResultsDirName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create test results directory %s, err = %w", testResultsDirName, err)
	}

	file, err := os.Create(testResultsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create test results file %s, err = %w", testResultsFilePath, err)
	}
	defer file.Close()
	fileLocation := sdk.NewStageLocation(testResultsStageId, fileName)

	if _, err := file.Write(testResults.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write test results to file %s, err = %w", testResultsFilePath, err)
	}

	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true;", testResultsFilePath, testResultsStageId.FullyQualifiedName())); err != nil {
		return nil, fmt.Errorf("failed to put test results file to stage, err = %w", err)
	}

	return fileLocation, nil
}

func transformTemporaryTableData(client *sdk.Client, temporaryTableId sdk.SchemaObjectIdentifier) error {
	// Remove all test results that are not valid (for now, we're only interested in the 'pass' and 'fail' actions).
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`
	delete from %s where
		test_name is null
		or
		not (action = 'pass' or action = 'fail');
`, temporaryTableId.FullyQualifiedName())); err != nil {
		return fmt.Errorf("failed to delete invalid test results from temporary table, err = %w", err)
	}

	return nil
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
