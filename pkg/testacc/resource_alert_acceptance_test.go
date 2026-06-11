//go:build non_account_level_tests

package testacc

import (
	"strconv"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Alert(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	alertModelBasic := model.Alert("test_alert", TestDatabaseName, TestSchemaName, id.Name(), "select 0 as c", "select 0 as c", TestWarehouseName).
		WithAlertScheduleInterval(5).
		WithEnabled(true).
		WithComment("dummy")

	alertModelUpdatedParams := model.Alert("test_alert", TestDatabaseName, TestSchemaName, id.Name(), "select 1 as c", "select 1 as c", TestWarehouseName).
		WithAlertScheduleInterval(15).
		WithEnabled(true).
		WithComment("test")

	alertModelFurtherUpdated := model.Alert("test_alert", TestDatabaseName, TestSchemaName, id.Name(), "select 2 as c", "select 2 as c", TestWarehouseName).
		WithAlertScheduleInterval(25).
		WithEnabled(true).
		WithComment("text")

	alertModelDisabled := model.Alert("test_alert", TestDatabaseName, TestSchemaName, id.Name(), "select 2 as c", "select 2 as c", TestWarehouseName).
		WithAlertScheduleInterval(5).
		WithEnabled(false).
		WithComment("")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, alertModelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", "select 0 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", "select 0 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", "dummy"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(5)),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModelUpdatedParams),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", "select 1 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", "select 1 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", "test"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(15)),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModelFurtherUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", "select 2 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", "select 2 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", "text"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(25)),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModelDisabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "enabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", "select 2 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", "select 2 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(5)),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModelBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", "select 0 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", "select 0 as c"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", "dummy"),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(5)),
				),
			},
		},
	})
}

// Can't reproduce the issue, leaving the test for now.
func TestAcc_Alert_Issue3117(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix("small caps with spaces")
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	alertModel := model.Alert("test_alert", id.DatabaseName(), id.SchemaName(), id.Name(), "select 0 as c", "select 0 as c", testClient().Ids.WarehouseId().Name()).
		WithAlertScheduleInterval(1).
		WithEnabled(true).
		WithComment("Alert config for GH issue 3117")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            providerConfig + accconfig.FromModels(t, alertModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, alertModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
				),
			},
		},
	})
}

// Can't reproduce the issue, leaving the test for now.
func TestAcc_Alert_Issue3117_PatternMatching(t *testing.T) {
	suffix := testClient().Ids.Alpha()
	id1 := testClient().Ids.NewSchemaObjectIdentifier("prefix1" + suffix)
	id2 := testClient().Ids.NewSchemaObjectIdentifier("prefix_" + suffix)

	warehouseId := testClient().Ids.WarehouseId()

	alertModel1 := model.Alert("test_alert_1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), "select 0 as c", "select 0 as c", warehouseId.Name()).
		WithAlertScheduleInterval(1).
		WithEnabled(true).
		WithComment("Alert config for GH issue 3117")

	alertModel2 := model.Alert("test_alert_2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), "select 0 as c", "select 0 as c", warehouseId.Name()).
		WithAlertScheduleInterval(1).
		WithEnabled(true).
		WithComment("Alert config for GH issue 3117")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, alertModel1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert_1", "name", id1.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, alertModel1, alertModel2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert_1", "name", id1.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert_2", "name", id2.Name()),
				),
			},
		},
	})
}

// Can't reproduce the issue, leaving the test for now.
func TestAcc_Alert_Issue3117_IgnoreQuotedIdentifierCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.NewSchemaObjectIdentifierInSchema("small_"+testClient().Ids.Alpha(), schema.ID())

	alertModel := model.Alert("test_alert", id.DatabaseName(), id.SchemaName(), id.Name(), "select 0 as c", "select 0 as c", testClient().Ids.WarehouseId().Name()).
		WithAlertScheduleInterval(1).
		WithEnabled(true).
		WithComment("Alert config for GH issue 3117")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Database.Alter(t, database.ID(), &sdk.AlterDatabaseOptions{
						Set: &sdk.DatabaseSet{
							QuotedIdentifiersIgnoreCase: sdk.Bool(true),
						},
					})
				},
				Config: accconfig.FromModels(t, alertModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
				),
			},
		},
	})
}
