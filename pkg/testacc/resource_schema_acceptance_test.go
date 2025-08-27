//go:build !account_level_tests

package testacc

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	acchelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Schema_basic(t *testing.T) {
	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	basicSchemaModel := model.Schema("test", id.DatabaseName(), id.Name())
	fullSchemaModel := model.Schema("test", id.DatabaseName(), id.Name()).
		WithComment(comment).
		WithWithManagedAccess(r.BooleanTrue).
		WithIsTransient(r.BooleanFalse).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(3).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithTraceLevel(string(sdk.TraceLevelPropagate)).
		WithSuspendTaskAfterNumFailures(20).
		WithTaskAutoRetryAttempts(20).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeXLarge)).
		WithUserTaskTimeoutMs(1200000).
		WithUserTaskMinimumTriggerIntervalInSeconds(120).
		WithQuotedIdentifiersIgnoreCase(true).
		WithEnableConsoleOutput(true).
		WithPipeExecutionPaused(true)

	schemaModelWithExplicitTransientFalse := model.Schema("test", id.DatabaseName(), id.Name()).
		WithIsTransient(r.BooleanFalse)
	schemaModelWithExplicitTransientTrue := model.Schema("test", id.DatabaseName(), id.Name()).
		WithIsTransient(r.BooleanTrue)

	var (
		testDatabaseDataRetentionTimeInDays            = new(string)
		testDatabaseMaxDataExtensionTimeInDays         = new(string)
		accountExternalVolume                          = new(string)
		testDatabaseCatalog                            = new(string)
		accountReplaceInvalidCharacters                = new(string)
		accountDefaultDdlCollation                     = new(string)
		accountStorageSerializationPolicy              = new(string)
		accountLogLevel                                = new(string)
		accountTraceLevel                              = new(string)
		accountSuspendTaskAfterNumFailures             = new(string)
		accountTaskAutoRetryAttempts                   = new(string)
		accountUserTaskMangedInitialWarehouseSize      = new(string)
		accountUserTaskTimeoutMs                       = new(string)
		accountUserTaskMinimumTriggerIntervalInSeconds = new(string)
		accountQuotedIdentifiersIgnoreCase             = new(string)
		accountEnableConsoleOutput                     = new(string)
		accountPipeExecutionPaused                     = new(string)
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					params := testClient().Parameter.ShowAccountParameters(t)
					*testDatabaseDataRetentionTimeInDays = strconv.Itoa(testClient().Database.TestDatabaseDataRetentionTimeInDays())
					*testDatabaseMaxDataExtensionTimeInDays = strconv.Itoa(testClient().Database.TestDatabaseMaxDataExtensionTimeInDays())
					*accountExternalVolume = acchelpers.FindParameter(t, params, sdk.AccountParameterExternalVolume).Value
					*testDatabaseCatalog = testClient().Database.TestDatabaseCatalog().Name()
					*accountReplaceInvalidCharacters = acchelpers.FindParameter(t, params, sdk.AccountParameterReplaceInvalidCharacters).Value
					*accountDefaultDdlCollation = acchelpers.FindParameter(t, params, sdk.AccountParameterDefaultDDLCollation).Value
					*accountStorageSerializationPolicy = acchelpers.FindParameter(t, params, sdk.AccountParameterStorageSerializationPolicy).Value
					*accountLogLevel = acchelpers.FindParameter(t, params, sdk.AccountParameterLogLevel).Value
					*accountTraceLevel = acchelpers.FindParameter(t, params, sdk.AccountParameterTraceLevel).Value
					*accountSuspendTaskAfterNumFailures = acchelpers.FindParameter(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures).Value
					*accountTaskAutoRetryAttempts = acchelpers.FindParameter(t, params, sdk.AccountParameterTaskAutoRetryAttempts).Value
					*accountUserTaskMangedInitialWarehouseSize = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize).Value
					*accountUserTaskTimeoutMs = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskTimeoutMs).Value
					*accountUserTaskMinimumTriggerIntervalInSeconds = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds).Value
					*accountQuotedIdentifiersIgnoreCase = acchelpers.FindParameter(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase).Value
					*accountEnableConsoleOutput = acchelpers.FindParameter(t, params, sdk.AccountParameterEnableConsoleOutput).Value
					*accountPipeExecutionPaused = acchelpers.FindParameter(t, params, sdk.AccountParameterPipeExecutionPaused).Value
				},
				Config: accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "database", TestDatabaseName),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "is_transient", r.BooleanDefault),

					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", testDatabaseDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "max_data_extension_time_in_days", testDatabaseMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "catalog", testDatabaseCatalog),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr(basicSchemaModel.ResourceReference(), "pipe_execution_paused", accountPipeExecutionPaused),

					resource.TestCheckResourceAttrSet(basicSchemaModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(basicSchemaModel.ResourceReference(), "show_output.0.is_current"),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "show_output.0.database_name", id.DatabaseId().Name()),
					resource.TestCheckResourceAttrSet(basicSchemaModel.ResourceReference(), "show_output.0.owner"),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "show_output.0.options", ""),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, basicSchemaModel),
				ResourceName: basicSchemaModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", TestDatabaseName),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "with_managed_access", "false"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "is_transient", "false"),
				),
			},
			// set other fields
			{
				Config: accconfig.FromModels(t, fullSchemaModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(fullSchemaModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "with_managed_access", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "is_transient", "false"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "data_retention_time_in_days", "5"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "max_data_extension_time_in_days", "3"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "trace_level", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "enable_console_output", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "pipe_execution_paused", "true"),

					resource.TestCheckResourceAttrSet(fullSchemaModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(fullSchemaModel.ResourceReference(), "show_output.0.is_current"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "show_output.0.database_name", id.DatabaseId().Name()),
					resource.TestCheckResourceAttrSet(fullSchemaModel.ResourceReference(), "show_output.0.owner"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "show_output.0.options", "MANAGED ACCESS"),

					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.data_retention_time_in_days.0.value", "5"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.max_data_extension_time_in_days.0.value", "3"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.external_volume.0.value", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.catalog.0.value", catalogId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.replace_invalid_characters.0.value", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.default_ddl_collation.0.value", "en_US"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.storage_serialization_policy.0.value", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.log_level.0.value", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.trace_level.0.value", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.suspend_task_after_num_failures.0.value", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.task_auto_retry_attempts.0.value", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.user_task_managed_initial_warehouse_size.0.value", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.user_task_timeout_ms.0.value", "1200000"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.user_task_minimum_trigger_interval_in_seconds.0.value", "120"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.quoted_identifiers_ignore_case.0.value", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.enable_console_output.0.value", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "parameters.0.pipe_execution_paused.0.value", "true"),
				),
			},
			{
				Config:                  accconfig.FromModels(t, fullSchemaModel),
				ResourceName:            fullSchemaModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"show_output.0.is_current"},
			},
			// unset
			{
				Config: accconfig.FromModels(t, schemaModelWithExplicitTransientFalse),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelWithExplicitTransientFalse.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientFalse.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientFalse.ResourceReference(), "database", id.DatabaseId().Name()),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientFalse.ResourceReference(), "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientFalse.ResourceReference(), "is_transient", "false"),

					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "data_retention_time_in_days", testDatabaseDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "max_data_extension_time_in_days", testDatabaseMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "catalog", testDatabaseCatalog),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientFalse.ResourceReference(), "pipe_execution_paused", accountPipeExecutionPaused),
				),
			},
			// set is_transient - recreate
			{
				Config: accconfig.FromModels(t, schemaModelWithExplicitTransientTrue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(schemaModelWithExplicitTransientTrue.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientTrue.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientTrue.ResourceReference(), "database", id.DatabaseId().Name()),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientTrue.ResourceReference(), "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr(schemaModelWithExplicitTransientTrue.ResourceReference(), "is_transient", "true"),

					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "data_retention_time_in_days", testDatabaseDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "max_data_extension_time_in_days", testDatabaseMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "catalog", testDatabaseCatalog),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr(schemaModelWithExplicitTransientTrue.ResourceReference(), "pipe_execution_paused", accountPipeExecutionPaused),
				),
			},
		},
	})
}

