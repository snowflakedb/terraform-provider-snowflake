package sdk_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func SendTestResults(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakeProdEnvironment {
		t.Skip("Skipping for non-prod Snowflake environments")
	}

	testResultsTableId := sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_TABLE")

	client := fourthTestClient(t)

	// Create a temporary stage for test results
	testResultsStageId := sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_STAGE")
	if err := client.Stages.CreateInternal(
		context.Background(),
		sdk.NewCreateInternalStageRequest(testResultsStageId).
			WithTemporary(true).
			WithFileFormat(
				*sdk.NewStageFileFormatRequest().WithFileFormatOptions(sdk.FileFormatOptions{
					CsvOptions: &sdk.FileFormatCsvOptions{},
				}),
			),
	); err != nil {
		t.Fatal("Failed to create test results stage:", err)
	}

	files, err := filepath.Glob("/tmp/test_results_*.csv")
	if err != nil {
		t.Fatal("Failed to find test results files:", err)
	}
	t.Logf("Found %v test results files", files)

	// Range over temporary files with test results and copy the data into target tables
	for _, file := range files {
		// Put file to the stage
		if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true", file, testResultsStageId.FullyQualifiedName())); err != nil {
			t.Fatal("failed to put test results file to stage:", err)
		}

		fileBase := filepath.Base(file)
		// Copy data from the stage to the target table
		if _, err := client.ExecUnsafe(
			context.Background(),
			fmt.Sprintf(
				`copy into %s from @%s/%s file_format = (type = csv parse_header = true field_optionally_enclosed_by = '"') match_by_column_name = case_sensitive`,
				testResultsTableId.FullyQualifiedName(),
				testResultsStageId.FullyQualifiedName(),
				fileBase,
			),
		); err != nil {
			t.Fatal("failed to copy test results file from stage to the target table:", err)
		}
	}

	t.Log("Successfully processed test results")
}
