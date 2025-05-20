//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserAuthenticationPolicyAttachment(t *testing.T) {
	// TODO [SNOW-1423486]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-1423486", "Error: 000606 (57P03): No active warehouse selected in the current session.  Select an active warehouse with the 'use warehouse' command.")

	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	userId := user1.ID()
	newUserId := user2.ID()
	authenticationPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()
	newAuthenticationPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		CheckDestroy:             CheckUserAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userAuthenticationPolicyAttachmentConfig(userId, authenticationPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "user_name", userId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "authentication_policy_name", authenticationPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), authenticationPolicyId.FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: userAuthenticationPolicyAttachmentConfig(newUserId, newAuthenticationPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "user_name", newUserId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "authentication_policy_name", newAuthenticationPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), newAuthenticationPolicyId.FullyQualifiedName())),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_authentication_policy_attachment.ppa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userAuthenticationPolicyAttachmentConfig(userId sdk.AccountObjectIdentifier, authenticationPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_authentication_policy" "ap" {
	database   = "%[2]s"
	schema     = "%[3]s"
	name       = "%[4]s"
}

resource "snowflake_user_authentication_policy_attachment" "apa" {
	authentication_policy_name = snowflake_authentication_policy.ap.fully_qualified_name
	user_name =  "%[1]s"
}
`, userId.Name(), authenticationPolicyId.DatabaseName(), authenticationPolicyId.SchemaName(), authenticationPolicyId.Name())
}
