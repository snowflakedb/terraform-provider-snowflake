package testfunctional_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var optionalWithBackingFieldHandler = common.NewDynamicHandlerWithInitialValueAndReplaceWithFunc[testfunctional.OptionalWithBackingFieldOpts](
	testfunctional.OptionalWithBackingFieldOpts{}, common.AlwaysReplace,
)

func init() {
	allTestHandlers["optional_with_backing_field"] = optionalWithBackingFieldHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_OptionalWithBackingField(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_optional_with_backing_field", PluginFrameworkFunctionalTestsProviderName)
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
				Config: optionalWithBackingFieldAllSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", "some_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", "some_value"),
				),
			},
		},
	})
}

func optionalWithBackingFieldAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  string_value = "some_value"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
