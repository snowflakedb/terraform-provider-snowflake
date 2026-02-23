package testfunctional_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_SdkV2Functional_TestResource_ListNullValueLogic tests whether SDKv2 can distinguish
// between three states of an Optional TypeList field:
// - null (field = null or not specified in config)
// - empty list (field = [])
// - filled list (field = ["a", "b"])
//
// Expected SDKv2 observation matrix:
//
//	| State  | d.GetOk ok | d.Get value | RawConfig value |
//	|--------|------------|-------------|-----------------|
//	| null   | false      | []          | null            |
//	| empty  | false      | []          | []              |
//	| filled | true       | [a b]       | [a b]           |
//
// d.GetOk and d.Get cannot distinguish null from empty.
// Only d.GetRawConfig() reliably differentiates null vs empty.
//
// SDKv2 also cannot detect transitions between null and empty because both
// produce nullable_list.# = 0 in state, so no plan diff is generated and
// Update is never called for null <-> empty changes.
func TestAcc_SdkV2Functional_TestResource_ListNullValueLogic(t *testing.T) {
	envName := fmt.Sprintf("%s_%s", testenvs.TestResourceDataTypeDiffHandlingEnv, strings.ToUpper(random.AlphaN(10)))
	resourceType := "snowflake_test_resource_list_null_value_logic"
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

	t.Run("create_with_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})

	t.Run("create_with_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						// d.GetOk returns false for empty list (zero value)
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						// d.Get returns the same empty slice as null — indistinguishable
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						// RawConfig distinguishes empty from null
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})

	t.Run("create_with_items", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "true"),
						resource.TestCheckResourceAttr(ref, "get_result", "[a b]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[a b]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "a"),
						resource.TestCheckResourceAttr(ref, "nullable_list.1", "b"),
					),
				},
			},
		})
	})

	// Transitions where Update IS triggered (state diff: #=N → #=M, N≠M)

	t.Run("filled_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
					),
				},
				{
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})

	t.Run("filled_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
					),
				},
				{
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})

	t.Run("empty_to_filled", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
				},
				{
					Config: configWithItems("x", "y", "z"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "true"),
						resource.TestCheckResourceAttr(ref, "get_result", "[x y z]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[x y z]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "3"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "x"),
						resource.TestCheckResourceAttr(ref, "nullable_list.1", "y"),
						resource.TestCheckResourceAttr(ref, "nullable_list.2", "z"),
					),
				},
			},
		})
	})

	t.Run("null_to_filled", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
				},
				{
					Config: configWithItems("x"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "get_ok_result", "true"),
						resource.TestCheckResourceAttr(ref, "get_result", "[x]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[x]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "x"),
					),
				},
			},
		})
	})

	// Transitions between null and empty: SDKv2 represents both as nullable_list.# = 0
	// in state, so no plan diff is generated and Update is never called. The observation
	// fields retain their values from the Create in step 1.

	t.Run("null_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
					),
				},
				{
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						// Observations are stale from Create(null) because SDKv2 detected
						// no diff (both null and empty produce #=0), so Update was not called.
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "null"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})

	t.Run("empty_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
					),
				},
				{
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						// Observations are stale from Create(empty) because SDKv2 detected
						// no diff (both null and empty produce #=0), so Update was not called.
						resource.TestCheckResourceAttr(ref, "get_ok_result", "false"),
						resource.TestCheckResourceAttr(ref, "get_result", "[]"),
						resource.TestCheckResourceAttr(ref, "raw_config_result", "[]"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "0"),
					),
				},
			},
		})
	})
}
