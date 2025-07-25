package testfunctional_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// This test proves that in the framework, sensitive attributes are not exposed by default in the output.
// They can be accessed in the output block, but it must be properly marked.
func TestAcc_TerraformPluginFrameworkFunctional_NestedSensitive(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_nested_sensitive", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: nestedSensitiveConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "output.0.string_sensitive", "SECRET"),
				),
				ExpectError: regexp.MustCompile("Output refers to sensitive values"),
			},
			{
				Config: nestedSensitiveConfigMarked(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "output.0.string_sensitive", "SECRET"),
				),
			},
		},
	})
}

func nestedSensitiveConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return nestedSensitiveResourceConfig(id, resourceType) + fmt.Sprintf(`
output "nested_sensitive" {
  value = resource.%[2]s.test.output[0].string_sensitive
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func nestedSensitiveConfigMarked(id sdk.AccountObjectIdentifier, resourceType string) string {
	return nestedSensitiveResourceConfig(id, resourceType) + fmt.Sprintf(`
output "nested_sensitive" {
  value = resource.%[2]s.test.output[0].string_sensitive
  sensitive = true
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func nestedSensitiveResourceConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
