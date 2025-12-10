//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_BasicUseCase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	configVariablesUpdated := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateNetworkPolicy)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
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
			// Import
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - configuration change
			{
				// always_apply is not tested here as it is covered in other tests and produces non-empty plans which may interfere with incorrect resource behavior
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateNetworkPolicy)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE NETWORK POLICY,CREATE DATABASE|OnAccount", roleFullyQualifiedName)),
				),
			},
			// Update - external change
			{
				PreConfig: func() {
					// We are not granting anything as new privileges won't be detected (authoritative grants would be used for this)
					testClient().Grant.RevokeGlobalPrivilegesFromAccountRole(t, role.ID(), []sdk.GlobalPrivilege{sdk.GlobalPrivilegeCreateDatabase})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateNetworkPolicy)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE NETWORK POLICY,CREATE DATABASE|OnAccount", roleFullyQualifiedName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_gh3153(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeManageShareTarget)),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_gh3153"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name":         tfconfig.StringVariable(roleFullyQualifiedName),
		"always_apply": tfconfig.BoolVariable(false),
	}
	configVariablesWithAlwaysApply := tfconfig.Variables{
		"name":         tfconfig.StringVariable(roleFullyQualifiedName),
		"always_apply": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.5"),
				Config:            grantAllPrivilegesToAccountRoleBasicConfig(role.ID()),
				ExpectError:       regexp.MustCompile(`Error: 003011 \(42501\): Grant partially executed`),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
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
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantablePrivilege1 := tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase))
	grantablePrivilege2 := tfconfig.StringVariable(string(sdk.GlobalPrivilegeApplyAggregationPolicy))
	nonGrantablePrivilege := tfconfig.StringVariable("MANAGE LISTING AUTO FULFILLMENT")
	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariablesWithGrantablePrivileges := tfconfig.Variables{
		"name":              tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges":        tfconfig.ListVariable(grantablePrivilege1),
		"with_grant_option": tfconfig.BoolVariable(true),
	}
	configVariablesWithNonGrantablePrivileges := tfconfig.Variables{
		"name":              tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges":        tfconfig.ListVariable(grantablePrivilege1, grantablePrivilege2, nonGrantablePrivilege),
		"with_grant_option": tfconfig.BoolVariable(true),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithNonGrantablePrivileges,
				ExpectError:     regexp.MustCompile("grant partially executed"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithGrantablePrivileges,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithNonGrantablePrivileges,
				ExpectError:     regexp.MustCompile("grant partially executed"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_ChangeListOfPrivilegesToAllPrivileges(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantablePrivilege1 := tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase))
	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariablesWithGrantablePrivileges := tfconfig.Variables{
		"name":              tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges":        tfconfig.ListVariable(grantablePrivilege1),
		"with_grant_option": tfconfig.BoolVariable(false),
	}
	configVariablesOnlyName := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariablesWithGrantablePrivileges,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount_AllPrivileges"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccountObject_BasicUseCase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name":     tfconfig.StringVariable(roleFullyQualifiedName),
		"database": tfconfig.StringVariable(databaseName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	configVariablesUpdated := tfconfig.Variables{
		"name":     tfconfig.StringVariable(roleFullyQualifiedName),
		"database": tfconfig.StringVariable(databaseName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeMonitor)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			// Import
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - config changes
			{
				// always_apply is not tested here as it is covered in other tests and produces non-empty plans which may interfere with incorrect resource behavior
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE ROLE,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			// Update - external changes
			{
				PreConfig: func() {
					// We are not granting anything as new privileges won't be detected (authoritative grants would be used for this)
					testClient().Grant.RevokePrivilegesOnDatabaseFromAccountRole(t, role.ID(), testClient().Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateDatabaseRole})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE ROLE,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccountObject_gh2717(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name":         tfconfig.StringVariable(roleFullyQualifiedName),
		"compute_pool": tfconfig.StringVariable(computePool.ID().Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject_gh2717"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", string(sdk.ObjectTypeComputePool)),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", computePool.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnAccountObject|%s|%s", roleFullyQualifiedName, sdk.ObjectTypeComputePool, computePool.ID().FullyQualifiedName())),
				),
			},
		},
	})
}

