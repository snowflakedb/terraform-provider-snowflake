//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_WarehouseAdaptive_BasicUseCase(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	newWarehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)
	newResourceMonitor, newResourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(newResourceMonitorCleanup)
	externalResourceMonitor, externalResourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(externalResourceMonitorCleanup)

	warehouseModel := model.WarehouseAdaptiveWithId(warehouseId)
	warehouseModelWithOptionals := model.WarehouseAdaptiveWithId(warehouseId).
		WithComment(comment).
		WithQueryThroughputMultiplier(2).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelLarge)).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithStatementQueuedTimeoutInSeconds(300).
		WithStatementTimeoutInSeconds(86400)
	warehouseModelWithZeroMultiplier := model.WarehouseAdaptiveWithId(warehouseId).
		WithQueryThroughputMultiplier(0)
	warehouseModelUpdated := model.WarehouseAdaptiveWithId(newWarehouseId).
		WithComment(newComment).
		WithQueryThroughputMultiplier(4).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelSmall)).
		WithResourceMonitor(newResourceMonitor.ID().Name()).
		WithStatementQueuedTimeoutInSeconds(600).
		WithStatementTimeoutInSeconds(43200)
	warehouseModelRenamedMinimal := model.WarehouseAdaptiveWithId(newWarehouseId)

	ref := warehouseModel.ResourceReference()
	externalQueryThroughputMultiplier := 10
	externalStatementTimeout := 99999
	externalStatementQueuedTimeout := 1200
	externalMaxQueryPerformanceLevel := sdk.MaxQueryPerformanceLevelMedium
	externalComment := random.Comment()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentEmpty().
			HasNoMaxQueryPerformanceLevel().
			HasQueryThroughputMultiplierString(r.IntDefaultString).
			HasResourceMonitorEmpty().
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800).
			HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasStateNotEmpty().
			HasCommentEmpty().
			HasResourceMonitorEmpty().
			HasOwnerNotEmpty().
			HasOwnerRoleTypeNotEmpty(),
	}

	withOptionalsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentString(comment).
			HasMaxQueryPerformanceLevelString(string(sdk.MaxQueryPerformanceLevelLarge)).
			HasQueryThroughputMultiplier(2).
			HasResourceMonitorString(resourceMonitor.ID().Name()).
			HasStatementQueuedTimeoutInSeconds(300).
			HasStatementTimeoutInSeconds(86400),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasComment(comment).
			HasResourceMonitor(resourceMonitor.ID()).
			HasQueryThroughputMultiplier(2),
	}

	updatedAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(newWarehouseId.Name()).
			HasCommentString(newComment).
			HasQueryThroughputMultiplier(4).
			HasResourceMonitorString(newResourceMonitor.ID().Name()).
			HasStatementQueuedTimeoutInSeconds(600).
			HasStatementTimeoutInSeconds(43200).
			HasFullyQualifiedNameString(newWarehouseId.FullyQualifiedName()),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(newWarehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasComment(newComment).
			HasResourceMonitor(newResourceMonitor.ID()).
			HasQueryThroughputMultiplier(4),
	}

	unsetAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(newWarehouseId.Name()).
			HasCommentEmpty().
			HasQueryThroughputMultiplierString(r.IntDefaultString).
			HasResourceMonitorEmpty().
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800).
			HasFullyQualifiedNameString(newWarehouseId.FullyQualifiedName()),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(newWarehouseId.Name()).
			HasCommentEmpty().
			HasResourceMonitorEmpty().
			HasQueryThroughputMultiplier(2),
	}

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModel),
				Check:                    assertThat(t, basicAssertions...),
			},
			// import - without optionals, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, warehouseModel),
				ResourceName:             warehouseModel.ResourceReference(),
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// set query_throughput_multiplier to 0 (explicit value, distinct from IntDefault sentinel -1)
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelWithZeroMultiplier),
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, ref).
						HasQueryThroughputMultiplier(0),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
						HasQueryThroughputMultiplier(0),
				),
			},
			// set all optional fields
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelWithOptionals),
				Check:                    assertThat(t, withOptionalsAssertions...),
			},
			// rename and update fields
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "name", tfjson.ActionUpdate, sdk.String(warehouseId.Name()), sdk.String(newWarehouseId.Name())),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), sdk.String(newComment)),
						planchecks.ExpectChange(ref, "query_throughput_multiplier", tfjson.ActionUpdate, sdk.String("2"), sdk.String("4")),
						planchecks.ExpectChange(ref, "max_query_performance_level", tfjson.ActionUpdate, sdk.String(string(sdk.MaxQueryPerformanceLevelLarge)), sdk.String(string(sdk.MaxQueryPerformanceLevelSmall))),
						planchecks.ExpectChange(ref, "resource_monitor", tfjson.ActionUpdate, sdk.String(resourceMonitor.ID().Name()), sdk.String(newResourceMonitor.ID().Name())),
						planchecks.ExpectChange(ref, "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("300"), sdk.String("600")),
						planchecks.ExpectChange(ref, "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// detect external changes to multiple fields
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelUpdated),
				PreConfig: func() {
					testClient().Warehouse.DropWarehouseFunc(t, newWarehouseId)()
					testClient().Warehouse.CreateAdaptiveWithRequest(
						t, sdk.NewCreateAdaptiveWarehouseRequest(newWarehouseId).
							WithComment(externalComment).
							WithQueryThroughputMultiplier(externalQueryThroughputMultiplier).
							WithStatementTimeoutInSeconds(externalStatementTimeout).
							WithStatementQueuedTimeoutInSeconds(externalStatementQueuedTimeout).
							WithMaxQueryPerformanceLevel(externalMaxQueryPerformanceLevel).
							WithResourceMonitor(externalResourceMonitor.ID()),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "query_throughput_multiplier", tfjson.ActionUpdate, sdk.String("10"), sdk.String("4")),
						planchecks.ExpectChange(ref, "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("99999"), sdk.String("43200")),
						planchecks.ExpectChange(ref, "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("1200"), sdk.String("600")),
						planchecks.ExpectChange(ref, "max_query_performance_level", tfjson.ActionUpdate, sdk.String(string(sdk.MaxQueryPerformanceLevelMedium)), sdk.String(string(sdk.MaxQueryPerformanceLevelSmall))),
						planchecks.ExpectChange(ref, "resource_monitor", tfjson.ActionUpdate, sdk.String(externalResourceMonitor.ID().Name()), sdk.String(newResourceMonitor.ID().Name())),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(newComment)),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// unset optional fields back to defaults
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelRenamedMinimal),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(newComment), nil),
						planchecks.ExpectChange(ref, "query_throughput_multiplier", tfjson.ActionUpdate, sdk.String("4"), sdk.String(r.IntDefaultString)),
						planchecks.ExpectChange(ref, "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("600"), nil),
						planchecks.ExpectChange(ref, "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("43200"), nil),
						planchecks.ExpectChange(ref, "max_query_performance_level", tfjson.ActionUpdate, sdk.String(string(sdk.MaxQueryPerformanceLevelSmall)), nil),
						planchecks.ExpectChange(ref, "resource_monitor", tfjson.ActionUpdate, sdk.String(newResourceMonitor.ID().Name()), nil),
					},
				},
				Check: assertThat(t, unsetAssertions...),
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_CompleteUseCase(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	warehouseModelComplete := model.WarehouseAdaptiveWithId(warehouseId).
		WithComment(comment).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelLarge)).
		WithQueryThroughputMultiplier(3).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithStatementQueuedTimeoutInSeconds(300).
		WithStatementTimeoutInSeconds(86400)

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				Config:                   accconfig.FromModels(t, warehouseModelComplete),
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, warehouseModelComplete.ResourceReference()).
						HasNameString(warehouseId.Name()).
						HasCommentString(comment).
						HasMaxQueryPerformanceLevelString(string(sdk.MaxQueryPerformanceLevelLarge)).
						HasQueryThroughputMultiplier(3).
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasStatementQueuedTimeoutInSeconds(300).
						HasStatementTimeoutInSeconds(86400).
						HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, warehouseModelComplete.ResourceReference()).
						HasName(warehouseId.Name()).
						HasType(sdk.WarehouseTypeAdaptive).
						HasComment(comment).
						HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
						HasQueryThroughputMultiplier(3).
						HasResourceMonitor(resourceMonitor.ID()),
					resourceparametersassert.WarehouseAdaptiveResourceParameters(t, warehouseModelComplete.ResourceReference()).
						HasStatementQueuedTimeoutInSeconds(300).
						HasStatementTimeoutInSeconds(86400),
				),
			},
			// import and verify state matches
			{
				ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
				ResourceName:             warehouseModelComplete.ResourceReference(),
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"show_output.0.running"},
			},
			// import - complete, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, warehouseModelComplete),
				ResourceName:             warehouseModelComplete.ResourceReference(),
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)

	basicModel := model.WarehouseAdaptiveWithId(id)
	ref := basicModel.ResourceReference()

	_, cleanup := testClient().Warehouse.CreateAdaptiveWithRequest(t, sdk.NewCreateAdaptiveWarehouseRequest(id))
	t.Cleanup(cleanup)

	warehouse, err := testClient().Warehouse.ShowByIDAdaptive(t, id)
	require.NoError(t, err)
	require.NotNil(t, warehouse.MaxQueryPerformanceLevel)
	require.NotNil(t, warehouse.QueryThroughputMultiplier)
	maxQueryPerformanceLevel := string(*warehouse.MaxQueryPerformanceLevel)
	queryThroughputMultiplier := *warehouse.QueryThroughputMultiplier

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// Import - without optionals
			{
				Config:             accconfig.FromModels(t, basicModel),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedWarehouseAdaptiveResource(t, resourceId).
						HasName(id.Name()).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasWarehouseType(string(sdk.WarehouseTypeAdaptive)).
						HasCommentEmpty().
						HasMaxQueryPerformanceLevel(maxQueryPerformanceLevel).
						HasQueryThroughputMultiplier(queryThroughputMultiplier).
						HasResourceMonitorEmpty().
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementTimeoutInSeconds(172800),
				),
			},
			// Expect empty plan after the plan after import is applied
			{
				Config: accconfig.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "max_query_performance_level", tfjson.ActionUpdate, new(maxQueryPerformanceLevel), nil),
						planchecks.ExpectChange(ref, "query_throughput_multiplier", tfjson.ActionUpdate, new(strconv.Itoa(queryThroughputMultiplier)), new(r.IntDefaultString)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, ref).
						HasName(id.Name()).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasWarehouseType(string(sdk.WarehouseTypeAdaptive)).
						HasCommentEmpty().
						HasMaxQueryPerformanceLevelEmpty().
						HasQueryThroughputMultiplier(r.IntDefault).
						HasResourceMonitorEmpty().
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementTimeoutInSeconds(172800),
				),
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_Import_WrongWarehouseType(t *testing.T) {
	adaptiveId := testClient().Ids.RandomAccountObjectIdentifier()
	interactiveId := testClient().Ids.RandomAccountObjectIdentifier()

	// Create an interactive warehouse outside of Terraform to use as the import target.
	_, interactiveCleanup := testClient().Warehouse.CreateInteractiveWithRequest(t, sdk.NewCreateInteractiveWarehouseRequest(interactiveId))
	t.Cleanup(interactiveCleanup)

	adaptiveModel := model.WarehouseAdaptiveWithId(adaptiveId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// Create an adaptive warehouse to have a resource in state.
			{
				Config: accconfig.FromModels(t, adaptiveModel),
			},
			// Attempt to import an interactive warehouse via the adaptive resource — expects a type mismatch error.
			{
				ResourceName:  adaptiveModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: interactiveId.Name(),
				ExpectError:   regexp.MustCompile("is an interactive warehouse and cannot be converted to ADAPTIVE"),
			},
		},
	})
}

