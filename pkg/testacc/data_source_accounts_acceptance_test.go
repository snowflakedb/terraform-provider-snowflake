//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Accounts_Complete(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	prefix := testClient().Ids.AlphaN(4)

	publicKey, _ := random.GenerateRSAPublicKey(t)
	id1 := sdk.NewAccountObjectIdentifier(fmt.Sprintf("%s_%s", prefix, random.AccountName()))
	account, accountCleanup := testClient().Account.CreateWithRequest(t, id1, &sdk.CreateAccountOptions{
		AdminName:         testClient().Ids.Alpha(),
		AdminRSAPublicKey: &publicKey,
		AdminUserType:     sdk.Pointer(sdk.UserTypeService),
		Email:             "test@example.com",
		Edition:           sdk.EditionStandard,
	})
	t.Cleanup(accountCleanup)

	id2 := sdk.NewAccountObjectIdentifier(fmt.Sprintf("%s_%s", prefix, random.AccountName()))
	_, account2Cleanup := testClient().Account.CreateWithRequest(t, id2, &sdk.CreateAccountOptions{
		AdminName:         testClient().Ids.Alpha(),
		AdminRSAPublicKey: &publicKey,
		AdminUserType:     sdk.Pointer(sdk.UserTypeService),
		Email:             "test@example.com",
		Edition:           sdk.EditionStandard,
	})
	t.Cleanup(account2Cleanup)

	provider := providermodel.SnowflakeProvider().WithRole(snowflakeroles.Orgadmin.Name())
	accountsWithPattern := datasourcemodel.Accounts("test").WithWithHistory(true).WithLike(prefix + "%")
	accountsWithAccountName := datasourcemodel.Accounts("test").WithWithHistory(true).WithLike(account.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, provider, accountsWithPattern),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_accounts.test", "accounts.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, provider, accountsWithAccountName),
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
