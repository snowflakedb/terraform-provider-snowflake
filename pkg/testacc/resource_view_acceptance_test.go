//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1423486): Fix using warehouse in all tests and remove unsetting testenvs.ConfigureClientOnce
func TestAcc_View_basic(t *testing.T) {
	rowAccessPolicy, rowAccessPolicyCleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	rowAccessPolicy2, rowAccessPolicy2Cleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicy2Cleanup)

	aggregationPolicy2, aggregationPolicy2Cleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicy2Cleanup)

	functionId := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "AVG")
	function2Id := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "MAX")

	cron := "10 * * * * UTC"
	cron2 := "20 * * * * UTC"

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())
	otherStatement := fmt.Sprintf("SELECT foo, id FROM %s", table.ID().FullyQualifiedName())
	comment := random.Comment()

	columnNames := []string{"ID", "FOO"}
	basicViewModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).WithColumnNames(columnNames...)
	viewModelRecursiveWithOtherStatement := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), otherStatement).WithColumnNames(columnNames...).WithIsRecursive(provider.BooleanTrue)
	viewModelWithOtherStatement := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), otherStatement).WithColumnNames(columnNames...)

	updatedViewModel := func(
		rowAccessPolicyId sdk.SchemaObjectIdentifier,
		aggregationPolicyId sdk.SchemaObjectIdentifier,
		dataMetricFunctionId sdk.SchemaObjectIdentifier,
		statement string,
		cron string,
		scheduleStatus sdk.DataMetricScheduleStatusOption,
	) *model.ViewModel {
		return model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).
			WithRowAccessPolicy(rowAccessPolicyId, "ID").
			WithAggregationPolicy(aggregationPolicyId, "ID").
			WithDataMetricFunction(dataMetricFunctionId, "ID", scheduleStatus).
			WithDataMetricSchedule(cron).
			WithComment(comment).
			WithColumnNames(columnNames...)
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: accconfig.ResourceFromModel(t, basicViewModel),
				Check: assertThat(t,
					resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.ResourceFromModel(t, basicViewModel),
				ResourceName: basicViewModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedViewResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasStatementString(statement).
						HasColumnLength(2),
				),
			},
			// set policies and dmfs externally
			{
				PreConfig: func() {
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []sdk.Column{{Value: "ID"}})))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*sdk.NewViewSetAggregationPolicyRequest(aggregationPolicy)))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithSetDataMetricSchedule(*sdk.NewViewSetDataMetricScheduleRequest(fmt.Sprintf("USING CRON %s", cron))))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{
						{
							DataMetricFunction: functionId,
							On:                 []sdk.Column{{Value: "ID"}},
						},
					})))
				},
				Config: accconfig.ResourceFromModel(t, basicViewModel),
				Check: assertThat(t,
					resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasAggregationPolicyLength(0).
						HasRowAccessPolicyLength(0).
						HasDataMetricScheduleLength(0).
						HasDataMetricFunctionLength(0),
				),
			},
			// set other fields
			{
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy.ID(), aggregationPolicy, functionId, statement, cron, sdk.DataMetricScheduleStatusStarted)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicViewModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// change policies and dmfs
			{
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy2.ID(), aggregationPolicy2, function2Id, statement, cron2, sdk.DataMetricScheduleStatusStarted)),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron2)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.schedule_status", string(sdk.DataMetricScheduleStatusStarted))),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", function2Id.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// change dmf status
			{
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy2.ID(), aggregationPolicy2, function2Id, statement, cron2, sdk.DataMetricScheduleStatusSuspended)),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron2)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.schedule_status", string(sdk.DataMetricScheduleStatusSuspended))),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", function2Id.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// change statement and policies
			{
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy.ID(), aggregationPolicy, functionId, otherStatement, cron, sdk.DataMetricScheduleStatusStarted)),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// change statements externally
			{
				PreConfig: func() {
					testClient().View.RecreateView(t, id, statement)
				},
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy.ID(), aggregationPolicy, functionId, otherStatement, cron, sdk.DataMetricScheduleStatusStarted)),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// unset policies externally
			{
				PreConfig: func() {
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithDropAllRowAccessPolicies(true))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest()))
				},
				Config: accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy.ID(), aggregationPolicy, functionId, otherStatement, cron, sdk.DataMetricScheduleStatusStarted)),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.using_cron", cron)),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicViewModel.ResourceReference(), "data_metric_function.0.on.0", "ID")),
				),
			},
			// import - with optionals
			{
				Config:       accconfig.ResourceFromModel(t, updatedViewModel(rowAccessPolicy.ID(), aggregationPolicy, functionId, otherStatement, cron, sdk.DataMetricScheduleStatusStarted)),
				ResourceName: basicViewModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "name", id.Name())),
					resourceassert.ImportedViewResource(t, resourceId).
						HasNameString(id.Name()).
						HasStatementString(otherStatement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasIsSecureString("false").
						HasIsTemporaryString("false").
						HasChangeTrackingString("false").
						HasAggregationPolicyLength(1).
						HasRowAccessPolicyLength(1),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "aggregation_policy.0.entity_key.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "aggregation_policy.0.entity_key.0", "ID")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "row_access_policy.0.on.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(resourceId, "row_access_policy.0.on.0", "ID")),
				),
			},
			// unset
			{
				Config:       accconfig.ResourceFromModel(t, viewModelWithOtherStatement),
				ResourceName: basicViewModel.ResourceReference(),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasNoAggregationPolicyByLength().
					HasNoRowAccessPolicyByLength().
					HasNoDataMetricScheduleByLength().
					HasNoDataMetricFunctionByLength(),
				),
			},
			// recreate - change is_recursive
			{
				Config: accconfig.ResourceFromModel(t, viewModelRecursiveWithOtherStatement),
				Check: assertThat(t, resourceassert.ViewResource(t, basicViewModel.ResourceReference()).
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasIsRecursiveString("true").
					HasIsTemporaryString("default").
					HasChangeTrackingString("default").
					HasNoAggregationPolicyByLength().
					HasNoRowAccessPolicyByLength().
					HasNoDataMetricScheduleByLength().
					HasNoDataMetricFunctionByLength(),
				),
			},
		},
	})
}

