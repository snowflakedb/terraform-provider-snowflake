// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type SchemaResourceAssert struct {
	*assert.ResourceAssert
}

func SchemaResource(t *testing.T, name string) *SchemaResourceAssert {
	t.Helper()

	return &SchemaResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedSchemaResource(t *testing.T, id string) *SchemaResourceAssert {
	t.Helper()

	return &SchemaResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (s *SchemaResourceAssert) HasCatalogString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("catalog", expected))
	return s
}

func (s *SchemaResourceAssert) HasCommentString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", expected))
	return s
}

func (s *SchemaResourceAssert) HasDataRetentionTimeInDaysString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("data_retention_time_in_days", expected))
	return s
}

func (s *SchemaResourceAssert) HasDatabaseString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("database", expected))
	return s
}

func (s *SchemaResourceAssert) HasDefaultDdlCollationString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("default_ddl_collation", expected))
	return s
}

func (s *SchemaResourceAssert) HasEnableConsoleOutputString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("enable_console_output", expected))
	return s
}

func (s *SchemaResourceAssert) HasExternalVolumeString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("external_volume", expected))
	return s
}

func (s *SchemaResourceAssert) HasFullyQualifiedNameString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return s
}

func (s *SchemaResourceAssert) HasIsTransientString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("is_transient", expected))
	return s
}

func (s *SchemaResourceAssert) HasLogLevelString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("log_level", expected))
	return s
}

func (s *SchemaResourceAssert) HasMaxDataExtensionTimeInDaysString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("max_data_extension_time_in_days", expected))
	return s
}

func (s *SchemaResourceAssert) HasNameString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

func (s *SchemaResourceAssert) HasPipeExecutionPausedString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("pipe_execution_paused", expected))
	return s
}

func (s *SchemaResourceAssert) HasQuotedIdentifiersIgnoreCaseString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("quoted_identifiers_ignore_case", expected))
	return s
}

func (s *SchemaResourceAssert) HasReplaceInvalidCharactersString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("replace_invalid_characters", expected))
	return s
}

func (s *SchemaResourceAssert) HasStorageSerializationPolicyString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_serialization_policy", expected))
	return s
}

func (s *SchemaResourceAssert) HasSuspendTaskAfterNumFailuresString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("suspend_task_after_num_failures", expected))
	return s
}

func (s *SchemaResourceAssert) HasTaskAutoRetryAttemptsString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("task_auto_retry_attempts", expected))
	return s
}

func (s *SchemaResourceAssert) HasTraceLevelString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("trace_level", expected))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskManagedInitialWarehouseSizeString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("user_task_managed_initial_warehouse_size", expected))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskMinimumTriggerIntervalInSecondsString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("user_task_minimum_trigger_interval_in_seconds", expected))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskTimeoutMsString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("user_task_timeout_ms", expected))
	return s
}

func (s *SchemaResourceAssert) HasWithManagedAccessString(expected string) *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("with_managed_access", expected))
	return s
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (s *SchemaResourceAssert) HasNoCatalog() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("catalog"))
	return s
}

func (s *SchemaResourceAssert) HasNoComment() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("comment"))
	return s
}

func (s *SchemaResourceAssert) HasNoDataRetentionTimeInDays() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("data_retention_time_in_days"))
	return s
}

func (s *SchemaResourceAssert) HasNoDatabase() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("database"))
	return s
}

func (s *SchemaResourceAssert) HasNoDefaultDdlCollation() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("default_ddl_collation"))
	return s
}

func (s *SchemaResourceAssert) HasNoEnableConsoleOutput() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("enable_console_output"))
	return s
}

func (s *SchemaResourceAssert) HasNoExternalVolume() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("external_volume"))
	return s
}

func (s *SchemaResourceAssert) HasNoFullyQualifiedName() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return s
}

func (s *SchemaResourceAssert) HasNoIsTransient() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("is_transient"))
	return s
}

func (s *SchemaResourceAssert) HasNoLogLevel() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("log_level"))
	return s
}

func (s *SchemaResourceAssert) HasNoMaxDataExtensionTimeInDays() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("max_data_extension_time_in_days"))
	return s
}

func (s *SchemaResourceAssert) HasNoName() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}

func (s *SchemaResourceAssert) HasNoPipeExecutionPaused() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("pipe_execution_paused"))
	return s
}

func (s *SchemaResourceAssert) HasNoQuotedIdentifiersIgnoreCase() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("quoted_identifiers_ignore_case"))
	return s
}

func (s *SchemaResourceAssert) HasNoReplaceInvalidCharacters() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("replace_invalid_characters"))
	return s
}

func (s *SchemaResourceAssert) HasNoStorageSerializationPolicy() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("storage_serialization_policy"))
	return s
}

func (s *SchemaResourceAssert) HasNoSuspendTaskAfterNumFailures() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("suspend_task_after_num_failures"))
	return s
}

