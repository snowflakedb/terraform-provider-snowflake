//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasColumnNullable(0, true). // PK column: state keeps the configured/default value (true); Snowflake-level NOT NULL is asserted via describe_output
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
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
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasRowsNil().
			HasBytesNil(),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelBasic.ResourceReference(), 0).
			HasName("ID").
			HasIsNullable(false).
			HasPrimaryKey(true).
			HasUniqueKey(false).
			HasDefault("").
			HasComment(""),
	}

	modelComplete := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithDataRetentionTimeInDays(1).
		WithMaxDataExtensionTimeInDays(10)

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString(comment).
			HasDataRetentionTimeInDaysString("1").
			HasMaxDataExtensionTimeInDaysString("10").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDays(1).
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
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasRowsNil().
			HasBytesNil(),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 0).
			HasName("ID").
			HasIsNullable(false).
			HasPrimaryKey(true).
			HasUniqueKey(false).
			HasDefault("").
			HasComment(""),
	}

	assertAfterUnset := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasColumnNullable(0, true). // PK column: state keeps the configured/default value (true); Snowflake-level NOT NULL is asserted via describe_output
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
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
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasRowsNil().
			HasBytesNil(),
		resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelBasic.ResourceReference(), 0).
			HasName("ID").
			HasIsNullable(false).
			HasPrimaryKey(true).
			HasUniqueKey(false).
			HasDefault("").
			HasComment(""),
	}

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column",
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
				Check:  assertThat(t, assertAfterUnset...),
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
				Check:  assertThat(t, assertAfterUnset...),
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
			// Destroy
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
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
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID"),
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

	columnConfigs := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
	}
	columnConfigsChanged := []model.HybridTableColumnConfig{
		{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "updated name column"},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
	}
	// colSigs extracts the name+type pairs needed for HybridTableFromId constructor.
	colSigs := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// FK and index are create-only; both models share the same values so the Update
	// step does not trigger ForceNew.
	uniqueConstraints := []model.HybridTableUniqueConstraintConfig{
		{Name: "my_uq", Columns: []string{"NAME"}},
		{Columns: []string{"EMAIL"}},
	}

	modelComplete := model.HybridTableFromId("test", id, colSigs, pk).
		WithColumnConfigs(columnConfigs).
		WithUniqueConstraints(uniqueConstraints...).
		WithNamedForeignKey("my_fk", []string{"ID"}, parentId.FullyQualifiedName(), []string{"ID"}).
		// IDX_NAME uses lowercase "name" to exercise the case-suppression path end-to-end
		// (SHOW INDEXES returns uppercase; indexHash normalises before hashing).
		WithIndex(model.HybridTableIndexConfig{Name: "IDX_NAME", Columns: []string{"name"}}).
		WithComment(comment).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10)

	modelChanged := model.HybridTableFromId("test", id, colSigs, pk).
		WithColumnConfigs(columnConfigsChanged).
		WithUniqueConstraints(uniqueConstraints...).
		WithNamedForeignKey("my_fk", []string{"ID"}, parentId.FullyQualifiedName(), []string{"ID"}).
		WithIndex(model.HybridTableIndexConfig{Name: "IDX_NAME", Columns: []string{"name"}}).
		WithComment(changedComment).
		WithDataRetentionTimeInDays(10).
		WithMaxDataExtensionTimeInDays(20)

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column",
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
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasDataRetentionTimeInDaysString("5").
						HasMaxDataExtensionTimeInDaysString("10").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasColumnConfigs(columnConfigs).
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(2).
						HasIndexCount(1),
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
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasRowsNil().
						HasBytesNil(),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 0).
						HasName("ID").
						HasIsNullable(false).
						HasPrimaryKey(true).
						HasUniqueKey(false).
						HasDefault("").
						HasComment(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 1).
						HasName("NAME").
						HasIsNullable(true).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasDefault("").
						HasComment("name column"),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelComplete.ResourceReference(), 2).
						HasName("EMAIL").
						HasIsNullable(false).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasDefault("").
						HasComment(""),
					// Constraint name readback: PK auto-generated, named UNIQUE round-trips,
					// anonymous UNIQUE gets SYS_CONSTRAINT_-prefixed name, FK name preserved.
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "primary_key.0.name")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "unique_constraint.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(modelComplete.ResourceReference(), "unique_constraint.*", map[string]string{
						"name": "my_uq",
					})),
					assert.Check(checkUniqueConstraintNamesReadBack(modelComplete.ResourceReference())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "foreign_key.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(modelComplete.ResourceReference(), "foreign_key.*", map[string]string{
						"name":                  "my_fk",
						"references.#":          "1",
						"references.0.table_id": parentId.FullyQualifiedName(),
					})),
					// Index readback: lowercase config column must round-trip (case-suppression).
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(modelComplete.ResourceReference(), "index.*", map[string]string{
						"name":      "IDX_NAME",
						"columns.#": "1",
						"columns.0": "name",
					})),
				),
			},
			// Empty plan — regression guard for constraint-name drift and index case-suppression.
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDataRetentionTimeInDaysString("10").
						HasMaxDataExtensionTimeInDaysString("20").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasColumnConfigs(columnConfigsChanged).
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(2).
						HasIndexCount(1),
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
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasRowsNil().
						HasBytesNil(),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 0).
						HasName("ID").
						HasIsNullable(false).
						HasPrimaryKey(true).
						HasUniqueKey(false).
						HasDefault("").
						HasComment(""),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 1).
						HasName("NAME").
						HasIsNullable(true).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasDefault("").
						HasComment("updated name column"),
					resourceshowoutputassert.HybridTableDescribeOutputRow(t, modelChanged.ResourceReference(), 2).
						HasName("EMAIL").
						HasIsNullable(false).
						HasPrimaryKey(false).
						HasUniqueKey(true).
						HasDefault("").
						HasComment(""),
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
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraint([]string{"NAME"})

	// Change the unique constraint to span two columns — forces recreation
	model2 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraint([]string{"NAME", "EMAIL"})

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
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(1),
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
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(1),
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
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithForeignKey([]string{"PARENT_ID"}, parentId.FullyQualifiedName(), []string{"ID"})

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
						HasPrimaryKeyKeys("ID").
						HasForeignKeyCount(1),
				),
			},
			// Remove the foreign key — any diff on foreign_key forces recreation
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
						HasPrimaryKeyKeys("ID").
						HasForeignKeyCount(0),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnDefault(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
	}

	zero := "0"

	// SCORE has a constant default of 0
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &zero}},
		})

	// SCORE has no default
	model2 := model.HybridTableFromId("test", id, baseCols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a constant default on SCORE
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols).
						HasColumnDefaultConstant(1, "0"),
				),
			},
			// Drop the default → in-place update, no recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols).
						HasColumnNoDefault(1),
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
//
// One subtest per variant per jmichalak's review comment (thread 3188827027):
// separate models for mutually-exclusive field sets. The "conflicting fields"
// subtest asserts the build-helper validation fires when more than one of
// {constant, expression, sequence} is set in the same default block.
func TestAcc_HybridTable_ColumnDefaultVariants(t *testing.T) {
	t.Run("constant", func(t *testing.T) {
		id := testClient().Ids.RandomSchemaObjectIdentifier()
		pk := []sdk.TableColumnSignature{{Name: "ID"}}
		cols := []sdk.TableColumnSignature{
			{Name: "ID", Type: testdatatypes.DataTypeInteger},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
		}
		zero := "0"
		m := model.HybridTableFromId("test", id, cols, pk).
			WithColumnConfigs([]model.HybridTableColumnConfig{
				{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
				{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &zero}},
			})

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
							HasColumns(cols).
							HasColumnDefaultConstant(1, "0"),
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
		m := model.HybridTableFromId("test", id, cols, pk).
			WithColumnConfigs([]model.HybridTableColumnConfig{
				{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
				{Name: "CREATED_AT", Type: testdatatypes.DataTypeTimestampLTZ.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Expression: &expr}},
			})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			CheckDestroy: CheckDestroy(t, resources.HybridTable),
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(t, m),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(m.ResourceReference(), "column.1.default.0.expression", expr),
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
		m := model.HybridTableFromId("test", id, cols, pk).
			WithColumnConfigs([]model.HybridTableColumnConfig{
				{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
				{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Sequence: &seqFQN}},
			})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireAbove(tfversion.Version1_5_0),
			},
			CheckDestroy: CheckDestroy(t, resources.HybridTable),
			Steps: []resource.TestStep{
				{
					Config: accconfig.FromModels(t, m),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(m.ResourceReference(), "column.1.default.#", "1"),
						resource.TestCheckResourceAttr(m.ResourceReference(), "column.1.default.0.sequence", seqFQN),
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
	model1 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}}).
		WithPrimaryKeyNames("ID")

	// Composite PK — any change to primary_key forces recreation
	model2 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}}).
		WithPrimaryKeyNames("ID", "NAME")

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
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID", "NAME"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnAdd(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	cols1 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	cols2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	cols3 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}
	// cols4 inserts MIDDLE_COL between NAME and EMAIL (not at the end of cols3).
	// Snowflake ADD COLUMN appends physically, so the post-apply state column order
	// differs from the config order, and the next plan must be non-empty.
	cols4 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "MIDDLE_COL", Type: testdatatypes.DataTypeInteger},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}

	model1 := model.HybridTableFromId("test", id, cols1, pk)
	model2 := model.HybridTableFromId("test", id, cols2, pk)
	model3 := model.HybridTableFromId("test", id, cols3, pk)
	model4 := model.HybridTableFromId("test", id, cols4, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with single column
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols1).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Add one column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Add two more columns in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(cols3).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Insert a column NOT at the end. Snowflake's ALTER TABLE ADD COLUMN
			// appends physically, so the resulting on-disk order
			// (ID, NAME, EMAIL, AGE, MIDDLE_COL) differs from the config order
			// (ID, NAME, MIDDLE_COL, EMAIL, AGE). The apply succeeds (the column
			// is added), but the post-apply plan is non-empty: subsequent indices
			// in the TypeList show drift. Achieving a true mid-list insertion
			// would require recreation, which the resource does not currently do.
			{
				Config:             accconfig.FromModels(t, model4),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model4.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func TestAcc_HybridTable_ColumnDrop(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	cols1 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}
	cols2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	cols3 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}

	model1 := model.HybridTableFromId("test", id, cols1, pk)
	model2 := model.HybridTableFromId("test", id, cols2, pk)
	model3 := model.HybridTableFromId("test", id, cols3, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with four columns
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols1).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Drop one column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Drop two more columns in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(cols3).
						HasPrimaryKeyKeys("ID"),
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
						HasPrimaryKeyKeys("ID"),
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

