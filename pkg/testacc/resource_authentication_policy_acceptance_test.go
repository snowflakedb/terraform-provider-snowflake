//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	acchelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	id2 := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	changedComment := random.Comment()

	samlIntegration, cleanupSamlIntegration := testClient().SecurityIntegration.CreateSaml2(t)
	t.Cleanup(cleanupSamlIntegration)

	completeModelClientPolicy := []sdk.AuthenticationPolicyClientPolicyRequest{
		{
			ClientType: sdk.ClientPolicyDriverTypeGoDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "1.14.1"},
		},
		{
			ClientType: sdk.ClientPolicyDriverTypeJdbcDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "3.25.0"},
		},
	}
	differentModelClientPolicy := []sdk.AuthenticationPolicyClientPolicyRequest{
		{
			ClientType: sdk.ClientPolicyDriverTypeGoDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "1.14.1"},
		},
		{
			ClientType: sdk.ClientPolicyDriverTypeJdbcDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "4.0.0"},
		},
	}

	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())

	basicModelWithDifferentName := model.AuthenticationPolicy("test", id2.DatabaseName(), id2.SchemaName(), id2.Name())

	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsOptionPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentOptionRequired).
		WithClientTypes(sdk.ClientTypesOptionSnowflakeUi, sdk.ClientTypesOptionDrivers).
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionAll).
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: sdk.MfaPolicyAllowedMethodsOptionPasskey},
					{Value: sdk.MfaPolicyAllowedMethodsOptionDuo},
				}),
		).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithDefaultExpiryInDays(1).
				WithMaxExpiryInDays(30).
				WithRequireRoleRestrictionForServiceUsers(false).
				WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationOptionNotEnforced),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAll},
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
		WithSecurityIntegrations("ALL").
		WithClientPolicy(completeModelClientPolicy...)

	completeModelWithDifferentValues := model.AuthenticationPolicy("test", id2.DatabaseName(), id2.SchemaName(), id2.Name()).
		WithComment(changedComment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsOptionSaml).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentOptionRequiredPasswordOnly).
		WithClientTypes(sdk.ClientTypesOptionSnowflakeCli, sdk.ClientTypesOptionDrivers).
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionNone).
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: sdk.MfaPolicyAllowedMethodsOptionTotp},
				}),
		).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithDefaultExpiryInDays(2).
				WithMaxExpiryInDays(40).
				WithRequireRoleRestrictionForServiceUsers(true).
				WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationOptionEnforcedNotRequired),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAzure},
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
		WithSecurityIntegrations(samlIntegration.ID().Name()).
		WithClientPolicy(differentModelClientPolicy...)
	completeWithDifferentValuesRef := completeModelWithDifferentValues.ResourceReference()

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModel),
				Check: assertThat(
					t,
					resourceassert.AuthenticationPolicyResource(t, basicModel.ResourceReference()).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.client_policy", "{}")),
				),
			},
			// import - without optionals, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, basicModel),
				ResourceName:             basicModel.ResourceReference(),
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// set optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, completeModel),
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasAuthenticationMethodsEnum(sdk.AuthenticationMethodsOptionPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentOptionRequired)).
						HasClientTypesEnum(sdk.ClientTypesOptionSnowflakeUi, sdk.ClientTypesOptionDrivers).
						HasSecurityIntegrations("ALL").
						HasClientPolicyEntries(completeModelClientPolicy...),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[DRIVERS, SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", string(sdk.MfaEnrollmentReadOptionRequired))),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=false, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_policy", "{JDBC_DRIVER={MINIMUM_VERSION=3.25.0}, GO_DRIVER={MINIMUM_VERSION=1.14.1}}")),
				),
			},
			// alter
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, completeModelWithDifferentValues),
				Check: assertThat(
					t,
					resourceassert.AuthenticationPolicyResource(t, completeModelWithDifferentValues.ResourceReference()).
						HasNameString(id2.Name()).
						HasDatabaseString(id2.DatabaseName()).
						HasSchemaString(id2.SchemaName()).
						HasCommentString(changedComment).
						HasAuthenticationMethodsEnum(sdk.AuthenticationMethodsOptionSaml).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentOptionRequiredPasswordOnly)).
						HasClientTypesEnum(sdk.ClientTypesOptionSnowflakeCli, sdk.ClientTypesOptionDrivers).
						HasSecurityIntegrations(samlIntegration.ID().Name()).
						HasClientPolicyEntries(differentModelClientPolicy...),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModelWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id2.Name()).
						HasDatabaseName(id2.DatabaseName()).
						HasSchemaName(id2.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.authentication_methods", "[SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_types", "[DRIVERS, SNOWFLAKE_CLI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.security_integrations", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_enrollment", string(sdk.MfaEnrollmentReadOptionRequiredPasswordOnly))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=2, MAX_EXPIRY_IN_DAYS=40, NETWORK_POLICY_EVALUATION=ENFORCED_NOT_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[444455556666], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v3.0], ALLOWED_OIDC_ISSUERS=[https://example2.com]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_policy", "{JDBC_DRIVER={MINIMUM_VERSION=4.0.0}, GO_DRIVER={MINIMUM_VERSION=1.14.1}}")),
				),
			},
			// change externally
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().AuthenticationPolicy.Alter(t, sdk.NewAlterAuthenticationPolicyRequest(id2).WithSet(
						*sdk.NewAuthenticationPolicySetRequest().
							WithComment(random.Comment()).
							WithAuthenticationMethods([]sdk.AuthenticationMethodsOptionListItem{{Value: sdk.AuthenticationMethodsOptionPassword}}).
							WithMfaEnrollment(sdk.MfaEnrollmentOptionRequired).
							WithClientTypes([]sdk.ClientTypesOptionListItem{{Value: sdk.ClientTypesOptionSnowflakeUi}, {Value: sdk.ClientTypesOptionDrivers}}).
							WithSecurityIntegrations(*sdk.NewAuthenticationPolicySecurityIntegrationsOptionRequest().WithAll(true)).
							WithMfaPolicy(*sdk.NewAuthenticationPolicyMfaPolicyRequest().
								WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionAll).
								WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{{Value: sdk.MfaPolicyAllowedMethodsOptionPasskey}}),
							).
							WithPatPolicy(*sdk.NewAuthenticationPolicyPatPolicyRequest().
								WithDefaultExpiryInDays(5).
								WithMaxExpiryInDays(50).
								WithRequireRoleRestrictionForServiceUsers(false).
								WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationOptionNotEnforced),
							).
							WithWorkloadIdentityPolicy(*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
								WithAllowedProviders([]sdk.AllowedProviderOptionListItem{{Value: sdk.AllowedProviderOptionAws}}).
								WithAllowedAwsAccounts([]sdk.StringListItemWrapper{{Value: "111122223333"}}).
								WithAllowedAzureIssuers([]sdk.StringListItemWrapper{{Value: "https://login.microsoftonline.com/tenantid/v2.0"}}).
								WithAllowedOidcIssuers([]sdk.StringListItemWrapper{{Value: "https://example.com"}}),
							).
							WithClientPolicy([]sdk.AuthenticationPolicyClientPolicyRequest{{
								ClientType: sdk.ClientPolicyDriverTypeGoDriver,
								Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "1.15.0"},
							}}),
					))
				},
				Config: accconfig.FromModels(t, completeModelWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeWithDifferentValuesRef, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "mfa_policy.0.enforce_mfa_on_external_authentication", tfjson.ActionUpdate, new(string(sdk.EnforceMfaOnExternalAuthenticationOptionAll)), new(string(sdk.EnforceMfaOnExternalAuthenticationOptionNone))),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "mfa_policy.0.allowed_methods.0", tfjson.ActionUpdate, new(string(sdk.MfaPolicyAllowedMethodsOptionPasskey)), new(string(sdk.MfaPolicyAllowedMethodsOptionTotp))),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "pat_policy.0.default_expiry_in_days", tfjson.ActionUpdate, new("5"), new("2")),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "pat_policy.0.max_expiry_in_days", tfjson.ActionUpdate, new("50"), new("40")),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "pat_policy.0.require_role_restriction_for_service_users", tfjson.ActionUpdate, new(r.BooleanFalse), new(r.BooleanTrue)),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "pat_policy.0.network_policy_evaluation", tfjson.ActionUpdate, new(string(sdk.NetworkPolicyEvaluationOptionNotEnforced)), new(string(sdk.NetworkPolicyEvaluationOptionEnforcedNotRequired))),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "workload_identity_policy.0.allowed_providers.0", tfjson.ActionUpdate, new(string(sdk.AllowedProviderOptionAws)), new(string(sdk.AllowedProviderOptionAzure))),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "workload_identity_policy.0.allowed_aws_accounts.0", tfjson.ActionUpdate, new("111122223333"), new("444455556666")),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "workload_identity_policy.0.allowed_azure_issuers.0", tfjson.ActionUpdate, new("https://login.microsoftonline.com/tenantid/v2.0"), new("https://login.microsoftonline.com/tenantid/v3.0")),
						planchecks.ExpectChange(completeWithDifferentValuesRef, "workload_identity_policy.0.allowed_oidc_issuers.0", tfjson.ActionUpdate, new("https://example.com"), new("https://example2.com")),
					},
				},
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, completeModelWithDifferentValues.ResourceReference()).
						HasNameString(id2.Name()).
						HasDatabaseString(id2.DatabaseName()).
						HasSchemaString(id2.SchemaName()).
						HasCommentString(changedComment).
						HasAuthenticationMethodsEnum(sdk.AuthenticationMethodsOptionSaml).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentOptionRequiredPasswordOnly)).
						HasClientTypesEnum(sdk.ClientTypesOptionSnowflakeCli, sdk.ClientTypesOptionDrivers).
						HasSecurityIntegrations(samlIntegration.ID().Name()).
						HasClientPolicyEntries(differentModelClientPolicy...),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModelWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id2.Name()).
						HasDatabaseName(id2.DatabaseName()).
						HasSchemaName(id2.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.authentication_methods", "[SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_types", "[DRIVERS, SNOWFLAKE_CLI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.security_integrations", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_enrollment", string(sdk.MfaEnrollmentReadOptionRequiredPasswordOnly))),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=2, MAX_EXPIRY_IN_DAYS=40, NETWORK_POLICY_EVALUATION=ENFORCED_NOT_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[444455556666], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v3.0], ALLOWED_OIDC_ISSUERS=[https://example2.com]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModelWithDifferentValues.ResourceReference(), "describe_output.0.client_policy", "{JDBC_DRIVER={MINIMUM_VERSION=4.0.0}, GO_DRIVER={MINIMUM_VERSION=1.14.1}}")),
				),
			},
			// unset
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basicModelWithDifferentName),
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, basicModelWithDifferentName.ResourceReference()).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
					assert.Check(resource.TestCheckResourceAttr(basicModelWithDifferentName.ResourceReference(), "describe_output.0.client_policy", "{}")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	completeModelClientPolicy := []sdk.AuthenticationPolicyClientPolicyRequest{
		{
			ClientType: sdk.ClientPolicyDriverTypeGoDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "1.14.1"},
		},
		{
			ClientType: sdk.ClientPolicyDriverTypeJdbcDriver,
			Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "3.25.0"},
		},
	}

	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(sdk.AuthenticationMethodsOptionPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentOptionRequired).
		WithClientTypes(sdk.ClientTypesOptionSnowflakeUi, sdk.ClientTypesOptionDrivers).
		WithSecurityIntegrations("ALL").
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionAll).
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: sdk.MfaPolicyAllowedMethodsOptionPasskey},
					{Value: sdk.MfaPolicyAllowedMethodsOptionDuo},
				}),
		).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithDefaultExpiryInDays(1).
				WithMaxExpiryInDays(30).
				WithRequireRoleRestrictionForServiceUsers(false).
				WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationOptionNotEnforced),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAll},
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
		WithComment(comment).
		WithClientPolicy(completeModelClientPolicy...)

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			// create
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAuthenticationMethodsEnum(sdk.AuthenticationMethodsOptionPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentOptionRequired)).
						HasClientTypesEnum(sdk.ClientTypesOptionSnowflakeUi, sdk.ClientTypesOptionDrivers).
						HasSecurityIntegrations("ALL").
						HasClientPolicyEntries(completeModelClientPolicy...),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, completeModel.ResourceReference()).
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[DRIVERS, SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", string(sdk.MfaEnrollmentReadOptionRequired))),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=false, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_policy", "{JDBC_DRIVER={MINIMUM_VERSION=3.25.0}, GO_DRIVER={MINIMUM_VERSION=1.14.1}}")),
				),
			},
			// import - complete
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, completeModel),
				ResourceName:             completeModel.ResourceReference(),
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedAuthenticationPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAuthenticationMethodsEnum(sdk.AuthenticationMethodsOptionPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentReadOptionRequired)).
						HasClientTypesEnum(sdk.ClientTypesOptionSnowflakeUi, sdk.ClientTypesOptionDrivers).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", comment)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.client_types", "[DRIVERS, SNOWFLAKE_UI]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.security_integrations", "[ALL]")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.mfa_enrollment", string(sdk.MfaEnrollmentReadOptionRequired))),
				),
			},
			// import - complete, permissive mode
			// TODO(SNOW-4392614): add suppressors for the authentication_methods, client_types, and security_integrations sets
			// and the mfa_policy, pat_policy, and workload_identity_policy blocks, then expect an empty plan.
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, completeModel),
				ResourceName:             completeModel.ResourceReference(),
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
				ExpectNonEmptyPlan:       true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	defaultMfaEnrollment := testClient().SnowflakeDefaults.DefaultMfaEnrollment(t)

	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	_, cleanup := testClient().AuthenticationPolicy.CreateWithOptions(t, id, sdk.NewCreateAuthenticationPolicyRequest(id))
	t.Cleanup(cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			// Import - without optionals
			{
				Config:             accconfig.FromModels(t, basicModel),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
			},
			// Expect empty plan after the plan after import is applied
			{
				Config: accconfig.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "authentication_methods", tfjson.ActionUpdate, new("[ALL]"), new("[]")),
						planchecks.ExpectChange(ref, "mfa_enrollment", tfjson.ActionUpdate, new(defaultMfaEnrollment), nil),
						planchecks.ExpectChange(ref, "client_types", tfjson.ActionUpdate, new("[ALL]"), new("[]")),
						planchecks.ExpectChange(ref, "security_integrations", tfjson.ActionUpdate, new("[ALL]"), new("[]")),
						planchecks.ExpectChange(ref, "client_policy", tfjson.ActionUpdate, new("[]"), new("[]")),
						planchecks.ExpectChange(ref, "mfa_policy", tfjson.ActionUpdate, new("[map[allowed_methods:[ALL] enforce_mfa_on_external_authentication:NONE]]"), new("[]")),
						planchecks.ExpectChange(ref, "pat_policy", tfjson.ActionUpdate, new("[map[default_expiry_in_days:15 max_expiry_in_days:365 network_policy_evaluation:ENFORCED_REQUIRED require_role_restriction_for_service_users:true]]"), new("[]")),
						planchecks.ExpectChange(ref, "workload_identity_policy", tfjson.ActionUpdate, new("[map[allowed_aws_accounts:[ALL] allowed_azure_issuers:[ALL] allowed_oidc_issuers:[ALL] allowed_providers:[ALL]]]"), new("[]")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.AuthenticationPolicyResource(t, ref).
						HasAuthenticationMethodsEmpty().
						HasMfaEnrollmentEmpty().
						HasClientTypesEmpty().
						HasSecurityIntegrationsEmpty(),
					assert.Check(resource.TestCheckResourceAttr(ref, "mfa_policy.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(ref, "pat_policy.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(ref, "workload_identity_policy.#", "0")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_enumSetCustomDiff(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelWithEnumSets := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(sdk.AuthenticationMethodsOptionPassword).
		WithClientTypes(sdk.ClientTypesOptionSnowflakeUi)
	modelWithLowercaseEnumSets := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(acchelpers.EnumToLower(sdk.AuthenticationMethodsOptionPassword)).
		WithClientTypes(acchelpers.EnumToLower(sdk.ClientTypesOptionSnowflakeUi))
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
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionNone).
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: sdk.MfaPolicyAllowedMethodsOptionTotp},
					{Value: sdk.MfaPolicyAllowedMethodsOptionDuo},
				}),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAzure},
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
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionNone),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAzure},
				}),
		)

	modelWithSwappedLists := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication(sdk.EnforceMfaOnExternalAuthenticationOptionNone).
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: sdk.MfaPolicyAllowedMethodsOptionDuo},
					{Value: sdk.MfaPolicyAllowedMethodsOptionTotp},
				}),
		).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: sdk.AllowedProviderOptionAzure},
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
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, ref).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[111111111111, 222222222222], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/one/v2.0, https://login.microsoftonline.com/two/v2.0], ALLOWED_OIDC_ISSUERS=[https://one.com, https://two.com]}")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithSomeListsEmpty),
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, ref).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithBasicLists),
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, ref).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
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
				Check: assertThat(
					t, resourceassert.AuthenticationPolicyResource(t, ref).
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
						HasOptions("").
						HasNoTargetScopes(),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_enrollment", testClient().SnowflakeDefaults.DefaultMfaEnrollment(t))),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.mfa_policy", "{ALLOWED_METHODS=[TOTP, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}")),
					assert.Check(resource.TestCheckResourceAttr(ref, "describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[AZURE], ALLOWED_AWS_ACCOUNTS=[111111111111, 222222222222], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/one/v2.0, https://login.microsoftonline.com/two/v2.0], ALLOWED_OIDC_ISSUERS=[https://one.com, https://two.com]}")),
				),
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
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithEnforceMfaOnExternalAuthentication("invalid"),
		)
	modelInvalidAllowedMethods := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaPolicy(
			*sdk.NewAuthenticationPolicyMfaPolicyRequest().
				WithAllowedMethods([]sdk.MfaPolicyAllowedMethodsOptionListItem{
					{Value: "invalid"},
				}),
		)
	modelInvalidDefaultExpiryInDays := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithDefaultExpiryInDays(0),
		)
	modelInvalidMaxExpiryInDays := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithMaxExpiryInDays(0),
		)
	modelInvalidNetworkPolicyEvaluation := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithPatPolicy(
			*sdk.NewAuthenticationPolicyPatPolicyRequest().
				WithNetworkPolicyEvaluation("invalid"),
		)
	modelInvalidAllowedProviders := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithWorkloadIdentityPolicy(
			*sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest().
				WithAllowedProviders([]sdk.AllowedProviderOptionListItem{
					{Value: "invalid"},
				}),
		)
	modelInvalidClientPolicyDriverType := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithClientTypes(sdk.ClientTypesOptionDrivers).
		WithClientPolicyValue(tfconfig.SetVariable(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"client_type":     tfconfig.StringVariable("invalid"),
				"minimum_version": tfconfig.StringVariable("1.0.0"),
			}),
		))
	modelClientPolicyIncompatibleClientTypes := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithClientTypes(sdk.ClientTypesOptionSnowflakeUi).
		WithClientPolicy(
			sdk.AuthenticationPolicyClientPolicyRequest{
				ClientType: sdk.ClientPolicyDriverTypeGoDriver,
				Params:     sdk.AuthenticationPolicyClientPolicyParamsRequest{MinimumVersion: "1.0.0"},
			},
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
				ExpectError: regexp.MustCompile(`invalid authentication methods option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMfaEnrollment),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid mfa enrollment option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidClientTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid client types option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidEnforceMfaOnExternalAuthentication),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid enforce mfa on external authentication option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAllowedMethods),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid mfa policy allowed methods option: INVALID`),
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
				ExpectError: regexp.MustCompile(`invalid allowed provider option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidClientPolicyDriverType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid client policy driver type: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelClientPolicyIncompatibleClientTypes),
				ExpectError: regexp.MustCompile(`can not contain CLIENT_POLICY`),
			},
		},
	})
}

// proves the fix for https://github.com/snowflakedb/terraform-provider-snowflake/issues/4557
func TestAcc_AuthenticationPolicy_migrateFromV2_14_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.AuthenticationPolicyResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.0"),
				Config:            accconfig.FromModels(t, providerModel, completeModel),
				ExpectError:       regexp.MustCompile("Error: object does not exist"),
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
