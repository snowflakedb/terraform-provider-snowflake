//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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

			assertThatObject(t, objectassert.HybridTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
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
							Identity: &sdk.ColumnIdentity{Start: 1, Increment: 1},
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
					{Name: "idx_status", Columns: []string{"STATUS"}},
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
			require.False(t, nameCol.IsNullable)
			require.True(t, nameCol.UniqueKey)

			statusCol := details[2]
			require.Equal(t, "STATUS", statusCol.Name)
			require.Equal(t, "VARCHAR(50)", statusCol.Type)
			require.NotEmpty(t, statusCol.Default)
			require.True(t, statusCol.IsNullable)

			counterCol := details[3]
			require.Equal(t, "COUNTER", counterCol.Name)
			require.Equal(t, "NUMBER(38,0)", counterCol.Type)
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
			require.Error(t, err)
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
				WithAddColumnAction(*sdk.NewHybridTableAddColumnActionRequest("EMAIL", sdk.DataType("VARCHAR(200)")).WithComment(colComment)))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "EMAIL", details[2].Name)
			require.Equal(t, colComment, details[2].Comment)

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithDropColumnAction(*sdk.NewHybridTableDropColumnActionRequest([]string{"NAME"})))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "ID", details[0].Name)
			require.Equal(t, "EMAIL", details[1].Name)
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

		// NOTE: Hybrid tables do not support ALTER TABLE ... UNSET.
		// To "unset" a property, use SET with a new value.

		t.Run("set comment", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithComment("new comment")))
			require.NoError(t, err)

			assertThatObject(t, objectassert.HybridTable(t, id).HasComment("new comment"))

			// Overwrite comment with a different value (UNSET is not supported for hybrid tables)
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithComment("updated comment")))
			require.NoError(t, err)

			assertThatObject(t, objectassert.HybridTable(t, id).HasComment("updated comment"))
		})

		t.Run("set data_retention_time_in_days", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithDataRetentionTimeInDays(7)))
			require.NoError(t, err)
		})

		t.Run("set max_data_extension_time_in_days", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().WithMaxDataExtensionTimeInDays(28)))
			require.NoError(t, err)
		})

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

			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithConstraintAction(*sdk.NewHybridTableConstraintActionRequest().
					WithDrop(*sdk.NewHybridTableConstraintActionDropRequest().WithConstraintName("uq_code_renamed"))))
			require.NoError(t, err)
		})
	})

	t.Run("show filter operations", func(t *testing.T) {
		t.Run("SHOW without filter", func(t *testing.T) {
			_, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest())
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(tables), 1)
		})

		t.Run("SHOW with LIKE", func(t *testing.T) {
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
			id, _ := testClientHelper().HybridTable.Create(t)
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
			id1, _ := testClientHelper().HybridTable.Create(t)
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id1).WithCascade(true))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id1)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)

			id2, _ := testClientHelper().HybridTable.Create(t)
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
