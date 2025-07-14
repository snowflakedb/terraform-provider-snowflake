package testfunctional_test

import (
	"fmt"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const enumHandlingDefaultValue = sdk.WarehouseType("STANDARD")

var enumHandlingHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.EnumHandlingOpts](
	testfunctional.EnumHandlingOpts{StringValue: sdk.Pointer(enumHandlingDefaultValue)}, enumHandlingOptsUseDefaultsForNil,
)

func enumHandlingOptsUseDefaultsForNil(base testfunctional.EnumHandlingOpts, defaults testfunctional.EnumHandlingOpts, replaceWith testfunctional.EnumHandlingOpts) testfunctional.EnumHandlingOpts {
	if replaceWith.StringValue == nil {
		base.StringValue = defaults.StringValue
	} else {
		base.StringValue = replaceWith.StringValue
	}
	return base
}

func init() {
	allTestHandlers["enum_handling"] = enumHandlingHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_EnumHandling(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_enum_handling", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	value := string(sdk.WarehouseTypeStandard)
	newValue := "new value"
	externalValue := "value changed externally"

	_ = newValue
	_ = externalValue
	_ = enumHandlingNotSetConfig

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
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseTypeStandard))),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, string(sdk.WarehouseTypeStandard)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", value),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.value", value),
				),
			},
		},
	})
}

func enumHandlingAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func enumHandlingNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
