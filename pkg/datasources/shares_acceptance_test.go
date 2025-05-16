//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Shares(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	shareId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	shareId2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: shares(shareId, shareId2, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.0.name"),
				),
			},
			{
				Config: sharesPattern(shareId, shareId.Name(), comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.#"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.0.kind", "OUTBOUND"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.0.comment", comment),
				),
			},
		},
	})
}

func shares(shareId sdk.AccountObjectIdentifier, shareId2 sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_share "test_share" {
			name = "%[1]s"
			comment = "%[3]s"
		}
		resource snowflake_share "test_share_2" {
			name = "%[2]s"
			comment = "%[3]s"
		}
		data snowflake_shares "r" {
			depends_on = [
				snowflake_share.test_share,
				snowflake_share.test_share_2,
			]
		}
	`, shareId.Name(), shareId2.Name(), comment)
}

func sharesPattern(shareId sdk.AccountObjectIdentifier, pattern string, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_share "test_share" {
			name = "%[1]s"
			comment = "%[3]s"
		}

		data snowflake_shares "r" {
			pattern = "%[2]s"
			depends_on = [
				snowflake_share.test_share,
			]
		}
	`, shareId.Name(), pattern, comment)
}
