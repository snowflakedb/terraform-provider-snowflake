//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_HybridTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create operations", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			role, err := client.ContextFunctions.CurrentRole(ctx)
			require.NoError(t, err)

			assertThatObject(t, objectassert.HybridTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(role.Name()).
				HasComment("").
				HasOwnerRoleType("ROLE"))
		})

		t.Run("complete - all column and constraint options", func(t *testing.T) {
			refId, refCleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "REF_ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
			})
			t.Cleanup(refCleanup)

			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			tableComment := "integration test table"
			columnComment := "primary key column"
			collation := "en-ci"
			defaultExpr := "'default_value'"
			notNull := true

			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{
						Name:             "ID",
						Type:             sdk.DataType("NUMBER(38,0)"),
						InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey},
						Comment:          &columnComment,
					},
					{
						Name:    "NAME",
						Type:    sdk.DataType("VARCHAR(100)"),
						NotNull: &notNull,
					},
					{
						Name:         "STATUS",
						Type:         sdk.DataType("VARCHAR(50)"),
						DefaultValue: &sdk.ColumnDefaultValue{Expression: &defaultExpr},
					},
					{
						Name: "COUNTER",
						Type: sdk.DataType("NUMBER(38,0)"),
						DefaultValue: &sdk.ColumnDefaultValue{
							Identity: &sdk.ColumnIdentity{Start: 100, Increment: 5, Order: sdk.Bool(true)},
						},
					},
					{
						Name:    "NOTES",
						Type:    sdk.DataType("VARCHAR(200)"),
						Collate: &collation,
					},
					{
						Name: "REF_FK_COL",
						Type: sdk.DataType("NUMBER(38,0)"),
					},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
					{
						Name:    sdk.String("uq_name"),
						Type:    sdk.ColumnConstraintTypeUnique,
						Columns: []string{"NAME"},
					},
					{
						Type:       sdk.ColumnConstraintTypeForeignKey,
						Columns:    []string{"REF_FK_COL"},
						ForeignKey: &sdk.OutOfLineForeignKey{TableName: refId, ColumnNames: []string{"REF_ID"}},
					},
				},
				OutOfLineIndex: []sdk.HybridTableOutOfLineIndexRequest{
					{Name: "idx_status", Columns: []string{"STATUS"}, IncludeColumns: []string{"NAME"}},
				},
			}

			req := sdk.NewCreateHybridTableRequest(id, columns).WithIfNotExists(true).WithComment(tableComment)
			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)
			t.Cleanup(testClientHelper().HybridTable.DropFunc(t, id))

			assertThatObject(t, objectassert.HybridTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasComment(tableComment).
				HasOwnerRoleType("ROLE"))

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 6)

			pk := details[0]
			require.Equal(t, "ID", pk.Name)
			require.Equal(t, "NUMBER(38,0)", pk.Type)
			require.Equal(t, "COLUMN", pk.Kind)
			require.False(t, pk.IsNullable)
			require.True(t, pk.PrimaryKey)
			require.False(t, pk.UniqueKey)
			require.Equal(t, columnComment, pk.Comment)

			nameCol := details[1]
			require.Equal(t, "NAME", nameCol.Name)
			require.Equal(t, "VARCHAR(100)", nameCol.Type)
			require.Equal(t, "COLUMN", nameCol.Kind)
			require.False(t, nameCol.IsNullable)
			require.False(t, nameCol.PrimaryKey)
			require.True(t, nameCol.UniqueKey)
			require.Empty(t, nameCol.Default)
			require.Empty(t, nameCol.Check)
			require.Empty(t, nameCol.Expression)
			require.Empty(t, nameCol.Comment)
			require.Empty(t, nameCol.PolicyName)
			require.Empty(t, nameCol.PrivacyDomain)
			require.Empty(t, nameCol.SchemaEvolutionRecord)

			statusCol := details[2]
			require.Equal(t, "STATUS", statusCol.Name)
			require.Equal(t, "VARCHAR(50)", statusCol.Type)
			require.Equal(t, "COLUMN", statusCol.Kind)
			require.True(t, statusCol.IsNullable)
			require.False(t, statusCol.PrimaryKey)
			require.False(t, statusCol.UniqueKey)
			require.NotEmpty(t, statusCol.Default)
			require.Empty(t, statusCol.Check)
			require.Empty(t, statusCol.Expression)
			require.Empty(t, statusCol.Comment)
			require.Empty(t, statusCol.PolicyName)
			require.Empty(t, statusCol.PrivacyDomain)
			require.Empty(t, statusCol.SchemaEvolutionRecord)

			counterCol := details[3]
			require.Equal(t, "COUNTER", counterCol.Name)
			require.Equal(t, "NUMBER(38,0)", counterCol.Type)
			require.Equal(t, "COLUMN", counterCol.Kind)
			require.True(t, counterCol.IsNullable)
			require.False(t, counterCol.PrimaryKey)
			require.False(t, counterCol.UniqueKey)
			require.NotEmpty(t, counterCol.Default)
			require.Empty(t, counterCol.Check)
			require.Empty(t, counterCol.Expression)
			require.Empty(t, counterCol.Comment)

			notesCol := details[4]
			require.Equal(t, "NOTES", notesCol.Name)
			require.Equal(t, "VARCHAR(200) COLLATE 'en-ci'", notesCol.Type)
			require.Equal(t, "COLUMN", notesCol.Kind)
			require.True(t, notesCol.IsNullable)
			require.False(t, notesCol.PrimaryKey)
			require.False(t, notesCol.UniqueKey)

			refFkCol := details[5]
			require.Equal(t, "REF_FK_COL", refFkCol.Name)
			require.Equal(t, "NUMBER(38,0)", refFkCol.Type)
			require.Equal(t, "COLUMN", refFkCol.Kind)
			require.True(t, refFkCol.IsNullable)
			require.False(t, refFkCol.PrimaryKey)
			require.False(t, refFkCol.UniqueKey)
		})

		t.Run("composite primary key", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "PART_A", Type: sdk.DataType("NUMBER(38,0)")},
					{Name: "PART_B", Type: sdk.DataType("NUMBER(38,0)")},
					{Name: "DATA", Type: sdk.DataType("VARCHAR(100)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
					{Type: sdk.ColumnConstraintTypePrimaryKey, Columns: []string{"PART_A", "PART_B"}},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().HybridTable.DropFunc(t, id))

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.True(t, details[0].PrimaryKey)
			require.True(t, details[1].PrimaryKey)
			require.False(t, details[2].PrimaryKey)
		})

		t.Run("or replace", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns).WithComment("original"))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().HybridTable.DropFunc(t, id))

			err = client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns).WithOrReplace(true).WithComment("replaced"))
			require.NoError(t, err)

			assertThatObject(t, objectassert.HybridTable(t, id).
				HasComment("replaced"))
		})

		t.Run("if not exists", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns).WithIfNotExists(true))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().HybridTable.DropFunc(t, id))

			err = client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns).WithIfNotExists(true))
			require.NoError(t, err)
		})
	})

	t.Run("create operations - error cases", func(t *testing.T) {
		t.Run("missing primary key", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)")},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.ErrorContains(t, err, "primary key")
		})
	})

	t.Run("alter operations", func(t *testing.T) {
		t.Run("rename", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithNewName(newId))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().HybridTable.DropFunc(t, newId))

			_, err = client.HybridTables.ShowByID(ctx, id)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)

			assertThatObject(t, objectassert.HybridTable(t, newId).
				HasName(newId.Name()))
		})

		t.Run("add and drop column", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				{Name: "NAME", Type: sdk.DataType("VARCHAR(100)")},
			})
			t.Cleanup(cleanup)

			colComment := "email column"
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAddColumnAction(*sdk.NewHybridTableAddColumnActionRequest("EMAIL", sdk.DataType("VARCHAR(200)")).WithIfNotExists(true).WithComment(colComment)))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "EMAIL", details[2].Name)
			require.Equal(t, colComment, details[2].Comment)

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithDropColumnAction(*sdk.NewHybridTableDropColumnActionRequest([]string{"NAME"}).WithIfExists(true)))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "ID", details[0].Name)
			require.Equal(t, "EMAIL", details[1].Name)

			// Test ADD COLUMN with Collate and DefaultValue
			defaultVal := "'N/A'"
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAddColumnAction(*sdk.NewHybridTableAddColumnActionRequest("NOTES", sdk.DataType("VARCHAR(200)")).
					WithCollate("en-ci").
					WithDefaultValue(sdk.ColumnDefaultValue{Expression: &defaultVal})))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "NOTES", details[2].Name)
			require.NotEmpty(t, details[2].Default)

			// NOTE: InlineConstraint on ADD COLUMN is not tested — hybrid tables reject
			// adding UNIQUE/FK constraints post-creation (same limitation as ADD CONSTRAINT).
		})

		t.Run("alter column - set data type and comment", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				{Name: "NAME", Type: sdk.DataType("VARCHAR(100)")},
			})
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").WithType(sdk.DataType("VARCHAR(500)")),
				}))
			require.NoError(t, err)

			columnComment := "widened column"
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").WithComment(columnComment),
				}))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Equal(t, "VARCHAR(500)", details[1].Type)
			require.Equal(t, columnComment, details[1].Comment)

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").WithUnsetComment(true),
				}))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Empty(t, details[1].Comment)
		})

		// NOTE: ALTER TABLE UNSET COMMENT succeeds on hybrid tables.
		// Other UNSET properties may or may not be supported.

		t.Run("set properties", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			// Set comment (also tests IfExists)
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithIfExists(true).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithComment("new comment")))
			require.NoError(t, err)
			assertThatObject(t, objectassert.HybridTable(t, id).HasComment("new comment"))

			// Overwrite comment (UNSET is not supported for hybrid tables)
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithComment("updated comment")))
			require.NoError(t, err)
			assertThatObject(t, objectassert.HybridTable(t, id).HasComment("updated comment"))

			// Set data retention
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithDataRetentionTimeInDays(7)))
			require.NoError(t, err)
			assertThatObject(t, objectparametersassert.HybridTableParameters(t, id).HasDataRetentionTimeInDays(7))

			// Set max data extension
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithMaxDataExtensionTimeInDays(28)))
			require.NoError(t, err)
			assertThatObject(t, objectparametersassert.HybridTableParameters(t, id).HasMaxDataExtensionTimeInDays(28))
		})

		t.Run("show parameters", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			// Parity: client.Parameters.ShowParameters with ParametersIn{Table: id} returns the same
			// payload as the HybridTables extension method.
			parametersDirect, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
				In: &sdk.ParametersIn{Table: id},
			})
			require.NoError(t, err)
			require.NotEmpty(t, parametersDirect)

			parametersExt, err := client.HybridTables.ShowParameters(ctx, id)
			require.NoError(t, err)
			require.Equal(t, parametersDirect, parametersExt)

			// After SET, the TABLE-level Level value is returned for that parameter.
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithDataRetentionTimeInDays(3)))
			require.NoError(t, err)

			parametersAfterSet, err := client.HybridTables.ShowParameters(ctx, id)
			require.NoError(t, err)

			retention, err := collections.FindFirst(parametersAfterSet, func(p *sdk.Parameter) bool {
				return p.Key == string(sdk.ObjectParameterDataRetentionTimeInDays)
			})
			require.NoError(t, err, "DATA_RETENTION_TIME_IN_DAYS parameter must be present after SET")
			require.Equal(t, "3", (*retention).Value)
			require.Equal(t, sdk.ParameterTypeHybridTable, (*retention).Level,
				"expected Level=%q (SHOW PARAMETERS returns TABLE for hybrid tables)", sdk.ParameterTypeHybridTable)
		})

		t.Run("unset properties", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			// Arrange: set all three unsettable properties so there is something to unset.
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().
					WithComment("to be unset").
					WithDataRetentionTimeInDays(3).
					WithMaxDataExtensionTimeInDays(7)))
			require.NoError(t, err, "SET must succeed so UNSET has non-default values to target")

			// Single-property UNSET — baseline (expected to succeed regardless of multi-property capability).
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().WithComment(true)))
			require.NoError(t, err, "single-property UNSET COMMENT must succeed")
			assertThatObject(t, objectassert.HybridTable(t, id).HasComment(""))

			// Multi-property UNSET is rejected by the SDK validator (exactly-one-of in
			// AlterHybridTableOptions.validate). The validator mirrors Snowflake's own
			// rejection — a live run against Snowflake produced
			// "001003 (42000): SQL compilation error: syntax error line 1 at position <n>
			// unexpected 'UNSET'" — which is why HybridTableUnsetProperties emits one
			// UNSET keyword per field and the resource Update path issues a separate
			// ALTER per property. Keep this assertion as a regression guard: if the
			// validator is ever loosened, this test will fail and the NOTE in
			// hybrid_tables_gen.go and hybrid_tables_validations_gen.go should be revisited.
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().
					WithDataRetentionTimeInDays(true).
					WithMaxDataExtensionTimeInDays(true)))
			require.Error(t, err, "multi-property UNSET must be rejected; single-field UNSET is the only supported shape")
			require.ErrorContains(t, err, "exactly one", "error must be the client-side validator form (guards against unrelated failures masking as success)")

			// Clean up the retention properties one-at-a-time (the supported shape).
			// Verify the TABLE-level override is cleared afterwards by fetching SHOW PARAMETERS
			// and checking Level is no longer TABLE (the fallback level — ACCOUNT/DATABASE/SCHEMA/default —
			// is environment-dependent, so we only assert what we changed, not the inherited level).
			require.NoError(t, client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().WithDataRetentionTimeInDays(true))))
			require.NoError(t, client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().WithMaxDataExtensionTimeInDays(true))))
			parametersAfterUnset, err := client.HybridTables.ShowParameters(ctx, id)
			require.NoError(t, err)
			for _, p := range parametersAfterUnset {
				if p.Key == string(sdk.ObjectParameterDataRetentionTimeInDays) || p.Key == string(sdk.ObjectParameterMaxDataExtensionTimeInDays) {
					require.NotEqual(t, sdk.ParameterTypeHybridTable, p.Level, "TABLE-level override for %s must be cleared after UNSET", p.Key)
				}
			}

			// UNSET with IfExists — ensures the ALTER wrapper works for UNSET too.
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithIfExists(true).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().WithComment(true)))
			require.NoError(t, err, "UNSET with IfExists must succeed on existing table")
		})

		// NOTE: The following ALTER TABLE SET properties are NOT supported on hybrid tables and are
		// therefore absent from HybridTableSetProperties in the SDK — no tests are added for them:
		//   - CHANGE_TRACKING: hybrid tables use an internal mechanism for change tracking; this
		//     property is not user-configurable on hybrid tables.
		//   - DEFAULT_DDL_COLLATION: not supported on hybrid tables per Snowflake documentation.
		//   - ENABLE_SCHEMA_EVOLUTION: schema evolution is not supported on hybrid tables.
		//   - CONTACT / CONTACT_PURPOSE: governance/data-sharing fields not applicable to hybrid tables.
		//   - ROW_TIMESTAMP: row-level timestamp designation is not supported on hybrid tables.
		// Attempting any of these via raw SQL results in a runtime error from Snowflake. The SDK
		// intentionally omits them from HybridTableSetProperties.

		// NOTE: Hybrid tables do not support ALTER TABLE ADD UNIQUE or ADD FOREIGN KEY.
		// Snowflake returns: "Unique and foreign-key constraints can only be defined at table creation time."

		t.Run("rename and drop constraint", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "CODE", Type: sdk.DataType("VARCHAR(50)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
					{Name: sdk.String("uq_code"), Type: sdk.ColumnConstraintTypeUnique, Columns: []string{"CODE"}},
				},
			}
			_, cleanup := testClientHelper().HybridTable.CreateWithRequest(t, id, columns)
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithConstraintAction(*sdk.NewHybridTableConstraintActionRequest().
					WithRename(*sdk.NewHybridTableConstraintActionRenameRequest("uq_code", "uq_code_renamed"))))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.True(t, details[1].UniqueKey)

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithConstraintAction(*sdk.NewHybridTableConstraintActionRequest().
					WithDrop(*sdk.NewHybridTableConstraintActionDropRequest().WithConstraintName("uq_code_renamed"))))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.False(t, details[1].UniqueKey)
		})

		t.Run("alter column - drop default", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.CreateWithRequest(t,
				testClientHelper().Ids.RandomSchemaObjectIdentifier(),
				sdk.HybridTableColumnsConstraintsAndIndexesRequest{
					Columns: []sdk.HybridTableColumnRequest{
						{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
						{Name: "STATUS", Type: sdk.DataType("VARCHAR(50)"), DefaultValue: &sdk.ColumnDefaultValue{Expression: sdk.String("'ACTIVE'")}},
					},
				})
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("STATUS").WithDropDefault(true),
				}))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Empty(t, details[1].Default)
		})

		// NOTE: AlterColumn.SetDefault (SET DEFAULT seq.NEXTVAL) is not tested — it is unclear
		// whether hybrid tables support sequence-based column defaults set post-creation.
		// This needs clarification with the Snowflake table team before a test can be added.

		t.Run("drop constraint - by type", func(t *testing.T) {
			// Drop UNIQUE by type
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "EMAIL", Type: sdk.DataType("VARCHAR(255)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
					{Type: sdk.ColumnConstraintTypeUnique, Columns: []string{"EMAIL"}},
				},
			}
			_, cleanup := testClientHelper().HybridTable.CreateWithRequest(t, id, columns)
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithConstraintAction(*sdk.NewHybridTableConstraintActionRequest().
					WithDrop(*sdk.NewHybridTableConstraintActionDropRequest().WithUnique(true).WithColumns([]string{"EMAIL"}).WithCascade(true))))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.False(t, details[1].UniqueKey)

			// Drop FOREIGN KEY by type
			parentId, parentCleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "PID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
			})
			t.Cleanup(parentCleanup)

			childId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			_, childCleanup := testClientHelper().HybridTable.CreateWithRequest(t, childId, sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "CID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "PARENT_REF", Type: sdk.DataType("NUMBER(38,0)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
					{
						Type: sdk.ColumnConstraintTypeForeignKey, Columns: []string{"PARENT_REF"},
						ForeignKey: &sdk.OutOfLineForeignKey{TableName: parentId, ColumnNames: []string{"PID"}},
					},
				},
			})
			t.Cleanup(childCleanup)

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(childId).
				WithConstraintAction(*sdk.NewHybridTableConstraintActionRequest().
					WithDrop(*sdk.NewHybridTableConstraintActionDropRequest().WithForeignKey(true).WithColumns([]string{"PARENT_REF"}).WithRestrict(true))))
			require.NoError(t, err)
		})

		t.Run("clustering operations", func(t *testing.T) {
			// Snowflake currently rejects CLUSTER BY on hybrid tables (error 391407).
			// We verify the SDK generates valid SQL and Snowflake returns the expected error.
			id, cleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				{Name: "CATEGORY", Type: sdk.DataType("VARCHAR(50)")},
			})
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithClusteringAction(*sdk.NewHybridTableClusteringActionRequest().WithClusterBy([]string{"CATEGORY"})))
			require.ErrorContains(t, err, "CLUSTER BY cannot be set for a hybrid table")
		})
	})

	t.Run("show filter operations", func(t *testing.T) {
		t.Run("SHOW with LIKE - single table", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
		})

		t.Run("SHOW with LIKE - excludes non-matching", func(t *testing.T) {
			id1, cleanup1 := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup1)
			id2, cleanup2 := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup2)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id1.Name(), tables[0].Name)
			// Verify the other table is NOT returned
			for _, tbl := range tables {
				require.NotEqual(t, id2.Name(), tbl.Name)
			}
		})

		t.Run("SHOW with IN DATABASE", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithIn(sdk.TableIn{In: sdk.In{Database: sdk.NewAccountObjectIdentifier(id.DatabaseName())}}).
				WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id.Name(), tables[0].Name)
		})

		t.Run("SHOW with IN SCHEMA", func(t *testing.T) {
			id1, cleanup1 := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup1)

			// Create a second schema with a same-named table
			schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
			t.Cleanup(schemaCleanup)

			id2InOtherSchema := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			_, cleanup2 := testClientHelper().HybridTable.CreateWithRequest(t, id2InOtherSchema, columns)
			t.Cleanup(cleanup2)

			// Query IN SCHEMA should return only the table in the original schema
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithIn(sdk.TableIn{In: sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(id1.DatabaseName(), id1.SchemaName())}}).
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id1.SchemaName(), tables[0].SchemaName)
		})

		t.Run("SHOW with STARTS WITH", func(t *testing.T) {
			prefix := "HTSWTEST"
			otherPrefix := "XOTHER"
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}

			id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
			_, c1 := testClientHelper().HybridTable.CreateWithRequest(t, id1, columns)
			t.Cleanup(c1)

			id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(otherPrefix)
			_, c2 := testClientHelper().HybridTable.CreateWithRequest(t, id2, columns)
			t.Cleanup(c2)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithStartsWith(prefix))
			require.NoError(t, err)

			var found1, found2 bool
			for _, tbl := range tables {
				if tbl.Name == id1.Name() {
					found1 = true
				}
				if tbl.Name == id2.Name() {
					found2 = true
				}
			}
			require.True(t, found1, "expected table with matching prefix to be returned")
			require.False(t, found2, "expected table with non-matching prefix to be excluded")
		})

		t.Run("SHOW with LIMIT", func(t *testing.T) {
			prefix := "HTLIMTEST"
			columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
			id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
			_, c1 := testClientHelper().HybridTable.CreateWithRequest(t, id1, columns)
			t.Cleanup(c1)
			_, c2 := testClientHelper().HybridTable.CreateWithRequest(t, id2, columns)
			t.Cleanup(c2)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithStartsWith(prefix).WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
		})

		t.Run("SHOW TERSE", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithTerse(true).WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id.Name(), tables[0].Name)
			require.NotZero(t, tables[0].CreatedOn)
			require.Nil(t, tables[0].Rows)
			require.Nil(t, tables[0].Bytes)
		})
	})

	t.Run("describe operations", func(t *testing.T) {
		t.Run("all fields validated", func(t *testing.T) {
			notNull := true
			id, cleanup := testClientHelper().HybridTable.CreateWithRequest(t,
				testClientHelper().Ids.RandomSchemaObjectIdentifier(),
				sdk.HybridTableColumnsConstraintsAndIndexesRequest{
					Columns: []sdk.HybridTableColumnRequest{
						{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
						{Name: "EMAIL", Type: sdk.DataType("VARCHAR(255)"), NotNull: &notNull},
					},
					OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraintRequest{
						{Type: sdk.ColumnConstraintTypeUnique, Columns: []string{"EMAIL"}},
					},
				})
			t.Cleanup(cleanup)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)

			// Validate all 13 fields for PK column
			pk := details[0]
			require.Equal(t, "ID", pk.Name)
			require.Equal(t, "NUMBER(38,0)", pk.Type)
			require.Equal(t, "COLUMN", pk.Kind)
			require.False(t, pk.IsNullable)
			require.True(t, pk.PrimaryKey)
			require.False(t, pk.UniqueKey)
			require.Empty(t, pk.Default)
			require.Empty(t, pk.Check)
			require.Empty(t, pk.Expression)
			require.Empty(t, pk.Comment)
			require.Empty(t, pk.PolicyName)
			require.Empty(t, pk.PrivacyDomain)
			require.Empty(t, pk.SchemaEvolutionRecord)

			// Validate all 13 fields for UNIQUE + NOT NULL column
			email := details[1]
			require.Equal(t, "EMAIL", email.Name)
			require.Equal(t, "VARCHAR(255)", email.Type)
			require.Equal(t, "COLUMN", email.Kind)
			require.False(t, email.IsNullable)
			require.False(t, email.PrimaryKey)
			require.True(t, email.UniqueKey)
			require.Empty(t, email.Default)
			require.Empty(t, email.Check)
			require.Empty(t, email.Expression)
			require.Empty(t, email.Comment)
			require.Empty(t, email.PolicyName)
			require.Empty(t, email.PrivacyDomain)
			require.Empty(t, email.SchemaEvolutionRecord)
		})

		t.Run("non-existent table", func(t *testing.T) {
			_, err := client.HybridTables.Describe(ctx, testClientHelper().Ids.RandomSchemaObjectIdentifier())
			require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	})

	t.Run("show_by_id operations", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id.Name(), ht.Name)
			require.Equal(t, id.DatabaseName(), ht.DatabaseName)
			require.Equal(t, id.SchemaName(), ht.SchemaName)
		})

		t.Run("non-existent table", func(t *testing.T) {
			_, err := client.HybridTables.ShowByID(ctx, testClientHelper().Ids.RandomSchemaObjectIdentifier())
			require.ErrorIs(t, err, collections.ErrObjectNotFound)
		})
	})

	t.Run("drop operations", func(t *testing.T) {
		t.Run("basic drop", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)
		})

		t.Run("drop non-existent with IF EXISTS", func(t *testing.T) {
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).WithIfExists(true))
			require.NoError(t, err)
		})

		t.Run("drop non-existent without IF EXISTS", func(t *testing.T) {
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()))
			require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})

		t.Run("drop with CASCADE and RESTRICT", func(t *testing.T) {
			id1, cleanup1 := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup1)
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id1).WithCascade(true))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id1)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)

			id2, cleanup2 := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup2)
			err = client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id2).WithRestrict(true))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id2)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)
		})
	})

	// NOTE: INDEX operations (CREATE INDEX, DROP INDEX, SHOW INDEXES) are blocked by an SDK design
	// issue — Snowflake expects unqualified index names but the SDK generates fully qualified
	// identifiers. Index tests are omitted until the SDK identifier handling is resolved.
}
