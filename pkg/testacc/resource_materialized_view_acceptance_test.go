//go:build non_account_level_tests

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

func TestAcc_MaterializedView(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewName := viewId.Name()

	queryEscaped := fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\"", tableId.Name())
	query := fmt.Sprintf(`SELECT ID, DATA FROM "%s"`, tableId.Name())
	otherQueryEscaped := fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\" WHERE ID LIKE 'foo%%'", tableId.Name())
	otherQuery := fmt.Sprintf(`SELECT ID, DATA FROM "%s" WHERE ID LIKE 'foo%%'`, tableId.Name())

	comment := random.Comment()
	otherComment := random.Comment()

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(tableId, viewId, queryEscaped, comment, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "true"),
				),
			},
			// update parameters
			{
				Config: materializedViewConfig(tableId, viewId, queryEscaped, otherComment, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// change statement
			{
				Config: materializedViewConfig(tableId, viewId, otherQueryEscaped, otherComment, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// change statement externally
			{
				PreConfig: func() {
					testClient().MaterializedView.CreateMaterializedViewWithName(t, viewId, query, true)
				},
				Config: materializedViewConfig(tableId, viewId, otherQueryEscaped, otherComment, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", otherComment),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "is_secure", "false"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse"},
			},
		},
	})
}

func TestAcc_MaterializedView_Tags(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	tag1Id := testClient().Ids.RandomSchemaObjectIdentifier()
	tag2Id := testClient().Ids.RandomSchemaObjectIdentifier()

	queryEscaped := fmt.Sprintf("SELECT ID FROM \\\"%s\\\"", tableId.Name())

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			// create tags
			{
				Config: materializedViewConfigWithTags(tableId, viewId, queryEscaped, "test_tag", tag1Id, tag2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag1Id.Name()),
				),
			},
			// update tags
			{
				Config: materializedViewConfigWithTags(tableId, viewId, queryEscaped, "test_tag_2", tag1Id, tag2Id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag2Id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse", "tag"},
			},
		},
	})
}

func TestAcc_MaterializedView_Rename(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	viewId := testClient().Ids.RandomSchemaObjectIdentifier()
	newViewId := testClient().Ids.RandomSchemaObjectIdentifier()

	queryEscaped := fmt.Sprintf("SELECT ID FROM \\\"%s\\\"", tableId.Name())
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaterializedView),
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(tableId, viewId, queryEscaped, comment, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "fully_qualified_name", viewId.FullyQualifiedName()),
				),
			},
			// rename with one param change
			{
				Config: materializedViewConfig(tableId, newViewId, queryEscaped, comment, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", newViewId.Name()),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "fully_qualified_name", newViewId.FullyQualifiedName()),
				),
			},
		},
	})
}

func materializedViewConfig(tableId sdk.SchemaObjectIdentifier, viewId sdk.SchemaObjectIdentifier, q string, comment string, isSecure bool, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	name     = "%[1]s"
	database = "%[2]s"
	schema   = "%[3]s"

	column {
		name = "ID"
		type = "NUMBER(38,0)"
	}

	column {
		name = "DATA"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_materialized_view" "test" {
	name      = "%[4]s"
	comment   = "%[5]s"
	database  = "%[2]s"
	schema    = "%[3]s"
	warehouse = "%[6]s"
	is_secure = %[7]t
	or_replace = %[8]t
	statement = "%[9]s"

	depends_on = [
  		snowflake_table.test
  	]
}
`, tableId.Name(), viewId.DatabaseName(), viewId.SchemaName(), viewId.Name(), comment, TestWarehouseName, isSecure, orReplace, q)
}

func materializedViewConfigWithTags(tableId sdk.SchemaObjectIdentifier, viewId sdk.SchemaObjectIdentifier, q string, tagResourceName string, tag1Id sdk.SchemaObjectIdentifier, tag2Id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "ID"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_tag" "test_tag" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[8]s"
}

resource "snowflake_tag" "test_tag_2" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[9]s"
}

resource "snowflake_materialized_view" "test" {
	database  = "%[1]s"
	schema    = "%[2]s"
	name      = "%[4]s"
	warehouse = "%[5]s"
	statement = "%[6]s"

	tag {
		name = snowflake_tag.%[7]s.name
		schema = snowflake_tag.%[7]s.schema
		database = snowflake_tag.%[7]s.database
		value = "some_value"
	}

	depends_on = [
		snowflake_table.test
	]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), viewId.Name(), TestWarehouseName, q, tagResourceName, tag1Id.Name(), tag2Id.Name())
}
