//go:build non_account_level_tests

package testacc

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Alert_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseName := testClient().Ids.WarehouseId().Name()

	basic := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5)

	complete := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5).
		WithEnabled(true).
		WithComment("initial")

	updated := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 3", "SELECT 2", warehouseName).
		WithAlertScheduleInterval(15).
		WithEnabled(true).
		WithComment("updated")

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(false),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateSuspended)),
	}

	initialAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("initial").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(true),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("initial").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateStarted),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("initial").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateStarted)),
	}

	updatedAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("updated").
			HasConditionString("SELECT 2").
			HasActionString("SELECT 3").
			HasEnabled(true),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("15 MINUTE").
			HasComment("updated").
			HasCondition("SELECT 2").
			HasAction("SELECT 3").
			HasState(sdk.AlertStateStarted),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("15 MINUTE").
			HasComment("updated").
			HasCondition("SELECT 2").
			HasAction("SELECT 3").
			HasState(string(sdk.AlertStateStarted)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			// Create - no options
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import - no options
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set comment and enable
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, initialAssertions...),
			},
			// Import - with comment
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - change comment, condition, action, and schedule
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedAssertions...),
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					client := atc.defaultTestEnv.client
					ctx := context.Background()
					require.NoError(t, client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).WithSet(
						*sdk.NewAlertSetRequest().
							WithSchedule("10 MINUTE").
							WithComment("external"),
					)))
					require.NoError(t, client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
						WithModifyCondition([]string{"SELECT 4"})))
					require.NoError(t, client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
						WithModifyAction("SELECT 5")))
					require.NoError(t, client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
						WithResume(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Create - recreate dropped alert
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Alerts.Drop(
						context.Background(),
						sdk.NewDropAlertRequest(id),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// No change
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, basic),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider()),
			},
		},
	})
}

func TestAcc_Alert_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseName := testClient().Ids.WarehouseId().Name()

	complete := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleCron("0 0 * * *", "UTC").
		WithComment("complete").
		WithEnabled(true)

	updated := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 3", "SELECT 2", warehouseName).
		WithAlertScheduleCron("0 1 * * *", "Europe/Warsaw").
		WithComment("updated").
		WithEnabled(false)

	basic := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5)

	ref := basic.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("complete").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(true).
			HasAlertScheduleCron("0 0 * * *", "UTC"),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("USING CRON 0 0 * * * UTC").
			HasComment("complete").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateStarted),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("USING CRON 0 0 * * * UTC").
			HasComment("complete").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateStarted)),
	}

	updatedAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("updated").
			HasConditionString("SELECT 2").
			HasActionString("SELECT 3").
			HasEnabled(false).
			HasAlertScheduleCron("0 1 * * *", "Europe/Warsaw"),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("USING CRON 0 1 * * * Europe/Warsaw").
			HasComment("updated").
			HasCondition("SELECT 2").
			HasAction("SELECT 3").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("USING CRON 0 1 * * * Europe/Warsaw").
			HasComment("updated").
			HasCondition("SELECT 2").
			HasAction("SELECT 3").
			HasState(string(sdk.AlertStateSuspended)),
	}

	basicAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(false).
			HasAlertScheduleInterval(5),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateSuspended)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			// Create - with cron schedule, comment, and enabled
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import - with cron schedule
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - change cron, comment, condition, and action
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedAssertions...),
			},
			// Import - with updated cron
			{
				Config:            config.FromModels(t, updated),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - switch cron to interval and unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Update - external schedule and comment
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Alerts.Alter(
						context.Background(),
						sdk.NewAlterAlertRequest(id).WithSet(
							*sdk.NewAlertSetRequest().
								WithSchedule("20 MINUTE").
								WithComment("external"),
						),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Create - recreate dropped alert
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Alerts.Drop(
						context.Background(),
						sdk.NewDropAlertRequest(id),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// No change
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, complete),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider()),
			},
		},
	})
}

func TestAcc_Alert_Replacement(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseName := testClient().Ids.WarehouseId().Name()
	warehouse, cleanupWarehouse := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(cleanupWarehouse)

	basic := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5)

	renamed := model.Alert("test", renamedId.DatabaseName(), renamedId.SchemaName(), renamedId.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5)

	changedWarehouse := model.Alert("test", renamedId.DatabaseName(), renamedId.SchemaName(), renamedId.Name(), "SELECT 1", "SELECT 1", warehouse.Name).
		WithAlertScheduleInterval(5)

	ref := basic.ResourceReference()

	renamedAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, renamedId).
			HasWarehouseString(warehouseName).
			HasCommentString("").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(false),
		newAlertObjectAssertions(t, renamedId).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, renamedId).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateSuspended)),
	}

	changedWarehouseAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, renamedId).
			HasWarehouseString(warehouse.Name).
			HasCommentString("").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(false),
		newAlertObjectAssertions(t, renamedId).
			HasWarehouse(warehouse.Name).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, renamedId).
			HasWarehouse(warehouse.Name).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateSuspended)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			// Create - no options
			{
				Config: config.FromModels(t, basic),
			},
			// Replace - rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, renamedAssertions...),
			},
			// Replace - change warehouse
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, changedWarehouse),
				Check:  assertThat(t, changedWarehouseAssertions...),
			},
		},
	})
}

