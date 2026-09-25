//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_GrantAccountRole_FailsOnPublicRole_Disabled verifies that granting
// the PUBLIC role to a parent role causes an inconsistent-result error when the experiment is disabled.
// PUBLIC is always implicitly granted so SHOW GRANTS doesn't list it, and Read clears state.
func TestAcc_GrantAccountRole_FailsOnPublicRole_Disabled(t *testing.T) {
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	providerModel := providermodel.SnowflakeProvider().
		WithAllEnabledByDefaultExperimentsDisabled()
	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: enabledByDefaultExperimentsDisabledProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providerModel, grantModel),
				ExpectError: regexp.MustCompile("Root object was present, but now absent"),
			},
		},
	})
}

// TestAcc_GrantAccountRole_SucceedsOnPublicRole verifies that granting
// the PUBLIC role to a parent role succeeds. Both create and refresh (plan) should succeed
// because Read skips the SHOW GRANTS check for PUBLIC.
func TestAcc_GrantAccountRole_SucceedsOnPublicRole(t *testing.T) {
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithParentRoleName(parentRole.ID().Name())

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRole.ID().Name()),
				),
			},
			{
				Config:            config.FromModels(t, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"PUBLIC"|%s|"%s"`, sdk.ObjectTypeRole, parentRole.ID().Name()),
			},
		},
	})
}

// TestAcc_GrantAccountRole_SucceedsOnPublicRoleToUser verifies that granting
// the PUBLIC role to a user succeeds.
func TestAcc_GrantAccountRole_SucceedsOnPublicRoleToUser(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithUserName(user.ID().Name())

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "user_name", user.ID().Name()),
				),
			},
			{
				Config:            config.FromModels(t, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"PUBLIC"|%s|"%s"`, sdk.ObjectTypeUser, user.ID().Name()),
			},
		},
	})
}
