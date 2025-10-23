//go:build account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Database_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	basic := model.Database("test", id.Name())

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Database(t, id).
			HasName(id.Name()).
			HasTransient(false).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasOptions("").
			HasComment("").
			HasRetentionTime(1).
			HasResourceGroup("").
			HasKind("STANDARD"),
		objectparametersassert.DatabaseParameters(t, id).
			HasDefaultDataRetentionTimeInDaysValueExplicit().
			HasDefaultMaxDataExtensionTimeInDaysValueExplicit().
			HasDefaultExternalVolumeValueExplicit().
			HasDefaultCatalogValueExplicit().
			HasDefaultReplaceInvalidCharactersValueExplicit().
			HasDefaultDefaultDdlCollationValueExplicit().
			HasDefaultStorageSerializationPolicyValueExplicit().
			HasDefaultLogLevelValueExplicit().
			HasDefaultTraceLevelValueExplicit().
			HasDefaultSuspendTaskAfterNumFailuresValueExplicit().
			HasDefaultTaskAutoRetryAttemptsValueExplicit().
			HasUserTaskManagedInitialWarehouseSize("Medium").
			HasDefaultUserTaskTimeoutMsValueExplicit().
			HasDefaultUserTaskMinimumTriggerIntervalInSecondsValueExplicit().
			HasDefaultQuotedIdentifiersIgnoreCaseValueExplicit().
			HasDefaultEnableConsoleOutputValueExplicit(),
		resourceassert.DatabaseResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasIsTransientString("false").
			HasReplicationEmpty().
			HasCommentEmpty().
			HasAllDefaultParameters(),
	}

	complete := model.Database("test", newId.Name()).
		WithDropPublicSchemaOnCreation(true).
		WithReplication(secondaryAccountId, true, true).
		WithComment(comment).
		WithDataRetentionTimeInDays(2).
		WithMaxDataExtensionTimeInDays(15).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithTraceLevel(string(sdk.TraceLevelAlways)).
		WithSuspendTaskAfterNumFailures(11).
		WithTaskAutoRetryAttempts(1).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeSmall)).
		WithUserTaskTimeoutMs(3600001).
		WithUserTaskMinimumTriggerIntervalInSeconds(31).
		WithQuotedIdentifiersIgnoreCase(false).
		WithEnableConsoleOutput(true)

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Database(t, newId).
			HasName(newId.Name()).
			HasTransient(false).
			HasIsDefault(false).
			HasIsCurrent(false).
			HasOptions("").
			HasComment(comment).
			HasRetentionTime(2).
			HasResourceGroup("").
			HasKind("STANDARD"),
		objectparametersassert.DatabaseParameters(t, newId).
			HasDataRetentionTimeInDays(2).
			HasMaxDataExtensionTimeInDays(15).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
			HasReplaceInvalidCharacters(true).
			HasDefaultDdlCollation("en_US").
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
			HasLogLevel(sdk.LogLevelInfo).
			HasTraceLevel(sdk.TraceLevelAlways).
			HasSuspendTaskAfterNumFailures(11).
			HasTaskAutoRetryAttempts(1).
			HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeSmall).
			HasUserTaskTimeoutMs(3600001).
			HasUserTaskMinimumTriggerIntervalInSeconds(31).
			HasQuotedIdentifiersIgnoreCase(false).
			HasEnableConsoleOutput(true),
		resourceassert.DatabaseResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasIsTransientString("false").
			HasReplication(secondaryAccountId, true, true).
			HasCommentString(comment).
			HasDataRetentionTimeInDaysString("2").
			HasMaxDataExtensionTimeInDaysString("15").
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasReplaceInvalidCharactersString("true").
			HasDefaultDdlCollationString("en_US").
			HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyCompatible)).
			HasLogLevelString(string(sdk.LogLevelInfo)).
			HasTraceLevelString(string(sdk.TraceLevelAlways)).
			HasSuspendTaskAfterNumFailuresString("11").
			HasTaskAutoRetryAttemptsString("1").
			HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeSmall)).
			HasUserTaskTimeoutMsString("3600001").
			HasUserTaskMinimumTriggerIntervalInSecondsString("31").
			HasQuotedIdentifiersIgnoreCaseString("false").
			HasEnableConsoleOutputString("true"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"drop_public_schema_on_creation",
				},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"drop_public_schema_on_creation",
					"replication.0.ignore_edition_check",
				},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Database.AlterReplication(t, id, &sdk.AlterDatabaseReplicationOptions{
						EnableReplication: &sdk.EnableReplication{
							ToAccounts:         []sdk.AccountIdentifier{secondaryAccountId},
							IgnoreEditionCheck: sdk.Bool(true),
						},
					})
					testClient().Database.AlterFailover(t, id, &sdk.AlterDatabaseFailoverOptions{
						EnableFailover: &sdk.EnableFailover{
							ToAccounts: []sdk.AccountIdentifier{secondaryAccountId},
						},
					})
					testClient().Database.Alter(t, id, &sdk.AlterDatabaseOptions{
						Set: &sdk.DatabaseSet{
							DataRetentionTimeInDays:                 sdk.Int(2),
							MaxDataExtensionTimeInDays:              sdk.Int(15),
							ExternalVolume:                          sdk.Pointer(externalVolumeId),
							Catalog:                                 sdk.Pointer(catalogId),
							ReplaceInvalidCharacters:                sdk.Bool(true),
							DefaultDDLCollation:                     sdk.String("en_US"),
							StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyCompatible),
							LogLevel:                                sdk.Pointer(sdk.LogLevelInfo),
							TraceLevel:                              sdk.Pointer(sdk.TraceLevelAlways),
							SuspendTaskAfterNumFailures:             sdk.Int(11),
							TaskAutoRetryAttempts:                   sdk.Int(1),
							UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeSmall),
							UserTaskTimeoutMs:                       sdk.Int(3600001),
							UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(31),
							EnableConsoleOutput:                     sdk.Bool(true),
							Comment:                                 sdk.String(random.Comment()),
						},
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Create - with optionals
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_ComputedValues(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	secondaryAccountIdentifier := secondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	var (
		accountDataRetentionTimeInDays                 = new(string)
		accountMaxDataExtensionTimeInDays              = new(string)
		accountExternalVolume                          = new(string)
		accountCatalog                                 = new(string)
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
	)

	completeConfigVariables := config.Variables{
		"name":                                     config.StringVariable(id.Name()),
		"comment":                                  config.StringVariable(comment),
		"transient":                                config.BoolVariable(false),
		"account_identifier":                       config.StringVariable(secondaryAccountIdentifier),
		"with_failover":                            config.BoolVariable(true),
		"ignore_edition_check":                     config.BoolVariable(true),
		"data_retention_time_in_days":              config.IntegerVariable(20),
		"max_data_extension_time_in_days":          config.IntegerVariable(30),
		"external_volume":                          config.StringVariable(externalVolumeId.Name()),
		"catalog":                                  config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":               config.BoolVariable(true),
		"default_ddl_collation":                    config.StringVariable("en_US"),
		"storage_serialization_policy":             config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                              config.StringVariable(string(sdk.TraceLevelPropagate)),
		"suspend_task_after_num_failures":          config.IntegerVariable(20),
		"task_auto_retry_attempts":                 config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size": config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                     config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					params := testClient().Parameter.ShowAccountParameters(t)
					*accountDataRetentionTimeInDays = helpers.FindParameter(t, params, sdk.AccountParameterDataRetentionTimeInDays).Value
					*accountMaxDataExtensionTimeInDays = helpers.FindParameter(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays).Value
					*accountExternalVolume = helpers.FindParameter(t, params, sdk.AccountParameterExternalVolume).Value
					*accountCatalog = helpers.FindParameter(t, params, sdk.AccountParameterCatalog).Value
					*accountReplaceInvalidCharacters = helpers.FindParameter(t, params, sdk.AccountParameterReplaceInvalidCharacters).Value
					*accountDefaultDdlCollation = helpers.FindParameter(t, params, sdk.AccountParameterDefaultDDLCollation).Value
					*accountStorageSerializationPolicy = helpers.FindParameter(t, params, sdk.AccountParameterStorageSerializationPolicy).Value
					*accountLogLevel = helpers.FindParameter(t, params, sdk.AccountParameterLogLevel).Value
					*accountTraceLevel = helpers.FindParameter(t, params, sdk.AccountParameterTraceLevel).Value
					*accountSuspendTaskAfterNumFailures = helpers.FindParameter(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures).Value
					*accountTaskAutoRetryAttempts = helpers.FindParameter(t, params, sdk.AccountParameterTaskAutoRetryAttempts).Value
					*accountUserTaskMangedInitialWarehouseSize = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize).Value
					*accountUserTaskTimeoutMs = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskTimeoutMs).Value
					*accountUserTaskMinimumTriggerIntervalInSeconds = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds).Value
					*accountQuotedIdentifiersIgnoreCase = helpers.FindParameter(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase).Value
					*accountEnableConsoleOutput = helpers.FindParameter(t, params, sdk.AccountParameterEnableConsoleOutput).Value
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_database.test", "trace_level", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr("snowflake_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
		},
	})
}

