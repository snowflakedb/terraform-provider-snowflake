//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalS3CompatStage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	modelBasic := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint)

	modelComplete := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		})

	modelUpdated := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          false,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		})

	modelNoDirectory := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithComment(comment)

	modelRenamed := model.ExternalS3CompatibleStageWithId(newId, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3CompatibleStage),
		Steps: []resource.TestStep{
			// Create with basic credentials
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString("").
						HasDirectoryEmpty().
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
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
			// Set complete options
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString(comment).
						HasDirectory(sdk.StageS3CommonDirectoryTableOptionsRequest{
							Enable:          true,
							AutoRefresh:     sdk.Bool(false),
							RefreshOnCreate: sdk.Bool(false),
						}).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
			// Import - after complete
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "url", s3CompatUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "endpoint", s3CompatEndpoint),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "stage_type", string(sdk.StageTypeExternal)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "cloud", string(sdk.StageCloudAws)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.auto_refresh", "false"),
				),
			},
			// Update directory
			{
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString(comment).
						HasDirectory(sdk.StageS3CommonDirectoryTableOptionsRequest{
							Enable:          false,
							AutoRefresh:     sdk.Bool(false),
							RefreshOnCreate: sdk.Bool(false),
						}).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelUpdated.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, id)()
					testClient().Stage.CreateStageOnS3CompatibleWithId(t, id, s3CompatUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
				},
				Config: accconfig.FromModels(t, modelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelUpdated.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString(comment).
						HasDirectory(sdk.StageS3CommonDirectoryTableOptionsRequest{
							Enable:          false,
							AutoRefresh:     sdk.Bool(false),
							RefreshOnCreate: sdk.Bool(false),
						}).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelUpdated.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
			// ForceNew - unset directory
			{
				Config: accconfig.FromModels(t, modelNoDirectory),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelNoDirectory.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelNoDirectory.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString(comment).
						HasDirectoryEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, modelNoDirectory.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasDirectoryEnabled(false).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelNoDirectory.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelNoDirectory.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
			// Rename
			{
				Config: accconfig.FromModels(t, modelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString("").
						HasDirectoryEmpty().
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
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
			// Detect changing stage type externally
			{
				PreConfig: func() {
					testClient().Stage.DropStageFunc(t, newId)()
					testClient().Stage.CreateStageOnS3WithId(t, newId, awsBucketUrl)
				},
				Config: accconfig.FromModels(t, modelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(modelRenamed.ResourceReference(), "url", sdk.Pointer(s3CompatUrl), sdk.Pointer(awsBucketUrl)),
						planchecks.ExpectDrift(modelRenamed.ResourceReference(), "endpoint", sdk.Pointer(s3CompatEndpoint), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString("").
						HasDirectoryEmpty().
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelRenamed.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
				),
			},
		},
	})
}

func TestAcc_ExternalS3CompatStage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	modelComplete := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithComment(comment).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(false),
			AutoRefresh:     sdk.Bool(false),
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3CompatibleStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(s3CompatUrl).
						HasEndpointString(s3CompatEndpoint).
						HasCommentString(comment).
						HasDirectory(sdk.StageS3CommonDirectoryTableOptionsRequest{
							Enable:          true,
							RefreshOnCreate: sdk.Bool(false),
							AutoRefresh:     sdk.Bool(false),
						}).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.location.0.url.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.location.0.url.0", s3CompatUrl)),
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
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "url", s3CompatUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "stage_type", string(sdk.StageTypeExternal)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "cloud", string(sdk.StageCloudAws)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.enable", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "directory.0.auto_refresh", "false"),
				),
			},
		},
	})
}

func TestAcc_ExternalS3CompatStage_FileFormat_SwitchBetweenTypes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	fileFormat, fileFormatCleanup := testClient().FileFormat.CreateFileFormat(t)
	t.Cleanup(fileFormatCleanup)

	modelBasic := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey)

	modelWithCsvFormat := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{})

	modelWithNamedFormat := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithFileFormatName(fileFormat.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3CompatibleStage),
		Steps: []resource.TestStep{
			// Start with inline CSV
			{
				Config: accconfig.FromModels(t, modelWithCsvFormat),
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelWithCsvFormat.ResourceReference()).
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
					resourceassert.ExternalS3CompatStageResource(t, modelWithNamedFormat.ResourceReference()).
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
				ImportStateVerifyIgnore: []string{"directory", "credentials"},
			},
			// Detect external change
			{
				Config: accconfig.FromModels(t, modelWithNamedFormat),
				PreConfig: func() {
					testClient().Stage.AlterExternalS3Stage(t, sdk.NewAlterExternalS3StageStageRequest(id).WithFileFormat(sdk.StageFileFormatRequest{FileFormatOptions: &sdk.FileFormatOptions{CsvOptions: &sdk.FileFormatCsvOptions{}}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNamedFormat.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ExternalS3CompatStageResource(t, modelWithNamedFormat.ResourceReference()).
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
					resourceassert.ExternalS3CompatStageResource(t, modelWithCsvFormat.ResourceReference()).
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
					resourceassert.ExternalS3CompatStageResource(t, modelBasic.ResourceReference()).
						HasFileFormatEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.file_format.0.format_name", "")),
				),
			},
		},
	})
}

func TestAcc_ExternalS3CompatStage_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	modelInvalidAutoRefresh := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithInvalidAutoRefresh()

	modelInvalidRefreshOnCreate := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithCredentials(awsKeyId, awsSecretKey).
		WithInvalidRefreshOnCreate()
	modelInvalidWithEmptyCredentials := model.ExternalS3CompatibleStageWithId(id, s3CompatUrl, s3CompatEndpoint).
		WithEmptyCredentials()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalS3CompatibleStage),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoRefresh),
				ExpectError: regexp.MustCompile(`expected .*auto_refresh.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidRefreshOnCreate),
				ExpectError: regexp.MustCompile(`expected .*refresh_on_create.* to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidWithEmptyCredentials),
				ExpectError: regexp.MustCompile(`The argument "aws_secret_key" is required, but no definition was found.`),
			},
		},
	})
}
