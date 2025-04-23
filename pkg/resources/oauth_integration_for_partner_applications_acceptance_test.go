//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForPartnerApplications_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"
	comment := random.Comment()

	basicModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl)
	completeModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl).
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithOauthIssueRefreshTokens(datasources.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "blocked_roles_list"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
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
				ResourceName: "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
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
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
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
				ResourceName: "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				Config: accconfig.FromModels(t, completeModel),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithBlockedRolesList(*sdk.NewBlockedRolesListRequest([]sdk.AccountObjectIdentifier{})).
							WithComment("").
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600),
					))
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest().
							WithEnabled(true).
							WithOauthUseSecondaryRoles(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift(completeModel.ResourceReference(), "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "comment", sdk.String(comment), sdk.String("")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_refresh_token_validity", sdk.String("86400"), sdk.String("3600")),

						planchecks.ExpectChange(completeModel.ResourceReference(), "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_refresh_token_validity", tfjson.ActionUpdate, sdk.String("3600"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
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
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
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

func TestAcc_OauthIntegrationForPartnerApplications_BasicTableauDesktop(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	basicModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")
	completeModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN", role.ID().Name()).
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue).
		WithOauthIssueRefreshTokens(datasources.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
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
				ResourceName: "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
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
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
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
				ResourceName: "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				Config: accconfig.FromModels(t, completeModel),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithBlockedRolesList(*sdk.NewBlockedRolesListRequest([]sdk.AccountObjectIdentifier{})).
							WithComment("").
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600),
					))
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest().
							WithEnabled(true).
							WithOauthUseSecondaryRoles(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift(completeModel.ResourceReference(), "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "comment", sdk.String(comment), sdk.String("")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "oauth_refresh_token_validity", sdk.String("86400"), sdk.String("3600")),

						planchecks.ExpectChange(completeModel.ResourceReference(), "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "oauth_refresh_token_validity", tfjson.ActionUpdate, sdk.String("3600"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
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
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "oauth_redirect_uri"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
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

func TestAcc_OauthIntegrationForPartnerApplications_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	completeModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauServer)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithEnabled(datasources.BooleanTrue).
		WithOauthIssueRefreshTokens(datasources.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "oauth_redirect_uri"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
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
				ResourceName: "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_Invalid(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	invalidUseSecondaryRolesModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithOauthUseSecondaryRoles("invalid")

	invalidOauthClientModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), "invalid").
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
				Config:      accconfig.FromModels(t, invalidOauthClientModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Error: invalid OauthSecurityIntegrationClientOption: INVALID`),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_WithQuotedName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForPartnerApplications("test", quotedId, string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, basicModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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

func TestAcc_OauthIntegrationForPartnerApplications_DetectExternalChangesForOauthRedirectUri(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	oauthRedirectUri := "https://example.com"
	otherOauthRedirectUri := "https://example2.com"
	configModel := model.OauthIntegrationForPartnerApplications(
		"test",
		id.Name(),
		string(sdk.OauthSecurityIntegrationClientLooker),
	).WithOauthRedirectUri(oauthRedirectUri)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.OauthIntegrationForPartnerApplicationsResource(t, configModel.ResourceReference()).
						HasOauthRedirectUriString(oauthRedirectUri),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().WithOauthRedirectUri(otherOauthRedirectUri),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.OauthIntegrationForPartnerApplicationsResource(t, configModel.ResourceReference()).
						HasOauthRedirectUriString(oauthRedirectUri),
				),
			},
		},
	})
}
