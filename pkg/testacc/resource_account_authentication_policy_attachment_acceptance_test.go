//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountAuthenticationPolicyAttachment_BasicUseCase(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	policy, policyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(policyCleanup)
	policyName := policy.ID().FullyQualifiedName()

	newPolicy, newPolicyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(newPolicyCleanup)
	newPolicyName := newPolicy.ID().FullyQualifiedName()

	basicModel := model.AccountAuthenticationPolicyAttachment("t", policyName)

	newPolicyModel := model.AccountAuthenticationPolicyAttachment("t", newPolicyName)

	ref := basicModel.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, ref).
			HasAuthenticationPolicy(policyName),
	}

	newPolicyAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, ref).
			HasAuthenticationPolicy(newPolicyName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: warehouseRequiredProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
			// Import
			{
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change policy
			{
				Config: config.FromModels(t, newPolicyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, newPolicyAssertions...),
			},
			// Unset policy externally
			{
				PreConfig: func() {
					testClient().Account.Alter(t, sdk.NewAlterAccountRequest().
						WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true))))
				},
				Config: config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
		},
	})
}