func TestAcc_Database_Update(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := testClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	secondaryAccountIdentifier := secondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	databaseModel := model.DatabaseWithParametersSet("test", id.Name()).
		WithComment(comment)

	completeConfigVariables := config.Variables{
		"name":                                     config.StringVariable(newId.Name()),
		"comment":                                  config.StringVariable(newComment),
		"transient":                                config.BoolVariable(false),
		"account_identifier":                       config.StringVariable(secondaryAccountIdentifier),
		"with_failover":                            config.BoolVariable(true),
		"ignore_edition_check":                     config.BoolVariable(true),
		"data_retention_time_in_days":              config.IntegerVariable(20),
		"max_data_extension_time_in_days":          config.IntegerVariable(30),
		"external_volume":                          config.StringVariable(externalVolumeId.Name()),
		"catalog":                                  config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":               config.BoolVariable(true),
		"default_ddl_collation":                    config.StringVariable("en_US"),
		"storage_serialization_policy":             config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                              config.StringVariable(string(sdk.TraceLevelPropagate)),
		"suspend_task_after_num_failures":          config.IntegerVariable(20),
		"task_auto_retry_attempts":                 config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size": config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                     config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(databaseModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_database.test", "trace_level", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr("snowflake_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "enable_console_output", "true"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(databaseModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
		},
	})
}

