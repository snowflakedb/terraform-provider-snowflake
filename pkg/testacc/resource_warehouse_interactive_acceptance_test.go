//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_WarehouseInteractive_BasicUseCase(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basic := model.WarehouseInteractiveWithId(warehouseId)
	// warehouse_size is not toggled in this set/unset cycle: removing it from config forces a resource
	// recreate (Snowflake has no UNSET WAREHOUSE_SIZE), which does not fit an update-only unset step.
	// The warehouse_size resize (update) and its removal (recreate) are exercised in the dedicated
	// TestAcc_WarehouseInteractive_WarehouseSize test.
	//
	// Every optional value below is chosen to differ from the interactive warehouse defaults
	// (auto_suspend=86400, auto_resume=true, min_cluster_count=1, max_cluster_count=1). If a config value
	// equals the current Snowflake value, IgnoreChangeToCurrentSnowflakeValueInShow correctly suppresses it
	// as a no-op and the attribute is never written to state, so the SET must be a genuine change.
	// auto_suspend must be >= 86400 for interactive warehouses, so we use a larger value to force a change.
	withOptionals := model.WarehouseInteractiveWithId(warehouseId).
		WithMaxClusterCount(3).
		WithMinClusterCount(2).
		WithAutoSuspend(172800).
		WithAutoResume(r.BooleanFalse).
		WithComment(comment).
		WithMaxConcurrencyLevel(8)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseInteractiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasWarehouseTypeString(string(sdk.WarehouseTypeInteractive)).
			HasNoWarehouseSize().
			HasAutoSuspendString(r.IntDefaultString).
			HasCommentEmpty().
			HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
		resourceshowoutputassert.WarehouseShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasStateNotEmpty(),
	}

	withOptionalsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseInteractiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasMaxClusterCount(3).
			HasMinClusterCount(2).
			HasAutoSuspend(172800).
			HasAutoResumeString(r.BooleanFalse).
			HasCommentString(comment).
			HasMaxConcurrencyLevel(8),
		resourceshowoutputassert.WarehouseShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// import after minimal config.
			// Optional fields left out of config are reconciled from SHOW during the import Read, so
			// they will not match the minimal state and are ignored here (same pattern as the adaptive resource).
			// show_output carries volatile runtime counters (running, queued, ...) that can change between
			// apply and import, so it is not verified here (it is asserted in the create/update steps).
			{
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"warehouse_size", "auto_suspend", "max_cluster_count", "min_cluster_count", "auto_resume", "show_output"},
			},
			// set all optional fields
			{
				Config: accconfig.FromModels(t, withOptionals),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, withOptionalsAssertions...),
			},
			// import after setting optional fields. warehouse_size is not in config, so it is reconciled
			// from SHOW on import and ignored here, as is the volatile show_output.
			{
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"warehouse_size", "show_output"},
			},
			// unset all optional fields
			{
				Config: accconfig.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_CompleteUseCase(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	fallback, fallbackCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(fallbackCleanup)

	table, tableCleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(tableCleanup)

	complete := model.WarehouseInteractiveWithId(warehouseId).
		WithWarehouseSize(string(sdk.WarehouseSizeSmall)).
		WithMaxClusterCount(2).
		WithMinClusterCount(1).
		WithAutoSuspend(86400).
		WithAutoResume(r.BooleanTrue).
		WithInitiallySuspended(true).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithFallbackWarehouse(fallback.ID().Name()).
		WithComment(comment).
		WithMaxConcurrencyLevel(8).
		WithStatementQueuedTimeoutInSeconds(30).
		WithStatementTimeoutInSeconds(45).
		WithTables(table.FullyQualifiedName())

	ref := complete.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasNameString(warehouseId.Name()).
						HasWarehouseTypeString(string(sdk.WarehouseTypeInteractive)).
						HasWarehouseSizeString(string(sdk.WarehouseSizeSmall)).
						HasMaxClusterCount(2).
						HasMinClusterCount(1).
						HasAutoSuspend(86400).
						HasAutoResumeString(r.BooleanTrue).
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasFallbackWarehouseString(fallback.ID().Name()).
						HasCommentString(comment).
						HasMaxConcurrencyLevel(8).
						HasStatementQueuedTimeoutInSeconds(30).
						HasStatementTimeoutInSeconds(45).
						HasTables(table.FullyQualifiedName()).
						HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
					resourceshowoutputassert.WarehouseShowOutput(t, ref).
						HasName(warehouseId.Name()).
						HasComment(comment).
						HasStateNotEmpty(),
				),
			},
			// import and verify state matches.
			// initially_suspended is a create-only field that is not read back from Snowflake, and show_output
			// carries volatile runtime counters, so both are ignored here.
			{
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initially_suspended", "show_output"},
			},
		},
	})
}

