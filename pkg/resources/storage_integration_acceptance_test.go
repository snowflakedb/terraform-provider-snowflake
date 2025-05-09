//go:build !account_level_tests

package resources_test

import (
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StorageIntegration_Empty_StorageAllowedLocations(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/Empty_StorageAllowedLocations"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Not enough list items"),
			},
		},
	})
}

func TestAcc_StorageIntegration_AWS_Create_WithExternalId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	awsRoleArn := "arn:aws:iam::000000000001:/role/test"
	awsExternalId := "test-create-external-id-12345"

	configVariables := config.Variables{
		"name":         config.StringVariable(id.Name()),
		"aws_role_arn": config.StringVariable(awsRoleArn),
		"external_id":  config.StringVariable(awsExternalId),
		"allowed_locations": config.SetVariable(
			config.StringVariable("s3://foo/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Create_WithExternalId"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_AWSObjectACL_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(awsObjectACLSet bool) config.Variables {
		variables := config.Variables{
			"name": config.StringVariable(id.Name()),
			"allowed_locations": config.SetVariable(
				config.StringVariable("s3://foo/"),
			),
		}
		if awsObjectACLSet {
			variables["aws_object_acl"] = config.StringVariable("bucket-owner-full-control")
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/before"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl"),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/after"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", "bucket-owner-full-control"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/before"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", ""),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_AWS_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	awsRoleArn := "arn:aws:iam::000000000001:/role/test"
	awsExternalId := "test-external-id-12345"

	configVariables := func(set bool) config.Variables {
		variables := config.Variables{
			"name":         config.StringVariable(id.Name()),
			"aws_role_arn": config.StringVariable(awsRoleArn),
			"allowed_locations": config.SetVariable(
				config.StringVariable("s3://foo/"),
			),
		}
		if set {
			variables["aws_object_acl"] = config.StringVariable("bucket-owner-full-control")
			variables["external_id"] = config.StringVariable(awsExternalId)
			variables["comment"] = config.StringVariable("some comment")
			variables["allowed_locations"] = config.SetVariable(
				config.StringVariable("s3://foo/"),
				config.StringVariable("s3://bar/"),
			)
			variables["blocked_locations"] = config.SetVariable(
				config.StringVariable("s3://foo/"),
				config.StringVariable("s3://bar/"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageIntegrationRequest(id).
						WithSet(*sdk.NewStorageIntegrationSetRequest().
							WithComment("altered comment outside of terraform"))
					acc.TestClient().StorageIntegration.Alter(t, alterRequest)
				},
				ConfigVariables: configVariables(true),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://bar/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.1", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.0", "s3://bar/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.1", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", "bucket-owner-full-control"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_Azure_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)

	azureTenantId, err := uuid.GenerateUUID()
	require.NoError(t, err)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(set bool) config.Variables {
		variables := config.Variables{
			"name":            config.StringVariable(id.Name()),
			"azure_tenant_id": config.StringVariable(azureTenantId),
			"allowed_locations": config.SetVariable(
				config.StringVariable(azureBucketUrl + "/foo"),
			),
		}
		if set {
			variables["comment"] = config.StringVariable("some comment")
			variables["allowed_locations"] = config.SetVariable(
				config.StringVariable(azureBucketUrl+"/foo"),
				config.StringVariable(azureBucketUrl+"/bar"),
			)
			variables["blocked_locations"] = config.SetVariable(
				config.StringVariable(azureBucketUrl+"/foo"),
				config.StringVariable(azureBucketUrl+"/bar"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/foo"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/bar"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.1", azureBucketUrl+"/foo"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.0", azureBucketUrl+"/bar"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.1", azureBucketUrl+"/foo"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/foo"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_GCP_Update(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(set bool) config.Variables {
		variables := config.Variables{
			"name": config.StringVariable(id.Name()),
			"allowed_locations": config.SetVariable(
				config.StringVariable("gcs://allowed_foo/"),
			),
		}
		if set {
			variables["comment"] = config.StringVariable("some comment")
			variables["allowed_locations"] = config.SetVariable(
				config.StringVariable("gcs://allowed_foo/"),
				config.StringVariable("gcs://allowed_bar/"),
			)
			variables["blocked_locations"] = config.SetVariable(
				config.StringVariable("gcs://blocked_foo/"),
				config.StringVariable("gcs://blocked_bar/"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "gcs://allowed_foo/"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "gcs://allowed_bar/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.1", "gcs://allowed_foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.0", "gcs://blocked_bar/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.1", "gcs://blocked_foo/"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "gcs://allowed_foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_BlockedLocations_issue2985(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(id.Name()),
		"allowed_locations": config.SetVariable(
			config.StringVariable("gcs://allowed_foo/"),
		),
		"comment": config.StringVariable("some comment"),
		"blocked_locations": config.SetVariable(
			config.StringVariable("gcs://blocked_foo/"),
			config.StringVariable("gcs://blocked_bar/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/set"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "gcs://allowed_foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.0", "gcs://blocked_bar/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.1", "gcs://blocked_foo/"),
				),
			},
		},
	})
}
