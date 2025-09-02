//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Pipes(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: pipes(tableId, stageId, pipeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "schema", TestSchemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_pipes.t", "pipes.#"),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "pipes.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "pipes.0.name", pipeId.Name()),
				),
			},
		},
	})
}

func pipes(tableId sdk.SchemaObjectIdentifier, stageId sdk.SchemaObjectIdentifier, pipeId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
  database = "%[1]s"
  schema   = "%[2]s"
  name     = "%[3]s"
  column {
	name = "id"
	type = "NUMBER(5,0)"
  }
  column {
    name = "data"
	type = "VARCHAR(16)"
  }
}

resource "snowflake_stage" "test" {
  database = "%[1]s"
  schema = "%[2]s"
  name = "%[4]s"
}

resource "snowflake_pipe" "test" {
  database       = "%[1]s"
  schema         = "%[2]s"
  name           = "%[5]s"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO ${snowflake_table.test.fully_qualified_name}
  FROM @${snowflake_stage.test.fully_qualified_name}
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}

data snowflake_pipes "t" {
	database = snowflake_pipe.test.database
	schema = snowflake_pipe.test.schema
}
`, pipeId.DatabaseName(), pipeId.SchemaName(), tableId.Name(), stageId.Name(), pipeId.Name())
}
