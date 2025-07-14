package testfunctional_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var stringWithBackingValueHandler = common.NewDynamicHandlerWithInitialValueAndReplaceWithFunc[testfunctional.StringWithMetadataOpts](
	testfunctional.StringWithMetadataOpts{}, common.AlwaysReplace,
)

func init() {
	allTestHandlers["string_with_metadata"] = stringWithBackingValueHandler
}

// Test proving that even optional computed fields cannot be altered if provided in the tf config.
func TestAcc_TerraformPluginFrameworkFunctional_CustomTypes_StringWithBackingField(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_string_with_metadata", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with known value
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionCreate),
					},
				},
				Config:      stringWithMetadataAllSetConfig(id, resourceType),
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
		},
	})
}

func stringWithMetadataAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  string_value = "some_value"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
