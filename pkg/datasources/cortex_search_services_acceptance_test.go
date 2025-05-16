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

func TestAcc_CortexSearchServices_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	dataSourceName := "data.snowflake_cortex_search_services.test"
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":  config.StringVariable(id.DatabaseName()),
			"schema":    config.StringVariable(id.SchemaName()),
			"table":     config.StringVariable(tableId.Name()),
			"name":      config.StringVariable(id.Name()),
			"on":        config.StringVariable("SOME_TEXT"),
			"warehouse": config.StringVariable(acc.TestWarehouseName),
			"query":     config.StringVariable(fmt.Sprintf("select SOME_TEXT from %s", tableId.FullyQualifiedName())),
			"comment":   config.StringVariable("Terraform acceptance test"),
		}
	}
	variableSet1 := m()

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
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cortex_search_services.0.created_on"),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.name", id.Name()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.comment", "Terraform acceptance test"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cortex_search_services.0.created_on"),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.name", id.Name()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(dataSourceName, "cortex_search_services.0.comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_CortexSearchServices_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      cortexSearchServicesDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_CortexSearchServices_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      cortexSearchServicesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func cortexSearchServicesDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_cortex_search_services" "test" {
  in {
    database = "%s"
    schema   = "%s"
  }
}
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func cortexSearchServicesDatasourceEmptyIn() string {
	return `
data "snowflake_cortex_search_services" "test" {
  in {
  }
}
`
}
