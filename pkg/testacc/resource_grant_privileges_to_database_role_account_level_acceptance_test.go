//go:build account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 (UBAC) doesn't affect the grant privileges to database role resource
func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase_WithPrivilegesGrantedOnDatabaseToUser(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	databaseId := testClient().Ids.DatabaseId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Grant.GrantPrivilegesOnDatabaseToUser(t, databaseId, user.ID(), sdk.AccountObjectPrivilegeUsage, sdk.AccountObjectPrivilegeMonitor)
				},
				Config: grantPrivilegesToDatabaseRoleOnDatabaseConfig(databaseRole.ID(), databaseId, sdk.AccountObjectPrivilegeCreateDatabaseRole, sdk.AccountObjectPrivilegeCreateSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "privileges.1", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "on_database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "id", fmt.Sprintf("%s|false|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), databaseId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToDatabaseRoleOnDatabaseConfig(databaseRoleId sdk.DatabaseObjectIdentifier, databaseId sdk.AccountObjectIdentifier, privileges ...sdk.AccountObjectPrivilege) string {
	quotedPrivileges := collections.Map(privileges, func(privilege sdk.AccountObjectPrivilege) string { return fmt.Sprintf("%q", privilege) })
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_database_role" "test" {
  	database_role_name 	= %[1]s
  	privileges        	= [%[2]s]
	on_database 		= "%[3]s"
}
`, strconv.Quote(databaseRoleId.FullyQualifiedName()), strings.Join(quotedPrivileges, ","), databaseId.Name())
}
