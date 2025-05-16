//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	m := func(minLength int, maxLength int, minUpperCaseChars int, minLowerCaseChars int, minNumericChars int, minSpecialChars int, minAgeDays int, maxAgeDays int, maxRetries int, lockoutTimeMins int, history int, comment string) map[string]config.Variable {
		return map[string]config.Variable{
			"name":                 config.StringVariable(id.Name()),
			"database":             config.StringVariable(id.DatabaseName()),
			"schema":               config.StringVariable(id.SchemaName()),
			"min_length":           config.IntegerVariable(minLength),
			"max_length":           config.IntegerVariable(maxLength),
			"min_upper_case_chars": config.IntegerVariable(minUpperCaseChars),
			"min_lower_case_chars": config.IntegerVariable(minLowerCaseChars),
			"min_numeric_chars":    config.IntegerVariable(minNumericChars),
			"min_special_chars":    config.IntegerVariable(minSpecialChars),
			"min_age_days":         config.IntegerVariable(minAgeDays),
			"max_age_days":         config.IntegerVariable(maxAgeDays),
			"max_retries":          config.IntegerVariable(maxRetries),
			"lockout_time_mins":    config.IntegerVariable(lockoutTimeMins),
			"history":              config.IntegerVariable(history),
			"comment":              config.StringVariable(comment),
		}
	}
	variables1 := m(10, 30, 2, 3, 4, 5, 6, 7, 8, 9, 10, comment)
	variables2 := m(20, 50, 1, 2, 3, 4, 5, 6, 7, 8, 9, comment)
	variables3 := m(20, 50, 1, 2, 3, 4, 5, 6, 7, 8, 9, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/basic"),
				ConfigVariables: variables1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "30"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_upper_case_chars", "2"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_lower_case_chars", "3"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_numeric_chars", "4"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_special_chars", "5"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_age_days", "6"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "7"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_retries", "8"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "lockout_time_mins", "9"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "history", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", comment),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/basic"),
				ConfigVariables: variables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "50"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_upper_case_chars", "1"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_lower_case_chars", "2"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_numeric_chars", "3"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_special_chars", "4"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_age_days", "5"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "6"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_retries", "7"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "lockout_time_mins", "8"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "history", "9"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", comment),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/basic"),
				ConfigVariables: variables3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", ""),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_PasswordPolicy/basic"),
				ConfigVariables:   variables3,
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_PasswordPolicyMaxAgeDays(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	oldId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	m := func(maxAgeDays int) map[string]config.Variable {
		return map[string]config.Variable{
			"name":         config.StringVariable(oldId.Name()),
			"database":     config.StringVariable(oldId.DatabaseName()),
			"schema":       config.StringVariable(oldId.SchemaName()),
			"max_age_days": config.IntegerVariable(maxAgeDays),
		}
	}

	configValueWithNewName := m(10)
	configValueWithNewName["name"] = config.StringVariable(newId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			// Creation sets zero properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/withMaxAgeDays"),
				ConfigVariables: m(0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "fully_qualified_name", oldId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/withMaxAgeDays"),
				ConfigVariables: m(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "10"),
				),
			},
			// Update sets zero properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/withMaxAgeDays"),
				ConfigVariables: m(0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			// Rename + Unsets properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newId.Name()),
					"database": config.StringVariable(newId.DatabaseName()),
					"schema":   config.StringVariable(newId.SchemaName()),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_password_policy.pa", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "90"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy/noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(oldId.Name()),
					"database": config.StringVariable(oldId.DatabaseName()),
					"schema":   config.StringVariable(oldId.SchemaName()),
				},
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_PasswordPolicy_migrateFromVersion_0_94_1(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resourceName := "snowflake_password_policy.pa"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            passwordPolicyBasicConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   passwordPolicyBasicConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}

func passwordPolicyBasicConfig(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_password_policy" "pa" {
  database = "%[1]s"
  schema   = "%[2]s"
  name     = "%[3]s"
}`, id.DatabaseName(), id.SchemaName(), id.Name())
}