func TestAcc_Schema_complete(t *testing.T) {
	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	fullSchemaModel := model.Schema("test", id.DatabaseName(), id.Name()).
		WithComment(comment).
		WithWithManagedAccess(r.BooleanTrue).
		WithIsTransient(r.BooleanTrue).
		WithDataRetentionTimeInDays(1).
		WithMaxDataExtensionTimeInDays(1).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithTraceLevel(string(sdk.TraceLevelPropagate)).
		WithSuspendTaskAfterNumFailures(20).
		WithTaskAutoRetryAttempts(20).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeXLarge)).
		WithUserTaskTimeoutMs(1200000).
		WithUserTaskMinimumTriggerIntervalInSeconds(120).
		WithQuotedIdentifiersIgnoreCase(true).
		WithEnableConsoleOutput(true).
		WithPipeExecutionPaused(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, fullSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "is_transient", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "with_managed_access", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "max_data_extension_time_in_days", "1"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "trace_level", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "enable_console_output", "true"),
					resource.TestCheckResourceAttr(fullSchemaModel.ResourceReference(), "pipe_execution_paused", "true"),
				),
			},
			{
				Config:                  accconfig.FromModels(t, fullSchemaModel),
				ResourceName:            fullSchemaModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"show_output.0.is_current"},
			},
		},
	})
}