// This proves that infinite plan is not produced as in snowflake_grant_privileges_to_role.
// More details can be found in the fix pr https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2364.
func TestAcc_GrantPrivilegesToApplicationRole_OnAccountObject_InfinitePlan(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleId := role.ID()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject_InfinitePlan"),
				ConfigVariables: tfconfig.Variables{
					"name":     tfconfig.StringVariable(roleId.Name()),
					"database": tfconfig.StringVariable(TestDatabaseName),
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

func TestAcc_GrantPrivilegesToAccountRole_OnSchema_BasicUseCase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	schemaId := testClient().Ids.SchemaId()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          tfconfig.StringVariable(schemaId.DatabaseName()),
		"schema":            tfconfig.StringVariable(schemaId.Name()),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	configVariablesUpdated := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateAlert)),
		),
		"database":          tfconfig.StringVariable(schemaId.DatabaseName()),
		"schema":            tfconfig.StringVariable(schemaId.Name()),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
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
			// Import
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - config changes
			{
				// always_apply is not tested here as it is covered in other tests and produces non-empty plans which may interfere with incorrect resource behavior
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateAlert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE ALERT,CREATE TABLE|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaId.FullyQualifiedName())),
				),
			},
			// Update - external changes
			{
				PreConfig: func() {
					// We are not granting anything as new privileges won't be detected (authoritative grants would be used for this)
					testClient().Grant.RevokePrivilegesOnSchemaFromAccountRole(t, role.ID(), testClient().Ids.SchemaId(), []sdk.SchemaPrivilege{sdk.SchemaPrivilegeCreateTable})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariablesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateAlert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE ALERT,CREATE TABLE|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaId.FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchema_ExactlyOneOf(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema_ExactlyOneOf"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAllSchemasInDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          tfconfig.StringVariable(databaseName),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnFutureSchemasInDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          tfconfig.StringVariable(databaseName),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	configVariables := tfconfig.Variables{
		"name":       tfconfig.StringVariable(roleFullyQualifiedName),
		"table_name": tfconfig.StringVariable(tableId.Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          tfconfig.StringVariable(tableId.DatabaseName()),
		"schema":            tfconfig.StringVariable(tableId.SchemaName()),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFunctionWithArguments(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	function := testClient().Function.CreateSecure(t, sdk.DataTypeFloat)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name":          tfconfig.StringVariable(roleFullyQualifiedName),
		"function_name": tfconfig.StringVariable(function.ID().Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          tfconfig.StringVariable(function.ID().DatabaseName()),
		"schema":            tfconfig.StringVariable(function.ID().SchemaName()),
		"with_grant_option": tfconfig.BoolVariable(false),
		"argument_type":     tfconfig.StringVariable(string(sdk.DataTypeFloat)),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFunctionWithoutArguments(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	function := testClient().Function.CreateSecure(t)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name":          tfconfig.StringVariable(roleFullyQualifiedName),
		"function_name": tfconfig.StringVariable(function.ID().Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          tfconfig.StringVariable(function.ID().DatabaseName()),
		"schema":            tfconfig.StringVariable(function.ID().SchemaName()),
		"with_grant_option": tfconfig.BoolVariable(false),
		"argument_type":     tfconfig.StringVariable(""),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFunction"),
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_OwnershipPrivilege"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           tfconfig.StringVariable(databaseName),
		"object_type_plural": tfconfig.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAllPipes(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
		),
		"database":          tfconfig.StringVariable(databaseName),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_InDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           tfconfig.StringVariable(databaseName),
		"object_type_plural": tfconfig.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           tfconfig.StringVariable(databaseName),
		"object_type_plural": tfconfig.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  tfconfig.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnFuture|STREAMLITS|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           tfconfig.StringVariable(databaseName),
		"object_type_plural": tfconfig.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) tfconfig.Variables {
		configVariables := tfconfig.Variables{
			"name":     tfconfig.StringVariable(roleFullyQualifiedName),
			"database": tfconfig.StringVariable(databaseName),
		}
		if allPrivileges {
			configVariables["all_privileges"] = tfconfig.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]tfconfig.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = tfconfig.StringVariable(string(privilege))
			}
			configVariables["privileges"] = tfconfig.ListVariable(configPrivileges...)
		}
		return configVariables
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleId := role.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(allPrivileges bool, privileges []string, schemaName string) tfconfig.Variables {
		configVariables := tfconfig.Variables{
			"name":     tfconfig.StringVariable(roleId.FullyQualifiedName()),
			"database": tfconfig.StringVariable(schemaId.DatabaseName()),
		}
		if allPrivileges {
			configVariables["all_privileges"] = tfconfig.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]tfconfig.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = tfconfig.StringVariable(privilege)
			}
			configVariables["privileges"] = tfconfig.ListVariable(configPrivileges...)
		}
		if len(schemaName) > 0 {
			configVariables["schema_name"] = tfconfig.StringVariable(schemaName)
		}
		return configVariables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) tfconfig.Variables {
		return tfconfig.Variables{
			"name":           tfconfig.StringVariable(roleFullyQualifiedName),
			"all_privileges": tfconfig.BoolVariable(true),
			"database":       tfconfig.StringVariable(databaseName),
			"always_apply":   tfconfig.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	externalShareId := createSharedDatabaseOnSecondaryAccount(t)

	databaseFromShare, databaseFromShareCleanup := testClient().Database.CreateDatabaseFromShare(t, externalShareId)
	t.Cleanup(databaseFromShareCleanup)

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToAccountObjectConfig(role.ID(), databaseFromShare.ID(), sdk.AccountObjectPrivilegeImportedPrivileges.String()),
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
				Config:            grantPrivilegesToAccountObjectConfig(role.ID(), databaseFromShare.ID(), sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivileges_Validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config:      grantPrivilegesToAccountObjectConfigInvalid(),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("IMPORTED PRIVILEGES cannot be used with other privileges"),
			},
		},
	})
}

// prove https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2803 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivileges_issue2803(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	externalShareId := createSharedDatabaseOnSecondaryAccount(t)

	databaseFromShare, databaseFromShareCleanup := testClient().Database.CreateDatabaseFromShare(t, externalShareId)
	t.Cleanup(databaseFromShareCleanup)

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.5.0"),
				Config:            grantPrivilegesToAccountObjectConfig(role.ID(), databaseFromShare.ID(), sdk.AccountObjectPrivilegeImportedPrivileges.String()),
			},
			// Expect an error when the import privilege is revoked externally in 2.5.0.
			{
				PreConfig: func() {
					testClient().Grant.RevokePrivilegesOnDatabaseFromAccountRole(t, role.ID(), databaseFromShare.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImportedPrivileges})
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.5.0"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceName, "privileges", tfjson.ActionUpdate, sdk.Pointer("[]"), sdk.Pointer(fmt.Sprintf("[%s]", string(sdk.AccountObjectPrivilegeImportedPrivileges)))),
					},
				},
				Config:      grantPrivilegesToAccountObjectConfig(role.ID(), databaseFromShare.ID(), sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				ExpectError: regexp.MustCompile("Failed to revoke privileges to add"),
			},
			// Prove the fix in later versions.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceName, "privileges", tfjson.ActionUpdate, sdk.Pointer("[]"), sdk.Pointer(fmt.Sprintf("[%s]", string(sdk.AccountObjectPrivilegeImportedPrivileges)))),
					},
				},
				Config: grantPrivilegesToAccountObjectConfig(role.ID(), databaseFromShare.ID(), sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				),
			},
		},
	})
}

