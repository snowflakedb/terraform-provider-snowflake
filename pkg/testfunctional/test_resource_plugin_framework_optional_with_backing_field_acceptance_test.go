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
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldAllSetConfig(id, resourceType, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", value),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", value),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.value", value),
				),
			},
			// import when known value
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring actions_log as they serve testing purpose; ignoring name as we do not fill it in read (import tests will be done separately).
				ImportStateVerifyIgnore: []string{"actions_log", "name"},
			},
			// remove value from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(value), nil),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", optionalWithBackingFieldDefaultValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.value", "nil"),
				),
			},
			// import when unset
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					// when we import, we fill this value with what's available in API
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "string_value", optionalWithBackingFieldDefaultValue),
				),
			},
			// change externally when absent in config
			{
				PreConfig: func() {
					optionalWithBackingFieldHandler.SetCurrentValue(testfunctional.OptionalWithBackingFieldOpts{
						StringValue: &externalValue,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", nil, sdk.String(externalValue)),
						planchecks.ExpectDrift(resourceReference, "string_value_backing_field", sdk.String(optionalWithBackingFieldDefaultValue), sdk.String(externalValue)),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(externalValue), nil),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", optionalWithBackingFieldDefaultValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.value", "nil"),
				),
			},
			// import when unset and external value is set
			{
				PreConfig: func() {
					optionalWithBackingFieldHandler.SetCurrentValue(testfunctional.OptionalWithBackingFieldOpts{
						StringValue: &externalValue,
					})
				},
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					// when we import, we fill this value with what's available in API
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "string_value", externalValue),
				),
			},
			// unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(externalValue), nil),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", optionalWithBackingFieldDefaultValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "4"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.value", "nil"),
				),
			},
			// set the value back again
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, nil, sdk.String(value)),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldAllSetConfig(id, resourceType, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", value),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", value),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "5"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.value", value),
				),
			},
			// change the value
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(value), sdk.String(newValue)),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", newValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "6"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.value", newValue),
				),
			},
			// react to external change
			{
				PreConfig: func() {
					optionalWithBackingFieldHandler.SetCurrentValue(testfunctional.OptionalWithBackingFieldOpts{
						StringValue: &externalValue,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", sdk.String(newValue), sdk.String(externalValue)),
						planchecks.ExpectDrift(resourceReference, "string_value_backing_field", sdk.String(newValue), sdk.String(externalValue)),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(externalValue), sdk.String(newValue)),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", newValue),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "7"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.value", newValue),
				),
			},
			// remove type from config but update externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					optionalWithBackingFieldHandler.SetCurrentValue(testfunctional.OptionalWithBackingFieldOpts{
						StringValue: sdk.Pointer(optionalWithBackingFieldDefaultValue),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "string_value", sdk.String(newValue), sdk.String(optionalWithBackingFieldDefaultValue)),
						planchecks.ExpectDrift(resourceReference, "string_value_backing_field", sdk.String(newValue), sdk.String(optionalWithBackingFieldDefaultValue)),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.String(optionalWithBackingFieldDefaultValue), nil),
						planchecks.ExpectComputed(resourceReference, "string_value_backing_field", true),
					},
				},
				Config: optionalWithBackingFieldNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "string_value_backing_field", optionalWithBackingFieldDefaultValue),

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

func optionalWithBackingFieldAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  string_value = "%[2]s"
}
`, id.Name(), value, resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func optionalWithBackingFieldNotSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
