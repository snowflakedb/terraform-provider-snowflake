//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "database", tableId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "schema", tableId.SchemaName()),
					resource.TestCheckResourceAttrSet("data.snowflake_tables.t", "tables.#"),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "tables.0.name", tableId.Name()),
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

	data snowflake_tables "t" {
		database = "%[1]s"
		schema = "%[2]s"
		depends_on = [snowflake_table.t, snowflake_external_table.et]
	}
	`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), stageId.Name(), externalTableId.Name())
}
