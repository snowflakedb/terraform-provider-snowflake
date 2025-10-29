//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicies(t *testing.T) {
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
				{Method: sdk.MfaPolicyPassAllowedMethodPassKey},
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

	authenticationPoliciesModel := datasourcemodel.AuthenticationPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())
	authenticationPoliciesModelWithoutOptionals := datasourcemodel.AuthenticationPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel, authenticationPoliciesModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.#", "1")),

					resourceshowoutputassert.AuthenticationPoliciesDatasourceShowOutput(t, "snowflake_authentication_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.authentication_methods", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.client_types", "[SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.security_integrations", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.mfa_enrollment", "REQUIRED")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.mfa_authentication_methods", "[PASSWORD, SAML]")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.mfa_policy", "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.pat_policy", "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=false}")),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModel.DatasourceReference(), "authentication_policies.0.describe_output.0.workload_identity_policy", "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[111122223333], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com]}")),
				),
			},
			{
				Config: accconfig.FromModels(t, completeModel, authenticationPoliciesModelWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModelWithoutOptionals.DatasourceReference(), "authentication_policies.#", "1")),
					resourceshowoutputassert.AuthenticationPoliciesDatasourceShowOutput(t, "snowflake_authentication_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(authenticationPoliciesModelWithoutOptionals.DatasourceReference(), "authentication_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicies_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()
	currentUserId := testClient().Context.CurrentUser(t)

	model1 := model.AuthenticationPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name())
	model2 := model.AuthenticationPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name())
	model3 := model.AuthenticationPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name())
	authenticationPoliciesModelLikeFirstOne := datasourcemodel.AuthenticationPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	authenticationPoliciesModelLikePrefix := datasourcemodel.AuthenticationPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	authenticationPoliciesModelOnUser := datasourcemodel.AuthenticationPolicies("test").
		WithOnUser(currentUserId).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	authenticationPoliciesModelOnAccount := datasourcemodel.AuthenticationPolicies("test").
		WithOnAccount().
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, authenticationPoliciesModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authenticationPoliciesModelLikeFirstOne.DatasourceReference(), "authentication_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, authenticationPoliciesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authenticationPoliciesModelLikePrefix.DatasourceReference(), "authentication_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, authenticationPoliciesModelOnUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authenticationPoliciesModelOnUser.DatasourceReference(), "authentication_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, authenticationPoliciesModelOnAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authenticationPoliciesModelOnAccount.DatasourceReference(), "authentication_policies.#", "1"),
				),
			},
		},
	})
}

func TestAcc_AuthenticationPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.AuthenticationPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_AuthenticationPolicies_emptyOn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.AuthenticationPolicies("test").WithEmptyOn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_AuthenticationPolicies_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_AuthenticationPolicies/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one authentication policy"),
			},
		},
	})
}
