//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_ExternalS3Stage_DirectoryTableTimestampParsing verifies that reading an external S3 stage
// with a directory table (which populates LAST_REFRESHED_ON) does not fail.
// This is a regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/4445.
func TestAcc_ExternalS3Stage_DirectoryTableTimestampParsing(t *testing.T) {
	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifier()
	awsUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	storageIntegrationId := ids.PrecreatedS3StorageIntegration

	revertParameter := secondaryTestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterTimestampOutputFormat, "YYYY-MM-DD HH24:MI:SS")
	t.Cleanup(revertParameter)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.ExternalS3StageResource))

	stageModel := model.ExternalS3StageWithId(id, awsUrl).
		WithStorageIntegration(storageIntegrationId.Name()).
		WithDirectoryEnabledAndOptions(sdk.StageS3CommonDirectoryTableOptionsRequest{
			Enable:          true,
			RefreshOnCreate: sdk.Bool(true),
		})

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:            accconfig.FromModels(t, providerModel, stageModel),
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				// Not matching the whole error message because in the issue the UTC time was used.
				ExpectError: regexp.MustCompile("Error: parsing time"),
			},
			{
				Config: accconfig.FromModels(t, providerModel, stageModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// This happens because the resource is marked as tainted. Unfortunately, in the testing framework we cannot assert this.
						plancheck.ExpectResourceAction(stageModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: s3StageProviderFactory,
				Check: assertThat(t,
					resourceassert.ExternalS3StageResource(t, stageModel.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUrlString(awsUrl).
						HasStorageIntegrationString(storageIntegrationId.Name()).
						HasDirectory(resourceassert.ExternalS3StageDirectoryTableAssert{
							Enable:          true,
							AutoRefresh:     sdk.Pointer(r.BooleanDefault),
							RefreshOnCreate: sdk.Bool(true),
						}).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCloudEnum(sdk.StageCloudAws).
						HasStageTypeEnum(sdk.StageTypeExternal),
					resourceshowoutputassert.StageShowOutput(t, stageModel.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.StageTypeExternal).
						HasDirectoryEnabled(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(stageModel.ResourceReference(), "describe_output.0.directory_table.0.enable", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(stageModel.ResourceReference(), "describe_output.0.directory_table.0.last_refreshed_on")),
				),
			},
		},
	})
}