func TestAcc_View_recursive(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	basicView := config.Variables{
		"name":         config.StringVariable(id.Name()),
		"database":     config.StringVariable(id.DatabaseName()),
		"schema":       config.StringVariable(id.SchemaName()),
		"statement":    config.StringVariable(statement),
		"is_recursive": config.BoolVariable(true),
		"column": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"column_name": config.StringVariable("ROLE_NAME"),
			}),
			config.MapVariable(map[string]config.Variable{
				"column_name": config.StringVariable("ROLE_OWNER"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_is_recursive"),
				ConfigVariables: basicView,
				Check: assertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasIsRecursiveString("true")),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_is_recursive"),
				ConfigVariables: basicView,
				ResourceName:    "snowflake_view.test",
				ImportState:     true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedViewResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasStatementString(statement).
						HasIsRecursiveString("true")),
			},
		},
	})
}

func TestAcc_View_temporary(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement)

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, viewModel.WithIsTemporary("true")),
				Check: assertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasIsTemporaryString("true")),
			},
		},
	})
}

func TestAcc_View_complete(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	rowAccessPolicy, rowAccessPolicyCleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	projectionPolicy, projectionPolicyCleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicyCleanup)

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyWithOptions(t,
		[]sdk.TableColumnSignature{
			{
				Name: "One",
				Type: testdatatypes.DataTypeNumber,
			},
			{
				Name: "Two",
				Type: testdatatypes.DataTypeNumber,
			},
		},
		testdatatypes.DataTypeNumber,
		`
case
	when One > 0 then One
	else Two
end;;
`,
		new(sdk.CreateMaskingPolicyOptions),
	)
	t.Cleanup(maskingPolicyCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())
	functionId := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "AVG")

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                            config.StringVariable(id.Name()),
			"database":                        config.StringVariable(id.DatabaseName()),
			"schema":                          config.StringVariable(id.SchemaName()),
			"comment":                         config.StringVariable("Terraform test resource"),
			"is_secure":                       config.BoolVariable(true),
			"is_temporary":                    config.BoolVariable(false),
			"copy_grants":                     config.BoolVariable(false),
			"change_tracking":                 config.BoolVariable(true),
			"row_access_policy":               config.StringVariable(rowAccessPolicy.ID().FullyQualifiedName()),
			"row_access_policy_on":            config.ListVariable(config.StringVariable("ID")),
			"aggregation_policy":              config.StringVariable(aggregationPolicy.FullyQualifiedName()),
			"aggregation_policy_entity_key":   config.ListVariable(config.StringVariable("ID")),
			"statement":                       config.StringVariable(statement),
			"warehouse":                       config.StringVariable(TestWarehouseName),
			"column1_name":                    config.StringVariable("ID"),
			"column1_comment":                 config.StringVariable("col comment"),
			"column2_name":                    config.StringVariable("FOO"),
			"column2_masking_policy":          config.StringVariable(maskingPolicy.ID().FullyQualifiedName()),
			"column2_masking_policy_using":    config.ListVariable(config.StringVariable("FOO"), config.StringVariable("ID")),
			"column2_projection_policy":       config.StringVariable(projectionPolicy.FullyQualifiedName()),
			"data_metric_function":            config.StringVariable(functionId.FullyQualifiedName()),
			"data_metric_function_on":         config.ListVariable(config.StringVariable("ID")),
			"data_metric_schedule_using_cron": config.StringVariable("5 * * * * UTC"),
		}
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/complete"),
				ConfigVariables: m(),
				Check: assertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("Terraform test resource").
					HasIsSecureString("true").
					HasIsTemporaryString("false").
					HasChangeTrackingString("true").
					HasDataMetricScheduleLength(1).
					HasDataMetricFunctionLength(1).
					HasAggregationPolicyLength(1).
					HasRowAccessPolicyLength(1).
					HasColumnLength(2),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "data_metric_schedule.0.using_cron", "5 * * * * UTC")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "data_metric_schedule.0.minutes", "0")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "data_metric_function.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "data_metric_function.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.0.column_name", "ID")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.0.masking_policy.#", "0")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.0.projection_policy.#", "0")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.0.comment", "col comment")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.1.column_name", "FOO")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.1.masking_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.1.masking_policy.0.policy_name", maskingPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.1.projection_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "column.1.projection_policy.0.policy_name", projectionPolicy.FullyQualifiedName())),
					resourceshowoutputassert.ViewShowOutput(t, "snowflake_view.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("Terraform test resource").
						HasIsSecure(true).
						HasChangeTracking("ON"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/complete"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_view.test",
				ImportState:     true,
				ImportStateCheck: assertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("Terraform test resource").
						HasIsSecureString("true").
						HasIsTemporaryString("false").
						HasChangeTrackingString("true").
						HasDataMetricScheduleLength(1).
						HasDataMetricFunctionLength(1).
						HasAggregationPolicyLength(1).
						HasRowAccessPolicyLength(1),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "data_metric_schedule.0.using_cron", "5 * * * * UTC")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "data_metric_schedule.0.minutes", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "data_metric_function.0.function_name", functionId.FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "data_metric_function.0.on.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "data_metric_function.0.on.0", "ID")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "aggregation_policy.0.entity_key.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "aggregation_policy.0.entity_key.0", "ID")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "row_access_policy.0.on.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "row_access_policy.0.on.0", "ID")),
				),
			},
		},
	})
}

