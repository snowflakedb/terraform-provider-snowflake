//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

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
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ApiIntegrationGitRepositoryOauth2_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const authorizationEndpoint = "https://auth.example.com/authorize"
	const tokenEndpoint = "https://auth.example.com/token"
	const clientId = "oauth-client-id-123"
	const clientSecret = "oauth-client-secret-456"
	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"

	comment := random.Comment()
	externalComment := random.Comment()

	const accessTokenValidity = 3600
	const refreshTokenValidity = 86400
	const oauthUsername = "test_user"

	const externalAuthorizationEndpoint = "https://different.example.com/authorize"
	const externalTokenEndpoint = "https://different.example.com/token"
	const externalClientId = "different-client-id"

	basic := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true, authorizationEndpoint, clientId, clientSecret, tokenEndpoint)
	withOptionals := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true, authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)
	withAllOptionals := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true, authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment).
		WithOauthAccessTokenValidity(accessTokenValidity).
		WithOauthRefreshTokenValidity(refreshTokenValidity).
		WithOauthAllowedScopes([]string{string(sdk.ApiIntegrationOauthAllowedScopeReadApi)}).
		WithOauthUsername(oauthUsername)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(authorizationEndpoint).
			HasOauthTokenEndpointString(tokenEndpoint).
			HasOauthClientIdString(clientId).
			HasOauthClientSecret(clientSecret).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasNoOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(allowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(authorizationEndpoint).
			HasOauthTokenEndpointString(tokenEndpoint).
			HasOauthClientIdString(clientId).
			HasOauthClientSecret(clientSecret).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasNoOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
	}

	assertWithAllOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(authorizationEndpoint).
			HasOauthTokenEndpointString(tokenEndpoint).
			HasOauthClientIdString(clientId).
			HasOauthClientSecret(clientSecret).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(oauthUsername),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(oauthUsername).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(sdk.ApiIntegrationOauthAllowedScopeReadApi).
			HasOauthUsername(oauthUsername).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryOauth2),
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
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes to comment
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
			// Update - external changes to non-ForceNew fields (enabled, api_allowed_prefixes, api_blocked_prefixes)
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().
							WithEnabled(false).
							WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/other-org/"}}).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/blocked-externally/"}}),
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
			// Plan check - external changes to ForceNew fields (oauth_authorization_endpoint, oauth_token_endpoint, oauth_client_id)
			{
				PreConfig: func() {
					testClient().ApiIntegration.DropApiIntegrationFunc(t, id)()
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().
								WithApiUserAuthentication(*sdk.NewOAuth2GitUserAuthenticationRequest(
									externalAuthorizationEndpoint, externalTokenEndpoint, externalClientId, clientSecret,
								))),
					)
					t.Cleanup(cleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_authorization_endpoint", tfjson.ActionDelete, new(externalAuthorizationEndpoint), new(authorizationEndpoint)),
						planchecks.ExpectChange(ref, "oauth_token_endpoint", tfjson.ActionDelete, new(externalTokenEndpoint), new(tokenEndpoint)),
						planchecks.ExpectChange(ref, "oauth_client_id", tfjson.ActionDelete, new(externalClientId), new(clientId)),
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
			// Create - with all optionals
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
				Config: config.FromModels(t, withAllOptionals),
				Check:  assertThat(t, assertWithAllOptionals...),
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryOauth2_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const authorizationEndpoint = "https://auth.example.com/authorize"
	const tokenEndpoint = "https://auth.example.com/token"
	const clientId = "oauth-client-id-123"
	const clientSecret = "oauth-client-secret-456"
	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"

	comment := random.Comment()

	const accessTokenValidity = 3600
	const refreshTokenValidity = 86400
	const oauthUsername = "test_user"

	const externalAuthorizationEndpoint = "https://different.example.com/authorize"
	const externalTokenEndpoint = "https://different.example.com/token"
	const externalClientId = "different-client-id"

	allAttributes := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true, authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment).
		WithOauthAccessTokenValidity(accessTokenValidity).
		WithOauthRefreshTokenValidity(refreshTokenValidity).
		WithOauthAllowedScopes([]string{string(sdk.ApiIntegrationOauthAllowedScopeReadApi)}).
		WithOauthUsername(oauthUsername)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(authorizationEndpoint).
			HasOauthTokenEndpointString(tokenEndpoint).
			HasOauthClientIdString(clientId).
			HasOauthClientSecret(clientSecret).
			HasApiAllowedPrefixes(allowedPrefix).
			HasApiBlockedPrefixes(blockedPrefix).
			HasCommentString(comment).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(oauthUsername),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(oauthUsername).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(clientId).
			HasOauthTokenEndpoint(tokenEndpoint).
			HasOauthAuthorizationEndpoint(authorizationEndpoint).
			HasOauthAccessTokenValidity(accessTokenValidity).
			HasOauthRefreshTokenValidity(refreshTokenValidity).
			HasOauthAllowedScopes(sdk.ApiIntegrationOauthAllowedScopeReadApi).
			HasOauthUsername(oauthUsername).
			HasAllowedPrefixes(allowedPrefix).
			HasBlockedPrefixes(blockedPrefix).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryOauth2),
		Steps: []resource.TestStep{
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Plan check - external changes to ForceNew fields (oauth_authorization_endpoint, oauth_token_endpoint, oauth_client_id)
			{
				PreConfig: func() {
					testClient().ApiIntegration.DropApiIntegrationFunc(t, id)()
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().
								WithApiUserAuthentication(*sdk.NewOAuth2GitUserAuthenticationRequest(
									externalAuthorizationEndpoint, externalTokenEndpoint, externalClientId, clientSecret,
								))),
					)
					t.Cleanup(cleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_authorization_endpoint", tfjson.ActionDelete, sdk.String(externalAuthorizationEndpoint), sdk.String(authorizationEndpoint)),
						planchecks.ExpectChange(ref, "oauth_token_endpoint", tfjson.ActionDelete, sdk.String(externalTokenEndpoint), sdk.String(tokenEndpoint)),
						planchecks.ExpectChange(ref, "oauth_client_id", tfjson.ActionDelete, sdk.String(externalClientId), sdk.String(clientId)),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import with all attributes (oauth_client_secret is write-only and not readable after import)
			{
				Config:                  config.FromModels(t, allAttributes),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryOauth2_Import_WrongAuthType(t *testing.T) {
	// Create a GitHub App integration outside Terraform to use as the import target.
	githubAppIntegration, githubAppCleanup := testClient().ApiIntegration.CreateGitGithubApp(t)
	t.Cleanup(githubAppCleanup)

	oauth2Id := testClient().Ids.RandomAccountObjectIdentifier()
	oauth2Model := model.ApiIntegrationGitRepositoryOauth2("t", oauth2Id.Name(),
		[]string{"https://github.com/my-org/"},
		true,
		"https://auth.example.com/authorize",
		"oauth-client-id-123",
		"oauth-client-secret-456",
		"https://auth.example.com/token",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryOauth2),
		Steps: []resource.TestStep{
			// Attempt to import a GitHub App integration via the OAuth2 resource — expects auth type mismatch error.
			{
				Config:        config.FromModels(t, oauth2Model),
				ResourceName:  oauth2Model.ResourceReference(),
				ImportState:   true,
				ImportStateId: githubAppIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_git_repository_oauth2"),
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryOauth2_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const allowedPrefix = "https://github.com/my-org/"
	const blockedPrefix = "https://github.com/my-org/blocked/"
	const authorizationEndpoint = "https://auth.example.com/authorize"
	const tokenEndpoint = "https://auth.example.com/token"
	const clientId = "oauth-client-id-123"
	clientSecret := random.AlphanumericN(20)
	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true,
		authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithApiBlockedPrefixes([]string{blockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryOauth2),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: allowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: blockedPrefix}}).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().
								WithApiUserAuthentication(*sdk.NewOAuth2GitUserAuthenticationRequest(
									authorizationEndpoint, tokenEndpoint, clientId, clientSecret,
								))),
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
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(testModel.ResourceReference(), "oauth_client_secret", tfjson.ActionDelete, nil, sdk.String(clientSecret)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryOauth2_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	const authorizationEndpoint = "https://auth.example.com/authorize"
	const tokenEndpoint = "https://auth.example.com/token"
	const clientId = "oauth-client-id-123"
	const clientSecret = "oauth-client-secret-456"
	const allowedPrefix = "https://github.com/my-org/"

	invalidScope := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true,
		authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithOauthAllowedScopes([]string{"INVALID_SCOPE"})

	invalidAccessTokenValidity := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true,
		authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithOauthAccessTokenValidity(0)

	invalidRefreshTokenValidity := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{allowedPrefix}, true,
		authorizationEndpoint, clientId, clientSecret, tokenEndpoint).
		WithOauthRefreshTokenValidity(0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidScope),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid api integration oauth allowed scope"),
			},
			{
				Config:      config.FromModels(t, invalidAccessTokenValidity),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected oauth_access_token_validity to be at least"),
			},
			{
				Config:      config.FromModels(t, invalidRefreshTokenValidity),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected oauth_refresh_token_validity to be at least"),
			},
		},
	})
}
