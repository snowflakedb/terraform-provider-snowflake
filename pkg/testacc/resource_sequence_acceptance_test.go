//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Sequence(t *testing.T) {
	oldId := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	sequenceModelBasic := model.SequenceWithId("test_sequence", oldId)

	sequenceModelWithComment := model.SequenceWithId("test_sequence", newId).
		WithComment(comment)

	sequenceModelWithIncrement := model.SequenceWithId("test_sequence", oldId).
		WithIncrement(32).
		WithOrdering("NOORDER")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Sequence),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: accconfig.FromModels(t, sequenceModelBasic),
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
				Config: accconfig.FromModels(t, sequenceModelWithComment),
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
				Config: accconfig.FromModels(t, sequenceModelWithIncrement),
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
