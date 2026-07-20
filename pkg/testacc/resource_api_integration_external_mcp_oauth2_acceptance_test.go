//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiIntegrationExternalMcpOAuth2_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauth2AuthorizationEndpoint, mcpOauth2ClientId, mcpOauth2ClientSecret, mcpOauth2TokenEndpoint)
	withOptionals := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauth2AuthorizationEndpoint, mcpOauth2ClientId, mcpOauth2ClientSecret, mcpOauth2TokenEndpoint).
		WithOauthClientAuthMethod(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity).
		WithApiBlockedPrefixes([]string{mcpBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpOAuth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthClientIdString(mcpOauth2ClientId).
			HasOauthClientSecretString(mcpOauth2ClientSecret).
			HasOauthTokenEndpointString(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpointString(mcpOauth2AuthorizationEndpoint).
			HasOauthClientAuthMethodEmpty().
			HasOauthRefreshTokenValidity(0).
			HasApiAllowedPrefixes(mcpAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationExternalMcpDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(mcpOauth2ClientId).
			HasOauthTokenEndpoint(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(mcpOauth2AuthorizationEndpoint).
			HasOauthClientAuthMethod("").
			HasNoBlockedPrefixes().
			HasComment(""),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProviderType(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthTypeEnum(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(mcpOauth2ClientId).
			HasOauthTokenEndpoint(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(mcpOauth2AuthorizationEndpoint).
			HasAllowedPrefixes(mcpAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationExternalMcpOAuth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasOauthClientIdString(mcpOauth2ClientId).
			HasOauthClientSecretString(mcpOauth2ClientSecret).
			HasOauthTokenEndpointString(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpointString(mcpOauth2AuthorizationEndpoint).
			HasOauthClientAuthMethodString(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
			HasOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity).
			HasApiAllowedPrefixes(mcpAllowedPrefix).
			HasApiBlockedPrefixes(mcpBlockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationExternalMcpDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(mcpOauth2ClientId).
			HasOauthTokenEndpoint(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(mcpOauth2AuthorizationEndpoint).
			HasOauthClientAuthMethod(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)).
			HasOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity).
			HasComment(comment),
		objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProviderType(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthTypeEnum(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(mcpOauth2ClientId).
			HasOauthTokenEndpoint(mcpOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(mcpOauth2AuthorizationEndpoint).
			HasOauthClientAuthMethodEnum(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost).
			HasOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity).
			HasAllowedPrefixes(mcpAllowedPrefix).
			HasBlockedPrefixes(mcpBlockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals
			{
				Config:                  config.FromModels(t, withOptionals),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
			// Update - unset optionals (requires recreate: Snowflake retains unspecified optional auth fields in ALTER)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_client_auth_method", tfjson.ActionDelete, sdk.String(string(sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost)), nil),
						planchecks.ExpectChange(ref, "oauth_refresh_token_validity", tfjson.ActionDelete, sdk.String("7200"), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with optionals
			{
				PreConfig: func() {
					_, err := testClient().ApiIntegration.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

// TestAcc_ApiIntegrationExternalMcpOAuth2_Import verifies that importing a resource created outside Terraform
// produces an in-place update (not destroy-recreate) to sync oauth_client_secret into state, then a noop.
func TestAcc_ApiIntegrationExternalMcpOAuth2_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	oauthClientAuthMethod := sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost
	comment := random.Comment()

	testModel := model.ApiIntegrationExternalMcpOAuth2("t", id.Name(), []string{mcpAllowedPrefix}, true, mcpOauth2AuthorizationEndpoint, mcpOauth2ClientId, mcpOauth2ClientSecret, mcpOauth2TokenEndpoint).
		WithOauthClientAuthMethod(string(oauthClientAuthMethod)).
		WithOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity).
		WithApiBlockedPrefixes([]string{mcpBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					auth := sdk.NewOAuth2McpUserAuthenticationRequest(mcpOauth2ClientId, mcpOauth2ClientSecret, mcpOauth2TokenEndpoint, mcpOauth2AuthorizationEndpoint).
						WithOauthClientAuthMethod(oauthClientAuthMethod).
						WithOauthRefreshTokenValidity(mcpOauth2RefreshTokenValidity)
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: mcpAllowedPrefix}}, true).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: mcpBlockedPrefix}}).
							WithComment(comment).
							WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest(*auth)),
					)
					t.Cleanup(cleanup)
				},
				Config:             config.FromModels(t, testModel),
				ResourceName:       testModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, testModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(testModel.ResourceReference(), "oauth_client_secret", tfjson.ActionUpdate, nil, new(mcpOauth2ClientSecret)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationExternalMcpOAuth2_Import_WrongProviderType(t *testing.T) {
	// Create a DynamicClient MCP integration outside Terraform to use as the import target.
	dynamicClientIntegration, dynamicClientCleanup := testClient().ApiIntegration.CreateMcpDynamicClient(t)
	t.Cleanup(dynamicClientCleanup)

	mcpOAuth2Id := testClient().Ids.RandomAccountObjectIdentifier()
	mcpOAuth2Model := model.ApiIntegrationExternalMcpOAuth2(
		"t", mcpOAuth2Id.Name(),
		[]string{mcpAllowedPrefix},
		true,
		mcpOauth2AuthorizationEndpoint,
		mcpOauth2ClientId,
		mcpOauth2ClientSecret,
		mcpOauth2TokenEndpoint,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationExternalMcpOAuth2),
		Steps: []resource.TestStep{
			// Attempt to import a DynamicClient MCP integration via the OAuth2 resource — expects a user auth type mismatch error.
			{
				Config:        config.FromModels(t, mcpOAuth2Model),
				ResourceName:  mcpOAuth2Model.ResourceReference(),
				ImportState:   true,
				ImportStateId: dynamicClientIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_external_mcp_oauth2"),
			},
		},
	})
}
