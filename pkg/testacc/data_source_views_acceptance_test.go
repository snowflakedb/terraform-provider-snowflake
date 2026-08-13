//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Views_CompleteUseCase(t *testing.T) {
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

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyWithRequest(
		t,
		[]sdk.CreateMaskingPolicySignatureRequest{
			*sdk.NewCreateMaskingPolicySignatureRequest("One", testdatatypes.DataTypeNumber),
			*sdk.NewCreateMaskingPolicySignatureRequest("Two", testdatatypes.DataTypeNumber),
		},
		testdatatypes.DataTypeNumber,
		`
case
	when One > 0 then One
	else Two
end;;
`,
	)
	t.Cleanup(maskingPolicyCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("SELECT id, foo FROM %s", table.ID().FullyQualifiedName())
	comment := random.Comment()
	functionId := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "AVG")

	viewModel := model.View("test", id.DatabaseName(), id.SchemaName(), id.Name(), statement).
		WithComment(comment).
		WithIsSecure("true").
		WithChangeTracking("true").
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID").
		WithDataMetricFunction(functionId, "ID", sdk.DataMetricScheduleStatusStarted).
		WithDataMetricSchedule("5 * * * * UTC").
		WithColumnValue(tfconfig.TupleVariable(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"column_name": tfconfig.StringVariable("ID"),
				"comment":     tfconfig.StringVariable("col comment"),
			}),
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"column_name": tfconfig.StringVariable("FOO"),
				"masking_policy": tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"policy_name": tfconfig.StringVariable(maskingPolicy.ID().FullyQualifiedName()),
					"using":       tfconfig.ListVariable(tfconfig.StringVariable("FOO"), tfconfig.StringVariable("ID")),
				}),
				"projection_policy": tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"policy_name": tfconfig.StringVariable(projectionPolicy.FullyQualifiedName()),
				}),
			}),
		))

	dataSourceModel := datasourcemodel.Views("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(viewModel.ResourceReference())
	dataSourceModelWithoutDescribe := datasourcemodel.Views("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(viewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// defaults (with_describe = true) - assert every field
			{
				Config: accconfig.FromModels(t, viewModel, dataSourceModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.#", "1"),

					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "views.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.kind", ""),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.reserved", ""),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "views.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "views.0.show_output.0.text"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.is_secure", "true"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.is_materialized", "false"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.owner_role_type", "ROLE"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.show_output.0.change_tracking", "ON"),

					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.#", "2"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.0.name", "ID"),
					resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "views.0.describe_output.0.type"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.0.kind", "COLUMN"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.0.comment", "col comment"),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.0.policy_name", ""),
					resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "views.0.describe_output.1.name", "FOO"),
					resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "views.0.describe_output.1.policy_name"),
				),
			},
			// with_describe = false - assert describe_output is empty
			{
				Config: accconfig.FromModels(t, viewModel, dataSourceModelWithoutDescribe),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelWithoutDescribe.DatasourceReference(), "views.#", "1"),
					resource.TestCheckResourceAttr(dataSourceModelWithoutDescribe.DatasourceReference(), "views.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(dataSourceModelWithoutDescribe.DatasourceReference(), "views.0.show_output.0.name", id.Name()),
				),
			},
		},
	})
}

func TestAcc_Views_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(6)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewModel1 := model.View("v1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), statement)
	viewModel2 := model.View("v2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), statement)
	viewModel3 := model.View("v3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), statement)

	dsLikeExact := datasourcemodel.Views("test").
		WithLike(id1.Name()).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())
	dsLikePrefix := datasourcemodel.Views("test").
		WithLike(prefix+"%").
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())
	dsStartsWith := datasourcemodel.Views("test").
		WithStartsWith(prefix).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())
	dsInDatabase := datasourcemodel.Views("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())
	dsInSchema := datasourcemodel.Views("test").
		WithInSchema(schema.ID()).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())
	dsLimit := datasourcemodel.Views("test").
		WithLike(prefix+"%").
		WithLimitRowsAndFrom(1, "").
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference(), viewModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: viewsProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// like (exact match)
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "views.#", "1"),
				),
			},
			// like (prefix pattern)
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLikePrefix.DatasourceReference(), "views.#", "2"),
				),
			},
			// starts_with
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "views.#", "2"),
				),
			},
			// in database
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsInDatabase.DatasourceReference(), "views.#", "2"),
				),
			},
			// in schema
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "views.#", "1"),
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "views.0.show_output.0.name", id3.Name()),
				),
			},
			// limit
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewModel3, dsLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "views.#", "1"),
				),
			},
		},
	})
}
