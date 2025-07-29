//go:build account_level_tests

// Tests in this file are all the current acceptance tests for the SDKv2 warehouse resource excluding migration and experimental tests.
// They were adjusted to verify Terraform Plugin Framework warehouse PoC resource implementation.
// Models used are the same, but with the resource type replaced.
// Assertions used are the same, but with the resource type replaced.
// Assertions using r.IntDefaultString or r.BooleanDefault were replaced (as such defaults are not needed in the Terraform Plugin Framework).
// Parameter values are not in state when the config does not contain them.
// Default parameter assertions can't be used because of above.
// WarehouseShowOutput assertions were removed or replaced with Snowflake object assertions.
// WarehouseResourceParameters assertions were removed or replaced with Snowflake parameters assertions.
// Default extensions were removed as they don't match.
// Expectations for tests utilizing IgnoreChangeToCurrentSnowflakeValueInShow were adjusted.
// Computed expectations for parameters were adjusted (ExpectComputed -> ExpectChange).
// Exchanged some old assertions with new ones.
// IgnoreAfterCreation is not implemented so assertions for initially_suspended were adjusted.
// Identifier suppression is not implemented, so adjusted steps with resource monitor.
// Some TestCheckResourceAttr were replaced with TestCheckNoResourceAttr as plugin framework handles null differently.
// Validation tests regex assertions were adjusted.
// Test verifying initially_suspended is skipped till the IgnoreAfterCreate is implemented.
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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
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

