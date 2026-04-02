//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Pipe(t *testing.T) {
	pipeId := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()

	tableModel := model.TableWithId("test", tableId, []sdk.TableColumnSignature{
		{Name: "id", Type: testdatatypes.DataTypeNumber},
		{Name: "data", Type: testdatatypes.DataTypeVarchar},
	})

	stageModel := model.Stage("test", stageId.DatabaseName(), stageId.SchemaName(), stageId.Name()).
		WithComment("Terraform acceptance test")

	pipeModel := model.PipeWithId("test", pipeId, "").
		WithComment("Terraform acceptance test").
		WithAutoIngest(false).
		WithCopyStatementCopyFromStageIntoTable(tableModel.ResourceReference(), stageModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel, stageModel, pipeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe.test", "name", pipeId.Name()),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "fully_qualified_name", pipeId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "auto_ingest", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "notification_channel", ""),
				),
			},
		},
	})
}
