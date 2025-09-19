//go:build !account_level_tests

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

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Saml2Integration_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer, issuer2 := testClient().Ids.Alpha(), testClient().Ids.Alpha()
	cert, cert2 := random.GenerateX509(t), random.GenerateX509(t)
	validUrl, validUrl2 := "https://example.com", "https://example2.com"
	acsURL := testClient().Context.ACSURL(t)
	issuerURL := testClient().Context.IssuerURL(t)
	comment := random.Comment()

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)
	_, configVariables2 := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert2)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)
	// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
	completeModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, temporaryVariableName).
		WithEnabled(r.BooleanTrue).
		WithComment(comment).
		WithSaml2SsoUrl(validUrl2).
		WithSaml2EnableSpInitiated(r.BooleanTrue).
		WithSaml2ForceAuthn(r.BooleanTrue).
		WithSaml2PostLogoutRedirectUrl(validUrl).
		WithSaml2RequestedNameidFormat(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)).
		WithSaml2SignRequest(r.BooleanTrue).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("foo").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")
	recreatesModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("foo").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_issuer", issuer),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_x509_cert", cert),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "saml2_sp_initiated_login_page_label"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_enable_sp_initiated", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_sign_request", r.BooleanDefault),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "saml2_requested_nameid_format"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_post_logout_redirect_url", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "saml2_snowflake_issuer_url"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "saml2_snowflake_acs_url"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "allowed_user_domains"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "allowed_email_patterns"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_x509_cert.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_enable_sp_initiated.0.value", "false"),
					resource.TestCheckNoResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_sign_request.0.value", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_requested_nameid_format.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[]"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[]"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.comment.0.value", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(basicModel.ResourceReference(), "show_output.0.created_on"),
				),
			},
			// import - without optionals
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ResourceName:    basicModel.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sso_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_x509_cert", cert),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sp_initiated_login_page_label", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_enable_sp_initiated", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sign_request", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_requested_nameid_format", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_post_logout_redirect_url", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_user_domains"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_email_patterns"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config:          accconfig.FromModels(t, completeModel, temporaryVariableModel),
				ConfigVariables: configVariables2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_issuer.0.value", issuer2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sso_url.0.value", validUrl2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),
				),
			},
			// import - complete
			{
				Config:          accconfig.FromModels(t, completeModel, temporaryVariableModel),
				ConfigVariables: configVariables2,
				ResourceName:    completeModel.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_issuer", issuer2),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sso_url", validUrl2),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_x509_cert", cert2),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sp_initiated_login_page_label", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_enable_sp_initiated", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sign_request", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_post_logout_redirect_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_user_domains.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_user_domains.0", "example.com"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_email_patterns.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change values externally
			{
				Config:          accconfig.FromModels(t, completeModel, temporaryVariableModel),
				ConfigVariables: configVariables2,
				PreConfig: func() {
					testClient().SecurityIntegration.UpdateSaml2(t, sdk.NewAlterSaml2SecurityIntegrationRequest(id).
						WithUnset(*sdk.NewSaml2IntegrationUnsetRequest().
							WithSaml2RequestedNameidFormat(true).
							WithSaml2PostLogoutRedirectUrl(true).
							WithSaml2ForceAuthn(true).
							WithComment(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(completeModel.ResourceReference(), "saml2_requested_nameid_format", sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)), sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress))),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "saml2_post_logout_redirect_url", sdk.String(validUrl), sdk.String("")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "saml2_force_authn", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift(completeModel.ResourceReference(), "comment", sdk.String(comment), sdk.String("")),

						planchecks.ExpectChange(completeModel.ResourceReference(), "saml2_requested_nameid_format", tfjson.ActionUpdate, sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)), sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified))),
						planchecks.ExpectChange(completeModel.ResourceReference(), "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(""), sdk.String(validUrl)),
						planchecks.ExpectChange(completeModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(completeModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_issuer.0.value", issuer2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sso_url.0.value", validUrl2),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),
				),
			},
			// unset
			{
				Config:          accconfig.FromModels(t, recreatesModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "enabled", r.BooleanDefault),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_issuer", issuer),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_x509_cert", cert),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_enable_sp_initiated", r.BooleanDefault),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_sign_request", r.BooleanDefault),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_requested_nameid_format", ""),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_post_logout_redirect_url", ""),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckNoResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_x509_cert.0.value"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_enable_sp_initiated.0.value", "false"),
					resource.TestCheckNoResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_sign_request.0.value", "false"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "describe_output.0.comment.0.value", ""),

					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(recreatesModel.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(recreatesModel.ResourceReference(), "show_output.0.created_on"),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_forceAuthn(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)
	saml2ConfigForceAuthnTrueModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2ForceAuthn(r.BooleanTrue)
	saml2ConfigForceAuthnFalseModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2ForceAuthn(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			// set up with concrete saml2_force_authn
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(saml2ConfigForceAuthnTrueModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange(saml2ConfigForceAuthnTrueModel.ResourceReference(), "saml2_force_authn", tfjson.ActionCreate, nil, sdk.String("true")),
						planchecks.ExpectComputed(saml2ConfigForceAuthnTrueModel.ResourceReference(), "describe_output", true),
					},
				},
				Config:          accconfig.FromModels(t, saml2ConfigForceAuthnTrueModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnTrueModel.ResourceReference(), "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnTrueModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnTrueModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "true"),
				),
			},
			// import when saml2_force_authn in config
			{
				ConfigVariables: configVariables,
				ResourceName:    saml2ConfigForceAuthnTrueModel.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.saml2_force_authn.0.value", "true"),
				),
			},
			// change saml2_force_authn in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(saml2ConfigForceAuthnFalseModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange(saml2ConfigForceAuthnFalseModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectComputed(saml2ConfigForceAuthnFalseModel.ResourceReference(), "describe_output", true),
					},
				},
				Config:          accconfig.FromModels(t, saml2ConfigForceAuthnFalseModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// change back to non-default
			{
				Config:          accconfig.FromModels(t, saml2ConfigForceAuthnTrueModel, temporaryVariableModel),
				ConfigVariables: configVariables,
			},
			// remove non-default saml2_force_authn from config
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(basicModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange(basicModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectComputed(basicModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// add saml2_force_authn to config (false - which is a default in Snowflake) - no changes expected
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(saml2ConfigForceAuthnFalseModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config:          accconfig.FromModels(t, saml2ConfigForceAuthnFalseModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(saml2ConfigForceAuthnFalseModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// change back to non-default
			{
				Config:          accconfig.FromModels(t, saml2ConfigForceAuthnTrueModel, temporaryVariableModel),
				ConfigVariables: configVariables,
			},
			// remove saml2_force_authn from config but update externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					testClient().SecurityIntegration.UpdateSaml2ForceAuthn(t, id, false)
				},
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(basicModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange(basicModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectComputed(basicModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// change the saml2_force_authn externally
			{
				PreConfig: func() {
					// we change the type to the type different from default, expecting action
					testClient().SecurityIntegration.UpdateSaml2ForceAuthn(t, id, true)
				},
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(basicModel.ResourceReference(), "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange(basicModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectComputed(basicModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// import when no saml2_force_authn in config
			{
				ConfigVariables: configVariables,
				ResourceName:    basicModel.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_complete(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"
	acsURL := testClient().Context.ACSURL(t)
	issuerURL := testClient().Context.IssuerURL(t)
	comment := random.Comment()

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
	completeModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithEnabled(r.BooleanTrue).
		WithComment(comment).
		WithSaml2EnableSpInitiated(r.BooleanTrue).
		WithSaml2ForceAuthn(r.BooleanTrue).
		WithSaml2PostLogoutRedirectUrl(validUrl).
		WithSaml2RequestedNameidFormat(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)).
		WithSaml2SignRequest(r.BooleanTrue).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SsoUrl(validUrl).
		WithSaml2SpInitiatedLoginPageLabel("foo").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, completeModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_issuer", issuer),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_x509_cert", cert),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckNoResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", validUrl),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment.0.value", comment),

					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(completeModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(completeModel.ResourceReference(), "show_output.0.created_on"),
				),
			},
			{
				Config:          accconfig.FromModels(t, completeModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ResourceName:    completeModel.ResourceReference(),
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sso_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_x509_cert", cert),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sp_initiated_login_page_label", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_enable_sp_initiated", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_sign_request", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_post_logout_redirect_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(resourcehelpers.EncodeResourceIdentifier(id), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_user_domains.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_user_domains.0", "example.com"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_email_patterns.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_InvalidNameIdFormat(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2RequestedNameidFormat("invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: invalid Saml2SecurityIntegrationSaml2RequestedNameidFormatOption: invalid"),
			},
		},
	})
}

func TestAcc_Saml2Integration_InvalidProvider(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, "invalid", validUrl, temporaryVariableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: invalid Saml2SecurityIntegrationSaml2ProviderOption: INVALID"),
			},
		},
	})
}

func TestAcc_Saml2Integration_ForceNewIfEmpty(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"
	acsURL := testClient().Context.ACSURL(t)
	issuerURL := testClient().Context.IssuerURL(t)

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	baseModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("label").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")
	withoutLoginPageLabelModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")
	withoutSnowflakeIssuerUrlModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl("").
		WithSaml2SpInitiatedLoginPageLabel("label").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")
	withoutAcsUrlModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl("").
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("label").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomains("example.com")
	withoutAllowedEmailPatternsModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("label").
		WithAllowedEmailPatternsValue(accconfig.EmptyListVariable()).
		WithAllowedUserDomains("example.com")
	withoutAllowedUserDomainsModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithSaml2SnowflakeAcsUrl(acsURL).
		WithSaml2SnowflakeIssuerUrl(issuerURL).
		WithSaml2SpInitiatedLoginPageLabel("label").
		WithAllowedEmailPatterns("^(.+dev)@example.com$").
		WithAllowedUserDomainsValue(accconfig.EmptyListVariable())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, baseModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),

					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr(baseModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
				),
			},
			{
				Config:          accconfig.FromModels(t, withoutLoginPageLabelModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withoutLoginPageLabelModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "saml2_sp_initiated_login_page_label", ""),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", ""),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(withoutLoginPageLabelModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				Config:          accconfig.FromModels(t, withoutSnowflakeIssuerUrlModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withoutSnowflakeIssuerUrlModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "saml2_snowflake_issuer_url", ""),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", strings.ToLower(issuerURL)),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(withoutSnowflakeIssuerUrlModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				Config:          accconfig.FromModels(t, withoutAcsUrlModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withoutAcsUrlModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "saml2_snowflake_acs_url", ""),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", strings.ToLower(acsURL)),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(withoutAcsUrlModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				Config:          accconfig.FromModels(t, withoutAllowedEmailPatternsModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withoutAllowedEmailPatternsModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "allowed_email_patterns.#", "0"),

					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr(withoutAllowedEmailPatternsModel.ResourceReference(), "describe_output.0.allowed_email_patterns.0.value", "[]"),
				),
			},
			{
				Config:          accconfig.FromModels(t, withoutAllowedUserDomainsModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(withoutAllowedUserDomainsModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "allowed_user_domains.#", "0"),

					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr(withoutAllowedUserDomainsModel.ResourceReference(), "describe_output.0.allowed_user_domains.0.value", "[]"),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_DefaultValues(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	issuer := testClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)
	withZeroValuesModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithEnabled(r.BooleanFalse).
		WithSaml2ForceAuthn(r.BooleanFalse).
		WithSaml2PostLogoutRedirectUrl("")
	withNonZeroValuesModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName).
		WithEnabled(r.BooleanTrue).
		WithSaml2ForceAuthn(r.BooleanTrue).
		WithSaml2PostLogoutRedirectUrl(validUrl)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// create with valid "zero" values
			{
				Config:          accconfig.FromModels(t, withZeroValuesModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "enabled", tfjson.ActionCreate, nil, sdk.String("false")),
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "saml2_force_authn", tfjson.ActionCreate, nil, sdk.String("false")),
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "saml2_post_logout_redirect_url", tfjson.ActionCreate, nil, sdk.String("")),
						planchecks.ExpectComputed(withZeroValuesModel.ResourceReference(), "show_output", true),
						planchecks.ExpectComputed(withZeroValuesModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				Config:          accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(basicModel.ResourceReference(), "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectChange(basicModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectChange(basicModel.ResourceReference(), "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectComputed(basicModel.ResourceReference(), "show_output", true),
						planchecks.ExpectComputed(basicModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "enabled", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_force_authn", r.BooleanDefault),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// set to "non-zero" values
			{
				Config:          accconfig.FromModels(t, withNonZeroValuesModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "enabled", "true"),
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "saml2_post_logout_redirect_url", validUrl),

					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "show_output.0.enabled", "true"),

					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", validUrl),
					resource.TestCheckResourceAttr(withNonZeroValuesModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "true"),
				),
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				Config:          accconfig.FromModels(t, withZeroValuesModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "enabled", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "saml2_force_authn", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
						planchecks.ExpectChange(withZeroValuesModel.ResourceReference(), "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(validUrl), sdk.String("")),
						planchecks.ExpectComputed(withZeroValuesModel.ResourceReference(), "show_output", true),
						planchecks.ExpectComputed(withZeroValuesModel.ResourceReference(), "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr(withZeroValuesModel.ResourceReference(), "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// import zero values
			{
				Config:          accconfig.FromModels(t, withZeroValuesModel, temporaryVariableModel),
				ConfigVariables: configVariables,
				ImportState:     true,
				ResourceName:    withZeroValuesModel.ResourceReference(),
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "saml2_post_logout_redirect_url", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "show_output.0.enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.saml2_force_authn.0.value", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	cert := random.GenerateX509(t)
	issuer := testClient().Ids.Alpha()
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables:   configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables:          configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	cert := random.GenerateX509(t)
	issuer := testClient().Ids.Alpha()
	validUrl := "https://example.com"

	temporaryVariableName := "saml2_x509_cert"
	temporaryVariableModel, configVariables := accconfig.SecretStringVariableModelWithConfigVariables(temporaryVariableName, cert)

	basicModel := model.Saml2SecurityIntegrationVar("test", quotedId, issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, temporaryVariableName)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables:    configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel, temporaryVariableModel),
				ConfigVariables:          configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "id", id.Name()),
				),
			},
		},
	})
}