func grantPrivilegesToAccountObjectConfig(roleName, databaseName sdk.AccountObjectIdentifier, privilege string) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
	account_role_name = "\"%s\""
	privileges = ["%s"]
	on_account_object {
		object_type = "DATABASE"
		object_name = "\"%s\""
	}
}
`, roleName.Name(), privilege, databaseName.Name())
}

func grantPrivilegesToAccountObjectConfigInvalid() string {
	return `
resource "snowflake_grant_privileges_to_account_role" "test" {
	account_role_name = "ROLE"
	privileges = ["IMPORTED PRIVILEGES", "APPLYBUDGET"]
	on_account_object {
		object_type = "DATABASE"
		object_name = "DB"
	}
}
`
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1998 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivilegesOnSnowflakeDatabase(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleName := role.ID().Name()
	configVariables := tfconfig.Variables{
		"role_name": tfconfig.StringVariable(roleName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(sdk.AccountObjectPrivilegeImportedPrivileges.String()),
		),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
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
	roleId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")
	_, roleCleanup := testClient().Role.CreateRoleWithIdentifier(t, roleId)
	t.Cleanup(roleCleanup)

	roleName := roleId.Name()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			tfconfig.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	externalVolumeId, cleanupExternalVolume := testClient().ExternalVolume.Create(t)
	t.Cleanup(cleanupExternalVolume)

	configVariables := tfconfig.Variables{
		"name":            tfconfig.StringVariable(role.ID().FullyQualifiedName()),
		"external_volume": tfconfig.StringVariable(externalVolumeId.Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnExternalVolume"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	schemaId := testClient().Ids.SchemaId()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			tfconfig.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
		),
		"database":          tfconfig.StringVariable(schemaId.DatabaseName()),
		"schema":            tfconfig.StringVariable(schemaId.Name()),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := tfconfig.Variables{
		"name":       tfconfig.StringVariable(roleFullyQualifiedName),
		"table_name": tfconfig.StringVariable(tableId.Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          tfconfig.StringVariable(tableId.DatabaseName()),
		"schema":            tfconfig.StringVariable(tableId.SchemaName()),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := tfconfig.Variables{
		"name":       tfconfig.StringVariable(roleFullyQualifiedName),
		"table_name": tfconfig.StringVariable(tableId.Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          tfconfig.StringVariable(tableId.DatabaseName()),
		"schema":            tfconfig.StringVariable(tableId.SchemaName()),
		"with_grant_option": tfconfig.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
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
	client := testClient()

	client.Grant.RevokePrivilegesOnSchemaObjectFromAccountRole(t, accountRoleId, sdk.ObjectTypeTable, tableName, privileges)
	client.Grant.GrantPrivilegesOnSchemaObjectToAccountRole(t, accountRoleId, sdk.ObjectTypeTable, tableName, privileges, withGrantOption)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveGrantedObjectOutsideTerraform(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()

	configVariables := tfconfig.Variables{
		"name":     tfconfig.StringVariable(roleFullyQualifiedName),
		"database": tfconfig.StringVariable(database.ID().Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseCleanup() },
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveAccountRoleOutsideTerraform(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseId := testClient().Ids.DatabaseId()
	configVariables := tfconfig.Variables{
		"name":     tfconfig.StringVariable(roleFullyQualifiedName),
		"database": tfconfig.StringVariable(databaseId.Name()),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			tfconfig.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": tfconfig.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { roleCleanup() },
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2689 is fixed
func TestAcc_GrantPrivilegesToAccountRole_AlwaysApply_SetAfterCreate(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) tfconfig.Variables {
		return tfconfig.Variables{
			"name":           tfconfig.StringVariable(roleFullyQualifiedName),
			"all_privileges": tfconfig.BoolVariable(true),
			"database":       tfconfig.StringVariable(databaseName),
			"always_apply":   tfconfig.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory:    ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
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

	database, databaseCleanup := secondaryTestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := secondaryTestClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	_ = secondaryTestClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage, sdk.ObjectPrivilegeReferenceUsage})

	accountName := testClient().Context.CurrentAccount(t)
	accountId := sdk.NewAccountIdentifierFromAccountLocator(accountName)
	secondaryTestClient().Share.SetAccountOnShare(t, accountId, share.ID())

	return sdk.NewExternalObjectIdentifier(secondaryTestClient().Account.GetAccountIdentifier(t), share.ID())
}

func queriedAccountRolePrivilegesEqualTo(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesEqualTo(func() ([]sdk.Grant, error) {
		return testClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, privileges...)
}

func queriedAccountRolePrivilegesContainAtLeast(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesContainAtLeast(func() ([]sdk.Grant, error) {
		return testClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, roleName, privileges...)
}

func queriedAccountRolePrivilegesDoNotContain(t *testing.T, roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesDoNotContain(func() ([]sdk.Grant, error) {
		return testClient().Grant.ShowGrantsToAccountRole(t, roleName)
	}, privileges...)
}

func TestAcc_GrantPrivilegesToAccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	schemaId := testClient().Ids.SchemaId()
	quotedSchemaId := fmt.Sprintf(`\"%s\".\"%s\"`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToAccountRoleBasicConfig(role.ID(), quotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	schemaId := testClient().Ids.SchemaId()
	unquotedSchemaId := fmt.Sprintf(`%s.%s`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToAccountRoleBasicConfig(role.ID(), unquotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "on_schema.0.schema_name", unquotedSchemaId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_account_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", role.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	roleFullyQualifiedName := role.ID().FullyQualifiedName()
	databaseName := testClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := tfconfig.Variables{
		"name": tfconfig.StringVariable(roleFullyQualifiedName),
		"privileges": tfconfig.ListVariable(
			tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           tfconfig.StringVariable(databaseName),
		"object_type_plural": tfconfig.StringVariable(sdk.PluralObjectTypeDatasets.String()),
		"with_grant_option":  tfconfig.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_account_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
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
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	accountRoleName := role.ID().Name()
	databaseName := testClient().Ids.DatabaseId().Name()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            grantPrivilegesToAccountRoleOnFutureInDatabaseConfig(accountRoleName, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseName),
				// Previously, we expected a non-empty plan, because Snowflake returned MODULE instead of MODEL in SHOW FUTURE GRANTS.
				// Now, this behavior is fixed in Snowflake, and the plan is empty.
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	hybridTableId, hybridTableCleanup := testClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	configVariables := func(objectType sdk.ObjectType) tfconfig.Variables {
		cfg := tfconfig.Variables{
			"account_role_name": tfconfig.StringVariable(role.ID().FullyQualifiedName()),
			"privileges": tfconfig.ListVariable(
				tfconfig.StringVariable(string(sdk.SchemaObjectPrivilegeApplyBudget)),
			),
			"hybrid_table_fully_qualified_name": tfconfig.StringVariable(hybridTableId.FullyQualifiedName()),
			"object_type":                       tfconfig.StringVariable(string(objectType)),
		}
		return cfg
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeHybridTable),
				ExpectError:     regexp.MustCompile("syntax error line 1 at position 28 unexpected 'TABLE"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeTable),
			},
		},
	})
}

