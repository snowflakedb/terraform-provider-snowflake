//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_InternalStage(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				Config: internalStageConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func internalStageConfig(stageId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	comment = "Terraform acceptance test"
}
`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name())
}
