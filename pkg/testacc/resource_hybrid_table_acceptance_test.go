//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTable_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	modelBasic := model.HybridTableFromId("test", id, columns, pk)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasComment("").
			HasColumns(columns).
			HasPrimaryKeyColumns("ID").
			HasUniqueConstraintEmpty().
			HasForeignKeyConstraintEmpty().
			HasIndexEmpty().
			HasFullyQualifiedName(id.FullyQualifiedName()),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDays(1).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
			HasMaxDataExtensionTimeInDays(1).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeDatabase),
		resourceshowoutputassert.HybridTableShowOutput(t, modelBasic.ResourceReference()).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasRows(0).
			HasBytes(0),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelBasic.ResourceReference(), 0).
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
		resourceshowoutputassert.HybridTableShowKeysOutputRow(t, modelBasic.ResourceReference(), 0).
			HasKind("PRIMARY KEY").
			HasNameNotEmpty().
			HasColumns("ID"),
	}

	modelComplete := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithDataRetentionTimeInDays(2).
		WithMaxDataExtensionTimeInDays(10)

	assertComplete := append([]assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasComment(comment).
			HasColumns(columns).
			HasPrimaryKeyColumns("ID").
			HasUniqueConstraintEmpty().
			HasForeignKeyConstraintEmpty().
			HasIndexEmpty().
			HasFullyQualifiedName(id.FullyQualifiedName()).
			HasDataRetentionTimeInDays(2).
			HasMaxDataExtensionTimeInDays(10),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDays(2).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasMaxDataExtensionTimeInDays(10).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable),
		resourceshowoutputassert.HybridTableShowOutput(t, modelComplete.ResourceReference()).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasRows(0).
			HasBytes(0),
	}, assertBasic[3:]...)

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column.0.type",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            modelBasic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithSet(
						*sdk.NewHybridTableSetPropertiesRequest().WithComment("external comment"),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// External deletion — resource dropped outside Terraform; Read detects absence
			// (ErrObjectNotFound → d.SetId("")) and the next plan recreates it.
			{
				PreConfig: func() {
					testClient().HybridTable.DropFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnBehavior(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	colsWith2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	colsWith4 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}
	// colsWith5MidInsert inserts MIDDLE_COL between NAME and EMAIL (not at the end).
	// Snowflake ADD COLUMN appends physically, so post-apply column order differs
	// from config order and the next plan is non-empty.
	colsWith5MidInsert := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "MIDDLE_COL", Type: testdatatypes.DataTypeInteger},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}
	colsWith3 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}

	modelBasic := model.HybridTableFromId("test", id, columns, pk)
	modelWith2Cols := model.HybridTableFromId("test", id, colsWith2, pk)
	modelWith4Cols := model.HybridTableFromId("test", id, colsWith4, pk)
	modelWith5ColsMidInsert := model.HybridTableFromId("test", id, colsWith5MidInsert, pk)
	modelWith3Cols := model.HybridTableFromId("test", id, colsWith3, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
						HasColumns(columns).
						HasPrimaryKeyColumns("ID"),
				),
			},
			// Add one column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWith2Cols.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWith2Cols),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelWith2Cols.ResourceReference()).
						HasColumns(colsWith2).
						HasPrimaryKeyColumns("ID"),
				),
			},
			// Add two more columns in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWith4Cols.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWith4Cols),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelWith4Cols.ResourceReference()).
						HasColumns(colsWith4).
						HasPrimaryKeyColumns("ID"),
				),
			},
			// Insert a column NOT at the end. Snowflake's ALTER TABLE ADD COLUMN appends
			// physically, so the resulting on-disk order (ID, NAME, EMAIL, AGE, MIDDLE_COL)
			// differs from the config order (ID, NAME, MIDDLE_COL, EMAIL, AGE). The apply
			// succeeds but the post-apply plan is non-empty (index drift on the TypeList).
			{
				Config:             accconfig.FromModels(t, modelWith5ColsMidInsert),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWith5ColsMidInsert.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
			},
			// Drop back to 3 columns (drops AGE and MIDDLE_COL)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWith3Cols.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWith3Cols),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelWith3Cols.ResourceReference()).
						HasColumns(colsWith3).
						HasPrimaryKeyColumns("ID"),
				),
			},
			// Drop back to single column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
						HasColumns(columns).
						HasPrimaryKeyColumns("ID"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_CompleteUseCase(t *testing.T) {
	// Create parent table externally for FK assertion.
	parentId := testClient().Ids.RandomSchemaObjectIdentifier()
	testClient().HybridTable.CreateWithRequest(t, parentId, sdk.HybridTableColumnsConstraintsAndIndexesRequest{
		Columns: []sdk.HybridTableColumnRequest{
			*sdk.NewHybridTableColumnRequest("ID", sdk.DataType("INTEGER")).
				WithInlineConstraint(sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}),
		},
	})
	t.Cleanup(testClient().HybridTable.DropFunc(t, parentId))

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()
	defaultConstant := "0"

	columnConfigs := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), NotNull: new(true)},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &defaultConstant}},
	}
	columnConfigsChanged := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "updated name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), NotNull: new(true)},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql()},
	}
	// colSigs extracts the name+type pairs needed for HybridTableFromId constructor.
	colSigs := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// FK and index are create-only; both models share the same values so the Update
	// step does not trigger ForceNew.
	uniqueConstraints := []model.HybridTableUniqueConstraintConfig{
		{Name: "my_uq", Columns: []string{"NAME"}},
		{Columns: []string{"EMAIL"}},
	}

	fkConstraints := []model.HybridTableForeignKeyConstraintConfig{
		{
			Name:       "my_fk",
			Columns:    []string{"ID"},
			TableName:  parentId.FullyQualifiedName(),
			RefColumns: []string{"ID"},
		},
	}

	indexes := []model.HybridTableIndexConfig{
		{Name: "IDX_NAME", Columns: []string{"NAME"}},
	}

	modelComplete := model.HybridTableFromId("test", id, colSigs, pk).
		WithColumnConfigs(columnConfigs).
		WithUniqueConstraints(uniqueConstraints...).
		WithForeignKeyConstraints(fkConstraints...).
		WithIndex(indexes...).
		WithComment(comment).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10)

	modelChanged := model.HybridTableFromId("test", id, colSigs, pk).
		WithColumnConfigs(columnConfigsChanged).
		WithUniqueConstraints(uniqueConstraints...).
		WithForeignKeyConstraints(fkConstraints...).
		WithIndex(indexes...).
		WithComment(changedComment).
		WithDataRetentionTimeInDays(10).
		WithMaxDataExtensionTimeInDays(20)

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column.0.type",
		"column.3.type",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create - with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasComment(comment).
						HasColumnConfigs(columnConfigs).
						HasPrimaryKeyColumns("ID").
						HasUniqueConstraints(uniqueConstraints...).
						HasForeignKeyConstraints(fkConstraints...).
						HasIndexes(indexes...).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasDataRetentionTimeInDays(5).
						HasMaxDataExtensionTimeInDays(10),
					objectparametersassert.HybridTableParameters(t, id).
						HasDataRetentionTimeInDays(5).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasMaxDataExtensionTimeInDays(10).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable),
					resourceshowoutputassert.HybridTableShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasRows(0).
						HasBytes(0),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 0).
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
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 1).
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
						HasComment("name column").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 2).
						HasName("EMAIL").
						HasType(testdatatypes.DefaultVarcharAsString).
						HasCollation("").
						HasKind("COLUMN").
						HasDefault("").
						HasIsNullable(false).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 3).
						HasName("SCORE").
						HasType("NUMBER(38,0)").
						HasCollation("").
						HasKind("COLUMN").
						HasDefault(defaultConstant).
						HasIsNullable(true).
						HasPrimaryKey(false).
						HasUniqueKey(false).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, modelComplete.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, modelComplete.ResourceReference(), 1).
						HasKind("UNIQUE").
						HasNameNotEmpty().
						HasColumns("EMAIL"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, modelComplete.ResourceReference(), 2).
						HasKind("UNIQUE").
						HasName("my_uq").
						HasColumns("NAME"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, modelComplete.ResourceReference(), 3).
						HasKind("FOREIGN KEY").
						HasName("my_fk").
						HasColumns("ID").
						HasReferencedTable(parentId.FullyQualifiedName()).
						HasReferencedColumns("ID"),
				),
			},
			// Import
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - change mutable properties
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelChanged),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelChanged.ResourceReference()).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasComment(changedComment).
						HasColumnConfigs(columnConfigsChanged).
						HasPrimaryKeyColumns("ID").
						HasUniqueConstraints(uniqueConstraints...).
						HasForeignKeyConstraints(fkConstraints...).
						HasIndexes(indexes...).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasDataRetentionTimeInDays(10).
						HasMaxDataExtensionTimeInDays(20),
					objectparametersassert.HybridTableParameters(t, id).
						HasDataRetentionTimeInDays(10).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasMaxDataExtensionTimeInDays(20).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable),
					resourceshowoutputassert.HybridTableShowOutput(t, modelChanged.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasRows(0).
						HasBytes(0),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 0).
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
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 1).
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
						HasComment("updated name column").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 2).
						HasName("EMAIL").
						HasType(testdatatypes.DefaultVarcharAsString).
						HasCollation("").
						HasKind("COLUMN").
						HasDefault("").
						HasIsNullable(false).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 3).
						HasName("SCORE").
						HasType("NUMBER(38,0)").
						HasCollation("").
						HasKind("COLUMN").
						HasDefault("").
						HasIsNullable(true).
						HasPrimaryKey(false).
						HasUniqueKey(false).
						HasCheck("").
						HasExpression("").
						HasComment("").
						HasPolicyName("").
						HasPrivacyDomain("").
						HasSchemaEvolutionRecord(""),
				),
			},
		},
	})
}

