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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StreamOnTable_BasicUseCase(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(tableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	basic := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName())

	complete := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).
		WithAppendOnly(r.BooleanTrue).
		WithShowInitialRows(r.BooleanTrue).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnTableResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTableString(table.ID().FullyQualifiedName()).
			HasAppendOnlyString(r.BooleanDefault).
			HasShowInitialRowsString(r.BooleanDefault).
			HasCommentString(""),

		resourceshowoutputassert.StreamShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(table.ID().FullyQualifiedName()).
			HasMode(sdk.StreamModeDefault).
			HasComment("").
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeTable).
			HasBaseTables(table.ID()).
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
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.table_name", table.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeDefault))),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnTableResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasTableString(table.ID().FullyQualifiedName()).
			HasAppendOnlyString(r.BooleanTrue).
			HasShowInitialRowsString(r.BooleanTrue).
			HasCommentString(comment),

		resourceshowoutputassert.StreamShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(table.ID().FullyQualifiedName()).
			HasMode(sdk.StreamModeAppendOnly).
			HasComment(comment).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeTable).
			HasBaseTables(table.ID()).
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
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.table_name", table.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
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
				ImportStateVerifyIgnore: []string{"append_only", "copy_grants", "show_initial_rows"},
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
				ImportStateVerifyIgnore: []string{"append_only", "copy_grants", "show_initial_rows"},
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
					invokeactionassert.StreamDoesNotExist(t, id),
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

func TestAcc_StreamOnTable_CopyGrants(t *testing.T) {
	table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamModelWithoutCopyGrants := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).WithCopyGrants(false)
	streamModelWithCopyGrants := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).WithCopyGrants(true)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamModelWithoutCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModelWithoutCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamModelWithoutCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						createdOn = value
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
		},
	})
}

func TestAcc_StreamOnTable_CheckGrantsAfterRecreation(t *testing.T) {
	table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	table2, cleanupTable2 := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable2)

	role, cleanupRole := testClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).
		WithCopyGrants(true)
	model1WithoutCopyGrants := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName())
	model2 := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table2.ID().FullyQualifiedName()).
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

func TestAcc_StreamOnTable_PermadiffWhenIsStaleAndHasNoRetentionTime(t *testing.T) {
	schema, cleanupSchema := testClient().Schema.CreateSchemaWithOpts(t,
		testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(testClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)

	table, cleanupTable := testClient().Table.CreateWithChangeTrackingInSchema(t, schema.ID())
	t.Cleanup(cleanupTable)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	streamModel := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName())

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
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
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
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
			// note that it is stale again because we still have schema parameters set to 0
			{
				Config: config.FromModels(t, streamModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(streamModel.ResourceReference(), "stale", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
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

func TestAcc_StreamOnTable_StaleWithExternalChanges(t *testing.T) {
	schema, cleanupSchema := testClient().Schema.CreateSchemaWithOpts(t,
		testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(testClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(1),
			MaxDataExtensionTimeInDays: sdk.Pointer(1),
		},
	)
	t.Cleanup(cleanupSchema)

	table, cleanupTable := testClient().Table.CreateWithChangeTrackingInSchema(t, schema.ID())
	t.Cleanup(cleanupTable)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	streamModel := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName())

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			// initial creation does not lead to stale stream
			{
				Config: config.FromModels(t, streamModel),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
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
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
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
func TestAcc_StreamOnTable_At(t *testing.T) {
	table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	testClient().Table.InsertInt(t, table.ID())
	lastQueryId := testClient().Context.LastQueryId(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	commonModel := func() *model.StreamOnTableModel {
		return model.StreamOnTableBase("test", id, table.ID()).
			WithComment("foo").
			WithAppendOnly(r.BooleanTrue).
			WithShowInitialRows(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))
	modelWithStatement := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"statement": pluginconfig.StringVariable(lastQueryId),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithOffset.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasTableString(table.ID().FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasShowInitialRowsString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, modelWithOffset.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.mode", "APPEND_ONLY")),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStream.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStatement),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStatement.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				ResourceName:    modelWithOffset.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStreamOnTableResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAppendOnlyString(r.BooleanTrue).
						HasTableString(table.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnTable_Before(t *testing.T) {
	table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	testClient().Table.InsertInt(t, table.ID())
	lastQueryId := testClient().Context.LastQueryId(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	commonModel := func() *model.StreamOnTableModel {
		return model.StreamOnTableBase("test", id, table.ID()).
			WithComment("foo").
			WithAppendOnly(r.BooleanTrue).
			WithShowInitialRows(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))
	modelWithStatement := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"statement": pluginconfig.StringVariable(lastQueryId),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithOffset.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasTableString(table.ID().FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasShowInitialRowsString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, modelWithOffset.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.mode", "APPEND_ONLY")),
					assert.Check(resource.TestCheckResourceAttrSet(modelWithOffset.ResourceReference(), "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(modelWithOffset.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStream.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStatement),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, modelWithStream.ResourceReference()).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps
		},
	})
}

func TestAcc_StreamOnTable_InvalidConfiguration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidTableId := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "invalid")

	modelWithBefore := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithAppendOnly(r.BooleanFalse).
		WithShowInitialRows(r.BooleanFalse).
		WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	modelWithAt := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithAppendOnly(r.BooleanFalse).
		WithShowInitialRows(r.BooleanFalse).
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithBefore),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// multiple excluding options - at
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StreamOnTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithAt),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// invalid table id
			{
				Config:      config.FromModels(t, modelWithInvalidTableId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_StreamOnTable_ExternalStreamTypeChange(t *testing.T) {
	table, cleanupTable := testClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamModel := model.StreamOnTableBase("test", id, table.ID())

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
						resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeTable)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeTable),
					),
				),
			},
			// external change with a different type
			{
				PreConfig: func() {
					statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
					view, cleanupView := testClient().View.CreateView(t, statement)
					t.Cleanup(cleanupView)
					testClient().Stream.DropFunc(t, id)()
					externalChangeStream, cleanup := testClient().Stream.CreateOnViewWithRequest(t, sdk.NewCreateOnViewStreamRequest(id, view.ID()))
					t.Cleanup(cleanup)
					require.Equal(t, sdk.StreamSourceTypeView, *externalChangeStream.SourceType)
				},
				Config: config.FromModels(t, streamModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.StreamOnTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeTable)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeTable),
					),
				),
			},
		},
	})
}