func TestAcc_Schema_Rename(t *testing.T) {
	oldId := testClient().Ids.RandomDatabaseObjectIdentifier()
	newId := testClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	oldModel := model.Schema("test", oldId.DatabaseName(), oldId.Name()).
		WithComment(comment)
	newModel := model.Schema("test", newId.DatabaseName(), newId.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, oldModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(oldModel.ResourceReference(), "name", oldId.Name()),
					resource.TestCheckResourceAttr(oldModel.ResourceReference(), "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(oldModel.ResourceReference(), "database", TestDatabaseName),
					resource.TestCheckResourceAttr(oldModel.ResourceReference(), "comment", comment),
				),
			},
			{
				Config: accconfig.FromModels(t, newModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(newModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(newModel.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(newModel.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(newModel.ResourceReference(), "database", TestDatabaseName),
					resource.TestCheckResourceAttr(newModel.ResourceReference(), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_Schema_ManagePublicVersion_0_94_0(t *testing.T) {
	// use a separate db because this test relies on schema history
	db, cleanupDb := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(cleanupDb)

	schemaId := testClient().Ids.NewDatabaseObjectIdentifierInDatabase("PUBLIC", db.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// PUBLIC can not be created in v0.93
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				Config:            schemaV093(schemaId),
				ExpectError:       regexp.MustCompile("Error: error creating schema PUBLIC"),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.94.0"),
				Config:            schemaV094WithPipeExecutionPaused(schemaId, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.94.0"),
				PreConfig: func() {
					// In v0.94 `CREATE OR REPLACE` was called, so we should see a DROP event.
					schemas := testClient().Schema.ShowWithOptions(t, &sdk.ShowSchemaOptions{
						History: sdk.Pointer(true),
						Like: &sdk.Like{
							Pattern: sdk.String(schemaId.Name()),
						},
					})
					require.Len(t, schemas, 2)
					slices.SortFunc(schemas, func(x, y sdk.Schema) int {
						return cmp.Compare(x.DroppedOn.Unix(), y.DroppedOn.Unix())
					})
					require.Zero(t, schemas[0].DroppedOn)
					require.NotZero(t, schemas[1].DroppedOn)
				},
				Config: schemaV094WithPipeExecutionPaused(schemaId, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "false"),
				),
			},
		},
	})
}

func TestAcc_Schema_ManagePublicVersion_0_94_1(t *testing.T) {
	// use a separate db because this test relies on schema history
	db, cleanupDb := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(cleanupDb)

	schemaId := testClient().Ids.NewDatabaseObjectIdentifierInDatabase("PUBLIC", db.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// PUBLIC can not be created in v0.93
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				Config:            schemaV093(schemaId),
				ExpectError:       regexp.MustCompile("Error: error creating schema PUBLIC"),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   schemaV094WithPipeExecutionPaused(schemaId, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", db.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					// In newer versions, ALTER was called, so we should not see a DROP event.
					schemas := testClient().Schema.ShowWithOptions(t, &sdk.ShowSchemaOptions{
						History: sdk.Pointer(true),
						Like: &sdk.Like{
							Pattern: sdk.String(schemaId.Name()),
						},
					})
					require.Len(t, schemas, 1)
					require.Zero(t, schemas[0].DroppedOn)
				},
				Config: schemaV094WithPipeExecutionPaused(schemaId, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", db.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
		},
	})
}

// TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2209 issue.
func TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases(t *testing.T) {
	// It seems like Snowflake orders the output of SHOW command based on names, so they do matter
	db1Id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("A")
	db2Id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("B")

	_, database1Cleanup := testClient().Database.CreateDatabaseWithParametersSetWithId(t, db1Id)
	t.Cleanup(database1Cleanup)

	_, database2Cleanup := testClient().Database.CreateDatabaseWithParametersSetWithId(t, db2Id)
	t.Cleanup(database2Cleanup)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db1Id)
	schemaId2 := testClient().Ids.NewDatabaseObjectIdentifierInDatabase(schemaId.Name(), db2Id)

	schema1Model := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	schema2Model := model.Schema("test_2", schemaId2.DatabaseName(), schemaId2.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schema1Model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schema1Model.ResourceReference(), "name", schemaId.Name()),
					resource.TestCheckResourceAttr(schema1Model.ResourceReference(), "database", schemaId.DatabaseName()),
				),
			},
			{
				Config: accconfig.FromModels(t, schema1Model, schema2Model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schema1Model.ResourceReference(), "name", schemaId.Name()),
					resource.TestCheckResourceAttr(schema1Model.ResourceReference(), "database", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr(schema2Model.ResourceReference(), "name", schemaId2.Name()),
					resource.TestCheckResourceAttr(schema2Model.ResourceReference(), "database", schemaId2.DatabaseName()),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime(t *testing.T) {
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())

	basicSchemaModel := model.Schema("test", id.DatabaseName(), id.Name())
	schemaModelWithDataRetentionInDays5 := model.Schema("test", id.DatabaseName(), id.Name()).
		WithDataRetentionTimeInDays(5)
	schemaModelWithDataRetentionInDays15 := model.Schema("test", id.DatabaseName(), id.Name()).
		WithDataRetentionTimeInDays(15)
	schemaModelWithDataRetentionInDays0 := model.Schema("test", id.DatabaseName(), id.Name()).
		WithDataRetentionTimeInDays(0)
	schemaModelWithDataRetentionInDays3 := model.Schema("test", id.DatabaseName(), id.Name()).
		WithDataRetentionTimeInDays(3)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", "1"),
				),
			},
			// change param value in database
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, db.ID(), 50)
				},
				Config: accconfig.FromModels(t, basicSchemaModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(basicSchemaModel.ResourceReference(), "data_retention_time_in_days"),
						planchecks.ExpectDrift(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", sdk.String("1"), sdk.String("50")),
						planchecks.ExpectChange(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("50"), sdk.String("50")),
						planchecks.ExpectComputed(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", "50"),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModelWithDataRetentionInDays5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithDataRetentionInDays5.ResourceReference(), "data_retention_time_in_days", "5"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 5),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModelWithDataRetentionInDays15),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithDataRetentionInDays15.ResourceReference(), "data_retention_time_in_days", "15"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 15),
				),
			},
			{
				Config: accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", "50"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 50),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModelWithDataRetentionInDays0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithDataRetentionInDays0.ResourceReference(), "data_retention_time_in_days", "0"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 0),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModelWithDataRetentionInDays3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithDataRetentionInDays3.ResourceReference(), "data_retention_time_in_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())

	basicSchemaModel := model.Schema("test", id.DatabaseName(), id.Name())
	schemaModelWithDataRetentionInDays3 := model.Schema("test", id.DatabaseName(), id.Name()).
		WithDataRetentionTimeInDays(3)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", "1"),
				),
			},
			{
				PreConfig: func() {
					testClient().Schema.UpdateDataRetentionTime(t, id, 20)
				},
				Config: accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "data_retention_time_in_days", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModelWithDataRetentionInDays3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemaModelWithDataRetentionInDays3.ResourceReference(), "data_retention_time_in_days", "3"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveSchemaOutsideOfTerraform(t *testing.T) {
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifier()

	basicSchemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basicSchemaModel),
			},
			{
				PreConfig: func() {
					testClient().Schema.DropSchemaFunc(t, schemaId)()
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicSchemaModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveDatabaseOutsideOfTerraform(t *testing.T) {
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())

	basicSchemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basicSchemaModel),
			},
			{
				PreConfig: func() {
					dbCleanup()
				},
				Config: accconfig.FromModels(t, basicSchemaModel),
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("Failed to create schema"),
			},
		},
	})
}

