//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToShare_OnDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				Config:            accconfig.FromModels(t, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchema(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())

	schemaModel := model.Schema("test", database.ID().Name(), schemaId.Name())
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeUsage.String()}, share.ID().Name()).
		WithOnSchema(schemaId.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schemaModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_schema", schemaId.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, schemaModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, schemaModel),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, share.ID().Name()).
		WithOnTable(tableId.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", tableId.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, tableModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, tableModel),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDynamicTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	baseTableId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), id.Name()+"_base")

	baseTableModel := model.TableWithId("base_table", baseTableId, []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}}).
		WithChangeTracking(true)
	dynamicTableQuery := fmt.Sprintf("with temp as (\n  select \"id\" from %s\n)\nselect * from temp", baseTableId.FullyQualifiedName())
	dynamicTableModel := model.DynamicTable("test_dynamic_table", id.DatabaseName(), id.SchemaName(), id.Name(),
		dynamicTableQuery, []sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithDependsOn(baseTableModel.ResourceReference())
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, share.ID().Name()).
		WithOnTable(id.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, baseTableModel, dynamicTableModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", id.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, baseTableModel, dynamicTableModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnAllTablesInSchema(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, share.ID().Name()).
		WithOnAllTablesInSchema(schema.ID().FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_all_tables_in_schema", schema.ID().FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnView(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	viewId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	viewStatement := fmt.Sprintf(`select "id" from %s`, tableId.FullyQualifiedName())
	viewModel := model.View("test", database.ID().Name(), schema.ID().Name(), viewId.Name(), viewStatement).
		WithIsSecure("true").
		WithColumnNames("id").
		WithDependsOn(tableModel.ResourceReference())
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, share.ID().Name()).
		WithOnView(viewId.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, viewModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_view", viewId.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, tableModel, viewModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, tableModel, viewModel),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTag(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	tagId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	tagModel := model.Tag("test", database.ID().Name(), schema.ID().Name(), tagId.Name())
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeRead.String()}, share.ID().Name()).
		WithOnTag(tagId.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: tagsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeRead.String()),
					resource.TestCheckResourceAttr(resourceName, "on_tag", tagId.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, tagModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, tagModel),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchemaObject_OnFunctionWithArguments(t *testing.T) {
	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	function := testClient().Function.CreateSecure(t, sdk.DataTypeFloat)

	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(TestDatabaseName)
	grantModel := model.GrantPrivilegesToShare("test", []string{string(sdk.SchemaObjectPrivilegeUsage)}, share.ID().Name()).
		WithOnFunction(function.ID().FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, setupGrantModel, grantModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_function", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|USAGE|OnFunction|%s", share.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				Config:            accconfig.FromModels(t, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchemaObject_OnFunctionWithoutArguments(t *testing.T) {
	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	function := testClient().Function.CreateSecure(t)

	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, share.ID().Name()).
		WithOnDatabase(TestDatabaseName)
	grantModel := model.GrantPrivilegesToShare("test", []string{string(sdk.SchemaObjectPrivilegeUsage)}, share.ID().Name()).
		WithOnFunction(function.ID().FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, setupGrantModel, grantModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_function", function.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|USAGE|OnFunction|%s", share.ID().FullyQualifiedName(), function.ID().FullyQualifiedName())),
				),
			},
			{
				Config:            accconfig.FromModels(t, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnPrivilegeUpdate(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	grantModelOnePrivilege := model.GrantPrivilegesToShare("test",
		[]string{sdk.ObjectPrivilegeUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())

	grantModelTwoPrivileges := model.GrantPrivilegesToShare("test",
		[]string{sdk.ObjectPrivilegeUsage.String(), sdk.ObjectPrivilegeReferenceUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModelOnePrivilege),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, grantModelTwoPrivileges),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				Config:            accconfig.FromModels(t, grantModelTwoPrivileges),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDatabaseWithReferenceUsagePrivilege(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	grantModel := model.GrantPrivilegesToShare("test",
		[]string{sdk.ObjectPrivilegeUsage.String(), sdk.ObjectPrivilegeReferenceUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", share.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", database.ID().Name()),
				),
			},
			{
				Config:            accconfig.FromModels(t, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoPrivileges(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = "some_share"
  on_database = "some_database"
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "privileges" is required, but no definition was found.`),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoOnOption(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = "some_share"
  privileges = ["USAGE"]
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToShare_RemoveShareOutsideTerraform(t *testing.T) {
	t.Skip("Should be addressed in SNOW-2048471 - without this task done, the test cannot be properly asserted in the terraform testing framework")

	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeUsage.String()}, share.ID().Name()).
		WithOnDatabase(database.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, grantModel),
			},
			{
				PreConfig: func() { shareCleanup() },
				Config:    accconfig.FromModels(t, grantModel),
				ExpectError: regexp.MustCompile(sdk.ErrObjectNotExistOrAuthorized.Error()),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShareWithNameContainingDots_OnTable(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	shareId := testClient().Ids.RandomAccountObjectIdentifierContaining(".foo.bar")
	_, shareCleanup := testClient().Share.CreateShareWithIdentifier(t, shareId)
	t.Cleanup(shareCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{{Name: "id", Type: testdatatypes.DataTypeNumber}})
	setupGrantModel := model.GrantPrivilegesToShare("test_setup", []string{"USAGE"}, shareId.Name()).
		WithOnDatabase(database.ID().Name())
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeSelect.String()}, shareId.Name()).
		WithOnTable(tableId.FullyQualifiedName()).
		WithDependsOn(setupGrantModel.ResourceReference())

	resourceName := "snowflake_grant_privileges_to_share.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, setupGrantModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", tableId.FullyQualifiedName()),
				),
			},
			{
				Config:            accconfig.FromModels(t, tableModel, setupGrantModel, grantModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, tableModel),
				Check:  CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + grantPrivilegesToShareBasicConfig(database.ID(), share.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, share.ID().FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToShareBasicConfig(database.ID(), share.ID()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, share.ID().FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToShareBasicConfig(databaseId sdk.AccountObjectIdentifier, shareId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = "%[2]s"
  privileges  = ["USAGE"]
  on_database = "%[1]s"
}
`, databaseId.Name(), shareId.Name())
}

func TestAcc_GrantPrivilegesToShare_IdentifierQuotingDiffSuppression(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	shareId := testClient().Ids.RandomAccountObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				ExpectError:       regexp.MustCompile("Error: Provider produced inconsistent final plan"),
				Config:            providerConfig + grantPrivilegesToShareQuotedIdentifiers(database.ID(), shareId),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   grantPrivilegesToShareQuotedIdentifiers(database.ID(), shareId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_share.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "to_share", shareId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "on_database", database.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_share.test", "id", fmt.Sprintf(`%s|USAGE|OnDatabase|%s`, shareId.FullyQualifiedName(), database.ID().FullyQualifiedName())),
				),
			},
		},
	})
}

func grantPrivilegesToShareQuotedIdentifiers(databaseId sdk.AccountObjectIdentifier, shareId sdk.AccountObjectIdentifier) string {
	quotedShareId := fmt.Sprintf(`\"%s\"`, shareId.Name())

	return fmt.Sprintf(`
resource "snowflake_share" "test" {
  name       = "%[2]s"
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = "%[1]s"
}
`, databaseId.Name(), quotedShareId)
}
