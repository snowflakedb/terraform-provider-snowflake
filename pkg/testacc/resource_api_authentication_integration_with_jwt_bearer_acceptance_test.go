//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_BasicUseCase(t *testing.T) {
	// Schema analysis (from pkg/resources/api_authentication_integration_with_jwt_bearer.go):
	// - name: ForceNew: true (cannot be renamed)
	// - api_provider: ForceNew: true (cannot be changed)
	// - api_key: NOT force-new (can be updated)
	// - api_secret: NOT force-new (can be updated)
	// - enabled: Optional, NOT force-new
	// - comment: Optional, NOT force-new
	// - oauth_token_endpoint: ForceNewIfChangeToEmptyString (exclude from complete to avoid force-new)
	// - oauth_authorization_endpoint: ForceNewIfChangeToEmptyString (exclude from complete to avoid force-new)
	// - oauth_client_auth_method: ForceNewIfChangeToEmptyString (exclude from complete to avoid force-new)
	// Result: Use same identifiers for basic/complete (name, api_provider are force-new), exclude force-new fields from complete model

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basic := model.ApiAuthenticationIntegrationWithJwtBearer("test", id.Name()).
		WithEnabled(true).
		WithOauthClientId("test_client_id").
		WithOauthClientSecret("test_client_secret").
		WithOauthAssertionIssuer("test_issuer")

	complete := model.ApiAuthenticationIntegrationWithJwtBearer("test", id.Name()).
		WithEnabled(true).
		WithOauthClientId("test_client_id").
		WithOauthClientSecret("test_client_secret").
		WithOauthAssertionIssuer("test_issuer").
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment("").
			HasCreatedOnNotEmpty(),

		resourceassert.ApiAuthenticationIntegrationWithJwtBearerResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasOauthAssertionIssuerString("test_issuer").
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
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "0")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer)),
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

		resourceassert.ApiAuthenticationIntegrationWithJwtBearerResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasEnabledString("true").
			HasOauthClientIdString("test_client_id").
			HasOauthClientSecretString("test_client_secret").
			HasOauthAssertionIssuerString("test_issuer").
			HasCommentString(comment),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("API_AUTHENTICATION").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment).
			HasCreatedOnNotEmpty(),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", "true")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_validity.0.value", "0")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "0")),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_auth_method.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.parent_integration.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.auth_type.0.value", "OAUTH2")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
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
					testClient().SecurityIntegration.AlterApiAuthenticationWithJwtBearerFlow(t, sdk.NewAlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id).
						WithSet(*sdk.NewApiAuthenticationWithJwtBearerFlowIntegrationSetRequest().
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

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_CompleteUseCase(t *testing.T) {
	// CompleteUseCase test for resources with ForceNewIfChangeToEmptyString fields
	// This test creates the resource with all optional fields set from the beginning
	// to avoid force-new conflicts between basic and complete models

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	complete := model.ApiAuthenticationIntegrationWithJwtBearer("test", id.Name()).
		WithEnabled(true).
		WithOauthClientId("test_client_id").
		WithOauthClientSecret("test_client_secret").
		WithOauthAssertionIssuer("test_issuer").
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
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
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

					resourceassert.ApiAuthenticationIntegrationWithJwtBearerResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasEnabledString("true").
						HasOauthClientIdString("test_client_id").
						HasOauthClientSecretString("test_client_secret").
						HasOauthAssertionIssuerString("test_issuer").
						HasCommentString(comment).
						HasOauthAccessTokenValidityString("42").
						HasOauthAuthorizationEndpointString("https://example.com").
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
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value", "https://example.com")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allowed_scopes.0.value", "[foo]")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer)),
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

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_basic(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

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
			c["oauth_authorization_endpoint"] = config.StringVariable("foo")
			c["oauth_client_auth_method"] = config.StringVariable("foo")
			c["oauth_refresh_token_validity"] = config.IntegerVariable(42)
			c["oauth_token_endpoint"] = config.StringVariable("foo")
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
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
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.oauth_client_id.0.value"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables:   m(true),
				ResourceName:      "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
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

	id := testClient().Ids.RandomAccountObjectIdentifier()
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")

	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            apiAuthenticationIntegrationWithJwtBearerBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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

	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiAuthenticationIntegrationWithJwtBearer),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             apiAuthenticationIntegrationWithJwtBearerBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