func TestAcc_Schema_RemoveDatabaseOutsideOfTerraform_dbInConfig(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)

	databaseModel := model.DatabaseWithParametersSet("test", databaseId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name()).
		WithDependsOn(databaseModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel, schemaModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(databaseModel.ResourceReference(), "name", databaseId.Name())),
					assert.Check(resource.TestCheckResourceAttr(schemaModel.ResourceReference(), "name", schemaId.Name())),
				),
			},
			{
				PreConfig: func() {
					err := testClient().Database.DropDatabase(t, databaseId)
					require.NoError(t, err)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(databaseModel.ResourceReference(), plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(schemaModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, databaseModel, schemaModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "name", databaseId.Name())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaId.Name())),
				),
			},
		},
	})
}

func checkDatabaseAndSchemaDataRetentionTime(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, expectedDatabaseRetentionsDays int, expectedSchemaRetentionDays int) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		schema, err := testClient().Schema.Show(t, schemaId)
		require.NoError(t, err)

		database, err := testClient().Database.Show(t, schemaId.DatabaseId())
		require.NoError(t, err)

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := schema.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			require.NoError(t, err)
		}

		if database.RetentionTime != expectedDatabaseRetentionsDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionsDays, database.RetentionTime)
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		return nil
	}
}

