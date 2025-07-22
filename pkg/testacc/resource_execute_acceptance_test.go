//go:build !account_level_tests

package testacc

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcc_Execute_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	secondId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	nameLowerCase := strings.ToLower(secondId.Name())
	secondIdLowerCased := sdk.NewAccountObjectIdentifier(nameLowerCase)
	nameLowerCaseEscaped := fmt.Sprintf(`"%s"`, nameLowerCase)
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }

	resourceName := "snowflake_execute.test"
	createConfigVariables := func(id string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(createDatabaseStatement(id)),
			"revert":  config.StringVariable(dropDatabaseStatement(id)),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckNoResourceAttr(resourceName, "query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, secondIdLowerCased, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(nameLowerCaseEscaped),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(nameLowerCaseEscaped)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(nameLowerCaseEscaped)),
					resource.TestCheckNoResourceAttr(resourceName, "query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, secondIdLowerCased, true),
				),
			},
		},
	})
}

func TestAcc_Execute_withRead(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }

	resourceName := "snowflake_execute.test"
	createConfigVariables := func(id string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(createDatabaseStatement(id)),
			"revert":  config.StringVariable(dropDatabaseStatement(id)),
			"query":   config.StringVariable(showDatabaseStatement(id)),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withRead"),
				ConfigVariables: createConfigVariables(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.created_on"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_Execute_readRemoved(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }
	resourceName := "snowflake_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(showDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(""),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", ""),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
				),
			},
		},
	})
}

func TestAcc_Execute_badQuery(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }
	resourceName := "snowflake_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable("bad query"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "query", "bad query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(showDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", name),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_Execute_invalidExecuteStatement(t *testing.T) {
	invalidCreateStatement := "create database"
	invalidDropStatement := "drop database"

	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(invalidCreateStatement),
			"revert":  config.StringVariable(invalidDropStatement),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
		},
	})
}

func TestAcc_Execute_invalidRevertStatement(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	updatedId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	updatedName := updatedId.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	invalidDropStatement := "drop database"

	resourceName := "snowflake_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckDatabaseExistence(t, id, false)(state)
			if err != nil {
				return err
			}
			err = testAccCheckDatabaseExistence(t, updatedId, false)(state)
			if err != nil {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(invalidDropStatement),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", invalidDropStatement),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedName)),
					"revert":  config.StringVariable(invalidDropStatement),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					testAccCheckDatabaseExistence(t, updatedId, false),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedName)),
					"revert":  config.StringVariable(dropDatabaseStatement(updatedName)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(updatedName)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(updatedName)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, false),
					testAccCheckDatabaseExistence(t, updatedId, true),
				),
			},
		},
	})
}

func TestAcc_Execute_revertUpdated(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	execute := fmt.Sprintf("create database %s", name)
	revert := fmt.Sprintf("drop database %s", name)
	notMatchingRevert := "select 1"
	var savedId string

	resourceName := "snowflake_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, notMatchingRevert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", notMatchingRevert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						savedId = value
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if savedId != value {
							return errors.New("different id after revert update")
						}
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_Execute_executeUpdated(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	name := id.Name()
	execute := fmt.Sprintf("create database %s", name)
	revert := fmt.Sprintf("drop database %s", name)

	newId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("EXECUTE_TEST_DATABASE_")
	newName := newId.Name()
	newExecute := fmt.Sprintf("create database %s", newName)
	newRevert := fmt.Sprintf("drop database %s", newName)

	var savedId string

	resourceName := "snowflake_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckDatabaseExistence(t, id, false)(state)
			if err != nil {
				return err
			}
			err = testAccCheckDatabaseExistence(t, newId, false)(state)
			if err != nil {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						savedId = value
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(newExecute, newRevert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", newExecute),
					resource.TestCheckResourceAttr(resourceName, "revert", newRevert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if savedId == value {
							return errors.New("same id after execute update")
						}
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, false),
					testAccCheckDatabaseExistence(t, newId, true),
				),
			},
		},
	})
}