// queriedAccountRolePrivilegesEqualTo will check if all the privileges specified in the argument are granted in Snowflake.
func queriedPrivilegesEqualTo(query func() ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		grants, err := query()
		if err != nil {
			return err
		}

		grantedPrivileges := collections.Map(grants, func(grant sdk.Grant) string {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				return ""
			}
			return grant.Privilege
		})
		grantedPrivileges = slices.DeleteFunc(grantedPrivileges, func(privilege string) bool { return privilege == "" })

		if !slices.Equal(grantedPrivileges, privileges) {
			return fmt.Errorf("granted privileges: %v, not equal to expected set: %v", grantedPrivileges, privileges)
		}

		return nil
	}
}

// queriedAccountRolePrivilegesContainAtLeast will check if all the privileges specified in the argument are granted in Snowflake.
// Any additional grants will be ignored.
func queriedPrivilegesContainAtLeast(query func() ([]sdk.Grant, error), roleName sdk.ObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		grants, err := query()
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		notAllPrivilegesInGrantedPrivileges := slices.ContainsFunc(privileges, func(privilege string) bool {
			return !slices.Contains(grantedPrivileges, privilege)
		})
		if notAllPrivilegesInGrantedPrivileges {
			return fmt.Errorf("not every privilege from the list: %v was found in grant privileges: %v, for role name: %s", privileges, grantedPrivileges, roleName.FullyQualifiedName())
		}

		return nil
	}
}

