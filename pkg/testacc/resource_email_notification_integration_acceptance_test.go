//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_EmailNotificationIntegration(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.EmailNotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: emailNotificationIntegrationConfig(id.Name(), helpers.VerifiedEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.0", helpers.VerifiedEmail),
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

func emailNotificationIntegrationConfig(name string, email string) string {
	s := `
resource "snowflake_email_notification_integration" "test" {
  name               = "%s"
  enabled            = true
  allowed_recipients = ["%s"]
  comment            = "test"
}
`
	return fmt.Sprintf(s, name, email)
}

// TestAcc_EmailNotificationIntegration_issue2223 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2223 issue.
// Snowflake allowed empty allowed recipients in https://docs.snowflake.com/en/release-notes/2023/7_40#email-notification-integrations-allowed-recipients-no-longer-required.
func TestAcc_EmailNotificationIntegration_issue2223(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	emailIntegrationName := id.Name()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.EmailNotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: emailNotificationIntegrationWithoutRecipientsConfig(emailIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.#", "0"),
				),
			},
			{
				Config: emailNotificationIntegrationConfig(emailIntegrationName, helpers.VerifiedEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.0", helpers.VerifiedEmail),
				),
			},
			{
				Config: emailNotificationIntegrationWithoutRecipientsConfig(emailIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.#", "0"),
				),
			},
		},
	})
}

func emailNotificationIntegrationWithoutRecipientsConfig(name string) string {
	s := `
resource "snowflake_email_notification_integration" "test" {
  name               = "%s"
  enabled            = true
}`
	return fmt.Sprintf(s, name)
}
