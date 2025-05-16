//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DynamicTables_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	dynamicTableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableId.Name()),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf("select \"id\" from %v", tableId.FullyQualifiedName())),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"table_name": config.StringVariable(tableId.Name()),
		}
	}
	variableSet1 := m()

	dataSourceName := "data.snowflake_dynamic_tables.dts"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "like.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "like.0.pattern", dynamicTableId.Name()),
					resource.TestCheckResourceAttr(dataSourceName, "in.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "in.0.database", dynamicTableId.DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "starts_with", dynamicTableId.Name()),
					resource.TestCheckResourceAttr(dataSourceName, "limit.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "limit.0.rows", "1"),

					// computed attributes
					resource.TestCheckResourceAttr(dataSourceName, "records.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.created_on"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.schema_name"),
					// unused by Snowflake API at this time (always empty)
					// resource.TestCheckResourceAttrSet(dataSourceName, "records.0.cluster_by"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.rows"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.target_lag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.refresh_mode"),
					// unused by Snowflake API at this time (always empty)
					// resource.TestCheckResourceAttrSet(dataSourceName, "records.0.refresh_mode_reason"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.warehouse"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.comment"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.text"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.automatic_clustering"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.scheduling_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.last_suspended_on"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.is_clone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.is_replica"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.data_timestamp"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "records.#", "1"),
				),
			},
		},
	})
}

func TestAcc_DynamicTables_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      dynamicTablesDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_DynamicTables_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      dynamicTablesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func dynamicTablesDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_dynamic_tables" "dts" {
  in {
    database = "%s"
    schema   = "%s"
  }
}
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func dynamicTablesDatasourceEmptyIn() string {
	return `
data "snowflake_dynamic_tables" "dts" {
  in {
  }
}
`
}
