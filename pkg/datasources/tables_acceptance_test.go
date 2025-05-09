//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tables(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	stageName := acc.TestClient().Ids.Alpha()
	externalTableName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tables(databaseName, schemaName, tableName, stageName, externalTableName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.name", tableName),
					resource.TestCheckResourceAttrSet("data.snowflake_tables.in_schema", "tables.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.database_name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.schema_name", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_tables.in_schema", "tables.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.comment", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_tables.in_schema", "tables.0.text"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.is_secure", "false"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.is_materialized", "false"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.owner_role_type", "ROLE"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.change_tracking", "OFF"),

					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.0.name", tableName),
				),
			},
		},
	})
}

func tables(tableId sdk.SchemaObjectIdentifier, stageId sdk.SchemaObjectIdentifier, externalTableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource snowflake_table "t"{
		database = "%[1]s"
		schema 	 = "%[2]s"
		name 	 = "%[3]s"
		column {
			name = "column2"
			type = "VARCHAR(16)"
		}
	}

	resource "snowflake_stage" "s" {
		database = "%[1]s"
		schema = "%[2]s"
		name = "%[4]s"
		url = "s3://snowflake-workshop-lab/weather-nyc"
	}

	resource "snowflake_external_table" "et" {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[5]s"
		column {
			name = "column1"
			type = "STRING"
			as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
			as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
		}
	    file_format = "TYPE = CSV"
	    location = "@${snowflake_stage.s.fully_qualified_name}"
	}

	data snowflake_tables "in_schema" {
		depends_on = [snowflake_table.t, snowflake_external_table.et]
		in {
			schema = snowflake_schema.s.fully_qualified_name
		}
	}

	data snowflake_tables "filtering" {
		depends_on = [snowflake_table.t, snowflake_external_table.et]
		in {
			database = snowflake_schema.s.database
		}
		like = "%[6]s"
	}
	`, databaseName, schemaName, tableName, stageName, externalTableName, tableName+"%", tableName+"%")
}
