//go:build non_account_level_tests

package testacc

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Stages_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := testClient().Ids.SchemaId()

	comment := random.Comment()

	stageModel1 := model.InternalStage("test1", idOne.DatabaseName(), idOne.SchemaName(), idOne.Name()).
		WithComment(comment)
	stageModel2 := model.InternalStage("test2", idTwo.DatabaseName(), idTwo.SchemaName(), idTwo.Name()).
		WithComment(comment)
	stageModel3 := model.InternalStage("test3", idThree.DatabaseName(), idThree.SchemaName(), idThree.Name()).
		WithComment(comment)

	stagesModelLikeFirst := datasourcemodel.Stages("test").
		WithLike(idOne.Name()).
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(stageModel1.ResourceReference(), stageModel2.ResourceReference(), stageModel3.ResourceReference())

	stagesModelLikePrefix := datasourcemodel.Stages("test").
		WithLike(prefix+"%").
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(stageModel1.ResourceReference(), stageModel2.ResourceReference(), stageModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.InternalStage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, stageModel1, stageModel2, stageModel3, stagesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stagesModelLikeFirst.DatasourceReference(), "stages.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, stageModel1, stageModel2, stageModel3, stagesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stagesModelLikePrefix.DatasourceReference(), "stages.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Stages_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := testClient().Ids.SchemaId()
	comment := random.Comment()
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	endpoint := "s3.us-west-2.amazonaws.com"
	prefix := random.AlphaN(4)
	s3StageId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	azureStageId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	gcsStageId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	s3CompatibleStageId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	fileFormat, fileFormatCleanup := testClient().FileFormat.CreateFileFormat(t)
	t.Cleanup(fileFormatCleanup)

	internalModel := model.InternalStage("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithFileFormatName(fileFormat.ID().FullyQualifiedName()).
		WithComment(comment)

	azureModel := model.ExternalAzureStage("test", azureStageId.DatabaseName(), azureStageId.SchemaName(), azureStageId.Name(), azureBucketUrl).
		WithFileFormatCsv(sdk.FileFormatCsvOptions{
			Compression: sdk.Pointer(sdk.CSVCompressionGzip),
		}).
		WithComment(comment)
	externalS3Model := model.ExternalS3Stage("test", s3StageId.DatabaseName(), s3StageId.SchemaName(), s3StageId.Name(), awsBucketUrl).
		WithFileFormatJson(sdk.FileFormatJsonOptions{
			Compression: sdk.Pointer(sdk.JSONCompressionGzip),
		}).
		WithComment(comment)
	externalGcsModel := model.ExternalGcsStage("test", gcsStageId.DatabaseName(), gcsStageId.SchemaName(), gcsStageId.Name(), ids.PrecreatedGcpStorageIntegration.FullyQualifiedName(), gcsBucketUrl).
		WithFileFormatAvro(sdk.FileFormatAvroOptions{
			Compression: sdk.Pointer(sdk.AvroCompressionGzip),
		}).
		WithComment(comment)
	externalS3CompatibleModel := model.ExternalS3CompatibleStage("test", s3CompatibleStageId.DatabaseName(), s3CompatibleStageId.SchemaName(), s3CompatibleStageId.Name(), endpoint, compatibleBucketUrl).
		WithFileFormatOrc(sdk.FileFormatOrcOptions{
			TrimSpace: sdk.Pointer(true),
		}).
		WithComment(comment)
	externalS3CompatibleModelWithParquet := model.ExternalS3CompatibleStage("test", s3CompatibleStageId.DatabaseName(), s3CompatibleStageId.SchemaName(), s3CompatibleStageId.Name(), endpoint, compatibleBucketUrl).
		WithFileFormatParquet(sdk.FileFormatParquetOptions{
			Compression: sdk.Pointer(sdk.ParquetCompressionLzo),
		}).
		WithComment(comment)
	externalS3CompatibleModelWithAvro := model.ExternalS3CompatibleStage("test", s3CompatibleStageId.DatabaseName(), s3CompatibleStageId.SchemaName(), s3CompatibleStageId.Name(), endpoint, compatibleBucketUrl).
		WithFileFormatAvro(sdk.FileFormatAvroOptions{
			Compression: sdk.Pointer(sdk.AvroCompressionGzip),
		}).
		WithComment(comment)

	stagesNoDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(internalModel.ResourceReference())

	stagesWithDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(true).
		WithDependsOn(internalModel.ResourceReference())

	azureStagesWithDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(true).
		WithDependsOn(azureModel.ResourceReference())
	externalS3StagesWithDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(true).
		WithDependsOn(externalS3Model.ResourceReference())
	externalGcsStagesWithDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(true).
		WithDependsOn(externalGcsModel.ResourceReference())
	externalS3CompatibleStagesWithDescribe := datasourcemodel.Stages("test").
		WithInSchema(schemaId).
		WithWithDescribe(true).
		WithDependsOn(externalS3CompatibleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Internal stage without describe
			{
				Config: accconfig.FromModels(t, internalModel, stagesNoDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, stagesNoDescribe.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrlEmpty().
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegionEmpty().
						HasCloudEmpty().
						HasStorageIntegrationEmpty().
						HasEndpointEmpty().
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeInternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(stagesNoDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(stagesNoDescribe.DatasourceReference(), "stages.0.describe_output.#", "0")),
				),
			},
			// Internal stage with describe
			{
				Config: accconfig.FromModels(t, internalModel, stagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, stagesWithDescribe.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrlEmpty().
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegionEmpty().
						HasCloudEmpty().
						HasStorageIntegrationEmpty().
						HasEndpointEmpty().
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeInternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(stagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", fileFormat.ID().FullyQualifiedName())),
				),
			},
			// Azure stage with describe
			{
				Config: accconfig.FromModels(t, azureModel, azureStagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, azureStagesWithDescribe.DatasourceReference()).
						HasName(azureStageId.Name()).
						HasDatabaseName(azureStageId.DatabaseName()).
						HasSchemaName(azureStageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(azureBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegion("eastus").
						HasCloud(string(sdk.StageCloudAzure)).
						HasStorageIntegrationEmpty().
						HasEndpointEmpty().
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(azureStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
			// External S3 stage with describe
			{
				Config: accconfig.FromModels(t, externalS3Model, externalS3StagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, externalS3StagesWithDescribe.DatasourceReference()).
						HasName(s3StageId.Name()).
						HasDatabaseName(s3StageId.DatabaseName()).
						HasSchemaName(s3StageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(awsBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegion("us-west-2").
						HasCloud(string(sdk.StageCloudAws)).
						HasStorageIntegrationEmpty().
						HasEndpointEmpty().
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.privatelink.0.use_privatelink_endpoint", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.json.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3StagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
			// External GCS stage with describe
			{
				Config: accconfig.FromModels(t, externalGcsModel, externalGcsStagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, externalGcsStagesWithDescribe.DatasourceReference()).
						HasName(gcsStageId.Name()).
						HasDatabaseName(gcsStageId.DatabaseName()).
						HasSchemaName(gcsStageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(gcsBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegionEmpty().
						HasCloud(string(sdk.StageCloudGcp)).
						HasStorageIntegration(ids.PrecreatedGcpStorageIntegration).
						HasEndpointEmpty().
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.avro.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalGcsStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
			// External S3 compatible stage with describe
			{
				Config: accconfig.FromModels(t, externalS3CompatibleModel, externalS3CompatibleStagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, externalS3CompatibleStagesWithDescribe.DatasourceReference()).
						HasName(s3CompatibleStageId.Name()).
						HasDatabaseName(s3CompatibleStageId.DatabaseName()).
						HasSchemaName(s3CompatibleStageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(compatibleBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegion("us-west-2").
						HasCloud(string(sdk.StageCloudAws)).
						HasStorageIntegrationEmpty().
						HasEndpoint(endpoint).
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.orc.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
			// External S3 compatible stage with describe + parquet file format
			{
				Config: accconfig.FromModels(t, externalS3CompatibleModelWithParquet, externalS3CompatibleStagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, externalS3CompatibleStagesWithDescribe.DatasourceReference()).
						HasName(s3CompatibleStageId.Name()).
						HasDatabaseName(s3CompatibleStageId.DatabaseName()).
						HasSchemaName(s3CompatibleStageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(compatibleBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegion("us-west-2").
						HasCloud(string(sdk.StageCloudAws)).
						HasStorageIntegrationEmpty().
						HasEndpoint(endpoint).
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.parquet.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
			// External S3 compatible stage with describe + avro file format
			{
				Config: accconfig.FromModels(t, externalS3CompatibleModelWithAvro, externalS3CompatibleStagesWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StagesDatasourceShowOutput(t, externalS3CompatibleStagesWithDescribe.DatasourceReference()).
						HasName(s3CompatibleStageId.Name()).
						HasDatabaseName(s3CompatibleStageId.DatabaseName()).
						HasSchemaName(s3CompatibleStageId.SchemaName()).
						HasCreatedOnNotEmpty().
						HasUrl(compatibleBucketUrl).
						HasHasCredentials(false).
						HasHasEncryptionKey(false).
						HasRegion("us-west-2").
						HasCloud(string(sdk.StageCloudAws)).
						HasStorageIntegrationEmpty().
						HasEndpoint(endpoint).
						HasOwnerRoleType("ROLE").
						HasDirectoryEnabled(false).
						HasType(sdk.StageTypeExternal).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.enable", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.directory_table.0.auto_refresh", "false")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.csv.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.avro.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(externalS3CompatibleStagesWithDescribe.DatasourceReference(), "stages.0.describe_output.0.file_format.0.format_name", "")),
				),
			},
		},
	})
}

func TestAcc_Stages_MultipleTypes(t *testing.T) {
	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)

	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "1")
	idTwo := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "2")
	schemaId := testClient().Ids.SchemaId()

	internalStage := model.InternalStage("w1", idOne.DatabaseName(), idOne.SchemaName(), idOne.Name())
	externalS3Stage := model.ExternalS3Stage("w2", idTwo.DatabaseName(), idTwo.SchemaName(), idTwo.Name(), awsUrl).
		WithStorageIntegration(ids.PrecreatedS3StorageIntegration.Name())

	stagesModel := datasourcemodel.Stages("test").
		WithLike(prefix+"%").
		WithInSchema(schemaId).
		WithDependsOn(internalStage.ResourceReference(), externalS3Stage.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: ComposeCheckDestroy(t, resources.InternalStage, resources.ExternalS3Stage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, internalStage, externalS3Stage, stagesModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.#", "2"),

					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.show_output.0.name", idOne.Name()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.show_output.0.database_name", idOne.DatabaseName()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.show_output.0.schema_name", idOne.SchemaName()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.show_output.0.type", string(sdk.StageTypeInternal)),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.show_output.0.cloud", ""),

					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.0.describe_output.#", "1"),

					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.show_output.0.name", idTwo.Name()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.show_output.0.database_name", idTwo.DatabaseName()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.show_output.0.schema_name", idTwo.SchemaName()),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.show_output.0.type", string(sdk.StageTypeExternal)),
					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.show_output.0.cloud", string(sdk.StageCloudAws)),

					resource.TestCheckResourceAttr(stagesModel.DatasourceReference(), "stages.1.describe_output.#", "1"),
				),
			},
		},
	})
}