func TestAcc_Execute_grants(t *testing.T) {
	client := testClient()

	database, databaseCleanup := client.Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	role, roleCleanup := client.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	privilege := sdk.AccountObjectPrivilegeCreateSchema
	execute := fmt.Sprintf("GRANT %s ON DATABASE %s TO ROLE %s", privilege, database.ID().FullyQualifiedName(), role.ID().FullyQualifiedName())
	revert := fmt.Sprintf("REVOKE %s ON DATABASE %s FROM ROLE %s", privilege, database.ID().FullyQualifiedName(), role.ID().FullyQualifiedName())

	resourceName := "snowflake_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := verifyGrantExists(t, role.ID(), privilege, false)(state)
			return err
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					verifyGrantExists(t, role.ID(), privilege, true),
				),
			},
		},
	})
}

// TestAcc_Execute_grantsComplex test fails with:
//
//	testing_new_config.go:156: unexpected index type (string) for "snowflake_execute.test[\"0\"]", for_each is not supported
//	testing_new.go:68: unexpected index type (string) for "snowflake_execute.test[\"0\"]", for_each is not supported
//
// Quick search unveiled this issue: https://github.com/hashicorp/terraform-plugin-sdk/issues/536.
//
// It also seems that it is working correctly underneath:
//
//	2023/11/26 17:16:03 [DEBUG] SQL "GRANT CREATE SCHEMA,MODIFY ON DATABASE EXECUTE_TEST_DATABASE_4397 TO ROLE EXECUTE_TEST_ROLE_1145" applied successfully
//	2023/11/26 17:16:03 [DEBUG] SQL "GRANT MODIFY,USAGE ON DATABASE EXECUTE_TEST_DATABASE_3740 TO ROLE EXECUTE_TEST_ROLE_3008" applied successfully
func TestAcc_Execute_grantsComplex(t *testing.T) {
	t.Skip("Skipping TestAcc_Execute_grantsComplex because of https://github.com/hashicorp/terraform-plugin-sdk/issues/536 issue")

	client := testClient()

	database1, database1Cleanup := client.Database.CreateDatabase(t)
	t.Cleanup(database1Cleanup)

	database2, database2Cleanup := client.Database.CreateDatabase(t)
	t.Cleanup(database2Cleanup)

	role1, role1Cleanup := client.Role.CreateRole(t)
	t.Cleanup(role1Cleanup)

	role2, role2Cleanup := client.Role.CreateRole(t)
	t.Cleanup(role2Cleanup)

	dbId1 := database1.ID()
	dbId2 := database2.ID()
	roleId1 := role1.ID()
	roleId2 := role2.ID()
	privilege1 := sdk.AccountObjectPrivilegeCreateSchema
	privilege2 := sdk.AccountObjectPrivilegeModify
	privilege3 := sdk.AccountObjectPrivilegeUsage

	// resourceName1 := "snowflake_execute.test.0"
	// resourceName2 := "snowflake_execute.test.1"
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database_grants": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId1.Name()),
				"role_id":       config.StringVariable(roleId1.Name()),
				"privileges":    config.ListVariable(config.StringVariable(privilege1.String()), config.StringVariable(privilege2.String())),
			}), config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId2.Name()),
				"role_id":       config.StringVariable(roleId2.Name()),
				"privileges":    config.ListVariable(config.StringVariable(privilege2.String()), config.StringVariable(privilege3.String())),
			})),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := verifyGrantExists(t, roleId1, privilege1, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId1, privilege2, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId1, privilege3, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege1, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege2, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege3, false)(state)
			if err != nil {
				return err
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet(resourceName1, "id"),
					// resource.TestCheckResourceAttrSet(resourceName2, "id"),
					verifyGrantExists(t, roleId1, privilege1, true),
					verifyGrantExists(t, roleId1, privilege2, true),
					verifyGrantExists(t, roleId1, privilege3, false),
					verifyGrantExists(t, roleId2, privilege1, false),
					verifyGrantExists(t, roleId2, privilege2, true),
					verifyGrantExists(t, roleId2, privilege3, true),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2491
func TestAcc_Execute_queryResultsBug(t *testing.T) {
	resourceName := "snowflake_execute.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: executeConfig(108),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", "SELECT 108"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.108"),
				),
			},
			{
				Config: executeConfig(96),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", "SELECT 96"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.96"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func executeConfig(queryNumber int) string {
	return fmt.Sprintf(`
resource "snowflake_execute" "test" {
  execute = "SELECT 18"
  revert  = "SELECT 36"
  query  = "SELECT %d"
}

output "query_results_output" {
  value = snowflake_execute.test.query_results
}
`, queryNumber)
}

func TestAcc_Execute_QueryResultsRecomputedWithoutQueryChanges(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resourceName := "snowflake_execute.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: executeConfigCreateDatabase(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query_results.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.comment", ""),
				),
			},
			{
				PreConfig: func() {
					testClient().Database.Alter(t, id, &sdk.AlterDatabaseOptions{
						Set: &sdk.DatabaseSet{
							Comment: sdk.String("some comment"),
						},
					})
				},
				Config: executeConfigCreateDatabase(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query_results.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.comment", "some comment"),
				),
			},
		},
	})
}

