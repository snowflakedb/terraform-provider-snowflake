//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiIntegrations_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	allowedPrefix := "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	awsRoleArn := "arn:aws:iam::000000000001:role/test"
	apiProvider := sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway

	resourceModel := model.ApiIntegrationAmazonApiGateway("test", id.Name(), []string{allowedPrefix}, awsRoleArn, string(apiProvider), true)

	datasourceModel := datasourcemodel.ApiIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())

	datasourceModelWithoutDescribe := datasourcemodel.ApiIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(resourceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel, datasourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "api_integrations.#", "1")),
					resourceshowoutputassert.ApiIntegrationsDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
						HasName(id.Name()).
						HasApiType("EXTERNAL_API").
						HasEnabled(true),
					resourceshowoutputassert.ApiIntegrationsDatasourceDescribeOutput(t, datasourceModel.DatasourceReference()).
						HasEnabled(true).
						HasApiProvider(string(apiProvider)).
						HasApiAwsRoleArn(awsRoleArn).
						HasApiAwsIamUserArnNotEmpty().
						HasApiAwsExternalIdNotEmpty(),
				),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, datasourceModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "api_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "api_integrations.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_ApiIntegrations_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomAccountObjectIdentifier()

	awsRoleArn := "arn:aws:iam::000000000001:role/test"
	allowedPrefix := "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	apiProvider := sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway

	model1 := model.ApiIntegrationAmazonApiGateway("test1", id1.Name(), []string{allowedPrefix}, awsRoleArn, string(apiProvider), true)
	model2 := model.ApiIntegrationAmazonApiGateway("test2", id2.Name(), []string{allowedPrefix}, awsRoleArn, string(apiProvider), true)
	model3 := model.ApiIntegrationAmazonApiGateway("test3", id3.Name(), []string{allowedPrefix}, awsRoleArn, string(apiProvider), true)

	datasourceModelLikeFirst := datasourcemodel.ApiIntegrations("test").
		WithLike(id1.Name()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.ApiIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, datasourceModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeFirst.DatasourceReference(), "api_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "api_integrations.#", "2"),
				),
			},
		},
	})
}
