//go:build non_account_level_tests

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

// TestAcc_Experimental_ObjectParameter_UnsetOnDelete verifies that with the
// OBJECT_PARAMETER_UNSET_ON_DELETE experiment enabled, deleting the resource
// uses UNSET instead of resetting to the default value.
func TestAcc_Experimental_ObjectParameter_UnsetOnDelete(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource)).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ObjectParameterUnsetOnDelete)

	config := accconfig.FromModels(t, providerModel) + fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[1]s"
	value = "123456"
	object_type = "DATABASE"
	object_identifier {
		name = "%[2]s"
	}
}
`, sdk.DatabaseParameterUserTaskTimeoutMs, database.ID().Name())

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
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.DatabaseParameterUserTaskTimeoutMs))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "123456")),
					objectparametersassert.DatabaseParameters(t, database.ID()).
						HasUserTaskTimeoutMs(123456),
				),
			},
			{
				Destroy: true,
				Config:  config,
				Check: assertThat(
					t,
					objectparametersassert.DatabaseParameters(t, database.ID()).
						HasDefaultUserTaskTimeoutMsValue(),
				),
			},
		},
	})
}