func executeConfigCreateDatabase(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE \"%[1]s\""
  revert  = "DROP DATABASE \"%[1]s\""
  query   = "SHOW DATABASES LIKE '%[1]s'"
}
`, id.Name())
}

func verifyGrantExists(t *testing.T, roleId sdk.AccountObjectIdentifier, privilege sdk.AccountObjectPrivilege, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		grants, err := testClient().Grant.ShowGrantsToAccountRole(t, roleId)
		if err != nil {
			return err
		}

		if shouldExist {
			require.Len(t, grants, 1)
			assert.Equal(t, privilege.String(), grants[0].Privilege)
			assert.Equal(t, sdk.ObjectTypeDatabase, grants[0].GrantedOn)
			assert.Equal(t, sdk.ObjectTypeRole, grants[0].GrantedTo)
			assert.Equal(t, roleId.FullyQualifiedName(), grants[0].GranteeName.FullyQualifiedName())
		} else {
			require.Empty(t, grants)
		}

		// it does not matter what we return, because we have assertions above
		return nil
	}
}

func TestAcc_Execute_ImportWithRandomId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup := testClient().Database.CreateDatabaseWithIdentifier(t, id)
					t.Cleanup(databaseCleanup)
				},
				Config:                  executeConfigCreateDatabase(id),
				ResourceName:            "snowflake_execute.test",
				ImportState:             true,
				ImportStatePersist:      true,
				ImportStateId:           "random_id",
				ImportStateVerifyIgnore: []string{"query_results"},
			},
			// filling the empty state fields (execute changed from empty)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_execute.test", plancheck.ResourceActionUpdate),
					},
				},
				Config: executeConfigCreateDatabase(id),
			},
			// change the id in every query to see if:
			// 1. execute will trigger force new behavior
			// 2. an old database is used in delete (it is)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_execute.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						resources.PlanCheckFunc(func(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
							_, err := testClient().Database.Show(t, id)
							if err == nil {
								resp.Error = fmt.Errorf("database %s still exist", id.FullyQualifiedName())
								t.Cleanup(testClient().Database.DropDatabaseFunc(t, id))
							}
						}),
					},
				},
				Config: executeConfigCreateDatabase(newId),
			},
		},
	})
}

// TODO [SNOW-1348121]: Move this to the file with check_destroy functions.
func testAccCheckDatabaseExistence(t *testing.T, id sdk.AccountObjectIdentifier, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		_, err := testClient().Database.Show(t, id)
		if shouldExist {
			if err != nil {
				return fmt.Errorf("error while retrieving database %s, err = %w", id, err)
			}
		} else {
			if err == nil {
				return fmt.Errorf("database %v still exists", id)
			}
		}
		return nil
	}
}

// Result of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3334.
func TestAcc_Execute_gh3334_allTimeouts(t *testing.T) {
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"execute":        config.StringVariable("CALL SYSTEM$WAIT(5, 'SECONDS');"),
			"revert":         config.StringVariable("select 2"),
			"query":          config.StringVariable("select 3"),
			"create_timeout": config.StringVariable("1m"),
			"read_timeout":   config.StringVariable("31m"),
			"update_timeout": config.StringVariable("32m"),
			"delete_timeout": config.StringVariable("33m"),
		}
	}

	resourceName := "snowflake_execute.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withTimeouts"),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "1m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "31m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "32m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.delete", "33m"),
				),
			},
		},
	})
}

// Result of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3334.
func TestAcc_Execute_gh3334_longRunningCreate(t *testing.T) {
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"execute":        config.StringVariable("CALL SYSTEM$WAIT(15, 'SECONDS');"),
			"revert":         config.StringVariable("select 2"),
			"query":          config.StringVariable("select 3"),
			"create_timeout": config.StringVariable("5s"),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Execute_withTimeouts"),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				ExpectError: regexp.MustCompile("Error: context deadline exceeded"),
			},
		},
	})
}
