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

func TestAcc_ExternalAccessIntegrations_BasicUseCase(t *testing.T) {
	networkRuleId := testClient().Ids.RandomSchemaObjectIdentifier()
	_, networkRuleCleanup := testClient().NetworkRule.CreateEgressWithIdentifier(t, networkRuleId)
	t.Cleanup(networkRuleCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceModel := model.ExternalAccessIntegrationFromId(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true)

	datasourceModel := datasourcemodel.ExternalAccessIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())

	datasourceModelWithoutDescribe := datasourcemodel.ExternalAccessIntegrations("test").
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
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "external_access_integrations.#", "1")),
					resourceshowoutputassert.ExternalAccessIntegrationsDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasEnabled(true),
					resourceshowoutputassert.ExternalAccessIntegrationsDatasourceDescribeOutput(t, datasourceModel.DatasourceReference()).
						HasId(id).
						HasEnabled(true).
						HasNoAllowedAuthenticationSecrets().
						HasNoAllowedApiAuthenticationIntegrations(),
				),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, datasourceModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "external_access_integrations.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_ExternalAccessIntegrations_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)

	networkRuleId := testClient().Ids.RandomSchemaObjectIdentifier()
	_, networkRuleCleanup := testClient().NetworkRule.CreateEgressWithIdentifier(t, networkRuleId)
	t.Cleanup(networkRuleCleanup)

	id1 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomAccountObjectIdentifier()

	nrFQN := networkRuleId.FullyQualifiedName()
	resourceModel1 := model.ExternalAccessIntegration("test1", id1.Name(), []string{nrFQN}, true)
	resourceModel2 := model.ExternalAccessIntegration("test2", id2.Name(), []string{nrFQN}, true)
	resourceModel3 := model.ExternalAccessIntegration("test3", id3.Name(), []string{nrFQN}, true)

	datasourceModelLikeFirst := datasourcemodel.ExternalAccessIntegrations("test").
		WithLike(id1.Name()).
		WithDependsOn(resourceModel1.ResourceReference(), resourceModel2.ResourceReference(), resourceModel3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.ExternalAccessIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(resourceModel1.ResourceReference(), resourceModel2.ResourceReference(), resourceModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel1, resourceModel2, resourceModel3, datasourceModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeFirst.DatasourceReference(), "external_access_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, resourceModel1, resourceModel2, resourceModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "external_access_integrations.#", "2"),
				),
			},
		},
	})
}
