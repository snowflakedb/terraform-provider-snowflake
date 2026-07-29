//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountAuthenticationPolicyAttachment(t *testing.T) {
	// TODO [SNOW-1763613]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-1763613", "Error: 003549 (23505): Object <account_name> already has a AUTHENTICATION_POLICY. Only one AUTHENTICATION_POLICY is allowed at a time")
	policyName := testClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountAuthenticationPolicyAttachmentConfig(TestDatabaseName, TestSchemaName, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_account_authentication_policy_attachment.att", "id"),
				),
			},
			{
				ResourceName:      "snowflake_account_authentication_policy_attachment.att",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func accountAuthenticationPolicyAttachmentConfig(databaseName, schemaName, policyName string) string {
	s := `
resource "snowflake_authentication_policy" "pa" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_account_authentication_policy_attachment" "att" {
	authentication_policy = snowflake_authentication_policy.pa.fully_qualified_name
}
`
	return fmt.Sprintf(s, databaseName, schemaName, policyName)
}
