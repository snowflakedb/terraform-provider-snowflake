//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Accounts_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	prefix := acc.TestClient().Ids.AlphaN(4)

	publicKey, _ := random.GenerateRSAPublicKey(t)
	account, accountCleanup := acc.TestClient().Account.CreateWithRequest(t, acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix), &sdk.CreateAccountOptions{
		AdminName:         acc.TestClient().Ids.Alpha(),
		AdminRSAPublicKey: &publicKey,
		AdminUserType:     sdk.Pointer(sdk.UserTypeService),
		Email:             "test@example.com",
		Edition:           sdk.EditionStandard,
	})
	t.Cleanup(accountCleanup)

	_, account2Cleanup := acc.TestClient().Account.CreateWithRequest(t, acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix), &sdk.CreateAccountOptions{
		AdminName:         acc.TestClient().Ids.Alpha(),
		AdminRSAPublicKey: &publicKey,
		AdminUserType:     sdk.Pointer(sdk.UserTypeService),
		Email:             "test@example.com",
		Edition:           sdk.EditionStandard,
	})
	t.Cleanup(account2Cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accountsConfig(prefix + "%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_accounts.test", "accounts.#", "2"),
				),
			},
			{
				Config: accountsConfig(account.ID().Name()),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_accounts.test", "accounts.#", "1")),
					resourceshowoutputassert.AccountDatasourceShowOutput(t, "snowflake_accounts.test").
						HasOrganizationName(account.OrganizationName).
						HasAccountName(account.AccountName).
						HasSnowflakeRegion(account.SnowflakeRegion).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						HasAccountUrlNotEmpty().
						HasCreatedOnNotEmpty().
						HasComment("SNOWFLAKE").
						HasAccountLocatorNotEmpty().
						HasAccountLocatorUrlNotEmpty().
						HasManagedAccounts(0).
						HasConsumptionBillingEntityNameNotEmpty().
						HasMarketplaceConsumerBillingEntityName("").
						HasMarketplaceProviderBillingEntityNameNotEmpty().
						HasOldAccountURL("").
						HasIsOrgAdmin(false).
						HasAccountOldUrlSavedOnEmpty().
						HasAccountOldUrlLastUsedEmpty().
						HasOrganizationOldUrlEmpty().
						HasOrganizationOldUrlSavedOnEmpty().
						HasOrganizationOldUrlLastUsedEmpty().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false).
						HasDroppedOnEmpty().
						HasScheduledDeletionTimeEmpty().
						HasRestoredOnEmpty().
						HasMovedToOrganizationEmpty().
						HasMovedOnEmpty().
						HasOrganizationUrlExpirationOnEmpty(),
				),
			},
		},
	})
}

func accountsConfig(pattern string) string {
	return fmt.Sprintf(`data "snowflake_accounts" "test" {
	with_history = true
	like = "%s"
}`, pattern)
}
