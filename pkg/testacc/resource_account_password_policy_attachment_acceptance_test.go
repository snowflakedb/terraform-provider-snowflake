//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountPasswordPolicyAttachment_BasicUseCase(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountPasswordPolicyAttachmentConfig(id),
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

func accountPasswordPolicyAttachmentConfig(id sdk.SchemaObjectIdentifier) string {
	s := `
resource "snowflake_password_policy" "pa" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_account_password_policy_attachment" "att" {
	password_policy = snowflake_password_policy.pa.fully_qualified_name
}
`
	return fmt.Sprintf(s, id.DatabaseName(), id.SchemaName(), id.Name())
}
