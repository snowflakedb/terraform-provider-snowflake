//go:build account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
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

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)

	databaseRoleId := acc.SecondaryTestClient().Ids.RandomDatabaseObjectIdentifier()
	parentRoleId := acc.SecondaryTestClient().Ids.RandomAccountObjectIdentifier()

	user, userCleanup := acc.SecondaryTestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:            config.FromModels(t, providerModel) + grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
				ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0"),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_account_role.test", "id", helpers.EncodeResourceIdentifier(parentRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleId.FullyQualifiedName()))),
				),
			},
			{
				PreConfig: func() {
					acc.SecondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_02")
					acc.SecondaryTestClient().Grant.GrantDatabaseRoleToUser(t, databaseRoleId, user.ID())
				},
				ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0"),
				Config:            config.FromModels(t, providerModel) + grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
				ExpectError:       regexp.MustCompile("Provider produced inconsistent result after apply"),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_account_role.test", "id", helpers.EncodeResourceIdentifier(parentRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleId.FullyQualifiedName()))),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providerModel) + grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_account_role.test", "id", helpers.EncodeResourceIdentifier(parentRoleId))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRoleId.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleId.FullyQualifiedName()))),
				),
			},
		},
	})
}

func grantDatabaseRoleIssue3629Config(databaseRoleId sdk.DatabaseObjectIdentifier, accountRoleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database_role" "test" {
    database = "%[1]s"
    name = "%[2]s"
}

resource "snowflake_account_role" "test" {
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  parent_role_name = snowflake_account_role.test.fully_qualified_name
}
`, databaseRoleId.DatabaseName(), databaseRoleId.Name(), accountRoleId.Name())
}
