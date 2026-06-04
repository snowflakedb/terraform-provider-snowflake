//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// webhookTestUrl is a valid Slack-format webhook URL accepted by Snowflake's URL format validation.
// The path segments (T.../B.../24-char-token) are not real; Snowflake does not call the URL during CREATE.
const (
	webhookTestUrl      = "https://hooks.slack.com/services/T00000000/B00000000/AAAAAAAAAAAAAAAAAAAAAAAA"
	webhookOtherTestUrl = "https://hooks.slack.com/services/T11111111/B11111111/BBBBBBBBBBBBBBBBBBBBBBBB"
)

func TestAcc_WebhookNotificationIntegration_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WebhookNotificationIntegration),
		Steps: []resource.TestStep{
			// create (minimal)
			{
				Config: webhookNotificationIntegrationConfig(id.Name(), webhookTestUrl, true, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_url", webhookTestUrl),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "comment", ""),
				),
			},
			// update enabled and comment
			{
				Config: webhookNotificationIntegrationConfig(id.Name(), webhookTestUrl, false, "test comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "comment", "test comment"),
				),
			},
			// update webhook_url
			{
				Config: webhookNotificationIntegrationConfig(id.Name(), webhookOtherTestUrl, false, "test comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_url", webhookOtherTestUrl),
				),
			},
			// import
			{
				ResourceName:      "snowflake_webhook_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_WebhookNotificationIntegration_withBodyTemplateAndHeaders(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	bodyTemplate := `{"text": "SNOWFLAKE_WEBHOOK_MESSAGE"}`
	otherBodyTemplate := `{"message": "SNOWFLAKE_WEBHOOK_MESSAGE", "channel": "#other"}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WebhookNotificationIntegration),
		Steps: []resource.TestStep{
			// create with body template and headers
			{
				Config: webhookNotificationIntegrationWithOptionsConfig(id.Name(), webhookTestUrl, bodyTemplate, map[string]string{"Content-Type": "application/json"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_url", webhookTestUrl),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_body_template", bodyTemplate),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_headers.Content-Type", "application/json"),
				),
			},
			// update body template and headers
			{
				Config: webhookNotificationIntegrationWithOptionsConfig(id.Name(), webhookTestUrl, otherBodyTemplate, map[string]string{"X-Custom": "value"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_body_template", otherBodyTemplate),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_headers.X-Custom", "value"),
				),
			},
			// unset body template and headers by going back to minimal config
			{
				Config: webhookNotificationIntegrationConfig(id.Name(), webhookTestUrl, true, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_body_template", ""),
					resource.TestCheckResourceAttr("snowflake_webhook_notification_integration.test", "webhook_headers.%", "0"),
				),
			},
			// import
			{
				ResourceName:      "snowflake_webhook_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func webhookNotificationIntegrationConfig(name, webhookUrl string, enabled bool, comment string) string {
	commentAttr := ""
	if comment != "" {
		commentAttr = fmt.Sprintf(`  comment = %q`, comment)
	}
	return fmt.Sprintf(`
resource "snowflake_webhook_notification_integration" "test" {
  name        = %q
  enabled     = %t
  webhook_url = %q
  %s
}
`, name, enabled, webhookUrl, commentAttr)
}

func webhookNotificationIntegrationWithOptionsConfig(name, webhookUrl, bodyTemplate string, headers map[string]string) string {
	headerLines := ""
	for k, v := range headers {
		headerLines += fmt.Sprintf("    %q = %q\n", k, v)
	}
	return fmt.Sprintf(`
resource "snowflake_webhook_notification_integration" "test" {
  name                  = %q
  enabled               = true
  webhook_url           = %q
  webhook_body_template = %q
  webhook_headers = {
%s  }
}
`, name, webhookUrl, bodyTemplate, headerLines)
}
