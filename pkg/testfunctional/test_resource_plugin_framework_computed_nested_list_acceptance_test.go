package testfunctional_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// proves https://github.com/hashicorp/terraform-plugin-framework/issues/1104
func TestAcc_TerraformPluginFrameworkFunctional_ComputedNestedList(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_computed_nested_list", PluginFrameworkFunctionalTestsProviderName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      computedNestedListConfig(id, resourceType),
				ExpectError: regexp.MustCompile("Value Conversion Error"),
			},
		},
	})
}

func computedNestedListConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
