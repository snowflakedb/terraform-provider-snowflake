package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

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

var (
	testResultsStageId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_STAGE")
	testResultsTableId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_TABLE")
)

func main() {
	testRunId, ok := os.LookupEnv("TEST_SF_TF_TEST_WORKFLOW_ID")
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_WORKFLOW_ID is not set")
	}

	dirName, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	testResultsDirName := dirName + "/test_results"

	clientConfig, err := sdk.ProfileConfig("fourth_test_account")
	if err != nil {
		log.Fatal("Failed to get client config:", err)
	}

	client, err := sdk.NewClient(clientConfig)
	if err != nil {
		log.Fatal("Failed to create a new client:", err)
	}

	if errs := errors.Join(
		processTestResults(TestTypeUnit, testRunId, client, testResultsStageId, testResultsTableId, testResultsDirName),
		//processTestResults(TestTypeIntegration, testRunId, client, testResultsStageId, testResultsTableId, testResultsDirName),
	); errs != nil {
		log.Fatal(errs)
	}

	// Remove all test results that are not valid
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`
delete from %s where
    test_name is null or
    not (action = 'pass' or action = 'fail');
`, testResultsTableId.FullyQualifiedName())); err != nil {
		log.Fatal("failed to put test results file to stage:", err)
	}

	log.Println("Successfully processed test results")
}

func processTestResults(testType TestType, testRunId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsTableId sdk.SchemaObjectIdentifier, testResultsDirName string) error {
	fileName := fmt.Sprintf("test_%s_output.json", testType)
	testResultsFilePath := testResultsDirName + "/" + fileName

	uniqueFileName := fmt.Sprintf("%s_test_%s_output.json", testRunId, testType)
	uniqueTestResultsFilePath := testResultsDirName + "/" + uniqueFileName

	// We have to rename them because it's not possible to pass different target file name in Snowflake,
	// and we need to have unique file names for each test run (to avoid collisions with other test runs).
	if err := os.Rename(testResultsFilePath, uniqueTestResultsFilePath); err != nil {
		return fmt.Errorf("failed to rename test results file, err = %w", err)
	}

	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true;", uniqueTestResultsFilePath, testResultsStageId.FullyQualifiedName())); err != nil {
		return fmt.Errorf("failed to put test results file to stage, err = %w", err)
	}
	defer func() {
		// Clean up the staged file after processing
		if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("delete @%s/%s;", testResultsStageId.FullyQualifiedName(), uniqueFileName)); err != nil {
			log.Printf("failed to remove test results from stage, err = %v", err)
		}
	}()

	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`
copy into %s(test_run_id, test_type, package_name, test_name, action, elapsed)
from (
    select '%s', '%s', $1:Package::text, $1:Test::text, $1:Action::text, $1:Elapsed::float
	from @%s/%s
)
on_error = 'continue';
`, testResultsTableId.FullyQualifiedName(), testRunId, testType, testResultsStageId.FullyQualifiedName(), uniqueFileName)); err != nil {
		return fmt.Errorf("failed to put test results file to stage, err = %w", err)
	}

	return nil
}
