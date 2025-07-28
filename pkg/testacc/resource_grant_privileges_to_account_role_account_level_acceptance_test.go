//go:build account_level_tests

package testacc

import (
	"fmt"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 (UBAC) doesn't affect the grant privileges to account role resource
func TestAcc_GrantPrivilegesToAccountRole_OnDatabase_WithPrivilegesGrantedOnDatabaseToUser(t *testing.T) {
	t.Skip("TODO(SNOW-2081651): re-enable this if the test is still relevant without the BCR bundle update as now it's enabled by default in Snowflake")
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	role, roleCleanup := secondaryTestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	user, userCleanup := secondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	databaseId := secondaryTestClient().Ids.DatabaseId()

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_02")
					secondaryTestClient().Grant.GrantPrivilegesOnDatabaseToUser(t, databaseId, user.ID(), sdk.AccountObjectPrivilegeUsage, sdk.AccountObjectPrivilegeMonitor)
				},
				Config: accconfig.FromModels(t, providerModel) + grantPrivilegesToAccountRoleOnDatabaseConfig(role.ID(), databaseId, sdk.AccountObjectPrivilegeCreateDatabaseRole, sdk.AccountObjectPrivilegeCreateSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "privileges.1", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_account_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_account_object.0.object_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnAccountObject|DATABASE|%s", role.ID().FullyQualifiedName(), databaseId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToAccountRoleOnDatabaseConfig(accountRoleId sdk.AccountObjectIdentifier, databaseId sdk.AccountObjectIdentifier, privileges ...sdk.AccountObjectPrivilege) string {
	quotedPrivileges := collections.Map(privileges, func(privilege sdk.AccountObjectPrivilege) string { return fmt.Sprintf("%q", privilege) })
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "%[1]s"
  privileges        = [%[2]s]
  on_account_object {
    object_type = "DATABASE"
    object_name = "%[3]s"
  }
}
`, accountRoleId.Name(), strings.Join(quotedPrivileges, ","), databaseId.Name())
}
