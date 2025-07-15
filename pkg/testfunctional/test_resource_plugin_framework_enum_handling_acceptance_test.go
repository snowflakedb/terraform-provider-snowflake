package testfunctional_test

import (
	"fmt"
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

	value := string(testfunctional.SomeEnumTypeVersion1)
	valueLowercased := strings.ToLower(value)
	newValue := string(testfunctional.SomeEnumTypeVersion2)
	externalValue := string(testfunctional.SomeEnumTypeVersion3)

	_ = externalValue

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
				Config: enumHandlingAllSetConfig(id, resourceType, value),
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
			// import when type in config
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "string_value", value),
				),
			},
			// change type in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, &value, &newValue),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, newValue),
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
			// remove type from config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, &newValue, nil),
					},
				},
				Config: enumHandlingNotSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),
					// TODO: backing field

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.value", "nil"),
				),
			},
			// import when no type in config
			{
				ResourceName: resourceReference,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.FullyQualifiedName(), "string_value", value),
				),
			},
			// add config (lower case)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, nil, &valueLowercased),
					},
				},
				Config: enumHandlingAllSetConfig(id, resourceType, valueLowercased),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", valueLowercased),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "4"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.value", value),
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
