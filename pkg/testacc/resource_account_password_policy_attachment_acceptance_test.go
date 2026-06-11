//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountPasswordPolicyAttachment_BasicUseCase(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	passwordPolicy, passwordPolicyCleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	attachmentModel := model.AccountPasswordPolicyAttachment("att", passwordPolicy.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, attachmentModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_account_password_policy_attachment.att", "id"),
				),
			},
			{
				ResourceName:      "snowflake_account_password_policy_attachment.att",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"initially_suspended",
					"wait_for_provisioning",
					"query_acceleration_max_scale_factor",
					"max_concurrency_level",
					"statement_queued_timeout_in_seconds",
					"statement_timeout_in_seconds",
				},
			},
		},
	})
}
