//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_GrantAccountRole_FailsOnPublicRole verifies that granting
// the PUBLIC role to a parent role causes an inconsistent-result error without the experiment.
// PUBLIC is always implicitly granted so SHOW GRANTS doesn't list it, and Read clears state.
func TestAcc_GrantAccountRole_FailsOnPublicRole(t *testing.T) {
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithParentRoleName(parentRole.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, grantModel),
				ExpectError: regexp.MustCompile("Root object was present, but now absent"),
			},
		},
	})
}

// TestAcc_GrantAccountRole_SucceedsOnPublicRoleWithExperiment verifies that granting
// the PUBLIC role to a parent role succeeds when the experiment is enabled. Both create and
// refresh (plan) should succeed because Read skips the SHOW GRANTS check for PUBLIC.
func TestAcc_GrantAccountRole_SucceedsOnPublicRoleWithExperiment(t *testing.T) {
	parentRole, parentRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithParentRoleName(parentRole.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantAccountRoleSafePublicRole)

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantAccountRoleSafePublicRoleProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRole.ID().Name()),
				),
			},
			// import
			{
				Config:            config.FromModels(t, experimentProviderModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"PUBLIC"|%s|"%s"`, sdk.ObjectTypeRole, parentRole.ID().Name()),
			},
		},
	})
}

// TestAcc_GrantAccountRole_SucceedsOnPublicRoleToUserWithExperiment verifies that granting
// the PUBLIC role to a user succeeds when the experiment is enabled.
func TestAcc_GrantAccountRole_SucceedsOnPublicRoleToUserWithExperiment(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	grantModel := model.GrantAccountRole("test", "PUBLIC").
		WithUserName(user.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantAccountRoleSafePublicRole)

	resourceName := "snowflake_grant_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: grantAccountRoleSafePublicRoleProviderFactory,
		CheckDestroy:             CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, experimentProviderModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "user_name", user.ID().Name()),
				),
			},
			// import
			{
				Config:            config.FromModels(t, experimentProviderModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf(`"PUBLIC"|%s|"%s"`, sdk.ObjectTypeUser, user.ID().Name()),
			},
		},
	})
}
