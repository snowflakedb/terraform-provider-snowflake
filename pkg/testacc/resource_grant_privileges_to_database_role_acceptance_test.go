//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeApplyBudget, sdk.AccountObjectPrivilegeCreateSchema, sdk.AccountObjectPrivilegeModify, sdk.AccountObjectPrivilegeUsage).
		WithOnDatabase(testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(true)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeUsage, sdk.AccountObjectPrivilegeModify, sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(true)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaPrivileges(sdk.SchemaPrivilegeCreateTable, sdk.SchemaPrivilegeModify).
		WithOnSchemaName(schemaId.FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: `resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "some_database.role_name"
  privileges         = ["USAGE"]
  on_schema {
    schema_name             = "some_database.schema_name"
    all_schemas_in_database = "some_database"
  }
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnAllSchemasInDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaPrivileges(sdk.SchemaPrivilegeCreateTable, sdk.SchemaPrivilegeModify).
		WithOnAllSchemasInDatabase(testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaPrivileges(sdk.SchemaPrivilegeCreateTable, sdk.SchemaPrivilegeModify).
		WithOnFutureSchemasInDatabase(testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber_38_0},
	})
	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeInsert, sdk.SchemaObjectPrivilegeUpdate).
		WithOnSchemaObjectObject(string(sdk.ObjectTypeTable), tableId.FullyQualifiedName()).
		WithWithGrantOption(false).
		WithDependsOn(tableModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, grantModel),
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
				Config:            accconfig.FromModels(t, tableModel, grantModel),
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
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: `resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"some_database\".\"some_name\""
  privileges         = ["OWNERSHIP"]
  with_grant_option  = false
  on_schema_object {
    object_type = "TABLE"
    object_name = "some_database.some_schema.some_table"
  }
}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeInsert, sdk.SchemaObjectPrivilegeUpdate).
		WithOnSchemaObjectAllInDatabase(sdk.PluralObjectTypeTables.String(), testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeMonitor).
		WithOnSchemaObjectAllInDatabase(sdk.PluralObjectTypePipes.String(), testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeInsert, sdk.SchemaObjectPrivilegeUpdate).
		WithOnSchemaObjectFutureInDatabase(sdk.PluralObjectTypeTables.String(), testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeUsage).
		WithOnSchemaObjectFutureInDatabase(sdk.PluralObjectTypeStreamlits.String(), testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", testClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnFuture|STREAMLITS|InDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeUsage).
		WithOnSchemaObjectAllInDatabase(sdk.PluralObjectTypeStreamlits.String(), testClient().Ids.DatabaseId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeUsage).
		WithOnSchemaObjectObject(string(sdk.ObjectTypeFunction), function.ID().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeUsage).
		WithOnSchemaObjectObject(string(sdk.ObjectTypeFunction), function.ID().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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
				Config:            accconfig.FromModels(t, grantModel),
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

	grantModel := func(allPrivileges bool, privileges ...sdk.AccountObjectPrivilege) *model.GrantPrivilegesToDatabaseRoleModel {
		m := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
			WithOnDatabase(testClient().Ids.DatabaseId().FullyQualifiedName())
		if allPrivileges {
			return m.WithAllPrivileges(true)
		}
		return m.WithAccountObjectPrivileges(privileges...)
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel(false, sdk.AccountObjectPrivilegeCreateSchema, sdk.AccountObjectPrivilegeModify)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				Config: accconfig.FromModels(t, grantModel(false, sdk.AccountObjectPrivilegeCreateSchema, sdk.AccountObjectPrivilegeMonitor, sdk.AccountObjectPrivilegeUsage)),
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
				Config: accconfig.FromModels(t, grantModel(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				Config: accconfig.FromModels(t, grantModel(false, sdk.AccountObjectPrivilegeModify, sdk.AccountObjectPrivilegeMonitor)),
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

	grantOnDatabaseModel := func(allPrivileges bool, privileges ...string) *model.GrantPrivilegesToDatabaseRoleModel {
		m := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
			WithOnDatabase(testClient().Ids.DatabaseId().FullyQualifiedName())
		if allPrivileges {
			return m.WithAllPrivileges(true)
		}
		return m.WithPrivileges(privileges...)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantOnDatabaseModel(false,
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				)),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					t,
					databaseRole.ID(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				Config: accconfig.FromModels(t, grantOnDatabaseModel(true)),
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
				Config: accconfig.FromModels(t, grantOnDatabaseModel(false,
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				)),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					t,
					databaseRole.ID(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				Config: accconfig.FromModels(t,
					model.Schema("test", schemaId.DatabaseName(), schemaId.Name()),
					model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
						WithPrivileges(sdk.SchemaPrivilegeCreateTask.String(), sdk.SchemaPrivilegeCreateExternalTable.String()).
						WithOnSchemaName(fmt.Sprintf("%s.%s", schemaId.DatabaseName(), schemaId.Name())).
						WithDependsOn(model.Schema("test", schemaId.DatabaseName(), schemaId.Name()).ResourceReference()),
				),
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

	unquotedName := fmt.Sprintf("%s.%s", databaseRole.ID().DatabaseName(), databaseRole.ID().Name())

	grantModel := func(alwaysApply bool) *model.GrantPrivilegesToDatabaseRoleModel {
		return model.GrantPrivilegesToDatabaseRole("test", unquotedName).
			WithAllPrivileges(true).
			WithOnDatabase(TestDatabaseName).
			WithAlwaysApply(alwaysApply)
	}

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel(false)),
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
				Config: accconfig.FromModels(t, grantModel(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: accconfig.FromModels(t, grantModel(true)),
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
				Config: accconfig.FromModels(t, grantModel(true)),
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
				Config: accconfig.FromModels(t, grantModel(false)),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaPrivileges(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection, sdk.SchemaPrivilegeCreateSnowflakeMlForecast).
		WithOnSchemaName(testClient().Ids.SchemaId().FullyQualifiedName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(databaseRole.ID().DatabaseName()).
		WithWithGrantOption(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
				Config: accconfig.FromModels(t, grantModel),
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
				Config: accconfig.FromModels(t, grantModel),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(databaseRole.ID().DatabaseName()).
		WithWithGrantOption(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
				Config: accconfig.FromModels(t, grantModel),
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
				Config: accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(databaseRole.ID().DatabaseName()).
		WithWithGrantOption(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
			},
			{
				PreConfig: func() { databaseCleanup() },
				Config:    accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(databaseRole.ID().DatabaseName()).
		WithWithGrantOption(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
			},
			{
				PreConfig: func() { databaseRoleCleanup() },
				Config:    accconfig.FromModels(t, grantModel),
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

	unquotedName := fmt.Sprintf("%s.%s", databaseRole.ID().DatabaseName(), databaseRole.ID().Name())

	grantModel := model.GrantPrivilegesToDatabaseRole("test", unquotedName).
		WithAllPrivileges(true).
		WithOnDatabase(TestDatabaseName).
		WithAlwaysApply(true)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config:             accconfig.FromModels(t, grantModel),
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

	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithSchemaPrivileges(sdk.SchemaPrivilegeCreateNotebook).
		WithOnAllSchemasInDatabase(databaseRole.ID().DatabaseName()).
		WithWithGrantOption(false)

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
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

	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + grantPrivilegesToDatabaseRoleBasicConfig(quotedDatabaseRoleId, quotedSchemaId),
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

	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + grantPrivilegesToDatabaseRoleBasicConfig(unquotedDatabaseRoleId, unquotedSchemaId),
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
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            providerConfig + grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseRoleId.DatabaseName()),
				// Previously, we expected a non-empty plan, because Snowflake returned MODULE instead of MODEL in SHOW FUTURE GRANTS.
				// Now, this behavior is fixed in Snowflake, and the plan is empty.
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModels, databaseRoleId.DatabaseName()),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnFutureModelMonitors_InDatabase_v2_17_0_NonEmptyPlan(t *testing.T) {
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders:  ExternalProviderWithExactVersion("2.17.0"),
				Config:             grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModelMonitors, databaseRoleId.DatabaseName()),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToDatabaseRoleOnFutureInDatabaseConfig(databaseRoleId, []string{"USAGE"}, sdk.PluralObjectTypeModelMonitors, databaseRoleId.DatabaseName()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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

	grantModel := func(objectType sdk.ObjectType) *model.GrantPrivilegesToDatabaseRoleModel {
		return model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
			WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeApplyBudget).
			WithOnSchemaObjectObject(string(objectType), hybridTableId.FullyQualifiedName())
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, grantModel(sdk.ObjectTypeHybridTable)),
				ExpectError: regexp.MustCompile("Unsupported feature"),
			},
			{
				Config: accconfig.FromModels(t, grantModel(sdk.ObjectTypeTable)),
			},
		},
	})
}

// proves that https://github.com/snowflakedb/terraform-provider-snowflake/issues/3690 is fixed
func TestAcc_GrantPrivileges_ToDatabaseRole_WithEmptyPrivileges(t *testing.T) {
	databaseRole, databaseRoleCleanup := testClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	grantModelWithPrivileges := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithAccountObjectPrivileges(sdk.AccountObjectPrivilegeUsage, sdk.AccountObjectPrivilegeCreateSchema).
		WithOnDatabase(testClient().Ids.DatabaseId().Name())

	grantModelEmpty := model.GrantPrivilegesToDatabaseRole("test", databaseRole.ID().FullyQualifiedName()).
		WithPrivilegesValue(tfconfig.ListVariable()).
		WithOnDatabase(testClient().Ids.DatabaseId().Name())

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModelWithPrivileges),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_database", testClient().Ids.DatabaseId().Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			// {
			//	ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
			//	Config:            grantPrivilegesToDatabaseRole3690Config(databaseRole.ID()),
			//	ExpectError:       regexp.MustCompile("Error: Failed to parse internal identifier"),
			// },
			//
			// The step above fails with:
			// │ Error: Failed to parse internal identifier
			// ...
			// │ Error: [grant_privileges_to_database_role_identifier.go:79] invalid Privileges value: , should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all
			// │ privileges
			//
			// and affects the next test steps
			{
				Config:      accconfig.FromModels(t, grantModelEmpty),
				ExpectError: regexp.MustCompile("Error: Not enough list items"),
			},
			{
				Config: accconfig.FromModels(t, grantModelWithPrivileges),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRole.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_database", testClient().Ids.DatabaseId().Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE|OnDatabase|%s", databaseRole.ID().FullyQualifiedName(), testClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}
