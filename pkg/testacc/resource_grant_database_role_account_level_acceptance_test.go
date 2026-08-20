//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
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
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), accountRole.ID().FullyQualifiedName()))),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   testConfig,
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", helpers.EncodeResourceIdentifier(databaseRole.ID().FullyQualifiedName(), sdk.ObjectTypeRole.String(), accountRole.ID().FullyQualifiedName()))),
				),
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_bcr2026_06_databaseRoleGrantee(t *testing.T) {
	// Regression test for SNOW-3953574 (PR https://github.com/snowflakedb/terraform-provider-snowflake/pull/5099):
	// the 2026_06 bundle (BCR-2371) makes SHOW GRANTS OF DATABASE ROLE return the grantee database role
	// without the database prefix. Before the fix, ReadGrantDatabaseRole could not parse the unprefixed
	// grantee, cleared the id, and the resource was perpetually recreated (apply failed with an
	// "inconsistent result after apply" error). The fixed provider reconstructs the fully qualified
	// grantee and matches the grant again.
	// TODO(SNOW-3953574): skip or remove this test once the 2026_06 bundle is enforced and can no longer be disabled.

	// A database role can only be granted to another database role in the same database, and
	// CreateDatabaseRole creates both in the same database. Run on the secondary account because toggling
	// a behavior change bundle is account-wide.
	childRole, childRoleCleanup := secondaryTestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(childRoleCleanup)

	parentRole, parentRoleCleanup := secondaryTestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(parentRoleCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	grantModel := model.GrantDatabaseRole("test", childRole.ID().FullyQualifiedName()).
		WithParentDatabaseRoleName(parentRole.ID().FullyQualifiedName())
	expectedId := helpers.EncodeResourceIdentifier(childRole.ID().FullyQualifiedName(), sdk.ObjectTypeDatabaseRole.String(), parentRole.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Latest released provider, bundle disabled: the grantee is returned fully qualified and parses fine.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.19.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2026_06")
				},
				Config: config.FromModels(t, providerModel, grantModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(grantModel.ResourceReference(), "id", expectedId)),
				),
			},
			// Same released provider, bundle enabled: the unprefixed grantee can no longer be parsed, so Read
			// clears the id and apply fails with an inconsistent-result error.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.19.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_06")
				},
				Config:      config.FromModels(t, providerModel, grantModel),
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
			// Fixed provider, bundle still enabled: the grantee is normalized back to a fully qualified
			// identifier, the grant matches on Read, and the resource is stable.
			// secondaryAccountProviderFactory is required because TestAccProtoV6ProviderFactories is
			// cached against the primary account.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   config.FromModels(t, providerModel, grantModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(grantModel.ResourceReference(), "id", expectedId)),
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
