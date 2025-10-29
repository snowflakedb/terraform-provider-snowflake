//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"
	pluginconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StreamOnExternalTable_BasicUseCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	basic := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue)

	complete := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnExternalTableResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasExternalTableString(externalTable.ID().FullyQualifiedName()).
			HasInsertOnlyString(r.BooleanTrue).
			HasCommentString(""),

		resourceshowoutputassert.StreamShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(externalTable.ID().FullyQualifiedName()).
			HasMode(sdk.StreamModeInsertOnly).
			HasComment("").
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeExternalTable).
			HasBaseTables(externalTable.ID()).
			HasType("DELTA").
			HasStale(false).
			HasStaleAfterNotEmpty().
			HasInvalidReason("N/A").
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.created_on")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.owner", testClient().Context.CurrentRole(t).Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnExternalTableResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasExternalTableString(externalTable.ID().FullyQualifiedName()).
			HasInsertOnlyString(r.BooleanTrue).
			HasCommentString(comment),

		resourceshowoutputassert.StreamShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(externalTable.ID().FullyQualifiedName()).
			HasMode(sdk.StreamModeInsertOnly).
			HasComment(comment).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeExternalTable).
			HasBaseTables(externalTable.ID()).
			HasType("DELTA").
			HasStale(false).
			HasStaleAfterNotEmpty().
			HasInvalidReason("N/A").
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.created_on")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.owner", testClient().Context.CurrentRole(t).Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment", comment)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            basic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"insert_only", "copy_grants"},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            complete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"insert_only", "copy_grants"},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Stream.Alter(t, sdk.NewAlterStreamRequest(id).WithSetComment(comment))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure stream is destroyed before the next step
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.StreamDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_CopyGrants(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamModelWithCopyGrants := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue).
		WithCopyGrants(true)
	streamModelWithoutCopyGrants := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue).
		WithCopyGrants(false)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamModelWithCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModelWithCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamModelWithCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, streamModelWithoutCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModelWithoutCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamModelWithoutCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, streamModelWithCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModelWithCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamModelWithCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_CheckGrantsAfterRecreation(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	externalTable2, externalTableCleanup2 := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup2)

	role, cleanupRole := testClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue).
		WithCopyGrants(true)
	model1WithoutCopyGrants := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue)
	model2 := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable2.ID().FullyQualifiedName()).
		WithInsertOnly(r.BooleanTrue).
		WithCopyGrants(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model1) + grantStreamPrivilegesConfig(model1.ResourceReference(), role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: config.FromModels(t, model2) + grantStreamPrivilegesConfig(model2.ResourceReference(), role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config:             config.FromModels(t, model1WithoutCopyGrants) + grantStreamPrivilegesConfig(model1WithoutCopyGrants.ResourceReference(), role.ID()),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.grant", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "1"),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_PermadiffWhenIsStaleAndHasNoRetentionTime(t *testing.T) {
	schema, cleanupSchema := testClient().Schema.CreateSchemaWithOpts(t,
		testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(testClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)

	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateInSchemaWithLocation(t, stage.Location(), schema.ID())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	streamModel := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).WithInsertOnly(r.BooleanTrue)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			// check that stale state is marked properly and forces an update
			{
				Config: config.FromModels(t, streamModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(streamModel.ResourceReference(), "stale", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanTrue),
					assert.Check(resource.TestCheckResourceAttr(streamModel.ResourceReference(), "show_output.0.stale", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(streamModel.ResourceReference(), "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// check that the resource was recreated
			// note that it is stale again because we still have schema parameters set to 0, this results in a permadiff
			{
				Config: config.FromModels(t, streamModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(streamModel.ResourceReference(), "stale", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanTrue),
					assert.Check(resource.TestCheckResourceAttr(streamModel.ResourceReference(), "show_output.0.stale", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(streamModel.ResourceReference(), "show_output.0.created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("stream was not recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_StaleWithExternalChanges(t *testing.T) {
	schema, cleanupSchema := testClient().Schema.CreateSchemaWithOpts(t,
		testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(testClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(1),
			MaxDataExtensionTimeInDays: sdk.Pointer(1),
		},
	)
	t.Cleanup(cleanupSchema)

	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateInSchemaWithLocation(t, stage.Location(), schema.ID())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	streamModel := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), externalTable.ID().FullyQualifiedName()).WithInsertOnly(r.BooleanTrue)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			// initial creation does not lead to stale stream
			{
				Config: config.FromModels(t, streamModel),
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(streamModel.ResourceReference(), "show_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttrWith(streamModel.ResourceReference(), "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// changing the value externally on schema
			{
				PreConfig: func() {
					testClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:    sdk.Int(0),
							MaxDataExtensionTimeInDays: sdk.Int(0),
						},
					})
					assertThatObject(t, objectassert.Stream(t, id).
						HasName(id.Name()).
						HasStale(true),
					)
					testClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:    sdk.Int(1),
							MaxDataExtensionTimeInDays: sdk.Int(1),
						},
					})
					assertThatObject(t, objectassert.Stream(t, id).
						HasName(id.Name()).
						HasStale(false),
					)
				},
				Config: config.FromModels(t, streamModel),
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(streamModel.ResourceReference(), "show_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttrWith(streamModel.ResourceReference(), "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnExternalTable_At(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	commonModel := func() *model.StreamOnExternalTableModel {
		return model.StreamOnExternalTableBase("test", id, externalTable.ID()).
			WithComment("foo").
			WithInsertOnly(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, modelWithOffset.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, modelWithOffset.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables(externalTable.ID()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStream.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps and statements
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnExternalTable_Before(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	commonModel := func() *model.StreamOnExternalTableModel {
		return model.StreamOnExternalTableBase("test", id, externalTable.ID()).
			WithComment("foo").
			WithInsertOnly(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assertThat(t, resourceassert.StreamOnExternalTableResource(t, modelWithOffset.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, modelWithOffset.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables(externalTable.ID()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStream.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps and statements
		},
	})
}

func TestAcc_StreamOnExternalTable_InvalidConfiguration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidExternalTableId := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "invalid")

	modelWithBefore := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithInsertOnly(r.BooleanTrue).
		WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	modelWithAt := model.StreamOnExternalTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithInsertOnly(r.BooleanTrue).
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// multiple excluding options - before
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithBefore),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// multiple excluding options - at
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithAt),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// invalid external table id
			{
				Config:      config.FromModels(t, modelWithInvalidExternalTableId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_ExternalStreamTypeChange(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := testClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamModel := model.StreamOnExternalTableBase("test", id, externalTable.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeExternalTable)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeExternalTable),
					),
				),
			},
			// external change with a different type
			{
				PreConfig: func() {
					table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
					t.Cleanup(cleanupTable)
					testClient().Stream.DropFunc(t, id)()
					externalChangeStream, cleanup := testClient().Stream.CreateOnTableWithRequest(t, sdk.NewCreateOnTableStreamRequest(id, table.ID()))
					t.Cleanup(cleanup)
					require.Equal(t, sdk.StreamSourceTypeTable, *externalChangeStream.SourceType)
				},
				Config: config.FromModels(t, streamModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.StreamOnExternalTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeExternalTable)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeExternalTable),
					),
				),
			},
		},
	})
}
