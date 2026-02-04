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

func TestAcc_ExternalS3Stage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	storageIntegrationId := ids.PrecreatedS3StorageIntegration

	masterKey := random.AwsCseMasterKey()
	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	modelBasic := model.ExternalS3StageWithId(id, awsUrl)
	modelAlter := model.ExternalS3StageWithId(newId, awsUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
		}).
		WithCredentialsAwsKey(awsKeyId, awsSecretKey).
		WithEncryptionAwsCse(masterKey)

	modelComplete := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionAwsCse(masterKey)

	modelUpdated := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          false,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionNone()

	modelEncryptionNoneWithComment := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(changedComment).
		WithEncryptionNone()

	modelRenamed := model.ExternalS3StageWithId(newId, awsUrl).
		WithCredentialsAwsKey(awsKeyId, awsSecretKey)

	modelWithStorageIntegration := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name())

	modelWithCredentials := model.ExternalS3StageWithId(id, awsUrl).
		WithCredentialsAwsKey(awsKeyId, awsSecretKey)

	modelWithUsePrivatelinkEndpoint := model.ExternalS3StageWithId(id, awsUrl).
		WithUsePrivatelinkEndpoint("true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3Stage),
		Steps: []resource.TestStep{
			// Create with empty optionals (basic)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
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
					resourceassert.ExternalS3StageResource(t, modelAlter.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:      true,
							AutoRefresh: sdk.Pointer(r.BooleanDefault),
						}).
						HasEncryptionAwsCse().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelAlter.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
			},
			// Import - after alter
			{
				Config:       accconfig.FromModels(t, modelAlter),
				ResourceName: modelAlter.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "name", newId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "database", newId.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "schema", newId.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "url", awsUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "fully_qualified_name", newId.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "use_privatelink_endpoint", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "stage_type", string(sdk.StageTypeExternal)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "cloud", string(sdk.StageCloudAws)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(newId), "directory.0.auto_refresh", "false"),
				),
			},
			// Set optionals (complete)
			{
				Config: accconfig.FromModels(t, modelComplete),

				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:          true,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasEncryptionAwsCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
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
					resourceassert.ExternalS3StageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:      true,
							AutoRefresh: sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionAwsCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
			},
			// Import - complete
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "url", awsUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "use_privatelink_endpoint", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "storage_integration", storageIntegrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "stage_type", string(sdk.StageTypeExternal)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "cloud", string(sdk.StageCloudAws)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.auto_refresh", "false"),
				),
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
					resourceassert.ExternalS3StageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:      false,
							AutoRefresh: sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, id)()
					testClient().Stage.CreateStageOnS3(t, awsUrl)
				},
				Config: accconfig.FromModels(t, modelUpdated),
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:          false,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(false),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
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
					resourceassert.ExternalS3StageResource(t, modelEncryptionNoneWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEmpty().
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelEncryptionNoneWithComment.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
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
					resourceassert.ExternalS3StageResource(t, modelWithCredentials.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasCredentialsAwsKey(awsKeyId, awsSecretKey).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelWithCredentials.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
			},
			// rename
			{
				Config: accconfig.FromModels(t, modelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
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
					resourceassert.ExternalS3StageResource(t, modelWithStorageIntegration.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
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
					resourceassert.ExternalS3StageResource(t, modelWithStorageIntegration.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckResourceAttr(modelWithStorageIntegration.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelBasic),
			},
			// use privatelink endpoint
			{
				Config: accconfig.FromModels(t, modelWithUsePrivatelinkEndpoint),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithUsePrivatelinkEndpoint.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelWithUsePrivatelinkEndpoint.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasUsePrivatelinkEndpointString(r.BooleanTrue).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelWithUsePrivatelinkEndpoint.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "true")),
				),
			},
			// change privatelink endpoint externally
			{
				Config: accconfig.FromModels(t, modelWithUsePrivatelinkEndpoint),
				PreConfig: func() {
					testClient().Stage.AlterExternalS3Stage(t, sdk.NewAlterExternalS3StageStageRequest(id).WithExternalStageParams(*sdk.NewExternalS3StageParamsRequest(awsUrl).WithUsePrivatelinkEndpoint(false)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithUsePrivatelinkEndpoint.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelWithUsePrivatelinkEndpoint.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasNoStorageIntegration().
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasUsePrivatelinkEndpointString(r.BooleanTrue).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelWithUsePrivatelinkEndpoint.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment("").
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelWithUsePrivatelinkEndpoint.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint", "true")),
				),
			},
		},
	})
}

func TestAcc_ExternalS3Stage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	storageIntegrationId := ids.PrecreatedS3StorageIntegration

	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	masterKey := random.AwsCseMasterKey()

	modelComplete := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionAwsCse(masterKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3Stage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:          true,
							RefreshOnCreate: sdk.Bool(true),
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionAwsCse().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
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
					assert.Check(resource.TestCheckNoResourceAttr(modelComplete.ResourceReference(), "describe_output.0.privatelink.0.use_privatelink_endpoint")),
				),
			},
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "url", awsUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "use_privatelink_endpoint", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "stage_type", string(sdk.StageTypeExternal)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "cloud", string(sdk.StageCloudAws)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.auto_refresh", "false"),
				),
			},
		},
	})
}

func TestAcc_ExternalS3Stage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	storageIntegrationId := testClient().Ids.RandomAccountObjectIdentifier()
	awsUrl := "s3://mybucket/mypath/"

	modelInvalidAutoRefresh := model.ExternalS3StageWithId(id, awsUrl).
		WithInvalidAutoRefresh()
	modelInvalidRefreshOnCreate := model.ExternalS3StageWithId(id, awsUrl).
		WithInvalidRefreshOnCreate()

	modelBothEncryptionTypes := model.ExternalS3StageWithId(id, awsUrl).
		WithEncryptionBothTypes()
	modelEncryptionNoneTypeSpecified := model.ExternalS3StageWithId(id, awsUrl).
		WithEncryptionNoneTypeSpecified()

	modelBothStorageIntegrationAndCredentials := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithCredentialsAwsKey("invalid_key", "invalid_secret")

	modelStorageIntegrationWithPrivatelink := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithUsePrivatelinkEndpoint(r.BooleanTrue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3Stage),
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
				ExpectError: regexp.MustCompile("encryption.0.aws_cse,encryption.0.aws_sse_kms,encryption.0.aws_sse_s3,encryption.0.none`\ncan be specified, but `encryption.0.aws_cse,encryption.0.none` were"),
			},
			{
				Config:      accconfig.FromModels(t, modelEncryptionNoneTypeSpecified),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("encryption.0.aws_cse,encryption.0.aws_sse_kms,encryption.0.aws_sse_s3,encryption.0.none`\nmust be specified"),
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
