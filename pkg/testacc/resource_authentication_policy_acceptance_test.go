//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	acchelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicy(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	id2 := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	changedComment := random.Comment()
	samlIntegration, cleanupSamlIntegration := testClient().SecurityIntegration.CreateSaml2(t)
	t.Cleanup(cleanupSamlIntegration)
	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	basicModelWithDifferentName := model.AuthenticationPolicy("test", id2.DatabaseName(), id2.SchemaName(), id2.Name())
	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationAll).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodPassKey},
				{Method: sdk.MfaPolicyAllowedMethodDuo},
			}),
		).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithDefaultExpiryInDays(1).
			WithMaxExpiryInDays(30).
			WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationNotEnforced),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAll},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "111122223333"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/tenantid/v2.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://example.com"},
			}),
		).
		WithSecurityIntegrations("ALL")
	completeModelWithDifferentValues := model.AuthenticationPolicy("test", id2.DatabaseName(), id2.SchemaName(), id2.Name()).
		WithComment(changedComment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsSaml).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequiredPasswordOnly).
		WithClientTypes(sdk.ClientTypesSnowflakeCli).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationNone).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodTotp},
			}),
		).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithDefaultExpiryInDays(2).
			WithMaxExpiryInDays(40).
			WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationEnforcedNotRequired),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAzure},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "444455556666"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/tenantid/v3.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://example2.com"},
			}),
		).
		WithSecurityIntegrations(samlIntegration.ID().Name())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, basicModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasNoMfaEnrollment().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, basicModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, basicModel),
				ResourceName: basicModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAuthenticationPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("").
						HasAuthenticationMethods(sdk.AuthenticationMethodsAll).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequiredPasswordOnly)).
						HasClientTypes(sdk.ClientTypesAll).
						HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.ImportedAuthenticationPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, completeModel),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAuthenticationMethods(sdk.AuthenticationMethodsPassword).
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequired)).
					HasClientTypes(sdk.ClientTypesSnowflakeUi).
					HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD, SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
				),
			},
			// alter
			{
				Config: accconfig.FromModels(t, completeModelWithDifferentValues),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, completeModelWithDifferentValues.ResourceReference()).
					HasNameString(id2.Name()).
					HasDatabaseString(id2.DatabaseName()).
					HasSchemaString(id2.SchemaName()).
					HasCommentString(changedComment).
					HasAuthenticationMethods(sdk.AuthenticationMethodsSaml).
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequiredPasswordOnly)).
					HasClientTypes(sdk.ClientTypesSnowflakeCli).
					HasSecurityIntegrations(samlIntegration.ID().Name()),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModelWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id2.Name()).
						HasDatabaseName(id2.DatabaseName()).
						HasSchemaName(id2.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.authentication_methods", "[SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_types", "[SNOWFLAKE_CLI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.security_integrations", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=2, MAX_EXPIRY_IN_DAYS=40, NETWORK_POLICY_EVALUATION=ENFORCED_NOT_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[444455556666], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v3.0], ALLOWED_OIDC_ISSUERS=[https://example2.com]}")),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().AuthenticationPolicy.Alter(t, sdk.NewAlterAuthenticationPolicyRequest(id2).WithSet(
						*sdk.NewAuthenticationPolicySetRequest().
							WithComment(random.Comment()).
							WithAuthenticationMethods([]sdk.AuthenticationMethods{{Method: sdk.AuthenticationMethodsPassword}}).
							WithMfaEnrollment(sdk.MfaEnrollmentRequired).
							WithClientTypes([]sdk.ClientTypes{{ClientType: sdk.ClientTypesSnowflakeUi}}).
							WithSecurityIntegrations(*sdk.NewSecurityIntegrationsOptionRequest().WithAll(true)),
						// Changes to mfa_policy, pat_policy, workload_identity_policy are not detected.
					))
				},
				Config: accconfig.FromModels(t, completeModelWithDifferentValues),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, completeModelWithDifferentValues.ResourceReference()).
					HasNameString(id2.Name()).
					HasDatabaseString(id2.DatabaseName()).
					HasSchemaString(id2.SchemaName()).
					HasCommentString(changedComment).
					HasAuthenticationMethods(sdk.AuthenticationMethodsSaml).
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequiredPasswordOnly)).
					HasClientTypes(sdk.ClientTypesSnowflakeCli).
					HasSecurityIntegrations(samlIntegration.ID().Name()),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModelWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id2.Name()).
						HasDatabaseName(id2.DatabaseName()).
						HasSchemaName(id2.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.authentication_methods", "[SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_types", "[SNOWFLAKE_CLI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.security_integrations", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=2, MAX_EXPIRY_IN_DAYS=40, NETWORK_POLICY_EVALUATION=ENFORCED_NOT_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[444455556666], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v3.0], ALLOWED_OIDC_ISSUERS=[https://example2.com]}")),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, basicModelWithDifferentName),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, basicModelWithDifferentName.ResourceReference()).
					HasNameString(id2.Name()).
					HasDatabaseString(id2.DatabaseName()).
					HasSchemaString(id2.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasMfaEnrollmentEmpty().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, basicModelWithDifferentName.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id2.Name()).
						HasDatabaseName(id2.DatabaseName()).
						HasSchemaName(id2.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithSecurityIntegrations("ALL").
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationAll).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodPassKey},
				{Method: sdk.MfaPolicyAllowedMethodDuo},
			}),
		).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithDefaultExpiryInDays(1).
			WithMaxExpiryInDays(30).
			WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationNotEnforced),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAll},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "111122223333"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/tenantid/v2.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://example.com"},
			}),
		).
		WithComment(comment)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel),
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAuthenticationMethods(sdk.AuthenticationMethodsPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequired)).
						HasClientTypes(sdk.ClientTypesSnowflakeUi).
						HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModel.ResourceReference()).
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD, SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=false}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
				),
			},
			{
				Config:       accconfig.FromModels(t, completeModel),
				ResourceName: completeModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAuthenticationPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAuthenticationMethods(sdk.AuthenticationMethodsPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequired)).
						HasClientTypes(sdk.ClientTypesSnowflakeUi).
						HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.ImportedAuthenticationPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", comment)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.client_types", "[SNOWFLAKE_UI]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.security_integrations", "[ALL]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.mfa_enrollment", "REQUIRED")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.mfa_authentication_methods", "[PASSWORD, SAML]")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_enumSetCustomDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelWithEnumSets := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithClientTypes(sdk.ClientTypesSnowflakeUi)
	modelWithLowercaseEnumSets := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(acchelpers.EnumToLower(sdk.AuthenticationMethodsPassword)).
		WithClientTypes(acchelpers.EnumToLower(sdk.ClientTypesSnowflakeUi))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithEnumSets),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithEnumSets.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelWithLowercaseEnumSets),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_handlingLists(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	modelWithBasicLists := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationNone).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodTotp},
				{Method: sdk.MfaPolicyAllowedMethodDuo},
			}),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAzure},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "111111111111"}, {Value: "222222222222"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/one/v2.0"}, {Value: "https://login.microsoftonline.com/two/v2.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://one.com"}, {Value: "https://two.com"},
			}),
		)
	modelWithSomeListsEmpty := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationNone),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAzure},
			}),
		)
	modelWithSwappedLists := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationNone).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodDuo},
				{Method: sdk.MfaPolicyAllowedMethodTotp},
			}),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAzure},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "222222222222"}, {Value: "111111111111"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/two/v2.0"}, {Value: "https://login.microsoftonline.com/one/v2.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://two.com"}, {Value: "https://one.com"},
			}),
		)
	ref := modelWithBasicLists.ResourceReference()
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithBasicLists),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasNoMfaEnrollment().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[111111111111, 222222222222], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/one/v2.0, https://login.microsoftonline.com/two/v2.0], ALLOWED_OIDC_ISSUERS=[https://one.com, https://two.com]}")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithSomeListsEmpty),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasNoMfaEnrollment().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithBasicLists),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasNoMfaEnrollment().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[111111111111, 222222222222], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/one/v2.0, https://login.microsoftonline.com/two/v2.0], ALLOWED_OIDC_ISSUERS=[https://one.com, https://two.com]}")),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithSwappedLists.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelWithSwappedLists),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasNoMfaEnrollment().
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[111111111111, 222222222222], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/one/v2.0, https://login.microsoftonline.com/two/v2.0], ALLOWED_OIDC_ISSUERS=[https://one.com, https://two.com]}")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_migrateFromV2_9_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithSecurityIntegrations("ALL")
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.AuthenticationPolicyResource))
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.9.0"),
				Config:            accconfig.FromModels(t, providerModel, completeModel),
				// This happens because the mfa_authentication_methods is not set in the config,
				// and the value returned from Snowflake is non-empty.
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, completeModel),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_migrateFromV2_9_0_setNewFields(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.AuthenticationPolicyResource))

	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	modelWithNewFields := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationAll).
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodPassKey},
				{Method: sdk.MfaPolicyAllowedMethodDuo},
			}),
		).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithDefaultExpiryInDays(1).
			WithMaxExpiryInDays(30).
			WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationNotEnforced),
		).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAll},
			}).
			WithAllowedAwsAccounts([]sdk.StringListItemWrapper{
				{Value: "111122223333"},
			}).
			WithAllowedAzureIssuers([]sdk.StringListItemWrapper{
				{Value: "https://login.microsoftonline.com/tenantid/v2.0"},
			}).
			WithAllowedOidcIssuers([]sdk.StringListItemWrapper{
				{Value: "https://example.com"},
			}),
		)
	ref := modelWithNewFields.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.9.0"),
				Config:            accconfig.FromModels(t, providerModel, basicModel),
				// This happens because the mfa_authentication_methods is not set in the config,
				// and the value returned from Snowflake is non-empty.
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequiredPasswordOnly)).
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasNoKind().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.mfa_policy")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.pat_policy")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.workload_identity_policy")),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				// These assertions are commented out because in `Before` objects inside plan checks, the values are nil.
				// This happens because the planchecks assume that the object schema trees are "roughly" the same.
				// Once the planchecks support having different hierarchies, we can uncomment these assertions.
				// ConfigPlanChecks: resource.ConfigPlanChecks{
				// 	PreApply: []plancheck.PlanCheck{
				// 		plancheck.ExpectResourceAction(modelWithNewFields.ResourceReference(), plancheck.ResourceActionUpdate),
				// 		planchecks.PrintPlanDetails(ref, "mfa_policy.0.enforce_mfa_on_external_authentication"),
				// 		planchecks.ExpectNoChangeOnField(ref, "pat_policy"),
				// 		planchecks.ExpectNoChangeOnField(ref, "workload_identity_policy"),
				// 	},
				// },
				Config: accconfig.FromModels(t, modelWithNewFields),
				// The values are the same as in the previous step, but the mfa_policy is set.
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasMfaEnrollmentString("").
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD, SAML]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_migrateFromV2_9_0_setOldFields(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithSecurityIntegrations("ALL")
	ref := completeModel.ResourceReference()
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.AuthenticationPolicyResource))
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.9.0"),
				Config:            accconfig.FromModels(t, providerModel, basicModel),
				// This happens because the mfa_authentication_methods is not set in the config,
				// and the value returned from Snowflake is non-empty.
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, ref).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasAuthenticationMethodsEmpty().
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequiredPasswordOnly)).
					HasClientTypesEmpty().
					HasSecurityIntegrationsEmpty(),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, ref).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasNoKind().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.mfa_policy")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.pat_policy")),
					assert.Check(resource.TestCheckNoResourceAttr(ref, "describe_output.0.workload_identity_policy")),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, completeModel),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidAuthenticationMethods := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods("invalid")
	modelInvalidMfaEnrollment := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaEnrollmentEnum("invalid")
	modelInvalidClientTypes := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithClientTypes("invalid")
	modelInvalidEnforceMfaOnExternalAuthentication := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithEnforceMfaOnExternalAuthentication("invalid"),
		)
	modelInvalidAllowedMethods := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
			WithAllowedMethods([]sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: "invalid"},
			}),
		)
	modelInvalidDefaultExpiryInDays := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithDefaultExpiryInDays(0),
		)
	modelInvalidMaxExpiryInDays := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithMaxExpiryInDays(0),
		)
	modelInvalidNetworkPolicyEvaluation := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
			WithNetworkPolicyEvaluation("invalid"),
		)
	modelInvalidAllowedProviders := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
			WithAllowedProviders([]sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: "invalid"},
			}),
		)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAuthenticationMethods),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid authentication method: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMfaEnrollment),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid MFA enrollment option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidClientTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid client type: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidEnforceMfaOnExternalAuthentication),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid enforce MFA on external authentication option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAllowedMethods),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid MFA policy allowed methods option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidDefaultExpiryInDays),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected default_expiry_in_days to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMaxExpiryInDays),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_expiry_in_days to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidNetworkPolicyEvaluation),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid network policy evaluation option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAllowedProviders),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid allowed provider: INVALID`),
			},
		},
	})
}
