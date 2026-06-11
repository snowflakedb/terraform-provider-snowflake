//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicyAttachmentUser(t *testing.T) {
	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	policyId := testClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModel := model.NetworkPolicy("test", policyId.Name()).WithAllowedIps("1.1.1.1", "2.2.2.2")

	attachmentModelSingleUser := model.NetworkPolicyAttachment("test", policyId.Name()).
		WithSetForAccount(false).
		WithUsers(user1.ID().Name()).
		WithDependsOn(networkPolicyModel.ResourceReference())

	attachmentModelMultipleUsers := model.NetworkPolicyAttachment("test", policyId.Name()).
		WithSetForAccount(false).
		WithUsers(user1.ID().Name(), user2.ID().Name()).
		WithDependsOn(networkPolicyModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, networkPolicyModel, attachmentModelSingleUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyId.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, networkPolicyModel, attachmentModelMultipleUsers),
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
	testClient().EnsureValidNonProdAccountIsUsed(t)

	policyNameAccount := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	networkPolicyModel := model.NetworkPolicy("test", policyNameAccount).WithAllowedIps("0.0.0.0/0")

	attachmentModel := model.NetworkPolicyAttachment("test", policyNameAccount).
		WithSetForAccount(true).
		WithDependsOn(networkPolicyModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckNetworkPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, networkPolicyModel, attachmentModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyNameAccount),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkPolicyAttachmentDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*provider.Context).Client

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
