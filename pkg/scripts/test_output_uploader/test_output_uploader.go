package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var testResultsTableId = sdk.NewSchemaObjectIdentifier("TEST_RESULTS_DATABASE", "TEST_RESULTS_SCHEMA", "TEST_RESULTS_TABLE")

func main() {
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
			WithTemporary(true).
			WithFileFormat(*sdk.NewStageFileFormatRequest().WithFileFormatType(sdk.FileFormatTypeCSV)),
	); err != nil {
		log.Fatal("Failed to create test results stage:", err)
	}

	temporaryFile, err := os.CreateTemp("", "test_results_*.csv")
	if err != nil {
		log.Fatal("Failed to create temporary file:", err)
	}
	defer os.Remove(temporaryFile.Name())

	if _, err := io.Copy(temporaryFile, os.Stdin); err != nil {
		log.Fatal("Failed to write test results to temporary file:", err)
	}

	temporaryFileName := filepath.Base(temporaryFile.Name())
	temporaryFilePath, err := filepath.Abs(temporaryFile.Name())
	if err != nil {
		log.Fatal("Failed to get absolute path of temporary file:", err)
	}

	if err := temporaryFile.Close(); err != nil {
		log.Fatal("Failed to close temporary file:", err)
	}

	// Put file to the stage
	if _, err := client.ExecUnsafe(context.Background(), fmt.Sprintf("put file://%s @%s auto_compress = true overwrite = true", temporaryFilePath, testResultsStageId.FullyQualifiedName())); err != nil {
		log.Fatal("failed to put test results file to stage:", err)
	}

	// Copy data from the stage to the target table
	if _, err := client.ExecUnsafe(
		context.Background(),
		fmt.Sprintf(
			`copy into %s from @%s/%s file_format = (type = csv parse_header = true field_optionally_enclosed_by = '"') match_by_column_name = case_sensitive`,
			testResultsTableId.FullyQualifiedName(),
			testResultsStageId.FullyQualifiedName(),
			temporaryFileName,
		),
	); err != nil {
		log.Fatal("failed to copy test results file from stage to the target table:", err)
	}

	log.Println("Successfully processed test results")
}
