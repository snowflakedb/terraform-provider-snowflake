//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnObject_Database_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Database_IdentifiersWithDots(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")
	_, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSetWithId(t, databaseId)
	t.Cleanup(databaseCleanup)

	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Schema_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()
	schemaFullyQualifiedName := schemaId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "SCHEMA"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|SCHEMA|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeSchema, accountRoleName, schemaFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Schema_ToDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()
	schemaFullyQualifiedName := schemaId.FullyQualifiedName()

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleName := databaseRoleId.Name()
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"database_role_name": config.StringVariable(databaseRoleName),
		"database_name":      config.StringVariable(databaseName),
		"schema_name":        config.StringVariable(schemaName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToDatabaseRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "SCHEMA"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToDatabaseRole|%s||OnObject|SCHEMA|%s", databaseRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							DatabaseRole: databaseRoleId,
						},
					}, sdk.ObjectTypeSchema, databaseRoleName, schemaFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToDatabaseRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Table_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
		"table_name":        config.StringVariable(tableId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleName, tableId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Table_ToDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	tableName := tableId.Name()
	tableFullyQualifiedName := tableId.FullyQualifiedName()

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleName := databaseRoleId.Name()
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"database_role_name": config.StringVariable(databaseRoleName),
		"database_name":      config.StringVariable(databaseName),
		"schema_name":        config.StringVariable(schemaName),
		"table_name":         config.StringVariable(tableName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToDatabaseRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", tableFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToDatabaseRole|%s||OnObject|TABLE|%s", databaseRoleFullyQualifiedName, tableFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							DatabaseRole: databaseRoleId,
						},
					}, sdk.ObjectTypeTable, databaseRoleName, tableFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToDatabaseRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_ProcedureWithArguments_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	procedureId := testClient().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(testClient().Ids.Alpha(), schemaId, sdk.DataTypeFloat)
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database_name":     config.StringVariable(databaseId.Name()),
		"schema_name":       config.StringVariable(schemaId.Name()),
		"procedure_name":    config.StringVariable(procedureId.Name()),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Procedure_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "PROCEDURE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", procedureId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|PROCEDURE|%s", accountRoleId.FullyQualifiedName(), procedureId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeProcedure, accountRoleId.Name(), procedureId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Procedure_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_ProcedureWithoutArguments_ToDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	procedureId := testClient().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(testClient().Ids.Alpha(), schemaId)
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)

	configVariables := config.Variables{
		"database_role_name": config.StringVariable(databaseRoleId.Name()),
		"database_name":      config.StringVariable(databaseId.Name()),
		"schema_name":        config.StringVariable(schemaId.Name()),
		"procedure_name":     config.StringVariable(procedureId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Procedure_ToDatabaseRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "PROCEDURE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", procedureId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToDatabaseRole|%s||OnObject|PROCEDURE|%s", databaseRoleId.FullyQualifiedName(), procedureId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							DatabaseRole: databaseRoleId,
						},
					}, sdk.ObjectTypeProcedure, databaseRoleId.Name(), procedureId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Procedure_ToDatabaseRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAll_InDatabase_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	secondTableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database_name":     config.StringVariable(databaseId.Name()),
		"schema_name":       config.StringVariable(schemaId.Name()),
		"table_name":        config.StringVariable(tableId.Name()),
		"second_table_name": config.StringVariable(secondTableId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InDatabase_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.in_database", databaseId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|TABLES|InDatabase|%s", accountRoleId.FullyQualifiedName(), databaseId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleId.Name(), tableId.FullyQualifiedName(), secondTableId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InDatabase_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAll_InSchema_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	secondTableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseId.Name()),
		"schema_name":       config.StringVariable(schemaId.Name()),
		"table_name":        config.StringVariable(tableId.Name()),
		"second_table_name": config.StringVariable(secondTableId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InSchema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.in_schema", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|TABLES|InSchema|%s", accountRoleId.FullyQualifiedName(), schemaId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleName, tableId.FullyQualifiedName(), secondTableId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InSchema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnFuture_InDatabase_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InDatabase_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnFuture|TABLES|InDatabase|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						Future: sdk.Bool(true),
						In: &sdk.ShowGrantsIn{
							Database: sdk.Pointer(databaseId),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf(`"%s"."<TABLE>"`, databaseName)),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InDatabase_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnFuture_InSchema_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()
	schemaFullyQualifiedName := schemaId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InSchema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.in_schema", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnFuture|TABLES|InSchema|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						Future: sdk.Bool(true),
						In: &sdk.ShowGrantsIn{
							Schema: sdk.Pointer(schemaId),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf(`"%s"."%s"."<TABLE>"`, databaseName, schemaName)),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InSchema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_InvalidConfiguration_EmptyObjectType(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(roleId.Name()),
		"database_name":     config.StringVariable(database.ID().Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/InvalidConfiguration_EmptyObjectType"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("expected on.0.object_type to be one of"),
			},
		},
	})
}

