//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_UserPublicKeys(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
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
