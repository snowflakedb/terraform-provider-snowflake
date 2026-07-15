//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IcebergTables_BasicUseCase(t *testing.T) {
	table, tableCleanup := testClient().IcebergTable.Create(t)
	t.Cleanup(tableCleanup)

	id := table.ID()

	datasourceModel := datasourcemodel.IcebergTables("test").
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	datasourceModelWithoutDescribe := datasourcemodel.IcebergTables("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	datasourceModelWithoutParameters := datasourcemodel.IcebergTables("test").
		WithWithParameters(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())

	showOutputAssertions := resourceshowoutputassert.IcebergTablesDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasIcebergTableType(sdk.IcebergTableTypeManaged).
		HasCanWriteMetadata(true).
		HasAutoRefreshStatusEmpty()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, datasourceModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "iceberg_tables.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.IcebergTablesDatasourceDescribeOutput(t, datasourceModel.DatasourceReference()).
						HasName("ID"),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "iceberg_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "iceberg_tables.0.describe_output.#", "0")),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModelWithoutParameters),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutParameters.DatasourceReference(), "iceberg_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutParameters.DatasourceReference(), "iceberg_tables.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_IcebergTables_Filtering(t *testing.T) {
	prefix := testClient().Ids.Alpha()
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	_, cleanup1 := testClient().IcebergTable.CreateWithRequest(
		t, sdk.NewCreateIcebergTableRequest(id1,
			sdk.IcebergTableColumnsAndConstraintsRequest{Columns: []sdk.IcebergTableColumnRequest{{Name: "ID", ColumnType: testdatatypes.DataTypeNumber}}}),
	)
	t.Cleanup(cleanup1)
	_, cleanup2 := testClient().IcebergTable.CreateWithRequest(
		t, sdk.NewCreateIcebergTableRequest(id2,
			sdk.IcebergTableColumnsAndConstraintsRequest{Columns: []sdk.IcebergTableColumnRequest{{Name: "ID", ColumnType: testdatatypes.DataTypeNumber}}}),
	)
	t.Cleanup(cleanup2)

	modelByPrefix := datasourcemodel.IcebergTables("by_prefix").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId())
	modelById1 := datasourcemodel.IcebergTables("by_id1").
		WithLike(id1.Name()).
		WithInSchema(id1.SchemaId())
	modelBySchema := datasourcemodel.IcebergTables("by_schema").
		WithInSchema(id1.SchemaId()).
		WithStartsWith(prefix)
	modelInDatabase := datasourcemodel.IcebergTables("in_database").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId())
	modelWithLimit := datasourcemodel.IcebergTables("with_limit").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId()).
		WithLimit(1)
	modelNoResults := datasourcemodel.IcebergTables("no_results").
		WithLike("non_existent_%").
		WithInSchema(id1.SchemaId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelByPrefix, modelById1, modelBySchema, modelInDatabase, modelWithLimit, modelNoResults),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(modelByPrefix.DatasourceReference(), "iceberg_tables.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelById1.DatasourceReference(), "iceberg_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBySchema.DatasourceReference(), "iceberg_tables.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelInDatabase.DatasourceReference(), "iceberg_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelWithLimit.DatasourceReference(), "iceberg_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelNoResults.DatasourceReference(), "iceberg_tables.#", "0")),
				),
			},
		},
	})
}

func TestAcc_IcebergTables_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.IcebergTables("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
