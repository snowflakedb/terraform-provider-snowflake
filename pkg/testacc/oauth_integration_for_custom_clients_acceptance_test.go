//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	resourcenames "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForCustomClients_Basic(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	blockedRole, blockedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(blockedRoleCleanup)

	validUrl := "https://example.com"
	id := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()

	basicModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl)
	completeModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithComment(comment).
		WithEnabled(resources.BooleanTrue).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN", blockedRole.ID().Name()).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithOauthAllowNonTlsRedirectUri(resources.BooleanTrue).
		WithOauthClientRsaPublicKey(key).
		WithOauthClientRsaPublicKey2(key).
		WithOauthEnforcePkce(resources.BooleanTrue).
		WithOauthIssueRefreshTokens(resources.BooleanTrue).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_enforce_pkce", resources.BooleanDefault),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", resources.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "network_policy", ""),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_client_rsa_public_key"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_client_rsa_public_key_2"),
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
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, basicModel),
				ResourceName: "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "0"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", ""),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, completeModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.Name),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, completeModel),
				ResourceName: "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				Config: accconfig.FromModels(t, completeModel),
				PreConfig: func() {
					testClient().SecurityIntegration.UpdateOauthForClients(t, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForCustomClientsIntegrationUnsetRequest().
							WithEnabled(true).
							WithNetworkPolicy(true).
							WithOauthUseSecondaryRoles(true).
							WithOauthClientRsaPublicKey(true).
							WithOauthClientRsaPublicKey2(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift(completeModel.ResourceReference(), "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "network_policy", sdk.String(networkPolicy.ID().Name()), sdk.String("")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_client_rsa_public_key", sdk.String(key), sdk.String(key)),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_client_rsa_public_key_2", sdk.String(key), sdk.String(key)),

						planchecks.ExpectChange(completeModel.ResourceReference(), "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "network_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String(networkPolicy.ID().Name())),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_client_rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(key)),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_client_rsa_public_key_2", tfjson.ActionUpdate, sdk.String(key), sdk.String(key)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
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

func TestAcc_OauthIntegrationForCustomClients_Complete(t *testing.T) {
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	blockedRole, blockedRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(blockedRoleCleanup)

	validUrl := "https://example.com"
	id := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()

	completeModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithComment(comment).
		WithEnabled(resources.BooleanTrue).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN", blockedRole.ID().Name()).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithOauthAllowNonTlsRedirectUri(resources.BooleanTrue).
		WithOauthClientRsaPublicKey(key).
		WithOauthClientRsaPublicKey2(key).
		WithOauthEnforcePkce(resources.BooleanTrue).
		WithOauthIssueRefreshTokens(resources.BooleanTrue).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			{
				Config:       accconfig.FromModels(t, completeModel),
				ResourceName: "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
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
		PreCheck:                 func() { TestAccPreCheck(t) },
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
		PreCheck:                 func() { TestAccPreCheck(t) },
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
		PreCheck: func() { TestAccPreCheck(t) },
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
		PreCheck: func() { TestAccPreCheck(t) },
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
