package main

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[SchemaRepresentation] = new(SchemaCsvRow)

type SchemaCsvRow struct {
	Comment                                      string `csv:"comment"`
	CreatedOn                                    string `csv:"created_on"`
	DatabaseName                                 string `csv:"database_name"`
	DroppedOn                                    string `csv:"dropped_on"`
	IsCurrent                                    string `csv:"is_current"`
	IsDefault                                    string `csv:"is_default"`
	Name                                         string `csv:"name"`
	Options                                      string `csv:"options"`
	Owner                                        string `csv:"owner"`
	OwnerRoleType                                string `csv:"owner_role_type"`
	RetentionTime                                string `csv:"retention_time"`
	CatalogValue                                 string `csv:"catalog_value"`
	CatalogLevel                                 string `csv:"catalog_level"`
	DataRetentionTimeInDaysValue                 string `csv:"data_retention_time_in_days_value"`
	DataRetentionTimeInDaysLevel                 string `csv:"data_retention_time_in_days_level"`
	DefaultDdlCollationValue                     string `csv:"default_ddl_collation_value"`
	DefaultDdlCollationLevel                     string `csv:"default_ddl_collation_level"`
	EnableConsoleOutputValue                     string `csv:"enable_console_output_value"`
	EnableConsoleOutputLevel                     string `csv:"enable_console_output_level"`
	ExternalVolumeValue                          string `csv:"external_volume_value"`
	ExternalVolumeLevel                          string `csv:"external_volume_level"`
	LogLevelValue                                string `csv:"log_level_value"`
	LogLevelLevel                                string `csv:"log_level_level"`
	MaxDataExtensionTimeInDaysValue              string `csv:"max_data_extension_time_in_days_value"`
	MaxDataExtensionTimeInDaysLevel              string `csv:"max_data_extension_time_in_days_level"`
	PipeExecutionPausedValue                     string `csv:"pipe_execution_paused_value"`
	PipeExecutionPausedLevel                     string `csv:"pipe_execution_paused_level"`
	QuotedIdentifiersIgnoreCaseValue             string `csv:"quoted_identifiers_ignore_case_value"`
	QuotedIdentifiersIgnoreCaseLevel             string `csv:"quoted_identifiers_ignore_case_level"`
	ReplaceInvalidCharactersValue                string `csv:"replace_invalid_characters_value"`
	ReplaceInvalidCharactersLevel                string `csv:"replace_invalid_characters_level"`
	StorageSerializationPolicyValue              string `csv:"storage_serialization_policy_value"`
	StorageSerializationPolicyLevel              string `csv:"storage_serialization_policy_level"`
	SuspendTaskAfterNumFailuresValue             string `csv:"suspend_task_after_num_failures_value"`
	SuspendTaskAfterNumFailuresLevel             string `csv:"suspend_task_after_num_failures_level"`
	TaskAutoRetryAttemptsValue                   string `csv:"task_auto_retry_attempts_value"`
	TaskAutoRetryAttemptsLevel                   string `csv:"task_auto_retry_attempts_level"`
	TraceLevelValue                              string `csv:"trace_level_value"`
	TraceLevelLevel                              string `csv:"trace_level_level"`
	UserTaskManagedInitialWarehouseSizeValue     string `csv:"user_task_managed_initial_warehouse_size_value"`
	UserTaskManagedInitialWarehouseSizeLevel     string `csv:"user_task_managed_initial_warehouse_size_level"`
	UserTaskMinimumTriggerIntervalInSecondsValue string `csv:"user_task_minimum_trigger_interval_in_seconds_value"`
	UserTaskMinimumTriggerIntervalInSecondsLevel string `csv:"user_task_minimum_trigger_interval_in_seconds_level"`
	UserTaskTimeoutMsValue                       string `csv:"user_task_timeout_ms_value"`
	UserTaskTimeoutMsLevel                       string `csv:"user_task_timeout_ms_level"`
}