func (s *SchemaResourceAssert) HasNoTaskAutoRetryAttempts() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("task_auto_retry_attempts"))
	return s
}

func (s *SchemaResourceAssert) HasNoTraceLevel() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("trace_level"))
	return s
}

func (s *SchemaResourceAssert) HasNoUserTaskManagedInitialWarehouseSize() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("user_task_managed_initial_warehouse_size"))
	return s
}

func (s *SchemaResourceAssert) HasNoUserTaskMinimumTriggerIntervalInSeconds() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("user_task_minimum_trigger_interval_in_seconds"))
	return s
}

func (s *SchemaResourceAssert) HasNoUserTaskTimeoutMs() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("user_task_timeout_ms"))
	return s
}

func (s *SchemaResourceAssert) HasNoWithManagedAccess() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueNotSet("with_managed_access"))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SchemaResourceAssert) HasCatalogEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("catalog", ""))
	return s
}

func (s *SchemaResourceAssert) HasCommentEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", ""))
	return s
}

func (s *SchemaResourceAssert) HasDatabaseEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("database", ""))
	return s
}

func (s *SchemaResourceAssert) HasDefaultDdlCollationEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("default_ddl_collation", ""))
	return s
}

func (s *SchemaResourceAssert) HasExternalVolumeEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("external_volume", ""))
	return s
}

func (s *SchemaResourceAssert) HasFullyQualifiedNameEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return s
}

func (s *SchemaResourceAssert) HasIsTransientEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("is_transient", ""))
	return s
}

func (s *SchemaResourceAssert) HasLogLevelEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("log_level", ""))
	return s
}

func (s *SchemaResourceAssert) HasNameEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("name", ""))
	return s
}

func (s *SchemaResourceAssert) HasStorageSerializationPolicyEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_serialization_policy", ""))
	return s
}

func (s *SchemaResourceAssert) HasTraceLevelEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("trace_level", ""))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskManagedInitialWarehouseSizeEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("user_task_managed_initial_warehouse_size", ""))
	return s
}

func (s *SchemaResourceAssert) HasWithManagedAccessEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValueSet("with_managed_access", ""))
	return s
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (s *SchemaResourceAssert) HasCatalogNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("catalog"))
	return s
}

func (s *SchemaResourceAssert) HasCommentNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("comment"))
	return s
}

func (s *SchemaResourceAssert) HasDataRetentionTimeInDaysNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("data_retention_time_in_days"))
	return s
}

func (s *SchemaResourceAssert) HasDatabaseNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("database"))
	return s
}

func (s *SchemaResourceAssert) HasDefaultDdlCollationNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("default_ddl_collation"))
	return s
}

func (s *SchemaResourceAssert) HasEnableConsoleOutputNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("enable_console_output"))
	return s
}

func (s *SchemaResourceAssert) HasExternalVolumeNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("external_volume"))
	return s
}

func (s *SchemaResourceAssert) HasFullyQualifiedNameNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return s
}

func (s *SchemaResourceAssert) HasIsTransientNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("is_transient"))
	return s
}

func (s *SchemaResourceAssert) HasLogLevelNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("log_level"))
	return s
}

func (s *SchemaResourceAssert) HasMaxDataExtensionTimeInDaysNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("max_data_extension_time_in_days"))
	return s
}

func (s *SchemaResourceAssert) HasNameNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("name"))
	return s
}

func (s *SchemaResourceAssert) HasPipeExecutionPausedNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("pipe_execution_paused"))
	return s
}

func (s *SchemaResourceAssert) HasQuotedIdentifiersIgnoreCaseNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("quoted_identifiers_ignore_case"))
	return s
}

func (s *SchemaResourceAssert) HasReplaceInvalidCharactersNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("replace_invalid_characters"))
	return s
}

func (s *SchemaResourceAssert) HasStorageSerializationPolicyNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("storage_serialization_policy"))
	return s
}

func (s *SchemaResourceAssert) HasSuspendTaskAfterNumFailuresNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("suspend_task_after_num_failures"))
	return s
}

func (s *SchemaResourceAssert) HasTaskAutoRetryAttemptsNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("task_auto_retry_attempts"))
	return s
}

func (s *SchemaResourceAssert) HasTraceLevelNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("trace_level"))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskManagedInitialWarehouseSizeNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("user_task_managed_initial_warehouse_size"))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskMinimumTriggerIntervalInSecondsNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("user_task_minimum_trigger_interval_in_seconds"))
	return s
}

func (s *SchemaResourceAssert) HasUserTaskTimeoutMsNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("user_task_timeout_ms"))
	return s
}

func (s *SchemaResourceAssert) HasWithManagedAccessNotEmpty() *SchemaResourceAssert {
	s.AddAssertion(assert.ValuePresent("with_managed_access"))
	return s
}
