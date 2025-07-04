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
func TestAcc_TerraformPluginFrameworkFunctional_Error(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_computed_nested_list", PluginFrameworkFunctionalTestsProviderName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      computedNestedListConfig(id, resourceType, "STRUCT"),
				ExpectError: regexp.MustCompile("Value Conversion Error"),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkFunctional_Explicit(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_computed_nested_list", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: computedNestedListConfig(id, resourceType, "EXPLICIT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.action", "SOME ACTION"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.field", "ON FIELD"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.value", "WITH VALUE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.action", "SOME OTHER ACTION"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.field", "ON OTHER FIELD"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.value", "WITH OTHER VALUE"),
				),
			},
		},
	})
}

func computedNestedListConfig(id sdk.AccountObjectIdentifier, resourceType string, option string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  option = "%[2]s"
}
`, id.Name(), option, resourceType, PluginFrameworkFunctionalTestsProviderName)
}
