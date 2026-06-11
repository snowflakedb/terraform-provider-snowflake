//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1007539]: use email of our service user (verified email address is required)
func TestAcc_EmailNotificationIntegration(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	verifiedEmail := "artur.sawicki@snowflake.com"

	integrationModel := model.EmailNotificationIntegration("test", id.Name(), true).
		WithAllowedRecipients(verifiedEmail).
		WithComment("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.EmailNotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, integrationModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.0", verifiedEmail),
				),
			},
			{
				ResourceName:      "snowflake_email_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAcc_EmailNotificationIntegration_issue2223 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2223 issue.
// Snowflake allowed empty allowed recipients in https://docs.snowflake.com/en/release-notes/2023/7_40#email-notification-integrations-allowed-recipients-no-longer-required.
func TestAcc_EmailNotificationIntegration_issue2223(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	verifiedEmail := "artur.sawicki@snowflake.com"

	integrationModelWithoutRecipients := model.EmailNotificationIntegration("test", id.Name(), true)

	integrationModelWithRecipients := model.EmailNotificationIntegration("test", id.Name(), true).
		WithAllowedRecipients(verifiedEmail)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.EmailNotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, integrationModelWithoutRecipients),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, integrationModelWithRecipients),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.0", verifiedEmail),
				),
			},
			{
				Config: accconfig.FromModels(t, integrationModelWithoutRecipients),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.#", "0"),
				),
			},
		},
	})
}