func TestAcc_View_columns(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("bar", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyWithOptions(t,
		[]sdk.TableColumnSignature{
			{
				Name: "One",
				Type: testdatatypes.DataTypeNumber,
			},
		},
		testdatatypes.DataTypeNumber,
		`
case
	when One > 0 then One
	else 0
end;;
`,
		new(sdk.CreateMaskingPolicyOptions),
	)
	t.Cleanup(maskingPolicyCleanup)

	projectionPolicy, projectionPolicyCleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicyCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())

	basicView := func(columns ...string) config.Variables {
		return config.Variables{
			"name":      config.StringVariable(id.Name()),
			"database":  config.StringVariable(id.DatabaseName()),
			"schema":    config.StringVariable(id.SchemaName()),
			"statement": config.StringVariable(statement),
			"column": config.SetVariable(
				collections.Map(columns, func(columnName string) config.Variable {
					return config.MapVariable(map[string]config.Variable{
						"column_name": config.StringVariable(columnName),
					})
				})...,
			),
		}
	}

	basicViewWithPolicies := func() config.Variables {
		conf := basicView("ID", "FOO")
		delete(conf, "column")
		conf["projection_name"] = config.StringVariable(projectionPolicy.FullyQualifiedName())
		conf["masking_name"] = config.StringVariable(maskingPolicy.ID().FullyQualifiedName())
		conf["masking_using"] = config.ListVariable(config.StringVariable("ID"))
		return conf
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			// Columns without policies
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables: basicView("ID", "FOO"),
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
				),
			},
			// Columns with policies added externally
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables: basicView("ID", "FOO"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				PreConfig: func() {
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithSetMaskingPolicyOnColumn(*sdk.NewViewSetColumnMaskingPolicyRequest("ID", maskingPolicy.ID()).WithUsing([]sdk.Column{{Value: "ID"}})))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithSetProjectionPolicyOnColumn(*sdk.NewViewSetProjectionPolicyRequest("ID", projectionPolicy).WithForce(true)))
				},
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
					objectassert.View(t, id).
						HasNoMaskingPolicyReferences(testClient()).
						HasNoProjectionPolicyReferences(testClient()),
				),
			},
			// With all policies on columns
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/columns"),
				ConfigVariables: basicViewWithPolicies(),
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
					objectassert.View(t, id).
						HasMaskingPolicyReferences(testClient(), 1).
						HasProjectionPolicyReferences(testClient(), 1),
				),
			},
			// Remove policies on columns externally
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/columns"),
				ConfigVariables: basicViewWithPolicies(),
				PreConfig: func() {
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetMaskingPolicyOnColumn(*sdk.NewViewUnsetColumnMaskingPolicyRequest("ID")))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetProjectionPolicyOnColumn(*sdk.NewViewUnsetProjectionPolicyRequest("ID")))
				},
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
					objectassert.View(t, id).
						HasMaskingPolicyReferences(testClient(), 1).
						HasProjectionPolicyReferences(testClient(), 1),
				),
			},
		},
	})
}

