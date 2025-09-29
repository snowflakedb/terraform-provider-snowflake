//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DynamicTable_basic(t *testing.T) {
	dynamicTableId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	newWarehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableId.Name()),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, TestDatabaseName, TestSchemaName, tableId.Name())),
			"comment":    config.StringVariable(comment),
			"table_name": config.StringVariable(tableId.Name()),
		}
	}
	variableSet2 := m()
	variableSet2["warehouse"] = config.StringVariable(newWarehouseId.Name())
	variableSet2["comment"] = config.StringVariable(newComment)

	variableSet3 := m()
	variableSet3["initialize"] = config.StringVariable(string(sdk.DynamicTableInitializeOnSchedule))

	variableSet4 := m()
	variableSet4["initialize"] = config.StringVariable(string(sdk.DynamicTableInitializeOnSchedule)) // keep the same setting from set 3
	variableSet4["refresh_mode"] = config.StringVariable(string(sdk.DynamicTableRefreshModeFull))

	// used to check whether a dynamic table was replaced
	var createdOn string

	resourceName := "snowflake_dynamic_table.dt"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dynamicTableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", dynamicTableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnCreate)),
					resource.TestCheckResourceAttr(resourceName, "refresh_mode", string(sdk.DynamicTableRefreshModeAuto)),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.maximum_duration", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select \"id\" from \"%v\".\"%v\".\"%v\"", TestDatabaseName, TestSchemaName, tableId.Name())),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),

					// computed attributes

					// - not used at this time
					//  resource.TestCheckResourceAttrSet(resourceName, "cluster_by"),
					resource.TestCheckResourceAttrSet(resourceName, "rows"),
					resource.TestCheckResourceAttrSet(resourceName, "bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
					// - not used at this time
					// resource.TestCheckResourceAttrSet(resourceName, "automatic_clustering"),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_state"),
					resource.TestCheckResourceAttrSet(resourceName, "last_suspended_on"),
					resource.TestCheckResourceAttrSet(resourceName, "is_clone"),
					resource.TestCheckResourceAttrSet(resourceName, "is_replica"),
					resource.TestCheckResourceAttrSet(resourceName, "data_timestamp"),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
				),
			},
			// test target lag to downstream and change comment

			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dynamicTableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", dynamicTableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "warehouse", newWarehouseId.Name()),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.downstream", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", newComment),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("created_on changed from %v to %v", createdOn, value)
						}
						return nil
					}),
				),
			},
			// test changing initialize setting
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnSchedule)),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						createdOn = value
						return nil
					}),
				),
			},
			// test changing refresh_mode setting
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnSchedule)),
					resource.TestCheckResourceAttr(resourceName, "refresh_mode", string(sdk.DynamicTableRefreshModeFull)),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						return nil
					}),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestStepDirectory(),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
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

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":         config.StringVariable(dynamicTableName),
			"database":     config.StringVariable(TestDatabaseName),
			"schema":       config.StringVariable(TestSchemaName),
			"warehouse":    config.StringVariable(TestWarehouseName),
			"query":        config.StringVariable(query),
			"comment":      config.StringVariable("Terraform acceptance test for GH issue 2173"),
			"table_name":   config.StringVariable(tableName),
			"other_schema": config.StringVariable(otherSchemaName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
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
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
			{
				// We use the same config here as in the previous step so the plan should be empty.
				ConfigDirectory: ConfigurationSameAsStepN(2),
				ConfigVariables: m(),
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
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Terraform acceptance test for GH issue 2134"),
			"table_name": config.StringVariable(tableName),
		}
	}
	m2 := m()
	m2["comment"] = config.StringVariable("Changed comment")

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
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
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
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
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

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Terraform acceptance test for GH issue 2276"),
			"table_name": config.StringVariable(tableName),
		}
	}
	m2 := m()
	m2["query"] = config.StringVariable(newQuery)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			{
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
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
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":      config.StringVariable(dynamicTableName),
			"database":  config.StringVariable(TestDatabaseName),
			"schema":    config.StringVariable(TestSchemaName),
			"warehouse": config.StringVariable(TestWarehouseName),
			// spaces added on purpose
			"query":      config.StringVariable("  " + query),
			"comment":    config.StringVariable("Comment with AS on purpose"),
			"table_name": config.StringVariable(tableName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			// No changes are expected
			{
				ConfigDirectory: ConfigurationSameAsStepN(1),
				ConfigVariables: m(),
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
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Comment with AS SELECT on purpose"),
			"table_name": config.StringVariable(tableName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			// If we match more than one time (in this case in comment) we raise an explanation error.
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_DynamicTable_issue2329/1"),
				ConfigVariables: m(),
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
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableId.Name()),
			"database":   config.StringVariable(TestDatabaseName),
			"schema":     config.StringVariable(TestSchemaName),
			"warehouse":  config.StringVariable(TestWarehouseName),
			"query":      config.StringVariable(query),
			"table_name": config.StringVariable(tableId.Name()),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_DynamicTable_timeouts"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableId.Name()),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
		},
	})
}

func TestAcc_DynamicTable_issue3931_EmptyText(t *testing.T) {
	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	query := fmt.Sprintf(`select * from %v`, table.ID().FullyQualifiedName())
	targetLag := sdk.TargetLag{
		MaximumDuration: sdk.String("2 minutes"),
	}
	dynamicTableModel := model.DynamicTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), query, []sdk.TargetLag{targetLag}, TestWarehouseName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DynamicTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, dynamicTableModel),
				Check: assertThat(t,
					resourceassert.DynamicTableResource(t, dynamicTableModel.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(dynamicTableModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				PreConfig: func() {
					testClient().Grant.GrantOwnershipToAccountRole(t, role.ID(), sdk.ObjectTypeDynamicTable, id)
					testClient().Grant.GrantPrivilegesOnSchemaObjectToAccountRole(t, snowflakeroles.Accountadmin, sdk.ObjectTypeDynamicTable, id, []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}, false)

					// Added just in case if the test fails and the last step is not executed
					t.Cleanup(testClient().DynamicTable.DropDynamicTableFunc(t, id))
					t.Cleanup(func() {
						testClient().Grant.GrantOwnershipToAccountRoleWithOwnershipOptions(t, snowflakeroles.Accountadmin, sdk.ObjectTypeDynamicTable, id, sdk.GrantOwnershipOptions{CurrentGrants: &sdk.OwnershipCurrentGrants{OutboundPrivileges: sdk.Revoke}})
					})
				},
				Config:      accconfig.FromModels(t, dynamicTableModel),
				ExpectError: regexp.MustCompile("the text is empty for dynamic table"),
			},
			// Step needed to clean up the ownership to correctly remove the dynamic table
			{
				PreConfig: func() {
					testClient().Grant.GrantOwnershipToAccountRole(t, snowflakeroles.Accountadmin, sdk.ObjectTypeDynamicTable, id)
				},
				Config: accconfig.FromModels(t, dynamicTableModel),
			},
		},
	})
}
