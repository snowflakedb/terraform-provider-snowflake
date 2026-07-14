//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DynamicTable_basic(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	newWarehouse, newWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(newWarehouseCleanup)
	comment := random.Comment()
	newComment := random.Comment()

	query := fmt.Sprintf(`select "id" from %s`, tableId.FullyQualifiedName())

	// used to check whether a dynamic table was replaced
	var createdOn string

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableId.Name(), []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
	}).WithChangeTracking(true)

	modelBasic := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableId.Name(), query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment(comment).
		WithDependsOn(tableModel.ResourceReference())

	modelWithDownstreamLag := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableId.Name(), query,
		[]sdk.TargetLag{{Downstream: sdk.Bool(true)}}, newWarehouse.ID().Name()).
		WithComment(newComment).
		WithDependsOn(tableModel.ResourceReference())

	modelWithInitialize := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableId.Name(), query,
		[]sdk.TargetLag{{Downstream: sdk.Bool(true)}}, TestWarehouseName).
		WithComment(comment).
		WithInitialize(string(sdk.DynamicTableInitializeOnSchedule)).
		WithDependsOn(tableModel.ResourceReference())

	modelWithRefreshMode := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableId.Name(), query,
		[]sdk.TargetLag{{Downstream: sdk.Bool(true)}}, TestWarehouseName).
		WithComment(comment).
		WithInitialize(string(sdk.DynamicTableInitializeOnSchedule)).
		WithRefreshMode(string(sdk.DynamicTableRefreshModeFull)).
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, modelBasic),
				Check: assertThat(t,
					resourceassert.DynamicTableResource(t, modelBasic.ResourceReference()).
						HasName(dynamicTableId.Name()).
						HasFullyQualifiedName(dynamicTableId.FullyQualifiedName()).
						HasDatabase(TestDatabaseName).
						HasSchema(TestSchemaName).
						HasWarehouse(TestWarehouseName).
						HasInitialize(string(sdk.DynamicTableInitializeOnCreate)).
						HasRefreshMode(string(sdk.DynamicTableRefreshModeAuto)).
						HasQuery(fmt.Sprintf("select \"id\" from \"%v\".\"%v\".\"%v\"", TestDatabaseName, TestSchemaName, tableId.Name())).
						HasComment(comment).
						HasRowsNotEmpty().
						HasBytesNotEmpty().
						HasOwnerNotEmpty().
						HasSchedulingStateNotEmpty().
						HasNoLastSuspendedOn().
						HasIsClone(false).
						HasIsReplica(false).
						HasDataTimestampNotEmpty(),
					// target_lag has no typed assertion
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "target_lag.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "target_lag.0.maximum_duration", "2 minutes")),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// test target lag to downstream and change comment
			{
				Config: accconfig.FromModels(t, tableModel, modelWithDownstreamLag),
				Check: assertThat(t,
					resourceassert.DynamicTableResource(t, modelWithDownstreamLag.ResourceReference()).
						HasName(dynamicTableId.Name()).
						HasFullyQualifiedName(dynamicTableId.FullyQualifiedName()).
						HasDatabase(TestDatabaseName).
						HasSchema(TestSchemaName).
						HasWarehouse(newWarehouse.ID().Name()).
						HasComment(newComment),
					assert.Check(resource.TestCheckResourceAttr(modelWithDownstreamLag.ResourceReference(), "target_lag.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithDownstreamLag.ResourceReference(), "target_lag.0.downstream", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(modelWithDownstreamLag.ResourceReference(), "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("created_on changed from %v to %v", createdOn, value)
						}
						return nil
					})),
				),
			},
			// test changing initialize setting
			{
				Config: accconfig.FromModels(t, tableModel, modelWithInitialize),
				Check: assertThat(t,
					resourceassert.DynamicTableResource(t, modelWithInitialize.ResourceReference()).
						HasInitialize(string(sdk.DynamicTableInitializeOnSchedule)),
					assert.Check(resource.TestCheckResourceAttrWith(modelWithInitialize.ResourceReference(), "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						createdOn = value
						return nil
					})),
				),
			},
			// test changing refresh_mode setting
			{
				Config: accconfig.FromModels(t, tableModel, modelWithRefreshMode),
				Check: assertThat(t,
					resourceassert.DynamicTableResource(t, modelWithRefreshMode.ResourceReference()).
						HasInitialize(string(sdk.DynamicTableInitializeOnSchedule)).
						HasRefreshMode(string(sdk.DynamicTableRefreshModeFull)),
					assert.Check(resource.TestCheckResourceAttrWith(modelWithRefreshMode.ResourceReference(), "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						return nil
					})),
				),
			},
			// test import
			{
				Config:            accconfig.FromModels(t, tableModel, modelWithDownstreamLag),
				ResourceName:      modelWithDownstreamLag.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAcc_DynamicTable_issue2173 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2173 issue.
func TestAcc_DynamicTable_issue2173(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	dynamicTableName := dynamicTableId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	query := fmt.Sprintf(`select "ID" from %s`, tableId.FullyQualifiedName())
	otherSchemaId := testClient().Ids.RandomDatabaseObjectIdentifier()
	otherSchemaName := otherSchemaId.Name()
	newDynamicTableId := testClient().Ids.NewSchemaObjectIdentifierInSchema(dynamicTableName, otherSchemaId)

	schemaModel := model.Schema("other_schema", TestDatabaseName, otherSchemaName).
		WithComment("Other schema")

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableName, []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber},
	}).WithChangeTracking(true)

	dtModel := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Terraform acceptance test for GH issue 2173").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schemaModel, tableModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.other_schema", "name", otherSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.t", "name", tableName),
				),
			},
			{
				PreConfig: func() {
					testClient().DynamicTable.CreateDynamicTableWithOptions(t, newDynamicTableId, testClient().Ids.WarehouseId(), tableId)
				},
				Config: accconfig.FromModels(t, schemaModel, tableModel, dtModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
			{
				// We use the same config here as in the previous step so the plan should be empty.
				Config: accconfig.FromModels(t, schemaModel, tableModel, dtModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					/*
					 * Before the fix this step resulted in
					 *     # snowflake_dynamic_table.dt will be updated in-place
					 *     ~ resource "snowflake_dynamic_table" "dt" {
					 *         + comment              = "Terraform acceptance test for GH issue 2173"
					 *           id                   = "terraform_test_database|terraform_test_schema|SFVNXKJFAA"
					 *           name                 = "SFVNXKJFAA"
					 *         ~ schema               = "MEYIYWUGGO" -> "terraform_test_schema"
					 *           # (14 unchanged attributes hidden)
					 *     }
					 * which matches the issue description exactly (issue mentioned also query being changed but here for simplicity the same underlying table and query were used).
					 */
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAcc_DynamicTable_issue2134 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2134 issue.
func TestAcc_DynamicTable_issue2134(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	dynamicTableName := dynamicTableId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	// whitespace (initial tab) is added on purpose here
	query := fmt.Sprintf(`	select "id" from "%v"."%v"."%v"`, TestDatabaseName, TestSchemaName, tableName)

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableName, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
	}).WithChangeTracking(true)

	dtModelInitial := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Terraform acceptance test for GH issue 2134").
		WithDependsOn(tableModel.ResourceReference())

	dtModelUpdatedComment := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Changed comment").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			/*
			 * Before the fix the first step resulted in not empty plan (as expected)
			 *     # snowflake_dynamic_table.dt will be updated in-place
			 *     ~ resource "snowflake_dynamic_table" "dt" {
			 *         id                   = "terraform_test_database|terraform_test_schema|IKLBYWKSOV"
			 *         name                 = "IKLBYWKSOV"
			 *         ~ query                = "select \"id\" from \"terraform_test_database\".\"terraform_test_schema\".\"IKLBYWKSOV_table\"" -> "\tselect \"id\" from \"terraform_test_database\".\"terraform_test_schema\".\"IKLBYWKSOV_table\""
			 *         # (15 unchanged attributes hidden)
			 *     }
			 * which matches the issue description exactly.
			 */
			{
				Config: accconfig.FromModels(t, tableModel, dtModelInitial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
			/*
			 * Before the fix the second step resulted in SQL error (as expected)
			 *     Error: 001003 (42000): SQL compilation error:
			 *         syntax error line 1 at position 86 unexpected '<EOF>'.
			 * which matches the issue description exactly.
			 */
			{
				Config: accconfig.FromModels(t, tableModel, dtModelUpdatedComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2276 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2276 issue.
func TestAcc_DynamicTable_issue2276(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	dynamicTableName := dynamicTableId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, TestDatabaseName, TestSchemaName, tableName)
	newQuery := fmt.Sprintf(`select "data" from "%v"."%v"."%v"`, TestDatabaseName, TestSchemaName, tableName)

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableName, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
		{Name: "data", Type: testdatatypes.DataTypeVarchar},
	}).WithChangeTracking(true)

	dtModelInitial := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Terraform acceptance test for GH issue 2276").
		WithDependsOn(tableModel.ResourceReference())

	dtModelUpdatedQuery := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, newQuery,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Terraform acceptance test for GH issue 2276").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, dtModelInitial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			{
				Config: accconfig.FromModels(t, tableModel, dtModelUpdatedQuery),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", newQuery),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2329 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2329 issue.
func TestAcc_DynamicTable_issue2329(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifierContaining("AS")
	dynamicTableName := dynamicTableId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, TestDatabaseName, TestSchemaName, tableName)

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableName, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
		{Name: "data", Type: testdatatypes.DataTypeVarchar},
	}).WithChangeTracking(true)

	dtModel := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName,
		// spaces added on purpose
		"  "+query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Comment with AS on purpose").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, dtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			// No changes are expected
			{
				Config: accconfig.FromModels(t, tableModel, dtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2329_with_matching_comment proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2329 issue.
func TestAcc_DynamicTable_issue2329_with_matching_comment(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifierContaining("AS")
	dynamicTableName := dynamicTableId.Name()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	query := fmt.Sprintf(`with temp as (select "id" from "%v"."%v"."%v") select * from temp`, TestDatabaseName, TestSchemaName, tableName)

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableName, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
		{Name: "data", Type: testdatatypes.DataTypeVarchar},
	}).WithChangeTracking(true)

	dtModel := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableName, query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithComment("Comment with AS SELECT on purpose").
		WithDependsOn(tableModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			// If we match more than one time (in this case in comment) we raise an explanation error.
			{
				Config: accconfig.FromModels(t, tableModel, dtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
		},
	})
}

func TestAcc_DynamicTable_issue3355_timeout(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	query := fmt.Sprintf(`with temp as (select "id" from %v) select * from temp`, tableId.FullyQualifiedName())

	tableModel := model.Table("t", TestDatabaseName, TestSchemaName, tableId.Name(), []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
	}).WithChangeTracking(true)

	dtModel := model.DynamicTable("dt", TestDatabaseName, TestSchemaName, dynamicTableId.Name(), query,
		[]sdk.TargetLag{{MaximumDuration: sdk.String("2 minutes")}}, TestWarehouseName).
		WithDependsOn(tableModel.ResourceReference()).
		WithTimeout(accconfig.Timeouts{Create: "50ms", Read: "50ms", Update: "50ms", Delete: "50ms"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, dtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableId.Name()),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
		},
	})
}
