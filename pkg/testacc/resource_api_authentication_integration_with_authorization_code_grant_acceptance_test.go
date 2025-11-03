//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basic := model.ApiAuthenticationIntegrationWithAuthorizationCodeGrant("test", id.Name(), true, "test_client_id", "test_client_secret")

	complete := model.ApiAuthenticationIntegrationWithAuthorizationCodeGrant("test", id.Name(), true, "test_client_id", "test_client_secret").
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthRefreshTokenValidity(12345).
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment("").
			HasCreatedOnNotEmpty(),

		resourceassert.ApiAuthenticationIntegrationWithAuthorizationCodeGrantResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasCommentString(""),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, basic.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment("").
			HasCreatedOnNotEmpty(),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.enabled.0.value", "true")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_access_token_validity.0.value", "0")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.parent_integration.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.auth_type.0.value", "OAUTH2")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment.0.value", "")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment).
			HasCreatedOnNotEmpty(),

		resourceassert.ApiAuthenticationIntegrationWithAuthorizationCodeGrantResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasCommentString(comment).
			HasOauthAccessTokenValidityString("42").
			HasOauthRefreshTokenValidityString("12345").
			HasOauthAllowedScopesLen(1).
			HasOauthAllowedScopesElem(0, "foo"),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment).
			HasCreatedOnNotEmpty(),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", "true")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_validity.0.value", "42")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "12345")),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.parent_integration.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.auth_type.0.value", "OAUTH2")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"oauth_access_token_validity",
					"oauth_authorization_endpoint",
					"oauth_client_auth_method",
					"oauth_client_secret",
					"oauth_refresh_token_validity",
					"oauth_token_endpoint",
				},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"oauth_access_token_validity",
					"oauth_authorization_endpoint",
					"oauth_client_auth_method",
					"oauth_client_secret",
					"oauth_refresh_token_validity",
					"oauth_token_endpoint",
				},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().SecurityIntegration.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(t, sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).
						WithSet(*sdk.NewApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSetRequest().
							WithComment(comment),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure api integration is destroyed before the next step
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
				Check: assertThat(t,
					invokeactionassert.ApiIntegrationDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	complete := model.ApiAuthenticationIntegrationWithAuthorizationCodeGrant("test", id.Name(), true, "test_client_id", "test_client_secret").
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthAuthorizationEndpoint("https://example.com").
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// Create - with all optionals (including force-new fields)
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					objectassert.SecurityIntegration(t, id).
						HasName(id.Name()).
						HasIntegrationType("API_AUTHENTICATION").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),

					resourceassert.ApiAuthenticationIntegrationWithAuthorizationCodeGrantResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasEnabledString("true").
						HasOauthClientIdString("test_client_id").
						HasOauthClientSecretString("test_client_secret").
						HasCommentString(comment).
						HasOauthAccessTokenValidityString("42").
						HasOauthAuthorizationEndpointString("https://example.com").
						HasOauthClientAuthMethodString(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
						HasOauthRefreshTokenValidityString("12345").
						HasOauthTokenEndpointString("https://example.com").
						HasOauthAllowedScopesLen(1).
						HasOauthAllowedScopesElem(0, "foo"),

					resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
						HasName(id.Name()).
						HasIntegrationType("API_AUTHENTICATION").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment).
						HasCreatedOnNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", "true")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_validity.0.value", "42")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "12345")),
					assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.parent_integration.0.value", "")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.auth_type.0.value", "OAUTH2")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
				),
			},
			// Import - with all optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"oauth_access_token_validity",
					"oauth_authorization_endpoint",
					"oauth_client_auth_method",
					"oauth_client_secret",
					"oauth_refresh_token_validity",
					"oauth_token_endpoint",
				},
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
