//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_InternalStage_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	modelBasic := model.InternalStage("test", TestDatabaseName, TestSchemaName, id.Name())

	modelComplete := model.InternalStage("test", TestDatabaseName, TestSchemaName, id.Name()).
		WithComment(comment).
		WithDirectoryEnabled(r.BooleanTrue).
		WithEncryptionSnowflakeFull()

	modelUpdated := model.InternalStage("test", TestDatabaseName, TestSchemaName, id.Name()).
		WithComment(changedComment).
		WithDirectoryEnabled(r.BooleanFalse).
		WithEncryptionSnowflakeFull()

	modelForceNewEncryption := model.InternalStage("test", TestDatabaseName, TestSchemaName, id.Name()).
		WithComment(changedComment).
		WithDirectoryEnabled(r.BooleanFalse).
		WithEncryptionSnowflakeSse()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			// Create with empty optionals (basic)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(TestDatabaseName).
						HasSchemaString(TestSchemaName).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeString("INTERNAL"),
					// Show output assertions
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.database_name", TestDatabaseName)),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.schema_name", TestSchemaName)),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.type", "INTERNAL")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.directory_enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "show_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "show_output.0.created_on")),
				),
			},
			// Import - without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedInternalStageResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(TestDatabaseName).
						HasSchemaString(TestSchemaName).
						HasCommentString("").
						HasDirectory(false, false).
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.directory_enabled", "false")),
				),
			},
			// Set optionals (complete)
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(TestDatabaseName).
						HasSchemaString(TestSchemaName).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "show_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "show_output.0.directory_enabled", "true")),
				),
			},
			// Import - complete
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
			},
			// Alter (update comment, directory.enable)
			{
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "show_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "show_output.0.directory_enabled", "false")),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.AlterInternalStage(t, sdk.NewAlterInternalStageStageRequest(id).
						WithComment(sdk.StringAllowEmpty{Value: changedComment}))
				},
				Config: accconfig.FromModels(t, modelUpdated),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelUpdated.ResourceReference()).
						HasCommentString(changedComment),
				),
			},
			// ForceNew - encryption change (requires recreate)
			{
				Config: accconfig.FromModels(t, modelForceNewEncryption),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelForceNewEncryption.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelForceNewEncryption.ResourceReference()).
						HasEncryptionSnowflakeSse(),
				),
			},
		},
	})
}

func TestAcc_InternalStage_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.InternalStage("test", TestDatabaseName, TestSchemaName, id.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndAutoRefresh(true, r.BooleanFalse).
		WithEncryptionSnowflakeFull()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.InternalStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(TestDatabaseName).
						HasSchemaString(TestSchemaName).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasDirectoryAutoRefreshString(r.BooleanFalse).
						HasEncryptionSnowflakeFull().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeString("INTERNAL"),
					resourceshowoutputassert.StageShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(TestDatabaseName).
						HasSchemaName(TestSchemaName).
						HasType(sdk.StageTypeInternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption"},
			},
		},
	})
}

func TestAcc_InternalStage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidAutoRefresh := model.InternalStageWithId(id).
		WithDirectoryEnabledAndAutoRefresh(true, "invalid")

	modelBothEncryptionTypes := model.InternalStageWithId(id).
		WithEncryptionBothTypes()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoRefresh),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*auto_refresh.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelBothEncryptionTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`encryption.0.snowflake_full,encryption.0.snowflake_sse.* can be specified`),
			},
		},
	})
}
