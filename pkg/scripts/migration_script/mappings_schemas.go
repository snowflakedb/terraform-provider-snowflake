package main

import (
	"fmt"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleSchemas(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[SchemaCsvRow, SchemaRepresentation](config, csvInput, MapSchemaToModel)
}

func MapSchemaToModel(schema SchemaRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	// Create the schema identifier
	schemaId := sdk.NewDatabaseObjectIdentifier(schema.DatabaseName, schema.Name)

	// Create a normalized resource ID
	resourceId := NormalizeResourceId(fmt.Sprintf("schema_%s", schemaId.FullyQualifiedName()))

	// Create the schema model
	resourceModel := model.Schema(resourceId, schema.DatabaseName, schema.Name)

	// Add optional fields
	if schema.Comment != "" {
		resourceModel.WithComment(schema.Comment)
	}

	if schema.IsTransient() {
		resourceModel.WithIsTransient("true")
	}

	if schema.IsManagedAccess() {
		resourceModel.WithWithManagedAccess("true")
	}

	handleOptionalFieldWithBuilder(schema.DataRetentionTimeInDays, resourceModel.WithDataRetentionTimeInDays)
	handleOptionalFieldWithBuilder(schema.MaxDataExtensionTimeInDays, resourceModel.WithMaxDataExtensionTimeInDays)
	handleOptionalFieldWithBuilder(schema.ExternalVolume, resourceModel.WithExternalVolume)
	handleOptionalFieldWithBuilder(schema.Catalog, resourceModel.WithCatalog)
	handleOptionalFieldWithBuilder(schema.PipeExecutionPaused, resourceModel.WithPipeExecutionPaused)
	handleOptionalFieldWithBuilder(schema.ReplaceInvalidCharacters, resourceModel.WithReplaceInvalidCharacters)
	handleOptionalFieldWithBuilder(schema.DefaultDdlCollation, resourceModel.WithDefaultDdlCollation)
	handleOptionalFieldWithBuilder(schema.StorageSerializationPolicy, resourceModel.WithStorageSerializationPolicy)
	handleOptionalFieldWithBuilder(schema.LogLevel, resourceModel.WithLogLevel)
	handleOptionalFieldWithBuilder(schema.TraceLevel, resourceModel.WithTraceLevel)
	handleOptionalFieldWithBuilder(schema.SuspendTaskAfterNumFailures, resourceModel.WithSuspendTaskAfterNumFailures)
	handleOptionalFieldWithBuilder(schema.TaskAutoRetryAttempts, resourceModel.WithTaskAutoRetryAttempts)
	handleOptionalFieldWithBuilder(schema.UserTaskManagedInitialWarehouseSize, resourceModel.WithUserTaskManagedInitialWarehouseSize)
	handleOptionalFieldWithBuilder(schema.UserTaskTimeoutMs, resourceModel.WithUserTaskTimeoutMs)
	handleOptionalFieldWithBuilder(schema.UserTaskMinimumTriggerIntervalInSeconds, resourceModel.WithUserTaskMinimumTriggerIntervalInSeconds)
	handleOptionalFieldWithBuilder(schema.QuotedIdentifiersIgnoreCase, resourceModel.WithQuotedIdentifiersIgnoreCase)
	handleOptionalFieldWithBuilder(schema.EnableConsoleOutput, resourceModel.WithEnableConsoleOutput)

	// Create the import model with the schema identifier
	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		schemaId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
