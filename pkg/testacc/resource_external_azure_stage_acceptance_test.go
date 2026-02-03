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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalAzureStage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	storageIntegrationId := ids.PrecreatedAzureStorageIntegration

	masterKey := random.AzureCseMasterKey()
	azureUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	// azureSasToken := testenvs.GetOrSkipTest(t, testenvs.AzureExternalSasToken)

	modelBasic := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name())

	modelComplete := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(true, r.BooleanTrue, r.BooleanFalse).
		WithEncryptionAzureCse(masterKey)

	modelUpdated := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithDirectoryEnabledAndOptions(false, r.BooleanTrue, r.BooleanFalse).
		WithEncryptionNone()

	modelEncryptionNoneWithComment := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithEncryptionNone()

	modelNoEncryption := model.ExternalAzureStageWithId(newId, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAzureStage),
		Steps: []resource.TestStep{
			// Create with empty optionals (basic)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
				),
			},
			// Import - without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "use_privatelink_endpoint", "false"),
				),
			},
			// Set optionals (complete)
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
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
			// Alter (update comment, directory.enable, encryption)
			{
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanFalse).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelUpdated.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(changedComment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.AlterExternalAzureStage(t, sdk.NewAlterExternalAzureStageStageRequest(id).
						WithComment(sdk.StringAllowEmpty{Value: changedComment}))
				},
				Config: accconfig.FromModels(t, modelUpdated),
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEnableString(r.BooleanFalse).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
				),
			},
			// ForceNew - unset directory
			{
				Config: accconfig.FromModels(t, modelEncryptionNoneWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelEncryptionNoneWithComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelEncryptionNoneWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEmpty().
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
				),
			},
			// ForceNew - unset encryption and rename
			{
				Config: accconfig.FromModels(t, modelNoEncryption),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelNoEncryption.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelNoEncryption.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
				),
			},
		},
	})
}

func TestAcc_ExternalAzureStage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	storageIntegrationId := ids.PrecreatedAzureStorageIntegration

	azureUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	masterKey := random.AzureCseMasterKey()

	modelComplete := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(true, r.BooleanTrue, r.BooleanFalse).
		WithEncryptionAzureCse(masterKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAzureStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectoryEnableString(r.BooleanTrue).
						HasDirectoryRefreshOnCreateString(r.BooleanTrue).
						HasDirectoryAutoRefreshString(r.BooleanFalse).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
				),
			},
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "use_privatelink_endpoint", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.auto_refresh", "false"),
				),
			},
		},
	})
}

func TestAcc_ExternalAzureStage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	storageIntegrationId := testClient().Ids.RandomAccountObjectIdentifier()
	azureUrl := "azure://myaccount.blob.core.windows.net/mycontainer"

	modelInvalidAutoRefresh := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithDirectoryEnabledAndOptions(true, "invalid", r.BooleanFalse)

	modelBothEncryptionTypes := model.ExternalAzureStageWithId(id, azureUrl).
		WithEncryptionBothTypes()
	modelEncryptionNoneTypeSpecified := model.ExternalAzureStageWithId(id, azureUrl).
		WithEncryptionNoneTypeSpecified()

	modelBothStorageIntegrationAndCredentials := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithCredentials("invalid")

	modelStorageIntegrationWithPrivatelink := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithUsePrivatelinkEndpoint(r.BooleanTrue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAzureStage),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoRefresh),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected .*refresh_on_create.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelBothEncryptionTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`encryption.0.azure_cse,encryption.0.none.* can be specified`),
			},
			{
				Config:      accconfig.FromModels(t, modelEncryptionNoneTypeSpecified),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("one of `encryption.0.azure_cse,encryption.0.none`"),
			},
			{
				Config:      accconfig.FromModels(t, modelBothStorageIntegrationAndCredentials),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`storage_integration": conflicts with credentials`),
			},
			{
				Config:      accconfig.FromModels(t, modelStorageIntegrationWithPrivatelink),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`storage_integration": conflicts with use_privatelink_endpoint`),
			},
		},
	})
}
