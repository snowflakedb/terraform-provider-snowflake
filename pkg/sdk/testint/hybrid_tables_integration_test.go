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

			assertions := objectassert.HybridTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasComment("").
				HasOwnerRoleType("ROLE")
			assertThatObject(t, assertions)
		})

		t.Run("complete - all column and constraint options", func(t *testing.T) {
			// Create a reference table for FK constraint
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
						Name: "STATUS",
						Type: sdk.DataType("VARCHAR(50)"),
						DefaultValue: &sdk.ColumnDefaultValue{
							Expression: &defaultExpr,
						},
					},
					{
						Name: "COUNTER",
						Type: sdk.DataType("NUMBER(38,0)"),
						DefaultValue: &sdk.ColumnDefaultValue{
							Identity: &sdk.ColumnIdentity{
								Start:     1,
								Increment: 1,
							},
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
						Type:    sdk.ColumnConstraintTypeForeignKey,
						Columns: []string{"REF_FK_COL"},
						ForeignKey: &sdk.OutOfLineForeignKey{
							TableName:   refId,
							ColumnNames: []string{"REF_ID"},
						},
					},
				},
				OutOfLineIndex: []sdk.HybridTableOutOfLineIndexRequest{
					{
						Name:    "idx_status",
						Columns: []string{"STATUS"},
					},
				},
			}

			req := sdk.NewCreateHybridTableRequest(id, columns).
				WithIfNotExists(true).
				WithComment(tableComment)
			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Verify SHOW output via generated assertions
			assertions := objectassert.HybridTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasComment(tableComment).
				HasOwnerRoleType("ROLE")
			assertThatObject(t, assertions)

			// Verify DESCRIBE output
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 6)

			// PK column
			pk := details[0]
			require.Equal(t, "ID", pk.Name)
			require.Contains(t, pk.Type, "NUMBER")
			require.Equal(t, "COLUMN", pk.Kind)
			require.False(t, pk.IsNullable)
			require.True(t, pk.PrimaryKey)
			require.False(t, pk.UniqueKey)
			require.Equal(t, columnComment, pk.Comment)

			// NOT NULL column with UNIQUE out-of-line constraint
			nameCol := details[1]
			require.Equal(t, "NAME", nameCol.Name)
			require.Contains(t, nameCol.Type, "VARCHAR")
			require.False(t, nameCol.IsNullable)
			require.False(t, nameCol.PrimaryKey)
			require.True(t, nameCol.UniqueKey)

			// DEFAULT expression column
			statusCol := details[2]
			require.Equal(t, "STATUS", statusCol.Name)
			require.NotEmpty(t, statusCol.Default)
			require.True(t, statusCol.IsNullable)

			// IDENTITY column
			counterCol := details[3]
			require.Equal(t, "COUNTER", counterCol.Name)
			require.Contains(t, counterCol.Type, "NUMBER")

			// COLLATE column
			notesCol := details[4]
			require.Equal(t, "NOTES", notesCol.Name)
			require.Contains(t, notesCol.Type, "VARCHAR")
			require.True(t, notesCol.IsNullable)
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

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// OR REPLACE with different comment
			err = client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns).
				WithOrReplace(true).
				WithComment("replaced"))
			require.NoError(t, err)

			assertions := objectassert.HybridTable(t, id).
				HasComment("replaced")
			assertThatObject(t, assertions)
		})
	})

	t.Run("alter operations", func(t *testing.T) {
		t.Run("rename", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithNewName(newId))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(newId).WithIfExists(true))
				require.NoError(t, err)
			})

			// Verify old name is gone
			_, err = client.HybridTables.ShowByID(ctx, id)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)

			// Verify new name works
			assertions := objectassert.HybridTable(t, newId).
				HasName(newId.Name())
			assertThatObject(t, assertions)
		})

		t.Run("add and drop column", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				{Name: "NAME", Type: sdk.DataType("VARCHAR(100)")},
			})
			t.Cleanup(cleanup)

			// Add column with comment
			colComment := "email column"
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAddColumnAction(*sdk.NewHybridTableAddColumnActionRequest("EMAIL", sdk.DataType("VARCHAR(200)")).
					WithComment(colComment)))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "EMAIL", details[2].Name)
			require.Contains(t, details[2].Type, "VARCHAR")
			require.Equal(t, colComment, details[2].Comment)

			// Drop column
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

			// Set data type (widen VARCHAR)
			columnComment := "widened column"
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").
						WithType(sdk.DataType("VARCHAR(500)")),
				}))
			require.NoError(t, err)

			// Set column comment
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").
						WithComment(columnComment),
				}))
			require.NoError(t, err)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Contains(t, details[1].Type, "VARCHAR(500)")
			require.Equal(t, columnComment, details[1].Comment)

			// Unset column comment
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{
					*sdk.NewHybridTableAlterColumnActionRequest("NAME").
						WithUnsetComment(true),
				}))
			require.NoError(t, err)

			details, err = client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Empty(t, details[1].Comment)
		})

		t.Run("set and unset properties", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			// Set properties
			newComment := "updated comment"
			retentionDays := 7
			err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().
					WithComment(newComment).
					WithDataRetentionTimeInDays(retentionDays)))
			require.NoError(t, err)

			// Verify comment via SHOW
			assertions := objectassert.HybridTable(t, id).
				HasComment(newComment)
			assertThatObject(t, assertions)

			// Update comment to a different value
			updatedComment := "second update"
			err = client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().
					WithComment(updatedComment)))
			require.NoError(t, err)

			assertions = objectassert.HybridTable(t, id).
				HasComment(updatedComment)
			assertThatObject(t, assertions)
		})
	})

	t.Run("show filter operations", func(t *testing.T) {
		t.Run("SHOW with LIKE pattern", func(t *testing.T) {
			id1, cleanup1 := testClientHelper().HybridTable.CreateWithColumns(t, []sdk.HybridTableColumnRequest{
				{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
			})
			t.Cleanup(cleanup1)

			id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("OTHER_PREFIX")
			columns2 := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
				Columns: []sdk.HybridTableColumnRequest{
					{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id2, columns2))
			require.NoError(t, err)
			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id2).WithIfExists(true))
				require.NoError(t, err)
			})

			// LIKE matching only the first table
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id1.Name(), tables[0].Name)
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
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithIn(sdk.TableIn{In: sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())}}).
				WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
			require.Equal(t, id.Name(), tables[0].Name)
		})

		t.Run("SHOW with STARTS WITH", func(t *testing.T) {
			prefix := "HTSWTEST"
			id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
			id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("XOTHER")

			for _, id := range []sdk.SchemaObjectIdentifier{id1, id2} {
				columns := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
					Columns: []sdk.HybridTableColumnRequest{
						{Name: "ID", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					},
				}
				err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
				require.NoError(t, err)
				t.Cleanup(func() {
					err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
					require.NoError(t, err)
				})
			}

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithStartsWith(prefix))
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
			require.True(t, found1, "expected to find table with matching prefix")
			require.False(t, found2, "should not find table with non-matching prefix")
		})

		t.Run("SHOW with LIMIT", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
			require.NoError(t, err)
			require.Len(t, tables, 1)
		})

		t.Run("SHOW TERSE", func(t *testing.T) {
			id, cleanup := testClientHelper().HybridTable.Create(t)
			t.Cleanup(cleanup)

			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().
				WithTerse(true).
				WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
			require.NoError(t, err)
			require.Len(t, tables, 1)

			ht := tables[0]
			require.Equal(t, id.Name(), ht.Name)
			require.Equal(t, id.DatabaseName(), ht.DatabaseName)
			require.Equal(t, id.SchemaName(), ht.SchemaName)
			require.NotZero(t, ht.CreatedOn)

			// TERSE-excluded fields must be nil pointers
			require.Nil(t, ht.Rows)
			require.Nil(t, ht.Bytes)
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
						{
							Type:    sdk.ColumnConstraintTypeUnique,
							Columns: []string{"EMAIL"},
						},
					},
				})
			t.Cleanup(cleanup)

			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)

			// Validate all HybridTableDetails fields for the PK column
			pk := details[0]
			require.Equal(t, "ID", pk.Name)
			require.Contains(t, pk.Type, "NUMBER")
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

			// Validate the UNIQUE + NOT NULL column
			email := details[1]
			require.Equal(t, "EMAIL", email.Name)
			require.Contains(t, email.Type, "VARCHAR")
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
		})

		t.Run("non-existent table", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			_, err := client.HybridTables.Describe(ctx, id)
			require.Error(t, err)
			require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	})

	t.Run("show_by_id operations", func(t *testing.T) {
		t.Run("non-existent table", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			_, err := client.HybridTables.ShowByID(ctx, id)
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
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
			require.NoError(t, err)
		})

		t.Run("drop non-existent without IF EXISTS", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id))
			require.Error(t, err)
			require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})

		t.Run("drop with CASCADE and RESTRICT", func(t *testing.T) {
			// CASCADE
			id1, _ := testClientHelper().HybridTable.Create(t)
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id1).WithCascade(true))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id1)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)

			// RESTRICT
			id2, _ := testClientHelper().HybridTable.Create(t)
			err = client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id2).WithRestrict(true))
			require.NoError(t, err)
			_, err = client.HybridTables.ShowByID(ctx, id2)
			require.ErrorIs(t, err, collections.ErrObjectNotFound)
		})
	})
}
