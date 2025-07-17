package testfunctional_test

import (
	"fmt"
	"regexp"
	"strings"
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

const enumHandlingDefaultValue = testfunctional.SomeEnumTypeVersion1

var enumHandlingHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.EnumHandlingOpts](
	testfunctional.EnumHandlingOpts{EnumValue: sdk.Pointer(enumHandlingDefaultValue)}, enumHandlingOptsUseDefaultsForNil,
)

func enumHandlingOptsUseDefaultsForNil(base testfunctional.EnumHandlingOpts, defaults testfunctional.EnumHandlingOpts, replaceWith testfunctional.EnumHandlingOpts) testfunctional.EnumHandlingOpts {
	if replaceWith.EnumValue == nil {
		base.EnumValue = defaults.EnumValue
	} else {
		base.EnumValue = replaceWith.EnumValue
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

	value := string(enumHandlingDefaultValue)
	newValue := string(testfunctional.SomeEnumTypeVersion2)
	newValueLowercased := strings.ToLower(newValue)
	externalValueEnum := testfunctional.SomeEnumTypeVersion3
	externalValue := string(externalValueEnum)

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
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionCreate, nil, sdk.String(value)),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "enum_value", value),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", value),
				),
			},
			// import when type in config
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "enum_value", value),
				),
			},
			// change type in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionUpdate, &value, &newValue),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "enum_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", newValue),
				),
			},
			// remove type from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionUpdate, &newValue, nil),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
					},
				},
				Config: enumHandlingNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "enum_value"),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", string(enumHandlingDefaultValue)),
				),
			},
			// import when no type in config
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "enum_value", value),
				),
			},
			// add config (lower case)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionUpdate, nil, &newValueLowercased),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, newValueLowercased),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "enum_value", newValueLowercased),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", newValue),
				),
			},
			// change config to upper case - expect no changes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", false),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "enum_value", newValueLowercased),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", newValue),
				),
			},
			// change the type externally
			{
				PreConfig: func() {
					enumHandlingHandler.SetCurrentValue(testfunctional.EnumHandlingOpts{
						EnumValue: &externalValueEnum,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "enum_value", &newValueLowercased, &externalValue),
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionUpdate, &externalValue, &newValue),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, newValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "enum_value", newValue),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", newValue),
				),
			},
			// remove type from config but update enum value externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					enumHandlingHandler.SetCurrentValue(testfunctional.EnumHandlingOpts{
						EnumValue: sdk.Pointer(enumHandlingDefaultValue),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(resourceReference, "enum_value_backing_field", true),
						planchecks.ExpectDrift(resourceReference, "enum_value", &newValue, &value),
						planchecks.ExpectChange(resourceReference, "enum_value", tfjson.ActionUpdate, &value, nil),
					},
				},
				Config: enumHandlingNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "enum_value"),
					resource.TestCheckResourceAttr(resourceReference, "enum_value_backing_field", string(enumHandlingDefaultValue)),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkFunctional_EnumHandling_Validations(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_enum_handling", PluginFrameworkFunctionalTestsProviderName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with invalid value
			{
				Config:      enumHandlingAllSetConfig(id, resourceType, "unknown"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid some enum type: unknown"),
			},
		},
	})
}

func enumHandlingAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string, value string) string {
	return fmt.Sprintf(`
resource "%[3]s" "test" {
  provider = "%[4]s"

  name = "%[1]s"
  enum_value = "%[2]s"
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
