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
)

func TestAcc_ApiIntegrationGitRepositoryOauth2_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true, gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint)
	withOptionals := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true, gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(gitOauth2AuthorizationEndpoint).
			HasOauthTokenEndpointString(gitOauth2TokenEndpoint).
			HasOauthClientIdString(gitOauth2ClientId).
			HasOauthClientSecret(gitOauth2ClientSecret).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasNoOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(gitAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(gitAllowedPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(gitOauth2AuthorizationEndpoint).
			HasOauthTokenEndpointString(gitOauth2TokenEndpoint).
			HasOauthClientIdString(gitOauth2ClientId).
			HasOauthClientSecret(gitOauth2ClientSecret).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasNoOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(0).
			HasOauthRefreshTokenValidity(0).
			HasOauthAllowedScopes().
			HasOauthUsername("").
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
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
						sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest(
								*sdk.NewOAuth2GitUserAuthenticationRequest(
									gitOauth2ExternalAuthorizationEndpoint, gitOauth2ExternalTokenEndpoint, gitOauth2ExternalClientId, gitOauth2ClientSecret,
								))),
					)
					t.Cleanup(cleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_authorization_endpoint", tfjson.ActionDelete, new(gitOauth2ExternalAuthorizationEndpoint), new(gitOauth2AuthorizationEndpoint)),
						planchecks.ExpectChange(ref, "oauth_token_endpoint", tfjson.ActionDelete, new(gitOauth2ExternalTokenEndpoint), new(gitOauth2TokenEndpoint)),
						planchecks.ExpectChange(ref, "oauth_client_id", tfjson.ActionDelete, new(gitOauth2ExternalClientId), new(gitOauth2ClientId)),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryOauth2_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()

	allAttributes := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true, gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment).
		WithOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
		WithOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
		WithOauthAllowedScopes([]string{string(sdk.ApiIntegrationOauthAllowedScopeReadApi)}).
		WithOauthUsername(gitOauth2Username)

	allAttributesUpdated := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true, gitOauth2UpdatedAuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2UpdatedTokenEndpoint).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment).
		WithOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
		WithOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
		WithOauthAllowedScopes([]string{string(sdk.ApiIntegrationOauthAllowedScopeReadApi)}).
		WithOauthUsername(gitOauth2Username)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(gitOauth2AuthorizationEndpoint).
			HasOauthTokenEndpointString(gitOauth2TokenEndpoint).
			HasOauthClientIdString(gitOauth2ClientId).
			HasOauthClientSecret(gitOauth2ClientSecret).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasCommentString(comment).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(gitOauth2Username),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(gitOauth2Username).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2TokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2AuthorizationEndpoint).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(sdk.ApiIntegrationOauthAllowedScopeReadApi).
			HasOauthUsername(gitOauth2Username).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
	}

	updatedCompleteAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryOauth2Resource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString("true").
			HasOauthAuthorizationEndpointString(gitOauth2UpdatedAuthorizationEndpoint).
			HasOauthTokenEndpointString(gitOauth2UpdatedTokenEndpoint).
			HasOauthClientIdString(gitOauth2ClientId).
			HasOauthClientSecret(gitOauth2ClientSecret).
			HasApiAllowedPrefixes(gitAllowedPrefix).
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasCommentString(comment).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(gitOauth2Username),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitHttpsApiDescribeOutput(t, ref).
			HasEnabled(true).
			HasUserAuthType(string(sdk.ApiIntegrationUserAuthTypeOauth2)).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2UpdatedTokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2UpdatedAuthorizationEndpoint).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(string(sdk.ApiIntegrationOauthAllowedScopeReadApi)).
			HasOauthUsername(gitOauth2Username).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthClientId(gitOauth2ClientId).
			HasOauthTokenEndpoint(gitOauth2UpdatedTokenEndpoint).
			HasOauthAuthorizationEndpoint(gitOauth2UpdatedAuthorizationEndpoint).
			HasOauthAccessTokenValidity(gitOauth2AccessTokenValidity).
			HasOauthRefreshTokenValidity(gitOauth2RefreshTokenValidity).
			HasOauthAllowedScopes(sdk.ApiIntegrationOauthAllowedScopeReadApi).
			HasOauthUsername(gitOauth2Username).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
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
						sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest(
								*sdk.NewOAuth2GitUserAuthenticationRequest(
									gitOauth2ExternalAuthorizationEndpoint, gitOauth2ExternalTokenEndpoint, gitOauth2ExternalClientId, gitOauth2ClientSecret,
								))),
					)
					t.Cleanup(cleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_authorization_endpoint", tfjson.ActionDelete, sdk.String(gitOauth2ExternalAuthorizationEndpoint), sdk.String(gitOauth2AuthorizationEndpoint)),
						planchecks.ExpectChange(ref, "oauth_token_endpoint", tfjson.ActionDelete, sdk.String(gitOauth2ExternalTokenEndpoint), sdk.String(gitOauth2TokenEndpoint)),
						planchecks.ExpectChange(ref, "oauth_client_id", tfjson.ActionDelete, sdk.String(gitOauth2ExternalClientId), sdk.String(gitOauth2ClientId)),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Update - ForceNew fields via config change (oauth_authorization_endpoint, oauth_token_endpoint) triggers destroy+recreate
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "oauth_authorization_endpoint", tfjson.ActionDelete, sdk.String(gitOauth2AuthorizationEndpoint), sdk.String(gitOauth2UpdatedAuthorizationEndpoint)),
						planchecks.ExpectChange(ref, "oauth_token_endpoint", tfjson.ActionDelete, sdk.String(gitOauth2TokenEndpoint), sdk.String(gitOauth2UpdatedTokenEndpoint)),
					},
				},
				Config: config.FromModels(t, allAttributesUpdated),
				Check:  assertThat(t, updatedCompleteAssertions...),
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
		[]string{gitAllowedPrefix},
		true,
		gitOauth2AuthorizationEndpoint,
		gitOauth2ClientId,
		gitOauth2ClientSecret,
		gitOauth2TokenEndpoint,
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

	clientSecret := random.AlphanumericN(20)
	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, clientSecret, gitOauth2TokenEndpoint).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
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
							[]sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitBlockedPrefix}}).
							WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest(
								*sdk.NewOAuth2GitUserAuthenticationRequest(
									gitOauth2AuthorizationEndpoint, gitOauth2TokenEndpoint, gitOauth2ClientId, clientSecret,
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

	invalidScope := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithOauthAllowedScopes([]string{"INVALID_SCOPE"})

	invalidAccessTokenValidity := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithOauthAccessTokenValidity(0)

	invalidRefreshTokenValidity := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithOauthRefreshTokenValidity(0)

	emptyAuthorizationEndpoint := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		"", gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint)

	emptyTokenEndpoint := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, "")

	emptyClientId := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, "", gitOauth2ClientSecret, gitOauth2TokenEndpoint)

	emptyUsername := model.ApiIntegrationGitRepositoryOauth2("t", id.Name(), []string{gitAllowedPrefix}, true,
		gitOauth2AuthorizationEndpoint, gitOauth2ClientId, gitOauth2ClientSecret, gitOauth2TokenEndpoint).
		WithOauthUsername("")

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
			{
				Config:      config.FromModels(t, emptyAuthorizationEndpoint),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_authorization_endpoint" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyTokenEndpoint),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_token_endpoint" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyClientId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_client_id" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyUsername),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_username" to not be an empty string`),
			},
		},
	})
}
