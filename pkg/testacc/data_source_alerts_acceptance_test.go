//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Alerts(t *testing.T) {
	alertId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigNoOptionals(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigDbOnly(alertId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.#", "1"),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigDbAndSchema(alertId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.#", "1"),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigAllOptionals(alertId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.0.name", alertId.Name()),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigSchemaOnly(alertId),
				Check: resource.ComposeTestCheckFunc(
					// TODO [SNOW-1348349]: currently, the schema is taken into consideration only if the database is set (this works differently in the stable datasources); address it during the rework.
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
		},
	})
}

func alertsResourceConfig(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_alert" "test_resource_alert" {
	name     	      = "%s"
	database  	      = "%s"
	schema   	      = "%s"
	warehouse 	      = "%s"
	condition         = "select 0 as c"
	action            = "select 0 as c"
	enabled  	      = false
	comment           = "some comment"
	alert_schedule 	  {
		interval = "60"
	}
}
`, alertId.Name(), alertId.DatabaseName(), alertId.SchemaName(), TestWarehouseName)
}

func alertsDatasourceConfigNoOptionals() string {
	return `
data "snowflake_alerts" "test_datasource_alert" {}
`
}

func alertsDatasourceConfigDbOnly(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
}
`, alertId.DatabaseName())
}

func alertsDatasourceConfigDbAndSchema(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
	schema  	      = "%s"
}
`, alertId.DatabaseName(), alertId.SchemaName())
}

func alertsDatasourceConfigAllOptionals(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
	schema  	      = "%s"
	pattern  	      = "%s"
}
`, alertId.DatabaseName(), alertId.SchemaName(), alertId.Name())
}

func alertsDatasourceConfigSchemaOnly(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	schema  	      = "%s"
}
`, alertId.SchemaName())
}
