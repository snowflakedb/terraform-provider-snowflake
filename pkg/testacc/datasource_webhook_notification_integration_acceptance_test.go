//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_WebhookNotificationIntegrationDatasource_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WebhookNotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: webhookNotificationIntegrationDatasourceConfig(id.Name(), webhookTestUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_webhook_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_webhook_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_webhook_notification_integration.test", "webhook_url", webhookTestUrl),
					resource.TestCheckResourceAttr("data.snowflake_webhook_notification_integration.test", "comment", ""),
					resource.TestCheckResourceAttr("data.snowflake_webhook_notification_integration.test", "webhook_headers.%", "0"),
				),
			},
		},
	})
}

func webhookNotificationIntegrationDatasourceConfig(name, webhookUrl string) string {
	return fmt.Sprintf(`
resource "snowflake_webhook_notification_integration" "test" {
  name        = %q
  enabled     = true
  webhook_url = %q
}

data "snowflake_webhook_notification_integration" "test" {
  name = snowflake_webhook_notification_integration.test.name
}
`, name, webhookUrl)
}
