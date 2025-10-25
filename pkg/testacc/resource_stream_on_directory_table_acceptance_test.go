//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StreamOnDirectoryTable_BasicUseCase(t *testing.T) {
	// Schema analysis (from pkg/resources/stream_common.go and stream_on_directory_table.go):
	// - name: ForceNew: true (cannot be renamed)
	// - database: ForceNew: true (cannot be changed)
	// - schema: ForceNew: true (cannot be changed)
	// - stage: NOT force-new (can be updated)
	// - comment: Optional, NOT force-new
	// Result: Use same identifiers for basic/complete (name, database, schema are force-new), no additional force-new fields to handle

	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	stage, stageCleanup := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(stageCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	basic := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName())

	complete := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName()).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnDirectoryTableResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasStageString(stage.ID().Name()).
			HasCommentString(""),

		resourceshowoutputassert.StreamShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(stage.ID().Name()).
			HasMode(sdk.StreamModeDefault).
			HasComment("").
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeStage).
			HasBaseTablesPartiallyQualified(stage.ID().Name()).
			HasType("DELTA").
			HasStale(false).
			HasStaleAfterNotEmpty().
			HasInvalidReason("N/A").
			HasOwnerRoleType("ROLE"),

		// Describe output assertions
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.created_on")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.owner", testClient().Context.CurrentRole(t).Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.table_name", stage.ID().Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.base_tables.0", stage.ID().Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeDefault))),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.StreamOnDirectoryTableResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasStageString(stage.ID().Name()).
			HasCommentString(comment),

		resourceshowoutputassert.StreamShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTableName(stage.ID().Name()).
			HasMode(sdk.StreamModeDefault).
			HasComment(comment).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasSourceType(sdk.StreamSourceTypeStage).
			HasBaseTablesPartiallyQualified(stage.ID().Name()).
			HasType("DELTA").
			HasStale(false).
			HasStaleAfterNotEmpty().
			HasInvalidReason("N/A").
			HasOwnerRoleType("ROLE"),

		// Describe output assertions
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.created_on")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.owner", testClient().Context.CurrentRole(t).Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment", comment)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.table_name", stage.ID().Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.base_tables.0", stage.ID().Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.type", "DELTA")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.stale", "false")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.mode", string(sdk.StreamModeDefault))),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.stale_after")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
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
				ImportStateVerifyIgnore: []string{"copy_grants"},
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
				ImportStateVerifyIgnore: []string{"copy_grants"},
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
			// Create - with optionals (from scratch via taint)
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_Basic(t *testing.T) {
	stage, cleanupStage := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	baseModel := func() *model.StreamOnDirectoryTableModel {
		return model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName())
	}

	modelWithExtraFields := baseModel().
		WithCopyGrants(true).
		WithComment("foo")

	modelWithExtraFieldsModified := baseModel().
		WithCopyGrants(true).
		WithComment("bar")

	resourceId := helpers.EncodeResourceIdentifier(id)
	resourceName := modelWithExtraFields.ResourceReference()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: config.FromModels(t, baseModel()),
				Check: assertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// import without optionals
			{
				Config:       config.FromModels(t, baseModel()),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStreamOnDirectoryTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageString(stage.ID().Name()),
				),
			},
			// set all fields
			{
				Config: config.FromModels(t, modelWithExtraFields),
				Check: assertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("foo").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// external change
			{
				PreConfig: func() {
					testClient().Stream.Alter(t, sdk.NewAlterStreamRequest(id).WithSetComment("bar"))
				},
				Config: config.FromModels(t, modelWithExtraFields),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("foo").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// update fields
			{
				Config: config.FromModels(t, modelWithExtraFieldsModified),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale(false).
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("bar").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "bar")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// import
			{
				Config:       config.FromModels(t, modelWithExtraFieldsModified),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStreamOnDirectoryTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageString(stage.ID().Name()).
						HasCommentString("bar"),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_CopyGrants(t *testing.T) {
	stage, cleanupStage := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamOnDirectoryModelWithCopyGrants := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName()).WithCopyGrants(true)
	streamOnDirectoryModelWithoutCopyGrants := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName()).WithCopyGrants(false)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamOnDirectoryModelWithCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamOnDirectoryModelWithCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamOnDirectoryModelWithCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, streamOnDirectoryModelWithoutCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamOnDirectoryModelWithoutCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamOnDirectoryModelWithoutCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, streamOnDirectoryModelWithCopyGrants),
				Check: assertThat(t, resourceassert.StreamOnTableResource(t, streamOnDirectoryModelWithCopyGrants.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(streamOnDirectoryModelWithCopyGrants.ResourceReference(), "show_output.0.created_on", func(value string) error {
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

func TestAcc_StreamOnDirectoryTable_CheckGrantsAfterRecreation(t *testing.T) {
	stage, cleanupStage := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	stage2, cleanupStage2 := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage2)

	role, cleanupRole := testClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName()).WithCopyGrants(true)
	model1WithoutCopyGrants := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName())
	model2 := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage2.ID().FullyQualifiedName()).WithCopyGrants(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
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

func grantStreamPrivilegesConfig(resourceName string, roleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "grant" {
  privileges        = ["SELECT"]
  account_role_name = %[1]s
  on_schema_object {
    object_type = "STREAM"
    object_name = %[2]s.fully_qualified_name
  }
}

data "snowflake_grants" "grants" {
  depends_on = [snowflake_grant_privileges_to_account_role.grant, %[2]s]
  grants_on {
    object_type = "STREAM"
    object_name = %[2]s.fully_qualified_name
  }
}`, roleId.FullyQualifiedName(), resourceName)
}

// TODO (SNOW-1737932): Setting schema parameters related to retention time seems to have no affect on streams on directory tables.
// Adjust this test after this is fixed on Snowflake side.
func TestAcc_StreamOnDirectoryTable_RecreateWhenStale(t *testing.T) {
	stage, cleanupStage := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	schema, cleanupSchema := testClient().Schema.CreateSchemaWithOpts(t,
		testClient().Ids.RandomDatabaseObjectIdentifier(),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	streamModel := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamModel),
				Check: assertThat(t, resourceassert.StreamOnDirectoryTableResource(t, streamModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(streamModel.ResourceReference(), "show_output.0.stale", "false")),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_InvalidConfiguration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidStageId := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), "invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// invalid stage id
			{
				Config:      config.FromModels(t, modelWithInvalidStageId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_ExternalStreamTypeChange(t *testing.T) {
	stage, cleanupStage := testClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamModel := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName())

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
						resourceassert.StreamOnDirectoryTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeStage)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeStage),
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
						resourceassert.StreamOnDirectoryTableResource(t, streamModel.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeStage)),
						resourceshowoutputassert.StreamShowOutput(t, streamModel.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeStage),
					),
				),
			},
		},
	})
}
