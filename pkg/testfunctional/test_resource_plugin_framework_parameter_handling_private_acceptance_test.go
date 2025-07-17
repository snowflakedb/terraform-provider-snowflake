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

const parameterHandlingPrivateDefaultValue = "default value - parameter handling private"

var parameterHandlingPrivateHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.ParameterHandlingPrivateOpts](
	testfunctional.ParameterHandlingPrivateOpts{
		StringValue: sdk.Pointer(parameterHandlingPrivateDefaultValue),
		Level:       string(sdk.ParameterTypeSnowflakeDefault),
	}, parameterHandlingPrivateOptsUseDefaultsForNil,
)

func parameterHandlingPrivateOptsUseDefaultsForNil(base testfunctional.ParameterHandlingPrivateOpts, defaults testfunctional.ParameterHandlingPrivateOpts, replaceWith testfunctional.ParameterHandlingPrivateOpts) testfunctional.ParameterHandlingPrivateOpts {
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
	allTestHandlers["parameter_handling_private"] = parameterHandlingPrivateHandler
}

// TODO [mux-PRs]: can we assert private bytes?
func TestAcc_TerraformPluginFrameworkFunctional_ParameterHandling_Private(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_parameter_handling_private", PluginFrameworkFunctionalTestsProviderName)
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
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, value),
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
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, value),
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
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, newValue),
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
					parameterHandlingPrivateHandler.SetCurrentValue(testfunctional.ParameterHandlingPrivateOpts{
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
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, newValue),
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
					parameterHandlingPrivateHandler.SetCurrentValue(testfunctional.ParameterHandlingPrivateOpts{
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
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, newValue),
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
			// remove the param from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(newValue), nil),
					},
				},
				Config: parameterHandlingPrivateNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "5"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.value", "nil"),
				),
			},
			// import when param not in config (API default)
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring actions_log as they serve testing purpose; ignoring name as we do not fill it in read (import tests will be done separately).
				ImportStateVerifyIgnore: []string{"actions_log", "name"},
			},
			// change the param value in config to API default (expecting action because of the different level)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, nil, sdk.String(parameterHandlingPrivateDefaultValue)),
					},
				},
				Config: parameterHandlingPrivateAllSetConfig(id, resourceType, parameterHandlingPrivateDefaultValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", parameterHandlingPrivateDefaultValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "6"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.value", parameterHandlingPrivateDefaultValue),
				),
			},
			// remove the param from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(parameterHandlingPrivateDefaultValue), nil),
					},
				},
				Config: parameterHandlingPrivateNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),

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
					parameterHandlingPrivateHandler.SetCurrentValue(testfunctional.ParameterHandlingPrivateOpts{
						StringValue: sdk.Pointer(parameterHandlingPrivateDefaultValue),
						Level:       "IN_BETWEEN",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
				Config: parameterHandlingPrivateNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),

					// no actions happened since last step (still 7)
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "7"),
				),
			},
			// change level externally to OBJECT level
			{
				PreConfig: func() {
					parameterHandlingPrivateHandler.SetCurrentValue(testfunctional.ParameterHandlingPrivateOpts{
						StringValue: sdk.Pointer(parameterHandlingPrivateDefaultValue),
						Level:       "OBJECT",
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", nil, sdk.String(parameterHandlingPrivateDefaultValue)),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(parameterHandlingPrivateDefaultValue), nil),
					},
				},
				Config: parameterHandlingPrivateNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "8"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.value", "nil"),
				),
			},
		},
	})
}

func parameterHandlingPrivateAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func parameterHandlingPrivateNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
