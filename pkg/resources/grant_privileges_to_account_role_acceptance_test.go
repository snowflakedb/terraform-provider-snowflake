//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToAccountRole_OnAccount(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateRole)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE,CREATE ROLE|OnAccount", roleFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_gh3153(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeManageShareTarget)),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_gh3153"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeManageShareTarget)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|%s|OnAccount", roleFullyQualifiedName, sdk.GlobalPrivilegeManageShareTarget)),
				),
			},
		},
	})
}

// Proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/3507 is fixed.
func TestAcc_GrantPrivilegesToAccountRole_OnAccount_gh3507(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name":         config.StringVariable(roleFullyQualifiedName),
		"always_apply": config.BoolVariable(false),
	}
	configVariablesWithAlwaysApply := config.Variables{
		"name":         config.StringVariable(roleFullyQualifiedName),
		"always_apply": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetLegacyConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("1.0.5"),
				Config:            grantAllPrivilegesToAccountRoleBasicConfig(role.ID()),
				ExpectError:       regexp.MustCompile(`Error: 003011 \(42501\): Grant partially executed`),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
				ConfigVariables:          configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckNoResourceAttr(resourceName, "privileges.#"),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccount", roleFullyQualifiedName)),
					queriedAccountRolePrivilegesContainAtLeast(t, role.ID(), string(sdk.GlobalPrivilegeCreateDatabase)),
					queriedAccountRolePrivilegesDoNotContain(t, role.ID(), string(sdk.GlobalPrivilegeManageListingAutoFulfillment), string(sdk.GlobalPrivilegeManageOrganizationSupportCases), string(sdk.GlobalPrivilegeManagePolarisConnections)),
				),
				// Due to limitations in the plugin SDK, returned warnings can not be asserted (see https://github.com/hashicorp/terraform-plugin-testing/issues/69).
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
				ConfigVariables:          configVariablesWithAlwaysApply,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckNoResourceAttr(resourceName, "privileges.#"),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccount", roleFullyQualifiedName)),
					queriedAccountRolePrivilegesContainAtLeast(t, role.ID(), string(sdk.GlobalPrivilegeCreateDatabase)),
					queriedAccountRolePrivilegesDoNotContain(t, role.ID(), string(sdk.GlobalPrivilegeManageListingAutoFulfillment), string(sdk.GlobalPrivilegeManageOrganizationSupportCases), string(sdk.GlobalPrivilegeManagePolarisConnections)),
				),
				// We expect the plan to be non-empty because in this step we set `always_apply`, causing the permadiff.
				ExpectNonEmptyPlan: true,
				// Due to limitations in the plugin SDK, returned warnings can not be asserted (see https://github.com/hashicorp/terraform-plugin-testing/issues/69).
			},
		},
	})
}

