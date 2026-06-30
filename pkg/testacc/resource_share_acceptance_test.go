//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1284394]: Unskip the test
// TODO [SNOW-1348347]: Add more tests or merge with other tests.
func TestAcc_Share(t *testing.T) {
	t.Skip("second and third account must be set for Share acceptance tests")

	var account2 string
	var account3 string

	shareComment := "Created by a Terraform acceptance test"
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config: shareConfig(id, shareComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_share.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_share.test", "comment", shareComment),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			{
				Config: shareConfigTwoAccounts(id, shareComment, account2, account3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.1", account3),
				),
			},
			{
				Config: shareConfigOneAccount(id, shareComment, account2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
				),
			},
			{
				Config: shareConfig(id, shareComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_share.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO [SNOW-1348347]: Add more tests or merge with other tests.
func TestAcc_Share_basic(t *testing.T) {
	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config: shareConfigOneAccount(id, "", account2.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_share.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2.Name()),
				),
			},
		},
	})
}

func TestAcc_Share_validateAccounts(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config:      shareConfigOneAccount(id, "any comment", "incorrect"),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
			{
				Config:      shareConfigTwoAccounts(id, "any comment", "correct.one", "incorrect"),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/4398 is fixed
func TestAcc_Share_issue4398(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.ShareResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            shareConfig(id, "") + accconfig.FromModels(t, providerModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			// Step 2: Grant USAGE on database to the share externally, then add accounts.
			// v2.13.0 fails because setShareAccounts creates a temp database and tries to
			// grant USAGE on it, conflicting with the already-granted USAGE on the real database.
			{
				PreConfig: func() {
					testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), id, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            shareConfigOneAccount(id, "", account2.Name()) + accconfig.FromModels(t, providerModel),
				ExpectError:       regexp.MustCompile("does not belong to the database that is being shared"),
			},
			// Step 3: Succeeds after fix (skip temp database when share already has a database granted).
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   shareConfigOneAccount(id, "", account2.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2.Name()),
				),
			},
		},
	})
}

func TestAcc_Share_issue4398_updateAccountsWithoutDatabaseGranted(t *testing.T) {
	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config: shareConfig(id, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			{
				Config: shareConfigOneAccount(id, "", account2.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2.Name()),
				),
			},
		},
	})
}

// TestAcc_Share_withGrantedDatabase proves that adding and updating accounts on a share
// that already has a database granted (via snowflake_grant_privileges_to_share) works correctly.
// This is the regression scenario reported after PR #4174 — previously, the provider would
// erroneously try to create a temp database even when a real DB was already in the share,
// causing a conflict on the USAGE grant.
// See: https://github.com/snowflakedb/terraform-provider-snowflake/pull/4174
func TestAcc_Share_withGrantedDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	shareModel := model.ShareWithDefaultMeta(id.Name())
	shareModelWithAccount := model.ShareWithDefaultMeta(id.Name()).WithAccounts(account2.Name())
	grantModel := model.GrantPrivilegesToShareWithDefaultMeta([]string{"USAGE"}, id.Name()).
		WithOnDatabase(database.ID().Name()).
		WithDependsOn(shareModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts + grant DB to share (no accounts yet).
			{
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(
					t,
					resourceassert.ShareResource(t, shareModel.ResourceReference()).
						HasName(id.Name()).
						HasAccountsEmpty(),
				),
			},
			// Step 2: Add an account while DB is already on the share.
			// This exercises the UpdateShare path that directly adds accounts
			// (skipping the temp DB workaround) when share.DatabaseName is set.
			{
				Config: accconfig.FromModels(t, shareModelWithAccount, grantModel),
				Check: assertThat(
					t,
					resourceassert.ShareResource(t, shareModelWithAccount.ResourceReference()).
						HasName(id.Name()).
						HasAccounts(account2.Name()),
				),
			},
			// Step 3: Remove account (DB still granted).
			{
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(
					t,
					resourceassert.ShareResource(t, shareModel.ResourceReference()).
						HasName(id.Name()).
						HasAccountsEmpty(),
				),
			},
		},
	})
}

func shareConfig(shareId sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%s"
	comment        = "%s"
}
`, shareId.Name(), comment)
}

func shareConfigTwoAccounts(shareId sdk.AccountObjectIdentifier, comment string, account string, account2 string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%s"
	comment        = "%s"
	accounts       = ["%s", "%s"]
}
`, shareId.Name(), comment, account, account2)
}