func TestAcc_GrantOwnership_InvalidConfiguration_MultipleTargets(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(roleId.Name()),
		"database_name":     config.StringVariable(database.ID().Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/InvalidConfiguration_MultipleTargets"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("only one of `on.0.all,on.0.future,on.0.object_name`"),
			},
		},
	})
}

func TestAcc_GrantOwnership_TargetObjectRemovedOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoDatabaseResource"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				PreConfig: func() {
					currentRole := testClient().Context.CurrentRole(t)
					testClient().Grant.GrantOwnershipToAccountRole(t, currentRole, sdk.ObjectTypeDatabase, databaseId)
					databaseCleanup()
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoDatabaseResource"),
				ConfigVariables: configVariables,
				// The error occurs in Create operation indicating the Read operation couldn't find the grant and set the resource as removed.
				ExpectError: regexp.MustCompile("An error occurred during grant ownership"),
			},
		},
	})
}

func TestAcc_GrantOwnership_AccountRoleRemovedOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	accountRole, cleanupAccountRole := testClient().Role.CreateRole(t)
	t.Cleanup(cleanupAccountRole)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := accountRole.ID()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoRoleResource"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				PreConfig: func() {
					cleanupAccountRole()
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoRoleResource"),
				ConfigVariables: configVariables,
				// The error occurs in Create operation indicating the Read operation couldn't find the grant and set the resource as removed.
				ExpectError: regexp.MustCompile("An error occurred during grant ownership"),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnMaterializedView(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaName := schemaId.Name()

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	tableName := tableId.Name()
	materializedViewId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name":      config.StringVariable(accountRoleName),
		"database_name":          config.StringVariable(databaseName),
		"schema_name":            config.StringVariable(schemaName),
		"table_name":             config.StringVariable(tableName),
		"materialized_view_name": config.StringVariable(materializedViewId.Name()),
		"warehouse_name":         config.StringVariable(TestWarehouseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_MaterializedView_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "MATERIALIZED VIEW"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", materializedViewId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|MATERIALIZED VIEW|%s", accountRoleId.FullyQualifiedName(), materializedViewId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeMaterializedView, accountRoleName, materializedViewId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_MaterializedView_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_RoleBasedAccessControlUseCase(t *testing.T) {
	t.Skip("Will be un-skipped in SNOW-1313849")

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseName := database.ID().Name()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	schemaName := schemaId.Name()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	userId := testClient().Context.CurrentUser(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// We have to make it in two steps, because provider blocks cannot contain depends_on meta-argument
			// that are needed to grant the role to the current user before it can be used.
			// Additionally, only the Config field can specify a configuration with custom provider blocks.
			{
				Config: roleBasedAccessControlUseCaseConfig(accountRoleName, databaseName, userId.Name(), schemaName, false),
			},
			{
				Config: roleBasedAccessControlUseCaseConfig(accountRoleName, databaseName, userId.Name(), schemaName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func roleBasedAccessControlUseCaseConfig(accountRoleName string, databaseName string, userName string, schemaName string, withSecondaryProvider bool) string {
	baseConfig := fmt.Sprintf(`
resource "snowflake_account_role" "test" {
  name = "%[1]s"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    object_type = "DATABASE"
    object_name = "%[2]s"
  }
}

resource "snowflake_grant_account_role" "test" {
  role_name = snowflake_role.test.name
  user_name = "%[3]s"
}
`, accountRoleName, databaseName, userName)

	// TODO [SNOW-1501905]: build these configs from builders
	secondaryProviderConfig := fmt.Sprintf(`
provider "snowflake" {
  profile = "default"
  alias = "secondary"
  role = snowflake_role.test.name
}

resource "snowflake_schema" "test" {
  depends_on = [snowflake_grant_ownership.test, snowflake_grant_account_role.test]
  provider = snowflake.secondary
  database = "%[1]s"
  name     = "%[2]s"
}
`, databaseName, schemaName)

	if withSecondaryProvider {
		return fmt.Sprintf("%s\n%s", baseConfig, secondaryProviderConfig)
	}

	return baseConfig
}

func TestAcc_GrantOwnership_MoveOwnershipOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	otherAccountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	otherAccountRoleName := otherAccountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name":       config.StringVariable(accountRoleName),
		"other_account_role_name": config.StringVariable(otherAccountRoleName),
		"database_name":           config.StringVariable(databaseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/MoveResourceOwnershipOutsideTerraform"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				PreConfig: func() {
					testClient().Grant.GrantOwnershipToAccountRole(t, otherAccountRoleId, sdk.ObjectTypeDatabase, databaseId)
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/MoveResourceOwnershipOutsideTerraform"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeDatabase,
								Name:       databaseId,
							},
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_ForceOwnershipTransferOnCreate(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	newRole, newRoleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(newRoleCleanup)

	testClient().Grant.GrantOwnershipToAccountRole(t, role.ID(), sdk.ObjectTypeDatabase, database.ID())

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	newDatabaseOwningAccountRoleId := newRole.ID()
	newDatabaseOwningAccountRoleName := newDatabaseOwningAccountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(newDatabaseOwningAccountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/ForceOwnershipTransferOnCreate"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", newDatabaseOwningAccountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|\"%s\"||OnObject|DATABASE|%s", newDatabaseOwningAccountRoleName, databaseFullyQualifiedName)),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnPipe(t *testing.T) {
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()
	stageName := stageId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(pipeId.DatabaseName()),
		"schema":            config.StringVariable(pipeId.SchemaName()),
		"stage":             config.StringVariable(stageName),
		"table":             config.StringVariable(tableName),
		"pipe":              config.StringVariable(pipeId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnPipe"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", sdk.ObjectTypePipe.String()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", pipeId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|PIPE|%s", accountRoleFullyQualifiedName, pipeId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypePipe,
								Name:       pipeId,
							},
						},
					}, sdk.ObjectTypePipe, accountRoleName, pipeId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAllPipes(t *testing.T) {
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()
	stageName := stageId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()
	secondPipeId := testClient().Ids.RandomSchemaObjectIdentifier()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(pipeId.DatabaseName()),
		"schema":            config.StringVariable(pipeId.SchemaName()),
		"stage":             config.StringVariable(stageName),
		"table":             config.StringVariable(tableName),
		"pipe":              config.StringVariable(pipeId.Name()),
		"second_pipe":       config.StringVariable(secondPipeId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|PIPES|InSchema|%s", accountRoleFullyQualifiedName, testClient().Ids.SchemaId().FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypePipe, accountRoleName, pipeId.FullyQualifiedName(), secondPipeId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnTask(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database":          config.StringVariable(taskId.DatabaseName()),
		"schema":            config.StringVariable(taskId.SchemaName()),
		"task":              config.StringVariable(taskId.Name()),
		"warehouse":         config.StringVariable(TestWarehouseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", sdk.ObjectTypeTask.String()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", taskId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TASK|%s", accountRoleId.FullyQualifiedName(), taskId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnTask_Discussion2877(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database":          config.StringVariable(taskId.DatabaseName()),
		"schema":            config.StringVariable(taskId.SchemaName()),
		"task":              config.StringVariable(taskId.Name()),
		"child":             config.StringVariable(childId.Name()),
		"warehouse":         config.StringVariable(TestWarehouseName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/1"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TASK|%s", accountRoleId.FullyQualifiedName(), taskId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/2"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("cannot have the given predecessor since they do not share the same owner role"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/3"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, testClient().Context.CurrentRole(t).Name(), taskId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask_Discussion2877/4"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test", "name", taskId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "name", childId.Name()),
					resource.TestCheckResourceAttr("snowflake_task.child", "after.0", taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       childId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), childId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAllTasks(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	secondTaskId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleId.Name()),
		"database":          config.StringVariable(taskId.DatabaseName()),
		"schema":            config.StringVariable(taskId.SchemaName()),
		"task":              config.StringVariable(taskId.Name()),
		"second_task":       config.StringVariable(secondTaskId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnAllTasks"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s|REVOKE|OnAll|TASKS|InSchema|%s", accountRoleId.FullyQualifiedName(), testClient().Ids.SchemaId().FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					},
						sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName(), secondTaskId.FullyQualifiedName()),
				),
			},
		},
	})
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/3750 is fixed
func TestAcc_GrantOwnership_OnServerlessTask(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	configVariables := config.Variables{
		"account_role_name":   config.StringVariable(accountRoleId.Name()),
		"database":            config.StringVariable(taskId.DatabaseName()),
		"schema":              config.StringVariable(taskId.SchemaName()),
		"task":                config.StringVariable(taskId.Name()),
		"warehouse_init_size": config.StringVariable("XSMALL"),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnTask"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", sdk.ObjectTypeTask.String()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", taskId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TASK|%s", accountRoleId.FullyQualifiedName(), taskId.FullyQualifiedName())),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       taskId,
							},
						},
					}, sdk.ObjectTypeTask, accountRoleId.Name(), taskId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name":  config.StringVariable(accountRoleId.Name()),
		"database_name":      config.StringVariable(databaseName),
		"database_role_name": config.StringVariable(databaseRoleId.Name()),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_DatabaseRole_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE ROLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE ROLE|%s", accountRoleFullyQualifiedName, databaseRoleFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeDatabaseRole,
								Name:       databaseRoleId,
							},
						},
					}, sdk.ObjectTypeRole, accountRoleId.Name(), databaseRoleFullyQualifiedName),
				),
			},
		},
	})
}