// queriedPrivilegesDoNotContain will check if all of the privileges specified in the argument are not granted in Snowflake.
func queriedPrivilegesDoNotContain(query func() ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		grants, err := query()
		if err != nil {
			return err
		}
		for _, grant := range grants {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				continue
			}
			if slices.Contains(privileges, grant.Privilege) {
				return fmt.Errorf("grant not expected, grants: %v should not contain any privilege from %v", grants, privileges)
			}
		}

		return nil
	}
}

func TestAcc_GrantPrivilegesToAccountRole_issue3992(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	databaseId := testClient().Ids.DatabaseId()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.6.0"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("snowflake_grant_privileges_to_account_role.test", tfjsonpath.New("privileges").AtSliceIndex(0), knownvalue.StringExact("USAGE")),
					},
				},
				Config: configIssue3992(role.ID(), databaseId),
				// It fails, even though the plan says it's known.
				ExpectError: regexp.MustCompile("panic: value is unknown"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("snowflake_grant_privileges_to_account_role.test", tfjsonpath.New("privileges").AtSliceIndex(0), knownvalue.StringExact("USAGE")),
					},
				},
				Config: configIssue3992(role.ID(), databaseId),
			},
		},
	})
}

func configIssue3992(roleId sdk.AccountObjectIdentifier, dbId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "%[1]s"
  privileges        = [var.privilege]

  on_account_object {
      object_type = "DATABASE"
      object_name = "%[2]s"
  }
}
variable "privilege" {
  type = string
  default = "USAGE"
}
`, roleId.Name(), dbId.Name())
}

func TestAcc_GrantPrivilegesToAccountRole_StrictRoleManagement_OnCreate(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	testClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(t, role.ID(), database.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
	resourceModelWithStrictRoleManagement := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithPrivileges(string(sdk.AccountObjectPrivilegeMonitor)).
		WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type": tfconfig.StringVariable(sdk.ObjectTypeDatabase.String()),
			"object_name": tfconfig.StringVariable(database.ID().Name()),
		})).
		WithStrictPrivilegeManagement(true)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, resourceModelWithStrictRoleManagement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "strict_privilege_management", "true"),
					queriedAccountRolePrivilegesEqualTo(t, role.ID(), string(sdk.AccountObjectPrivilegeMonitor), string(sdk.AccountObjectPrivilegeUsage)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceModelWithStrictRoleManagement.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges", tfjson.ActionUpdate,
							sdk.String(fmt.Sprintf("[%s %s]", sdk.AccountObjectPrivilegeMonitor, sdk.AccountObjectPrivilegeUsage)),
							sdk.String(fmt.Sprintf("[%s]", sdk.AccountObjectPrivilegeMonitor)),
						),
					},
				},
				Config: config.FromModels(t, providerModel, resourceModelWithStrictRoleManagement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "strict_privilege_management", "true"),
					queriedAccountRolePrivilegesEqualTo(t, role.ID(), string(sdk.AccountObjectPrivilegeMonitor)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_StrictRoleManagement(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
	resourceModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithPrivileges(string(sdk.AccountObjectPrivilegeMonitor)).
		WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type": tfconfig.StringVariable(sdk.ObjectTypeDatabase.String()),
			"object_name": tfconfig.StringVariable(database.ID().Name()),
		}))
	resourceModelWithStrictRoleManagement := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
		WithPrivileges(string(sdk.AccountObjectPrivilegeMonitor)).
		WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type": tfconfig.StringVariable(sdk.ObjectTypeDatabase.String()),
			"object_name": tfconfig.StringVariable(database.ID().Name()),
		})).
		WithStrictPrivilegeManagement(true)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, resourceModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "account_role_name", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "strict_privilege_management", "false"),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "on_account_object.0.object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "on_account_object.0.object_name", database.ID().Name()),
				),
			},
			{
				PreConfig: func() {
					testClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(t, role.ID(), database.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, resourceModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceModel.ResourceReference(), "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					queriedAccountRolePrivilegesContainAtLeast(t, role.ID(), string(sdk.AccountObjectPrivilegeMonitor), string(sdk.AccountObjectPrivilegeUsage)),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceModel.ResourceReference(), plancheck.ResourceActionUpdate),
						// TODO: change to: from 2 to 1 privilege
						planchecks.ExpectChange(resourceModel.ResourceReference(), "privileges", tfjson.ActionUpdate, sdk.String(fmt.Sprintf("[%s]", sdk.AccountObjectPrivilegeMonitor)), sdk.String(fmt.Sprintf("[%s]", sdk.AccountObjectPrivilegeMonitor))),
						//planchecks.ExpectChange(resourceName, "strict_privilege_management", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
					},
				},
				Config: config.FromModels(t, providerModel, resourceModelWithStrictRoleManagement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceModelWithStrictRoleManagement.ResourceReference(), "strict_privilege_management", "true"),
					queriedAccountRolePrivilegesEqualTo(t, role.ID(), string(sdk.AccountObjectPrivilegeMonitor)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_StrictRoleManagement_Validation(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	databaseId := testClient().Ids.DatabaseId()
	schemaId := testClient().Ids.SchemaId()

	providerModelWithExperiment := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
	resourceModelMissingExperiment := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_missing_experiment", role.ID().Name()).
		WithPrivileges(string(sdk.GlobalPrivilegeCreateDatabase)).
		WithOnAccount(true).
		WithStrictPrivilegeManagement(true)
	resourceModelAllPrivileges := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_all_privileges", role.ID().Name()).
		WithOnAccount(true).
		WithAllPrivileges(true).
		WithStrictPrivilegeManagement(true)
	resourceModelOnSchemaAll := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_on_schema_all", role.ID().Name()).
		WithPrivileges(string(sdk.SchemaPrivilegeUsage)).
		WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"all_schemas_in_database": tfconfig.StringVariable(databaseId.FullyQualifiedName()),
		})).
		WithStrictPrivilegeManagement(true)
	resourceModelOnSchemaFuture := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_on_schema_future", role.ID().Name()).
		WithPrivileges(string(sdk.SchemaPrivilegeUsage)).
		WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"future_schemas_in_database": tfconfig.StringVariable(databaseId.FullyQualifiedName()),
		})).
		WithStrictPrivilegeManagement(true)
	resourceModelOnSchemaObjectAll := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_on_schema_object_all", role.ID().Name()).
		WithPrivileges(string(sdk.SchemaObjectPrivilegeSelect)).
		WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"all": tfconfig.ListVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type_plural": tfconfig.StringVariable(string(sdk.PluralObjectTypeTables)),
					"in_schema":          tfconfig.StringVariable(schemaId.FullyQualifiedName()),
				}),
			),
		})).
		WithStrictPrivilegeManagement(true)
	resourceModelOnSchemaObjectFuture := model.GrantPrivilegesToAccountRole("test_strict_privileges_validation_on_schema_object_future", role.ID().Name()).
		WithPrivileges(string(sdk.SchemaObjectPrivilegeSelect)).
		WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"future": tfconfig.ListVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type_plural": tfconfig.StringVariable(string(sdk.PluralObjectTypeTables)),
					"in_schema":          tfconfig.StringVariable(schemaId.FullyQualifiedName()),
				}),
			),
		})).
		WithStrictPrivilegeManagement(true)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, resourceModelMissingExperiment),
				ExpectError: regexp.MustCompile(
					"`strict_privilege_management`.*`GRANTS_STRICT_PRIVILEGE_MANAGEMENT`"),
			},
			{
				Config:      config.FromModels(t, providerModelWithExperiment, resourceModelAllPrivileges),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"strict_privilege_management": conflicts with all_privileges`),
			},
			{
				Config:      config.FromModels(t, providerModelWithExperiment, resourceModelOnSchemaAll),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"strict_privilege_management": conflicts with on_schema\.0\.all_schemas_in_database`),
			},
			{
				Config:      config.FromModels(t, providerModelWithExperiment, resourceModelOnSchemaFuture),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"strict_privilege_management": conflicts with on_schema\.0\.future_schemas_in_database`),
			},
			{
				Config:      config.FromModels(t, providerModelWithExperiment, resourceModelOnSchemaObjectAll),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"strict_privilege_management": conflicts with on_schema_object\.0\.all`),
			},
			{
				Config:      config.FromModels(t, providerModelWithExperiment, resourceModelOnSchemaObjectFuture),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"strict_privilege_management": conflicts with on_schema_object\.0\.future`),
			},
		},
	})
}

