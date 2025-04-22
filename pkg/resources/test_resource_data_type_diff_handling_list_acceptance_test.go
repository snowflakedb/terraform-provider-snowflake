//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-2054208]: merge setups/test cases with TestAcc_TestResource_DataTypeDiffHandling during the package cleanup.
func TestAcc_TestResource_DataTypeDiffHandlingList(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	envName := fmt.Sprintf("%s_%s", testenvs.TestResourceDataTypeDiffHandlingEnv, strings.ToUpper(random.AlphaN(10)))
	resourceType := "snowflake_test_resource_data_type_diff_handling_list"
	resourceName := "test"
	resourceReference := fmt.Sprintf("%s.%s", resourceType, resourceName)
	listPropertyName := "nesting_list"
	propertyName := "nested_datatype"
	nestedPropertyAddress := fmt.Sprintf("%s.0.%s", listPropertyName, propertyName)

	testConfig := func(configValue string) string {
		return fmt.Sprintf(`
resource "%[3]s" "%[4]s" {
	env_name = "%[2]s"
	%[5]s {
		%[6]s = "%[1]s"
	}
}
`, configValue, envName, resourceType, resourceName, listPropertyName, propertyName)
	}

	type DataTypeDiffHandlingTestCase struct {
		ConfigValue    string
		NewConfigValue string
		ExternalValue  string
		ExpectChanges  bool
	}

	changeInConfig := func(configValue string, newConfigValue string, expectChanges bool) DataTypeDiffHandlingTestCase {
		return DataTypeDiffHandlingTestCase{
			ConfigValue:    configValue,
			NewConfigValue: newConfigValue,
			ExpectChanges:  expectChanges,
		}
	}

	externalChange := func(configValue string, externalValue string, expectChanges bool) DataTypeDiffHandlingTestCase {
		return DataTypeDiffHandlingTestCase{
			ConfigValue:   configValue,
			ExternalValue: externalValue,
			ExpectChanges: expectChanges,
		}
	}

	testCases := []DataTypeDiffHandlingTestCase{
		// different data type
		changeInConfig("NUMBER(20, 4)", "VARCHAR(20)", true),
		changeInConfig("NUMBER(20, 4)", "VARCHAR", true),
		changeInConfig("NUMBER(20)", "VARCHAR(20)", true),
		changeInConfig("NUMBER", "VARCHAR(20)", true),
		changeInConfig("NUMBER", "VARCHAR", true),

		// same data type - no attributes before
		changeInConfig("NUMBER", "NUMBER", false),
		changeInConfig("NUMBER", "NUMBER(20)", true),
		changeInConfig("NUMBER", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER", "NUMBER(38)", false),
		changeInConfig("NUMBER", "NUMBER(38, 0)", false),

		// same data type - one attribute before
		changeInConfig("NUMBER(20)", "NUMBER(20)", false),
		changeInConfig("NUMBER(20)", "NUMBER", true),
		changeInConfig("NUMBER(20)", "NUMBER(21)", true),
		changeInConfig("NUMBER(20)", "NUMBER(20, 0)", false),
		changeInConfig("NUMBER(20)", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER(20)", "NUMBER(21, 4)", true),

		// same data type - two attributes before
		changeInConfig("NUMBER(20, 3)", "NUMBER(20, 3)", false),
		changeInConfig("NUMBER(20, 3)", "NUMBER", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(20)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21, 3)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21, 4)", true),

		// same data type - one attribute but default before
		changeInConfig("NUMBER(38)", "NUMBER(38)", false),
		changeInConfig("NUMBER(38)", "NUMBER", false),
		changeInConfig("NUMBER(38)", "NUMBER(20)", true),
		changeInConfig("NUMBER(38)", "NUMBER(20, 3)", true),
		changeInConfig("NUMBER(38)", "NUMBER(38, 2)", true),
		changeInConfig("NUMBER(38)", "NUMBER(38, 0)", false),

		// same data type - two attributes but default before
		changeInConfig("NUMBER(38, 0)", "NUMBER(38, 0)", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER(38)", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER(20)", true),
		changeInConfig("NUMBER(38, 0)", "NUMBER(20, 3)", true),
		changeInConfig("NUMBER(38, 0)", "NUMBER(38, 2)", true),

		// different data type
		externalChange("NUMBER(20, 4)", "VARCHAR(20)", true),
		externalChange("NUMBER(20, 4)", "VARCHAR", true),
		externalChange("NUMBER(20)", "VARCHAR(20)", true),
		externalChange("NUMBER", "VARCHAR(20)", true),
		externalChange("NUMBER", "VARCHAR", true),

		// same data type - no attributes before
		externalChange("NUMBER", "NUMBER", false),
		externalChange("NUMBER", "NUMBER(20)", true),
		externalChange("NUMBER", "NUMBER(20, 4)", true),
		externalChange("NUMBER", "NUMBER(38)", false),
		externalChange("NUMBER", "NUMBER(38, 0)", false),

		// same data type - one attribute before
		externalChange("NUMBER(20)", "NUMBER(20)", false),
		externalChange("NUMBER(20)", "NUMBER", false),
		externalChange("NUMBER(20)", "NUMBER(21)", true),
		externalChange("NUMBER(20)", "NUMBER(20, 0)", false),
		externalChange("NUMBER(20)", "NUMBER(20, 4)", true),
		externalChange("NUMBER(20)", "NUMBER(21, 4)", true),

		// same data type - two attributes before
		externalChange("NUMBER(20, 3)", "NUMBER(20, 3)", false),
		externalChange("NUMBER(20, 3)", "NUMBER", false),
		externalChange("NUMBER(20, 3)", "NUMBER(20)", false),
		externalChange("NUMBER(20, 3)", "NUMBER(20, 4)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21, 3)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21, 4)", true),

		// same data type - one attribute but default before
		externalChange("NUMBER(38)", "NUMBER(38)", false),
		externalChange("NUMBER(38)", "NUMBER", false),
		externalChange("NUMBER(38)", "NUMBER(20)", true),
		externalChange("NUMBER(38)", "NUMBER(20, 3)", true),
		externalChange("NUMBER(38)", "NUMBER(38, 2)", true),
		externalChange("NUMBER(38)", "NUMBER(38, 0)", false),

		// same data type - two attributes but default before
		externalChange("NUMBER(38, 0)", "NUMBER(38, 0)", false),
		externalChange("NUMBER(38, 0)", "NUMBER", false),
		externalChange("NUMBER(38, 0)", "NUMBER(38)", false),
		externalChange("NUMBER(38, 0)", "NUMBER(20)", true),
		externalChange("NUMBER(38, 0)", "NUMBER(20, 3)", true),
		externalChange("NUMBER(38, 0)", "NUMBER(38, 2)", true),
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(fmt.Sprintf("TestAcc_TestResource_DataTypeDiffHandlingList config value: %s, new config value: %s, external value: %s, expecting changes: %t", tc.ConfigValue, tc.NewConfigValue, tc.ExternalValue, tc.ExpectChanges), func(t *testing.T) {
			configValueDataType, err := datatypes.ParseDataType(tc.ConfigValue)
			require.NoError(t, err)

			newConfigValue := tc.ConfigValue
			if tc.NewConfigValue != "" {
				newConfigValue = tc.NewConfigValue
			}

			expectedStateFirstStep := configValueDataType.ToSql()
			expectedStateSecondStep := expectedStateFirstStep
			if tc.ExpectChanges {
				dt, err := datatypes.ParseDataType(newConfigValue)
				require.NoError(t, err)
				expectedStateSecondStep = dt.ToSql()
			}

			var checks []plancheck.PlanCheck
			if tc.ExpectChanges {
				if tc.ExternalValue != "" {
					checks = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(resourceReference, nestedPropertyAddress, sdk.String(expectedStateFirstStep), sdk.String(tc.ExternalValue)),
						// TODO [SNOW-1473409]: expecting delete as currently this plan check does not offer setting multiple actions; we expect destroy and create here
						planchecks.ExpectChange(resourceReference, nestedPropertyAddress, tfjson.ActionDelete, sdk.String(tc.ExternalValue), sdk.String(expectedStateFirstStep)),
					}
				} else {
					checks = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionDestroyBeforeCreate),
						// TODO [SNOW-1473409]: expecting delete as currently this plan check does not offer setting multiple actions; we expect destroy and create here
						planchecks.ExpectChange(resourceReference, nestedPropertyAddress, tfjson.ActionDelete, sdk.String(expectedStateFirstStep), sdk.String(expectedStateSecondStep)),
					}
				}
			} else {
				checks = []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				}
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				Steps: []resource.TestStep{
					{
						// our test resource manages this env, so we remove it before the test start
						PreConfig: func() {
							t.Setenv(envName, "")
						},
						Config: testConfig(tc.ConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, nestedPropertyAddress, expectedStateFirstStep),
						),
					},
					{
						PreConfig: func() {
							if tc.ExternalValue != "" {
								t.Setenv(envName, tc.ExternalValue)
							}
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: checks,
						},
						Config: testConfig(newConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, nestedPropertyAddress, expectedStateSecondStep),
						),
					},
				},
			})
		})
	}
}
