//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StorageIntegrations_BasicUseCase_DifferentFiltering(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	allowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
	}

	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	awsModel1 := model.StorageIntegrationAws("test1", idOne.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))
	awsModel2 := model.StorageIntegrationAws("test2", idTwo.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))
	awsModel3 := model.StorageIntegrationAws("test3", idThree.Name(), false, allowedLocations, awsRoleArn, string(sdk.RegularS3Protocol))

	storageIntegrationsModelLikeFirst := datasourcemodel.StorageIntegrations("test").
		WithWithDescribe(false).
		WithLike(idOne.Name()).
		WithDependsOn(awsModel1.ResourceReference(), awsModel2.ResourceReference(), awsModel3.ResourceReference())

	storageIntegrationsModelLikePrefix := datasourcemodel.StorageIntegrations("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithDependsOn(awsModel1.ResourceReference(), awsModel2.ResourceReference(), awsModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, awsModel1, awsModel2, awsModel3, storageIntegrationsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(storageIntegrationsModelLikeFirst.DatasourceReference(), "storage_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, awsModel1, awsModel2, awsModel3, storageIntegrationsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(storageIntegrationsModelLikePrefix.DatasourceReference(), "storage_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_StorageIntegrations_CompleteUseCase(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	awsAllowedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "allowed-location/"},
		{Path: awsBucketUrl + "allowed-location2/"},
	}
	awsBlockedLocations := []sdk.StorageLocation{
		{Path: awsBucketUrl + "blocked-location/"},
		{Path: awsBucketUrl + "blocked-location2/"},
	}
	azureAllowedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "allowed-location/"},
		{Path: azureBucketUrl + "allowed-location2/"},
	}
	azureBlockedLocations := []sdk.StorageLocation{
		{Path: azureBucketUrl + "blocked-location/"},
		{Path: azureBucketUrl + "blocked-location2/"},
	}
	gcsAllowedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "allowed-location/"},
		{Path: gcsBucketUrl + "allowed-location2/"},
	}
	gcsBlockedLocations := []sdk.StorageLocation{
		{Path: gcsBucketUrl + "blocked-location/"},
		{Path: gcsBucketUrl + "blocked-location2/"},
	}

	prefix := random.AlphaN(4)
	awsIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	azureIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	gcsIntegrationId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	comment := random.Comment()

	awsExternalId := "some_external_id"

	storageIntegrationAws := model.StorageIntegrationAws("w", awsIntegrationId.Name(), false, awsAllowedLocations, awsRoleArn, string(sdk.RegularS3Protocol)).
		WithStorageBlockedLocations(awsBlockedLocations).
		WithComment(comment).
		WithStorageAwsExternalId(awsExternalId).
		WithStorageAwsObjectAcl("bucket-owner-full-control")

	storageIntegrationAzure := model.StorageIntegrationAzure("w", azureIntegrationId.Name(), azureTenantId, false, azureAllowedLocations).
		WithStorageBlockedLocations(azureBlockedLocations).
		WithComment(comment)

	storageIntegrationGcs := model.StorageIntegrationGcs("w", gcsIntegrationId.Name(), false, gcsAllowedLocations).
		WithStorageBlockedLocations(gcsBlockedLocations).
		WithComment(comment)

	awsNoDescribe := datasourcemodel.StorageIntegrations("test").
		WithLike(awsIntegrationId.Name()).
		WithWithDescribe(false).
		WithDependsOn(storageIntegrationAws.ResourceReference())

	awsWithDescribe := datasourcemodel.StorageIntegrations("test").
		WithLike(awsIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(storageIntegrationAws.ResourceReference())

	azureWithDescribe := datasourcemodel.StorageIntegrations("test").
		WithLike(azureIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(storageIntegrationAzure.ResourceReference())

	gcsWithDescribe := datasourcemodel.StorageIntegrations("test").
		WithLike(gcsIntegrationId.Name()).
		WithWithDescribe(true).
		WithDependsOn(storageIntegrationGcs.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegrationAws),
		Steps: []resource.TestStep{
			// AWS without describe
			{
				Config: accconfig.FromModels(t, storageIntegrationAws, awsNoDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StorageIntegrationsDatasourceShowOutput(t, awsNoDescribe.DatasourceReference()).
						HasName(awsIntegrationId.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.0.describe_output.#", "0")),
				),
			},
			// AWS with describe
			{
				Config: accconfig.FromModels(t, storageIntegrationAws, awsWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StorageIntegrationsDatasourceShowOutput(t, awsWithDescribe.DatasourceReference()).
						HasName(awsIntegrationId.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.0.describe_output.#", "1")),
					// TODO [this PR]: add all checks
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(awsNoDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.comment", comment)),
				),
			},
			// Azure with describe
			{
				Config: accconfig.FromModels(t, storageIntegrationAzure, azureWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StorageIntegrationsDatasourceShowOutput(t, azureWithDescribe.DatasourceReference()).
						HasName(azureIntegrationId.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					assert.Check(resource.TestCheckResourceAttr(azureWithDescribe.DatasourceReference(), "storage_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(azureWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.#", "1")),
					// TODO [this PR]: add all checks
					assert.Check(resource.TestCheckResourceAttr(azureWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(azureWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.comment", comment)),
				),
			},
			// GCS with describe
			{
				Config: accconfig.FromModels(t, storageIntegrationGcs, gcsWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.StorageIntegrationsDatasourceShowOutput(t, gcsWithDescribe.DatasourceReference()).
						HasName(gcsIntegrationId.Name()).
						HasEnabled(false).
						HasComment(comment).
						HasStorageType("EXTERNAL_STAGE").
						HasCategory("STORAGE"),
					assert.Check(resource.TestCheckResourceAttr(gcsWithDescribe.DatasourceReference(), "storage_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(gcsWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.#", "1")),
					// TODO [this PR]: add all checks
					assert.Check(resource.TestCheckResourceAttr(gcsWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(gcsWithDescribe.DatasourceReference(), "storage_integrations.0.describe_output.0.comment", comment)),
				),
			},
		},
	})
}
