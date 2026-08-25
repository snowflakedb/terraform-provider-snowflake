//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTables_CompleteUseCase(t *testing.T) {
	parentId, parentCleanup := testClient().HybridTable.CreateWithRequest(t, testClient().Ids.RandomSchemaObjectIdentifier(), sdk.HybridTableColumnsConstraintsAndIndexesRequest{
		Columns: []sdk.HybridTableColumnRequest{
			*sdk.NewHybridTableColumnRequest("ID", sdk.DataType("INTEGER")).
				WithInlineConstraint(sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}),
		},
	})
	t.Cleanup(parentCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	objectModel := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithUniqueConstraints(model.HybridTableUniqueConstraintConfig{
			Name:    new("my_uq"),
			Columns: []string{"NAME"},
		}).
		WithForeignKeyConstraints(model.HybridTableForeignKeyConstraintConfig{
			Name:       new("my_fk"),
			Columns:    []string{"ID"},
			TableName:  parentId.FullyQualifiedName(),
			RefColumns: []string{"ID"},
		})

	datasourceModel := datasourcemodel.HybridTables("test").
		WithLike(id.Name()).
		WithInSchema(id.SchemaId()).
		WithWithKeys(true).
		WithWithIndexes(true).
		WithDependsOn(objectModel.ResourceReference())
	datasourceModelWithoutOptionals := datasourcemodel.HybridTables("test").
		WithLike(id.Name()).
		WithInSchema(id.SchemaId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(objectModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.HybridTablesDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, objectModel, datasourceModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "hybrid_tables.#", "1")),
					resourceparametersassert.HybridTablesDatasourceParameters(t, datasourceModel.DatasourceReference()).
						HasDataRetentionTimeInDays(1).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
						HasMaxDataExtensionTimeInDays(1).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeDatabase),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "hybrid_tables.0.describe_output.#", "2")),
					resourceshowoutputassert.HybridTablesDatasourceDescribeOutputRow(t, datasourceModel.DatasourceReference(), 0).
						HasName("ID").
						HasType("NUMBER(38,0)").
						HasCollation("").
						HasKind("COLUMN").
						HasDefault("").
						HasIsNullable(false).
						HasPrimaryKey(true).
						HasUniqueKey(false).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTablesDatasourceDescribeOutputRow(t, datasourceModel.DatasourceReference(), 1).
						HasName("NAME").
						HasType(testdatatypes.DefaultVarcharAsString).
						HasCollation("").
						HasKind("COLUMN").
						HasDefault("").
						HasIsNullable(true).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "hybrid_tables.0.show_keys_output.#", "3")),
					resourceshowoutputassert.HybridTablesDatasourceShowKeysOutputRow(t, datasourceModel.DatasourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
					resourceshowoutputassert.HybridTablesDatasourceShowKeysOutputRow(t, datasourceModel.DatasourceReference(), 1).
						HasKind("UNIQUE").
						HasName("my_uq").
						HasColumns("NAME"),
					resourceshowoutputassert.HybridTablesDatasourceShowKeysOutputRow(t, datasourceModel.DatasourceReference(), 2).
						HasKind("FOREIGN KEY").
						HasName("my_fk").
						HasColumns("ID").
						HasReferencedTable(parentId.FullyQualifiedName()).
						HasReferencedColumns("ID"),
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "hybrid_tables.0.show_indexes.#", "3")),
				),
			},
			{
				Config: accconfig.FromModels(t, objectModel, datasourceModelWithoutOptionals),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutOptionals.DatasourceReference(), "hybrid_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutOptionals.DatasourceReference(), "hybrid_tables.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutOptionals.DatasourceReference(), "hybrid_tables.0.parameters.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutOptionals.DatasourceReference(), "hybrid_tables.0.show_keys_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutOptionals.DatasourceReference(), "hybrid_tables.0.show_indexes.#", "0")),
				),
			},
		},
	})
}

func TestAcc_HybridTables_BasicUseCase_DifferentFiltering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	model1 := model.HybridTableFromId("t1", id1, columns, pk)
	model2 := model.HybridTableFromId("t2", id2, columns, pk)
	model3 := model.HybridTableFromId("t3", id3, columns, pk)

	hybridTablesModelLikeFirst := datasourcemodel.HybridTables("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	hybridTablesModelLikePrefix := datasourcemodel.HybridTables("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	hybridTablesModelStartsWith := datasourcemodel.HybridTables("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	hybridTablesModelLimit := datasourcemodel.HybridTables("test").
		WithRowsAndFrom(1, prefix).
		WithInDatabase(id1.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	hybridTablesModelInSchema := datasourcemodel.HybridTables("test").
		WithInSchema(id1.SchemaId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, hybridTablesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(hybridTablesModelLikeFirst.DatasourceReference(), "hybrid_tables.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, hybridTablesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(hybridTablesModelLikePrefix.DatasourceReference(), "hybrid_tables.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, hybridTablesModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(hybridTablesModelStartsWith.DatasourceReference(), "hybrid_tables.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, hybridTablesModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(hybridTablesModelLimit.DatasourceReference(), "hybrid_tables.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, hybridTablesModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(hybridTablesModelInSchema.DatasourceReference(), "hybrid_tables.#", "2"),
				),
			},
		},
	})
}

func TestAcc_HybridTables_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.HybridTables("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
				PlanOnly:    true,
			},
		},
	})
}
