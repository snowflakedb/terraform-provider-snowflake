package testfunctional_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_SdkV2Functional_TestResource_ListNullValueLogicWithHelperField tests that adding a
// nullable_list_presence computed helper field with CustomizeDiff solves the null <-> empty
// limitation demonstrated in TestAcc_SdkV2Functional_TestResource_ListNullValueLogic.
//
// The helper field stores the config's list presence (null, empty, or items) as a computed string.
// CustomizeDiff compares the raw config presence against the state value and forces a plan diff
// when they diverge, triggering Update even for null <-> empty transitions.
//
// Expected behavior with helper field (compare with base resource):
//
//	| Transition       | Base Resource | With Helper Field |
//	|------------------|--------------|-------------------|
//	| null  → empty    | Noop         | Update            |
//	| empty → null     | Noop         | Update            |
//	| null  → filled   | Update       | Update            |
//	| empty → filled   | Update       | Update            |
//	| filled → empty   | Update       | Update            |
//	| filled → null    | Update       | Update            |
//	| ext. null→empty  | Noop         | Update            |
//	| ext. empty→null  | Noop         | Update            |
func TestAcc_SdkV2Functional_TestResource_ListNullValueLogicWithHelperField(t *testing.T) {
	envName := fmt.Sprintf("%s_%s", testenvs.TestResourceNullListHandlingEnv, strings.ToUpper(random.AlphaN(10)))
	resourceType := "snowflake_test_resource_list_null_value_logic_with_helper_field"
	resourceName := "test"
	ref := fmt.Sprintf("%s.%s", resourceType, resourceName)

	configWithNull := fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
	provider = "%[3]s"
	env_name = "%[4]s"
}
`, resourceType, resourceName, SdkV2FunctionalTestsProviderName, envName)

	configWithEmpty := fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
	provider = "%[3]s"
	env_name = "%[4]s"
	nullable_list = []
}
`, resourceType, resourceName, SdkV2FunctionalTestsProviderName, envName)

	configWithItems := func(items ...string) string {
		quoted := make([]string, len(items))
		for i, item := range items {
			quoted[i] = fmt.Sprintf("%q", item)
		}
		return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
	provider = "%[3]s"
	env_name = "%[4]s"
	nullable_list = [%[5]s]
}
`, resourceType, resourceName, SdkV2FunctionalTestsProviderName, envName, strings.Join(quoted, ", "))
	}

	t.Run("transition_null_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
			},
		})
	})

	t.Run("transition_empty_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
			},
		})
	})

	t.Run("transition_empty_to_filled", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithItems("x"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "true"),
						resource.TestCheckResourceAttr(ref, "get_result", "[x]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[x]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "x"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
			},
		})
	})

	t.Run("transition_null_to_filled", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithItems("a"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "a"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
			},
		})
	})

	t.Run("transition_filled_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
			},
		})
	})

	t.Run("transition_filled_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
				{
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
			},
		})
	})

	t.Run("external_change_removes_items_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "__NULL__") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithItems("a"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "a"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
			},
		})
	})

	t.Run("external_change_removes_items_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "__EMPTY__") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithItems("a"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "a"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "items"),
					),
				},
			},
		})
	})

	t.Run("external_change_adds_items_when_config_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "x,y") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
			},
		})
	})

	t.Run("external_change_adds_items_when_config_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "x,y") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
			},
		})
	})

	t.Run("external_change_nullifies_when_config_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "__NULL__") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "empty"),
					),
				},
			},
		})
	})

	t.Run("external_change_empties_when_config_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
				{
					PreConfig: func() { t.Setenv(envName, "__EMPTY__") },
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						},
					},
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
						resource.TestCheckResourceAttr(ref, "nullable_list_presence", "null"),
					),
				},
			},
		})
	})
}
