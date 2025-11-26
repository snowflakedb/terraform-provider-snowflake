package main

import (
	"fmt"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleDatabases(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[DatabaseCsvRow, DatabaseRepresentation](config, csvInput, MapDatabaseToModel)
}

func MapDatabaseToModel(database DatabaseRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	databaseId := sdk.NewAccountObjectIdentifier(database.Name)
	resourceId := NormalizeResourceId(fmt.Sprintf("database_%s", databaseId.FullyQualifiedName()))
	resourceModel := model.Database(resourceId, database.Name)

	handleIfNotEmpty(database.Comment, resourceModel.WithComment)
	if database.Transient {
		resourceModel.WithIsTransient(true)
	}

	handleOptionalFieldWithBuilder(database.DataRetentionTimeInDays, resourceModel.WithDataRetentionTimeInDays)
	handleOptionalFieldWithBuilder(database.MaxDataExtensionTimeInDays, resourceModel.WithMaxDataExtensionTimeInDays)
	handleOptionalFieldWithBuilder(database.ExternalVolume, resourceModel.WithExternalVolume)
	handleOptionalFieldWithBuilder(database.Catalog, resourceModel.WithCatalog)
	handleOptionalFieldWithBuilder(database.ReplaceInvalidCharacters, resourceModel.WithReplaceInvalidCharacters)
	handleOptionalFieldWithBuilder(database.DefaultDDLCollation, resourceModel.WithDefaultDdlCollation)
	handleOptionalFieldWithBuilder(database.StorageSerializationPolicy, resourceModel.WithStorageSerializationPolicy)
	handleOptionalFieldWithBuilder(database.LogLevel, resourceModel.WithLogLevel)
	handleOptionalFieldWithBuilder(database.TraceLevel, resourceModel.WithTraceLevel)
	handleOptionalFieldWithBuilder(database.SuspendTaskAfterNumFailures, resourceModel.WithSuspendTaskAfterNumFailures)
	handleOptionalFieldWithBuilder(database.TaskAutoRetryAttempts, resourceModel.WithTaskAutoRetryAttempts)
	handleOptionalFieldWithBuilder(database.UserTaskManagedInitialWarehouseSize, resourceModel.WithUserTaskManagedInitialWarehouseSize)
	handleOptionalFieldWithBuilder(database.UserTaskTimeoutMs, resourceModel.WithUserTaskTimeoutMs)
	handleOptionalFieldWithBuilder(database.UserTaskMinimumTriggerIntervalInSeconds, resourceModel.WithUserTaskMinimumTriggerIntervalInSeconds)
	handleOptionalFieldWithBuilder(database.QuotedIdentifiersIgnoreCase, resourceModel.WithQuotedIdentifiersIgnoreCase)
	handleOptionalFieldWithBuilder(database.EnableConsoleOutput, resourceModel.WithEnableConsoleOutput)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		databaseId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