func TestAcc_Alert_SQLFormattingAndEmptyComment(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseName := testClient().Ids.WarehouseId().Name()

	basic := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5)

	formatted := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleInterval(5).
		WithCondition("select  1").
		WithActionValue(config.MultilineWrapperVariable("select\n1")).
		WithComment("").
		WithEnabled(false)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		newAlertResourceAssertions(t, ref, id).
			HasWarehouseString(warehouseName).
			HasCommentString("").
			HasConditionString("SELECT 1").
			HasActionString("SELECT 1").
			HasEnabled(false),
		newAlertObjectAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(sdk.AlertStateSuspended),
		newAlertDetailAssertions(t, id).
			HasWarehouse(warehouseName).
			HasSchedule("5 MINUTE").
			HasComment("").
			HasCondition("SELECT 1").
			HasAction("SELECT 1").
			HasState(string(sdk.AlertStateSuspended)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			// Create - no options
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// No change - formatted condition and action and empty comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, formatted),
			},
			// Import - formatted condition and action
			{
				Config:            config.FromModels(t, formatted),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Alert_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	warehouseName := testClient().Ids.WarehouseId().Name()

	bothSchedules := model.Alert("test", id.DatabaseName(), id.SchemaName(), id.Name(), "SELECT 1", "SELECT 1", warehouseName).
		WithAlertScheduleIntervalAndCron(5, "0 0 * * *", "UTC")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan - interval and cron together
			{
				Config:      config.FromModels(t, bothSchedules),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"alert_schedule\.0\.(?:cron|interval)": conflicts with\s+alert_schedule\.0\.(?:interval|cron)`),
			},
		},
	})
}

// Can't reproduce the issue, leaving the test for now.
func TestAcc_Alert_Issue3117(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix("small caps with spaces")
	providerConfig := providermodel.V097CompatibleProviderConfig(t)
	// The legacy provider emits the warehouse as an unquoted identifier.
	warehouseId := sdk.NewAccountObjectIdentifier(strings.ToUpper(testClient().Ids.RandomAccountObjectIdentifier().Name()))
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId))
	t.Cleanup(warehouseCleanup)

	alertModel := model.Alert("test_alert", id.DatabaseName(), id.SchemaName(), id.Name(), "select 0 as c", "select 0 as c", warehouse.Name).
		WithAlertScheduleInterval(1).
		WithEnabled(true).
		WithComment("Alert config for GH issue 3117")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			// Create - legacy provider
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.92.0"),
				Config:            providerConfig + config.FromModels(t, alertModel),
				Check: assertThat(t,
					resourceassert.AlertResource(t, "snowflake_alert.test_alert").
						HasNameString(id.Name()),
				),
			},
			// Upgrade - current provider
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, alertModel),
				Check: assertThat(t,
					resourceassert.AlertResource(t, "snowflake_alert.test_alert").
						HasNameString(id.Name()),
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
			// Create - first alert
			{
				Config: config.FromModels(t, alertModel1),
				Check: assertThat(t,
					resourceassert.AlertResource(t, "snowflake_alert.test_alert_1").
						HasNameString(id1.Name()),
				),
			},
			// Create - second alert
			{
				Config: config.FromModels(t, alertModel1, alertModel2),
				Check: assertThat(t,
					resourceassert.AlertResource(t, "snowflake_alert.test_alert_1").
						HasNameString(id1.Name()),
					resourceassert.AlertResource(t, "snowflake_alert.test_alert_2").
						HasNameString(id2.Name()),
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
			// Create - quoted identifiers ignore case
			{
				PreConfig: func() {
					testClient().Database.Alter(t, sdk.NewAlterDatabaseRequest(database.ID()).WithSet(
						*sdk.NewDatabaseSetRequest().
							WithQuotedIdentifiersIgnoreCase(true),
					))
				},
				Config: config.FromModels(t, alertModel),
				Check: assertThat(t,
					resourceassert.AlertResource(t, "snowflake_alert.test_alert").
						HasNameString(id.Name()),
				),
			},
		},
	})
}

// Each call returns a new assert. Chain step-specific checks on the returned value.
func newAlertResourceAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier) *resourceassert.AlertResourceAssert {
	t.Helper()
	return resourceassert.AlertResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
}

func newAlertObjectAssertions(t *testing.T, id sdk.SchemaObjectIdentifier) *objectassert.AlertAssert {
	t.Helper()
	return objectassert.Alert(t, id).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwnerRoleType("ROLE")
}

func newAlertDetailAssertions(t *testing.T, id sdk.SchemaObjectIdentifier) *objectassert.AlertDetailsAssert {
	t.Helper()
	return objectassert.AlertDetails(t, id).
		HasName(id.Name())
}
