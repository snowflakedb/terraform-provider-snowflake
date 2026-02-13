package testfunctional_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_SdkV2Functional_TestResource_ListNullValueLogic tests whether SDKv2 can distinguish
// between three states of an Optional TypeList field:
// - null (field not specified in config)
// - empty list (field = [])
// - filled list (field = ["a", "b"])
//
// It observes the results of d.Get, d.GetOk, and d.GetRawConfig for each state and transition.
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

	// Test 1: Create with null list (field not specified)
	t.Run("create_with_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ref, "get_ok_result"),
						resource.TestCheckResourceAttrSet(ref, "get_length_result"),
						resource.TestCheckResourceAttrSet(ref, "raw_config_is_null_result"),
						resource.TestCheckResourceAttrSet(ref, "raw_config_length_result"),
						// Log all observations for analysis
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})

	// Test 2: Create with empty list
	t.Run("create_with_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})

	// Test 3: Create with filled list
	t.Run("create_with_items", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithItems("a", "b"),
					Check: resource.ComposeTestCheckFunc(
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "2"),
						resource.TestCheckResourceAttr(ref, "nullable_list.0", "a"),
						resource.TestCheckResourceAttr(ref, "nullable_list.1", "b"),
					),
				},
			},
		})
	})

	// Test 4: Transition from filled to empty
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
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})

	// Test 5: Transition from filled to null
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
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})

	// Test 6: Transition from null to empty
	t.Run("null_to_empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithNull,
				},
				{
					Config: configWithEmpty,
					Check: resource.ComposeTestCheckFunc(
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})

	// Test 7: Transition from empty to filled
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
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "3"),
					),
				},
			},
		})
	})

	// Test 8: Transition from null to filled
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
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
						resource.TestCheckResourceAttr(ref, "nullable_list.#", "1"),
					),
				},
			},
		})
	})

	// Test 9: Transition from empty to null
	t.Run("empty_to_null", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerForSdkV2FunctionalTestsFactories,
			TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.RequireAbove(tfversion.Version1_5_0)},
			Steps: []resource.TestStep{
				{
					PreConfig: func() { t.Setenv(envName, "") },
					Config:    configWithEmpty,
				},
				{
					Config: configWithNull,
					Check: resource.ComposeTestCheckFunc(
						logResourceAttr(ref, "get_ok_result"),
						logResourceAttr(ref, "get_length_result"),
						logResourceAttr(ref, "raw_config_is_null_result"),
						logResourceAttr(ref, "raw_config_length_result"),
						logResourceAttr(ref, "nullable_list.#"),
					),
				},
			},
		})
	})
}

// logResourceAttr is a test check function that logs the value of a resource attribute for observation.
func logResourceAttr(resourceReference string, attrName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceReference)
		}
		value, ok := rs.Primary.Attributes[attrName]
		if !ok {
			value = "<not set>"
		}
		fmt.Printf("[OBSERVATION] %s.%s = %s\n", resourceReference, attrName, value)
		return nil
	}
}
