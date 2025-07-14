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

const parameterHandlingResourcePlanModifierDefaultValue = "default value"

var parameterHandlingResourcePlanModifierHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.ParameterHandlingResourcePlanModifierOpts](
	testfunctional.ParameterHandlingResourcePlanModifierOpts{
		StringValue: sdk.Pointer(parameterHandlingResourcePlanModifierDefaultValue),
		Level:       string(sdk.ParameterTypeSnowflakeDefault),
	}, parameterHandlingResourcePlanModifierOptsUseDefaultsForNil,
)

func parameterHandlingResourcePlanModifierOptsUseDefaultsForNil(base testfunctional.ParameterHandlingResourcePlanModifierOpts, defaults testfunctional.ParameterHandlingResourcePlanModifierOpts, replaceWith testfunctional.ParameterHandlingResourcePlanModifierOpts) testfunctional.ParameterHandlingResourcePlanModifierOpts {
	if replaceWith.StringValue == nil {
		base.StringValue = defaults.StringValue
		base.Level = string(sdk.ParameterTypeSnowflakeDefault)
	} else {
		base.StringValue = replaceWith.StringValue
		base.Level = "OBJECT"
	}
	return base
}

func init() {
	allTestHandlers["parameter_handling_resource_plan_modifier"] = parameterHandlingResourcePlanModifierHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_ParameterHandling_ResourcePlanModifier(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_parameter_handling_resource_plan_modifier", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	value := "some value"
	newValue := "new value"
	externalValue := "value changed externally"

	_, _ = newValue, externalValue
	_ = parameterHandlingResourcePlanModifierNotSetConfig(id, resourceType)

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
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionCreate, nil, sdk.String(value)),
					},
				},
				Config: parameterHandlingResourcePlanModifierAllSetConfig(id, resourceType, value),
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

func parameterHandlingResourcePlanModifierAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func parameterHandlingResourcePlanModifierNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