func TestAcc_Schema_migrateFromVersion093WithoutManagedAccess(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()

	resourceName := "snowflake_schema.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				Config:            schemaV093(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   schemaV094(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "with_managed_access", r.BooleanDefault),
				),
			},
		},
	})
}

func TestAcc_Schema_migrateFromVersion093(t *testing.T) {
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifier()

	resourceName := "snowflake_schema.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				Config:            schemaV093WithIsManagedAndDataRetentionDays(id, tag.ID(), "foo", true, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.0.name", tag.Name),
					resource.TestCheckResourceAttr(resourceName, "tag.0.value", "foo"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   schemaV094WithManagedAccessAndDataRetentionTimeInDays(id, true, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "is_managed"),
					resource.TestCheckResourceAttr(resourceName, "with_managed_access", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "data_retention_days"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_time_in_days", "10"),
					resource.TestCheckNoResourceAttr(resourceName, "tag.#"),
				),
			},
		},
	})
}

func schemaV093WithIsManagedAndDataRetentionDays(schemaId sdk.DatabaseObjectIdentifier, tagId sdk.SchemaObjectIdentifier, tagValue string, isManaged bool, dataRetentionDays int) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database				= "%[1]s"
	name					= "%[2]s"
	is_managed				= %[7]t
	data_retention_days		= %[8]d
	tag {
		database = "%[3]s"
		schema = "%[4]s"
		name = "%[5]s"
		value = "%[6]s"
	}
}
`, schemaId.DatabaseName(), schemaId.Name(), tagId.DatabaseName(), tagId.SchemaName(), tagId.Name(), tagValue, isManaged, dataRetentionDays)
}

func schemaV093(schemaId sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database				= "%[1]s"
	name					= "%[2]s"
}
`, schemaId.DatabaseName(), schemaId.Name())
}

