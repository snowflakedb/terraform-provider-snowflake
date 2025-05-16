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

func TestAcc_ExternalTables(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	externalTableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalTables(stageId, externalTableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "database", stageId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "schema", stageId.SchemaName()),
					resource.TestCheckResourceAttrSet("data.snowflake_external_tables.t", "external_tables.#"),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "external_tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "external_tables.0.name", externalTableId.Name()),
				),
			},
		},
	})
}

func externalTables(stageId sdk.SchemaObjectIdentifier, externalTableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_stage" "test" {
		name = "%[3]s"
		url = "s3://snowflake-workshop-lab/weather-nyc"
		database = "%[1]s"
		schema = "%[2]s"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_external_table" "test_table" {
		database = "%[1]s"
		schema   = "%[2]s"
		name     = "%[4]s"
		comment  = "Terraform acceptance test"
		column {
			name = "column1"
			type = "STRING"
		as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
		}
	    file_format = "TYPE = CSV"
	    location = "@${snowflake_stage.test.fully_qualified_name}"
	}

	data snowflake_external_tables "t" {
		database = "%[1]s"
		schema = "%[2]s"
		depends_on = [snowflake_external_table.test_table]
	}
	`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(), externalTableId.Name())
}