// For now, this test can sometimes fail (if MAX_DATA_EXTENSION_TIME_IN_DAYS parameter is changed in the meantime).
// We need to test behavior for the database when it is not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_HierarchicalValues(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	paramDefault := new(string)
	var revertAccountParameterToDefault func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					*paramDefault = testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterMaxDataExtensionTimeInDays).Default
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault = testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "50")
					t.Cleanup(revertAccountParameterToDefault)
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "50"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault()
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
		},
	})
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_Replication(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountIdentifier := secondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	configVariables := func(id sdk.AccountObjectIdentifier, withReplication bool, withFailover bool, ignoreEditionCheck bool) config.Variables {
		if withReplication {
			return config.Variables{
				"name":                 config.StringVariable(id.Name()),
				"account_identifier":   config.StringVariable(secondaryAccountIdentifier),
				"with_failover":        config.BoolVariable(withFailover),
				"ignore_edition_check": config.BoolVariable(ignoreEditionCheck),
			}
		}
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(""),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables:         configVariables(id, true, true, true),
				ResourceName:            "snowflake_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
			},
		},
	})
}

func TestAcc_Database_IntParameter(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	databaseBasicConfig := config.Variables{
		"name": config.StringVariable(id.Name()),
	}

	databaseWithIntParameterConfig := func(dataRetentionTimeInDays int) config.Variables {
		return config.Variables{
			"name":                        config.StringVariable(id.Name()),
			"data_retention_time_in_days": config.IntegerVariable(dataRetentionTimeInDays),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			// create with setting one param
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigVariables: databaseWithIntParameterConfig(50),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionCreate, nil, sdk.String("50")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50"),
				),
			},
			// do not make any change (to check if there is no drift)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigVariables: databaseWithIntParameterConfig(50),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// import when param in config
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseWithIntParameterConfig(50),
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "data_retention_time_in_days", "50"),
				),
			},
			// change the param value in config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25"),
				),
			},
			// change param value on account - expect no changes
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
					revert := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "50")
					t.Cleanup(revert)
				},
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("25"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25"),
				),
			},
			// change the param value externally
			{
				PreConfig: func() {
					// clean after the previous step
					testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					// update externally
					testClient().Database.UpdateDataRetentionTime(t, id, 50)
				},
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("25"), sdk.String("50")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(25).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("25"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(1).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault),
				),
			},
			// import when param not in config (snowflake default)
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseBasicConfig,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "data_retention_time_in_days", "1"),
				),
			},
			// change the param value in config to snowflake default
			{
				ConfigVariables: databaseWithIntParameterConfig(1),
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("1"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(1).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("1"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1")), // Database default
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(1).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault),
				),
			},
			// change param value on account - change expected to be noop
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
					revert := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "50")
					t.Cleanup(revert)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("1"), sdk.String("50")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("50"), sdk.String("50")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(50).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeAccount),
				),
			},
			// import when param not in config (set on account)
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseBasicConfig,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "data_retention_time_in_days", "50"),
				),
				Check: assertThat(t,
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(50).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeAccount),
				),
			},
			// change param value on database
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, id, 50)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(50).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeAccount),
				),
			},
			// unset param on account
			{
				PreConfig: func() {
					testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("50"), sdk.String("1")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("1"), sdk.String("1")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1")),
					objectparametersassert.DatabaseParameters(t, id).HasDataRetentionTimeInDays(1).HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault),
				),
			},
		},
	})
}