func TestAcc_HybridTable_ColumnAlterComment(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// No comment on NAME
	model1 := model.HybridTableFromId("test", id, baseCols, pk)

	// Set comment on NAME
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "this is a name column"},
		})

	// Change comment on NAME
	model3 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "updated comment"},
		})

	// Unset comment on NAME (back to no comment)
	model4 := model.HybridTableFromId("test", id, baseCols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with two columns, NAME has no comment
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, ""),
				),
			},
			// Set comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, "this is a name column"),
				),
			},
			// Change comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, "updated comment"),
				),
			},
			// Unset comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model4.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model4),
				Check: assertThat(
					t,
					resourceassert.HybridTableResource(t, model4.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, ""),
				),
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
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(true)},
		})

	// NAME changed to nullable=false — must force recreation
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
		})

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
						HasColumns(baseCols).
						HasColumnNullable(1, true),
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
						HasColumns(baseCols).
						HasColumnNullable(1, false),
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
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "en"},
		})

	// NAME collate changed to 'fr' — must force recreation
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "fr"},
		})

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
						HasColumns(baseCols),
				),
			},
			// Change NAME collate='fr' — expect DestroyBeforeCreate (ForceNew)
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
						HasColumns(baseCols),
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
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					objectassert.HybridTable(t, id).
						HasName(id.Name()),
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
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCommentString(renamedComment),
					objectassert.HybridTable(t, newId).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
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

