package testfunctional_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var computedFieldCustomDiffHandler = common.NewDynamicHandler[testfunctional.ComputedFieldState]()

func init() {
	allTestHandlers[testfunctional.ComputedFieldCustomDiffPath] = computedFieldCustomDiffHandler
}

// TestAcc_SdkV2Functional_SetNewComputedOnComputedFieldTriggersUpdate verifies that
// calling SetNewComputed on a purely computed field from a CustomizeDiff triggers an Update
// when no config-driven fields actually change.
//
// The resource has:
//   - values: TypeSet (config-driven) — reordering in config produces no diff
//   - computed_order: TypeList (computed) — CustomizeDiff calls SetNewComputed on this conditionally
//   - update_count: tracks how many times Update was called
//
// This mimics the tag resource's allowed_values_order pattern where the server-side ordering
// of a TypeSet may change but the set itself (unordered) reports no diff.
func TestAcc_SdkV2Functional_SetNewComputedOnComputedFieldTriggersUpdate(t *testing.T) {
	resourceType := "snowflake_test_resource_computed_field_custom_diff"
	resourceName := "test"
	ref := fmt.Sprintf("%s.%s", resourceType, resourceName)

	testConfig := fmt.Sprintf(`
resource "%s" "%s" {
	provider = "%s"
	name     = "test_computed_diff"
	values   = ["a", "b", "c"]
}
`, resourceType, resourceName, SdkV2FunctionalTestsProviderName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Step 1: Create — update_count should be 0, trigger is off
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "0"),
					resource.TestCheckResourceAttr(ref, "computed_order.#", "3"),
				),
			},
			// Step 2: Same config. Set trigger_update=true on the server so
			// CustomizeDiff calls SetNewComputed("computed_order").
			// No config-driven fields change — only the computed field is marked.
			// Expect: Terraform plans an Update and calls the Update function.
			{
				PreConfig: func() {
					computedFieldCustomDiffHandler.SetCurrentValue(testfunctional.ComputedFieldState{
						Values:        []string{"a", "b", "c"},
						UpdateCount:   0,
						TriggerUpdate: true,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "1"),
				),
			},
			// Step 3: Same config again. trigger_update was cleared by Step 2's Update.
			// No config change, no trigger — expect no plan diff (noop).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "1"),
				),
			},
		},
	})
}

// TestAcc_SdkV2Functional_SetNewComputedOnComputedFieldTriggersUpdateRepeatedly verifies that
// SetNewComputed can trigger multiple consecutive Updates. After each Update clears the trigger,
// re-enabling it should produce another Update plan.
func TestAcc_SdkV2Functional_SetNewComputedOnComputedFieldTriggersUpdateRepeatedly(t *testing.T) {
	resourceType := "snowflake_test_resource_computed_field_custom_diff"
	resourceName := "test"
	ref := fmt.Sprintf("%s.%s", resourceType, resourceName)

	testConfig := fmt.Sprintf(`
resource "%s" "%s" {
	provider = "%s"
	name     = "test_computed_diff_repeated"
	values   = ["x", "y"]
}
`, resourceType, resourceName, SdkV2FunctionalTestsProviderName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Step 1: Create
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "0"),
				),
			},
			// Step 2: First trigger — expect Update
			{
				PreConfig: func() {
					computedFieldCustomDiffHandler.SetCurrentValue(testfunctional.ComputedFieldState{
						Values:        []string{"x", "y"},
						UpdateCount:   0,
						TriggerUpdate: true,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "1"),
				),
			},
			// Step 3: Re-enable trigger — expect another Update
			{
				PreConfig: func() {
					computedFieldCustomDiffHandler.SetCurrentValue(testfunctional.ComputedFieldState{
						Values:        []string{"x", "y"},
						UpdateCount:   1,
						TriggerUpdate: true,
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "update_count", "2"),
				),
			},
		},
	})
}
