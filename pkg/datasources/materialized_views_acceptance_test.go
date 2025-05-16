//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaterializedViews(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViews(acc.TestWarehouseName, tableId, viewId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_materialized_views.v", "materialized_views.#"),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "materialized_views.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "materialized_views.0.name", viewId.Name()),
				),
			},
		},
	})
}

func materializedViews(warehouseName string, tableId sdk.SchemaObjectIdentifier, viewId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource snowflake_table "t"{
		name 	 = "%[4]v"
		database = "%[2]s"
		schema 	 = "%[3]s"
		column {
			name = "column2"
			type = "VARCHAR(16)"
		}
	}

	resource snowflake_materialized_view "v"{
		name 	   = "%[5]v"
		comment    = "Terraform test resource"
		database   = "%[2]s"
		schema 	   = "%[3]s"
		is_secure  = true
		or_replace = false
		statement  = "SELECT * FROM ${snowflake_table.t.name}"
		warehouse  = "%[1]s"
	}

	data snowflake_materialized_views "v" {
		database = "%[2]s"
		schema = "%[3]s"
		depends_on = [snowflake_materialized_view.v]
	}
	`, warehouseName, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), viewId.Name())
}