func TestAcc_View_columnsWithMaskingPolicyWithoutUsing(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("bar", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyWithOptions(t,
		[]sdk.TableColumnSignature{
			{
				Name: "One",
				Type: testdatatypes.DataTypeNumber,
			},
		},
		testdatatypes.DataTypeNumber,
		`
case
	when One > 0 then One
	else 0
end;;
`,
		new(sdk.CreateMaskingPolicyOptions),
	)
	t.Cleanup(maskingPolicyCleanup)

	projectionPolicy, projectionPolicyCleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicyCleanup)

	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())

	viewWithPolicies := func() config.Variables {
		conf := config.Variables{
			"name":      config.StringVariable(id.Name()),
			"database":  config.StringVariable(id.DatabaseName()),
			"schema":    config.StringVariable(id.SchemaName()),
			"statement": config.StringVariable(statement),
		}
		conf["projection_name"] = config.StringVariable(projectionPolicy.FullyQualifiedName())
		conf["masking_name"] = config.StringVariable(maskingPolicy.ID().FullyQualifiedName())
		return conf
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			// With all policies on columns
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/columns"),
				ConfigVariables: viewWithPolicies(),
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
					objectassert.View(t, id).
						HasMaskingPolicyReferences(testClient(), 1).
						HasProjectionPolicyReferences(testClient(), 1),
				),
			},
			// Remove policies on columns externally
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/columns"),
				ConfigVariables: viewWithPolicies(),
				PreConfig: func() {
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetMaskingPolicyOnColumn(*sdk.NewViewUnsetColumnMaskingPolicyRequest("ID")))
					testClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetProjectionPolicyOnColumn(*sdk.NewViewUnsetProjectionPolicyRequest("ID")))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ViewResource(t, "snowflake_view.test").
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasColumnLength(2),
					objectassert.View(t, id).
						HasMaskingPolicyReferences(testClient(), 1).
						HasProjectionPolicyReferences(testClient(), 1),
				),
			},
		},
	})
}

