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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
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
	}

	modelComplete := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithDataRetentionTimeInDays(2).
		WithMaxDataExtensionTimeInDays(10)

	assertComplete := []assert.TestCheckFuncProvider{
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
	}

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column.0.type",
		// PK columns come back from DESCRIBE as NOT NULL.
		"column.0.nullable",
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
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &defaultConstant}},
	}
	columnConfigsChanged := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "updated name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
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
		// PK columns come back from DESCRIBE as NOT NULL.
		"column.0.nullable",
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
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
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
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
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
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
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
				{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
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

func TestAcc_HybridTable_ColumnNullableForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// NAME is explicitly nullable=true
	columnConfigs1 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(true)},
	}
	model1 := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(columnConfigs1)

	// NAME changed to nullable=false — must force recreation
	columnConfigs2 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
	}
	model2 := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(columnConfigs2)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with NAME nullable=true
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumnConfigs(columnConfigs1),
				),
			},
			// Change NAME nullable=false — expect DestroyBeforeCreate (ForceNew)
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

func TestAcc_HybridTable_ColumnCollateForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// NAME with collate='en'
	columnConfigs1 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "en"},
	}
	model1 := model.HybridTableFromId("test", id, baseCols, pk).WithColumnConfigs(columnConfigs1)

	// NAME collate changed to 'FR' — must force recreation
	columnConfigs2 := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
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

// TestAcc_HybridTable_PKNullableNoSpurious verifies that a primary-key column
// declared as nullable=true (the schema default) does not produce a spurious
// diff after Read, even though Snowflake silently enforces NOT NULL on PK
// columns and DESCRIBE reports null="N". The reconciliation happens via
// Read-time substitution in buildHybridColumnStateFromDescribe.
func TestAcc_HybridTable_PKNullableNoSpurious(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}
	model := model.HybridTableFromId("test", id, columns, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model),
			},
			// A second apply with no config change must produce a no-op plan.
			{
				Config: accconfig.FromModels(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
