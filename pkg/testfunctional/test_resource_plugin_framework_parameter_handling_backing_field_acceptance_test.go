package testfunctional_test

import (
	"fmt"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const parameterHandlingBackingFieldDefaultValue = "default value - parameter handling backing field"

var parameterHandlingBackingFieldHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.ParameterHandlingBackingFieldOpts](
	testfunctional.ParameterHandlingBackingFieldOpts{
		StringValue: sdk.Pointer(parameterHandlingBackingFieldDefaultValue),
		Level:       string(sdk.ParameterTypeSnowflakeDefault),
	}, parameterHandlingBackingFieldOptsUseDefaultsForNil,
)

func parameterHandlingBackingFieldOptsUseDefaultsForNil(base testfunctional.ParameterHandlingBackingFieldOpts, defaults testfunctional.ParameterHandlingBackingFieldOpts, replaceWith testfunctional.ParameterHandlingBackingFieldOpts) testfunctional.ParameterHandlingBackingFieldOpts {
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
	allTestHandlers["parameter_handling_backing_field"] = parameterHandlingBackingFieldHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_ParameterHandling_BackingField(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_parameter_handling_backing_field", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	value := "some value"
	newValue := "new value"
	externalValue := "value changed externally"

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
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", value),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", value),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "OBJECT"),

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
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, value),
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
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "OBJECT"),

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
					parameterHandlingBackingFieldHandler.SetCurrentValue(testfunctional.ParameterHandlingBackingFieldOpts{
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
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "OBJECT"),

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
					parameterHandlingBackingFieldHandler.SetCurrentValue(testfunctional.ParameterHandlingBackingFieldOpts{
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
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "OBJECT"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "4"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.value", newValue),
				),
			},
			// remove the param from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(newValue), nil),
					},
				},
				Config: parameterHandlingBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", parameterHandlingBackingFieldDefaultValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", ""),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "5"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.value", "nil"),
				),
			},
			// import when param not in config (API default)
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					// when we import, we fill this value with what's available in API
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "string_value", optionalWithBackingFieldDefaultValue),
				),
			},
			// change the param value in config to API default (expecting action because of the different level)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, nil, sdk.String(optionalWithBackingFieldDefaultValue)),
					},
				},
				Config: parameterHandlingBackingFieldAllSetConfig(id, resourceType, optionalWithBackingFieldDefaultValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", parameterHandlingBackingFieldDefaultValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", parameterHandlingBackingFieldDefaultValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "OBJECT"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "6"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.value", optionalWithBackingFieldDefaultValue),
				),
			},
			// remove the param from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(optionalWithBackingFieldDefaultValue), nil),
					},
				},
				Config: parameterHandlingBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", parameterHandlingBackingFieldDefaultValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", ""),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "7"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.value", "nil"),
				),
			},
			// change level externally - change expected to be noop
			{
				PreConfig: func() {
					parameterHandlingBackingFieldHandler.SetCurrentValue(testfunctional.ParameterHandlingBackingFieldOpts{
						StringValue: sdk.Pointer(parameterHandlingBackingFieldDefaultValue),
						Level:       "IN_BETWEEN",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
				Config: parameterHandlingBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.value", parameterHandlingBackingFieldDefaultValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field.level", "IN_BETWEEN"),

					// no actions happened since last step (still 7)
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "7"),
				),
			},
		},
	})
}

func parameterHandlingBackingFieldAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func parameterHandlingBackingFieldNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