func replaceResourceReference(originalReference string) string {
	replaced := strings.ReplaceAll(originalReference, `snowflake_warehouse`, `snowflake_warehouse_poc`)
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
	warehouseModelRenamedFullWithParametersMediumSize := model.WarehouseSnowflakeDefaultWithoutParameters(warehouseId2, comment).
		WithMaxConcurrencyLevel(8).
		WithStatementQueuedTimeoutInSeconds(0).
		WithStatementTimeoutInSeconds(172800).
		WithWarehouseSizeEnum(sdk.WarehouseSizeMedium)
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
			// create with only required fields present in config
			{
				Config: replaceWithWarehousePoCResourceType(t, replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel))),
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, replaceResourceReference(warehouseModel.ResourceReference())).
						HasNameString(warehouseId.Name()).
						HasNoWarehouseType().
						HasNoWarehouseSize().
						HasNoMaxClusterCount().
						HasNoMinClusterCount().
						HasNoScalingPolicy().
						HasNoAutoSuspend().
						HasNoAutoResume().
						HasNoInitiallySuspended().
						HasNoResourceMonitor().
						HasCommentString(comment).
						HasNoEnableQueryAcceleration().
						HasNoQueryAccelerationMaxScaleFactor().
						HasNoMaxConcurrencyLevel().
						HasNoStatementQueuedTimeoutInSeconds().
						HasNoStatementTimeoutInSeconds(),
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
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "name", warehouseId.Name())),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "fully_qualified_name", warehouseId.FullyQualifiedName())),
				),
			},
			// IMPORT after empty config (in this method, most of the attributes will be filled with the defaults acquired from Snowflake)
			{
				ResourceName: replaceResourceReference(warehouseModel.ResourceReference()),
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
						HasNoMaxConcurrencyLevel().
						HasNoStatementQueuedTimeoutInSeconds().
						HasNoStatementTimeoutInSeconds(),
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
						plancheck.ExpectResourceAction(replaceResourceReference(warehouseModelRenamed.ResourceReference()), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelRenamed.ResourceReference()), "name", warehouseId2.Name()),
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelRenamed.ResourceReference()), "fully_qualified_name", warehouseId2.FullyQualifiedName()),
				),
			},
			// Change config but use defaults for every attribute (but not the parameters) - expect no changes (because these are already SF values)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullWithoutParameters)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelRenamedFullWithoutParameters.ResourceReference()), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, replaceResourceReference(warehouseModelRenamedFullWithoutParameters.ResourceReference())).
						HasNameString(warehouseId2.Name()).
						HasWarehouseTypeString(string(sdk.WarehouseTypeStandard)).
						HasWarehouseSizeString(string(sdk.WarehouseSizeXSmall)).
						HasMaxClusterCountString("1").
						HasMinClusterCountString("1").
						HasScalingPolicyString(string(sdk.ScalingPolicyStandard)).
						HasAutoSuspendString("600").
						HasAutoResumeString(r.BooleanTrue).
						HasInitiallySuspendedString(r.BooleanFalse).
						HasNoResourceMonitor().
						HasCommentString(comment).
						HasEnableQueryAccelerationString(r.BooleanFalse).
						HasQueryAccelerationMaxScaleFactorString("8").
						HasNoMaxConcurrencyLevel().
						HasNoStatementQueuedTimeoutInSeconds().
						HasNoStatementTimeoutInSeconds(),
				),
			},
			// add parameters - update expected (different level even with same values)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullWithParameters)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelRenamedFullWithParameters.ResourceReference()), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),

						// this is this only situation in which there will be a strange output in the plan
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFullWithParameters.ResourceReference()), "max_concurrency_level", tfjson.ActionUpdate, nil, sdk.String("8")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFullWithParameters.ResourceReference()), "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFullWithParameters.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("172800")),
					},
				},
				Check: assertThat(t,
					// no changes in the attributes, only for parameters
					resourceassert.WarehouseResource(t, replaceResourceReference(warehouseModel.ResourceReference())).
						HasNameString(warehouseId2.Name()).
						HasWarehouseTypeString(string(sdk.WarehouseTypeStandard)).
						HasWarehouseSizeString(string(sdk.WarehouseSizeXSmall)).
						HasMaxClusterCountString("1").
						HasMinClusterCountString("1").
						HasScalingPolicyString(string(sdk.ScalingPolicyStandard)).
						HasAutoSuspendString("600").
						HasAutoResumeString(r.BooleanTrue).
						HasInitiallySuspendedString(r.BooleanFalse).
						HasNoResourceMonitor().
						HasCommentString(comment).
						HasEnableQueryAccelerationString(r.BooleanFalse).
						HasQueryAccelerationMaxScaleFactorString("8").
						HasMaxConcurrencyLevelString("8").
						HasStatementQueuedTimeoutInSecondsString("0").
						HasStatementTimeoutInSecondsString("172800"),
				),
			},
			// additional step to tackle
			//  001425 (22023): SQL compilation error:
			//  invalid property combination 'RESOURCE_CONSTRAINT'='MEMORY_16X' and
			//  'WAREHOUSE_SIZE'='X-Small'
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullWithParametersMediumSize)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFullWithParametersMediumSize.ResourceReference()), "warehouse_size", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseSizeXSmall)), sdk.String(string(sdk.WarehouseSizeMedium))),
					},
				},
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, replaceResourceReference(warehouseModelRenamedFullWithParametersMediumSize.ResourceReference())).
						HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)),
					objectassert.Warehouse(t, warehouseId2).
						HasSize(sdk.WarehouseSizeMedium),
				),
			},
			// CHANGE PROPERTIES (normal and parameters)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFull)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "enable_query_acceleration", "query_acceleration_max_scale_factor", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),

						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "max_cluster_count", tfjson.ActionUpdate, sdk.String("1"), sdk.String("4")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "min_cluster_count", tfjson.ActionUpdate, sdk.String("1"), sdk.String("2")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "scaling_policy", tfjson.ActionUpdate, sdk.String(string(sdk.ScalingPolicyStandard)), sdk.String(string(sdk.ScalingPolicyEconomy))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, sdk.String("600"), sdk.String("1200")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "enable_query_acceleration", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("8"), sdk.String("4")),

						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "max_concurrency_level", tfjson.ActionUpdate, sdk.String("8"), sdk.String("4")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), sdk.String("5")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), sdk.String("86400")),
					},
				},
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, replaceResourceReference(warehouseModelRenamedFullWithParametersMediumSize.ResourceReference())).
						HasWarehouseTypeString(string(sdk.WarehouseTypeSnowparkOptimized)).
						HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)).
						HasMaxClusterCountString("4").
						HasMinClusterCountString("2").
						HasScalingPolicyString(string(sdk.ScalingPolicyEconomy)).
						HasAutoSuspendString("1200").
						HasAutoResumeString("false").
						// TODO [mux-PR]: change after IgnoreAfterCreation is added
						HasInitiallySuspendedString("false").
						HasResourceMonitorString(resourceMonitor.ID().Name()).
						HasCommentString(newComment).
						HasEnableQueryAccelerationString("true").
						HasQueryAccelerationMaxScaleFactorString("4").
						HasMaxConcurrencyLevelString("4").
						HasStatementQueuedTimeoutInSecondsString("5").
						HasStatementTimeoutInSecondsString("86400"),
					objectassert.Warehouse(t, warehouseId2).
						HasType(sdk.WarehouseTypeSnowparkOptimized).
						HasSize(sdk.WarehouseSizeMedium).
						HasMaxClusterCount(4).
						HasMinClusterCount(2).
						HasScalingPolicy(sdk.ScalingPolicyEconomy).
						HasAutoSuspend(1200).
						HasAutoResume(false).
						HasResourceMonitor(resourceMonitor.ID()).
						HasComment(newComment).
						HasEnableQueryAcceleration(true).
						HasQueryAccelerationMaxScaleFactor(4),
				),
			},
			// TODO [mux-PR]: expect no changes with identifier suppression (using plancheck.ExpectNonEmptyPlan() temporarily)
			// change resource monitor - wrap in quotes (no change expected)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFullResourceMonitorInQuotes)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// CHANGE max_concurrency_level EXTERNALLY (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318)
			{
				Config:    replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelRenamedFull)),
				PreConfig: func() { testClient().Warehouse.UpdateMaxConcurrencyLevel(t, warehouseId2, 10) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "max_concurrency_level", sdk.String("4"), sdk.String("10")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "max_concurrency_level", tfjson.ActionUpdate, sdk.String("10"), sdk.String("4")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "name", warehouseId2.Name()),
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelRenamedFull.ResourceReference()), "max_concurrency_level", "4"),
				),
			},
			// IMPORT
			{
				ResourceName:      replaceResourceReference(warehouseModelRenamedFull.ResourceReference()),
				ImportState:       true,
				ImportStateVerify: true,
				// TODO[mux-PR]: adjust when handling IgnoreAfterCreate
				ImportStateVerifyIgnore: []string{"initially_suspended"},
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelStandard.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelStandard.ResourceReference()), "warehouse_type", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseTypeStandard))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelStandard.ResourceReference()), "warehouse_type", string(sdk.WarehouseTypeStandard))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when type in config
			{
				ResourceName: replaceResourceReference(warehouseModelStandard.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_type", string(sdk.WarehouseTypeStandard)),
				),
			},
			// change type in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSnowparkOptimized)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelSnowparkOptimized.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelSnowparkOptimized.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelSnowparkOptimized.ResourceReference()), "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized))),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeSnowparkOptimized),
				),
			},
			// remove type from config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoType)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(replaceResourceReference(warehouseModelNoType.ResourceReference()), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type")),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// add config (lower case)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSnowparkOptimizedLowercase)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelSnowparkOptimizedLowercase.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelSnowparkOptimizedLowercase.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelSnowparkOptimizedLowercase.ResourceReference()), "warehouse_type", strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized))), sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type")),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", nil, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoType.ResourceReference()), "warehouse_type")),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: replaceResourceReference(warehouseModelNoType.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_type", string(sdk.WarehouseTypeStandard)),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelSmall.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelSmall.ResourceReference()), "warehouse_size", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseSizeSmall))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelSmall.ResourceReference()), "warehouse_size", string(sdk.WarehouseSizeSmall))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeSmall),
				),
			},
			// import when size in config
			{
				ResourceName: replaceResourceReference(warehouseModelSmall.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_size", string(sdk.WarehouseSizeSmall)),
				),
			},
			// change size in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelMedium)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelMedium.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelMedium.ResourceReference()), "warehouse_size", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseSizeSmall)), sdk.String(string(sdk.WarehouseSizeMedium))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelMedium.ResourceReference()), "warehouse_size", string(sdk.WarehouseSizeMedium))),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeMedium),
				),
			},
			// remove size from config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelNoSize)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(replaceResourceReference(warehouseModelNoSize.ResourceReference()), plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeMedium)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size")),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeXSmall),
				),
			},
			// add config (lower case)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelSmallLowercase)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelSmallLowercase.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelSmallLowercase.ResourceReference()), "warehouse_size", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall)))),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelSmallLowercase.ResourceReference()), "warehouse_size", strings.ToLower(string(sdk.WarehouseSizeSmall)))),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall))), sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeXSmall)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size")),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", nil, sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeSmall)), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelNoSize.ResourceReference()), "warehouse_size")),
					objectassert.Warehouse(t, id).HasSize(sdk.WarehouseSizeXSmall),
				),
			},
			// import when no size in config
			{
				ResourceName: replaceResourceReference(warehouseModelNoSize.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "warehouse_size", string(sdk.WarehouseSizeXSmall)),
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
				ExpectError: regexp.MustCompile(`Attribute max_cluster_count value must be at least 1, got: 0`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidMinClusterCount)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute min_cluster_count value must be at least 1, got: 0`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidScalingPolicy)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid scaling policy: unknown"),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidAutoResume)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Inappropriate value for attribute "auto_resume": a bool is required.`),
			},
			{
				Config:      replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelInvalidMaxConcurrencyLevel)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute max_concurrency_level value must be at least 1, got: -2`),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelAutoResumeTrue.ResourceReference()), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelAutoResumeTrue.ResourceReference()), "auto_resume", tfjson.ActionCreate, nil, sdk.String("true")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelAutoResumeTrue.ResourceReference()), "auto_resume", "true")),
					objectassert.Warehouse(t, id).HasAutoResume(true),
				),
			},
			// import when type in config
			{
				ResourceName: replaceResourceReference(warehouseModelAutoResumeTrue.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_resume", "true"),
				),
			},
			// change value in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoResumeFalse)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelAutoResumeFalse.ResourceReference()), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelAutoResumeFalse.ResourceReference()), "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelAutoResumeFalse.ResourceReference()), "auto_resume", "false")),
					objectassert.Warehouse(t, id).HasAutoResume(false),
				),
			},
			// remove type from config (expecting non-empty plan because we do not know the default)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoResume)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume", tfjson.ActionUpdate, sdk.String("false"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume")),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume", nil, sdk.String("false")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume", tfjson.ActionUpdate, sdk.String("false"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()), "auto_resume")),
					objectassert.Warehouse(t, id).HasType(sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: replaceResourceReference(warehouseModelWithoutAutoResume.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_resume", "true"),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelAutoSuspend1200.ResourceReference()), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelAutoSuspend1200.ResourceReference()), "auto_suspend", tfjson.ActionCreate, nil, sdk.String("1200")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelAutoSuspend1200.ResourceReference()), "auto_suspend", "1200")),
					objectassert.Warehouse(t, id).HasAutoSuspend(1200),
				),
			},
			// import when auto suspend in config
			{
				ResourceName: replaceResourceReference(warehouseModelAutoSuspend1200.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "1200"),
				),
			},
			// change value in config to Snowflake default
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelAutoSuspend600)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelAutoSuspend600.ResourceReference()), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelAutoSuspend600.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, sdk.String("1200"), sdk.String("600")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelAutoSuspend600.ResourceReference()), "auto_suspend", "600")),
					objectassert.Warehouse(t, id).HasAutoSuspend(600),
				),
			},
			// remove auto suspend from config (expecting non-empty plan because we do not know the default)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithoutAutoSuspend)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, sdk.String("600"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend")),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend", r.ShowOutputAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend", nil, sdk.String("2400")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, sdk.String("2400"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()), "auto_suspend")),
					objectassert.Warehouse(t, id).HasAutoSuspend(600),
				),
			},
			// import when no type in config
			{
				ResourceName: replaceResourceReference(warehouseModelWithoutAutoSuspend.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "600"),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "query_acceleration_max_scale_factor", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_queued_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "query_acceleration_max_scale_factor", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_queued_timeout_in_seconds", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_timeout_in_seconds", "0")),
					objectassert.Warehouse(t, id).
						HasAutoSuspend(0).
						HasQueryAccelerationMaxScaleFactor(0),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasStatementTimeoutInSeconds(0).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModel)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, sdk.String("0"), nil),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("0"), nil),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), nil),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "auto_suspend")),
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "query_acceleration_max_scale_factor")),
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_queued_timeout_in_seconds")),
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectassert.Warehouse(t, id).
						HasAutoSuspend(600).
						HasQueryAccelerationMaxScaleFactor(8),
					objectparametersassert.WarehouseParameters(t, id).
						HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
						HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault),
				),
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithAllValidZeroValues)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", r.ShowOutputAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", tfjson.ActionUpdate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "query_acceleration_max_scale_factor", tfjson.ActionUpdate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("0")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("0")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "auto_suspend", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "query_acceleration_max_scale_factor", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_queued_timeout_in_seconds", "0")),
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()), "statement_timeout_in_seconds", "0")),
					objectassert.Warehouse(t, id).
						HasAutoSuspend(0).
						HasQueryAccelerationMaxScaleFactor(0),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementQueuedTimeoutInSeconds(0).
						HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasStatementTimeoutInSeconds(0).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
				),
			},
			// import zero values
			{
				ResourceName: replaceResourceReference(warehouseModelWithAllValidZeroValues.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),

					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "auto_suspend", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "query_acceleration_max_scale_factor", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_queued_timeout_in_seconds", "0"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "0"),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("86400")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference()), "statement_timeout_in_seconds", "86400")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(86400).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
				),
			},
			// import when param in config
			{
				ResourceName: replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds86400.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds", "86400"),
				),
			},
			// change the param value in config
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds43200)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", "43200")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(43200).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", "43200")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(43200).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", sdk.String("43200"), sdk.String("86400")),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", "43200")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(43200).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectDrift(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", sdk.String("43200"), nil),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("43200")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds43200.ResourceReference()), "statement_timeout_in_seconds", "43200")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(43200).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("43200"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectparametersassert.WarehouseParameters(t, id).
						HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault),
				),
			},
			// import when param not in config (snowflake default)
			{
				ResourceName: replaceResourceReference(warehouseModel.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds"),
				),
			},
			// change the param value in config to snowflake default (expecting action because of the different level)
			{
				Config: replaceWithWarehousePoCResourceType(t, config.FromModels(t, warehouseModelWithStatementTimeoutInSeconds172800)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, nil, sdk.String("172800")),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithStatementTimeoutInSeconds172800.ResourceReference()), "statement_timeout_in_seconds", "172800")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(172800).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectparametersassert.WarehouseParameters(t, id).
						HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(86400).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeAccount),
				),
			},
			// import when param not in config (set on account)
			{
				ResourceName: replaceResourceReference(warehouseModel.ResourceReference()),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(helpers.EncodeResourceIdentifier(id), "statement_timeout_in_seconds"),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						planchecks.ExpectChange(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), nil),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectparametersassert.WarehouseParameters(t, id).
						HasStatementTimeoutInSeconds(86400).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeAccount),
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
						planchecks.PrintPlanDetails(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds", r.ParametersAttributeName),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					assert.Check(resource.TestCheckNoResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "statement_timeout_in_seconds")),
					objectparametersassert.WarehouseParameters(t, id).
						HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault),
				),
			},
		},
	})
}

// TODO [mux-PR]: address with IgnoreAfterCreate
func TestAcc_TerraformPluginFrameworkPoc_WarehousePoc_InitiallySuspendedChangesPostCreation(t *testing.T) {
	t.Skip("IgnoreAfterCreate not supported yet")
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
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedTrue.ResourceReference()), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedTrue.ResourceReference()), "show_output.#", "1"),
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedTrue.ResourceReference()), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
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
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedFalse.ResourceReference()), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedFalse.ResourceReference()), "show_output.#", "1"),
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModelWithInitiallySuspendedFalse.ResourceReference()), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
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
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "initially_suspended", "true"),

					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "show_output.#", "1"),
					resource.TestCheckResourceAttr(replaceResourceReference(warehouseModel.ResourceReference()), "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
		},
	})
}
