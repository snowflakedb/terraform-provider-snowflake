//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantDatabaseRole_databaseRole(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleName := databaseRoleId.Name()
	parentDatabaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	parentDatabaseRoleName := parentDatabaseRoleId.Name()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(databaseRoleId.DatabaseName()),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}

	resourceName := "snowflake_grant_database_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%[1]v"."%[2]v"|DATABASE ROLE|"%[1]v"."%[3]v"`, acc.TestDatabaseName, databaseRoleName, parentDatabaseRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_databaseRoleMixedQuoting(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleName := databaseRoleId.Name()
	parentDatabaseRoleId := acc.TestClient().Ids.NewDatabaseObjectIdentifier(strings.ToUpper(acc.TestClient().Ids.Alpha()))
	parentDatabaseRoleName := parentDatabaseRoleId.Name()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}

	resourceName := "snowflake_grant_database_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%[1]v"."%[2]v"|DATABASE ROLE|"%[1]v"."%[3]v"`, acc.TestDatabaseName, databaseRoleName, parentDatabaseRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_issue2402(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	databaseRoleName := databaseRoleId.Name()
	parentDatabaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	parentDatabaseRoleName := parentDatabaseRoleId.Name()
	databaseName := database.ID().Name()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(databaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}

	resourceName := "snowflake_grant_database_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantDatabaseRole/issue2402"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%[1]v"."%[2]v"|DATABASE ROLE|"%[1]v"."%[3]v"`, databaseName, databaseRoleName, parentDatabaseRoleName)),
				),
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_accountRole(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleName := databaseRoleId.Name()
	parentRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	parentRoleName := parentRoleId.Name()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":           config.StringVariable(acc.TestDatabaseName),
			"database_role_name": config.StringVariable(databaseRoleName),
			"parent_role_name":   config.StringVariable(parentRoleName),
		}
	}

	resourceName := "snowflake_grant_database_role.g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", fmt.Sprintf("%v", parentRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ROLE|"%v"`, acc.TestDatabaseName, databaseRoleName, parentRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2410 is fixed
func TestAcc_GrantDatabaseRole_share(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	shareId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func() config.Variables {
		return config.Variables{
			"database":           config.StringVariable(database.ID().Name()),
			"database_role_name": config.StringVariable(databaseRoleId.Name()),
			"share_name":         config.StringVariable(shareId.Name()),
		}
	}

	resourceName := "snowflake_grant_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables: configVariables(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "share_name", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`%v|%v|%v`, databaseRoleId.FullyQualifiedName(), "SHARE", shareId.FullyQualifiedName())),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables:   configVariables(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_shareWithDots(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	shareId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")

	configVariables := func() config.Variables {
		return config.Variables{
			"database":           config.StringVariable(database.ID().Name()),
			"database_role_name": config.StringVariable(databaseRoleId.Name()),
			"share_name":         config.StringVariable(shareId.Name()),
		}
	}

	resourceName := "snowflake_grant_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables: configVariables(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "share_name", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`%v|%v|%v`, databaseRoleId.FullyQualifiedName(), "SHARE", shareId.FullyQualifiedName())),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables:   configVariables(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	parentRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantDatabaseRoleBasicConfigQuoted(databaseRoleId, parentRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, databaseRoleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantDatabaseRoleBasicConfigQuoted(databaseRoleId, parentRoleId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, databaseRoleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantDatabaseRoleBasicConfigQuoted(databaseRoleId sdk.DatabaseObjectIdentifier, parentRoleId sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database_role" "role" {
  database = "%[1]s"
  name = "%[2]s"
}

resource "snowflake_database_role" "parent_role" {
  database = "%[1]s"
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name        = "\"%[1]s\".\"${snowflake_database_role.role.name}\""
  parent_database_role_name = "\"%[1]s\".\"${snowflake_database_role.parent_role.name}\""
}
`, databaseRoleId.DatabaseName(), databaseRoleId.Name(), parentRoleId.Name())
}

func TestAcc_GrantDatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	parentRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantDatabaseRoleBasicConfigUnquoted(databaseRoleId, parentRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "database_role_name", fmt.Sprintf("%s.%s", databaseRoleId.DatabaseName(), databaseRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "parent_database_role_name", fmt.Sprintf("%s.%s", parentRoleId.DatabaseName(), parentRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, databaseRoleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantDatabaseRoleBasicConfigUnquoted(databaseRoleId, parentRoleId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "database_role_name", fmt.Sprintf("%s.%s", databaseRoleId.DatabaseName(), databaseRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "parent_database_role_name", fmt.Sprintf("%s.%s", parentRoleId.DatabaseName(), parentRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, databaseRoleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantDatabaseRoleBasicConfigUnquoted(databaseRoleId sdk.DatabaseObjectIdentifier, parentRoleId sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database_role" "role" {
  database = "%[1]s"
  name = "%[2]s"
}

resource "snowflake_database_role" "parent_role" {
  database = "%[1]s"
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name        = "%[1]s.${snowflake_database_role.role.name}"
  parent_database_role_name = "%[1]s.${snowflake_database_role.parent_role.name}"
}
`, databaseRoleId.DatabaseName(), databaseRoleId.Name(), parentRoleId.Name())
}

func TestAcc_GrantDatabaseRole_Issue_3629(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	parentRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:            grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
				ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0"),
			},
			{
				PreConfig: func() {
					acc.TestClient().BcrBundles.EnableBcrBundle(t, "2025_02")
					acc.TestClient().Execute.SQL(t, fmt.Sprintf("GRANT DATABASE ROLE %s TO USER %s", databaseRoleId.FullyQualifiedName(), user.ID().FullyQualifiedName()))
				},
				ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0"),
				Config:            grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
				ExpectError:       regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Config: grantDatabaseRoleIssue3629Config(databaseRoleId, parentRoleId),
			},
		},
	})
}

func grantDatabaseRoleIssue3629Config(databaseRoleId sdk.DatabaseObjectIdentifier, accountRoleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database_role" "test" {
    database = "%[1]s"
    name = "%[2]s"
}

resource "snowflake_account_role" "test" {
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  parent_role_name = snowflake_account_role.test.fully_qualified_name
}
`, databaseRoleId.DatabaseName(), databaseRoleId.Name(), accountRoleId.Name())
}
