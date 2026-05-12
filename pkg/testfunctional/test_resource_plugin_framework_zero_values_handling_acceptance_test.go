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

var zeroValuesHandler = common.NewDynamicHandlerWithInitialValueAndReplaceWithFunc[testfunctional.ZeroValuesOpts](
	testfunctional.ZeroValuesOpts{}, common.AlwaysReplace,
)

// TODO [SNOW-2296391]: handle by reflection or generate (keeping it for the Optional+Computed test)
func zeroValuesOptsReplaceWithNonNil(base testfunctional.ZeroValuesOpts, replaceWith testfunctional.ZeroValuesOpts) testfunctional.ZeroValuesOpts {
	if replaceWith.BoolValue != nil {
		base.BoolValue = replaceWith.BoolValue
	}
	if replaceWith.IntValue != nil {
		base.IntValue = replaceWith.IntValue
	}
	if replaceWith.StringValue != nil {
		base.StringValue = replaceWith.StringValue
	}
	return base
}

func init() {
	allTestHandlers["zero_values_handling"] = zeroValuesHandler
}

// This test verifies the behavior of Optional fields with possible zero-values (based on TestAcc_Warehouse_ZeroValues).
func TestAcc_TerraformPluginFrameworkFunctional_ZeroValues_Optional(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_zero_values", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with valid "zero" values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionCreate),
					},
				},
				Config: zeroValuesAllSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "bool_value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "int_value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "string_value", ""),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.0.value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.1.value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.action", "CREATE"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.2.value", ""),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: zeroValuesNoneSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceReference, "bool_value"),
					resource.TestCheckNoResourceAttr(resourceReference, "int_value"),
					resource.TestCheckNoResourceAttr(resourceReference, "string_value"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "6"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.3.value", "nil"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.4.value", "nil"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.action", "UPDATE - UNSET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.5.value", "nil"),
				),
			},
			// import when empty
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring actions_log as they serve testing purpose; ignoring name as we do not fill it in read (import tests will be done separately).
				ImportStateVerifyIgnore: []string{"actions_log", "name"},
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: zeroValuesAllSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "bool_value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "int_value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "string_value", ""),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "9"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.6.value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.7.value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.8.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.8.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.8.value", ""),
				),
			},
			// import zero values
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring actions_log as they serve testing purpose; ignoring name as we do not fill it in read (import tests will be done separately).
				ImportStateVerifyIgnore: []string{"actions_log", "name"},
			},
			// set externally to non-zero values
			{
				PreConfig: func() {
					zeroValuesHandler.SetCurrentValue(testfunctional.ZeroValuesOpts{
						BoolValue:   sdk.Pointer(true),
						IntValue:    sdk.Pointer(10),
						StringValue: sdk.Pointer("some external text"),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: zeroValuesAllSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "bool_value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "int_value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "string_value", ""),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "12"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.9.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.9.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.9.value", "false"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.10.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.10.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.10.value", "0"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.11.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.11.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.11.value", ""),
				),
			},
			// set to non-zero values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: zeroValuesNonZeroValuesConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "bool_value", "true"),
					resource.TestCheckResourceAttr(resourceReference, "int_value", "10"),
					resource.TestCheckResourceAttr(resourceReference, "string_value", "some text"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "15"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.12.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.12.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.12.value", "true"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.13.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.13.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.13.value", "10"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.14.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.14.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.14.value", "some text"),
				),
			},
			// set externally to zero values
			{
				PreConfig: func() {
					zeroValuesHandler.SetCurrentValue(testfunctional.ZeroValuesOpts{
						BoolValue:   sdk.Pointer(false),
						IntValue:    sdk.Pointer(0),
						StringValue: sdk.Pointer(""),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Config: zeroValuesNonZeroValuesConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "bool_value", "true"),
					resource.TestCheckResourceAttr(resourceReference, "int_value", "10"),
					resource.TestCheckResourceAttr(resourceReference, "string_value", "some text"),

					// check actions
					resource.TestCheckResourceAttr(resourceReference, "actions_log.#", "18"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.15.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.15.field", "bool_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.15.value", "true"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.16.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.16.field", "int_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.16.value", "10"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.17.action", "UPDATE - SET"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.17.field", "string_value"),
					resource.TestCheckResourceAttr(resourceReference, "actions_log.17.value", "some text"),
				),
			},
		},
	})
}

func zeroValuesAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  bool_value = false
  int_value = 0
  string_value = ""
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func zeroValuesNoneSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func zeroValuesNonZeroValuesConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  bool_value = true
  int_value = 10
  string_value = "some text"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