func TestAcc_WarehouseInteractive_WarehouseSize(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	small := model.WarehouseInteractiveWithId(warehouseId).
		WithWarehouseSize(string(sdk.WarehouseSizeSmall))
	medium := model.WarehouseInteractiveWithId(warehouseId).
		WithWarehouseSize(string(sdk.WarehouseSizeMedium))
	noSize := model.WarehouseInteractiveWithId(warehouseId)

	ref := small.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			// create with a concrete size
			{
				Config: accconfig.FromModels(t, small),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasWarehouseSizeString(string(sdk.WarehouseSizeSmall)),
				),
			},
			// changing the size updates in place: an interactive warehouse rejects an ALTER resize while
			// running, so the provider suspends it, applies the resize, and resumes it (no recreation).
			{
				Config: accconfig.FromModels(t, medium),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)),
				),
			},
			// removing the size recreates the warehouse: Snowflake has no UNSET WAREHOUSE_SIZE, so the
			// provider force-recreates (letting Snowflake apply its default size) instead.
			{
				Config: accconfig.FromModels(t, noSize),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasNoWarehouseSize(),
				),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_TablesDelta(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	table1, table1Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table1Cleanup)
	table2, table2Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table2Cleanup)
	table3, table3Cleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(table3Cleanup)

	modelWithTwoTables := model.WarehouseInteractiveWithId(warehouseId).
		WithTables(table1.FullyQualifiedName(), table2.FullyQualifiedName())
	// Drop table2, add table3 — should result in a single ADD + single DROP, no full replace.
	modelWithSwappedTable := model.WarehouseInteractiveWithId(warehouseId).
		WithTables(table1.FullyQualifiedName(), table3.FullyQualifiedName())

	ref := modelWithTwoTables.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithTwoTables),
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasTables(table1.FullyQualifiedName(), table2.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithSwappedTable),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseInteractiveResource(t, ref).
						HasTables(table1.FullyQualifiedName(), table3.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_Import_WrongWarehouseType(t *testing.T) {
	interactiveId := testClient().Ids.RandomAccountObjectIdentifier()
	regularId := testClient().Ids.RandomAccountObjectIdentifier()

	// Create a regular (non-interactive) warehouse outside of Terraform to use as the import target.
	_, regularCleanup := testClient().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(regularId))
	t.Cleanup(regularCleanup)

	interactiveModel := model.WarehouseInteractiveWithId(interactiveId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, interactiveModel),
			},
			{
				ResourceName:  interactiveModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: regularId.Name(),
				ExpectError:   regexp.MustCompile("is not an interactive warehouse"),
			},
		},
	})
}

func TestAcc_WarehouseInteractive_Validations(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	// NOTE: resource_monitor and fallback_warehouse also have validators, but both use
	// IsValidIdentifier[sdk.AccountObjectIdentifier], which intentionally skips validation
	// (see SNOW-1495079) because an account object identifier may legitimately contain dots.
	// There is therefore no config value that makes them fail at plan time, so they are not
	// exercised here.
	modelInvalidWarehouseSize := model.WarehouseInteractiveWithId(warehouseId).
		WithWarehouseSize("unknown")
	modelInvalidMaxClusterCount := model.WarehouseInteractiveWithId(warehouseId).
		WithMaxClusterCount(0)
	modelInvalidMinClusterCount := model.WarehouseInteractiveWithId(warehouseId).
		WithMinClusterCount(0)
	modelInvalidAutoSuspend := model.WarehouseInteractiveWithId(warehouseId).
		WithAutoSuspend(0)
	modelInvalidAutoResume := model.WarehouseInteractiveWithId(warehouseId).
		WithAutoResume("other")
	modelInvalidTables := model.WarehouseInteractiveWithId(warehouseId).
		WithTables("db.schema.table.column")
	modelInvalidMaxConcurrencyLevel := model.WarehouseInteractiveWithId(warehouseId).
		WithMaxConcurrencyLevel(0)
	modelInvalidStatementQueuedTimeout := model.WarehouseInteractiveWithId(warehouseId).
		WithStatementQueuedTimeoutInSeconds(-1)
	modelInvalidStatementTimeout := model.WarehouseInteractiveWithId(warehouseId).
		WithStatementTimeoutInSeconds(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidWarehouseSize),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid warehouse size: UNKNOWN"),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMaxClusterCount),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_cluster_count to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMinClusterCount),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_cluster_count to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoSuspend),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected auto_suspend to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidAutoResume),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected \[\{\{} auto_resume}] to be one of \["true" "false"], got other`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidTables),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Expected SchemaObjectIdentifier identifier type`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMaxConcurrencyLevel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_concurrency_level to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidStatementQueuedTimeout),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected statement_queued_timeout_in_seconds to be at least \(0\), got -1`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidStatementTimeout),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected statement_timeout_in_seconds to be in the range \(0 - 604800\), got -1`),
			},
		},
	})
}
