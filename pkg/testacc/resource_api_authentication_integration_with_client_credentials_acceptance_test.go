//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

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

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basic := model.ApiAuthenticationIntegrationWithClientCredentials("test", id.Name(), true, "test_client_id", "test_client_secret")

	complete := model.ApiAuthenticationIntegrationWithClientCredentials("test", id.Name(), true, "test_client_id", "test_client_secret").
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

		resourceassert.ApiAuthenticationIntegrationWithClientCredentialsResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasCommentString("").
			HasOauthAccessTokenValidityString("-1").
			HasNoOauthClientAuthMethod().
			// TODO(SNOW-2457144): HasNoOauthRefreshTokenValidity().
			HasNoOauthTokenEndpoint().
			HasOauthAllowedScopesLen(0),

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
		// TODO(SNOW-2457144): assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "0")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		// TODO(SNOW-2457144): assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials)),
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

		resourceassert.ApiAuthenticationIntegrationWithClientCredentialsResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasCommentString(comment).
			HasOauthAccessTokenValidityString("42").
			HasNoOauthClientAuthMethod().
			HasOauthRefreshTokenValidityString("12345").
			HasNoOauthTokenEndpoint().
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
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.parent_integration.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.auth_type.0.value", "OAUTH2")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
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
					testClient().SecurityIntegration.AlterApiAuthenticationWithClientCredentialsFlow(t, sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).
						WithSet(*sdk.NewApiAuthenticationWithClientCredentialsFlowIntegrationSetRequest().
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
			// Create - with optionals (from scratch via taint)
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	complete := model.ApiAuthenticationIntegrationWithClientCredentials("test", id.Name(), true, "test_client_id", "test_client_secret").
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
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

					resourceassert.ApiAuthenticationIntegrationWithClientCredentialsResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasEnabledString("true").
						HasOauthClientIdString("test_client_id").
						HasOauthClientSecretString("test_client_secret").
						HasCommentString(comment).
						HasOauthAccessTokenValidityString("42").
						HasOauthClientAuthMethodString(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
						HasOauthRefreshTokenValidityString("12345").
						HasOauthTokenEndpointString("https://example.com").
						HasOauthAllowedScopesString("[foo]"),

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
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials)),
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
					"oauth_client_auth_method",
					"oauth_client_secret",
					"oauth_refresh_token_validity",
					"oauth_token_endpoint",
				},
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_invalidIncomplete(t *testing.T) {
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            apiAuthenticationIntegrationWithClientCredentialsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithClientCredentialsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             apiAuthenticationIntegrationWithClientCredentialsBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   apiAuthenticationIntegrationWithClientCredentialsBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_client_credentials.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_api_authentication_integration_with_client_credentials.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "id", id.Name()),
				),
			},
		},
	})
}

func apiAuthenticationIntegrationWithClientCredentialsBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  enabled             = true
  name                = "%s"
  oauth_client_id     = "foo"
  oauth_client_secret = "foo"
}
`, name)
}
