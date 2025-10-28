//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	resourcenames "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForCustomClients_BasicUseCase(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	validUrl := "https://example.com/callback"
	id := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()

	basic := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl)

	complete := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithOauthAllowNonTlsRedirectUri(resources.BooleanTrue).
		WithOauthClientRsaPublicKey(key).
		WithOauthClientRsaPublicKey2(key).
		WithOauthEnforcePkce(resources.BooleanTrue).
		WithEnabled(resources.BooleanTrue).
		WithOauthIssueRefreshTokens(resources.BooleanTrue).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID()).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - CUSTOM").
			HasCategory("SECURITY").
			HasEnabled(false).
			HasComment(""),

		resourceassert.OauthIntegrationForCustomClientsResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasOauthClientTypeString("CONFIDENTIAL").
			HasOauthRedirectUriString(validUrl).
			HasEnabledString(resources.BooleanDefault).
			HasOauthAllowNonTlsRedirectUriString(resources.BooleanDefault).
			HasOauthEnforcePkceString(resources.BooleanDefault).
			HasOauthIssueRefreshTokensString(resources.BooleanDefault).
			HasOauthAllowNonTlsRedirectUriString(resources.BooleanDefault).
			HasPreAuthorizedRolesListLen(0).
			HasRelatedParametersNotEmpty().
			HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(resources.BooleanTrue).
			HasBlockedRolesListEmpty(),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, basic.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - CUSTOM").
			HasCategory("SECURITY").
			HasEnabled(false).
			HasComment(""),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "CONFIDENTIAL")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.enabled.0.value", resources.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", resources.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.network_policy.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment.0.value", "")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - CUSTOM").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment),

		resourceassert.OauthIntegrationForCustomClientsResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasOauthClientTypeString("CONFIDENTIAL").
			HasOauthRedirectUriString(validUrl).
			HasEnabledString(resources.BooleanTrue).
			HasOauthAllowNonTlsRedirectUriString(resources.BooleanTrue).
			HasOauthEnforcePkceString(resources.BooleanTrue).
			HasOauthIssueRefreshTokensString(resources.BooleanTrue).
			HasOauthRefreshTokenValidityString("86400").
			HasOauthUseSecondaryRolesString("IMPLICIT").
			HasPreAuthorizedRolesListLen(1).
			HasPreAuthorizedRolesListElem(0, preAuthorizedRole.ID().Name()).
			HasNetworkPolicyString(networkPolicy.ID().Name()).
			HasOauthClientRsaPublicKeyString(key).
			HasOauthClientRsaPublicKey2String(key).
			HasCommentString(comment).
			HasRelatedParametersNotEmpty().
			HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(resources.BooleanTrue).
			HasBlockedRolesListEmpty(),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - CUSTOM").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "CONFIDENTIAL")),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", resources.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.ID().Name())),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
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
					"blocked_roles_list",
					"enabled",
					"oauth_allow_non_tls_redirect_uri",
					"oauth_enforce_pkce",
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
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
					"blocked_roles_list",
					"enabled",
					"oauth_allow_non_tls_redirect_uri",
					"oauth_enforce_pkce",
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
					"oauth_use_secondary_roles",
					"oauth_client_rsa_public_key",
					"oauth_client_rsa_public_key_2",
				},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().SecurityIntegration.UpdateOauthForClients(t, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForCustomClientsIntegrationSetRequest().
							WithEnabled(true).
							WithNetworkPolicy(networkPolicy.ID()).
							WithOauthEnforcePkce(true).
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(86400).
							WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit).
							WithPreAuthorizedRolesList(*sdk.NewPreAuthorizedRolesListRequest().WithPreAuthorizedRolesList([]sdk.AccountObjectIdentifier{preAuthorizedRole.ID()})).
							WithComment(comment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure external oauth integration is deleted
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.SecurityIntegrationDoesNotExist(t, id),
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

func TestAcc_OauthIntegrationForCustomClients_CompleteUseCase(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	validUrl := "https://example.com/callback"
	id := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()

	complete := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypePublic), validUrl).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithOauthAllowNonTlsRedirectUri(resources.BooleanTrue).
		WithOauthClientRsaPublicKey(key).
		WithOauthClientRsaPublicKey2(key).
		WithOauthEnforcePkce(resources.BooleanTrue).
		WithEnabled(resources.BooleanTrue).
		WithOauthIssueRefreshTokens(resources.BooleanTrue).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional ForceNew parameters)
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					objectassert.SecurityIntegration(t, id).
						HasName(id.Name()).
						HasIntegrationType("OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment),

					resourceassert.OauthIntegrationForCustomClientsResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasOauthClientTypeString("PUBLIC").
						HasOauthRedirectUriString(validUrl).
						HasEnabledString(resources.BooleanTrue).
						HasOauthAllowNonTlsRedirectUriString(resources.BooleanTrue).
						HasOauthEnforcePkceString(resources.BooleanTrue).
						HasOauthIssueRefreshTokensString(resources.BooleanTrue).
						HasOauthRefreshTokenValidityString("86400").
						HasOauthUseSecondaryRolesString("IMPLICIT").
						HasPreAuthorizedRolesListLen(1).
						HasPreAuthorizedRolesListElem(0, preAuthorizedRole.ID().Name()).
						HasNetworkPolicyString(networkPolicy.ID().Name()).
						HasOauthClientRsaPublicKeyString(key).
						HasOauthClientRsaPublicKey2String(key).
						HasCommentString(comment).
						HasRelatedParametersNotEmpty().
						HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(resources.BooleanTrue).
						HasBlockedRolesListEmpty(),

					resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
						HasName(id.Name()).
						HasIntegrationType("OAUTH - CUSTOM").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "PUBLIC")),
					assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", resources.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", resources.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", resources.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", resources.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
					assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value")),
				),
			},
			// Import - with all optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"blocked_roles_list",
					"enabled",
					"oauth_allow_non_tls_redirect_uri",
					"oauth_enforce_pkce",
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
					"oauth_use_secondary_roles",
					"related_parameters",
					"oauth_client_rsa_public_key",
					"oauth_client_rsa_public_key_2",
				},
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_DefaultValues(t *testing.T) {
	validUrl := "https://example.com"
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl)
	defaultValuesModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithComment("").
		WithEnabled(resources.BooleanFalse).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithNetworkPolicy("").
		WithOauthAllowNonTlsRedirectUri(resources.BooleanFalse).
		WithOauthClientRsaPublicKeyEmpty().
		WithOauthClientRsaPublicKey2Empty().
		WithOauthEnforcePkce(resources.BooleanFalse).
		WithOauthIssueRefreshTokens(resources.BooleanFalse).
		WithOauthRefreshTokenValidity(7776000).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)).
		WithPreAuthorizedRolesEmpty()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, defaultValuesModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_enforce_pkce", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_refresh_token_validity", "7776000"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "network_policy", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_client_rsa_public_key", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "oauth_client_rsa_public_key_2", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(defaultValuesModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.comment.0.value", ""),
					resource.TestCheckNoResourceAttr(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(defaultValuesModel.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_enforce_pkce", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "network_policy", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client_rsa_public_key", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client_rsa_public_key_2", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.comment.0.value", ""),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_Invalid(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"

	invalidUseSecondaryRolesModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithOauthUseSecondaryRoles("invalid")
	invalidClientTypesModel := model.OauthIntegrationForCustomClients("test", id.Name(), "invalid", validUrl)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, invalidUseSecondaryRolesModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Error: invalid OauthSecurityIntegrationUseSecondaryRolesOption: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, invalidClientTypesModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Error: invalid OauthSecurityIntegrationClientTypeOption: INVALID`),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForCustomClients("test", quotedId, string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}
