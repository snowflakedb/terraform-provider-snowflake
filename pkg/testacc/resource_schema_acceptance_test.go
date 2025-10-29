//go:build non_account_level_tests

package testacc

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
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

func TestAcc_Schema_BasicUseCase(t *testing.T) {
	db, cleanupDb := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDb)

	id := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())
	newId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())
	comment := random.Comment()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	basic := model.Schema("test", id.DatabaseName(), id.Name())

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Schema(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasIsDefault(false).
			HasDatabaseName(id.DatabaseName()).
			HasOwnerNotEmpty().
			HasComment("").
			HasOptions("").
			HasRetentionTime("1").
			HasOwnerRoleTypeNotEmpty(),

		objectparametersassert.SchemaParameters(t, id).
			HasAllDefaultsExplicit(),

		resourceassert.SchemaResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasWithManagedAccessString(r.BooleanDefault).
			HasIsTransientString(r.BooleanDefault).
			HasCommentString(""),

		resourceshowoutputassert.SchemaShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasIsDefault(false).
			HasDatabaseName(id.DatabaseName()).
			HasOwnerNotEmpty().
			HasComment("").
			HasOptions("").
			HasOwnerRoleTypeNotEmpty(),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "0")),
	}

	complete := model.Schema("test", newId.DatabaseName(), newId.Name()).
		WithComment(comment).
		WithWithManagedAccess(r.BooleanTrue).
		WithDataRetentionTimeInDays(15).
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

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Schema(t, newId).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasIsDefault(false).
			HasDatabaseName(newId.DatabaseName()).
			HasOwnerNotEmpty().
			HasComment(comment).
			HasOptions("MANAGED ACCESS").
			HasRetentionTime("15").
			HasOwnerRoleTypeNotEmpty(),

		// TODO(SNOW-1501905): update assertions after updating schema parameters
		objectparametersassert.SchemaParameters(t, newId).
			HasDefaultDdlCollation("en_US"),

		resourceassert.SchemaResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasDatabaseString(newId.DatabaseName()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasWithManagedAccessString(r.BooleanTrue).
			HasIsTransientString(r.BooleanDefault).
			HasCommentString(comment).
			HasDataRetentionTimeInDaysString("15").
			HasMaxDataExtensionTimeInDaysString("3").
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasReplaceInvalidCharactersString("true").
			HasDefaultDdlCollationString("en_US").
			HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyCompatible)).
			HasLogLevelString(string(sdk.LogLevelInfo)).
			HasTraceLevelString(string(sdk.TraceLevelPropagate)).
			HasSuspendTaskAfterNumFailuresString("20").
			HasTaskAutoRetryAttemptsString("20").
			HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeXLarge)).
			HasUserTaskTimeoutMsString("1200000").
			HasUserTaskMinimumTriggerIntervalInSecondsString("120").
			HasQuotedIdentifiersIgnoreCaseString("true").
			HasEnableConsoleOutputString("true").
			HasPipeExecutionPausedString("true"),

		resourceshowoutputassert.SchemaShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasIsDefault(false).
			HasDatabaseName(newId.DatabaseName()).
			HasOwnerNotEmpty().
			HasComment(comment).
			HasOptions("MANAGED ACCESS").
			HasOwnerRoleTypeNotEmpty(),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "0")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
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
					"is_transient",
					"show_output.0.is_current",
					"with_managed_access",
				},
			},
			// Update - set optionals (including rename)
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
					"is_transient",
					"show_output.0.is_current",
					"with_managed_access",
				},
			},
			// Update - unset optionals (back to basic, with rename back)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Schema.Alter(t, id, &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:                 sdk.Int(2),
							MaxDataExtensionTimeInDays:              sdk.Int(15),
							ExternalVolume:                          sdk.Pointer(externalVolumeId),
							Catalog:                                 sdk.Pointer(catalogId),
							ReplaceInvalidCharacters:                sdk.Bool(true),
							DefaultDDLCollation:                     &sdk.StringAllowEmpty{Value: "en_US"},
							StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyCompatible),
							LogLevel:                                sdk.Pointer(sdk.LogLevelInfo),
							TraceLevel:                              sdk.Pointer(sdk.TraceLevelAlways),
							SuspendTaskAfterNumFailures:             sdk.Int(11),
							TaskAutoRetryAttempts:                   sdk.Int(1),
							UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeSmall),
							UserTaskTimeoutMs:                       sdk.Int(3600001),
							UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(31),
							EnableConsoleOutput:                     sdk.Bool(true),
							PipeExecutionPaused:                     sdk.Bool(true),
							QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
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
			// Empty config - ensure schema is destroyed
			{
				Config: " ",
				Check: assertThat(t,
					objectassert.SchemaIsMissing(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_Schema_CompleteUseCase(t *testing.T) {
	db, cleanupDb := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDb)

	id := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())
	comment := random.Comment()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	complete := model.Schema("test", id.DatabaseName(), id.Name()).
		WithComment(comment).
		WithWithManagedAccess(r.BooleanTrue).
		WithIsTransient(r.BooleanTrue).
		WithDataRetentionTimeInDays(1).
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional force-new fields)
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					objectassert.Schema(t, id).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIsDefault(false).
						HasDatabaseName(id.DatabaseName()).
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("TRANSIENT, MANAGED ACCESS").
						HasRetentionTime("1").
						HasOwnerRoleTypeNotEmpty(),

					resourceassert.SchemaResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasWithManagedAccessString(r.BooleanTrue).
						HasIsTransientString(r.BooleanTrue).
						HasCommentString(comment).
						HasDataRetentionTimeInDaysString("1").
						HasMaxDataExtensionTimeInDaysString("3").
						HasExternalVolumeString(externalVolumeId.Name()).
						HasCatalogString(catalogId.Name()).
						HasReplaceInvalidCharactersString("true").
						HasDefaultDdlCollationString("en_US").
						HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyCompatible)).
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasTraceLevelString(string(sdk.TraceLevelPropagate)).
						HasSuspendTaskAfterNumFailuresString("20").
						HasTaskAutoRetryAttemptsString("20").
						HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeXLarge)).
						HasUserTaskTimeoutMsString("1200000").
						HasUserTaskMinimumTriggerIntervalInSecondsString("120").
						HasQuotedIdentifiersIgnoreCaseString("true").
						HasEnableConsoleOutputString("true").
						HasPipeExecutionPausedString("true"),

					// TODO(SNOW-1501905): update assertions after updating schema parameters
					objectparametersassert.SchemaParameters(t, id).
						HasDefaultDdlCollation("en_US"),

					resourceshowoutputassert.SchemaShowOutput(t, complete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIsDefault(false).
						HasDatabaseName(id.DatabaseName()).
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("TRANSIENT, MANAGED ACCESS").
						HasOwnerRoleTypeNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "0")),
				),
			},
			// Import - with all optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"show_output.0.is_current",
				},
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
		CheckDestroy: CheckDestroyUsingLegacyIdParsing(t, resources.Schema),
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
