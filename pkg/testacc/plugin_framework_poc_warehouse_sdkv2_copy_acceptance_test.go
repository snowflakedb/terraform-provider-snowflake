// Tests in this file are all the current acceptance tests for the SDKv2 warehouse resource excluding migration tests.
// They were adjusted to verify Terraform Plugin Framework warehouse PoC resource implementation.
// Models used are the same but with the resource type replaced.
package testacc

import (
	"regexp"
	"strings"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func replaceWithWarehousePoCResourceType(t *testing.T, originalConfig string) string {
	replaced := strings.ReplaceAll(originalConfig, `resource "snowflake_warehouse"`, `resource "snowflake_warehouse_poc"`)
	t.Logf("Replaced config:\n%s", replaced)
	return replaced
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_BasicFlows(t *testing.T) {
	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	warehouseId2 := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	warehouseModel := model.Warehouse("test", warehouseId.Name()).WithComment(comment)
	warehouseModelRenamed := model.BasicWarehouseModel(warehouseId2, comment)
	warehouseModelRenamedFullWithoutParameters := model.WarehouseSnowflakeDefaultWithoutParameters(warehouseId2, comment)
	warehouseModelRenamedFullWithParameters := model.WarehouseSnowflakeDefaultWithoutParameters(warehouseId2, comment).
		WithMaxConcurrencyLevel(8).
		WithStatementQueuedTimeoutInSeconds(0).
		WithStatementTimeoutInSeconds(172800)
	warehouseModelRenamedFull := model.BasicWarehouseModel(warehouseId2, newComment).
		WithWarehouseTypeEnum(sdk.WarehouseTypeSnowparkOptimized).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium).
		WithMaxClusterCount(4).
		WithMinClusterCount(2).
		WithScalingPolicyEnum(sdk.ScalingPolicyEconomy).
		WithAutoSuspend(1200).
		WithAutoResume(r.BooleanFalse).
		WithInitiallySuspended(false).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithEnableQueryAcceleration(r.BooleanTrue).
		WithQueryAccelerationMaxScaleFactor(4).
		WithMaxConcurrencyLevel(4).
		WithStatementQueuedTimeoutInSeconds(5).
		WithStatementTimeoutInSeconds(86400)
	warehouseModelRenamedFullResourceMonitorInQuotes := model.BasicWarehouseModel(warehouseId2, newComment).
		WithWarehouseTypeEnum(sdk.WarehouseTypeSnowparkOptimized).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium).
		WithMaxClusterCount(4).
		WithMinClusterCount(2).
		WithScalingPolicyEnum(sdk.ScalingPolicyEconomy).
		WithAutoSuspend(1200).
		WithAutoResume(r.BooleanFalse).
		WithInitiallySuspended(false).
		WithResourceMonitor(resourceMonitor.ID().FullyQualifiedName()).
		WithEnableQueryAcceleration(r.BooleanTrue).
		WithQueryAccelerationMaxScaleFactor(4).
		WithMaxConcurrencyLevel(4).
		WithStatementQueuedTimeoutInSeconds(5).
		WithStatementTimeoutInSeconds(86400)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: replaceWithWarehousePoCResourceType(t, replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel))),
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, warehouseModel.ResourceReference()).
						HasNameString(warehouseId.Name()).
						HasNoWarehouseType().
						HasNoWarehouseSize().
						HasNoMaxClusterCount().
						HasNoMinClusterCount().
						HasNoScalingPolicy().
						HasAutoSuspendString(r.IntDefaultString).
						HasAutoResumeString(r.BooleanDefault).
						HasNoInitiallySuspended().
						HasNoResourceMonitor().
						HasCommentString(comment).
						HasEnableQueryAccelerationString(r.BooleanDefault).
						HasQueryAccelerationMaxScaleFactorString(r.IntDefaultString).
						HasMaxConcurrencyLevelString("8").
						HasStatementQueuedTimeoutInSecondsString("0").
						HasStatementTimeoutInSecondsString("172800").
						// alternatively extensions possible:
						HasDefaultMaxConcurrencyLevel().
						HasDefaultStatementQueuedTimeoutInSeconds().
						HasDefaultStatementTimeoutInSeconds().
						// alternatively extension possible
						HasAllDefault(),
					resourceshowoutputassert.WarehouseShowOutput(t, warehouseModel.ResourceReference()).
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMaxClusterCount(1).
						HasMinClusterCount(1).
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasResourceMonitor(sdk.AccountObjectIdentifier{}).
						HasComment(comment).
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8),
					resourceparametersassert.WarehouseResourceParameters(t, warehouseModel.ResourceReference()).
						HasMaxConcurrencyLevel(8).
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementTimeoutInSeconds(172800).
						// alternatively extensions possible:
						HasDefaultMaxConcurrencyLevel().
						HasDefaultStatementQueuedTimeoutInSeconds().
						HasDefaultStatementTimeoutInSeconds(),
					objectassert.Warehouse(t, warehouseId).
						HasName(warehouseId.Name()).
						HasState(sdk.WarehouseStateStarted).
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMaxClusterCount(1).
						HasMinClusterCount(1).
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasResourceMonitor(sdk.AccountObjectIdentifier{}).
						HasComment(comment).
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8),
					objectparametersassert.WarehouseParameters(t, warehouseId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					// we can still use normal checks
					assert.Check(resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "name", warehouseId.Name())),
					assert.Check(resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "fully_qualified_name", warehouseId.FullyQualifiedName())),
				),
			},
			// IMPORT after empty config (in this method, most of the attributes will be filled with the defaults acquired from Snowflake)
			{
				ResourceName: warehouseModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(warehouseId), "name", warehouseId.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(warehouseId), "fully_qualified_name", warehouseId.FullyQualifiedName())),
					resourceassert.ImportedWarehouseResource(t, helpers.EncodeResourceIdentifier(warehouseId)).
						HasNameString(warehouseId.Name()).
						HasWarehouseTypeString(string(sdk.WarehouseTypeStandard)).
						HasWarehouseSizeString(string(sdk.WarehouseSizeXSmall)).
						HasMaxClusterCountString("1").
						HasMinClusterCountString("1").
						HasScalingPolicyString(string(sdk.ScalingPolicyStandard)).
						HasAutoSuspendString("600").
						HasAutoResumeString("true").
						HasResourceMonitorString("").
						HasCommentString(comment).
						HasEnableQueryAccelerationString("false").
						HasQueryAccelerationMaxScaleFactorString("8").
						HasDefaultMaxConcurrencyLevel().
						HasDefaultStatementQueuedTimeoutInSeconds().
						HasDefaultStatementTimeoutInSeconds(),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(warehouseId)),
					resourceparametersassert.ImportedWarehouseResourceParameters(t, helpers.EncodeResourceIdentifier(warehouseId)).
						HasMaxConcurrencyLevel(8).
						HasMaxConcurrencyLevelLevel("").
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementQueuedTimeoutInSecondsLevel("").
						HasStatementTimeoutInSeconds(172800).
						HasStatementTimeoutInSecondsLevel(""),
					objectassert.Warehouse(t, warehouseId).
						HasName(warehouseId.Name()).
						HasState(sdk.WarehouseStateStarted).
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMaxClusterCount(1).
						HasMinClusterCount(1).
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasResourceMonitor(sdk.AccountObjectIdentifier{}).
						HasComment(comment).
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8),
					objectparametersassert.WarehouseParameters(t, warehouseId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
				),
			},
			// RENAME
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamed)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelRenamed.ResourceReference(), "name", warehouseId2.Name()),
					resource.TestCheckResourceAttr(warehouseModelRenamed.ResourceReference(), "fully_qualified_name", warehouseId2.FullyQualifiedName()),
				),
			},
			// Change config but use defaults for every attribute (but not the parameters) - expect no changes (because these are already SF values)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullWithoutParameters)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelRenamedFullWithoutParameters.ResourceReference(), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// add parameters - update expected (different level even with same values)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullWithParameters)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelRenamedFullWithParameters.ResourceReference(), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),

						// this is this only situation in which there will be a strange output in the plan
						planchecks.ExpectComputed(warehouseModelRenamedFullWithParameters.ResourceReference(), "max_concurrency_level", true),
						planchecks.ExpectComputed(warehouseModelRenamedFullWithParameters.ResourceReference(), "statement_queued_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModelRenamedFullWithParameters.ResourceReference(), "statement_timeout_in_seconds", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// no changes in the attributes
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "warehouse_type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "warehouse_size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "max_cluster_count", "1"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "min_cluster_count", "1"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "auto_suspend", "600"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "auto_resume", "true"),
					resource.TestCheckNoResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "initially_suspended"),
					resource.TestCheckNoResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "resource_monitor"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "comment", comment),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "query_acceleration_max_scale_factor", "8"),

					// parameters have the same values...
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					// ... but are set on different level
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.max_concurrency_level.0.value", "8"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.max_concurrency_level.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFullWithParameters.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// CHANGE PROPERTIES (normal and parameters)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFull)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelRenamedFull.ResourceReference(), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),

						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "warehouse_size", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseSizeXSmall)), sdk.String(string(sdk.WarehouseSizeMedium))),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "max_cluster_count", tfjson.ActionUpdate, sdk.String("1"), sdk.String("4")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "min_cluster_count", tfjson.ActionUpdate, sdk.String("1"), sdk.String("2")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "scaling_policy", tfjson.ActionUpdate, sdk.String(string(sdk.ScalingPolicyStandard)), sdk.String(string(sdk.ScalingPolicyEconomy))),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String("600"), sdk.String("1200")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "enable_query_acceleration", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("8"), sdk.String("4")),

						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "max_concurrency_level", tfjson.ActionUpdate, sdk.String("8"), sdk.String("4")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), sdk.String("5")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "warehouse_size", string(sdk.WarehouseSizeMedium)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "max_cluster_count", "4"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "min_cluster_count", "2"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "scaling_policy", string(sdk.ScalingPolicyEconomy)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "auto_suspend", "1200"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "auto_resume", "false"),
					resource.TestCheckNoResourceAttr(warehouseModelRenamedFull.ResourceReference(), "initially_suspended"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "resource_monitor", resourceMonitor.ID().Name()),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "comment", newComment),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "enable_query_acceleration", "true"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "query_acceleration_max_scale_factor", "4"),

					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "max_concurrency_level", "4"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "statement_queued_timeout_in_seconds", "5"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "statement_timeout_in_seconds", "86400"),

					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.max_concurrency_level.0.value", "4"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.max_concurrency_level.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.value", "5"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change resource monitor - wrap in quotes (no change expected)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullResourceMonitorInQuotes)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// CHANGE max_concurrency_level EXTERNALLY (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318)
			{
				Config:    replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFull)),
				PreConfig: func() { testClient().Warehouse.UpdateMaxConcurrencyLevel(t, warehouseId2, 10) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(warehouseModelRenamedFull.ResourceReference(), "max_concurrency_level", sdk.String("4"), sdk.String("10")),
						planchecks.ExpectChange(warehouseModelRenamedFull.ResourceReference(), "max_concurrency_level", tfjson.ActionUpdate, sdk.String("10"), sdk.String("4")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "name", warehouseId2.Name()),
					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "max_concurrency_level", "4"),

					resource.TestCheckResourceAttr(warehouseModelRenamedFull.ResourceReference(), "parameters.0.max_concurrency_level.0.value", "4"),
				),
			},
			// IMPORT
			{
				ResourceName:      warehouseModelRenamedFull.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_WarehouseType(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelStandard := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium).
		WithWarehouseTypeEnum(sdk.WarehouseTypeStandard)
	warehouseModelSnowparkOptimized := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium).
		WithWarehouseTypeEnum(sdk.WarehouseTypeSnowparkOptimized)
	warehouseModelNoType := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium)
	warehouseModelSnowparkOptimizedLowercase := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium).
		WithWarehouseType(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// set up with concrete type
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelStandard)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelStandard.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelStandard.ResourceReference(), "warehouse_type", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectComputed(warehouseModelStandard.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelStandard.ResourceReference(), "warehouse_type", string(sdk.WarehouseTypeStandard))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelStandard.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelStandard.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeStandard))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when type in config
			{
				ResourceName: warehouseModelStandard.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_type", string(sdk.WarehouseTypeStandard)),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.type", string(sdk.WarehouseTypeStandard)),
				),
			},
			// change type in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSnowparkOptimized)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelSnowparkOptimized.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelSnowparkOptimized.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectComputed(warehouseModelSnowparkOptimized.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimized.ResourceReference(), "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimized.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimized.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeSnowparkOptimized))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeSnowparkOptimized),
				),
			},
			// remove type from config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoType)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModelNoType.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(warehouseModelNoType.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelNoType.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
						planchecks.ExpectComputed(warehouseModelNoType.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "warehouse_type", "")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeStandard))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// add config (lower case)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSnowparkOptimizedLowercase)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))),
						planchecks.ExpectComputed(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), "warehouse_type", strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSnowparkOptimizedLowercase.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeSnowparkOptimized))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeSnowparkOptimized),
				),
			},
			// remove type from config but update warehouse externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					testClient().Warehouse.UpdateWarehouseType(t, id, sdk.WarehouseTypeStandard)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoType)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelNoType.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelNoType.ResourceReference(), "warehouse_type", sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized))), sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectDrift(warehouseModelNoType.ResourceReference(), "show_output.0.type", sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectChange(warehouseModelNoType.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), nil),
						planchecks.ExpectComputed(warehouseModelNoType.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "warehouse_type", "")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeStandard))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// change the type externally
			{
				PreConfig: func() {
					// we change the type to the type different from default, expecting action
					testClient().Warehouse.UpdateWarehouseType(t, id, sdk.WarehouseTypeSnowparkOptimized)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoType)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelNoType.ResourceReference(), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelNoType.ResourceReference(), "warehouse_type", nil, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectDrift(warehouseModelNoType.ResourceReference(), "show_output.0.type", sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectChange(warehouseModelNoType.ResourceReference(), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
						planchecks.ExpectComputed(warehouseModelNoType.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "warehouse_type", "")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoType.ResourceReference(), "show_output.0.type", string(sdk.WarehouseTypeStandard))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: warehouseModelNoType.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_type", string(sdk.WarehouseTypeStandard)),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.type", string(sdk.WarehouseTypeStandard)),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_WarehouseSizes(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelSmall := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeSmall)
	warehouseModelMedium := model.Warehouse("test", id.Name()).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium)
	warehouseModelNoSize := model.Warehouse("test", id.Name())
	warehouseModelSmallLowercase := model.Warehouse("test", id.Name()).
		WithWarehouseSize(strings.ToLower(string(sdk.WarehouseSizeSmall)))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// set up with concrete size
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSmall)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelSmall.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelSmall.ResourceReference(), "warehouse_size", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectComputed(warehouseModelSmall.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmall.ResourceReference(), "warehouse_size", string(sdk.WarehouseSizeSmall))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmall.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmall.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeSmall),
				),
			},
			// import when size in config
			{
				ResourceName: warehouseModelSmall.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_size", string(sdk.WarehouseSizeSmall)),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.size", string(sdk.WarehouseSizeSmall)),
				),
			},
			// change size in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelMedium)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelMedium.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelMedium.ResourceReference(), "warehouse_size", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseSizeSmall)), sdk.String(string(sdk.WarehouseSizeMedium))),
						planchecks.ExpectComputed(warehouseModelMedium.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelMedium.ResourceReference(), "warehouse_size", string(sdk.WarehouseSizeMedium))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelMedium.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelMedium.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeMedium))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeMedium),
				),
			},
			// remove size from config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoSize)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModelNoSize.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.PrintPlanDetails(warehouseModelNoSize.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelNoSize.ResourceReference(), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeMedium)), nil),
						planchecks.ExpectComputed(warehouseModelNoSize.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(warehouseModelNoSize.ResourceReference(), "warehouse_size")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeXSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeXSmall),
				),
			},
			// add config (lower case)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSmallLowercase)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelSmallLowercase.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelSmallLowercase.ResourceReference(), "warehouse_size", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall)))),
						planchecks.ExpectComputed(warehouseModelSmallLowercase.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmallLowercase.ResourceReference(), "warehouse_size", strings.ToLower(string(sdk.WarehouseSizeSmall)))),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmallLowercase.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelSmallLowercase.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeSmall),
				),
			},
			// remove size from config but update warehouse externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					testClient().Warehouse.UpdateWarehouseSize(t, id, sdk.WarehouseSizeXSmall)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoSize)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelNoSize.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelNoSize.ResourceReference(), "warehouse_size", sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall))), sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectDrift(warehouseModelNoSize.ResourceReference(), "show_output.0.size", sdk.String(string(sdk.WarehouseSizeSmall)), sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectChange(warehouseModelNoSize.ResourceReference(), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeXSmall)), nil),
						planchecks.ExpectComputed(warehouseModelNoSize.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(warehouseModelNoSize.ResourceReference(), "warehouse_size")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeXSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeXSmall),
				),
			},
			// change the size externally
			{
				PreConfig: func() {
					// we change the size to the size different from default, expecting action
					testClient().Warehouse.UpdateWarehouseSize(t, id, sdk.WarehouseSizeSmall)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoSize)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelNoSize.ResourceReference(), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelNoSize.ResourceReference(), "warehouse_size", nil, sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectDrift(warehouseModelNoSize.ResourceReference(), "show_output.0.size", sdk.String(string(sdk.WarehouseSizeXSmall)), sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectChange(warehouseModelNoSize.ResourceReference(), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeSmall)), nil),
						planchecks.ExpectComputed(warehouseModelNoSize.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(warehouseModelNoSize.ResourceReference(), "warehouse_size")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelNoSize.ResourceReference(), "show_output.0.size", string(sdk.WarehouseSizeXSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeXSmall),
				),
			},
			// import when no size in config
			{
				ResourceName: warehouseModelNoSize.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_size", string(sdk.WarehouseSizeXSmall)),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelInvalidType := model.Warehouse("test", id.Name()).
		WithWarehouseType("unknown")
	warehouseModelInvalidSize := model.Warehouse("test", id.Name()).
		WithWarehouseSize("SMALLa")
	warehouseModelInvalidMaxClusterCount := model.Warehouse("test", id.Name()).
		WithMaxClusterCount(0)
	warehouseModelInvalidMinClusterCount := model.Warehouse("test", id.Name()).
		WithMinClusterCount(0)
	warehouseModelInvalidScalingPolicy := model.Warehouse("test", id.Name()).
		WithScalingPolicy("unknown")
	warehouseModelInvalidAutoResume := model.Warehouse("test", id.Name()).
		WithAutoResume("other")
	warehouseModelInvalidMaxConcurrencyLevel := model.Warehouse("test", id.Name()).
		WithMaxConcurrencyLevel(-2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidType)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid warehouse type: unknown"),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidSize)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid warehouse size: SMALLa"),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidMaxClusterCount)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_cluster_count to be at least \(1\), got 0`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidMinClusterCount)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_cluster_count to be at least \(1\), got 0`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidScalingPolicy)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid scaling policy: unknown"),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidAutoResume)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected \[\{\{} auto_resume}] to be one of \["true" "false"], got other`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidMaxConcurrencyLevel)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_concurrency_level to be at least \(1\), got -2`),
			},
		},
	})
}