func schemaV094WithManagedAccessAndDataRetentionTimeInDays(schemaId sdk.DatabaseObjectIdentifier, isManaged bool, dataRetentionDays int) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database		 				= "%[1]s"
	name             				= "%[2]s"
	with_managed_access				= %[3]t
	data_retention_time_in_days		= %[4]d
}
`, schemaId.DatabaseName(), schemaId.Name(), isManaged, dataRetentionDays)
}

func schemaV094(schemaId sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database		 				= "%[1]s"
	name             				= "%[2]s"
}
`, schemaId.DatabaseName(), schemaId.Name())
}

func schemaV094WithPipeExecutionPaused(schemaId sdk.DatabaseObjectIdentifier, pipeExecutionPaused bool) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database		 				= "%[1]s"
	name             				= "%[2]s"
	pipe_execution_paused			= %[3]t
}
`, schemaId.DatabaseName(), schemaId.Name(), pipeExecutionPaused)
}

func TestAcc_Schema_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()

	basicSchemaModel := model.Schema("test", id.DatabaseName(), id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicSchemaModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModel.ResourceReference(), "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Schema_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	quotedDatabaseName := fmt.Sprintf(`"%s"`, id.DatabaseName())
	quotedName := fmt.Sprintf(`"%s"`, id.Name())

	basicSchemaModelWithQuotes := model.Schema("test", quotedDatabaseName, quotedName)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, basicSchemaModelWithQuotes),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "id", fmt.Sprintf(`"%s"|"%s"`, id.DatabaseName(), id.Name())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicSchemaModelWithQuotes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicSchemaModelWithQuotes.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicSchemaModelWithQuotes.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicSchemaModelWithQuotes.ResourceReference(), "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Schema_EmptyParameterAsDefaultValue(t *testing.T) {
	id := testClient().Ids.RandomDatabaseObjectIdentifier()
	parameterSet := model.Schema("test", id.DatabaseName(), id.Name()).WithDefaultDdlCollation("en_nz")
	parameterSetToNull := model.Schema("test", id.DatabaseName(), id.Name()).WithDefaultDdlCollationValue(accconfig.ReplacementPlaceholderVariable(accconfig.SnowflakeProviderConfigNull))
	parameterSetToEmptyString := model.Schema("test", id.DatabaseName(), id.Name()).WithDefaultDdlCollation("")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, parameterSet),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "default_ddl_collation", "en_nz"),
				),
			},
			{
				Config: accconfig.FromModels(t, parameterSetToNull),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(parameterSetToNull.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(parameterSetToNull.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(parameterSetToNull.ResourceReference(), "default_ddl_collation", ""),
				),
			},
			{
				Config: accconfig.FromModels(t, parameterSetToEmptyString),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(parameterSetToEmptyString.ResourceReference(), "database", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(parameterSetToEmptyString.ResourceReference(), "name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(parameterSetToEmptyString.ResourceReference(), "default_ddl_collation", "")),
					// Prove the parameter is set to empty string and on the right level.
					objectparametersassert.SchemaParameters(t, id).
						HasDefaultDdlCollation("").
						HasDefaultDdlCollationLevel(sdk.ParameterTypeSchema),
				),
			},
			// We set it back to demonstrate that going from set value straight to empty string won't work (an intermediate step with null value is required).
			{
				Config: accconfig.FromModels(t, parameterSet),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "default_ddl_collation", "en_nz"),
				),
			},
			{
				Config: accconfig.FromModels(t, parameterSetToEmptyString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(parameterSet.ResourceReference(), "default_ddl_collation", "en_nz"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
