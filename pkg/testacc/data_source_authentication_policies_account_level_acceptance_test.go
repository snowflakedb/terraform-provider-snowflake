//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicies_handling_with_builtin_policy_set_on_current_account(t *testing.T) {
	basicModel := datasourcemodel.AuthenticationPolicies("test").
		WithOnAccount()
	providerModel := providermodel.SnowflakeProvider().
		WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.AuthenticationPoliciesDatasource))

	policy := secondaryTestClient().AuthenticationPolicy.ShowOnCurrentAccount(t)
	if policy != nil && policy.Name != "BUILT-IN" {
		secondaryTestClient().Account.Alter(t, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{AuthenticationPolicy: sdk.Bool(true)},
		})
		t.Cleanup(func() {
			secondaryTestClient().Account.Alter(t, &sdk.AlterAccountOptions{
				Set: &sdk.AccountSet{AuthenticationPolicy: sdk.Pointer(policy.ID())},
			})
		})
	}

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
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, providerModel, basicModel),
			},
		},
	})
}
