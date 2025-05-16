//go:build !account_level_tests

package resources_test

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

func TestAcc_UserPublicKeys(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: uPublicKeysConfig(user.ID(), key1, key2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key", key1),
					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key_2", key2),
					resource.TestCheckNoResourceAttr("snowflake_user_public_keys.foobar", "has_rsa_public_key"),
				),
			},
		},
	})
}

func uPublicKeysConfig(userId sdk.AccountObjectIdentifier, key1 string, key2 string) string {
	return fmt.Sprintf(`
resource "snowflake_user_public_keys" "foobar" {
	name = "%[1]s"
	rsa_public_key = <<KEY
%[2]s
	KEY

	rsa_public_key_2 = <<KEY
%[3]s
	KEY
}
`, userId.Name(), key1, key2)
}
