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

const optionalComputedHandlingPrivateDefaultValue = "default value - optional computed handling"

var optionalComputedHandler = common.NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[testfunctional.OptionalComputedOpts](
	testfunctional.OptionalComputedOpts{
		StringValue: sdk.Pointer(optionalComputedHandlingPrivateDefaultValue),
	}, optionalComputedOptsReplaceWithNonNil,
)

func optionalComputedOptsReplaceWithNonNil(base testfunctional.OptionalComputedOpts, defaultValue testfunctional.OptionalComputedOpts, replaceWith testfunctional.OptionalComputedOpts) testfunctional.OptionalComputedOpts {
	if replaceWith.StringValue != nil {
		base.StringValue = replaceWith.StringValue
	} else {
		base.StringValue = defaultValue.StringValue
	}
	return base
}

func init() {
	allTestHandlers["optional_computed_handling"] = optionalComputedHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_OptionalComputed(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_optional_computed", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	value := "some value"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// create with value
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionCreate),
					},
				},
				Config: optionalComputedAllSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", value),
				),
			},
			// remove all from config (to validate that unset is run correctly and default value is in state)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, &value, nil),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.Pointer(optionalComputedHandlingPrivateDefaultValue), nil),
					},
				},
				Config: optionalComputedNoneSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", optionalComputedHandlingPrivateDefaultValue),
				),
				// Because our plan modifier react to the situation config == null and state != null, it will be a perma diff.
				ExpectNonEmptyPlan: true,
			},
			// Adding the same step to show it ends the same way (hence - perma diff)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.Pointer(optionalComputedHandlingPrivateDefaultValue), nil),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "string_value", tfjson.ActionUpdate, sdk.Pointer(optionalComputedHandlingPrivateDefaultValue), nil),
					},
				},
				Config: optionalComputedNoneSetConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "string_value", optionalComputedHandlingPrivateDefaultValue),
				),
				// Because our plan modifier react to the situation config == null and state != null, it will be a perma diff.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func optionalComputedAllSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
  string_value = "some value"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}

func optionalComputedNoneSetConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
