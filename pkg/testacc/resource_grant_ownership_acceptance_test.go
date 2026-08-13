//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_BasicUseCase_OnObject_Database_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	roleModel := model.AccountRole("test", accountRoleName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName).
		WithDependsOn(roleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, roleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
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
				Config:            accconfig.FromModels(t, roleModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_Regression_IdentifiersWithDots(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")
	_, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSetWithId(t, databaseId)
	t.Cleanup(databaseCleanup)

	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifierContaining(".")
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	roleModel := model.AccountRole("test", accountRoleName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName).
		WithDependsOn(roleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, roleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
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
				Config:            accconfig.FromModels(t, roleModel, grantModel),
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

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaFullyQualifiedName := schemaId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	roleModel := model.AccountRole("test", accountRoleName)
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeSchema, schemaFullyQualifiedName).
		WithDependsOn(roleModel.ResourceReference(), schemaModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, roleModel, schemaModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeSchema)),
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
				Config:            accconfig.FromModels(t, roleModel, schemaModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnObject_Schema_ToDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	schemaFullyQualifiedName := schemaId.FullyQualifiedName()

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleName := databaseRoleId.Name()
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	dbRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithDatabaseRoleName(databaseRoleFullyQualifiedName).
		WithOnObject(sdk.ObjectTypeSchema, schemaFullyQualifiedName).
		WithDependsOn(dbRoleModel.ResourceReference(), schemaModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, dbRoleModel, schemaModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeSchema)),
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
				Config:            accconfig.FromModels(t, dbRoleModel, schemaModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnObject_Table_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	tableModel := model.Table("test", databaseId.Name(), schemaId.Name(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeTable, tableId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), tableModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeTable)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, grantModel),
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
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	tableFullyQualifiedName := tableId.FullyQualifiedName()

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleName := databaseRoleId.Name()
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	dbRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	tableModel := model.Table("test", databaseId.Name(), schemaId.Name(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithDatabaseRoleName(databaseRoleFullyQualifiedName).
		WithOnObject(sdk.ObjectTypeTable, tableFullyQualifiedName).
		WithDependsOn(dbRoleModel.ResourceReference(), tableModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, dbRoleModel, schemaModel, tableModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeTable)),
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
				Config:            accconfig.FromModels(t, dbRoleModel, schemaModel, tableModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnObject_ProcedureWithArguments_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	procedureId := testClient().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(testClient().Ids.Alpha(), schemaId, sdk.DataTypeFloat)
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	procedureModel := model.ProcedureJavascriptBasicInline("test", procedureId, testdatatypes.DataTypeFloat, "var X=1\nreturn X").
		WithArgument("ARG1", testdatatypes.DataTypeFloat).
		WithExecuteAs("CALLER").
		WithNullInputBehavior(string(sdk.NullInputBehaviorReturnsNullInput)).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeProcedure, procedureId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), procedureModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, procedureModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeProcedure)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, procedureModel, grantModel),
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

	databaseRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	procedureModel := model.ProcedureJavascriptBasicInline("test", procedureId, testdatatypes.DataTypeFloat, "var X=1\nreturn X").
		WithExecuteAs("CALLER").
		WithNullInputBehavior(string(sdk.NullInputBehaviorReturnsNullInput)).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithDatabaseRoleName(databaseRoleId.FullyQualifiedName()).
		WithOnObject(sdk.ObjectTypeProcedure, procedureId.FullyQualifiedName()).
		WithDependsOn(databaseRoleModel.ResourceReference(), procedureModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseRoleModel, schemaModel, procedureModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeProcedure)),
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
				Config:            accconfig.FromModels(t, databaseRoleModel, schemaModel, procedureModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnAll_InDatabase_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	secondTableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	tableModel := model.Table("test", databaseId.Name(), schemaId.Name(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	table2Model := model.Table("test2", databaseId.Name(), schemaId.Name(), secondTableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnAllInDatabase(sdk.PluralObjectTypeTables, databaseId.Name()).
		WithDependsOn(tableModel.ResourceReference(), table2Model.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, table2Model, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, table2Model, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnAll_InSchema_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	secondTableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	tableModel := model.Table("test", databaseId.Name(), schemaId.Name(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	table2Model := model.Table("test2", databaseId.Name(), schemaId.Name(), secondTableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnAllInSchema(sdk.PluralObjectTypeTables, schemaId.FullyQualifiedName()).
		WithDependsOn(tableModel.ResourceReference(), table2Model.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, table2Model, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, table2Model, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnFuture_InDatabase_ToAccountRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnFutureInDatabase(sdk.PluralObjectTypeTables, databaseName).
		WithDependsOn(accountRoleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_BasicUseCase_OnFuture_InSchema_ToAccountRole(t *testing.T) {
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

	accountRoleModel := model.AccountRole("test", accountRoleName)
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnFutureInSchema(sdk.PluralObjectTypeTables, schemaFullyQualifiedName).
		WithDependsOn(accountRoleModel.ResourceReference(), schemaModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_Validations_EmptyObjectType(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleId := testClient().Ids.RandomAccountObjectIdentifier()

	accountRoleModel := model.AccountRole("test", roleId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(roleId.Name()).
		WithOnObject(sdk.ObjectType(""), database.ID().Name()).
		WithDependsOn(accountRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, accountRoleModel, grantModel),
				ExpectError: regexp.MustCompile("expected on.0.object_type to be one of"),
			},
		},
	})
}

func TestAcc_GrantOwnership_Validations_MultipleTargets(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	roleId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseName := database.ID().Name()

	accountRoleModel := model.AccountRole("test", roleId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(roleId.Name()).
		WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type": tfconfig.StringVariable(string(sdk.ObjectTypeDatabase)),
			"object_name": tfconfig.StringVariable(databaseName),
			"all": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type_plural": tfconfig.StringVariable(string(sdk.PluralObjectTypeTables)),
				"in_database":        tfconfig.StringVariable(databaseName),
			})),
			"future": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type_plural": tfconfig.StringVariable(string(sdk.PluralObjectTypeTables)),
				"in_database":        tfconfig.StringVariable(databaseName),
			})),
		})).
		WithDependsOn(accountRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, accountRoleModel, grantModel),
				ExpectError: regexp.MustCompile("only one of `on.0.all,on.0.future,on.0.object_name`"),
			},
		},
	})
}

