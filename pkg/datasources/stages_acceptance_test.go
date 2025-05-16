//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Stages(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	storageIntegrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: stages(stageId, storageIntegrationId, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "database", stageId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "schema", stageId.SchemaName()),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.name", stageId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.storage_integration", storageIntegrationId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.comment", comment),
				),
			},
		},
	})
}

func stages(stageId sdk.SchemaObjectIdentifier, storageIntegrationId sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
	resource "snowflake_storage_integration" "test" {
		name = "%[4]s"
		storage_allowed_locations = ["s3://foo/"]
		storage_provider = "S3"

		storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
	}

	resource "snowflake_stage" "test"{
		database 				= "%[1]s"
		schema 	 				= "%[2]s"
		name 	 				= "%[3]s"
		url 					= "s3://foo/"
		storage_integration 	= snowflake_storage_integration.test.name
		comment  				= "%[5]s"
	}

	data "snowflake_stages" "test" {
		database = snowflake_stage.test.database
		schema = snowflake_stage.test.schema
	}
	`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(), storageIntegrationId.Name(), comment)
}