// checkUniqueConstraintNamesReadBack asserts that every unique_constraint element in
// state has a name that is either the explicit "my_uq" or a server-generated name with
// the "SYS_CONSTRAINT_" prefix. An empty name (which would mean GetConstraints failed to
// read back the auto-generated name of the unnamed UNIQUE) fails the check.
func checkUniqueConstraintNamesReadBack(resourceRef string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceRef]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceRef)
		}

		// TypeSet name attributes look like "unique_constraint.<hash>.name".
		nameKey := regexp.MustCompile(`^unique_constraint\.\d+\.name$`)
		var found int
		for key, value := range rs.Primary.Attributes {
			if !nameKey.MatchString(key) {
				continue
			}
			found++
			if value != "my_uq" && !strings.HasPrefix(value, "SYS_CONSTRAINT_") {
				return fmt.Errorf("unique_constraint name %q (at %s) is neither %q nor SYS_CONSTRAINT_-prefixed; auto-generated name was not read back", value, key, "my_uq")
			}
		}
		if found != 2 {
			return fmt.Errorf("expected 2 unique_constraint name attributes, found %d", found)
		}
		return nil
	}
}

// TestAcc_HybridTable_CollateCaseInsensitive verifies that a config-supplied
// collate of "en-ci" produces no spurious diff even if DESCRIBE returns it
// as "EN-CI" or some other case variant. Reconciliation comes from
// ignoreCaseSuppressFunc on the field plus Read-time substitution that
// preserves the user's spelling when DESCRIBE is case-equivalent.
func TestAcc_HybridTable_CollateCaseInsensitive(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}
	// HybridTableFromId / WithColumn does not expose per-column collate, so use
	// the richer WithColumnConfigs builder that does.
	tableModel := model.HybridTableFromId("test", id, columns, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "en-ci"},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel),
			},
			{
				Config: accconfig.FromModels(t, tableModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
