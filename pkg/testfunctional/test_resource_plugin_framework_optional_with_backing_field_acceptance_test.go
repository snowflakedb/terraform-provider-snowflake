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

const optionalWithBackingFieldDefaultValue = "default value"

var optionalWithBackingFieldHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.OptionalWithBackingFieldOpts](
	testfunctional.OptionalWithBackingFieldOpts{StringValue: sdk.Pointer(optionalWithBackingFieldDefaultValue)}, optionalWithBackingFieldOptsUseDefaultsForNil,
)

func optionalWithBackingFieldOptsUseDefaultsForNil(base testfunctional.OptionalWithBackingFieldOpts, defaults testfunctional.OptionalWithBackingFieldOpts, replaceWith testfunctional.OptionalWithBackingFieldOpts) testfunctional.OptionalWithBackingFieldOpts {
	if replaceWith.StringValue == nil {
		base.StringValue = defaults.StringValue
	} else {
		base.StringValue = replaceWith.StringValue
	}
	return base
}

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
			// remove value from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: optionalWithBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", "default value"),
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

func optionalWithBackingFieldNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
