//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalStage(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	stageModel := model.StageWithId("test", id).
		WithUrl("s3://com.example.bucket/prefix").
		WithComment("Terraform acceptance test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, stageModel),
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
