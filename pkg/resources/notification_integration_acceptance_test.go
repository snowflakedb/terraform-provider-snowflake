//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NotificationIntegration_AutoGoogle(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const gcpOtherPubsubSubscriptionName = "projects/project-1234/subscriptions/other"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: googleAutoConfig(id, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// change parameters
			{
				Config: googleAutoConfig(id, gcpOtherPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpOtherPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NotificationIntegration_AutoAzure(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const azureStorageQueuePrimaryUri = "azure://great-bucket/great-path/"
	const azureOtherStorageQueuePrimaryUri = "azure://great-bucket/other-great-path/"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureOtherTenantId = "11111111-1111-1111-1111-111111111111"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: azureAutoConfig(id, azureStorageQueuePrimaryUri, azureTenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", azureStorageQueuePrimaryUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
			// change parameters
			{
				Config: azureAutoConfig(id, azureOtherStorageQueuePrimaryUri, azureOtherTenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", azureOtherStorageQueuePrimaryUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", azureOtherTenantId),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
				// it is not returned in DESCRIBE for azure automated data load
				ImportStateVerifyIgnore: []string{"azure_tenant_id"},
			},
		},
	})
}

func TestAcc_NotificationIntegration_PushAmazon(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsOtherSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:OtherTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const awsOtherSnsRoleArn = "arn:aws:iam::000000000001:/role/other"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: amazonPushConfig(id, awsSnsTopicArn, awsSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
			// change parameters
			{
				Config: amazonPushConfig(id, awsOtherSnsTopicArn, awsOtherSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsOtherSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsOtherSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NotificationIntegration_changeNotificationProvider(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: googleAutoConfig(id, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// change provider to AWS
			{
				Config: amazonPushConfig(id, awsSnsTopicArn, awsSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
		},
	})
}

// TODO [SNOW-1017802]: handle after "create and describe notification integration - push google" test passes
func TestAcc_NotificationIntegration_PushGoogle(t *testing.T) {
	t.Skip("Skipping because can't be currently created. Check 'create and describe notification integration - push google' test in the SDK.")
}

// TODO [SNOW-1017802]: handle after "create and describe notification integration - push azure" test passes
// TODO [SNOW-1348345]: handle after it's added to the resource
func TestAcc_NotificationIntegration_PushAzure(t *testing.T) {
	t.Skip("Skipping because can't be currently created. Check 'create and describe notification integration - push azure' test in the SDK.")
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2501
func TestAcc_NotificationIntegration_migrateFromVersion085(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.85.0"),
				Config:            googleAutoConfig(id, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   googleAutoConfigWithoutDirection(id, gcpPubsubSubscriptionName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
		},
	})
}

func TestAcc_NotificationIntegration_migrateFromVersion085_explicitType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.85.0"),
				Config:            googleAutoConfigWithExplicitType(id, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   googleAutoConfig(id, gcpPubsubSubscriptionName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", id.Name()),
				),
			},
		},
	})
}

func googleAutoConfig(id sdk.AccountObjectIdentifier, gcpPubsubSubscriptionName string) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  name                            = "%[1]s"
  notification_provider           = "GCP_PUBSUB"
  gcp_pubsub_subscription_name    = "%[2]s"
  direction                       = "INBOUND"
}
`, id.Name(), gcpPubsubSubscriptionName)
}

func googleAutoConfigWithoutDirection(id sdk.AccountObjectIdentifier, gcpPubsubSubscriptionName string) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  name                            = "%[1]s"
  notification_provider           = "GCP_PUBSUB"
  gcp_pubsub_subscription_name    = "%[2]s"
}
`, id.Name(), gcpPubsubSubscriptionName)
}

func googleAutoConfigWithExplicitType(id sdk.AccountObjectIdentifier, gcpPubsubSubscriptionName string) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  type                            = "QUEUE"
  name                            = "%[1]s"
  notification_provider           = "GCP_PUBSUB"
  gcp_pubsub_subscription_name    = "%[2]s"
  direction                       = "INBOUND"
}
`, id.Name(), gcpPubsubSubscriptionName)
}

func azureAutoConfig(id sdk.AccountObjectIdentifier, azureStorageQueuePrimaryUri string, azureTenantId string) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  name                            = "%[1]s"
  notification_provider			  = "AZURE_STORAGE_QUEUE"
  azure_storage_queue_primary_uri = "%[2]s"
  azure_tenant_id                 = "%[3]s"
}
`, id.Name(), azureStorageQueuePrimaryUri, azureTenantId)
}

func amazonPushConfig(id sdk.AccountObjectIdentifier, awsSnsTopicArn string, awsSnsRoleArn string) string {
	return fmt.Sprintf(`
resource "snowflake_notification_integration" "test" {
  name                            = "%[1]s"
  notification_provider           = "AWS_SNS"
  aws_sns_topic_arn               = "%[2]s"
  aws_sns_role_arn                = "%[3]s"
  direction                       = "OUTBOUND"
}
`, id.Name(), awsSnsTopicArn, awsSnsRoleArn)
}
