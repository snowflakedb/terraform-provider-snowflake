//go:build account_level_tests

package testacc

import (
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForPartnerApplications_WithPrivilegedRolesBlockedList(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	// Use an identifier with this prefix to have this role in the end.
	roleId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("Z")
	role, roleCleanup := testClient().Role.CreateRoleWithIdentifier(t, roleId)
	t.Cleanup(roleCleanup)

	allRoles := []string{snowflakeroles.Accountadmin.Name(), snowflakeroles.SecurityAdmin.Name(), role.ID().Name()}
	onlyPrivilegedRoles := []string{snowflakeroles.Accountadmin.Name(), snowflakeroles.SecurityAdmin.Name()}
	customRoles := []string{role.ID().Name()}

	paramCleanup := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList, "true")
	t.Cleanup(paramCleanup)

	modelWithoutBlockedRole := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop))
	modelWithBlockedRole := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList(role.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.*", role.ID().Name()),
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(allRoles, ",")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "blocked_roles_list.#", "0"),
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(onlyPrivilegedRoles, ",")),
				),
			},
			{
				PreConfig: func() {
					// Do not revert, because the revert is setup above.
					testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList, "false")
				},
				Config: accconfig.FromModels(t, modelWithBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.*", role.ID().Name()),
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(customRoles, ",")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "blocked_roles_list.#", "0"),
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", ""),
				),
			},
		},
	})
}
