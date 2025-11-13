//go:build account_level_tests

// Tests for OAuth integration for partner applications are marked as account_level_tests,
// because only one integration of this type can exist in the account.

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForPartnerApplications_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basic := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri("https://example.com/callback")

	complete := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri("https://example.com/callback").
		WithEnabled(r.BooleanTrue).
		WithOauthIssueRefreshTokens(r.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.SecurityIntegration(t, id).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - LOOKER").
			HasCategory("SECURITY").
			HasEnabled(false).
			HasComment(""),

		resourceassert.OauthIntegrationForPartnerApplicationsResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasOauthClientString("LOOKER").
			HasOauthRedirectUriString("https://example.com/callback").
			HasEnabledString(r.BooleanDefault).
			HasOauthIssueRefreshTokensString(r.BooleanDefault).
			HasOauthRefreshTokenValidityString(r.IntDefaultString).
			HasCommentString("").
			HasBlockedRolesListEmpty().
			HasRelatedParametersNotEmpty().
			HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(r.BooleanTrue),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, basic.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - LOOKER").
			HasCategory("SECURITY").
			HasEnabled(false).
			HasComment(""),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "CONFIDENTIAL")),
		assert.Check(resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "NONE")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", r.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "7776000")),
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
			HasIntegrationType("OAUTH - LOOKER").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment),

		resourceassert.OauthIntegrationForPartnerApplicationsResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasOauthClientString("LOOKER").
			HasOauthRedirectUriString("https://example.com/callback").
			HasEnabledString(r.BooleanTrue).
			HasOauthIssueRefreshTokensString(r.BooleanFalse).
			HasOauthRefreshTokenValidityString("86400").
			HasOauthUseSecondaryRolesString("IMPLICIT").
			HasCommentString(comment).
			HasBlockedRolesListEmpty().
			HasRelatedParametersNotEmpty().
			HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(r.BooleanTrue),

		resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
			HasName(id.Name()).
			HasIntegrationType("OAUTH - LOOKER").
			HasCategory("SECURITY").
			HasEnabled(true).
			HasComment(comment),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "CONFIDENTIAL")),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanTrue)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", r.BooleanFalse)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
		assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
					"oauth_use_secondary_roles",
					"related_parameters",
				},
			},
			// Update - set optionals
			{
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
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
					"oauth_use_secondary_roles",
					"related_parameters",
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
					testClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithEnabled(true).
							WithComment(comment).
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600).
							WithOauthUseSecondaryRoles(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
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
			// Destroy - ensure security integration is destroyed before the last step
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
				Check: assertThat(t,
					invokeactionassert.SecurityIntegrationDoesNotExist(t, id),
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

func TestAcc_OauthIntegrationForPartnerApplications_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	complete := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauServer)).
		WithEnabled(r.BooleanTrue).
		WithOauthIssueRefreshTokens(r.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional ForceNew parameters)
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					objectassert.SecurityIntegration(t, id).
						HasName(id.Name()).
						HasIntegrationType("OAUTH - TABLEAU_SERVER").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment),

					resourceassert.OauthIntegrationForPartnerApplicationsResource(t, complete.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasOauthClientString("TABLEAU_SERVER").
						HasEnabledString(r.BooleanTrue).
						HasOauthIssueRefreshTokensString(r.BooleanFalse).
						HasOauthRefreshTokenValidityString("86400").
						HasOauthUseSecondaryRolesString("IMPLICIT").
						HasCommentString(comment),

					resourceassert.OauthIntegrationForPartnerApplicationsResource(t, complete.ResourceReference()).
						HasBlockedRolesListEmpty(),

					resourceshowoutputassert.SecurityIntegrationShowOutput(t, complete.ResourceReference()).
						HasName(id.Name()).
						HasIntegrationType("OAUTH - TABLEAU_SERVER").
						HasCategory("SECURITY").
						HasEnabled(true).
						HasComment(comment),

					// Describe output assertions
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_type.0.value", "PUBLIC")),
					assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_redirect_uri.0.value")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.enabled.0.value", r.BooleanTrue)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_issue_refresh_tokens.0.value", r.BooleanFalse)),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_validity.0.value", "86400")),
					assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
					assert.Check(resource.TestCheckNoResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_client_id.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_authorization_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_token_endpoint.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_authorization_endpoints.0.value")),
					assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_allowed_token_endpoints.0.value")),

					// Related parameters assertions
					resourceassert.OauthIntegrationForPartnerApplicationsResource(t, complete.ResourceReference()).
						HasRelatedParametersNotEmpty().
						HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(r.BooleanTrue),
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
					"oauth_issue_refresh_tokens",
					"oauth_refresh_token_validity",
					"oauth_use_secondary_roles",
					"related_parameters",
				},
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_BasicTableauDesktop(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	role, roleCleanup := testClient().Role.CreateRole(t)
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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
					testClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithBlockedRolesList(*sdk.NewBlockedRolesListRequest([]sdk.AccountObjectIdentifier{})).
							WithComment("").
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600),
					))
					testClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(
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
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	completeModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauServer)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithEnabled(datasources.BooleanTrue).
		WithOauthIssueRefreshTokens(datasources.BooleanFalse).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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
	id := testClient().Ids.RandomAccountObjectIdentifier()

	invalidUseSecondaryRolesModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN").
		WithOauthUseSecondaryRoles("invalid")

	invalidOauthClientModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), "invalid").
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

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
				Config:      accconfig.FromModels(t, invalidOauthClientModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Error: invalid OauthSecurityIntegrationClientOption: INVALID`),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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

func TestAcc_OauthIntegrationForPartnerApplications_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	validUrl := "https://example.com"

	basicModel := model.OauthIntegrationForPartnerApplications("test", quotedId, string(sdk.OauthSecurityIntegrationClientLooker)).
		WithOauthRedirectUri(validUrl).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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

func TestAcc_OauthIntegrationForPartnerApplications_DetectExternalChangesForOauthRedirectUri(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	oauthRedirectUri := "https://example.com"
	otherOauthRedirectUri := "https://example2.com"
	configModel := model.OauthIntegrationForPartnerApplications(
		"test",
		id.Name(),
		string(sdk.OauthSecurityIntegrationClientLooker),
	).WithOauthRedirectUri(oauthRedirectUri)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
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
					testClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().WithOauthRedirectUri(otherOauthRedirectUri),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModel.ResourceReference(), plancheck.ResourceActionNoop),
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

func TestAcc_OauthIntegrationForPartnerApplications_WithPrivilegedRolesBlockedList(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	// Use an identifier with this prefix to have this role in the end.
	roleId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("Z")
	role, roleCleanup := testClient().Role.CreateRoleWithIdentifier(t, roleId)
	t.Cleanup(roleCleanup)

	allRoles := []string{snowflakeroles.Accountadmin.Name(), snowflakeroles.SecurityAdmin.Name(), role.ID().Name()}
	onlyPrivilegedRoles := []string{snowflakeroles.Accountadmin.Name(), snowflakeroles.SecurityAdmin.Name()}
	customRoles := []string{role.ID().Name()}

	paramCleanup := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList, "true")
	t.Cleanup(paramCleanup)

	modelWithoutBlockedRole := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop))
	modelWithBlockedRole := model.OauthIntegrationForPartnerApplications("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTableauDesktop)).
		WithBlockedRolesList(role.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.*", role.ID().Name()),
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(allRoles, ",")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "blocked_roles_list.#", "0"),
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(onlyPrivilegedRoles, ",")),
				),
			},
			{
				PreConfig: func() {
					// Do not revert, because the revert is setup above.
					testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList, "false")
				},
				Config: accconfig.FromModels(t, modelWithBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(modelWithBlockedRole.ResourceReference(), "blocked_roles_list.*", role.ID().Name()),
					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", strings.Join(customRoles, ",")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutBlockedRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "blocked_roles_list.#", "0"),
					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "name", id.Name()),

					resource.TestCheckResourceAttr(modelWithoutBlockedRole.ResourceReference(), "describe_output.0.blocked_roles_list.0.value", ""),
				),
			},
		},
	})
}