func TestAcc_GrantOwnership_Regression_TargetObjectRemovedOutsideTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName).
		WithDependsOn(accountRoleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
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
				Config: accconfig.FromModels(t, accountRoleModel, grantModel),
				// The error occurs in Create operation indicating the Read operation couldn't find the grant and set the resource as removed.
				ExpectError: regexp.MustCompile("An error occurred during grant ownership"),
			},
		},
	})
}

func TestAcc_GrantOwnership_Regression_AccountRoleRemovedOutsideTerraform(t *testing.T) {
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

	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName)

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
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
				Config: accconfig.FromModels(t, grantModel),
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
	materializedViewId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	schemaModel := model.Schema("test", databaseName, schemaName)
	tableModel := model.Table("test", databaseName, schemaName, tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	materializedViewModel := model.MaterializedViewWithId("test", materializedViewId, fmt.Sprintf("select * from %s", tableId.FullyQualifiedName()), TestWarehouseName).
		WithDependsOn(tableModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeMaterializedView, materializedViewId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), materializedViewModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, materializedViewModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeMaterializedView)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, materializedViewModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_Regression_RoleBasedAccessControl(t *testing.T) {
	t.Skip("Will be un-skipped in SNOW-1313849")

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseName := database.ID().Name()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	schemaName := schemaId.Name()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	userId := testClient().Context.CurrentUser(t)

	accountRoleModel := model.AccountRole("test", accountRoleName)
	grantOwnershipModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName).
		WithDependsOn(accountRoleModel.ResourceReference())
	grantAccountRoleModel := model.GrantAccountRole("test", accountRoleName).
		WithUserName(userId.Name()).
		WithDependsOn(accountRoleModel.ResourceReference())

	// Resource-level `provider = snowflake.secondary` is not yet supported on models (see TODO [SNOW-1501905] in resource_model.go),
	// so the schema managed by the secondary provider is built from a small raw HCL block.
	secondaryProviderModel := providermodel.SnowflakeProviderAlias("secondary").
		WithProfile("default").
		WithRole(accountRoleName)
	secondarySchemaConfig := fmt.Sprintf(`
resource "snowflake_schema" "test" {
  depends_on = [snowflake_grant_ownership.test, snowflake_grant_account_role.test]
  provider   = snowflake.secondary
  database   = "%[1]s"
  name       = "%[2]s"
}
`, databaseName, schemaName)

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
				Config: accconfig.FromModels(t, accountRoleModel, grantOwnershipModel, grantAccountRoleModel),
			},
			{
				Config: accconfig.FromModels(t, accountRoleModel, grantOwnershipModel, grantAccountRoleModel, secondaryProviderModel) + secondarySchemaConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
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

	accountRoleModel := model.AccountRole("test", accountRoleName)
	otherRoleModel := model.AccountRole("other_role", otherAccountRoleName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName).
		WithDependsOn(accountRoleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, otherRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				PreConfig: func() {
					testClient().Grant.GrantOwnershipToAccountRole(t, otherAccountRoleId, sdk.ObjectTypeDatabase, databaseId)
				},
				Config: accconfig.FromModels(t, accountRoleModel, otherRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
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

func TestAcc_GrantOwnership_CompleteUseCase_ForceOwnershipTransferOnCreate(t *testing.T) {
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

	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(newDatabaseOwningAccountRoleName).
		WithOnObject(sdk.ObjectTypeDatabase, databaseName)

	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", newDatabaseOwningAccountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|\"%s\"||OnObject|DATABASE|%s", newDatabaseOwningAccountRoleName, databaseFullyQualifiedName)),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnPipe(t *testing.T) {
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	stageModel := model.Stage("test", stageId.DatabaseName(), stageId.SchemaName(), stageId.Name())
	tableModel := model.Table("test", tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	pipeModel := model.PipeWithIdCopyFromStageIntoTable("test", pipeId, tableModel.ResourceReference(), stageModel.ResourceReference()).
		WithDependsOn(tableModel.ResourceReference(), stageModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypePipe, pipeId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), pipeModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, stageModel, tableModel, pipeModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypePipe)),
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
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()
	secondPipeId := testClient().Ids.RandomSchemaObjectIdentifier()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	stageModel := model.Stage("test", stageId.DatabaseName(), stageId.SchemaName(), stageId.Name())
	tableModel := model.Table("test", tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	pipeModel := model.PipeWithIdCopyFromStageIntoTable("test", pipeId, tableModel.ResourceReference(), stageModel.ResourceReference()).
		WithDependsOn(tableModel.ResourceReference(), stageModel.ResourceReference())
	secondPipeModel := model.PipeWithIdCopyFromStageIntoTable("second_test", secondPipeId, tableModel.ResourceReference(), stageModel.ResourceReference()).
		WithDependsOn(tableModel.ResourceReference(), stageModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnAllInSchema(sdk.PluralObjectTypePipes, testClient().Ids.SchemaId().FullyQualifiedName()).
		WithDependsOn(pipeModel.ResourceReference(), secondPipeModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, stageModel, tableModel, pipeModel, secondPipeModel, grantModel),
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

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	taskModel := model.TaskWithId("test", taskId, false, "SELECT CURRENT_TIMESTAMP").
		WithWarehouse(TestWarehouseName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTask, taskId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), taskModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, taskModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeTask)),
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

func TestAcc_GrantOwnership_OnAllTasks(t *testing.T) {
	taskId := testClient().Ids.RandomSchemaObjectIdentifier()
	secondTaskId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	taskModel := model.TaskWithId("test", taskId, false, "SELECT CURRENT_TIMESTAMP")
	secondTaskModel := model.TaskWithId("second_test", secondTaskId, false, "SELECT CURRENT_TIMESTAMP")
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnAllInSchema(sdk.PluralObjectTypeTasks, testClient().Ids.SchemaId().FullyQualifiedName()).
		WithOutboundPrivileges("REVOKE").
		WithDependsOn(taskModel.ResourceReference(), secondTaskModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, taskModel, secondTaskModel, grantModel),
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

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	taskModel := model.TaskWithId("test", taskId, false, "SELECT CURRENT_TIMESTAMP").
		WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTask, taskId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), taskModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, taskModel, grantModel),
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

	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	dbRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeDatabaseRole, databaseRoleFullyQualifiedName).
		WithDependsOn(accountRoleModel.ResourceReference(), dbRoleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, dbRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabaseRole)),
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

