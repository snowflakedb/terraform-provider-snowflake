//go:build account_level_tests

package testacc

import (
	"fmt"
	"testing"
	"time"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// TestAcc_SecondaryDatabase_BasicUseCase is the standard comprehensive acceptance test following the BasicUseCase pattern.
// This test consolidates the functionality of the old _basic and _complete tests.
//
// Resource Schema Analysis (from pkg/resources/secondary_database.go):
// - name: NOT force-new (line 22) - can be renamed in place
// - as_replica_of: ForceNew: true (line 30) - CRITICAL: must stay the same
// - is_transient: ForceNew: true (line 38) - CRITICAL: must stay the same
// - comment: Optional, NOT force-new (line 41)
// - database parameters: NOT force-new (inherited from databaseParametersSchema)
//
// ID Strategy: Since name is NOT force-new, use id and newId with DIFFERENT names but SAME as_replica_of
func TestAcc_SecondaryDatabase_BasicUseCase(t *testing.T) {
	// Setup: Generate identifiers and test values
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	// Create primary database on secondary account
	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		testClient().Account.GetAccountIdentifier(t),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	// Create external volume and catalog for complete configuration
	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	// Basic model - only required fields
	basic := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName())

	// Assertions for basic configuration
	assertBasic := []assert.TestCheckFuncProvider{
		// Check actual Snowflake object state
		objectassert.Database(t, id).
			HasName(id.Name()).
			HasTransient(false).
			HasComment(""),

		// Check parameters on the Snowflake object
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

		// Check Terraform resource state
		resourceassert.SecondaryDatabaseResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasAsReplicaOfString(externalPrimaryId.FullyQualifiedName()).
			HasCommentString(""),
	}

	// Complete model - all optional fields set
	// Using newId.Name() because name is NOT force-new, but same as_replica_of because it IS force-new
	complete := model.SecondaryDatabase("test", newId.Name(), externalPrimaryId.FullyQualifiedName()).
		WithComment(comment).
		WithDataRetentionTimeInDays(20).
		WithMaxDataExtensionTimeInDays(25).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithLogLevel(string(sdk.LogLevelDebug)).
		WithTraceLevel(string(sdk.TraceLevelAlways)).
		WithSuspendTaskAfterNumFailures(20).
		WithTaskAutoRetryAttempts(20).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeLarge)).
		WithUserTaskTimeoutMs(1200000).
		WithUserTaskMinimumTriggerIntervalInSeconds(60).
		WithQuotedIdentifiersIgnoreCase(true).
		WithEnableConsoleOutput(true)

	// Assertions for complete configuration
	assertComplete := []assert.TestCheckFuncProvider{
		// Check actual Snowflake object state
		objectassert.Database(t, newId).
			HasName(newId.Name()).
			HasTransient(false).
			HasComment(comment).
			HasRetentionTime(20),

		// Check parameters on the Snowflake object
		objectparametersassert.DatabaseParameters(t, newId).
			HasDataRetentionTimeInDays(20).
			HasMaxDataExtensionTimeInDays(25).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
			HasReplaceInvalidCharacters(true).
			HasDefaultDdlCollation("en_US").
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
			HasLogLevel(sdk.LogLevelDebug).
			HasTraceLevel(sdk.TraceLevelAlways).
			HasSuspendTaskAfterNumFailures(20).
			HasTaskAutoRetryAttempts(20).
			HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeLarge).
			HasUserTaskTimeoutMs(1200000).
			HasUserTaskMinimumTriggerIntervalInSeconds(60).
			HasQuotedIdentifiersIgnoreCase(true).
			HasEnableConsoleOutput(true),

		// Check Terraform resource state
		resourceassert.SecondaryDatabaseResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasAsReplicaOfString(externalPrimaryId.FullyQualifiedName()).
			HasCommentString(comment).
			HasDataRetentionTimeInDaysString("20").
			HasMaxDataExtensionTimeInDaysString("25").
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasReplaceInvalidCharactersString("true").
			HasDefaultDdlCollationString("en_US").
			HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyCompatible)).
			HasLogLevelString(string(sdk.LogLevelDebug)).
			HasTraceLevelString(string(sdk.TraceLevelAlways)).
			HasSuspendTaskAfterNumFailuresString("20").
			HasTaskAutoRetryAttemptsString("20").
			HasUserTaskManagedInitialWarehouseSizeString("LARGE").
			HasUserTaskTimeoutMsString("1200000").
			HasUserTaskMinimumTriggerIntervalInSecondsString("60").
			HasQuotedIdentifiersIgnoreCaseString("true").
			HasEnableConsoleOutputString("true"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			// Step 1: Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},

			// Step 2: Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Step 3: Update - set optionals (including rename)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},

			// Step 4: Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Step 5: Update - unset optionals (back to basic, with rename back)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},

			// Step 6: Update - detect external changes
			{
				PreConfig: func() {
					testClient().Database.Alter(t, id, &sdk.AlterDatabaseOptions{
						Set: &sdk.DatabaseSet{
							Comment: sdk.String(random.Comment()),
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

			// Step 7: Create - with optionals (from scratch via taint)
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

func TestAcc_CreateSecondaryDatabase_Basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		testClient().Account.GetAccountIdentifier(t),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	newId := testClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

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

	secondaryDatabaseModel := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName()).
		WithComment(comment)
	renamedSecondaryDatabaseModel := model.SecondaryDatabase("test", newId.Name(), externalPrimaryId.FullyQualifiedName()).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SharedDatabase),
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
				Config: accconfig.FromModels(t, secondaryDatabaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Rename + comment update
			{
				Config: accconfig.FromModels(t, renamedSecondaryDatabaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(renamedSecondaryDatabaseModel.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(renamedSecondaryDatabaseModel.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(renamedSecondaryDatabaseModel.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(renamedSecondaryDatabaseModel.ResourceReference(), "comment", newComment),

					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(renamedSecondaryDatabaseModel.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Import all values
			{
				Config:            accconfig.FromModels(t, renamedSecondaryDatabaseModel),
				ResourceName:      renamedSecondaryDatabaseModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSecondaryDatabase_complete(t *testing.T) {
	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	newExternalVolumeId, newExternalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(newExternalVolumeCleanup)

	newCatalogId, newCatalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(newCatalogCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(testClient().GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

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

	secondaryDatabaseModel := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName())
	secondaryDatabaseModelComplete := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName()).
		WithComment(comment).
		WithDataRetentionTimeInDays(20).
		WithMaxDataExtensionTimeInDays(25).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithLogLevel(string(sdk.LogLevelDebug)).
		WithTraceLevel(string(sdk.TraceLevelAlways)).
		WithSuspendTaskAfterNumFailures(20).
		WithTaskAutoRetryAttempts(20).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeLarge)).
		WithUserTaskTimeoutMs(1200000).
		WithUserTaskMinimumTriggerIntervalInSeconds(60).
		WithQuotedIdentifiersIgnoreCase(true).
		WithEnableConsoleOutput(true)
	secondaryDatabaseModelCompleteUpdated := model.SecondaryDatabase("test", newId.Name(), externalPrimaryId.FullyQualifiedName()).
		WithComment(newComment).
		WithDataRetentionTimeInDays(40).
		WithMaxDataExtensionTimeInDays(45).
		WithExternalVolume(newExternalVolumeId.Name()).
		WithCatalog(newCatalogId.Name()).
		WithReplaceInvalidCharacters(false).
		WithDefaultDdlCollation("en_GB").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithTraceLevel(string(sdk.TraceLevelPropagate)).
		WithSuspendTaskAfterNumFailures(40).
		WithTaskAutoRetryAttempts(40).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeXLarge)).
		WithUserTaskTimeoutMs(2400000).
		WithUserTaskMinimumTriggerIntervalInSeconds(120).
		WithQuotedIdentifiersIgnoreCase(false).
		WithEnableConsoleOutput(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryDatabase),
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
				Config: accconfig.FromModels(t, secondaryDatabaseModelComplete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "is_transient", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "max_data_extension_time_in_days", "25"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "log_level", string(sdk.LogLevelDebug)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "trace_level", string(sdk.TraceLevelAlways)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_managed_initial_warehouse_size", "LARGE"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "60"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "enable_console_output", "true"),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModelCompleteUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "is_transient", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "comment", newComment),

					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "data_retention_time_in_days", "40"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "max_data_extension_time_in_days", "45"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "external_volume", newExternalVolumeId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "catalog", newCatalogId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "default_ddl_collation", "en_GB"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "trace_level", string(sdk.TraceLevelPropagate)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "suspend_task_after_num_failures", "40"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "task_auto_retry_attempts", "40"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "user_task_managed_initial_warehouse_size", "XLARGE"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "user_task_timeout_ms", "2400000"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "quoted_identifiers_ignore_case", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelCompleteUpdated.ResourceReference(), "enable_console_output", "false"),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "is_transient", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(secondaryDatabaseModel.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModelComplete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "is_transient", "false"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "max_data_extension_time_in_days", "25"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "log_level", string(sdk.LogLevelDebug)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "trace_level", string(sdk.TraceLevelAlways)),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_managed_initial_warehouse_size", "LARGE"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "60"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr(secondaryDatabaseModelComplete.ResourceReference(), "enable_console_output", "true"),
				),
			},
			// Import all values
			{
				Config:            accconfig.FromModels(t, secondaryDatabaseModelComplete),
				ResourceName:      secondaryDatabaseModelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSecondaryDatabase_DataRetentionTimeInDays(t *testing.T) {
	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(testClient().GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	accountDataRetentionTimeInDays := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)

	secondaryDatabaseModel := func(
		dataRetentionTimeInDays *int,
	) *model.SecondaryDatabaseModel {
		secondaryDatabaseModel := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName()).
			WithMaxDataExtensionTimeInDays(10).
			WithExternalVolume(externalVolumeId.Name()).
			WithCatalog(catalogId.Name()).
			WithReplaceInvalidCharacters(true).
			WithDefaultDdlCollation("en_US").
			WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
			WithLogLevel(string(sdk.LogLevelOff)).
			WithTraceLevel(string(sdk.LogLevelOff)).
			WithSuspendTaskAfterNumFailures(10).
			WithTaskAutoRetryAttempts(10).
			WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeSmall)).
			WithUserTaskTimeoutMs(1200000).
			WithUserTaskMinimumTriggerIntervalInSeconds(120).
			WithQuotedIdentifiersIgnoreCase(true).
			WithEnableConsoleOutput(true)

		if dataRetentionTimeInDays != nil {
			secondaryDatabaseModel.WithDataRetentionTimeInDays(*dataRetentionTimeInDays)
		}

		return secondaryDatabaseModel
	}

	var revertAccountParameterChange func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel(sdk.Int(2))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel(sdk.Int(1))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel(nil)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterChange = testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "3")
					t.Cleanup(revertAccountParameterChange)
				},
				Config: accconfig.FromModels(t, secondaryDatabaseModel(nil)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "3"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterChange()
				},
				Config: accconfig.FromModels(t, secondaryDatabaseModel(nil)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel(sdk.Int(3))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "3"),
				),
			},
			{
				Config: accconfig.FromModels(t, secondaryDatabaseModel(nil)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
		},
	})
}

func TestAcc_SecondaryDatabase_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(testClient().GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	secondaryDatabaseModel := model.SecondaryDatabase("test", id.Name(), externalPrimaryId.FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, secondaryDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, secondaryDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_SecondaryDatabase_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())

	primaryDatabase, externalPrimaryId, _ := secondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(testClient().GetAccountLocator()),
	})
	unquotedExternalPrimaryId := fmt.Sprintf("%s.%s.%s", externalPrimaryId.AccountIdentifier().OrganizationName(), externalPrimaryId.AccountIdentifier().AccountName(), externalPrimaryId.Name())
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return secondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	secondaryDatabaseModel := model.SecondaryDatabase("test", quotedId, unquotedExternalPrimaryId)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, secondaryDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, secondaryDatabaseModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondaryDatabaseModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondaryDatabaseModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(secondaryDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}
