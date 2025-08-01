package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	testTypesInWorkflow, ok := os.LookupEnv("TEST_SF_TF_TEST_TYPES_IN_WORKFLOW")
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_TYPES_IN_WORKFLOW is not set")
	}

	testTypesInWorkflowMapped, err := collections.MapErr(strings.Split(testTypesInWorkflow, ","), ToTestType)
	if err != nil {
		log.Fatal("Failed to parse TEST_SF_TF_TEST_TYPES_IN_WORKFLOW:", err)
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

	if errs := errors.Join(collections.Map(testTypesInWorkflowMapped, func(testType TestType) error {
		log.Printf("Processing test results from  %s test type", testType)
		return processTestResults(testType, testWorkflowId, client, testResultsStageId, testResultsDirName)
	})...); errs != nil {
		log.Fatal(errs)
	}

	log.Println("Successfully processed test results")
}

func processTestResults(testType TestType, testWorkflowId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsDirName string) error {
	err, fileLocation := stageTestResults(testType, testWorkflowId, client, testResultsStageId, testResultsDirName)
	if err != nil {
		return fmt.Errorf("failed to stage test results for test type %s, err = %w", testType, err)
	}

	// Create a temporary table for storing and transforming test results before inserting them into the final table
	temporaryTableId := sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", fmt.Sprintf("TEST_RESULTS_TEMP_%s", testType))
	if err := client.Tables.Create(context.Background(), sdk.NewCreateTableRequest(temporaryTableId, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("test_run_id", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("test_type", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("package_name", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("test_name", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("action", sdk.DataTypeVARCHAR),
		*sdk.NewTableColumnRequest("elapsed", sdk.DataTypeFloat),
		*sdk.NewTableColumnRequest("finished_at", sdk.DataTypeTimestampNTZ),
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

func stageTestResults(testType TestType, testWorkflowId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsDirName string) (error, sdk.Location) {
	fileName := fmt.Sprintf("test_%s_output.json", testType)
	testResultsFilePath := testResultsDirName + "/" + fileName

	uniqueFileName := fmt.Sprintf("%s_test_%s_output.json", testWorkflowId, testType)
	uniqueTestResultsFilePath := testResultsDirName + "/" + uniqueFileName

	fileLocation := sdk.NewStageLocation(testResultsStageId, uniqueFileName)

	// We have to rename them because it's not possible to pass different target file name in Snowflake,
	// and we need to have unique file names for each test run (to avoid collisions with other test runs).
	if err := os.Rename(testResultsFilePath, uniqueTestResultsFilePath); err != nil {
		return fmt.Errorf("failed to rename test results file %s, err = %w", testResultsFilePath, err), fileLocation
	}

	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true;", uniqueTestResultsFilePath, testResultsStageId.FullyQualifiedName())); err != nil {
		return fmt.Errorf("failed to put test results file to stage, err = %w", err), fileLocation
	}

	return nil, fileLocation
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