// Proves that an existing snowflake_warehouse can be migrated to snowflake_warehouse_adaptive in a single apply,
// by combining a removed block, the new resource, and an import block. See https://github.com/snowflakedb/terraform-provider-snowflake/issues/5201.
func TestAcc_WarehouseAdaptive_MigrateFromStandardWarehouseInSingleApply(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	standardModel := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium)
	adaptiveModel := model.WarehouseAdaptiveWithId(id).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelLarge))

	standardRef := standardModel.ResourceReference()
	adaptiveRef := adaptiveModel.ResourceReference()

	migrationConfig := accconfig.FromModels(t, adaptiveModel) + fmt.Sprintf(`
removed {
  from = %[1]s

  lifecycle {
    destroy = false
  }
}

import {
  to = %[2]s
  id = %[3]q
}
`, standardRef, adaptiveRef, id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// removed blocks require Terraform 1.7 or later
			tfversion.RequireAbove(tfversion.Version1_7_0),
		},
		CheckDestroy: ComposeCheckDestroy(t, resources.Warehouse, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// Start from a standard warehouse managed by snowflake_warehouse.
			{
				Config: accconfig.FromModels(t, standardModel),
				Check: assertThat(
					t,
					objectassert.Warehouse(t, id).
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeMedium),
				),
			},
			// Forget the standard resource, adopt the adaptive one, and change the type - all in one apply.
			{
				Config: migrationConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// an update (not a recreation) means the warehouse and everything depending on it is preserved
						plancheck.ExpectResourceAction(adaptiveRef, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(adaptiveRef, "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeAdaptive))),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, adaptiveRef).
						HasWarehouseTypeString(string(sdk.WarehouseTypeAdaptive)).
						HasMaxQueryPerformanceLevelString(string(sdk.MaxQueryPerformanceLevelLarge)),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, adaptiveRef).
						HasType(sdk.WarehouseTypeAdaptive).
						HasState(sdk.WarehouseStateEnabled),
					objectassert.WarehouseAdaptive(t, id).
						HasType(sdk.WarehouseTypeAdaptive).
						HasState(sdk.WarehouseStateEnabled),
					objectassert.Warehouse(t, id).
						HasNoSize(),
				),
			},
			// The migration blocks are no longer needed.
			{
				Config: accconfig.FromModels(t, adaptiveModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_Validations(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelInvalidMaxQueryPerfLevel := model.WarehouseAdaptiveWithId(warehouseId).
		WithMaxQueryPerformanceLevel("unknown")
	warehouseModelInvalidQueryThroughput := model.WarehouseAdaptiveWithId(warehouseId).
		WithQueryThroughputMultiplier(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, warehouseModelInvalidMaxQueryPerfLevel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid max query performance level: UNKNOWN"),
			},
			{
				Config:      accconfig.FromModels(t, warehouseModelInvalidQueryThroughput),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected query_throughput_multiplier to be at least \(0\), got -1`),
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_ExternalTypeChange(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	warehouseModel := model.WarehouseAdaptiveWithId(id)
	ref := warehouseModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, warehouseModel),
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, ref).
						HasWarehouseTypeString(string(sdk.WarehouseTypeAdaptive)),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
						HasType(sdk.WarehouseTypeAdaptive),
				),
			},
			{
				PreConfig: func() {
					testClient().Warehouse.UpdateWarehouseType(t, id, sdk.WarehouseTypeStandard)
				},
				Config: accconfig.FromModels(t, warehouseModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeAdaptive))),
					},
				},
				Check: assertThat(
					t,
					resourceassert.WarehouseAdaptiveResource(t, ref).
						HasWarehouseTypeString(string(sdk.WarehouseTypeAdaptive)),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
						HasType(sdk.WarehouseTypeAdaptive),
				),
			},
		},
	})
}

func TestAcc_WarehouseAdaptive_ExternalTypeChangeToInteractive(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel := model.WarehouseAdaptiveWithId(id)
	ref := warehouseModel.ResourceReference()

	assertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasWarehouseType(string(sdk.WarehouseTypeAdaptive)),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasType(sdk.WarehouseTypeAdaptive),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: adaptiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, warehouseModel),
				Check:  assertThat(t, assertions...),
			},
			{
				PreConfig: func() {
					testClient().Warehouse.CreateInteractiveWithRequest(t,
						sdk.NewCreateInteractiveWarehouseRequest(id).WithOrReplace(true))
				},
				Config: accconfig.FromModels(t, warehouseModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "warehouse_type", new(string(sdk.WarehouseTypeAdaptive)), new(string(sdk.WarehouseTypeInteractive))),
						planchecks.ExpectChange(ref, "warehouse_type", tfjson.ActionDelete, new(string(sdk.WarehouseTypeInteractive)), nil),
					},
				},
				Check: assertThat(t, assertions...),
			},
		},
	})
}
