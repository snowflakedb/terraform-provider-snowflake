//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

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
