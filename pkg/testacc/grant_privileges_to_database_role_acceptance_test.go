//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeApplyBudget)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeModify)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeApplyBudget)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.3", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|APPLYBUDGET,CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase_PrivilegesReversed(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeModify)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"with_grant_option": config.BoolVariable(true),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchema(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	schemaId := testClient().Ids.SchemaId()
	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(schemaId.DatabaseName()),
		"schema":            config.StringVariable(schemaId.Name()),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnSchema|%s", databaseRole.ID().FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchema_ExactlyOneOf(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema_ExactlyOneOf"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnAllSchemasInDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnAllSchemasInDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnFutureSchemasInDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnFutureSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.future_schemas_in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnFutureSchemasInDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnFutureSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnObject(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	databaseRoleName := databaseRole.ID().Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := config.Variables{
		"name":       config.StringVariable(databaseRoleName),
		"table_name": config.StringVariable(tableId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          config.StringVariable(tableId.DatabaseName()),
		"schema":            config.StringVariable(tableId.SchemaName()),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnObject|TABLE|%s", databaseRole.ID().FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnObject_OwnershipPrivilege(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject_OwnershipPrivilege"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnAll|TABLES|InDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAllPipes(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypePipes)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MONITOR|OnSchemaObject|OnAll|PIPES|InDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFuture_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnFuture|TABLES|InDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO [SNOW-1272222]: fix the test when it starts working on Snowflake side
func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	t.Skip("Fix after it starts working on Snowflake side, reference: SNOW-1272222")

	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported feature 'STREAMLIT'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnAll|STREAMLITS|InDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFunctionWithArguments(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	function := testClient().Function.CreateSecure(t, sdk.DataTypeFloat)
	configVariables := config.Variables{
		"name":          config.StringVariable(databaseRole.ID().FullyQualifiedName()),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          config.StringVariable(function.ID().DatabaseName()),
		"schema":            config.StringVariable(function.ID().SchemaName()),
		"with_grant_option": config.BoolVariable(false),
		"argument_type":     config.StringVariable(string(sdk.DataTypeFloat)),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeFunction)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|%s", databaseRole.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFunctionWithoutArguments(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	function := testClient().Function.CreateSecure(t)
	configVariables := config.Variables{
		"name":          config.StringVariable(databaseRole.ID().FullyQualifiedName()),
		"function_name": config.StringVariable(function.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":          config.StringVariable(function.ID().DatabaseName()),
		"schema":            config.StringVariable(function.ID().SchemaName()),
		"with_grant_option": config.BoolVariable(false),
		"argument_type":     config.StringVariable(""),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFunction"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeFunction)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|%s", databaseRole.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFunction"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(databaseRole.ID().Name()),
			"database": config.StringVariable(databaseRole.ID().DatabaseName()),
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

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeModify,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeModify,
					sdk.AccountObjectPrivilegeMonitor,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MODIFY,MONITOR|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges_SnowflakeChecked(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(allPrivileges bool, privileges []string, schemaName string) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(databaseRole.ID().Name()),
			"database": config.StringVariable(databaseRole.ID().DatabaseName()),
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				}, ""),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					t,
					databaseRole.ID(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
				ConfigVariables: configVariables(true, []string{}, ""),
				Check: queriedPrivilegesToDatabaseRoleContainAtLeast(
					t,
					databaseRole.ID(),
					sdk.AccountObjectPrivilegeCreateDatabaseRole.String(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
					sdk.AccountObjectPrivilegeUsage.String(),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				}, ""),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					t,
					databaseRole.ID(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
				ConfigVariables: configVariables(false, []string{
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				}, schemaId.Name()),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					t,
					databaseRole.ID(),
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_AlwaysApply(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(databaseRole.ID().Name()),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(TestDatabaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

// proved https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2651
func TestAcc_GrantPrivilegesToDatabaseRole_MLPrivileges(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
		),
		"database":          config.StringVariable(TestDatabaseName),
		"schema":            config.StringVariable(TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", testClient().Ids.SchemaId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SNOWFLAKE.ML.ANOMALY_DETECTION,CREATE SNOWFLAKE.ML.FORECAST|OnSchema|OnSchema|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.SchemaId().FullyQualifiedName())),
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
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithGrantOptions(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(databaseRole.ID().DatabaseName()),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
						t, databaseRole.ID(),
						testClient().Ids.DatabaseId(),
						[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
						false,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(databaseRole.ID().DatabaseName()),
		"with_grant_option": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
						t, databaseRole.ID(),
						testClient().Ids.DatabaseId(),
						[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
						true,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToDatabaseRole_RemoveGrantedObjectOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name":     config.StringVariable(databaseRole.ID().Name()),
		"database": config.StringVariable(databaseRole.ID().DatabaseName()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseCleanup() },
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to database role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToDatabaseRole_RemoveDatabaseRoleOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name":     config.StringVariable(databaseRole.ID().Name()),
		"database": config.StringVariable(databaseRole.ID().DatabaseName()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseRoleCleanup() },
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to database role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2689 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_AlwaysApply_SetAfterCreate(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(databaseRole.ID().Name()),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(TestDatabaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory:    ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables:    configVariables(true),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2960
func TestAcc_GrantPrivilegesToDatabaseRole_CreateNotebooks(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateNotebook)),
		),
		"database":          config.StringVariable(databaseRole.ID().DatabaseName()),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateNotebook)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE NOTEBOOK|OnSchema|OnAllSchemasInDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

// TODO [SNOW-1431726]: Move to helpers
func queriedPrivilegesToDatabaseRoleEqualTo(t *testing.T, databaseRoleName sdk.DatabaseObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesEqualTo(func() ([]sdk.Grant, error) {
		return testClient().Grant.ShowGrantsToDatabaseRole(t, databaseRoleName)
	}, privileges...)
}

func queriedPrivilegesToDatabaseRoleContainAtLeast(t *testing.T, databaseRoleName sdk.DatabaseObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	t.Helper()
	return queriedPrivilegesContainAtLeast(func() ([]sdk.Grant, error) {
		return testClient().Grant.ShowGrantsToDatabaseRole(t, databaseRoleName)
	}, databaseRoleName, privileges...)
}

func revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
	t *testing.T,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	client := testClient()

	client.Grant.RevokePrivilegesOnDatabaseFromDatabaseRole(t, databaseRoleId, databaseId, privileges)
	client.Grant.GrantPrivilegesOnDatabaseToDatabaseRole(t, databaseRoleId, databaseId, privileges, withGrantOption)
}

func TestAcc_GrantPrivilegesToDatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	databaseRoleId := databaseRole.ID()
	quotedDatabaseRoleId := fmt.Sprintf(`\"%s\".\"%s\"`, databaseRoleId.DatabaseName(), databaseRoleId.Name())

	schemaId := testClient().Ids.SchemaId()
	quotedSchemaId := fmt.Sprintf(`\"%s\".\"%s\"`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToDatabaseRoleBasicConfig(quotedDatabaseRoleId, quotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToDatabaseRoleBasicConfig(quotedDatabaseRoleId, quotedSchemaId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToDatabaseRoleBasicConfig(fullyQualifiedDatabaseRoleName string, fullyQualifiedSchemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "%[1]s"
  privileges         = ["USAGE"]

  on_schema {
    schema_name = "%[2]s"
  }
}
`, fullyQualifiedDatabaseRoleName, fullyQualifiedSchemaName)
}

func TestAcc_GrantPrivilegesToDatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	databaseRoleId := databaseRole.ID()
	unquotedDatabaseRoleId := fmt.Sprintf(`%s.%s`, databaseRoleId.DatabaseName(), databaseRoleId.Name())

	schemaId := testClient().Ids.SchemaId()
	unquotedSchemaId := fmt.Sprintf(`%s.%s`, schemaId.DatabaseName(), schemaId.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantPrivilegesToDatabaseRoleBasicConfig(unquotedDatabaseRoleId, unquotedSchemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "database_role_name", unquotedDatabaseRoleId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "on_schema.0.schema_name", unquotedSchemaId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToDatabaseRoleBasicConfig(unquotedDatabaseRoleId, unquotedSchemaId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "database_role_name", unquotedDatabaseRoleId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "on_schema.0.schema_name", unquotedSchemaId),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_database_role.test", "id", fmt.Sprintf("%s|false|false|USAGE|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), schemaId.FullyQualifiedName())),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3050
func TestAcc_GrantPrivilegesToDatabaseRole_OnFutureModels_issue3050(t *testing.T) {
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.95.0"),
				Config:             grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseRoleId.DatabaseName()),
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseRoleId.DatabaseName()),
			},
		},
	})
}

func grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId sdk.DatabaseObjectIdentifier, privileges []string, objectTypePlural sdk.PluralObjectType, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_database_role" "test" {
	name = "%[1]s"
	database = "%[2]s"
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  privileges        = [ %[3]s ]

  on_schema_object {
    future {
      object_type_plural = "%[4]s"
      in_database        = "%[5]s"
    }
  }
}
`, databaseRoleId.Name(), databaseRoleId.DatabaseName(), strings.Join(collections.Map(privileges, strconv.Quote), ","), objectTypePlural, databaseName)
}

// This test proves that managing grants on HYBRID TABLE is not supported in Snowflake. TABLE should be used instead.
func TestAcc_GrantPrivileges_OnObject_HybridTable_ToDatabaseRole_Fails(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	hybridTableId, hybridTableCleanup := testClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	configVariables := func(objectType sdk.ObjectType) config.Variables {
		cfg := config.Variables{
			"database_role_name": config.StringVariable(databaseRole.ID().FullyQualifiedName()),
			"privileges": config.ListVariable(
				config.StringVariable(string(sdk.SchemaObjectPrivilegeApplyBudget)),
			),
			"hybrid_table_fully_qualified_name": config.StringVariable(hybridTableId.FullyQualifiedName()),
			"object_type":                       config.StringVariable(string(objectType)),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeHybridTable),
				ExpectError:     regexp.MustCompile("syntax error line 1 at position 28 unexpected 'TABLE"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject_HybridTable"),
				ConfigVariables: configVariables(sdk.ObjectTypeTable),
			},
		},
	})
}

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3690 is fixed
func TestAcc_GrantPrivileges_ToDatabaseRole_WithEmptyPrivileges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRole, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRole.ID().Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			// Previously, this would throw:
			// │ Error: Failed to parse internal identifier
			// ...
			// │ Error: [grant_privileges_to_database_role_identifier.go:79] invalid Privileges value: , should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all
			// │ privileges
			//
			// After that, the state file would be corrupted and only manual state manipulation would be able to fix it.
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabaseWithEmptyPrivileges"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectError: regexp.MustCompile("Error: Not enough list items"),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}