// Just for the experimental purposes
func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_ValidateDriftForCurrentWarehouse(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	secondId := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel := model.Warehouse("test", id.Name())
	secondWarehouseModel := model.Warehouse("test2", secondId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.is_current", "true"),
				),
			},
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel, secondWarehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModel.ResourceReference(), plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(secondWarehouseModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.is_current", "true"),

					resource.TestCheckResourceAttr(secondWarehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(secondWarehouseModel.ResourceReference(), "show_output.0.is_current", "true"),
				),
			},
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel, secondWarehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(warehouseModel.ResourceReference(), "show_output.0.is_current", sdk.String("true"), sdk.String("false")),
						plancheck.ExpectResourceAction(warehouseModel.ResourceReference(), plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(secondWarehouseModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.is_current", "false"),
				),
			},
		},
	})
}

// TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_AutoResume validates behavior for falling back to Snowflake default for boolean attribute
func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_AutoResume(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelWithoutAutoResume := model.Warehouse("test", id.Name())
	warehouseModelAutoResumeTrue := model.Warehouse("test", id.Name()).WithAutoResume(r.BooleanTrue)
	warehouseModelAutoResumeFalse := model.Warehouse("test", id.Name()).WithAutoResume(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// set up with auto resume set in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoResumeTrue)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelAutoResumeTrue.ResourceReference(), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelAutoResumeTrue.ResourceReference(), "auto_resume", tfjson.ActionCreate, nil, sdk.String("true")),
						planchecks.ExpectComputed(warehouseModelAutoResumeTrue.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeTrue.ResourceReference(), "auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeTrue.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeTrue.ResourceReference(), "show_output.0.auto_resume", "true")),
					objectassert.Warehouse(t, id).HasAutoResume(true),
				),
			},
			// import when type in config
			{
				ResourceName: warehouseModelAutoResumeTrue.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_resume", "true"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.auto_resume", "true"),
				),
			},
			// change value in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoResumeFalse)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelAutoResumeFalse.ResourceReference(), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelAutoResumeFalse.ResourceReference(), "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectComputed(warehouseModelAutoResumeFalse.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeFalse.ResourceReference(), "auto_resume", "false")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeFalse.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoResumeFalse.ResourceReference(), "show_output.0.auto_resume", "false")),
					objectassert.Warehouse(t, id).HasAutoResume(false),
				),
			},
			// remove type from config (expecting non-empty plan because we do not know the default)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoResume)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModelWithoutAutoResume.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", tfjson.ActionUpdate, sdk.String("false"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectComputed(warehouseModelWithoutAutoResume.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", r.BooleanDefault)),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "show_output.0.auto_resume", "true")),
					objectassert.Warehouse(t, id).HasAutoResume(true),
				),
			},
			// change auto resume externally
			{
				PreConfig: func() {
					// we change the auto resume to the type different from default, expecting action
					testClient().Warehouse.UpdateAutoResume(t, id, false)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoResume)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", sdk.String(r.BooleanDefault), sdk.String("false")),
						planchecks.ExpectDrift(warehouseModelWithoutAutoResume.ResourceReference(), "show_output.0.auto_resume", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", tfjson.ActionUpdate, sdk.String("false"), sdk.String(r.BooleanDefault)),
						planchecks.ExpectComputed(warehouseModelWithoutAutoResume.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "auto_resume", r.BooleanDefault)),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoResume.ResourceReference(), "show_output.0.auto_resume", "true")),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: warehouseModelWithoutAutoResume.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_resume", "true"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.auto_resume", "true"),
				),
			},
		},
	})
}

// TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_AutoSuspend validates behavior for falling back to Snowflake default for the integer attribute
func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_AutoSuspend(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModelWithoutAutoSuspend := model.Warehouse("test", id.Name())
	warehouseModelAutoSuspend1200 := model.Warehouse("test", id.Name()).WithAutoSuspend(1200)
	warehouseModelAutoSuspend600 := model.Warehouse("test", id.Name()).WithAutoSuspend(600)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// set up with auto suspend set in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoSuspend1200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelAutoSuspend1200.ResourceReference(), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelAutoSuspend1200.ResourceReference(), "auto_suspend", tfjson.ActionCreate, nil, sdk.String("1200")),
						planchecks.ExpectComputed(warehouseModelAutoSuspend1200.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend1200.ResourceReference(), "auto_suspend", "1200")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend1200.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend1200.ResourceReference(), "show_output.0.auto_suspend", "1200")),
					objectassert.Warehouse(t, id).HasAutoSuspend(1200),
				),
			},
			// import when auto suspend in config
			{
				ResourceName: warehouseModelAutoSuspend1200.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "1200"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.auto_suspend", "1200"),
				),
			},
			// change value in config to Snowflake default
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoSuspend600)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelAutoSuspend600.ResourceReference(), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelAutoSuspend600.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String("1200"), sdk.String("600")),
						planchecks.ExpectComputed(warehouseModelAutoSuspend600.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend600.ResourceReference(), "auto_suspend", "600")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend600.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelAutoSuspend600.ResourceReference(), "show_output.0.auto_suspend", "600")),
					objectassert.Warehouse(t, id).HasAutoSuspend(600),
				),
			},
			// remove auto suspend from config (expecting non-empty plan because we do not know the default)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoSuspend)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(warehouseModelWithoutAutoSuspend.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String("600"), sdk.String(r.IntDefaultString)),
						planchecks.ExpectComputed(warehouseModelWithoutAutoSuspend.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", r.IntDefaultString)),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "show_output.0.auto_suspend", "600")),
					objectassert.Warehouse(t, id).HasAutoSuspend(600),
				),
			},
			// change auto suspend externally
			{
				PreConfig: func() {
					// we change the max cluster count to the type different from default, expecting action
					testClient().Warehouse.UpdateAutoSuspend(t, id, 2400)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoSuspend)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", sdk.String(r.IntDefaultString), sdk.String("2400")),
						planchecks.ExpectDrift(warehouseModelWithoutAutoSuspend.ResourceReference(), "show_output.0.auto_suspend", sdk.String("600"), sdk.String("2400")),
						planchecks.ExpectChange(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String("2400"), sdk.String(r.IntDefaultString)),
						planchecks.ExpectComputed(warehouseModelWithoutAutoSuspend.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "auto_suspend", r.IntDefaultString)),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "show_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(warehouseModelWithoutAutoSuspend.ResourceReference(), "show_output.0.auto_suspend", "600")),
					objectassert.Warehouse(t, id).HasAutoSuspend(600),
				),
			},
			// import when no type in config
			{
				ResourceName: warehouseModelWithoutAutoSuspend.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "600"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.auto_suspend", "600"),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_ZeroValues(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel := model.Warehouse("test", id.Name())
	warehouseModelWithAllValidZeroValues := model.Warehouse("test", id.Name()).
		WithAutoSuspend(0).
		WithQueryAccelerationMaxScaleFactor(0).
		WithStatementQueuedTimeoutInSeconds(0).
		WithStatementTimeoutInSeconds(0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// create with valid "zero" values
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithAllValidZeroValues)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "query_acceleration_max_scale_factor", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_queued_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectComputed(warehouseModelWithAllValidZeroValues.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "query_acceleration_max_scale_factor", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_timeout_in_seconds", "0"),

					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.0.auto_suspend", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.0.query_acceleration_max_scale_factor", "0"),

					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String("0"), sdk.String(r.IntDefaultString)),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("0"), sdk.String(r.IntDefaultString)),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), "statement_queued_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "auto_suspend", r.IntDefaultString),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "query_acceleration_max_scale_factor", r.IntDefaultString),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.query_acceleration_max_scale_factor", "8"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.level", ""),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithAllValidZeroValues)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", tfjson.ActionUpdate, sdk.String(r.IntDefaultString), sdk.String("0")),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String(r.IntDefaultString), sdk.String("0")),
						planchecks.ExpectComputed(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_queued_timeout_in_seconds", true),
						planchecks.ExpectChange(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), sdk.String("0")),
						planchecks.ExpectComputed(warehouseModelWithAllValidZeroValues.ResourceReference(), r.ShowOutputAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "auto_suspend", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "query_acceleration_max_scale_factor", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "statement_timeout_in_seconds", "0"),

					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.0.auto_suspend", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "show_output.0.query_acceleration_max_scale_factor", "0"),

					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr(warehouseModelWithAllValidZeroValues.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// import zero values
			{
				ResourceName: warehouseModelWithAllValidZeroValues.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),

					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "query_acceleration_max_scale_factor", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_queued_timeout_in_seconds", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "0"),

					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.auto_suspend", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "show_output.0.query_acceleration_max_scale_factor", "0"),

					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.value", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_Parameter(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel := model.Warehouse("test", id.Name())
	warehouseModelWithStatementTimeoutInSeconds86400 := model.Warehouse("test", id.Name()).WithStatementTimeoutInSeconds(86400)
	warehouseModelWithStatementTimeoutInSeconds43200 := model.Warehouse("test", id.Name()).WithStatementTimeoutInSeconds(43200)
	warehouseModelWithStatementTimeoutInSeconds172800 := model.Warehouse("test", id.Name()).WithStatementTimeoutInSeconds(172800)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// create with setting one param
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds86400)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("86400")),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "statement_timeout_in_seconds", "86400"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// do not make any change (to check if there is no drift)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds86400)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// import when param in config
			{
				ResourceName: warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds43200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change param value on account - expect no changes
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
					revert := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "86400")
					t.Cleanup(revert)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds43200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", tfjson.ActionNoop, sdk.String("43200"), sdk.String("43200")),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value externally
			{
				PreConfig: func() {
					// clean after previous step
					testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					// update externally
					testClient().Warehouse.UpdateStatementTimeoutInSeconds(t, id, 86400)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds43200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectDrift(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", sdk.String("43200"), sdk.String("86400")),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value on account to the value from config (but on different level)
			{
				PreConfig: func() {
					testClient().Warehouse.UnsetStatementTimeoutInSeconds(t, id)
					testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43200")
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds43200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("43200"), nil),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					// clean after previous step
					testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("43200"), nil),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// import when param not in config (snowflake default)
			{
				ResourceName: warehouseModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "172800"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// change the param value in config to snowflake default (expecting action because of the different level)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds172800)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), nil),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "statement_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), nil),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// change param value on account - change expected to be noop
			{
				PreConfig: func() {
					param := testClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
					revert := testClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "86400")
					t.Cleanup(revert)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectDrift(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", sdk.String("172800"), sdk.String("86400")),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", tfjson.ActionNoop, sdk.String("86400"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "86400"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// import when param not in config (set on account)
			{
				ResourceName: warehouseModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// change param value on warehouse
			{
				PreConfig: func() {
					testClient().Warehouse.UpdateStatementTimeoutInSeconds(t, id, 86400)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), nil),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", true),
						planchecks.ExpectComputed(warehouseModel.ResourceReference(), r.ParametersAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "86400"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// unset param on account
			{
				PreConfig: func() {
					testClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
				},
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectDrift(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", sdk.String("86400"), sdk.String("172800")),
						planchecks.ExpectDrift(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", sdk.String(string(sdk.ParameterTypeAccount)), sdk.String("")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
		},
	})
}

func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_InitiallySuspendedChangesPostCreation(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel := model.Warehouse("test", id.Name())
	warehouseModelWithInitiallySuspendedTrue := model.Warehouse("test", id.Name()).WithInitiallySuspended(true)
	warehouseModelWithInitiallySuspendedFalse := model.Warehouse("test", id.Name()).WithInitiallySuspended(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithInitiallySuspendedTrue)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedTrue.ResourceReference(), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedTrue.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedTrue.ResourceReference(), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithInitiallySuspendedFalse)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedFalse.ResourceReference(), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedFalse.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModelWithInitiallySuspendedFalse.ResourceReference(), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(warehouseModel.ResourceReference(), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
		},
	})
}
