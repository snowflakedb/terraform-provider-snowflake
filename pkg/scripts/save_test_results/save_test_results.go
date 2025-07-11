package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"log"
	"os"
	"time"
)

type TestType string

const (
	TestTypeUnit        TestType = "unit"
	TestTypeIntegration TestType = "integration"
)

var (
	unitTestFileName        = "test_unit_output.json"
	integrationTestFileName = "test_integration_output.json"

	testResultsStageId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_STAGE")
	testResultsTableId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_TABLE")
)

func main() {
	testObjectSuffix, ok := os.LookupEnv(string(testenvs.TestObjectsSuffix))
	if !ok {
		log.Fatal("Environment variable TEST_SF_TF_TEST_OBJECT_SUFFIX is not set")
	}

	dirName, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	testResultsDirName := dirName + "/test_results"

	client, err := sdk.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create SDK client:", err)
	}

	testRunId := fmt.Sprintf("%s-%s", testObjectSuffix, time.Now().Format(time.RFC3339))

	if errs := errors.Join(
		processTestResults(TestTypeUnit, testRunId, client, testResultsStageId, testResultsTableId, testResultsDirName, unitTestFileName),
		processTestResults(TestTypeIntegration, testRunId, client, testResultsStageId, testResultsTableId, testResultsDirName, integrationTestFileName),
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

func processTestResults(testType TestType, testRunId string, client *sdk.Client, testResultsStageId sdk.SchemaObjectIdentifier, testResultsTableId sdk.SchemaObjectIdentifier, testResultsDirName string, fileName string) error {
	testResultsFilePath := testResultsDirName + "/" + fileName

	// TODO: This shouldn't be overwritten, make sure that each test run has it's own id
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true;", testResultsFilePath, testResultsStageId.FullyQualifiedName())); err != nil {
		return fmt.Errorf("failed to put test results file to stage, err = %w", err)
	}

	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf(`
copy into %s(test_run_id, test_type, package_name, test_name, action, elapsed)
from (
    select '%s', '%s', $1:Package::text, $1:Test::text, $1:Action::text, $1:Elapsed::float
	from @%s/%s
)
on_error = 'continue';
`, testResultsTableId.FullyQualifiedName(), testRunId, testType, testResultsStageId.FullyQualifiedName(), fileName)); err != nil {
		return fmt.Errorf("failed to put test results file to stage, err = %w", err)
	}

	return nil
}
