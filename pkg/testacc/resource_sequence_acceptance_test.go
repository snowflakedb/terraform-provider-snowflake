//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Sequence(t *testing.T) {
	oldId := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Sequence),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: sequenceConfig(oldId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "ordering", "ORDER"),
				),
			},
			// Set comment and rename
			{
				Config: sequenceConfigWithComment(newId, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
			// Unset comment and set increment
			{
				Config: sequenceConfigWithIncrement(oldId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "increment", "32"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "ordering", "NOORDER"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", oldId.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_sequence.test_sequence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func sequenceConfig(sequenceId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_sequence" {
	database   = "%[1]s"
	schema     = "%[2]s"
	name       = "%[3]s"
}
`, sequenceId.DatabaseName(), sequenceId.SchemaName(), sequenceId.Name())
}

func sequenceConfigWithComment(sequenceId sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_sequence" {
	database   = "%[1]s"
	schema     = "%[2]s"
	name       = "%[3]s"
    comment    = "%[4]s"
}
`, sequenceId.DatabaseName(), sequenceId.SchemaName(), sequenceId.Name(), comment)
}

func sequenceConfigWithIncrement(sequenceId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_sequence" {
	database   = "%[1]s"
	schema     = "%[2]s"
	name       = "%[3]s"
    increment  = 32
	ordering   = "NOORDER"
}
`, sequenceId.DatabaseName(), sequenceId.SchemaName(), sequenceId.Name())
}
