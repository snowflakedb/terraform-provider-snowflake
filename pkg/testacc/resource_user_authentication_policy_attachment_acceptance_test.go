//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserAuthenticationPolicyAttachment(t *testing.T) {
	// TODO [SNOW-2661409]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-2661409", "Error: 000606 (57P03): No active warehouse selected in the current session.  Select an active warehouse with the 'use warehouse' command.")

	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	userId := user1.ID()
	newUserId := user2.ID()
	authenticationPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()
	newAuthenticationPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: userAuthenticationPoliciesProviderFactory,
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

// prove that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3672 is fixed
func TestAcc_UserAuthenticationPolicyAttachment_MissingUser(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	authPolicy, authPolicyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(authPolicyCleanup)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				Config:            userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "user_name", user.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "authentication_policy_name", authPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "id", fmt.Sprintf("%s|%s", user.ID().FullyQualifiedName(), authPolicy.ID().FullyQualifiedName())),
				),
			},
			{
				PreConfig: func() {
					userCleanup()
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				Config:            userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError:       regexp.MustCompile("object does not exist or not authorized"),
				Check:             resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "id", fmt.Sprintf("%s|%s", user.ID().FullyQualifiedName(), authPolicy.ID().FullyQualifiedName())),
			},
			{
				ProtoV6ProviderFactories: userAuthenticationPoliciesProviderFactory,
				Config:                   userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError:              regexp.MustCompile("object does not exist or not authorized"),
				Check:                    resource.TestCheckNoResourceAttr("snowflake_user_authentication_policy_attachment.test", "id"),
			},
			{
				ProtoV6ProviderFactories: userAuthenticationPoliciesProviderFactory,
				Config:                   userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError:              regexp.MustCompile("error while creating authentication policy attachment"),
				Check:                    resource.TestCheckNoResourceAttr("snowflake_user_authentication_policy_attachment.test", "id"),
			},
		},
	})
}

// prove that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3672 is fixed
func TestAcc_UserAuthenticationPolicyAttachment_MissingAuthPolicy(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	authPolicy, authPolicyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(authPolicyCleanup)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				Config:            userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "user_name", user.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "authentication_policy_name", authPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "id", fmt.Sprintf("%s|%s", user.ID().FullyQualifiedName(), authPolicy.ID().FullyQualifiedName())),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				PreConfig: func() {
					testClient().User.Alter(t, user.ID(), &sdk.AlterUserOptions{Unset: &sdk.UserUnset{AuthenticationPolicy: sdk.Bool(true)}})
					authPolicyCleanup()
				},
				Config:      userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
				Check:       resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.test", "id", fmt.Sprintf("%s|%s", user.ID().FullyQualifiedName(), authPolicy.ID().FullyQualifiedName())),
			},
			{
				ProtoV6ProviderFactories: userAuthenticationPoliciesProviderFactory,
				Config:                   userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError:              regexp.MustCompile("object does not exist or not authorized"),
				Check:                    resource.TestCheckNoResourceAttr("snowflake_user_authentication_policy_attachment.test", "id"),
			},
			{
				ProtoV6ProviderFactories: userAuthenticationPoliciesProviderFactory,
				Config:                   userauthenticationpolicyattachmentMissingDependentObjectsConfig(user.ID(), authPolicy.ID()),
				ExpectError:              regexp.MustCompile("error while creating authentication policy attachment"),
				Check:                    resource.TestCheckNoResourceAttr("snowflake_user_authentication_policy_attachment.test", "id"),
			},
		},
	})
}

func userauthenticationpolicyattachmentMissingDependentObjectsConfig(userId sdk.AccountObjectIdentifier, authenticationPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
provider "snowflake" {
	preview_features_enabled = [ "snowflake_user_authentication_policy_attachment_resource" ]
}

resource "snowflake_user_authentication_policy_attachment" "test" {
	user_name =  "%[1]s"
	authentication_policy_name = %[2]s
}
`, userId.Name(), strconv.Quote(authenticationPolicyId.FullyQualifiedName()))
}
