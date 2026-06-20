//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_GrantAccountRole_ShowCaching_RoleToRole verifies that granting a role to a parent role
// works correctly with the GRANT_ACCOUNT_ROLE_SHOW_CACHING experiment enabled. With the experiment,
// the trailing Read at the end of Create is skipped, so this test also asserts (via an empty plan on
// the second step) that the resulting state is correct and drift-free on the subsequent refresh.
func TestAcc_GrantAccountRole_ShowCaching_RoleToRole(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantAccountRoleShowCaching)
	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithParentRoleName(parentRole.ID().Name())

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantAccountRoleShowCachingProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRole.ID().Name()),
				),
			},
			// re-apply the same config: the refresh (Read) must observe the grant created above and
			// produce an empty plan, proving the skipped trailing Read on Create did not corrupt state.
			{
				Config: config.FromModels(t, experimentProviderModel, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// import
			{
				Config:            config.FromModels(t, experimentProviderModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"%s"|%s|"%s"`, role.ID().Name(), sdk.ObjectTypeRole, parentRole.ID().Name()),
			},
		},
	})
}

// TestAcc_GrantAccountRole_ShowCaching_RoleToUser verifies the USER object-type path works correctly
// with the GRANT_ACCOUNT_ROLE_SHOW_CACHING experiment enabled.
func TestAcc_GrantAccountRole_ShowCaching_RoleToUser(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantAccountRoleShowCaching)
	grantModel := model.GrantAccountRole("test", role.ID().Name()).
		WithUserName(user.ID().Name())

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantAccountRoleShowCachingProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "user_name", user.ID().Name()),
				),
			},
			// import
			{
				Config:            config.FromModels(t, experimentProviderModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"%s"|%s|"%s"`, role.ID().Name(), sdk.ObjectTypeUser, user.ID().Name()),
			},
		},
	})
}

// TestAcc_GrantAccountRole_ShowCaching_SharedRoleMultipleGrantees exercises the actual cache scenario:
// the same role is granted to several users within a single configuration, so one plan performs
// multiple Reads of the same role and hits the cache for all but the first. All grants must be created
// correctly and the configuration must converge to an empty plan.
func TestAcc_GrantAccountRole_ShowCaching_SharedRoleMultipleGrantees(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)
	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantAccountRoleShowCaching)
	grantToUser1 := model.GrantAccountRole("to_user_1", role.ID().Name()).WithUserName(user1.ID().Name())
	grantToUser2 := model.GrantAccountRole("to_user_2", role.ID().Name()).WithUserName(user2.ID().Name())
	grantToParent := model.GrantAccountRole("to_parent", role.ID().Name()).WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantAccountRoleShowCachingProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantToUser1, grantToUser2, grantToParent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_user_1", "user_name", user1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_user_2", "user_name", user2.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_parent", "parent_role_name", parentRole.ID().Name()),
				),
			},
			// the second refresh must converge: every resource Reads the same role and must still find
			// its own grant in the (cached) SHOW GRANTS result, yielding an empty plan.
			{
				Config: config.FromModels(t, experimentProviderModel, grantToUser1, grantToUser2, grantToParent),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
