//go:build !account_level_tests

package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicyAttachmentUser(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	user1, user1Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	policyId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyAttachmentConfigSingle(user1.ID(), policyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyId.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "1"),
				),
			},
			{
				Config: networkPolicyAttachmentConfig(user1.ID(), user2.ID(), policyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyId.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "2"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_network_policy_attachment.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAcc_NetworkPolicyAttachmentAccount(t *testing.T) {
	// TODO [SNOW-2010844]: unskip
	t.Skip("Skipping as it messes with the account level setting. Should be moved to manual tests and later invoked on a brand new account.")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	policyNameAccount := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckNetworkPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyAttachmentConfigAccount(policyNameAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyNameAccount),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkPolicyAttachmentDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_network_policy_attachment" {
			continue
		}
		ctx := context.Background()
		parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterNetworkPolicy)
		if err != nil {
			fmt.Printf("[WARN] network policy (%s) not found on account", rs.Primary.Attributes["Id"])
			return nil
		}
		if parameter.Level == "ACCOUNT" && parameter.Key == "NETWORK_POLICY" && parameter.Value == rs.Primary.Attributes["network_policy_name"] {
			return fmt.Errorf("network policy attachment %v still exists", rs.Primary.Attributes["Id"])
		}
	}
	return nil
}

func networkPolicyAttachmentConfigSingle(user1Id sdk.AccountObjectIdentifier, policyId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = %[2]s
	allowed_ip_list = ["1.1.1.1", "2.2.2.2"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = false
	users               = [%[1]s]
}
`, user1Id.FullyQualifiedName(), policyId.FullyQualifiedName())
}

func networkPolicyAttachmentConfig(user1Id sdk.AccountObjectIdentifier, user2Id sdk.AccountObjectIdentifier, policyId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = %[3]s
	allowed_ip_list = ["1.1.1.1", "2.2.2.2"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = false
	users               = [%[1]s, %[2]s]
}
`, user1Id.FullyQualifiedName(), user2Id.FullyQualifiedName(), policyId.FullyQualifiedName())
}

func networkPolicyAttachmentConfigAccount(policyName string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	allowed_ip_list = ["0.0.0.0/0"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = true
}
`, policyName)
}
