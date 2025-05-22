//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalStage(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				Config: externalStageConfig(id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func externalStageConfig(n, databaseName, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	name = "%v"
	url = "s3://com.example.bucket/prefix"
	database = "%s"
	schema = "%s"
	comment = "Terraform acceptance test"
}
`, n, databaseName, schemaName)
}