func TestAcc_Database_StringValueSetOnDifferentParameterLevelWithSameValue(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := config.Variables{
		"name":    config.StringVariable(id.Name()),
		"catalog": config.StringVariable(catalogId.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/catalog"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
				),
			},
			{
				PreConfig: func() {
					require.Empty(t, testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterCatalog).Level)
					testClient().Database.UnsetCatalog(t, id)
					t.Cleanup(testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterCatalog, catalogId.Name()))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "catalog"),
						planchecks.ExpectChange("snowflake_database.test", "catalog", tfjson.ActionUpdate, sdk.String(catalogId.Name()), nil),
						planchecks.ExpectComputed("snowflake_database.test", "catalog", true),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Database/catalog"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
				),
			},
		},
	})
}

// For now, this test can sometimes fail (if other account parameters are changed in the meantime).
// We could set the known parameters here.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_UpgradeWithDataRetentionSet(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            databaseStateUpgraderDataRetentionSet(id, comment, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "10"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            databaseStateUpgraderDataRetentionSet(id, comment, 10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "10"),
				),
			},
		},
	})
}

func databaseStateUpgraderDataRetentionSet(id sdk.AccountObjectIdentifier, comment string, dataRetention int) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	comment = "%s"
	data_retention_time_in_days = %d
}
`, id.Name(), comment, dataRetention)
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_WithReplication(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountLocator := secondaryTestClient().GetAccountLocator()
	secondaryAccountIdentifier := secondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            databaseStateUpgraderWithReplicationOld(id, secondaryAccountLocator),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.0", secondaryAccountLocator),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderWithReplicationNew(id, secondaryAccountIdentifier),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "replication"),
						// Updates in place (no ALTER DATABASE is called)
						planchecks.ExpectChange("snowflake_database.test", "replication.0.ignore_edition_check", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_database.test", "replication.0.enable_to_account", tfjson.ActionUpdate, sdk.String(fmt.Sprintf("[map[account_identifier:%s with_failover:false]]", secondaryAccountIdentifier)), sdk.String(fmt.Sprintf("[map[account_identifier:%s with_failover:false]]", secondaryAccountIdentifier))),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
					resource.TestCheckNoResourceAttr("snowflake_database.test", "replication_configuration"),
				),
			},
		},
	})
}

func databaseStateUpgraderWithReplicationOld(id sdk.AccountObjectIdentifier, enableToAccount string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	replication_configuration {
		accounts = ["%s"]
		ignore_edition_check = true
	}
}
`, id.Name(), enableToAccount)
}