// TestAcc_GrantPrivilegesToAccountRole_StrictRoleManagement_PrivilegesThatCannotBeRevoked tests edge cases
// for strict_privilege_management with privileges that have special behavior or cannot be revoked.
//
// Edge cases covered:
//   - OWNERSHIP: Cannot be specified in this resource (blocked at validation level, tested separately in
//     TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject_OwnershipPrivilege)
//   - External OWNERSHIP grants: When ownership is granted externally to the role, it should be ignored
//     by strict_privilege_management because OWNERSHIP grants have empty grantedBy field and the resource
//     filters these out (see ReadGrantPrivilegesToAccountRole implementation)
//   - USAGE privilege: May come from IMPORTED_PRIVILEGES and has special handling in the resource
//   - Privileges granted externally with different grant option: Should be handled correctly
func TestAcc_GrantPrivilegesToAccountRole_StrictRoleManagement_PrivilegesThatCannotBeRevoked(t *testing.T) {
	// Test 1: External ownership grant should be ignored by strict_privilege_management
	// because ownership grants have empty grantedBy and are filtered out
	t.Run("ExternalOwnershipGrantIsIgnored", func(t *testing.T) {
		role, roleCleanup := testClient().Role.CreateRole(t)
		t.Cleanup(roleCleanup)

		database, databaseCleanup := testClient().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		// Grant OWNERSHIP on the database to the role externally
		testClient().Grant.GrantOwnershipToAccountRole(t, role.ID(), sdk.ObjectTypeDatabase, database.ID())

		providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
		// Only manage MONITOR privilege, not OWNERSHIP
		resourceModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
			WithPrivileges(string(sdk.AccountObjectPrivilegeMonitor)).
			WithOnAccountObjectValue(tfconfig.StringVariable(database.ID().Name())).
			WithStrictPrivilegeManagement(true)
		resourceName := resourceModel.ResourceReference()

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					// Create the grant with strict_privilege_management enabled
					// The external OWNERSHIP grant should be ignored (not added to state)
					Config: config.FromModels(t, providerModel, resourceModel),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
						resource.TestCheckResourceAttr(resourceName, "strict_privilege_management", "true"),
					),
				},
				{
					// Verify no changes on re-apply (OWNERSHIP should not cause drift)
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
						},
					},
					Config: config.FromModels(t, providerModel, resourceModel),
				},
			},
		})
	})

	// Test 2: External privilege grant with different WITH GRANT OPTION should be detected
	// and revoked when strict_privilege_management is enabled
	t.Run("ExternalPrivilegeWithDifferentGrantOption", func(t *testing.T) {
		role, roleCleanup := testClient().Role.CreateRole(t)
		t.Cleanup(roleCleanup)

		database, databaseCleanup := testClient().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
		// Configure resource with MONITOR privilege without grant option
		resourceModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
			WithPrivileges(string(sdk.AccountObjectPrivilegeMonitor)).
			WithOnAccountObjectValue(tfconfig.StringVariable(database.ID().Name())).
			WithWithGrantOption(false).
			WithStrictPrivilegeManagement(true)
		resourceName := resourceModel.ResourceReference()

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
			Steps: []resource.TestStep{
				{
					// Create the resource first
					Config: config.FromModels(t, providerModel, resourceModel),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
					),
				},
				{
					PreConfig: func() {
						// Grant USAGE privilege externally (this should be detected and revoked)
						testClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(
							t,
							role.ID(),
							database.ID(),
							[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
							false,
						)
					},
					// With strict_privilege_management, the external USAGE grant should be detected
					// and the plan should show it needs to be revoked
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						},
					},
					Config: config.FromModels(t, providerModel, resourceModel),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeMonitor)),
						// Verify the external USAGE privilege was revoked
						queriedAccountRolePrivilegesEqualTo(
							t,
							role.ID(),
							string(sdk.AccountObjectPrivilegeMonitor),
						),
					),
				},
			},
		})
	})

	// Test 3: USAGE privilege from IMPORTED_PRIVILEGES handling
	// When IMPORTED_PRIVILEGES is granted, it shows as USAGE in the grants
	// This test verifies the special handling in the resource works correctly
	//t.Run("ImportedPrivilegesUsageHandling", func(t *testing.T) {
	//	role, roleCleanup := testClient().Role.CreateRole(t)
	//	t.Cleanup(roleCleanup)
	//
	//	// Create a database from a share (which allows IMPORTED_PRIVILEGES)
	//	shareExternalId := createShareableDatabase(t)
	//	databaseFromShare, databaseFromShareCleanup := testClient().Database.CreateDatabaseFromShare(t, shareExternalId)
	//	t.Cleanup(databaseFromShareCleanup)
	//
	//	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
	//	// Use IMPORTED_PRIVILEGES (which internally maps to USAGE)
	//	resourceModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
	//		WithPrivileges(string(sdk.AccountObjectPrivilegeImportedPrivileges)).
	//		WithOnAccountObjectValue(tfconfig.StringVariable(databaseFromShare.ID().Name())).
	//		WithStrictPrivilegeManagement(true)
	//	resourceName := resourceModel.ResourceReference()
	//
	//	resource.Test(t, resource.TestCase{
	//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
	//			tfversion.RequireAbove(tfversion.Version1_5_0),
	//		},
	//		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
	//		CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
	//		Steps: []resource.TestStep{
	//			{
	//				Config: config.FromModels(t, providerModel, resourceModel),
	//				Check: resource.ComposeAggregateTestCheckFunc(
	//					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
	//					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeImportedPrivileges)),
	//					resource.TestCheckResourceAttr(resourceName, "strict_privilege_management", "true"),
	//				),
	//			},
	//			{
	//				// Verify no changes on re-apply (IMPORTED_PRIVILEGES/USAGE mapping should be stable)
	//				ConfigPlanChecks: resource.ConfigPlanChecks{
	//					PreApply: []plancheck.PlanCheck{
	//						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
	//					},
	//				},
	//				Config: config.FromModels(t, providerModel, resourceModel),
	//			},
	//		},
	//	})
	//})

	// Test 4: Multiple privileges with one added externally
	// Strict privilege management should detect and revoke only the external one
	t.Run("MultiplePrivilegesWithExternalAddition", func(t *testing.T) {
		role, roleCleanup := testClient().Role.CreateRole(t)
		t.Cleanup(roleCleanup)

		database, databaseCleanup := testClient().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsStrictPrivilegeManagement)
		resourceModel := model.GrantPrivilegesToAccountRole("test", role.ID().Name()).
			WithPrivileges(
				string(sdk.AccountObjectPrivilegeMonitor),
				string(sdk.AccountObjectPrivilegeUsage),
			).
			WithOnAccountObjectValue(tfconfig.StringVariable(database.ID().Name())).
			WithStrictPrivilegeManagement(true)
		resourceName := resourceModel.ResourceReference()

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			CheckDestroy:             CheckAccountRolePrivilegesRevoked(t),
			Steps: []resource.TestStep{
				{
					Config: config.FromModels(t, providerModel, resourceModel),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					),
				},
				{
					PreConfig: func() {
						// Grant CREATE_SCHEMA privilege externally
						testClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(
							t,
							role.ID(),
							database.ID(),
							[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
							false,
						)
					},
					// The external CREATE_SCHEMA grant should be detected and the plan should update
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						},
					},
					Config: config.FromModels(t, providerModel, resourceModel),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
						// Verify only the configured privileges remain, CREATE_SCHEMA was revoked
						queriedAccountRolePrivilegesEqualTo(
							t,
							role.ID(),
							string(sdk.AccountObjectPrivilegeMonitor),
							string(sdk.AccountObjectPrivilegeUsage),
						),
					),
				},
			},
		})
	})
}