type SchemaRepresentation struct {
	sdk.Schema

	// parameters
	Catalog                                 *string
	DataRetentionTimeInDays                 *int
	DefaultDdlCollation                     *string
	EnableConsoleOutput                     *bool
	ExternalVolume                          *string
	LogLevel                                *string
	MaxDataExtensionTimeInDays              *int
	PipeExecutionPaused                     *bool
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

func (row SchemaCsvRow) convert() (*SchemaRepresentation, error) {
	schemaRepresentation := &SchemaRepresentation{
		Schema: sdk.Schema{
			Name:          row.Name,
			IsDefault:     row.IsDefault == "Y",
			IsCurrent:     row.IsCurrent == "Y",
			DatabaseName:  row.DatabaseName,
			Owner:         row.Owner,
			Comment:       row.Comment,
			RetentionTime: row.RetentionTime,
			OwnerRoleType: row.OwnerRoleType,
		},
	}
	if row.Options != "" {
		schemaRepresentation.Options = &row.Options
	}

	handler := newParameterHandler(sdk.ParameterTypeSchema)
	errs := errors.Join(
		handler.handleIntegerParameter(row.DataRetentionTimeInDaysLevel, row.DataRetentionTimeInDaysValue, &schemaRepresentation.DataRetentionTimeInDays),
		handler.handleIntegerParameter(row.MaxDataExtensionTimeInDaysLevel, row.MaxDataExtensionTimeInDaysValue, &schemaRepresentation.MaxDataExtensionTimeInDays),
		handler.handleStringParameter(row.ExternalVolumeLevel, row.ExternalVolumeValue, &schemaRepresentation.ExternalVolume),
		handler.handleStringParameter(row.CatalogLevel, row.CatalogValue, &schemaRepresentation.Catalog),
		handler.handleBooleanParameter(row.PipeExecutionPausedLevel, row.PipeExecutionPausedValue, &schemaRepresentation.PipeExecutionPaused),
		handler.handleBooleanParameter(row.ReplaceInvalidCharactersLevel, row.ReplaceInvalidCharactersValue, &schemaRepresentation.ReplaceInvalidCharacters),
		handler.handleStringParameter(row.DefaultDdlCollationLevel, row.DefaultDdlCollationValue, &schemaRepresentation.DefaultDdlCollation),
		handler.handleStringParameter(row.StorageSerializationPolicyLevel, row.StorageSerializationPolicyValue, &schemaRepresentation.StorageSerializationPolicy),
		handler.handleStringParameter(row.LogLevelLevel, row.LogLevelValue, &schemaRepresentation.LogLevel),
		handler.handleStringParameter(row.TraceLevelLevel, row.TraceLevelValue, &schemaRepresentation.TraceLevel),
		handler.handleIntegerParameter(row.SuspendTaskAfterNumFailuresLevel, row.SuspendTaskAfterNumFailuresValue, &schemaRepresentation.SuspendTaskAfterNumFailures),
		handler.handleIntegerParameter(row.TaskAutoRetryAttemptsLevel, row.TaskAutoRetryAttemptsValue, &schemaRepresentation.TaskAutoRetryAttempts),
		handler.handleStringParameter(row.UserTaskManagedInitialWarehouseSizeLevel, row.UserTaskManagedInitialWarehouseSizeValue, &schemaRepresentation.UserTaskManagedInitialWarehouseSize),
		handler.handleIntegerParameter(row.UserTaskTimeoutMsLevel, row.UserTaskTimeoutMsValue, &schemaRepresentation.UserTaskTimeoutMs),
		handler.handleIntegerParameter(row.UserTaskMinimumTriggerIntervalInSecondsLevel, row.UserTaskMinimumTriggerIntervalInSecondsValue, &schemaRepresentation.UserTaskMinimumTriggerIntervalInSeconds),
		handler.handleBooleanParameter(row.QuotedIdentifiersIgnoreCaseLevel, row.QuotedIdentifiersIgnoreCaseValue, &schemaRepresentation.QuotedIdentifiersIgnoreCase),
		handler.handleBooleanParameter(row.EnableConsoleOutputLevel, row.EnableConsoleOutputValue, &schemaRepresentation.EnableConsoleOutput),
	)
	if errs != nil {
		return nil, errs
	}

	return schemaRepresentation, nil
}