func TestAcc_View_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewConfig := func(identifier sdk.SchemaObjectIdentifier) config.Variables {
		return config.Variables{
			"name":      config.StringVariable(identifier.Name()),
			"database":  config.StringVariable(identifier.DatabaseName()),
			"schema":    config.StringVariable(identifier.SchemaName()),
			"statement": config.StringVariable(statement),
			"column": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("ROLE_NAME"),
				}),
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("ROLE_OWNER"),
				}),
			),
		}
	}

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables: viewConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			// rename with one param changed
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables: viewConfig(newId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_View_Issue3073(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement)
	viewModelWithColumns := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).WithColumnValue(config.SetVariable(
		config.MapVariable(map[string]config.Variable{
			"column_name": config.StringVariable("ROLE_NAME"),
		}),
		config.MapVariable(map[string]config.Variable{
			"column_name": config.StringVariable("ROLE_OWNER"),
		}),
	))

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, viewModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
				),
			},
			// specify the columns
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, viewModelWithColumns),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, viewModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
				),
			},
		},
	})
}

// fixes https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3073#issuecomment-2392250469
func TestAcc_View_IncorrectColumnsWithOrReplace(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := `SELECT ROLE_NAME as "role_name", ROLE_OWNER as "role_owner" FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	statementUnquotedColumns := `SELECT ROLE_NAME as role_name, ROLE_OWNER as role_owner FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	statementUnquotedColumns3 := `SELECT ROLE_NAME as role_name, ROLE_OWNER as role_owner, IS_GRANTABLE as is_grantable FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`

	viewModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement)
	viewLowercaseStatementModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statementUnquotedColumns)
	viewLowercaseStatementModel3 := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statementUnquotedColumns3)

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, viewModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.0.column_name", "role_name"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.1.column_name", "role_owner"),
				),
			},
			// use columns without quotes in the statement
			{
				Config: accconfig.FromModels(t, viewLowercaseStatementModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.0.column_name", "ROLE_NAME"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.1.column_name", "ROLE_OWNER"),
				),
			},
			// add a new column to the statement
			{
				Config: accconfig.FromModels(t, viewLowercaseStatementModel3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.0.column_name", "ROLE_NAME"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.1.column_name", "ROLE_OWNER"),
					resource.TestCheckResourceAttr("snowflake_view.test", "column.2.column_name", "IS_GRANTABLE"),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewConfig := func(copyGrants bool) config.Variables {
		return config.Variables{
			"name":        config.StringVariable(id.Name()),
			"database":    config.StringVariable(id.DatabaseName()),
			"schema":      config.StringVariable(id.SchemaName()),
			"statement":   config.StringVariable(statement),
			"copy_grants": config.BoolVariable(copyGrants),
			"is_secure":   config.BoolVariable(true),
			"column": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("ID"),
				}),
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("FOO"),
				}),
			),
		}
	}

	var createdOn string

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_copy_grants"),
				ConfigVariables: viewConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					resource.TestCheckResourceAttr("snowflake_view.test", "is_secure", "true"),
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					}),
				),
			},
			// Checks that copy_grants changes don't trigger a drop
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_copy_grants"),
				ConfigVariables: viewConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("snowflake_view.test", "is_secure", "true"),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrantsReversed(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewConfig := func(copyGrants bool) config.Variables {
		return config.Variables{
			"name":        config.StringVariable(id.Name()),
			"database":    config.StringVariable(id.DatabaseName()),
			"schema":      config.StringVariable(id.SchemaName()),
			"statement":   config.StringVariable(statement),
			"copy_grants": config.BoolVariable(copyGrants),
			"is_secure":   config.BoolVariable(true),
			"column": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("ID"),
				}),
				config.MapVariable(map[string]config.Variable{
					"column_name": config.StringVariable("FOO"),
				}),
			),
		}
	}

	var createdOn string

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_copy_grants"),
				ConfigVariables: viewConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					resource.TestCheckResourceAttr("snowflake_view.test", "is_secure", "true"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_View/basic_copy_grants"),
				ConfigVariables: viewConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("snowflake_view.test", "is_secure", "true"),
				),
			},
		},
	})
}

