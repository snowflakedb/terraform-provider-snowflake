// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type SchemaModel struct {
	Catalog                                 tfconfig.Variable `json:"catalog,omitempty"`
	Comment                                 tfconfig.Variable `json:"comment,omitempty"`
	DataRetentionTimeInDays                 tfconfig.Variable `json:"data_retention_time_in_days,omitempty"`
	Database                                tfconfig.Variable `json:"database,omitempty"`
	DefaultDdlCollation                     tfconfig.Variable `json:"default_ddl_collation,omitempty"`
	EnableConsoleOutput                     tfconfig.Variable `json:"enable_console_output,omitempty"`
	ExternalVolume                          tfconfig.Variable `json:"external_volume,omitempty"`
	FullyQualifiedName                      tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	IsTransient                             tfconfig.Variable `json:"is_transient,omitempty"`
	LogLevel                                tfconfig.Variable `json:"log_level,omitempty"`
	MaxDataExtensionTimeInDays              tfconfig.Variable `json:"max_data_extension_time_in_days,omitempty"`
	Name                                    tfconfig.Variable `json:"name,omitempty"`
	PipeExecutionPaused                     tfconfig.Variable `json:"pipe_execution_paused,omitempty"`
	QuotedIdentifiersIgnoreCase             tfconfig.Variable `json:"quoted_identifiers_ignore_case,omitempty"`
	ReplaceInvalidCharacters                tfconfig.Variable `json:"replace_invalid_characters,omitempty"`
	StorageSerializationPolicy              tfconfig.Variable `json:"storage_serialization_policy,omitempty"`
	SuspendTaskAfterNumFailures             tfconfig.Variable `json:"suspend_task_after_num_failures,omitempty"`
	TaskAutoRetryAttempts                   tfconfig.Variable `json:"task_auto_retry_attempts,omitempty"`
	TraceLevel                              tfconfig.Variable `json:"trace_level,omitempty"`
	UserTaskManagedInitialWarehouseSize     tfconfig.Variable `json:"user_task_managed_initial_warehouse_size,omitempty"`
	UserTaskMinimumTriggerIntervalInSeconds tfconfig.Variable `json:"user_task_minimum_trigger_interval_in_seconds,omitempty"`
	UserTaskTimeoutMs                       tfconfig.Variable `json:"user_task_timeout_ms,omitempty"`
	WithManagedAccess                       tfconfig.Variable `json:"with_managed_access,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Schema(
	resourceName string,
	database string,
	name string,
) *SchemaModel {
	s := &SchemaModel{ResourceModelMeta: config.Meta(resourceName, resources.Schema)}
	s.WithDatabase(database)
	s.WithName(name)
	return s
}

func SchemaWithDefaultMeta(
	database string,
	name string,
) *SchemaModel {
	s := &SchemaModel{ResourceModelMeta: config.DefaultMeta(resources.Schema)}
	s.WithDatabase(database)
	s.WithName(name)
	return s
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (s *SchemaModel) MarshalJSON() ([]byte, error) {
	type Alias SchemaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(s),
		DependsOn: s.DependsOn(),
	})
}

func (s *SchemaModel) WithDependsOn(values ...string) *SchemaModel {
	s.SetDependsOn(values...)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (s *SchemaModel) WithCatalog(catalog string) *SchemaModel {
	s.Catalog = tfconfig.StringVariable(catalog)
	return s
}

func (s *SchemaModel) WithComment(comment string) *SchemaModel {
	s.Comment = tfconfig.StringVariable(comment)
	return s
}

func (s *SchemaModel) WithDataRetentionTimeInDays(dataRetentionTimeInDays int) *SchemaModel {
	s.DataRetentionTimeInDays = tfconfig.IntegerVariable(dataRetentionTimeInDays)
	return s
}

func (s *SchemaModel) WithDatabase(database string) *SchemaModel {
	s.Database = tfconfig.StringVariable(database)
	return s
}

func (s *SchemaModel) WithDefaultDdlCollation(defaultDdlCollation string) *SchemaModel {
	s.DefaultDdlCollation = tfconfig.StringVariable(defaultDdlCollation)
	return s
}

func (s *SchemaModel) WithEnableConsoleOutput(enableConsoleOutput bool) *SchemaModel {
	s.EnableConsoleOutput = tfconfig.BoolVariable(enableConsoleOutput)
	return s
}

func (s *SchemaModel) WithExternalVolume(externalVolume string) *SchemaModel {
	s.ExternalVolume = tfconfig.StringVariable(externalVolume)
	return s
}

func (s *SchemaModel) WithFullyQualifiedName(fullyQualifiedName string) *SchemaModel {
	s.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return s
}

func (s *SchemaModel) WithIsTransient(isTransient string) *SchemaModel {
	s.IsTransient = tfconfig.StringVariable(isTransient)
	return s
}

func (s *SchemaModel) WithLogLevel(logLevel string) *SchemaModel {
	s.LogLevel = tfconfig.StringVariable(logLevel)
	return s
}

func (s *SchemaModel) WithMaxDataExtensionTimeInDays(maxDataExtensionTimeInDays int) *SchemaModel {
	s.MaxDataExtensionTimeInDays = tfconfig.IntegerVariable(maxDataExtensionTimeInDays)
	return s
}

func (s *SchemaModel) WithName(name string) *SchemaModel {
	s.Name = tfconfig.StringVariable(name)
	return s
}

func (s *SchemaModel) WithPipeExecutionPaused(pipeExecutionPaused bool) *SchemaModel {
	s.PipeExecutionPaused = tfconfig.BoolVariable(pipeExecutionPaused)
	return s
}

func (s *SchemaModel) WithQuotedIdentifiersIgnoreCase(quotedIdentifiersIgnoreCase bool) *SchemaModel {
	s.QuotedIdentifiersIgnoreCase = tfconfig.BoolVariable(quotedIdentifiersIgnoreCase)
	return s
}

func (s *SchemaModel) WithReplaceInvalidCharacters(replaceInvalidCharacters bool) *SchemaModel {
	s.ReplaceInvalidCharacters = tfconfig.BoolVariable(replaceInvalidCharacters)
	return s
}

func (s *SchemaModel) WithStorageSerializationPolicy(storageSerializationPolicy string) *SchemaModel {
	s.StorageSerializationPolicy = tfconfig.StringVariable(storageSerializationPolicy)
	return s
}

func (s *SchemaModel) WithSuspendTaskAfterNumFailures(suspendTaskAfterNumFailures int) *SchemaModel {
	s.SuspendTaskAfterNumFailures = tfconfig.IntegerVariable(suspendTaskAfterNumFailures)
	return s
}

func (s *SchemaModel) WithTaskAutoRetryAttempts(taskAutoRetryAttempts int) *SchemaModel {
	s.TaskAutoRetryAttempts = tfconfig.IntegerVariable(taskAutoRetryAttempts)
	return s
}

func (s *SchemaModel) WithTraceLevel(traceLevel string) *SchemaModel {
	s.TraceLevel = tfconfig.StringVariable(traceLevel)
	return s
}

func (s *SchemaModel) WithUserTaskManagedInitialWarehouseSize(userTaskManagedInitialWarehouseSize string) *SchemaModel {
	s.UserTaskManagedInitialWarehouseSize = tfconfig.StringVariable(userTaskManagedInitialWarehouseSize)
	return s
}

func (s *SchemaModel) WithUserTaskMinimumTriggerIntervalInSeconds(userTaskMinimumTriggerIntervalInSeconds int) *SchemaModel {
	s.UserTaskMinimumTriggerIntervalInSeconds = tfconfig.IntegerVariable(userTaskMinimumTriggerIntervalInSeconds)
	return s
}

func (s *SchemaModel) WithUserTaskTimeoutMs(userTaskTimeoutMs int) *SchemaModel {
	s.UserTaskTimeoutMs = tfconfig.IntegerVariable(userTaskTimeoutMs)
	return s
}

func (s *SchemaModel) WithWithManagedAccess(withManagedAccess string) *SchemaModel {
	s.WithManagedAccess = tfconfig.StringVariable(withManagedAccess)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *SchemaModel) WithCatalogValue(value tfconfig.Variable) *SchemaModel {
	s.Catalog = value
	return s
}

func (s *SchemaModel) WithCommentValue(value tfconfig.Variable) *SchemaModel {
	s.Comment = value
	return s
}

func (s *SchemaModel) WithDataRetentionTimeInDaysValue(value tfconfig.Variable) *SchemaModel {
	s.DataRetentionTimeInDays = value
	return s
}

func (s *SchemaModel) WithDatabaseValue(value tfconfig.Variable) *SchemaModel {
	s.Database = value
	return s
}

func (s *SchemaModel) WithDefaultDdlCollationValue(value tfconfig.Variable) *SchemaModel {
	s.DefaultDdlCollation = value
	return s
}

func (s *SchemaModel) WithEnableConsoleOutputValue(value tfconfig.Variable) *SchemaModel {
	s.EnableConsoleOutput = value
	return s
}

func (s *SchemaModel) WithExternalVolumeValue(value tfconfig.Variable) *SchemaModel {
	s.ExternalVolume = value
	return s
}

func (s *SchemaModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *SchemaModel {
	s.FullyQualifiedName = value
	return s
}

func (s *SchemaModel) WithIsTransientValue(value tfconfig.Variable) *SchemaModel {
	s.IsTransient = value
	return s
}

func (s *SchemaModel) WithLogLevelValue(value tfconfig.Variable) *SchemaModel {
	s.LogLevel = value
	return s
}

func (s *SchemaModel) WithMaxDataExtensionTimeInDaysValue(value tfconfig.Variable) *SchemaModel {
	s.MaxDataExtensionTimeInDays = value
	return s
}

func (s *SchemaModel) WithNameValue(value tfconfig.Variable) *SchemaModel {
	s.Name = value
	return s
}

func (s *SchemaModel) WithPipeExecutionPausedValue(value tfconfig.Variable) *SchemaModel {
	s.PipeExecutionPaused = value
	return s
}

func (s *SchemaModel) WithQuotedIdentifiersIgnoreCaseValue(value tfconfig.Variable) *SchemaModel {
	s.QuotedIdentifiersIgnoreCase = value
	return s
}

func (s *SchemaModel) WithReplaceInvalidCharactersValue(value tfconfig.Variable) *SchemaModel {
	s.ReplaceInvalidCharacters = value
	return s
}

func (s *SchemaModel) WithStorageSerializationPolicyValue(value tfconfig.Variable) *SchemaModel {
	s.StorageSerializationPolicy = value
	return s
}

func (s *SchemaModel) WithSuspendTaskAfterNumFailuresValue(value tfconfig.Variable) *SchemaModel {
	s.SuspendTaskAfterNumFailures = value
	return s
}

func (s *SchemaModel) WithTaskAutoRetryAttemptsValue(value tfconfig.Variable) *SchemaModel {
	s.TaskAutoRetryAttempts = value
	return s
}

func (s *SchemaModel) WithTraceLevelValue(value tfconfig.Variable) *SchemaModel {
	s.TraceLevel = value
	return s
}

func (s *SchemaModel) WithUserTaskManagedInitialWarehouseSizeValue(value tfconfig.Variable) *SchemaModel {
	s.UserTaskManagedInitialWarehouseSize = value
	return s
}

func (s *SchemaModel) WithUserTaskMinimumTriggerIntervalInSecondsValue(value tfconfig.Variable) *SchemaModel {
	s.UserTaskMinimumTriggerIntervalInSeconds = value
	return s
}

func (s *SchemaModel) WithUserTaskTimeoutMsValue(value tfconfig.Variable) *SchemaModel {
	s.UserTaskTimeoutMs = value
	return s
}

func (s *SchemaModel) WithWithManagedAccessValue(value tfconfig.Variable) *SchemaModel {
	s.WithManagedAccess = value
	return s
}
