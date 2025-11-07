//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
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
			c["oauth_authorization_endpoint"] = config.StringVariable("https://example.com")
			c["oauth_client_auth_method"] = config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))
			c["oauth_refresh_token_validity"] = config.IntegerVariable(12345)
			c["oauth_token_endpoint"] = config.StringVariable("https://example.com")
			c["oauth_allowed_scopes"] = config.SetVariable(config.StringVariable("foo"))
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),

					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_access_token_validity.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_refresh_token_validity.0.value", "0"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_client_auth_method.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_authorization_endpoint.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_token_endpoint.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_allowed_scopes.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.comment.0.value", ""),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables: m(false),
				ResourceName:    "snowflake_api_authentication_integration_with_authorization_code_grant.test",
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
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_client_id.0.value"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_client_auth_method.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_token_endpoint.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_allowed_scopes.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.parent_integration.0.value", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.auth_type.0.value", "OAUTH2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.comment.0.value", ""),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_allowed_scopes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_allowed_scopes.0", "foo"),

					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_access_token_validity.0.value", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_refresh_token_validity.0.value", "12345"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_token_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.comment.0.value", "foo")),
			},
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables:         m(true),
				ResourceName:            "snowflake_api_authentication_integration_with_authorization_code_grant.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
			// unset
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "-1"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity"),
					resource.TestCheckNoResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint"),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_complete(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_access_token_validity":  config.IntegerVariable(42),
			"oauth_authorization_endpoint": config.StringVariable("https://example.com"),
			"oauth_client_auth_method":     config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
			"oauth_client_id":              config.StringVariable("foo"),
			"oauth_client_secret":          config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(12345),
			"oauth_token_endpoint":         config.StringVariable("https://example.com"),
			"oauth_allowed_scopes":         config.SetVariable(config.StringVariable("foo")),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_allowed_scopes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_allowed_scopes.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
				),
			},
			{
				ConfigDirectory:         ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_api_authentication_integration_with_authorization_code_grant.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_invalidIncomplete(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			`The argument "oauth_client_secret" is required, but no definition was found.`,
			`The argument "oauth_client_id" is required, but no definition was found.`,
			`The argument "enabled" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            apiAuthenticationIntegrationWithAuthorizationCodeGrantBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithAuthorizationCodeGrantBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             apiAuthenticationIntegrationWithAuthorizationCodeGrantBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithAuthorizationCodeGrantBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_authorization_code_grant.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_authorization_code_grant.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "id", id.Name()),
				),
			},
		},
	})
}

func apiAuthenticationIntegrationWithAuthorizationCodeGrantBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_api_authentication_integration_with_authorization_code_grant" "test" {
  enabled             = true
  name                = "%s"
  oauth_client_id     = "foo"
  oauth_client_secret = "foo"
}
`, name)
}
