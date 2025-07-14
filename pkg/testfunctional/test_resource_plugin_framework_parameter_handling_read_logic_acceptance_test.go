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

const parameterHandlingReadLogicDefaultValue = "default value"

var parameterHandlingReadLogicHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.ParameterHandlingReadLogicOpts](
	testfunctional.ParameterHandlingReadLogicOpts{
		StringValue: sdk.Pointer(parameterHandlingReadLogicDefaultValue),
		Level:       string(sdk.ParameterTypeSnowflakeDefault),
	}, parameterHandlingReadLogicOptsUseDefaultsForNil,
)

func parameterHandlingReadLogicOptsUseDefaultsForNil(base testfunctional.ParameterHandlingReadLogicOpts, defaults testfunctional.ParameterHandlingReadLogicOpts, replaceWith testfunctional.ParameterHandlingReadLogicOpts) testfunctional.ParameterHandlingReadLogicOpts {
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
	allTestHandlers["parameter_handling_read_logic"] = parameterHandlingReadLogicHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_ParameterHandling_ReadLogic(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_parameter_handling_read_logic", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	value := "some value"
	newValue := "new value"
	externalValue := "value changed externally"

	_, _ = newValue, externalValue
	_ = parameterHandlingReadLogicNotSetConfig(id, resourceType)

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
				Config: parameterHandlingReadLogicAllSetConfig(id, resourceType, value),
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
			// do not make any change (to check if there is no drift)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: parameterHandlingReadLogicAllSetConfig(id, resourceType, value),
			},
			// import when known value
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring actions_log as they serve testing purpose; ignoring name as we do not fill it in read (import tests will be done separately).
				ImportStateVerifyIgnore: []string{"actions_log", "name"},
			},
			// change the param value in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(value), sdk.String(newValue)),
					},
				},
				Config: parameterHandlingReadLogicAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.value", newValue),
				),
			},
			// change the param value externally
			{
				PreConfig: func() {
					parameterHandlingReadLogicHandler.SetCurrentValue(testfunctional.ParameterHandlingReadLogicOpts{
						StringValue: &externalValue,
						Level:       "OBJECT",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", sdk.String(newValue), sdk.String(externalValue)),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(externalValue), sdk.String(newValue)),
					},
				},
				Config: parameterHandlingReadLogicAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.value", newValue),
				),
			},
			//  change the param value externally to the value from config (but on different level)
			{
				PreConfig: func() {
					parameterHandlingReadLogicHandler.SetCurrentValue(testfunctional.ParameterHandlingReadLogicOpts{
						StringValue: &newValue,
						Level:       "OTHER",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", sdk.String(newValue), nil),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, nil, sdk.String(newValue)),
						planchecks.ExpectComputed(resourceReference, "string_value", false),
					},
				},
				Config: parameterHandlingReadLogicAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "4"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.value", newValue),
				),
			},
		},
	})
}

func parameterHandlingReadLogicAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func parameterHandlingReadLogicNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