func TestAcc_View_CheckGrantsAfterRecreation(t *testing.T) {
	table, cleanupTable := testClient().Table.Create(t)
	t.Cleanup(cleanupTable)

	role, cleanupRole := testClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithGrants(id, table.ID(), "id", role.ID(), true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: viewConfigWithGrants(id, table.ID(), "*", role.ID(), true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			// Recreate without copy grants. Now we expect changes because the grants are still in the config.
			{
				Config:             viewConfigWithGrants(id, table.ID(), "id", role.ID(), false),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.grant", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "1"),
				),
			},
		},
	})
}

func viewConfigWithGrants(viewId, tableId sdk.SchemaObjectIdentifier, selectStatement string, roleId sdk.AccountObjectIdentifier, copyGrants bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[4]s"
  comment = "created by terraform"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "select %[5]s from \"%[1]s\".\"%[2]s\".\"%[3]s\""
  copy_grants = %[7]t
  is_secure = true

  column {
    column_name = "%[5]s"
  }
}

resource "snowflake_grant_privileges_to_account_role" "grant" {
  privileges        = ["SELECT"]
  account_role_name = %[6]s
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.test.fully_qualified_name
  }
}

data "snowflake_grants" "grants" {
  depends_on = [snowflake_grant_privileges_to_account_role.grant, snowflake_view.test]
  grants_on {
    object_type = "VIEW"
    object_name = snowflake_view.test.fully_qualified_name
  }
}
	`, viewId.DatabaseName(), viewId.SchemaName(), tableId.Name(), viewId.Name(), selectStatement, roleId.FullyQualifiedName(), copyGrants)
}

func TestAcc_View_Issue2640(t *testing.T) {
	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	part1 := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	part2 := "SELECT ROLE_OWNER, ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	statement := fmt.Sprintf("%s\n\tunion\n%s\n", part1, part2)

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithMultilineUnionStatement(id, part1, part2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", statement),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", TestSchemaName),
				),
			},
			// try to import secure view without being its owner (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2640)
			{
				PreConfig: func() {
					testClient().Grant.GrantOwnershipOnSchemaObjectToAccountRole(t, role.ID(), sdk.ObjectTypeView, id, sdk.Revoke)
				},
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ExpectError:  regexp.MustCompile("`text` is missing; if the view is secure then the role used by the provider must own the view"),
			},
			// import with the proper role
			{
				PreConfig: func() {
					testClient().Grant.GrantOwnershipOnSchemaObjectToAccountRole(t, snowflakeroles.Accountadmin, sdk.ObjectTypeView, id, sdk.Revoke)
				},
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ImportStateCheck: assertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()),
				),
			},
		},
	})
}

func viewConfigWithMultilineUnionStatement(id sdk.SchemaObjectIdentifier, part1 string, part2 string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = <<-SQL
%[4]s
	union
%[5]s
SQL
  is_secure = true
  column {
    column_name = "ROLE_OWNER"
  }
  column {
    column_name = "ROLE_NAME"
  }
}
	`, id.DatabaseName(), id.SchemaName(), id.Name(), part1, part2)
}

