//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Sequences(t *testing.T) {
	sequenceId := testClient().Ids.RandomSchemaObjectIdentifier()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sequences(sequenceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "database", sequenceId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "schema", sequenceId.SchemaName()),
					resource.TestCheckResourceAttrSet("data.snowflake_sequences.t", "sequences.#"),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "sequences.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "sequences.0.name", sequenceId.Name()),
				),
			},
		},
	})
}

func sequences(sequenceId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource snowflake_sequence "t"{
		database = "%[1]s"
		schema 	 = "%[2]s"
		name 	 = "%[3]s"
	}

	data snowflake_sequences "t" {
		database = snowflake_sequence.t.database
		schema = snowflake_sequence.t.schema
	}
	`, sequenceId.DatabaseName(), sequenceId.SchemaName(), sequenceId.Name())
}
