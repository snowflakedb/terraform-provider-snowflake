//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	passwordPolicy1, passwordPolicy1Cleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicy1Cleanup)

	passwordPolicy2, passwordPolicy2Cleanup := testClient().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicy2Cleanup)

	attachmentModel := model.UserPasswordPolicyAttachment("ppa", passwordPolicy1.ID().FullyQualifiedName(), user1.ID().Name())

	updatedAttachmentModel := model.UserPasswordPolicyAttachment("ppa", passwordPolicy2.ID().FullyQualifiedName(), user2.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: userPasswordPoliciesProviderFactory,
		CheckDestroy:             CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: accconfig.FromModels(t, attachmentModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", user1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicy1.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", user1.ID().FullyQualifiedName(), passwordPolicy1.ID().FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: accconfig.FromModels(t, updatedAttachmentModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", user2.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicy2.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", user2.ID().FullyQualifiedName(), passwordPolicy2.ID().FullyQualifiedName())),
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
