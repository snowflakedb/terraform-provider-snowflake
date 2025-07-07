package testfunctional_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TerraformPluginFrameworkFunctional_ZeroValues_Basic(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_zero_values", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: zeroValuesConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.value", "true"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.value", "0"),
				),
			},
		},
	})
}

func zeroValuesConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  bool_value = "true"
  int_value = 0
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
