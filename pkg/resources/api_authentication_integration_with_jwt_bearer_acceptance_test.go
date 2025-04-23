//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_basic(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"enabled":             config.BoolVariable(true),
			"name":                config.StringVariable(id.Name()),
			"oauth_client_id":     config.StringVariable("foo"),
			"oauth_client_secret": config.StringVariable("foo"),
		}
		if complete {
			c["comment"] = config.StringVariable("foo")
			c["oauth_access_token_validity"] = config.IntegerVariable(42)
			c["oauth_authorization_endpoint"] = config.StringVariable("foo")
			c["oauth_client_auth_method"] = config.StringVariable("foo")
			c["oauth_refresh_token_validity"] = config.IntegerVariable(42)
			c["oauth_token_endpoint"] = config.StringVariable("foo")
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables: m(false),
				ResourceName:    "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_id", "foo"),

					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.integration_type", "API_AUTHENTICATION"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.category", "SECURITY"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.comment", ""),

					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.enabled.0.value", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_access_token_validity.0.value", "0"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_refresh_token_validity.0.value", "0"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_client_auth_method.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_token_endpoint.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_allowed_scopes.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_grant.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.parent_integration.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.auth_type.0.value", "OAUTH2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.comment.0.value", ""),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables:   m(true),
				ResourceName:      "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_complete(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_access_token_validity":  config.IntegerVariable(42),
			"oauth_authorization_endpoint": config.StringVariable("foo"),
			"oauth_client_auth_method":     config.StringVariable("foo"),
			"oauth_client_id":              config.StringVariable("foo"),
			"oauth_client_secret":          config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(42),
			"oauth_token_endpoint":         config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_invalidIncomplete(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			`The argument "enabled" is required, but no definition was found.`,
			// this one is trimmed because of inconsistent \n behavior in error message
			`The argument "oauth_assertion_issuer" is required, but no definition`,
			`The argument "oauth_client_id" is required, but no definition was found.`,
			`The argument "oauth_client_secret" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            apiAuthenticationIntegrationWithJwtBearerBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithJwtBearerBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_IdentifierQuotingDiffSuppression(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             apiAuthenticationIntegrationWithJwtBearerBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithJwtBearerBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_jwt_bearer.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_jwt_bearer.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
		},
	})
}

func apiAuthenticationIntegrationWithJwtBearerBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  enabled             = true
  name                = "%s"
  oauth_client_id     = "foo"
  oauth_client_secret = "foo"
  oauth_assertion_issuer = "foo"
}
`, name)
}