func databaseStateUpgraderWithReplicationNew(id sdk.AccountObjectIdentifier, enableToAccount string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	replication {
		enable_to_account {
			account_identifier = %s
			with_failover = false
		}
		ignore_edition_check = true
	}
}
`, id.Name(), strconv.Quote(enableToAccount))
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_WithoutPublicSchema(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: databaseWithDropPublicSchemaConfig(id, true),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name())),
					objectassert.DatabaseDescribe(t, id).DoesNotContainPublicSchema(),
				),
			},
			// Change in parameter shouldn't change the state Snowflake
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionNoop),
					},
				},
				Config: databaseWithDropPublicSchemaConfig(id, false),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name())),
					objectassert.DatabaseDescribe(t, id).DoesNotContainPublicSchema(),
				),
			},
		},
	})
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_WithPublicSchema(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: databaseWithDropPublicSchemaConfig(id, false),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name())),
					objectassert.DatabaseDescribe(t, id).ContainsPublicSchema(),
				),
			},
			// Change in parameter shouldn't change the state Snowflake
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionNoop),
					},
				},
				Config: databaseWithDropPublicSchemaConfig(id, true),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name())),
					objectassert.DatabaseDescribe(t, id).ContainsPublicSchema(),
				),
			},
		},
	})
}

func databaseWithDropPublicSchemaConfig(id sdk.AccountObjectIdentifier, withDropPublicSchema bool) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	drop_public_schema_on_creation = %s
}
`, id.Name(), strconv.FormatBool(withDropPublicSchema))
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            databaseConfigBasic(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   databaseConfigBasic(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
				),
			},
		},
	})
}

func databaseConfigBasic(name string) string {
	return fmt.Sprintf(`resource "snowflake_database" "test" {
		name = "%v"
	}`, name)
}

// For now, this test can sometimes fail (if account parameters are changed in the meantime).
// We could set the known parameters here, however, we need to test behavior for the database when they are not set.
// We could try ignoring the changes to parameters too.
func TestAcc_Database_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)
	quotedExternalVolumeId := fmt.Sprintf(`\"%s\"`, externalVolumeId.Name())

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)
	quotedCatalogId := fmt.Sprintf(`\"%s\"`, catalogId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             databaseConfigBasicWithExternalVolumeAndCatalog(quotedId, quotedExternalVolumeId, quotedCatalogId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   databaseConfigBasicWithExternalVolumeAndCatalog(quotedId, quotedExternalVolumeId, quotedCatalogId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
				),
			},
		},
	})
}

func databaseConfigBasicWithExternalVolumeAndCatalog(databaseName string, externalVolumeName string, catalogName string) string {
	return fmt.Sprintf(`resource "snowflake_database" "test" {
		name = "%v"
		external_volume = "%v"
		catalog = "%v"
	}`, databaseName, externalVolumeName, catalogName)
}
