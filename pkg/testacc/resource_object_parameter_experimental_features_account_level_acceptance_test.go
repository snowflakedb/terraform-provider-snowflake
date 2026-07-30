//go:build account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_ObjectParameter_UnsetOnDelete_OnAccount verifies that with the
// OBJECT_PARAMETER_UNSET_ON_DELETE experiment enabled, deleting an account-level parameter
// uses UNSET instead of resetting to the default value.
func TestAcc_Experimental_ObjectParameter_UnsetOnDelete_OnAccount(t *testing.T) {
	t.Cleanup(func() { testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterUserTaskTimeoutMs) })

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource)).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ObjectParameterUnsetOnDelete)

	config := accconfig.FromModels(t, providerModel) + fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[1]s"
	value = "123123"
	on_account = true
}
`, sdk.ObjectParameterUserTaskTimeoutMs)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: objectParameterUnsetOnDeleteProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.ObjectParameterUserTaskTimeoutMs))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "123123")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "true")),
					objectparametersassert.AccountParameters(t, testClient().Context.CurrentAccountId(t)).
						HasUserTaskTimeoutMs(123123),
				),
			},
			{
				Destroy: true,
				Config:  config,
				Check: assertThat(
					t,
					objectparametersassert.AccountParameters(t, testClient().Context.CurrentAccountId(t)).
						HasDefaultUserTaskTimeoutMsValue(),
				),
			},
		},
	})
}
