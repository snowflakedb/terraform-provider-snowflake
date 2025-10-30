//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicies_AccountLevel(t *testing.T) {
	client := secondaryTestClient()
	client.BcrBundles.DisableBcrBundle(t, "2025_06")

	id := client.Ids.RandomSchemaObjectIdentifier()
	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	completeModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaAuthenticationMethods(sdk.MfaAuthenticationMethodsPassword, sdk.MfaAuthenticationMethodsSaml)
	providerModel := providermodel.SnowflakeProvider().
		WithProfile(testprofiles.Secondary)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		// TODO [SNOW-2324320]: secondary
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, providerModel, basicModel),
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, basicModel.ResourceReference()).
						HasMfaAuthenticationMethodsEmpty(),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[PASSWORD]")),
				),
			},
			{
				Config: accconfig.FromModels(t, providerModel, completeModel),
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasMfaAuthenticationMethods(sdk.MfaAuthenticationMethodsPassword, sdk.MfaAuthenticationMethodsSaml),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[SAML, PASSWORD]")),
				),
			},
			{
				Config: accconfig.FromModels(t, providerModel, completeModel),
				Taint:  []string{completeModel.ResourceReference()},
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasMfaAuthenticationMethods(sdk.MfaAuthenticationMethodsPassword, sdk.MfaAuthenticationMethodsSaml),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[SAML, PASSWORD]")),
				),
			},
			{
				PreConfig: func() {
					client.AuthenticationPolicy.Alter(t, sdk.NewAlterAuthenticationPolicyRequest(id).WithSet(
						*sdk.NewAuthenticationPolicySetRequest().
							WithMfaAuthenticationMethods([]sdk.MfaAuthenticationMethods{{Method: sdk.MfaAuthenticationMethodsSaml}}),
					))
				},
				Config: accconfig.FromModels(t, providerModel, completeModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, completeModel.ResourceReference()).
						HasMfaAuthenticationMethods(sdk.MfaAuthenticationMethodsPassword, sdk.MfaAuthenticationMethodsSaml),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(completeModel.ResourceReference(), "describe_output.0.mfa_authentication_methods", "[SAML, PASSWORD]")),
				),
			},
		},
	})
}
