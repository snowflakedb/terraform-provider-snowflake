package main

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[DatabaseRepresentation] = new(DatabaseCsvRow)

type DatabaseCsvRow struct {
	Comment                                      string `csv:"comment"`
	CreatedOn                                    string `csv:"created_on"`
	DroppedOn                                    string `csv:"dropped_on"`
	IsCurrent                                    string `csv:"is_current"`
	IsDefault                                    string `csv:"is_default"`
	Kind                                         string `csv:"kind"`
	Name                                         string `csv:"name"`
	Options                                      string `csv:"options"`
	Origin                                       string `csv:"origin"`
	Owner                                        string `csv:"owner"`
	OwnerRoleType                                string `csv:"owner_role_type"`
	ResourceGroup                                string `csv:"resource_group"`
	RetentionTime                                string `csv:"retention_time"`
	CatalogLevel                                 string `csv:"catalog_level"`
	CatalogValue                                 string `csv:"catalog_value"`
	DataRetentionTimeInDaysLevel                 string `csv:"data_retention_time_in_days_level"`
	DataRetentionTimeInDaysValue                 string `csv:"data_retention_time_in_days_value"`
	DefaultDDLCollationLevel                     string `csv:"default_ddl_collation_level"`
	DefaultDDLCollationValue                     string `csv:"default_ddl_collation_value"`
	EnableConsoleOutputLevel                     string `csv:"enable_console_output_level"`
	EnableConsoleOutputValue                     string `csv:"enable_console_output_value"`
	ExternalVolumeLevel                          string `csv:"external_volume_level"`
	ExternalVolumeValue                          string `csv:"external_volume_value"`
	LogLevelLevel                                string `csv:"log_level_level"`
	LogLevelValue                                string `csv:"log_level_value"`
	MaxDataExtensionTimeInDaysLevel              string `csv:"max_data_extension_time_in_days_level"`
	MaxDataExtensionTimeInDaysValue              string `csv:"max_data_extension_time_in_days_value"`
	QuotedIdentifiersIgnoreCaseLevel             string `csv:"quoted_identifiers_ignore_case_level"`
	QuotedIdentifiersIgnoreCaseValue             string `csv:"quoted_identifiers_ignore_case_value"`
	ReplaceInvalidCharactersLevel                string `csv:"replace_invalid_characters_level"`
	ReplaceInvalidCharactersValue                string `csv:"replace_invalid_characters_value"`
	StorageSerializationPolicyLevel              string `csv:"storage_serialization_policy_level"`
	StorageSerializationPolicyValue              string `csv:"storage_serialization_policy_value"`
	SuspendTaskAfterNumFailuresLevel             string `csv:"suspend_task_after_num_failures_level"`
	SuspendTaskAfterNumFailuresValue             string `csv:"suspend_task_after_num_failures_value"`
	TaskAutoRetryAttemptsLevel                   string `csv:"task_auto_retry_attempts_level"`
	TaskAutoRetryAttemptsValue                   string `csv:"task_auto_retry_attempts_value"`
	TraceLevelLevel                              string `csv:"trace_level_level"`
	TraceLevelValue                              string `csv:"trace_level_value"`
	UserTaskManagedInitialWarehouseSizeLevel     string `csv:"user_task_managed_initial_warehouse_size_level"`
	UserTaskManagedInitialWarehouseSizeValue     string `csv:"user_task_managed_initial_warehouse_size_value"`
	UserTaskMinimumTriggerIntervalInSecondsLevel string `csv:"user_task_minimum_trigger_interval_in_seconds_level"`
	UserTaskMinimumTriggerIntervalInSecondsValue string `csv:"user_task_minimum_trigger_interval_in_seconds_value"`
	UserTaskTimeoutMsLevel                       string `csv:"user_task_timeout_ms_level"`
	UserTaskTimeoutMsValue                       string `csv:"user_task_timeout_ms_value"`
}

