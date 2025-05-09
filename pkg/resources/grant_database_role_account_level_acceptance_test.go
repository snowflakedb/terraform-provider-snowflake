//go:build account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629 is fixed
func TestAcc_GrantDatabaseRole_Issue_3629(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRole, databaseRoleCleanup := acc.SecondaryTestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	accountRole, accountRoleCleanup := acc.SecondaryTestClient().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	user, userCleanup := acc.SecondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	testConfig := accconfig.FromModels(t, providerModel) + grantDatabaseRoleIssue3629Config(databaseRole.ID(), accountRole.ID())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.SecondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_02")
					acc.SecondaryTestClient().Grant.GrantDatabaseRoleToUser(t, databaseRole.ID(), user.ID())
				},
				ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0"),
				Config:            testConfig,
				ExpectError:       regexp.MustCompile("Provider produced inconsistent result after apply"),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), accountRole.ID().FullyQualifiedName()))),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
