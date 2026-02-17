//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTables_BasicFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	table1 := model.HybridTable("test1", testClient().Ids.DatabaseName(), testClient().Ids.SchemaName(), id1.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
		}).
		WithPrimaryKeyColumns("id")

	table2 := model.HybridTable("test2", testClient().Ids.DatabaseName(), testClient().Ids.SchemaName(), id2.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
		}).
		WithPrimaryKeyColumns("id")

	dsLikeExact := datasourcemodel.HybridTables("test").
		WithLike(id1.Name()).
		WithDependsOn(table1.ResourceReference(), table2.ResourceReference())

	dsLikePrefix := datasourcemodel.HybridTables("test").
		WithLike(prefix+"%").
		WithDependsOn(table1.ResourceReference(), table2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// LIKE with wildcard (prefix%)
			{
				Config: config.FromModels(t, table1, table2, dsLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLikePrefix.DatasourceReference(), "hybrid_tables.#", "2"),
				),
			},
			// LIKE exact match
			{
				Config: config.FromModels(t, table1, table2, dsLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "hybrid_tables.#", "1"),
					resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "hybrid_tables.0.name", id1.Name()),
				),
			},
		},
	})
}

func TestAcc_HybridTables_InFilters(t *testing.T) {
	dbName := testClient().Ids.DatabaseName()
	schemaName := testClient().Ids.SchemaName()
	tableName := testClient().Ids.Alpha()

	table := model.HybridTable("test", dbName, schemaName, tableName).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
		}).
		WithPrimaryKeyColumns("id")

	dsInSchema := datasourcemodel.HybridTables("test").
		WithInSchema(fmt.Sprintf("%s.%s", dbName, schemaName)).
		WithLike(tableName).
		WithDependsOn(table.ResourceReference())

	dsInDatabase := datasourcemodel.HybridTables("test").
		WithInDatabase(dbName).
		WithLike(tableName).
		WithDependsOn(table.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// IN SCHEMA
			{
				Config: config.FromModels(t, table, dsInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "hybrid_tables.#", "1"),
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "hybrid_tables.0.name", tableName),
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "hybrid_tables.0.database_name", dbName),
					resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "hybrid_tables.0.schema_name", schemaName),
				),
			},
			// IN DATABASE
			{
				Config: config.FromModels(t, table, dsInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsInDatabase.DatasourceReference(), "hybrid_tables.#", "1"),
					resource.TestCheckResourceAttr(dsInDatabase.DatasourceReference(), "hybrid_tables.0.name", tableName),
				),
			},
		},
	})
}

func TestAcc_HybridTables_StartsWith(t *testing.T) {
	prefix := "TEST_" + random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	table := model.HybridTable("test", testClient().Ids.DatabaseName(), testClient().Ids.SchemaName(), id1.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
		}).
		WithPrimaryKeyColumns("id")

	dsStartsWith := datasourcemodel.HybridTables("test").
		WithStartsWith(prefix).
		WithDependsOn(table.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, table, dsStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "hybrid_tables.#", "1"),
					resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "hybrid_tables.0.name", id1.Name()),
				),
			},
		},
	})
}

func TestAcc_HybridTables_Limit(t *testing.T) {
	prefix := random.AlphaN(4)
	ids := make([]string, 3)
	tables := make([]*model.HybridTableModel, 3)

	for i := 0; i < 3; i++ {
		id := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
		ids[i] = id.Name()
		tables[i] = model.HybridTable(fmt.Sprintf("test%d", i), testClient().Ids.DatabaseName(), testClient().Ids.SchemaName(), id.Name()).
			WithColumnDescs([]model.ColumnDesc{
				{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
			}).
			WithPrimaryKeyColumns("id")
	}

	dsLimit := datasourcemodel.HybridTables("test").
		WithLike(prefix+"%").
		WithLimit(2).
		WithDependsOn(tables[0].ResourceReference(), tables[1].ResourceReference(), tables[2].ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, tables[0], tables[1], tables[2], dsLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "hybrid_tables.#", "2"),
				),
			},
		},
	})
}

func TestAcc_HybridTables_CompleteUseCase(t *testing.T) {
	tableName := testClient().Ids.Alpha()
	dbName := testClient().Ids.DatabaseName()
	schemaName := testClient().Ids.SchemaName()
	comment := random.Comment()

	table := model.HybridTable("test", dbName, schemaName, tableName).
		WithComment(comment).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
			{Name: "name", DataType: "VARCHAR(100)", Nullable: model.Bool(true)},
		}).
		WithPrimaryKeyColumns("id")

	ds := datasourcemodel.HybridTables("test").
		WithLike(tableName).
		WithDependsOn(table.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, table, ds),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.#", "1"),
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.0.name", tableName),
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.0.database_name", dbName),
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.0.comment", comment),
					resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "hybrid_tables.0.created_on"),
					resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "hybrid_tables.0.owner"),
					resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "hybrid_tables.0.rows"),
					resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "hybrid_tables.0.bytes"),
					resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "hybrid_tables.0.owner_role_type"),
				),
			},
		},
	})
}

func TestAcc_HybridTables_EmptyResults(t *testing.T) {
	nonExistentPattern := "NONEXISTENT_TABLE_" + random.AlphaN(10)

	ds := datasourcemodel.HybridTables("test").
		WithLike(nonExistentPattern)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, ds),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds.DatasourceReference(), "hybrid_tables.#", "0"),
				),
			},
		},
	})
}