func checkResourceOwnershipIsGranted(opts *sdk.ShowGrantOptions, grantOn sdk.ObjectType, roleName string, objectNames ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := client.Grants.Show(ctx, opts)
		if err != nil {
			return err
		}

		found := make([]string, 0)
		for _, grant := range grants {
			if grant.Privilege == "OWNERSHIP" &&
				(grant.GrantedOn == grantOn || grant.GrantOn == grantOn) &&
				grant.GranteeName.Name() == roleName &&
				slices.Contains(objectNames, grant.Name.FullyQualifiedName()) {
				found = append(found, grant.Name.FullyQualifiedName())
			}
		}

		if len(found) != len(objectNames) {
			return fmt.Errorf("unable to find ownership privilege on %s granted to %s, expected names: %v, found: %v", grantOn, roleName, objectNames, found)
		}

		return nil
	}
}

func TestAcc_GrantOwnership_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	escapedFullyQualifiedName := fmt.Sprintf(`\"%s\".\"%s\".\"%s\"`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantOwnershipOnTableBasicConfig(TestDatabaseName, TestSchemaName, tableId.Name(), accountRoleId.Name(), escapedFullyQualifiedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantOwnershipOnTableBasicConfig(TestDatabaseName, TestSchemaName, tableId.Name(), accountRoleId.Name(), escapedFullyQualifiedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantOwnershipOnTableBasicConfig(databaseName string, schemaName string, tableName string, roleName string, fullTableName string) string {
	return fmt.Sprintf(`
resource "snowflake_account_role" "test" {
	name = "%[4]s"
}

resource "snowflake_table" "test" {
	name     = "%[3]s"
	database = "%[1]s"
	schema   = "%[2]s"

	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_grant_ownership" "test" {
	depends_on = [snowflake_table.test]
	account_role_name = snowflake_account_role.test.name
	on {
		object_type = "TABLE"
		object_name = "%[5]s"
	}
}
`, databaseName, schemaName, tableName, roleName, fullTableName)
}

func TestAcc_GrantOwnership_IdentifierQuotingDiffSuppression(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	unescapedFullyQualifiedName := fmt.Sprintf(`%s.%s.%s`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            grantOwnershipOnTableBasicConfigWithManagedDatabaseAndSchema(databaseId.Name(), schemaId.Name(), tableId.Name(), accountRoleId.Name(), unescapedFullyQualifiedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", unescapedFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantOwnershipOnTableBasicConfigWithManagedDatabaseAndSchema(databaseId.Name(), schemaId.Name(), tableId.Name(), accountRoleId.Name(), unescapedFullyQualifiedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", unescapedFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantOwnershipOnTableBasicConfigWithManagedDatabaseAndSchema(databaseName string, schemaName string, tableName string, roleName string, fullTableName string) string {
	return fmt.Sprintf(`
resource "snowflake_account_role" "test" {
	name = "%[4]s"
}

resource "snowflake_schema" "test" {
	database = "%[1]s"
	name = "%[2]s"
}

resource "snowflake_table" "test" {
	name     = "%[3]s"
	database = "%[1]s"
	schema   = snowflake_schema.test.name

	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_grant_ownership" "test" {
	depends_on = [snowflake_table.test]
	account_role_name = snowflake_account_role.test.name
	on {
		object_type = "TABLE"
		object_name = "%[5]s"
	}
}
`, databaseName, schemaName, tableName, roleName, fullTableName)
}

// confirms addition of resource monitor as part of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3318
func TestAcc_GrantOwnership_OnObject_ResourceMonitor_ToAccountRole(t *testing.T) {
	resourceMonitorId := testClient().Ids.RandomAccountObjectIdentifier()
	resourceMonitorName := resourceMonitorId.Name()
	resourceMonitorIdFullyQualifiedName := resourceMonitorId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name":     config.StringVariable(accountRoleName),
		"resource_monitor_name": config.StringVariable(resourceMonitorName),
	}

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_ResourceMonitor_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "RESOURCE MONITOR"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", resourceMonitorName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|RESOURCE MONITOR|%s", accountRoleFullyQualifiedName, resourceMonitorIdFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeResourceMonitor, accountRoleName, resourceMonitorIdFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_ResourceMonitor_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This test proves that managing grants on HYBRID TABLE is not supported in Snowflake. TABLE should be used instead.
func TestAcc_GrantOwnership_OnObject_HybridTable_ToAccountRole_Fails(t *testing.T) {
	hybridTableId, hybridTableCleanup := testClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	configVariables := func(objectType sdk.ObjectType) config.Variables {
		cfg := config.Variables{
			"account_role_name":                 config.StringVariable(accountRoleName),
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
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_HybridTable_ToAccountRole"),
				ConfigVariables: configVariables(sdk.ObjectTypeHybridTable),
				ExpectError:     regexp.MustCompile("syntax error line 1 at position 26 unexpected 'TABLE"),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_HybridTable_ToAccountRole"),
				ConfigVariables: configVariables(sdk.ObjectTypeTable),
			},
		},
	})
}
