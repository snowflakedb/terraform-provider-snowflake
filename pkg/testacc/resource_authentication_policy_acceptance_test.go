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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
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
		WithSecurityIntegrations("ALL")
	completeModelWithDifferentValues := model.AuthenticationPolicy("test", id2.DatabaseName(), id2.SchemaName(), id2.Name()).
		WithComment(changedComment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsSaml).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequiredPasswordOnly).
		WithClientTypes(sdk.ClientTypesSnowflakeCli).
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
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", "null")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
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
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
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
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.name", id2.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.authentication_methods", "[SAML]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.client_types", "[SNOWFLAKE_CLI]")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.security_integrations", fmt.Sprintf("[%s]", samlIntegration.ID().Name()))),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_enrollment", "REQUIRED_PASSWORD_ONLY")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
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
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
				),
			},
			{
				Config:            accconfig.FromModels(t, completeModel),
				ResourceName:      completeModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAcc_AuthenticationPolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidAuthenticationMethods := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods("invalid")
	modelInvalidMfaEnrollment := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaEnrollmentEnum("invalid")
	modelInvalidClientTypes := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithClientTypes("invalid")

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
		},
	})
}