type DatabaseRepresentation struct {
	sdk.Database

	// parameters
	Catalog                                 *string
	DataRetentionTimeInDays                 *int
	DefaultDDLCollation                     *string
	EnableConsoleOutput                     *bool
	ExternalVolume                          *string
	LogLevel                                *string
	MaxDataExtensionTimeInDays              *int
	QuotedIdentifiersIgnoreCase             *bool
	ReplaceInvalidCharacters                *bool
	StorageSerializationPolicy              *string
	SuspendTaskAfterNumFailures             *int
	TaskAutoRetryAttempts                   *int
	TraceLevel                              *string
	UserTaskManagedInitialWarehouseSize     *string
	UserTaskMinimumTriggerIntervalInSeconds *int
	UserTaskTimeoutMs                       *int
}

func (row DatabaseCsvRow) convert() (*DatabaseRepresentation, error) {
	databaseRepresentation := &DatabaseRepresentation{
		Database: sdk.Database{
			Name:          row.Name,
			IsCurrent:     row.IsCurrent == "Y",
			IsDefault:     row.IsDefault == "Y",
			Owner:         row.Owner,
			Comment:       row.Comment,
			Kind:          row.Kind,
			OwnerRoleType: row.OwnerRoleType,
			ResourceGroup: row.ResourceGroup,
		},
	}
	if row.Options != "" {
		databaseRepresentation.Options = row.Options
		databaseRepresentation.Database.SetTransient(row.Options)
	}

	handler := newParameterHandler(sdk.ParameterTypeDatabase)
	errs := errors.Join(
		handler.handleStringParameter(row.CatalogLevel, row.CatalogValue, &databaseRepresentation.Catalog),
		handler.handleIntegerParameter(row.DataRetentionTimeInDaysLevel, row.DataRetentionTimeInDaysValue, &databaseRepresentation.DataRetentionTimeInDays),
		handler.handleStringParameter(row.DefaultDDLCollationLevel, row.DefaultDDLCollationValue, &databaseRepresentation.DefaultDDLCollation),
		handler.handleBooleanParameter(row.EnableConsoleOutputLevel, row.EnableConsoleOutputValue, &databaseRepresentation.EnableConsoleOutput),
		handler.handleStringParameter(row.ExternalVolumeLevel, row.ExternalVolumeValue, &databaseRepresentation.ExternalVolume),
		handler.handleStringParameter(row.LogLevelLevel, row.LogLevelValue, &databaseRepresentation.LogLevel),
		handler.handleIntegerParameter(row.MaxDataExtensionTimeInDaysLevel, row.MaxDataExtensionTimeInDaysValue, &databaseRepresentation.MaxDataExtensionTimeInDays),
		handler.handleBooleanParameter(row.QuotedIdentifiersIgnoreCaseLevel, row.QuotedIdentifiersIgnoreCaseValue, &databaseRepresentation.QuotedIdentifiersIgnoreCase),
		handler.handleBooleanParameter(row.ReplaceInvalidCharactersLevel, row.ReplaceInvalidCharactersValue, &databaseRepresentation.ReplaceInvalidCharacters),
		handler.handleStringParameter(row.StorageSerializationPolicyLevel, row.StorageSerializationPolicyValue, &databaseRepresentation.StorageSerializationPolicy),
		handler.handleIntegerParameter(row.SuspendTaskAfterNumFailuresLevel, row.SuspendTaskAfterNumFailuresValue, &databaseRepresentation.SuspendTaskAfterNumFailures),
		handler.handleIntegerParameter(row.TaskAutoRetryAttemptsLevel, row.TaskAutoRetryAttemptsValue, &databaseRepresentation.TaskAutoRetryAttempts),
		handler.handleStringParameter(row.TraceLevelLevel, row.TraceLevelValue, &databaseRepresentation.TraceLevel),
		handler.handleStringParameter(row.UserTaskManagedInitialWarehouseSizeLevel, row.UserTaskManagedInitialWarehouseSizeValue, &databaseRepresentation.UserTaskManagedInitialWarehouseSize),
		handler.handleIntegerParameter(row.UserTaskMinimumTriggerIntervalInSecondsLevel, row.UserTaskMinimumTriggerIntervalInSecondsValue, &databaseRepresentation.UserTaskMinimumTriggerIntervalInSeconds),
		handler.handleIntegerParameter(row.UserTaskTimeoutMsLevel, row.UserTaskTimeoutMsValue, &databaseRepresentation.UserTaskTimeoutMs),
	)
	if errs != nil {
		return nil, errs
	}

	return databaseRepresentation, nil
}
