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

// TODO(SNOW-2356128): Test notification integration
func TestAcc_ExternalGcsStage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	storageIntegrationId := ids.PrecreatedGcpStorageIntegration

	gcsUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	s3Url := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)

	modelBasic := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl)
	modelAlter := model.ExternalGcsStageWithId(newId, storageIntegrationId.Name(), gcsUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalGCSDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
		}).
		WithEncryptionNone()

	modelComplete := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalGCSDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionNone()
	modelCompleteWithoutAutoRefresh := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalGCSDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
		}).
		WithEncryptionNone()

	modelUpdated := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithComment(changedComment).
		WithDirectoryEnabledAndOptions(sdk.ExternalGCSDirectoryTableOptionsRequest{
			Enable:          false,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionNone()

	modelEncryptionNoneWithComment := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithComment(changedComment).
		WithEncryptionNone()

	modelRenamed := model.ExternalGcsStageWithId(newId, storageIntegrationId.Name(), gcsUrl)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalGcsStage),
		Steps: []resource.TestStep{
			// Create with required fields only (basic)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
				ImportStateVerifyIgnore: []string{"directory", "file_format"},
			},
			// Alter - rename
			{
				Config: accconfig.FromModels(t, modelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString("").
						HasDirectoryEmpty().
						HasEncryptionEmpty().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
			// alter
			{
				Config: accconfig.FromModels(t, modelAlter),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelAlter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelAlter.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:      true,
							AutoRefresh: sdk.Pointer(r.BooleanDefault),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
				ImportStateVerifyIgnore: []string{"directory", "encryption"},
			},
			// Set optionals (complete)
			{
				Config: accconfig.FromModels(t, modelComplete),

				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          true,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
				ImportStateVerifyIgnore: []string{"directory", "encryption"},
			},
			// unset auto_refresh
			{
				Config: accconfig.FromModels(t, modelCompleteWithoutAutoRefresh),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithoutAutoRefresh.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelCompleteWithoutAutoRefresh.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          true,
							RefreshOnCreate: sdk.Bool(true),
							AutoRefresh:     sdk.Pointer(r.BooleanDefault),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithoutAutoRefresh.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithoutAutoRefresh.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
				),
			},
			// Bring back the complete model
			{
				Config: accconfig.FromModels(t, modelComplete),
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
					resourceassert.ExternalGcsStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          false,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
					testClient().Stage.CreateStageOnGCSWithId(t, id, gcsUrl)
				},
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          false,
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
					resourceassert.ExternalGcsStageResource(t, modelEncryptionNoneWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEmpty().
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
			// Detect changing stage type externally
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, id)()
					testClient().Stage.CreateStageOnS3WithRequest(t,
						sdk.NewCreateOnS3StageRequest(id,
							*sdk.NewExternalS3StageParamsRequest(s3Url)))
				},
				Config: accconfig.FromModels(t, modelEncryptionNoneWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelEncryptionNoneWithComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelEncryptionNoneWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(changedComment).
						HasDirectoryEmpty().
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
		},
	})
}

func TestAcc_ExternalGcsStage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	storageIntegrationId := ids.PrecreatedGcpStorageIntegration

	gcsUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)

	modelComplete := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.ExternalGCSDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
			AutoRefresh:     sdk.Bool(false),
		}).
		WithEncryptionNone()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalGcsStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(gcsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasCommentString(comment).
						HasDirectory(resourceassert.ExternalStageDirectoryTableAssert{
							Enable:          true,
							RefreshOnCreate: sdk.Pointer(true),
							AutoRefresh:     sdk.Pointer(r.BooleanFalse),
						}).
						HasEncryptionNone().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudGcp).
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
				ImportStateVerifyIgnore: []string{"directory", "encryption"},
			},
		},
	})
}

func TestAcc_ExternalGcsStage_FileFormat_SwitchBetweenTypes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	storageIntegrationId := ids.PrecreatedGcpStorageIntegration
	gcsUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)

	fileFormat, fileFormatCleanup := testClient().FileFormat.CreateFileFormat(t)
	t.Cleanup(fileFormatCleanup)

	modelBasic := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl)

	modelWithCsvFormat := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{})

	modelWithNamedFormat := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithFileFormatName(fileFormat.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalGcsStage),
		Steps: []resource.TestStep{
			// Start with inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelWithCsvFormat.ResourceReference()).
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
					resourceassert.ExternalGcsStageResource(t, modelWithNamedFormat.ResourceReference()).
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
				ImportStateVerifyIgnore: []string{"encryption", "directory"},
			},
			// Detect external change
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				PreConfig: func() {
					testClient().Stage.AlterExternalGCSStage(t, sdk.NewAlterExternalGCSStageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalGcsStageResource(t, modelWithNamedFormat.ResourceReference()).
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
					resourceassert.ExternalGcsStageResource(t, modelWithCsvFormat.ResourceReference()).
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
					resourceassert.ExternalGcsStageResource(t, modelBasic.ResourceReference()).
						HasFileFormatEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
		},
	})
}

func TestAcc_ExternalGcsStage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	storageIntegrationId := testClient().Ids.RandomAccountObjectIdentifier()
	gcsUrl := "gcs://mybucket/mypath/"

	modelInvalidAutoRefresh := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithInvalidAutoRefresh()
	modelInvalidRefreshOnCreate := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithInvalidRefreshOnCreate()

	modelBothEncryptionTypes := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithEncryptionBothTypes()
	modelEncryptionNoneTypeSpecified := model.ExternalGcsStageWithId(id, storageIntegrationId.Name(), gcsUrl).
		WithEncryptionNoneTypeSpecified()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalGcsStage),
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
				ExpectError: regexp.MustCompile(`encryption.0.gcs_sse_kms,encryption.0.none.* can be specified`),
			},
			{
				Config:      accconfig.FromModels(t, modelEncryptionNoneTypeSpecified),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("one of `encryption.0.gcs_sse_kms,encryption.0.none`"),
			},
		},
	})
}
