//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 is fixed
func TestAcc_GrantDatabaseRole_Issue_3629(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	accountRole, accountRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	testConfig := grantDatabaseRoleIssue3629Config(databaseRole.ID(), accountRole.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Grant.GrantDatabaseRoleToUser(t, databaseRole.ID(), user.ID())
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.0.0"),
				Config:            testConfig,
				ExpectError:       regexp.MustCompile("Provider produced inconsistent result after apply"),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), accountRole.ID().FullyQualifiedName()))),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   testConfig,
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), accountRole.ID().FullyQualifiedName()))),
				),
			},
		},
	})
}

func grantDatabaseRoleIssue3629Config(databaseRoleId sdk.DatabaseObjectIdentifier, accountRoleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_database_role" "test" {
  database_role_name = %[1]s
  parent_role_name = "%[2]s"
}
`, strconv.Quote(databaseRoleId.FullyQualifiedName()), accountRoleId.Name())
}
