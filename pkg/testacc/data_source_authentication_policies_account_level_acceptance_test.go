//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicies_handling_with_builtin_policy_set_on_current_account(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakeNonProdEnvironment {
		t.Skip("Missing snowflake defaults configuration on prod environment")
	}

	basicModel := datasourcemodel.AuthenticationPolicies("test").
		WithOnAccount()
	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.SnowflakeDefaults)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.11.0"),
				Config:            accconfig.FromModels(t, providerModel, basicModel),
				ExpectError:       regexp.MustCompile("Error: sql: Scan error on column index 0, name \"created_on\""),
			},
			{
				ProtoV6ProviderFactories: snowflakeDefaultsAccountProviderFactory,
				Config:                   accconfig.FromModels(t, providerModel, basicModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(basicModel.DatasourceReference(), "authentication_policies.0.show_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicModel.DatasourceReference(), "authentication_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}
