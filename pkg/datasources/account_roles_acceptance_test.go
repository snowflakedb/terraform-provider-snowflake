//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountRoles_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	accountRoleNamePrefix := random.AlphaN(10)
	accountRoleName1 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix + "1")
	accountRoleName2 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix + "2")
	accountRoleName3 := acc.TestClient().Ids.Alpha()
	dbRoleName := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix + "db")
	comment := random.Comment()

	// Proof that database role with the same prefix is not in the output of SHOW ROLES.
	dbRole, dbRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, dbRoleName)
	t.Cleanup(dbRoleCleanup)

	likeVariables := config.Variables{
		"account_role_name_1": config.StringVariable(accountRoleName1),
		"account_role_name_2": config.StringVariable(accountRoleName2),
		"account_role_name_3": config.StringVariable(accountRoleName3),
		"comment":             config.StringVariable(comment),
		"like":                config.StringVariable(accountRoleNamePrefix + "%"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: likeVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_account_roles.test", "account_roles.#", "2"),
					accountRolesDataSourceContainsRole(accountRoleName1, comment),
					accountRolesDataSourceContainsRole(accountRoleName2, comment),
					accountRolesDataSourceDoesNotContainRole(accountRoleName3, comment),
					accountRolesDataSourceDoesNotContainRole(dbRole.ID().FullyQualifiedName(), comment),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: config.Variables{},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.snowflake_account_roles.test", "account_roles.#", func(value string) error {
						numberOfRoles, err := strconv.ParseInt(value, 10, 8)
						if err != nil {
							return err
						}

						if numberOfRoles == 0 {
							return fmt.Errorf("expected roles to be non-empty")
						}

						return nil
					}),
				),
			},
		},
	})
}

func accountRolesDataSourceDoesNotContainRole(name string, comment string) func(s *terraform.State) error {
	return func(state *terraform.State) error {
		err := accountRolesDataSourceContainsRole(name, comment)(state)
		if err != nil && err.Error() == fmt.Sprintf("role %s not found", name) {
			return nil
		}
		return fmt.Errorf("expected %s not to be present", name)
	}
}

func accountRolesDataSourceContainsRole(name string, comment string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_account_roles" {
				continue
			}

			iter, err := strconv.ParseInt(rs.Primary.Attributes["account_roles.#"], 10, 32)
			if err != nil {
				return err
			}

			for i := 0; i < int(iter); i++ {
				if rs.Primary.Attributes[fmt.Sprintf("account_roles.%d.show_output.0.name", i)] == name {
					actualComment := rs.Primary.Attributes[fmt.Sprintf("account_roles.%d.show_output.0.comment", i)]
					if actualComment != comment {
						return fmt.Errorf("expected comment: %s, but got: %s", comment, actualComment)
					}

					return nil
				}
			}
		}

		return fmt.Errorf("role %s not found", name)
	}
}
