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

func TestAcc_UserAuthenticationPolicyAttachment_BasicUseCase(t *testing.T) {
	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)
	userName := user1.ID().Name()

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)
	newUserName := user2.ID().Name()

	authenticationPolicy, authenticationPolicyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(authenticationPolicyCleanup)
	authenticationPolicyName := authenticationPolicy.ID().FullyQualifiedName()

	newAuthenticationPolicy, newAuthenticationPolicyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(newAuthenticationPolicyCleanup)
	newAuthenticationPolicyName := newAuthenticationPolicy.ID().FullyQualifiedName()

	basicModel := model.UserAuthenticationPolicyAttachment("t", authenticationPolicyName, userName)

	newUserModel := model.UserAuthenticationPolicyAttachment("t", authenticationPolicyName, newUserName)

	newPolicyModel := model.UserAuthenticationPolicyAttachment("t", newAuthenticationPolicyName, newUserName)

	ref := basicModel.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserAuthenticationPolicyAttachmentResource(t, ref).
			HasUserName(userName).
			HasAuthenticationPolicyName(authenticationPolicyName),
	}

	newUserAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserAuthenticationPolicyAttachmentResource(t, ref).
			HasUserName(newUserName).
			HasAuthenticationPolicyName(authenticationPolicyName),
	}

	newPolicyAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserAuthenticationPolicyAttachmentResource(t, ref).
			HasUserName(newUserName).
			HasAuthenticationPolicyName(newAuthenticationPolicyName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: warehouseRequiredProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserAuthenticationPolicyAttachmentDestroy(t),
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
			// Change user
			{
				Config: config.FromModels(t, newUserModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, newUserAssertions...),
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
			// Import
			{
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Drop user externally and remove attachment from config - expect empty plan
			{
				PreConfig: user2Cleanup,
				Config:    " ",
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: config.FromModels(t, basicModel),
			},
			// Unset policy externally
			{
				PreConfig: func() {
					testClient().User.Alter(t, sdk.NewAlterUserRequest(user1.ID()).
						WithUnset(*sdk.NewUserUnsetRequest().WithAuthenticationPolicy(true)))
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
