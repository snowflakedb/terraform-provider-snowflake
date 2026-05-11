//go:build non_account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NotificationIntegration_AWS_SQS(t *testing.T) {
	// TODO [SNOW-1017580]: Use real SQS queue ARN from test environment
	sqsArn := "arn:aws:sqs:us-east-2:123456789012:test-queue"
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create AWS_SQS notification integration
			{
				Config: awsSqsNotificationIntegrationConfig(name, sqsArn, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SQS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sqs_arn", sqsArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sqs_iam_user_arn"),
				),
			},
			// Import
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update enabled
			{
				Config: awsSqsNotificationIntegrationConfig(name, sqsArn, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "false"),
				),
			},
		},
	})
}

func awsSqsNotificationIntegrationConfig(name string, sqsArn string, enabled bool) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  name                  = "%s"
  notification_provider = "AWS_SQS"
  aws_sqs_arn           = "%s"
  enabled               = %t
  comment               = "Terraform acceptance test for AWS_SQS"
}
`, name, sqsArn, enabled)
}
