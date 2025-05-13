//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tables(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	externalTableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tables(tableId, stageId, externalTableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.show_output.0.schema_name", tableId.SchemaName()),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.0.show_output.0.name", externalTableId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.1.show_output.0.schema_name", tableId.SchemaName()),
					resource.TestCheckResourceAttr("data.snowflake_tables.in_schema", "tables.1.show_output.0.name", tableId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "like", tableId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.0.show_output.0.database_name", tableId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.0.show_output.0.name", tableId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.0.describe_output.*", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.filtering", "tables.0.describe_output.0.name", "column2"),
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
		}
	    file_format = "TYPE = CSV"
	    location = "@${snowflake_stage.s.fully_qualified_name}"
	}

	data snowflake_tables "in_schema" {
		depends_on = [snowflake_table.t, snowflake_external_table.et]
		in {
			schema = "%[2]s"
		}
	}

	data snowflake_tables "filtering" {
		depends_on = [snowflake_table.t, snowflake_external_table.et]
		in {
			database = "%[1]s"
		}
		like = "%[3]s"
	}
	`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), stageId.Name(), externalTableId.Name())
}
