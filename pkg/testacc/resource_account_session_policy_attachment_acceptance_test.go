//go:build non_account_level_tests

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

func TestAcc_AccountSessionPolicyAttachment_BasicUseCase(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	sessionPolicy, sessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup)
	sessionPolicyName := sessionPolicy.ID().FullyQualifiedName()

	sessionPolicy2, sessionPolicyCleanup2 := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup2)
	sessionPolicyName2 := sessionPolicy2.ID().FullyQualifiedName()

	basic := model.AccountSessionPolicyAttachment("t", sessionPolicyName)

	newPolicy := model.AccountSessionPolicyAttachment("t", sessionPolicyName2)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountSessionPolicyAttachmentResource(t, ref).
			HasSessionPolicyName(sessionPolicyName),
	}

	newPolicyAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountSessionPolicyAttachmentResource(t, ref).
			HasSessionPolicyName(sessionPolicyName2),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change policy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, newPolicy),
				Check:  assertThat(t, newPolicyAssertions...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, newPolicy),
				Destroy: true,
			},
			{
				Config: config.FromModels(t, basic),
			},
			// Unset policy externally
			{
				PreConfig: func() {
					testClient().Account.Alter(t, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{SessionPolicy: sdk.Bool(true)}})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
		},
	})
}
