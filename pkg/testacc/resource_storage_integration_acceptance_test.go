//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StorageIntegration_Empty_StorageAllowedLocations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Empty_StorageAllowedLocations"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Not enough list items"),
			},
		},
	})
}

func TestAcc_StorageIntegration_AWS_Create_WithExternalId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	awsRoleArn := "arn:aws:iam::000000000001:/role/test"
	awsExternalId := "test-create-external-id-12345"

	configVariables := tfconfig.Variables{
		"name":         tfconfig.StringVariable(id.Name()),
		"aws_role_arn": tfconfig.StringVariable(awsRoleArn),
		"external_id":  tfconfig.StringVariable(awsExternalId),
		"allowed_locations": tfconfig.SetVariable(
			tfconfig.StringVariable("s3://foo/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Create_WithExternalId"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.0.storage_aws_external_id.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.0.storage_aws_external_id.0.value", awsExternalId),
				),
			},
			// external change
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Create_WithExternalId"),
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageIntegrationRequest(id).
						WithSet(*sdk.NewStorageIntegrationSetRequest().
							WithS3Params(*sdk.NewSetS3StorageParamsRequest(awsRoleArn).
								WithStorageAwsExternalId("new-external-id")))
					testClient().StorageIntegration.Alter(t, alterRequest)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.0.storage_aws_external_id.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "describe_output.0.storage_aws_external_id.0.value", awsExternalId),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_AWSObjectACL_Update(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(awsObjectACLSet bool) tfconfig.Variables {
		variables := tfconfig.Variables{
			"name": tfconfig.StringVariable(id.Name()),
			"allowed_locations": tfconfig.SetVariable(
				tfconfig.StringVariable("s3://foo/"),
			),
		}
		if awsObjectACLSet {
			variables["aws_object_acl"] = tfconfig.StringVariable("bucket-owner-full-control")
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/before"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl"),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/after"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", "bucket-owner-full-control"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWSObjectACL_Update/before"),
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
	id := testClient().Ids.RandomAccountObjectIdentifier()
	awsRoleArn := "arn:aws:iam::000000000001:/role/test"
	awsExternalId := "test-external-id-12345"

	configVariables := func(set bool) tfconfig.Variables {
		variables := tfconfig.Variables{
			"name":         tfconfig.StringVariable(id.Name()),
			"aws_role_arn": tfconfig.StringVariable(awsRoleArn),
			"allowed_locations": tfconfig.SetVariable(
				tfconfig.StringVariable("s3://foo/"),
			),
		}
		if set {
			variables["aws_object_acl"] = tfconfig.StringVariable("bucket-owner-full-control")
			variables["external_id"] = tfconfig.StringVariable(awsExternalId)
			variables["comment"] = tfconfig.StringVariable("some comment")
			variables["use_privatelink_endpoint"] = tfconfig.StringVariable(r.BooleanTrue)
			variables["allowed_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable("s3://foo/"),
				tfconfig.StringVariable("s3://bar/"),
			)
			variables["blocked_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable("s3://foo/"),
				tfconfig.StringVariable("s3://bar/"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/unset"),
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
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanTrue),
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
				PreConfig: func() {
					unsetRequest := sdk.NewAlterStorageIntegrationRequest(id).
						WithUnset(*sdk.NewStorageIntegrationUnsetRequest().
							WithStorageAwsExternalId(true).
							WithStorageAwsObjectAcl(true).
							WithStorageBlockedLocations(true))
					testClient().StorageIntegration.Alter(t, unsetRequest)
				},
				ConfigVariables: configVariables(true),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", awsExternalId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", "bucket-owner-full-control"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", r.BooleanFalse),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_role_arn", awsRoleArn),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", "s3://foo/"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
				),
			},
			{
				PreConfig: func() {
					setRequest := sdk.NewAlterStorageIntegrationRequest(id).
						WithSet(*sdk.NewStorageIntegrationSetRequest().
							WithS3Params(*sdk.NewSetS3StorageParamsRequest(awsRoleArn).
								WithStorageAwsExternalId(awsExternalId).
								WithStorageAwsObjectAcl("bucket-owner-full-control").
								WithUsePrivateLinkEndpoint(true)).
							WithStorageBlockedLocations([]sdk.StorageLocation{
								{Path: "s3://external-blocked/"},
							}))
					testClient().StorageIntegration.Alter(t, setRequest)
				},
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/AWS_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_external_id", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_aws_object_acl", ""),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
				),
			},
		},
	})
}