func TestAcc_GrantOwnership_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	tableModel := model.Table("test", tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTable, tableId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference(), tableModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + accconfig.FromModels(t, accountRoleModel, tableModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, accountRoleModel, tableModel, grantModel),
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

func TestAcc_GrantOwnership_IdentifierQuotingDiffSuppression(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	unescapedFullyQualifiedName := fmt.Sprintf(`%s.%s.%s`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	schemaModel := model.Schema("test", databaseId.Name(), schemaId.Name())
	tableModel := model.Table("test", databaseId.Name(), schemaId.Name(), tableId.Name(), []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithDependsOn(schemaModel.ResourceReference())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeTable, unescapedFullyQualifiedName).
		WithDependsOn(accountRoleModel.ResourceReference(), tableModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", unescapedFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, accountRoleModel, schemaModel, tableModel, grantModel),
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

// confirms addition of resource monitor as part of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3318
func TestAcc_GrantOwnership_OnObject_ResourceMonitor_ToAccountRole(t *testing.T) {
	resourceMonitorId := testClient().Ids.RandomAccountObjectIdentifier()
	resourceMonitorName := resourceMonitorId.Name()
	resourceMonitorIdFullyQualifiedName := resourceMonitorId.FullyQualifiedName()

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	resourceMonitorModel := model.ResourceMonitor("test", resourceMonitorName)
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeResourceMonitor, resourceMonitorName).
		WithDependsOn(accountRoleModel.ResourceReference(), resourceMonitorModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, resourceMonitorModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeResourceMonitor)),
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
				Config:            accconfig.FromModels(t, accountRoleModel, resourceMonitorModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_SnowflakeIntelligence_ToAccountRole(t *testing.T) {
	snowflakeIntelligenceId, snowflakeIntelligenceCleanup := testClient().SnowflakeIntelligence.Create(t)
	t.Cleanup(snowflakeIntelligenceCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	accountRoleId := role.ID()

	resourceModel := model.GrantOwnership("test", []sdk.OwnershipGrantOn{
		{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeSnowflakeIntelligence,
				Name:       snowflakeIntelligenceId,
			},
		},
	}).WithAccountRoleName(accountRoleId.FullyQualifiedName())
	ref := resourceModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel),
				Check: assertThat(
					t,
					resourceassert.GrantOwnershipResource(t, ref).
						HasAccountRoleName(accountRoleId.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(ref, "on.0.object_type", string(sdk.ObjectTypeSnowflakeIntelligence))),
					assert.Check(resource.TestCheckResourceAttr(ref, "on.0.object_name", snowflakeIntelligenceId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(ref, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|%s|%s", accountRoleId.FullyQualifiedName(), sdk.ObjectTypeSnowflakeIntelligence, snowflakeIntelligenceId.FullyQualifiedName()))),
					assert.Check(checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeSnowflakeIntelligence, accountRoleId.Name(), snowflakeIntelligenceId.FullyQualifiedName())),
				),
			},
			{
				Config:       accconfig.FromModels(t, resourceModel),
				ResourceName: ref,
				ImportState:  true,
			},
		},
	})
}

// This test proves that managing grants on HYBRID TABLE is not supported in Snowflake. TABLE should be used instead.
func TestAcc_GrantOwnership_Validations_HybridTable_ToAccountRole_Fails(t *testing.T) {
	hybridTableId, hybridTableCleanup := testClient().HybridTable.Create(t)
	t.Cleanup(hybridTableCleanup)

	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()

	accountRoleModel := model.AccountRole("test", accountRoleName)
	grantOnHybridTableModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeHybridTable, hybridTableId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference())
	grantOnTableModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleName).
		WithOnObject(sdk.ObjectTypeTable, hybridTableId.FullyQualifiedName()).
		WithDependsOn(accountRoleModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, accountRoleModel, grantOnHybridTableModel),
				ExpectError: regexp.MustCompile("Unsupported feature"),
			},
			{
				Config: accconfig.FromModels(t, accountRoleModel, grantOnTableModel),
			},
		},
	})
}
