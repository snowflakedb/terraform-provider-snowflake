//go:build account_level_tests

package resources_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 (UBAC) doesn't affect the grant privileges to database role resource
func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase_WithPrivilegesGrantedOnDatabaseToUser(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	databaseRole, databaseRoleCleanup := acc.SecondaryTestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	user, userCleanup := acc.SecondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	databaseId := acc.SecondaryTestClient().Ids.DatabaseId()

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.SecondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_02")
					acc.SecondaryTestClient().Grant.GrantPrivilegesOnDatabaseToUser(t, databaseId, user.ID(), sdk.AccountObjectPrivilegeUsage, sdk.AccountObjectPrivilegeMonitor)
				},
				Config: accconfig.FromModels(t, providerModel) + grantPrivilegesToDatabaseRoleOnDatabaseConfig(databaseRole.ID(), databaseId, sdk.AccountObjectPrivilegeCreateDatabaseRole, sdk.AccountObjectPrivilegeCreateSchema),
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
