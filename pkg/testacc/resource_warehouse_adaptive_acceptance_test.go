//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_WarehouseAdaptive_BasicUseCase(t *testing.T) {
	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	warehouseModel := model.WarehouseAdaptiveWithId(warehouseId)
	warehouseModelWithOptionals := model.WarehouseAdaptiveWithId(warehouseId).
		WithComment(comment).
		WithQueryThroughputMultiplier(2).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelLarge)).
		WithStatementQueuedTimeoutInSeconds(300).
		WithStatementTimeoutInSeconds(86400)
	warehouseModelUpdated := model.WarehouseAdaptiveWithId(warehouseId).
		WithComment(newComment).
		WithQueryThroughputMultiplier(4).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelSmall)).
		WithStatementQueuedTimeoutInSeconds(600).
		WithStatementTimeoutInSeconds(43200)

	ref := warehouseModel.ResourceReference()
	externalQueryThroughputMultiplier := 10
	externalStatementTimeout := 99999
	externalStatementQueuedTimeout := 1200
	externalMaxQueryPerformanceLevel := sdk.MaxQueryPerformanceLevelMedium
	externalComment := random.Comment()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentString("").
			HasNoMaxQueryPerformanceLevel().
			HasNoQueryThroughputMultiplier().
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800).
			HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasStateNotEmpty().
			HasCommentEmpty().
			HasOwnerNotEmpty().
			HasOwnerRoleTypeNotEmpty(),
	}

	withOptionalsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentString(comment).
			HasMaxQueryPerformanceLevelString(string(sdk.MaxQueryPerformanceLevelLarge)).
			HasQueryThroughputMultiplier(2).
			HasStatementQueuedTimeoutInSeconds(300).
			HasStatementTimeoutInSeconds(86400),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasComment(comment).
			HasQueryThroughputMultiplier(2),
	}

	updatedAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentString(newComment).
			HasQueryThroughputMultiplier(4).
			HasStatementQueuedTimeoutInSeconds(600).
			HasStatementTimeoutInSeconds(43200),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasName(warehouseId.Name()).
			HasType(sdk.WarehouseTypeAdaptive).
			HasComment(newComment).
			HasQueryThroughputMultiplier(4),
	}

	unsetAssertions := []assert.TestCheckFuncProvider{
		resourceassert.WarehouseAdaptiveResource(t, ref).
			HasNameString(warehouseId.Name()).
			HasCommentString("").
			HasNoQueryThroughputMultiplier().
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800),
		resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, ref).
			HasCommentEmpty().
			HasQueryThroughputMultiplier(0),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				Config: accconfig.FromModels(t, warehouseModel),
				Check:  assertThat(t, basicAssertions...),
			},
			// import after minimal config
			{
				ResourceName:            warehouseModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_query_performance_level", "query_throughput_multiplier"},
			},
			// set all optional fields
			{
				Config: accconfig.FromModels(t, warehouseModelWithOptionals),
				Check:  assertThat(t, withOptionalsAssertions...),
			},
			// import after setting optional fields
			{
				ResourceName:      warehouseModelWithOptionals.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update non-recreating fields
			{
				Config: accconfig.FromModels(t, warehouseModelUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChangeDeleteCreate(ref, "comment", sdk.String(comment), sdk.String(newComment)),
						planchecks.ExpectChangeDeleteCreate(ref, "query_throughput_multiplier", sdk.String("2"), sdk.String("4")),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// detect external changes to multiple fields
			{
				Config: accconfig.FromModels(t, warehouseModelUpdated),
				PreConfig: func() {
					testClient().Warehouse.DropWarehouseFunc(t, warehouseId)()
					testClient().Warehouse.CreateAdaptiveWithOptions(t, warehouseId, &sdk.CreateAdaptiveWarehouseOptions{
						Comment:                         sdk.String(externalComment),
						QueryThroughputMultiplier:       sdk.Int(externalQueryThroughputMultiplier),
						StatementTimeoutInSeconds:       sdk.Int(externalStatementTimeout),
						StatementQueuedTimeoutInSeconds: sdk.Int(externalStatementQueuedTimeout),
						MaxQueryPerformanceLevel:        sdk.Pointer(externalMaxQueryPerformanceLevel),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChangeDeleteCreate(ref, "query_throughput_multiplier", sdk.String("10"), sdk.String("4")),
						planchecks.ExpectChangeDeleteCreate(ref, "statement_timeout_in_seconds", sdk.String("99999"), sdk.String("43200")),
						planchecks.ExpectChangeDeleteCreate(ref, "statement_queued_timeout_in_seconds", sdk.String("1200"), sdk.String("600")),
						planchecks.ExpectChangeDeleteCreate(ref, "max_query_performance_level", sdk.String(string(sdk.MaxQueryPerformanceLevelMedium)), sdk.String(string(sdk.MaxQueryPerformanceLevelSmall))),
						planchecks.ExpectChangeDeleteCreate(ref, "comment", sdk.String(externalComment), sdk.String(newComment)),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// unset optional fields back to defaults
			{
				Config: accconfig.FromModels(t, warehouseModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChangeDeleteCreate(ref, "comment", sdk.String(newComment), nil),
						planchecks.ExpectChangeDeleteCreate(ref, "query_throughput_multiplier", sdk.String("4"), nil),
						planchecks.ExpectChangeDeleteCreate(ref, "statement_queued_timeout_in_seconds", sdk.String("600"), nil),
						planchecks.ExpectChangeDeleteCreate(ref, "statement_timeout_in_seconds", sdk.String("43200"), nil),
						planchecks.ExpectChangeDeleteCreate(ref, "max_query_performance_level", sdk.String(string(sdk.MaxQueryPerformanceLevelSmall)), nil),
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

	warehouseModelComplete := model.WarehouseAdaptiveWithId(warehouseId).
		WithComment(comment).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelLarge)).
		WithQueryThroughputMultiplier(3).
		WithStatementQueuedTimeoutInSeconds(300).
		WithStatementTimeoutInSeconds(86400)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: accconfig.FromModels(t, warehouseModelComplete),
				Check: assertThat(t,
					resourceassert.WarehouseAdaptiveResource(t, warehouseModelComplete.ResourceReference()).
						HasNameString(warehouseId.Name()).
						HasCommentString(comment).
						HasMaxQueryPerformanceLevelString(string(sdk.MaxQueryPerformanceLevelLarge)).
						HasQueryThroughputMultiplier(3).
						HasStatementQueuedTimeoutInSeconds(300).
						HasStatementTimeoutInSeconds(86400).
						HasFullyQualifiedNameString(warehouseId.FullyQualifiedName()),
					resourceshowoutputassert.WarehouseAdaptiveShowOutput(t, warehouseModelComplete.ResourceReference()).
						HasName(warehouseId.Name()).
						HasType(sdk.WarehouseTypeAdaptive).
						HasComment(comment).
						HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
						HasQueryThroughputMultiplier(3),
					resourceparametersassert.WarehouseAdaptiveResourceParameters(t, warehouseModelComplete.ResourceReference()).
						HasStatementQueuedTimeoutInSeconds(300).
						HasStatementTimeoutInSeconds(86400),
				),
			},
			// import and verify state matches
			{
				ResourceName:      warehouseModelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
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