func grantAllPrivilegesToAccountRoleBasicConfig(roleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
	account_role_name = "%s"
	all_privileges    = true
	on_account        = true
}
`, roleId.Name())
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_ErrorOnPrivilegesNotGranted(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantablePrivilege1 := config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase))
	grantablePrivilege2 := config.StringVariable(string(sdk.GlobalPrivilegeApplyAggregationPolicy))
	nonGrantablePrivilege := config.StringVariable("MANAGE LISTING AUTO FULFILLMENT")
	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariablesWithGrantablePrivileges := config.Variables{
		"name":              config.StringVariable(roleFullyQualifiedName),
		"privileges":        config.ListVariable(grantablePrivilege1),
		"with_grant_option": config.BoolVariable(true),
	}
	configVariablesWithNonGrantablePrivileges := config.Variables{
		"name":              config.StringVariable(roleFullyQualifiedName),
		"privileges":        config.ListVariable(grantablePrivilege1, grantablePrivilege2, nonGrantablePrivilege),
		"with_grant_option": config.BoolVariable(true),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithNonGrantablePrivileges,
				ExpectError:     regexp.MustCompile("grant partially executed"),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithGrantablePrivileges,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithNonGrantablePrivileges,
				ExpectError:     regexp.MustCompile("grant partially executed"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_ChangeListOfPrivilegesToAllPrivileges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantablePrivilege1 := config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase))
	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariablesWithGrantablePrivileges := config.Variables{
		"name":              config.StringVariable(roleFullyQualifiedName),
		"privileges":        config.ListVariable(grantablePrivilege1),
		"with_grant_option": config.BoolVariable(false),
	}
	configVariablesOnlyName := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithGrantablePrivileges,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
				ConfigVariables: configVariablesOnlyName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckNoResourceAttr(resourceName, "privileges.#"),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccount", roleFullyQualifiedName)),
				),
				// Due to limitations in the plugin SDK, returned warnings can not be asserted (see https://github.com/hashicorp/terraform-plugin-testing/issues/69).
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_PrivilegesReversed(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateRole)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE,CREATE ROLE|OnAccount", roleFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccountObject(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccountObject_gh2717(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	computePoolId, computePoolCleanup := acc.TestClient().ComputePool.CreateComputePool(t)
	t.Cleanup(computePoolCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name":         config.StringVariable(roleFullyQualifiedName),
		"compute_pool": config.StringVariable(computePoolId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject_gh2717"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", string(sdk.ObjectTypeComputePool)),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", computePoolId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnAccountObject|%s|%s", roleFullyQualifiedName, sdk.ObjectTypeComputePool, computePoolId.FullyQualifiedName())),
				),
			},
		},
	})
}

// This proves that infinite plan is not produced as in snowflake_grant_privileges_to_role.
// More details can be found in the fix pr https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2364.
func TestAcc_GrantPrivilegesToApplicationRole_OnAccountObject_InfinitePlan(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleId := role.ID()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject_InfinitePlan"),
				ConfigVariables: config.Variables{
					"name":     config.StringVariable(roleId.Name()),
					"database": config.StringVariable(acc.TestDatabaseName),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchema(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	schemaId := acc.TestClient().Ids.SchemaId()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(schemaId.DatabaseName()),
		"schema":            config.StringVariable(schemaId.Name()),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaId.FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchema_ExactlyOneOf(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema_ExactlyOneOf"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAllSchemasInDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnAllSchemasInDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnFutureSchemasInDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.future_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnFutureSchemasInDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tableId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          config.StringVariable(tableId.DatabaseName()),
		"schema":            config.StringVariable(tableId.SchemaName()),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnObject|TABLE|%s", roleFullyQualifiedName, tableId.FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFunctionWithArguments(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	function := acc.TestClient().Function.CreateSecure(t, sdk.DataTypeFloat)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name":          config.StringVariable(roleFullyQualifiedName),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          config.StringVariable(function.ID().DatabaseName()),
		"schema":            config.StringVariable(function.ID().SchemaName()),
		"with_grant_option": config.BoolVariable(false),
		"argument_type":     config.StringVariable(string(sdk.DataTypeFloat)),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeFunction)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|%s", roleFullyQualifiedName, function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFunctionWithoutArguments(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	function := acc.TestClient().Function.CreateSecure(t)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := config.Variables{
		"name":          config.StringVariable(roleFullyQualifiedName),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          config.StringVariable(function.ID().DatabaseName()),
		"schema":            config.StringVariable(function.ID().SchemaName()),
		"with_grant_option": config.BoolVariable(false),
		"argument_type":     config.StringVariable(""),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeFunction)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|%s", roleFullyQualifiedName, function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject_OwnershipPrivilege(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_OwnershipPrivilege"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnAll|TABLES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAllPipes(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypePipes)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MONITOR|OnSchemaObject|OnAll|PIPES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_InDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnFuture|TABLES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO [SNOW-1272222]: fix the test when it starts working on Snowflake side
func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	t.Skip("Fix after it starts working on Snowflake side, reference: SNOW-1272222")

	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported feature 'STREAMLIT'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnAll|STREAMLITS|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_UpdatePrivileges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(roleFullyQualifiedName),
			"database": config.StringVariable(databaseName),
		}
		if allPrivileges {
			configVariables["all_privileges"] = config.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]config.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = config.StringVariable(string(privilege))
			}
			configVariables["privileges"] = config.ListVariable(configPrivileges...)
		}
		return configVariables
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeModify,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeMonitor,
					sdk.AccountObjectPrivilegeUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeModify,
					sdk.AccountObjectPrivilegeMonitor,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MODIFY,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_UpdatePrivileges_SnowflakeChecked(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleId := role.ID()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(allPrivileges bool, privileges []string, schemaName string) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(roleId.FullyQualifiedName()),
			"database": config.StringVariable(schemaId.DatabaseName()),
		}
		if allPrivileges {
			configVariables["all_privileges"] = config.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]config.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = config.StringVariable(privilege)
			}
			configVariables["privileges"] = config.ListVariable(configPrivileges...)
		}
		if len(schemaName) > 0 {
			configVariables["schema_name"] = config.StringVariable(schemaName)
		}
		return configVariables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				}, ""),
				Check: queriedAccountRolePrivilegesEqualTo(
					t,
					roleId,
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
				ConfigVariables: configVariables(true, []string{}, ""),
				Check: queriedAccountRolePrivilegesContainAtLeast(
					t,
					roleId,
					sdk.AccountObjectPrivilegeCreateDatabaseRole.String(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
					sdk.AccountObjectPrivilegeUsage.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				}, ""),
				Check: queriedAccountRolePrivilegesEqualTo(
					t,
					roleId,
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
				ConfigVariables: configVariables(false, []string{
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				}, schemaId.Name()),
				Check: queriedAccountRolePrivilegesEqualTo(
					t,
					roleId,
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_AlwaysApply(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(roleFullyQualifiedName),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(databaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivileges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	externalShareId := createSharedDatabaseOnSecondaryAccount(t)

	databaseFromShare, databaseFromShareCleanup := acc.TestClient().Database.CreateDatabaseFromShare(t, externalShareId)
	t.Cleanup(databaseFromShareCleanup)

	configVariables := config.Variables{
		"role_name":            config.StringVariable(role.ID().Name()),
		"shared_database_name": config.StringVariable(databaseFromShare.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.AccountObjectPrivilegeImportedPrivileges.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivileges"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivileges"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1998 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivilegesOnSnowflakeDatabase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleName := role.ID().Name()
	configVariables := config.Variables{
		"role_name": config.StringVariable(roleName),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.AccountObjectPrivilegeImportedPrivileges.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", "\"SNOWFLAKE\""),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO(SNOW-1213622): Add test for custom applications using on_account_object.object_type = "DATABASE"

func TestAcc_GrantPrivilegesToAccountRole_MultiplePartsInRoleName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")
	_, roleCleanup := acc.TestClient().Role.CreateRoleWithIdentifier(t, roleId)
	t.Cleanup(roleCleanup)

	roleName := roleId.Name()
	configVariables := config.Variables{
		"name": config.StringVariable(roleName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleName),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2533 is fixed
func TestAcc_GrantPrivilegesToAccountRole_OnExternalVolume(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	externalVolumeId, cleanupExternalVolume := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(cleanupExternalVolume)

	configVariables := config.Variables{
		"name":            config.StringVariable(role.ID().FullyQualifiedName()),
		"external_volume": config.StringVariable(externalVolumeId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnExternalVolume"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", role.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "EXTERNAL VOLUME"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|USAGE|OnAccountObject|EXTERNAL VOLUME|%s", role.ID().FullyQualifiedName(), externalVolumeId.FullyQualifiedName())),
				),
			},
		},
	})
}

// proved https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2651
func TestAcc_GrantPrivilegesToAccountRole_MLPrivileges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	schemaId := acc.TestClient().Ids.SchemaId()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
		),
		"database":          config.StringVariable(schemaId.DatabaseName()),
		"schema":            config.StringVariable(schemaId.Name()),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SNOWFLAKE.ML.ANOMALY_DETECTION,CREATE SNOWFLAKE.ML.FORECAST|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaId.FullyQualifiedName())),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ChangeWithGrantOptionsOutsideOfTerraform_WithGrantOptions(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tableId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          config.StringVariable(tableId.DatabaseName()),
		"schema":            config.StringVariable(tableId.SchemaName()),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnTableToAccountRole(
						t,
						role.ID(),
						tableId,
						[]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeTruncate},
						false,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tableId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          config.StringVariable(tableId.DatabaseName()),
		"schema":            config.StringVariable(tableId.SchemaName()),
		"with_grant_option": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnTableToAccountRole(
						t,
						role.ID(),
						tableId,
						[]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeTruncate},
						true,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// TODO [SNOW-1431726]: Move to helpers
func revokeAndGrantPrivilegesOnTableToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	tableName sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	client := acc.TestClient()

	client.Grant.RevokePrivilegesOnSchemaObjectFromAccountRole(t, accountRoleId, sdk.ObjectTypeTable, tableName, privileges)
	client.Grant.GrantPrivilegesOnSchemaObjectToAccountRole(t, accountRoleId, sdk.ObjectTypeTable, tableName, privileges, withGrantOption)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveGrantedObjectOutsideTerraform(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(database.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveAccountRoleOutsideTerraform(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseId := acc.TestClient().Ids.DatabaseId()
	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(databaseId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { roleCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2689 is fixed
func TestAcc_GrantPrivilegesToAccountRole_AlwaysApply_SetAfterCreate(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(roleFullyQualifiedName),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(databaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory:    acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables:    configVariables(true),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

// TODO [SNOW-1431726]: Move to helpers
func createSharedDatabaseOnSecondaryAccount(t *testing.T) sdk.ExternalObjectIdentifier {
	t.Helper()

	database, databaseCleanup := acc.SecondaryTestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := acc.SecondaryTestClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	_ = acc.SecondaryTestClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeReferenceUsage})

	accountName := acc.TestClient().Context.CurrentAccount(t)
	accountId := sdk.NewAccountIdentifierFromAccountLocator(accountName)
	acc.SecondaryTestClient().Share.SetAccountOnShare(t, accountId, share.ID())

	return sdk.NewExternalObjectIdentifier(acc.SecondaryTestClient().Account.GetAccountIdentifier(t), share.ID())
}

func queriedAccountRolePrivilegesEqualTo(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesEqualTo(func() ([]sdk.Grant, error) {
		return acc.TestClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, privileges...)
}

func queriedAccountRolePrivilegesContainAtLeast(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesContainAtLeast(func() ([]sdk.Grant, error) {
		return acc.TestClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, roleName, privileges...)
}

func queriedAccountRolePrivilegesDoNotContain(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesDoNotContain(func() ([]sdk.Grant, error) {
		return acc.TestClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, privileges...)
}

func TestAcc_GrantPrivilegesToAccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	schemaId := acc.TestClient().Ids.SchemaId()
	quotedSchemaId := fmt.Sprintf(`\"%s\".\"%s\"`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToAccountRoleBasicConfig(role.ID(), quotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToAccountRoleBasicConfig(role.ID(), quotedSchemaId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToAccountRoleBasicConfig(roleId sdk.AccountObjectIdentifier, fullyQualifiedSchemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "%[1]s"
  privileges         = ["USAGE"]

  on_schema {
    schema_name = "%[2]s"
  }
}
`, roleId.Name(), fullyQualifiedSchemaName)
}

func TestAcc_GrantPrivilegesToAccountRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	schemaId := acc.TestClient().Ids.SchemaId()
	unquotedSchemaId := fmt.Sprintf(`%s.%s`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToAccountRoleBasicConfig(role.ID(), unquotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_schema.0.schema_name", unquotedSchemaId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToAccountRoleBasicConfig(role.ID(), unquotedSchemaId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_schema.0.schema_name", unquotedSchemaId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2807
func TestAcc_GrantPrivilegesToAccountRole_OnDataset_issue2807(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeDatasets.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeDatasets)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnFuture|DATASETS|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3050
func TestAcc_GrantPrivilegesToAccountRole_OnFutureModels_issue3050(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	accountRoleName := role.ID().Name()
	databaseName := acc.TestClient().Ids.DatabaseId().Name()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.95.0"),
				Config:             grantPrivilegesToAccountRoleOnFutureInDatabaseConfig(accountRoleName, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseName),
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToAccountRoleOnFutureInDatabaseConfig(accountRoleName, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseName),
			},
		},
	})
}

func grantPrivilegesToAccountRoleOnFutureInDatabaseConfig(accountRoleName string, privileges []string, objectTypePlural sdk.PluralObjectType, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "%[1]s"
  privileges        = [ %[2]s ]

  on_schema_object {
    future {
      object_type_plural = "%[3]s"
      in_database        = "%[4]s"
    }
  }
}
`, accountRoleName, strings.Join(collections.Map(privileges, strconv.Quote), ","), objectTypePlural, databaseName)
}

// This test proves that managing grants on HYBRID TABLE is not supported in Snowflake. TABLE should be used instead.
func TestAcc_GrantPrivileges_OnObject_HybridTable_ToAccountRole_Fails(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	hybridTableId, hybridTableCleanup := acc.TestClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	configVariables := func(objectType sdk.ObjectType) config.Variables {
		cfg := config.Variables{
			"account_role_name": config.StringVariable(role.ID().FullyQualifiedName()),
			"privileges": config.ListVariable(
				config.StringVariable(string(sdk.SchemaObjectPrivilegeApplyBudget)),
			),
			"hybrid_table_fully_qualified_name": config.StringVariable(hybridTableId.FullyQualifiedName()),
			"object_type":                       config.StringVariable(string(objectType)),
		}
		return cfg
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeHybridTable),
				ExpectError:     regexp.MustCompile("syntax error line 1 at position 28 unexpected 'TABLE"),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeTable),
			},
		},
	})
}
