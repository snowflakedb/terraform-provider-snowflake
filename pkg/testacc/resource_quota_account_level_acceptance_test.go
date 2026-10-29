//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// An account can have only one Snowflake Intelligence object, so this can't run in parallel with other tests creating one.
func TestAcc_Quota_SnowflakeIntelligenceInstance(t *testing.T) {
	intelligenceId, intelligenceCleanup := testClient().SnowflakeIntelligence.Create(t)
	t.Cleanup(intelligenceCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	sharedResource := model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainSnowflakeIntelligence, Instance: intelligenceId.Name()}
	quotaModel := model.QuotaFromId("test", id).WithSharedResources(sharedResource)
	ref := quotaModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Quota),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, quotaModel),
				Check: assertThat(t,
					resourceassert.QuotaResource(t, ref).
						HasSharedResourcesCount(1).
						HasSharedResource(sharedResource.Domain, sharedResource.Instance),
				),
			},
			{
				Config: accconfig.FromModels(t, quotaModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config:            accconfig.FromModels(t, quotaModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
