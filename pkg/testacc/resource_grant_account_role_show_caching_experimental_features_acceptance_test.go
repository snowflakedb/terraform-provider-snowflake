//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// countShowGrantsOfRoleQueries returns how many `SHOW GRANTS OF ROLE <id>` statements were issued
// (across all sessions of the test user) for the given, test-unique role. Because the role name is
// randomly generated per test, this count is isolated to the queries produced by the test under
// inspection. Used to prove the GRANT_ACCOUNT_ROLE_SHOW_CACHING experiment collapses N identical
// SHOW calls into one.
func countShowGrantsOfRoleQueries(t *testing.T, roleId sdk.AccountObjectIdentifier) int {
	t.Helper()
	queryHistory := testClient().InformationSchema.GetQueryHistory(t, 1000)
	needle := fmt.Sprintf("SHOW GRANTS OF ROLE %s", roleId.FullyQualifiedName())
	return len(collections.Filter(queryHistory, func(h helpers.QueryHistory) bool {
		return strings.Contains(h.QueryText, needle)
	}))
}

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

// TestAcc_GrantAccountRole_ShowCaching_SharedRoleIssuesSingleShow proves the cache effect: with the
// experiment enabled, granting the same role to several grantees results in exactly one
// `SHOW GRANTS OF ROLE` call for that role over the provider's lifetime, instead of one per instance.
// The query-history count is asserted in the second step, by which point at least one Read pass has
// run; because the provider (and its cache) is reused across steps, the role is shown at most once.
func TestAcc_GrantAccountRole_ShowCaching_SharedRoleIssuesSingleShow(t *testing.T) {
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
			// the second refresh must converge (every resource Reads the same, cached role) and the
			// role must have been shown exactly once across both steps.
			{
				Config: config.FromModels(t, experimentProviderModel, grantToUser1, grantToUser2, grantToParent),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					if got := countShowGrantsOfRoleQueries(t, role.ID()); got != 1 {
						return fmt.Errorf("expected exactly 1 `SHOW GRANTS OF ROLE %s` with caching enabled, got %d", role.ID().FullyQualifiedName(), got)
					}
					return nil
				},
			},
		},
	})
}

// TestAcc_GrantAccountRole_SharedRoleIssuesShowPerInstanceWithoutExperiment is the contrast to
// TestAcc_GrantAccountRole_ShowCaching_SharedRoleIssuesSingleShow: without the experiment, the same
// role shared by several grantees is shown multiple times (once per Read of each instance), which is
// exactly the redundancy the experiment removes.
func TestAcc_GrantAccountRole_SharedRoleIssuesShowPerInstanceWithoutExperiment(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)
	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantToUser1 := model.GrantAccountRole("to_user_1", role.ID().Name()).WithUserName(user1.ID().Name())
	grantToUser2 := model.GrantAccountRole("to_user_2", role.ID().Name()).WithUserName(user2.ID().Name())
	grantToParent := model.GrantAccountRole("to_parent", role.ID().Name()).WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, grantToUser1, grantToUser2, grantToParent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_user_1", "user_name", user1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_user_2", "user_name", user2.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.to_parent", "parent_role_name", parentRole.ID().Name()),
				),
			},
			{
				Config: config.FromModels(t, grantToUser1, grantToUser2, grantToParent),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: func(_ *terraform.State) error {
					// Without the experiment, each instance issues its own SHOW on every Read pass, so the
					// shared role is shown well more than once. We assert "more than one" rather than an
					// exact number to stay robust against Terraform's refresh cadence.
					if got := countShowGrantsOfRoleQueries(t, role.ID()); got <= 1 {
						return fmt.Errorf("expected more than 1 `SHOW GRANTS OF ROLE %s` without the experiment, got %d", role.ID().FullyQualifiedName(), got)
					}
					return nil
				},
			},
		},
	})
}