// TestAcc_HybridTable_InvalidConfig verifies that schema-level validators reject
// out-of-range values before a Snowflake connection is needed.
func TestAcc_HybridTable_InvalidConfig(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	cols := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeInteger}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(
					t,
					model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs([]model.HybridTableColumnConfig{
						{Name: "ID", Type: "INVALIDTYPE"},
					}),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid data type`),
			},
		},
	})
}

// TestAcc_HybridTable_PrimaryKeyRequiresNotNull verifies that the plan-time
// validation (requireNotNullOnPrimaryKeyColumns) rejects primary key columns
// that do not set not_null = true.
func TestAcc_HybridTable_PrimaryKeyRequiresNotNull(t *testing.T) {
	t.Run("not_null omitted on the primary key column", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeInteger}}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(
						t,
						model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs([]model.HybridTableColumnConfig{
							{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
						}),
					),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`primary key column "ID" must set not_null = true`),
				},
			},
		})
	})

	t.Run("not_null explicitly false on the primary key column", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeInteger}}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(
						t,
						model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs([]model.HybridTableColumnConfig{
							{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(false)},
						}),
					),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`primary key column "ID" must set not_null = true`),
				},
			},
		})
	})
}

func TestAcc_HybridTable_UniqueConstraint(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// Single-column unique constraint
	uq1 := model.HybridTableUniqueConstraintConfig{Columns: []string{"NAME"}}
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraints(uq1)

	// Change the unique constraint to span two columns — forces recreation
	uq2 := model.HybridTableUniqueConstraintConfig{Columns: []string{"NAME", "EMAIL"}}
	model2 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraints(uq2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a single-column unique constraint
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID").
						HasUniqueConstraints(uq1),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model1.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model1.ResourceReference(), 1).
						HasKind("UNIQUE").
						HasNameNotEmpty().
						HasColumns("NAME"),
				),
			},
			// Change the unique constraint columns — any diff on unique_constraint forces recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID").
						HasUniqueConstraints(uq2),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model2.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model2.ResourceReference(), 1).
						HasKind("UNIQUE").
						HasNameNotEmpty().
						HasColumns("NAME", "EMAIL"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ForeignKey(t *testing.T) {
	// Create parent hybrid table externally — it is not managed by Terraform in this test.
	parentId := testClient().Ids.RandomSchemaObjectIdentifier()
	testClient().HybridTable.CreateWithRequest(t, parentId, sdk.HybridTableColumnsConstraintsAndIndexesRequest{
		Columns: []sdk.HybridTableColumnRequest{
			*sdk.NewHybridTableColumnRequest("ID", sdk.DataType("INTEGER")).
				WithInlineConstraint(sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}),
		},
	})
	t.Cleanup(testClient().HybridTable.DropFunc(t, parentId))

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "PARENT_ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// Child table with FK → parent.ID
	fk := model.HybridTableForeignKeyConstraintConfig{
		Columns:    []string{"PARENT_ID"},
		TableName:  parentId.FullyQualifiedName(),
		RefColumns: []string{"ID"},
	}
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithForeignKeyConstraints(fk)

	// Child table without FK
	model2 := model.HybridTableFromId("test", id, cols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a foreign key referencing the parent table
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID").
						HasForeignKeyConstraints(fk),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model1.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model1.ResourceReference(), 1).
						HasKind("FOREIGN KEY").
						HasNameNotEmpty().
						HasColumns("PARENT_ID").
						HasReferencedTable(parentId.FullyQualifiedName()).
						HasReferencedColumns("ID"),
				),
			},
			// Remove the foreign key — any diff on foreign_key_constraint forces recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID").
						HasForeignKeyConstraintEmpty(),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model2.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
				),
			},
		},
	})
}

// TestAcc_HybridTable_ColumnDefaultVariants exercises each mutually-exclusive variant
// of the column `default` block with its own model. The `default` block has three
// sub-variants (constant, expression, sequence) and exactly one must be set per
// column. Mutual exclusivity is enforced inside buildHybridColumnDefaultValue
// in pkg/resources/hybrid_table.go (the declarative ExactlyOneOf/ConflictsWith
// options on schema fields cannot be used inside a multi-element TypeList
// because terraform-plugin-sdk/v2 rejects paths with non-zero indices at
// provider boot). Validation fires at apply time before any Snowflake call.
func TestAcc_HybridTable_ColumnDefaultVariants(t *testing.T) {
	t.Run("constant", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{
			{Name: "ID", Type: testdatatypes.DataTypeInteger},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
		}
		zero := "0"
		columnConfigs := []model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &zero}},
		}
		m := model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs(columnConfigs)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			CheckDestroy: CheckDestroy(t, resources.HybridTable),
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(t, m),
					Check: assertThat(
						t,
						resourceassert.HybridTableResource(t, m.ResourceReference()).
							HasColumnConfigs(columnConfigs),
					),
				},
			},
		})
	})

	t.Run("expression", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{
			{Name: "ID", Type: testdatatypes.DataTypeInteger},
			{Name: "CREATED_AT", Type: testdatatypes.DataTypeTimestampLTZ},
		}
		expr := "CURRENT_TIMESTAMP()"
		columnConfigs := []model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
			{Name: "CREATED_AT", Type: testdatatypes.DataTypeTimestampLTZ.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Expression: &expr}},
		}
		m := model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs(columnConfigs)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			CheckDestroy: CheckDestroy(t, resources.HybridTable),
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(t, m),
					Check: assertThat(
						t,
						resourceassert.HybridTableResource(t, m.ResourceReference()).
							HasColumnConfigs(columnConfigs),
					),
				},
			},
		})
	})

	t.Run("sequence", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{
			{Name: "ID", Type: testdatatypes.DataTypeInteger},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
		}
		seqId, cleanup := testClient().Sequence.Create(t)
		t.Cleanup(cleanup)
		seqFQN := seqId.FullyQualifiedName()
		columnConfigs := []model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Sequence: &seqFQN}},
		}
		m := model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs(columnConfigs)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			CheckDestroy: CheckDestroy(t, resources.HybridTable),
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(t, m),
					Check: assertThat(
						t,
						resourceassert.HybridTableResource(t, m.ResourceReference()).
							HasColumnConfigs(columnConfigs),
					),
				},
			},
		})
	})

	// Negative test: setting more than one of {constant, expression, sequence}
	// in the same default block must be rejected by buildHybridColumnDefaultValue.
	// Validation fires at apply time (Create runs the build helper before any
	// Snowflake call), so the apply errors out before any resource is created.
	t.Run("conflicting fields", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{
			{Name: "ID", Type: testdatatypes.DataTypeInteger},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
		}
		zero := "0"
		expr := "0"
		m := model.HybridTableFromId("test", id, cols, pk).
			WithColumnConfigs([]model.HybridTableColumnConfig{
				{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
				{
					Name: "SCORE",
					Type: testdatatypes.DataTypeInteger.ToSql(),
					Default: &model.HybridTableColumnDefaultConfig{
						Constant:   &zero,
						Expression: &expr,
					},
				},
			})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			Steps: []resource.TestStep{
				{
					Config:      accconfig.FromModels(t, m),
					ExpectError: regexp.MustCompile(`default block must have exactly one of "constant", "expression", or "sequence" set`),
				},
			},
		})
	})
}

func TestAcc_HybridTable_PrimaryKeyForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// Single-column PK
	model1 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}})

	// Composite PK — any change to primary_key_constraint forces recreation
	model2 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}, {Name: "NAME"}})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with single-column PK
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model1.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID"),
				),
			},
			// Change to composite PK → ForceNew (DestroyBeforeCreate)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyColumns("ID", "NAME"),
					resourceshowoutputassert.HybridTableShowKeysOutputRow(t, model2.ResourceReference(), 0).
						HasKind("PRIMARY KEY").
						HasNameNotEmpty().
						HasColumns("ID", "NAME"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ExternalColumnChanges(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	cols2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	cols3WithEmail := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}

	model2 := model.HybridTableFromId("test", id, cols2, pk)
	model3 := model.HybridTableFromId("test", id, cols3WithEmail, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// 1. Create with 2 columns via Terraform.
			{
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2).
						HasPrimaryKeyColumns("ID"),
				),
			},
			// 2. Externally ADD a column (EMAIL). Config at cols2; expect Update to drop EMAIL.
			{
				PreConfig: func() {
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithAddColumnAction(
						*sdk.NewHybridTableAddColumnActionRequest("EMAIL", sdk.DataType("VARCHAR")),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2),
				),
			},
			// 3. Externally DROP a column (NAME). Config at cols2; expect Update to re-add NAME.
			{
				PreConfig: func() {
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithDropColumnAction(
						*sdk.NewHybridTableDropColumnActionRequest([]string{"NAME"}),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2),
				),
			},
			// 4. Externally MODIFY a column comment (non-ForceNew). Config at cols2;
			//    expect Update to reset the comment.
			{
				PreConfig: func() {
					newComment := "external comment"
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
						*sdk.NewHybridTableAlterColumnActionRequest("NAME").WithComment(newComment),
					}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2),
				),
			},
			// 5. Config moves to cols3WithEmail (adds EMAIL) while externally a column comment
			//    is changed — complex combined-change scenario from jmichalak's review comment.
			{
				PreConfig: func() {
					newComment := "external comment 2"
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
						*sdk.NewHybridTableAlterColumnActionRequest("NAME").WithComment(newComment),
					}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(cols3WithEmail),
				),
			},
			// 6. At cols3WithEmail baseline, externally DROP NAME (a non-trailing column).
			//    Config still at model3; expect Update to re-add NAME. Snowflake's ADD COLUMN
			//    appends physically, so the resulting on-disk order is (ID, EMAIL, NAME)
			//    while the config order is (ID, NAME, EMAIL) — the next plan is non-empty.
			//    Verifies drift detection still fires on a non-trailing column drop in a
			//    larger table (covers the "more complex" external-drift scenarios from
			//    jmichalak's review comment).
			{
				PreConfig: func() {
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithDropColumnAction(
						*sdk.NewHybridTableDropColumnActionRequest([]string{"NAME"}),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config:             accconfig.FromModels(t, model3),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAcc_HybridTable_ColumnNotNullForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// not_null omitted on NAME -> resolves to the schema default (false, nullable).
	// ID is the primary key, so it must declare not_null = true.
	defaultColumnConfigs := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql()},
	}
	defaultModel := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(defaultColumnConfigs)
	assertDefault := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, defaultModel.ResourceReference()).
			HasColumnConfigs(defaultColumnConfigs),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, defaultModel.ResourceReference(), 1).
			HasIsNullable(true),
	}

	// not_null explicitly false on NAME -> same meaning as omitted (nullable).
	nullableColumnConfigs := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), NotNull: new(false)},
	}
	nullableModel := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(nullableColumnConfigs)
	assertNullable := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, nullableModel.ResourceReference()).
			HasColumnConfigs(nullableColumnConfigs),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, nullableModel.ResourceReference(), 1).
			HasIsNullable(true),
	}

	// not_null explicitly true on NAME -> NOT NULL column.
	nonNullColumnConfigs := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), NotNull: new(true)},
	}
	nonNullModel := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(nonNullColumnConfigs)
	assertNonNull := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, nonNullModel.ResourceReference()).
			HasColumnConfigs(nonNullColumnConfigs),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, nonNullModel.ResourceReference(), 1).
			HasIsNullable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create - not_null omitted (nullable).
			{
				Config: accconfig.FromModels(t, defaultModel),
				Check:  assertThat(t, assertDefault...),
			},
			// Omitted -> explicit false: same meaning, no diff.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, nullableModel),
				Check:  assertThat(t, assertNullable...),
			},
			// false -> true: SDK cannot ALTER COLUMN SET NOT NULL, so recreate.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(nonNullModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, nonNullModel),
				Check:  assertThat(t, assertNonNull...),
			},
			// true -> false: not_null change on an existing column always forces recreation.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(nullableModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, nullableModel),
				Check:  assertThat(t, assertNullable...),
			},
			// External change detection.
			{
				PreConfig: func() {
					testClient().HybridTable.DropFunc(t, id)()
					testClient().HybridTable.CreateWithRequest(t, id, sdk.HybridTableColumnsConstraintsAndIndexesRequest{
						Columns: []sdk.HybridTableColumnRequest{
							*sdk.NewHybridTableColumnRequest("ID", sdk.DataType(testdatatypes.DataTypeInteger.ToSql())).
								WithInlineConstraint(sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}),
							*sdk.NewHybridTableColumnRequest("NAME", sdk.DataType(testdatatypes.DataTypeVarchar.ToSql())).
								WithNotNull(true),
						},
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(nullableModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, nullableModel),
				Check:  assertThat(t, assertNullable...),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnCollateForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// NAME with collate='en'
	columnConfigs1 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "en"},
	}
	model1 := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(columnConfigs1)

	// NAME collate changed to 'FR' — must force recreation
	columnConfigs2 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql(), NotNull: new(true)},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "FR"},
	}
	model2 := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(columnConfigs2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with NAME collate='en'
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumnConfigs(columnConfigs1),
				),
			},
			// Change NAME collate='FR' — expect DestroyBeforeCreate (ForceNew)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumnConfigs(columnConfigs2),
				),
			},
		},
	})
}

func TestAcc_HybridTable_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	renamedComment := random.Comment()
	modelBasic := model.HybridTableFromId("test", id, columns, pk)
	modelRenamed := model.HybridTableFromId("test", newId, columns, pk).WithComment(renamedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasFullyQualifiedName(id.FullyQualifiedName()),
				),
			},
			// Rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelRenamed),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelRenamed.ResourceReference()).
						HasName(newId.Name()).
						HasDatabase(newId.DatabaseName()).
						HasSchema(newId.SchemaName()).
						HasFullyQualifiedName(newId.FullyQualifiedName()).
						HasComment(renamedComment),
				),
			},
		},
	})
}

func TestAcc_Experimental_HybridTable_HierarchyRenames_MoveToAnotherSchema(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaXName := testClient().Ids.Alpha()
	schemaYName := testClient().Ids.Alpha()
	tableName := testClient().Ids.Alpha()

	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)
	databaseModel := model.DatabaseWithParametersSet("db", databaseId.Name())

	schemaModelX := model.SchemaWithImplicitDatabaseDependency("schemaX", schemaXName, databaseModel)
	schemaModelY := model.SchemaWithImplicitDatabaseDependency("schemaY", schemaYName, databaseModel)

	columns := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber}}
	pks := []sdk.TableColumnSignature{{Name: "ID"}}
	tableModelBefore := model.HybridTableWithImplicitDependencies("test", tableName, columns, pks, schemaModelX, databaseModel)
	tableModelAfter := model.HybridTableWithImplicitDependencies("test", tableName, columns, pks, schemaModelY, databaseModel)

	id := sdk.NewSchemaObjectIdentifier(databaseId.Name(), schemaXName, tableName)
	expectedNewTableId := sdk.NewSchemaObjectIdentifier(databaseId.Name(), schemaYName, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create database, both schemas, and the hybrid table in schema X
			{
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelX, schemaModelY, tableModelBefore),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, tableModelBefore.ResourceReference()).
						HasName(tableName).
						HasDatabase(databaseId.Name()).
						HasSchema(schemaXName).
						HasFullyQualifiedName(id.FullyQualifiedName()),
				),
			},
			// Move the hybrid table to schema Y
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, databaseModel, schemaModelX, schemaModelY, tableModelAfter),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, tableModelAfter.ResourceReference()).
						HasName(tableName).
						HasDatabase(databaseId.Name()).
						HasSchema(schemaYName).
						HasFullyQualifiedName(expectedNewTableId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Experimental_HybridTable_HierarchyRenames_Disabled_ForceRecreation(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	newId := testClient().Ids.NewSchemaObjectIdentifierInSchema(id.Name(), schemaId)

	columns := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber}}
	pks := []sdk.TableColumnSignature{{Name: "ID"}}
	tableModelBefore := model.HybridTableFromId("test", id, columns, pks)
	tableModelAfter := model.HybridTableFromId("test", newId, columns, pks)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, tableModelBefore),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, tableModelBefore.ResourceReference()).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasFullyQualifiedName(id.FullyQualifiedName()),
				),
			},
			// Change database — should force recreation.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tableModelAfter.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:             accconfig.FromModels(t, tableModelAfter),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
