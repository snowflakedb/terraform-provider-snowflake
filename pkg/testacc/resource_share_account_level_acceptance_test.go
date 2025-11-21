//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Share_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	// NOTE: In this case, we swap the test client to the secondary test client from the TestAcc_Share_basic test.
	// We enable the BCR bundle on the secondary test client, and make a share to the primary test client.
	account2 := testClient().Account.GetAccountIdentifier(t)

	id := secondaryTestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: secondaryAccountProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config: shareConfigOneAccount(id, "", account2.Name()) + config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)),
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