func TestAcc_view_migrateFromVersion_0_94_1(t *testing.T) {
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewConfig := config.Variables{
		"name":      config.StringVariable(id.Name()),
		"database":  config.StringVariable(id.DatabaseName()),
		"schema":    config.StringVariable(id.SchemaName()),
		"statement": config.StringVariable(statement),
		"column": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"column_name": config.StringVariable("ROLE_NAME"),
			}),
			config.MapVariable(map[string]config.Variable{
				"column_name": config.StringVariable("ROLE_OWNER"),
			}),
		),
	}

	resourceName := "snowflake_view.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            viewV0941WithTags(id, tag.ID(), "foo", statement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.0.name", tag.Name),
					resource.TestCheckResourceAttr(resourceName, "tag.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "or_replace", "true"),
				),
			},
			{
				PreConfig: func() { UnsetConfigPathEnv(t) },
				// TODO [SNOW-2324320]: warehouse needed?
				ProtoV6ProviderFactories: providerFactoryWithoutCache(),
				ConfigDirectory:          ConfigurationDirectory("TestAcc_View/basic"),
				ConfigVariables:          viewConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "tag.#"),
					resource.TestCheckNoResourceAttr(resourceName, "or_replace"),
				),
			},
		},
	})
}

func viewV0941WithTags(id sdk.SchemaObjectIdentifier, tagId sdk.SchemaObjectIdentifier, tagValue, statement string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
	database				= "%[1]s"
	schema				    = "%[2]s"
	name					= "%[3]s"
	statement				= "%[6]s"
	or_replace				= true
	tag {
		database = "%[1]s"
		schema = "%[2]s"
		name = "%[4]s"
		value = "%[5]s"
	}
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), tagId.Name(), tagValue, statement)
}

func TestAcc_View_Issue3676_proof(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())
	otherStatement := fmt.Sprintf("SELECT foo, id FROM %s", table.ID().FullyQualifiedName())
	columnNames := []string{"ID", "FOO"}

	initialModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).WithColumnNames(columnNames...)
	renamedAndAlteredModel := model.View("test", newId.DatabaseName(), newId.SchemaName(), newId.Name(), otherStatement).WithColumnNames(columnNames...)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				Config:            accconfig.FromModels(t, initialModel),
				Check: assertThat(t, resourceassert.ViewResource(t, initialModel.ResourceReference()).
					HasDatabaseString(TestDatabaseName).
					HasSchemaString(TestSchemaName).
					HasNameString(id.Name()).
					HasStatementString(statement),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.1.0"),
				Config:            accconfig.FromModels(t, renamedAndAlteredModel),
				Check: assertThat(t, resourceassert.ViewResource(t, initialModel.ResourceReference()).
					HasDatabaseString(TestDatabaseName).
					HasSchemaString(TestSchemaName).
					HasNameString(newId.Name()).
					HasStatementString(otherStatement),
					objectassert.View(t, newId).HasName(newId.Name()),
					// existence of the old view confirms the issue https://github.com/snowflakedb/terraform-provider-snowflake/issues/3676
					objectassert.View(t, id).HasName(id.Name()),
				),
			},
		},
	})
}

func TestAcc_View_Issue3676_fix(t *testing.T) {
	table, tableCleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("foo", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())
	otherStatement := fmt.Sprintf("SELECT foo, id FROM %s", table.ID().FullyQualifiedName())
	columnNames := []string{"ID", "FOO"}

	initialModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).WithColumnNames(columnNames...)
	renamedAndAlteredModel := model.View("test", newId.DatabaseName(), newId.SchemaName(), newId.Name(), otherStatement).WithColumnNames(columnNames...)

	resource.Test(t, resource.TestCase{
		// TODO [SNOW-2324320]: warehouse needed?
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, initialModel),
				Check: assertThat(t, resourceassert.ViewResource(t, initialModel.ResourceReference()).
					HasDatabaseString(TestDatabaseName).
					HasSchemaString(TestSchemaName).
					HasNameString(id.Name()).
					HasStatementString(statement),
				),
			},
			{
				Config: accconfig.FromModels(t, renamedAndAlteredModel),
				Check: assertThat(t, resourceassert.ViewResource(t, initialModel.ResourceReference()).
					HasDatabaseString(TestDatabaseName).
					HasSchemaString(TestSchemaName).
					HasNameString(newId.Name()).
					HasStatementString(otherStatement),
					objectassert.View(t, newId).HasName(newId.Name()),
					// view with the old id does not exist
					invokeactionassert.ViewDoesNotExist(t, id),
				),
			},
		},
	})
}
