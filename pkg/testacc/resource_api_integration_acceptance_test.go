//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiIntegration_aws(t *testing.T) {
	const dummyAwsPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const dummyAwsOtherPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	const dummyAwsApiRoleArn = "arn:aws:iam::000000000001:/role/test"
	const dummyAwsOtherApiRoleArn = "arn:aws:iam::000000000001:/role/other"

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := "acceptance test"
	key := "12345"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":             config.StringVariable(id.Name()),
			"api_provider":     config.StringVariable("aws_api_gateway"),
			"api_aws_role_arn": config.StringVariable(dummyAwsApiRoleArn),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsPrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsOtherPrefix),
			),
			"api_key": config.StringVariable(key),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}
	m2 := m()
	m2["api_aws_role_arn"] = config.StringVariable(dummyAwsOtherApiRoleArn)
	m2["api_key"] = config.StringVariable("other_key")
	m2["api_blocked_prefixes"] = config.ListVariable()
	m2["api_allowed_prefixes"] = config.ListVariable(
		config.StringVariable(dummyAwsOtherPrefix),
	)
	m2["comment"] = config.StringVariable("different comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_aws_role_arn", dummyAwsApiRoleArn),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.0", dummyAwsPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_blocked_prefixes.0", dummyAwsOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_external_id"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_key", key),
				),
			},
			// change parameters
			{
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_aws_role_arn", dummyAwsOtherApiRoleArn),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.0", dummyAwsOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_blocked_prefixes.#", "0"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "comment", "different comment"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_external_id"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_key", "other_key"),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m2,
				ResourceName:            "snowflake_api_integration.test_aws_int",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAcc_ApiIntegration_azure(t *testing.T) {
	const dummyAzurePrefix = "https://apim-hello-world.azure-api.net/dev"
	const dummyAzureOtherPrefix = "https://apim-hello-world.azure-api.net/prod"
	const dummyAzureTenantId = "00000000-0000-0000-0000-000000000000"
	const dummyAzureOtherTenantId = "11111111-1111-1111-1111-111111111111"
	const dummyAzureAdApplicationId = "22222222-2222-2222-2222-222222222222"
	const dummyAzureOtherAdApplicationId = "33333333-3333-3333-3333-333333333333"

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := "acceptance test"
	key := "12345"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                    config.StringVariable(id.Name()),
			"azure_tenant_id":         config.StringVariable(dummyAzureTenantId),
			"azure_ad_application_id": config.StringVariable(dummyAzureAdApplicationId),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyAzurePrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyAzureOtherPrefix),
			),
			"api_key": config.StringVariable(key),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}
	m2 := m()
	m2["azure_ad_application_id"] = config.StringVariable(dummyAzureOtherAdApplicationId)
	m2["azure_tenant_id"] = config.StringVariable(dummyAzureOtherTenantId)
	m2["api_key"] = config.StringVariable("other_key")
	m2["api_blocked_prefixes"] = config.ListVariable()
	m2["api_allowed_prefixes"] = config.ListVariable(
		config.StringVariable(dummyAzureOtherPrefix),
	)
	m2["comment"] = config.StringVariable("different comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_provider", "azure_api_management"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "azure_tenant_id", dummyAzureTenantId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "azure_ad_application_id", dummyAzureAdApplicationId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_allowed_prefixes.0", dummyAzurePrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_blocked_prefixes.0", dummyAzureOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_multi_tenant_app_name"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_consent_url"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_key", key),
				),
			},
			// change parameters
			{
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_provider", "azure_api_management"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "azure_tenant_id", dummyAzureOtherTenantId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "azure_ad_application_id", dummyAzureOtherAdApplicationId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_allowed_prefixes.0", dummyAzureOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_blocked_prefixes.#", "0"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "comment", "different comment"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_multi_tenant_app_name"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_consent_url"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_key", "other_key"),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m2,
				ResourceName:            "snowflake_api_integration.test_azure_int",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAcc_ApiIntegration_google(t *testing.T) {
	const dummyGooglePrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const dummyGoogleOtherPrefix = "https://gateway-id-123456.uc.gateway.dev/dev"
	const dummyGoogleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const dummyGoogleOtherAudience = "api-gateway-id-666777.apigateway.gcp-project.cloud.goog"

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := "acceptance test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":            config.StringVariable(id.Name()),
			"google_audience": config.StringVariable(dummyGoogleAudience),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyGooglePrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyGoogleOtherPrefix),
			),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}
	m2 := m()
	m2["google_audience"] = config.StringVariable(dummyGoogleOtherAudience)
	m2["api_blocked_prefixes"] = config.ListVariable()
	m2["api_allowed_prefixes"] = config.ListVariable(
		config.StringVariable(dummyGoogleOtherPrefix),
	)
	m2["comment"] = config.StringVariable("different comment")
	m2["api_aws_role_arn"] = config.StringVariable(dummyGoogleOtherAudience)
	m2["api_key"] = config.StringVariable("other_key")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_provider", "google_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "google_audience", dummyGoogleAudience),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_allowed_prefixes.0", dummyGooglePrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_blocked_prefixes.0", dummyGoogleOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "google_audience"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "api_gcp_service_account"),
					resource.TestCheckNoResourceAttr("snowflake_api_integration.test_gcp_int", "api_key"),
				),
			},
			// change parameters
			{
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_provider", "google_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "google_audience", dummyGoogleOtherAudience),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_allowed_prefixes.0", dummyGoogleOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_blocked_prefixes.#", "0"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "comment", "different comment"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "google_audience"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "api_gcp_service_account"),
					resource.TestCheckNoResourceAttr("snowflake_api_integration.test_gcp_int", "api_key"),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m2,
				ResourceName:            "snowflake_api_integration.test_gcp_int",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAcc_ApiIntegration_changeApiProvider(t *testing.T) {
	const dummyAwsPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const dummyAwsOtherPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	const dummyAwsApiRoleArn = "arn:aws:iam::000000000001:/role/test"
	const dummyAzurePrefix = "https://apim-hello-world.azure-api.net/dev"
	const dummyAzureOtherPrefix = "https://apim-hello-world.azure-api.net/prod"
	const dummyAzureTenantId = "00000000-0000-0000-0000-000000000000"
	const dummyAzureAdApplicationId = "22222222-2222-2222-2222-222222222222"

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := "acceptance test"
	key := "12345"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":             config.StringVariable(id.Name()),
			"api_provider":     config.StringVariable("aws_api_gateway"),
			"api_aws_role_arn": config.StringVariable(dummyAwsApiRoleArn),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsPrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsOtherPrefix),
			),
			"api_key": config.StringVariable(key),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}
	m2 := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                    config.StringVariable(id.Name()),
			"azure_tenant_id":         config.StringVariable(dummyAzureTenantId),
			"azure_ad_application_id": config.StringVariable(dummyAzureAdApplicationId),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyAzurePrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyAzureOtherPrefix),
			),
			"api_key": config.StringVariable(key),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_aws_role_arn", dummyAwsApiRoleArn),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_allowed_prefixes.0", dummyAwsPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_blocked_prefixes.0", dummyAwsOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "api_aws_external_id"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_key", key),
				),
			},
			// change parameters
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_provider", "azure_api_management"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "azure_tenant_id", dummyAzureTenantId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "azure_ad_application_id", dummyAzureAdApplicationId),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_allowed_prefixes.0", dummyAzurePrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_blocked_prefixes.0", dummyAzureOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "azure_multi_tenant_app_name"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_change", "azure_consent_url"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_change", "api_key", key),
				),
			},
		},
	})
}
