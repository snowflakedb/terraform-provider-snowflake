//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TODO [SNOW-1423486]: Fix using warehouse; remove unsetting testenvs.ConfigureClientOnce
func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user1, user1Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	userId := user1.ID()
	newUserId := user2.ID()
	passwordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newPasswordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userPasswordPolicyAttachmentConfig(userId, passwordPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", userId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig(newUserId, newPasswordPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", newUserId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", newPasswordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", newUserId.FullyQualifiedName(), newPasswordPolicyId.FullyQualifiedName())),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_password_policy_attachment.ppa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userPasswordPolicyAttachmentConfig(userId sdk.AccountObjectIdentifier, passwordPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_password_policy" "pp" {
	database   = "%[2]s"
	schema     = "%[3]s"
	name       = "%[4]s"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	password_policy_name = snowflake_password_policy.pp.fully_qualified_name
	user_name = "%[1]s"
}
`, userId.Name(), passwordPolicyId.DatabaseName(), passwordPolicyId.SchemaName(), passwordPolicyId.Name())
}