// TODO(SNOW-2356128): Add test for use_privatelink_endpoint
func TestAcc_StorageIntegration_Azure_Update(t *testing.T) {
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	resourceName := "snowflake_storage_integration.test"

	azureTenantId, err := uuid.GenerateUUID()
	require.NoError(t, err)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(set bool) tfconfig.Variables {
		variables := tfconfig.Variables{
			"name":            tfconfig.StringVariable(id.Name()),
			"azure_tenant_id": tfconfig.StringVariable(azureTenantId),
			"allowed_locations": tfconfig.SetVariable(
				tfconfig.StringVariable(azureBucketUrl + "/foo"),
			),
		}
		if set {
			variables["comment"] = tfconfig.StringVariable("some comment")
			variables["allowed_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable(azureBucketUrl+"/foo"),
				tfconfig.StringVariable(azureBucketUrl+"/bar"),
			)
			variables["blocked_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable(azureBucketUrl+"/foo"),
				tfconfig.StringVariable(azureBucketUrl+"/bar"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/foo"),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "describe_output.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse),
				),
			},
			{
				ConfigVariables: configVariables(true),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", "some comment"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/bar"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.1", azureBucketUrl+"/foo"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "2"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.0", azureBucketUrl+"/bar"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.1", azureBucketUrl+"/foo"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse),
				),
			},
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/Azure_Update/unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "enabled", r.BooleanFalse),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "use_privatelink_endpoint", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_allowed_locations.0", azureBucketUrl+"/foo"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "storage_blocked_locations.#", "0"),
					resource.TestCheckResourceAttr("snowflake_storage_integration.test", "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "describe_output.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_GCP_Update(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(set bool) tfconfig.Variables {
		variables := tfconfig.Variables{
			"name": tfconfig.StringVariable(id.Name()),
			"allowed_locations": tfconfig.SetVariable(
				tfconfig.StringVariable("gcs://allowed_foo/"),
			),
		}
		if set {
			variables["comment"] = tfconfig.StringVariable("some comment")
			variables["allowed_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable("gcs://allowed_foo/"),
				tfconfig.StringVariable("gcs://allowed_bar/"),
			)
			variables["blocked_locations"] = tfconfig.SetVariable(
				tfconfig.StringVariable("gcs://blocked_foo/"),
				tfconfig.StringVariable("gcs://blocked_bar/"),
			)
		}
		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(false),
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/unset"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/set"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/unset"),
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
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(id.Name()),
		"allowed_locations": tfconfig.SetVariable(
			tfconfig.StringVariable("gcs://allowed_foo/"),
		),
		"comment": tfconfig.StringVariable("some comment"),
		"blocked_locations": tfconfig.SetVariable(
			tfconfig.StringVariable("gcs://blocked_foo/"),
			tfconfig.StringVariable("gcs://blocked_bar/"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: ConfigurationDirectory("TestAcc_StorageIntegration/GCP_Update/set"),
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

func TestAcc_StorageIntegration_UsePrivateLinkEndpoint_MigrateManuallySetDefaultValue(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	awsRoleArn := "arn:aws:iam::000000000001:/role/test"

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.StorageIntegrationResource))

	s3Model := model.S3StorageIntegration("test", id.Name(), awsRoleArn)
	s3ModelWithPrivatelink := model.S3StorageIntegration("test", id.Name(), awsRoleArn).
		WithUsePrivatelinkEndpoint(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.7.0"),
				Config:            config.FromModels(t, providerModel, s3Model),
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3Model.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasNoUsePrivatelinkEndpoint(),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse)),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(s3Model.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, s3ModelWithPrivatelink),
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3ModelWithPrivatelink.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasUsePrivatelinkEndpointString(r.BooleanDefault),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse)),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_UsePrivateLinkEndpoint_MigrateWithoutValue(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	awsRoleArn := "arn:aws:iam::000000000001:/role/test"

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.StorageIntegrationResource))

	s3Model := model.S3StorageIntegration("test", id.Name(), awsRoleArn)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.7.0"),
				Config:            config.FromModels(t, providerModel, s3Model),
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3Model.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasNoUsePrivatelinkEndpoint(),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse)),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(s3Model.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, s3Model),
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3Model.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasUsePrivatelinkEndpointString(r.BooleanDefault),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse)),
				),
			},
		},
	})
}

func TestAcc_StorageIntegration_UsePrivateLinkEndpoint_MigrateWithoutValue_UpdatedExternally(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	awsRoleArn := "arn:aws:iam::000000000001:/role/test"

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.StorageIntegrationResource))

	s3Model := model.S3StorageIntegration("test", id.Name(), awsRoleArn)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageIntegration),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.7.0"),
				Config:            config.FromModels(t, providerModel, s3Model),
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3Model.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasNoUsePrivatelinkEndpoint(),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanFalse)),
				),
			},
			{
				PreConfig: func() {
					testClient().StorageIntegration.Alter(t, sdk.NewAlterStorageIntegrationRequest(id).WithSet(*sdk.NewStorageIntegrationSetRequest().WithS3Params(*sdk.NewSetS3StorageParamsRequest(awsRoleArn).WithUsePrivateLinkEndpoint(true))))
				},
				Config:                   config.FromModels(t, s3Model),
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(s3Model.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.StorageIntegrationResource(t, s3Model.ResourceReference()).
						HasNameString(id.Name()).
						HasStorageProviderString(string(sdk.StorageProviderS3)).
						HasUsePrivatelinkEndpointString(r.BooleanDefault),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(s3Model.ResourceReference(), "describe_output.0.use_privatelink_endpoint.0.value", r.BooleanTrue)),
				),
			},
		},
	})
}
