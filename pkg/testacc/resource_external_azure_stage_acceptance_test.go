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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-2356128): Test use_privatelink_endpoint and notification integration
func TestAcc_ExternalAzureStage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	storageIntegrationId := ids.PrecreatedAzureStorageIntegration

	masterKey := random.AzureCseMasterKey()
	azureUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureSasToken := testenvs.GetOrSkipTest(t, testenvs.AzureExternalSasToken)

	modelBasic := model.ExternalAzureStageWithId(id, azureUrl)
	modelAlter := model.ExternalAzureStageWithId(newId, azureUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalAzureDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
		}).
		WithCredentials(azureSasToken).
		WithEncryptionAzureCse(masterKey)

	modelComplete := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalAzureDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionAzureCse(masterKey)

	modelUpdated := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithDirectoryEnabledAndOptions(sdk.ExternalAzureDirectoryTableOptionsRequest{
			Enable:          false,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionNone()

	modelEncryptionNoneWithComment := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithEncryptionNone()

	modelRenamed := model.ExternalAzureStageWithId(newId, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name())

	modelWithStorageIntegration := model.ExternalAzureStageWithId(id, azureUrl).
		WithStorageIntegration(storageIntegrationId.Name())

	modelWithCredentials := model.ExternalAzureStageWithId(id, azureUrl).
		WithCredentials(azureSasToken)

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
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Import - without optionals
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            modelBasic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"directory", "file_format", "use_privatelink_endpoint"},
			},

			// alter
			{
				Config: accconfig.FromModels(t, modelAlter),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelAlter.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(azureUrl).
						HasNoStorageIntegration().
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:      true,
							AutoRefresh: sdk.Pointer(r.BooleanDefault),
						}).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelAlter.ResourceReference()).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelAlter.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelAlter.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Import - after alter
			{
				Config:                  accconfig.FromModels(t, modelAlter),
				ResourceName:            modelAlter.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials", "encryption", "use_privatelink_endpoint", "directory"},
			},
			// Set optionals (complete)
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          true,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// External change - disable directory
			{
				Config: accconfig.FromModels(t, modelComplete),
				PreConfig: func() {
					testClient().Stage.AlterDirectoryTable(t, sdk.NewAlterDirectoryTableStageRequest(id).WithSetDirectory(sdk.DirectoryTableSetRequest{
						Enable: false,
					}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:      true,
							AutoRefresh: sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Import - complete
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "use_privatelink_endpoint"},
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
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:      false,
							AutoRefresh: sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, id)()
					testClient().Stage.CreateStageOnAzure(t, azureUrl)
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
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          false,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(false),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
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
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelEncryptionNoneWithComment.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(changedComment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelEncryptionNoneWithComment.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelEncryptionNoneWithComment.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// ForceNew - unset encryption and rename
			{
				Config: accconfig.FromModels(t, modelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelRenamed.ResourceReference()).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// set credentials
			{
				Config: accconfig.FromModels(t, modelWithCredentials),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithCredentials.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithCredentials.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasCredentials(azureSasToken).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelWithCredentials.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelWithCredentials.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCredentials.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// set back to storage integration
			{
				Config: accconfig.FromModels(t, modelWithStorageIntegration),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithStorageIntegration.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithStorageIntegration.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelWithStorageIntegration.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasStorageIntegration(storageIntegrationId).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Detect changing stage type externally
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, id)()
					testClient().Stage.CreateInternalStageWithId(t, id)
				},
				Config: accconfig.FromModels(t, modelWithStorageIntegration),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithStorageIntegration.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithStorageIntegration.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelWithStorageIntegration.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasStorageIntegration(storageIntegrationId).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// unset storage integration
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(azureUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
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
		WithDirectoryEnabledAndOptions(sdk.ExternalAzureDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
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
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          true,
							RefreshOnCreate: sdk.Bool(true),
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionAzureCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAzure).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials", "encryption", "use_privatelink_endpoint", "directory"},
			},
		},
	})
}

func TestAcc_ExternalAzureStage_FileFormat_SwitchBetweenTypes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	azureUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)

	fileFormat, fileFormatCleanup := testClient().FileFormat.CreateFileFormat(t)
	t.Cleanup(fileFormatCleanup)

	modelBasic := model.ExternalAzureStageWithId(id, azureUrl)

	modelWithCsvFormat := model.ExternalAzureStageWithId(id, azureUrl).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{})

	modelWithNamedFormat := model.ExternalAzureStageWithId(id, azureUrl).
		WithFileFormatName(fileFormat.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalAzureStage),
		Steps: []resource.TestStep{
			// Start with inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithCsvFormat.ResourceReference()).
						HasFileFormatCsv(),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
			// Switch to named format
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithNamedFormat.ResourceReference()).
						HasFileFormatFormatName(fileFormat.ID().FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", fileFormat.ID().FullyQualifiedName())),
				),
			},
			// import named format
			{
				Config:                  accconfig.FromModels(t, modelWithNamedFormat),
				ResourceName:            modelWithNamedFormat.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption", "directory", "credentials", "use_privatelink_endpoint"},
			},
			// Detect external change
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				PreConfig: func() {
					testClient().Stage.AlterExternalAzureStage(t, sdk.NewAlterExternalAzureStageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithNamedFormat.ResourceReference()).
						HasFileFormatFormatName(fileFormat.ID().FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelWithNamedFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", fileFormat.ID().FullyQualifiedName())),
				),
			},
			// Switch back to inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithCsvFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelWithCsvFormat.ResourceReference()).
						HasFileFormatCsv(),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithCsvFormat.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
			// Switch back to default
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalAzureStageResource(t, modelBasic.ResourceReference()).
						HasFileFormatEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
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
		WithInvalidAutoRefresh()
	modelInvalidRefreshOnCreate := model.ExternalAzureStageWithId(id, azureUrl).
		WithInvalidRefreshOnCreate()

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
				ExpectError: regexp.MustCompile(`expected .*auto_refresh.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidRefreshOnCreate),
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
