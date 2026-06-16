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

func TestAcc_UserSessionPolicyAttachment_BasicUseCase(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)
	userName := user.ID().Name()

	user2, userCleanup2 := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup2)
	userName2 := user2.ID().Name()

	sessionPolicy, sessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup)
	sessionPolicyName := sessionPolicy.ID().FullyQualifiedName()

	sessionPolicy2, sessionPolicyCleanup2 := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup2)
	sessionPolicyName2 := sessionPolicy2.ID().FullyQualifiedName()

	basic := model.UserSessionPolicyAttachment("t", sessionPolicyName, userName)

	newUser := model.UserSessionPolicyAttachment("t", sessionPolicyName, userName2)

	newPolicy := model.UserSessionPolicyAttachment("t", sessionPolicyName2, userName2)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserSessionPolicyAttachmentResource(t, ref).
			HasUserName(userName).
			HasSessionPolicyName(sessionPolicyName),
	}

	newUserAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserSessionPolicyAttachmentResource(t, ref).
			HasUserName(userName2).
			HasSessionPolicyName(sessionPolicyName),
	}

	newPolicyAssertions := []assert.TestCheckFuncProvider{
		resourceassert.UserSessionPolicyAttachmentResource(t, ref).
			HasUserName(userName2).
			HasSessionPolicyName(sessionPolicyName2),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: activeWarehouseSetOnUserProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserSessionPolicyAttachmentDestroy(t),
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
			// Change user
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, newUser),
				Check:  assertThat(t, newUserAssertions...),
			},
			// Change policy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
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
			// Drop user externally and remove attachment from config - expect empty plan
			{
				PreConfig: func() {
					testClient().User.DropUserFunc(t, user.ID())()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: " ",
			},
			{
				Config: config.FromModels(t, newUser),
			},
			// Unset policy externally
			{
				PreConfig: func() {
					testClient().User.Alter(t, user2.ID(), &sdk.AlterUserOptions{Unset: &sdk.UserUnset{SessionPolicy: sdk.Bool(true)}})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, newUser),
				Check:  assertThat(t, newUserAssertions...),
			},
		},
	})
}
