//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Pipe(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	pipeId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			{
				Config: pipeConfig(pipeId, tableId, stageId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe.test", "name", pipeId.Name()),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "fully_qualified_name", pipeId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "auto_ingest", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "notification_channel", ""),
				),
			},
		},
	})
}

// whitespace in copy_statement matters for the tests, change with caution!
func pipeConfig(pipeId sdk.SchemaObjectIdentifier, tableId sdk.SchemaObjectIdentifier, stageId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
  	schema   = "%[2]s"
	name     = "%[4]s"

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
	name = "%[5]s"
	comment = "Terraform acceptance test"
}

resource "snowflake_pipe" "test" {
  database       = "%[1]s"
  schema         = "%[2]s"
  name           = "%[3]s"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
  	COPY INTO ${snowflake_table.test.fully_qualified_name}
  FROM @${snowflake_stage.test.fully_qualified_name}
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`, pipeId.DatabaseName(), pipeId.SchemaName(), pipeId.Name(), tableId.Name(), stageId.Name())
}
