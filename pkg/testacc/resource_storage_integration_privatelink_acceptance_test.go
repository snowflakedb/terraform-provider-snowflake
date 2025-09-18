//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StorageIntegration_PrivateLink_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	TestAccPreCheck(t)

	name := testClient().Ids.RandomAccountObjectIdentifier()
	awsRoleArn := "arn:aws:iam::000000000001:/role/test"

	configVariables := config.Variables{
		"name":         config.StringVariable(name.Name()),
		"aws_role_arn": config.StringVariable(awsRoleArn),
		"allowed_locations": config.SetVariable(
			config.StringVariable("s3://foo/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/S3_PrivateLinkEndpoint/set_true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_provider", "S3"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_private_link_endpoint", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
				),
			},
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/S3_PrivateLinkEndpoint/set_false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_provider", "S3"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_private_link_endpoint", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
				),
			},
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/S3_PrivateLinkEndpoint/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_provider", "S3"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_private_link_endpoint", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_Azure_PrivateLink_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	TestAccPreCheck(t)

	name := testClient().Ids.RandomAccountObjectIdentifier()
	azureTenantId := "11111111-2222-3333-4444-555555555555"

	configVariables := config.Variables{
		"name":            config.StringVariable(name.Name()),
		"azure_tenant_id": config.StringVariable(azureTenantId),
		"allowed_locations": config.SetVariable(
			config.StringVariable("azure://myaccount.blob.core.windows.net/mycontainer/path1/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_PrivateLinkEndpoint/set_true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_provider", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "use_private_link_endpoint", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.0", "azure://myaccount.blob.core.windows.net/mycontainer/path1/"),
				),
			},
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_PrivateLinkEndpoint/set_false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_provider", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "use_private_link_endpoint", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.0", "azure://myaccount.blob.core.windows.net/mycontainer/path1/"),
				),
			},
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_PrivateLinkEndpoint/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "name", name.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_provider", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "use_private_link_endpoint", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test_azure", "storage_allowed_locations.0", "azure://myaccount.blob.core.windows.net/mycontainer/path1/"),
				),
			},
		},
	})
}
