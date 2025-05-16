//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1284394]: Unskip the test
func TestAcc_Share(t *testing.T) {
	t.Skip("second and third account must be set for Share acceptance tests")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	var account2 string
	var account3 string

	shareComment := "Created by a Terraform acceptance test"
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Share),
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

func TestAcc_Share_validateAccounts(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Share),
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

func shareConfig(shareId sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%s"
	comment        = "%s"
}
`, shareId.Name(), comment)
}

func shareConfigOneAccount(shareId sdk.AccountObjectIdentifier, comment string, account string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%s"
	comment        = "%s"
	accounts       = ["%s"]
}
`, shareId.Name(), comment, account)
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
