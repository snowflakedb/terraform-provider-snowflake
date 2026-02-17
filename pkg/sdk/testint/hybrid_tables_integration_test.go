//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_HybridTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// Helper: Assert SHOW HYBRID TABLES result with all 9 columns
	assertHybridTableShowResult := func(t *testing.T, ht *sdk.HybridTable, id sdk.SchemaObjectIdentifier, expectedComment string) {
		t.Helper()
		require.NotZero(t, ht.CreatedOn)
		require.Equal(t, id.Name(), ht.Name)
		require.Equal(t, id.DatabaseName(), ht.DatabaseName)
		require.Equal(t, id.SchemaName(), ht.SchemaName)
		require.NotEmpty(t, ht.Owner)
		require.GreaterOrEqual(t, ht.Rows, 0)
		require.GreaterOrEqual(t, ht.Bytes, 0)
		require.Equal(t, expectedComment, ht.Comment)
		require.NotEmpty(t, ht.OwnerRoleType)
	}

	// Helper: Create basic hybrid table with single PK column
	createHybridTableBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		columns := sdk.HybridTableColumnsConstraintsAndIndexes{
			Columns: []sdk.HybridTableColumn{
				{
					Name: "id",
					Type: sdk.DataType("NUMBER(38,0)"),
					InlineConstraint: &sdk.HybridTableColumnInlineConstraint{
						Type: sdk.ColumnConstraintTypePrimaryKey,
					},
				},
			},
		}

		req := sdk.NewCreateHybridTableRequest(id, columns)
		err := client.HybridTables.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
			require.NoError(t, err)
		})

		return id
	}

	t.Run("create operations", func(t *testing.T) {
		t.Run("basic with single primary key", func(t *testing.T) {
			id := createHybridTableBasic(t)

			// Validate via SHOW
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			assertHybridTableShowResult(t, ht, id, "")

			// Validate via DESCRIBE - should have 1 column
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 1)
			require.Equal(t, "ID", details[0].Name) // Snowflake returns uppercase
			require.Contains(t, details[0].Type, "NUMBER")
			require.Equal(t, "COLUMN", details[0].Kind)
			require.Equal(t, "N", details[0].IsNullable) // NOT NULL due to PK
			require.Equal(t, "Y", details[0].PrimaryKey)
		})

		t.Run("with table COMMENT", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			comment := "Test hybrid table with comment"

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{
						Name: "id",
						Type: sdk.DataType("NUMBER(38,0)"),
						InlineConstraint: &sdk.HybridTableColumnInlineConstraint{
							Type: sdk.ColumnConstraintTypePrimaryKey,
						},
					},
				},
			}

			req := sdk.NewCreateHybridTableRequest(id, columns).
				WithComment(comment)

			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate comment via SHOW
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			assertHybridTableShowResult(t, ht, id, comment)
		})

		t.Run("composite primary key", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "order_id", Type: sdk.DataType("NUMBER(38,0)")},
					{Name: "line_item", Type: sdk.DataType("NUMBER(38,0)")},
					{Name: "product", Type: sdk.DataType("VARCHAR(100)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraint{
					{
						Type:    sdk.ColumnConstraintTypePrimaryKey,
						Columns: []string{"order_id", "line_item"},
					},
				},
			}

			req := sdk.NewCreateHybridTableRequest(id, columns)
			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE - both columns should have primary key
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "Y", details[0].PrimaryKey)
			require.Equal(t, "Y", details[1].PrimaryKey)
			require.Equal(t, "N", details[2].PrimaryKey)
		})

		t.Run("column with DEFAULT", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "status", Type: sdk.DataType("VARCHAR(50)"), DefaultValue: &sdk.ColumnDefaultValue{Expression: sdk.String("'PENDING'")}},
					{Name: "created_at", Type: sdk.DataType("TIMESTAMP_NTZ"), DefaultValue: &sdk.ColumnDefaultValue{Expression: sdk.String("CURRENT_TIMESTAMP()")}},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 3)
			require.Equal(t, "STATUS", details[1].Name)
			require.NotEmpty(t, details[1].Default)
			require.Contains(t, details[1].Default, "PENDING")
			require.Equal(t, "CREATED_AT", details[2].Name)
			require.NotEmpty(t, details[2].Default)
		})

		t.Run("column with NOT NULL", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			notNullValue := true
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "required_field", Type: sdk.DataType("VARCHAR(100)"), NotNull: &notNullValue},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "REQUIRED_FIELD", details[1].Name)
			require.Equal(t, "N", details[1].IsNullable) // NOT NULL
		})

		t.Run("column with IDENTITY", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			startNum := 100
			incrementNum := 10
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{
						Name: "id",
						Type: sdk.DataType("NUMBER(38,0)"),
						DefaultValue: &sdk.ColumnDefaultValue{
							Identity: &sdk.ColumnIdentity{
								Start:     startNum,
								Increment: incrementNum,
							},
						},
						InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey},
					},
					{Name: "data", Type: sdk.DataType("VARCHAR(100)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "ID", details[0].Name)
			require.NotEmpty(t, details[0].Default)
			require.Contains(t, details[0].Default, "IDENTITY")
		})

		t.Run("multiple data types", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "varchar_col", Type: sdk.DataType("VARCHAR(255)")},
					{Name: "date_col", Type: sdk.DataType("DATE")},
					{Name: "timestamp_col", Type: sdk.DataType("TIMESTAMP_NTZ")},
					{Name: "boolean_col", Type: sdk.DataType("BOOLEAN")},
					{Name: "variant_col", Type: sdk.DataType("VARIANT")},
					{Name: "decimal_col", Type: sdk.DataType("DECIMAL(10,2)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE - check all data types
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 7)
			require.Equal(t, "ID", details[0].Name)
			require.Contains(t, details[0].Type, "NUMBER")
			require.Equal(t, "VARCHAR_COL", details[1].Name)
			require.Contains(t, details[1].Type, "VARCHAR")
			require.Equal(t, "DATE_COL", details[2].Name)
			require.Equal(t, "DATE", details[2].Type)
			require.Equal(t, "TIMESTAMP_COL", details[3].Name)
			require.Contains(t, details[3].Type, "TIMESTAMP")
			require.Equal(t, "BOOLEAN_COL", details[4].Name)
			require.Equal(t, "BOOLEAN", details[4].Type)
			require.Equal(t, "VARIANT_COL", details[5].Name)
			require.Equal(t, "VARIANT", details[5].Type)
		})

		t.Run("out-of-line UNIQUE constraint", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "email", Type: sdk.DataType("VARCHAR(255)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraint{
					{
						Type:    sdk.ColumnConstraintTypeUnique,
						Columns: []string{"email"},
					},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "EMAIL", details[1].Name)
			require.Equal(t, "Y", details[1].UniqueKey) // UNIQUE constraint
		})

		t.Run("out-of-line FOREIGN KEY", func(t *testing.T) {
			// Create parent table first
			parentId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			parentColumns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "name", Type: sdk.DataType("VARCHAR(100)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(parentId, parentColumns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(parentId).WithIfExists(true))
				require.NoError(t, err)
			})

			// Create child table with foreign key
			childId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			childColumns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "parent_id", Type: sdk.DataType("NUMBER(38,0)")},
				},
				OutOfLineConstraint: []sdk.HybridTableOutOfLineConstraint{
					{
						Type:    sdk.ColumnConstraintTypeForeignKey,
						Columns: []string{"parent_id"},
						ForeignKey: &sdk.OutOfLineForeignKey{
							TableName:   parentId,
							ColumnNames: []string{"id"},
						},
					},
				},
			}

			err = client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(childId, childColumns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(childId).WithIfExists(true))
				require.NoError(t, err)
			})

			// Validate via DESCRIBE - check foreign key relationship
			details, err := client.HybridTables.Describe(ctx, childId)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "PARENT_ID", details[1].Name)
			// Note: Foreign key info might not show in DESCRIBE output
		})
	})

	t.Run("alter operations", func(t *testing.T) {
		t.Run("ALTER COLUMN SET COMMENT", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "name", Type: sdk.DataType("VARCHAR(100)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Set column comment
			columnComment := "Name column comment"
			alterReq := sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction(*sdk.NewHybridTableAlterColumnActionRequest("NAME"). // Snowflake uses uppercase
														WithComment(columnComment))

			err = client.HybridTables.Alter(ctx, alterReq)
			require.NoError(t, err)

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "NAME", details[1].Name)
			require.Equal(t, columnComment, details[1].Comment)
		})

		t.Run("ALTER COLUMN UNSET COMMENT", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			commentValue := "Initial comment"
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "name", Type: sdk.DataType("VARCHAR(100)"), Comment: &commentValue},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Unset column comment
			alterReq := sdk.NewAlterHybridTableRequest(id).
				WithAlterColumnAction(*sdk.NewHybridTableAlterColumnActionRequest("NAME").
					WithUnsetComment(true))

			err = client.HybridTables.Alter(ctx, alterReq)
			require.NoError(t, err)

			// Validate via DESCRIBE
			details, err := client.HybridTables.Describe(ctx, id)
			require.NoError(t, err)
			require.Len(t, details, 2)
			require.Equal(t, "NAME", details[1].Name)
			require.Empty(t, details[1].Comment)
		})

		t.Run("SET DATA_RETENTION_TIME_IN_DAYS", func(t *testing.T) {
			id := createHybridTableBasic(t)
			newRetention := 7

			req := sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().
					WithDataRetentionTimeInDays(newRetention))

			err := client.HybridTables.Alter(ctx, req)
			require.NoError(t, err)

			// Validate via SHOW - note: DATA_RETENTION_TIME_IN_DAYS not in SHOW output
			// This test verifies the command succeeds
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id.Name(), ht.Name)
		})

		t.Run("SET COMMENT", func(t *testing.T) {
			id := createHybridTableBasic(t)
			newComment := "Updated table comment"

			req := sdk.NewAlterHybridTableRequest(id).
				WithSet(*sdk.NewHybridTableSetPropertiesRequest().
					WithComment(newComment))

			err := client.HybridTables.Alter(ctx, req)
			require.NoError(t, err)

			// Validate via SHOW
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, newComment, ht.Comment)
		})

		t.Run("UNSET COMMENT", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}

			req := sdk.NewCreateHybridTableRequest(id, columns).
				WithComment("Initial comment")

			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Unset comment
			unsetReq := sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().
					WithComment(true))

			err = client.HybridTables.Alter(ctx, unsetReq)
			require.NoError(t, err)

			// Validate via SHOW
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Empty(t, ht.Comment)
		})

		t.Run("UNSET DATA_RETENTION_TIME_IN_DAYS", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}

			retentionDays := 7
			req := sdk.NewCreateHybridTableRequest(id, columns).
				WithDataRetentionTimeInDays(retentionDays)

			err := client.HybridTables.Create(ctx, req)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// Unset retention time
			unsetReq := sdk.NewAlterHybridTableRequest(id).
				WithUnset(*sdk.NewHybridTableUnsetPropertiesRequest().
					WithDataRetentionTimeInDays(true))

			err = client.HybridTables.Alter(ctx, unsetReq)
			require.NoError(t, err)

			// Validate via SHOW - command should succeed
			ht, err := client.HybridTables.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id.Name(), ht.Name)
		})
	})

	t.Run("index operations", func(t *testing.T) {
		t.Run("CREATE INDEX basic", func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			indexId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			// Create hybrid table first
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "status", Type: sdk.DataType("VARCHAR(50)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(tableId, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(tableId).WithIfExists(true))
				require.NoError(t, err)
			})

			// Create index
			indexReq := sdk.NewCreateHybridTableIndexRequest(indexId, tableId, []string{"status"})
			err = client.HybridTables.CreateIndex(ctx, indexReq)
			require.NoError(t, err)

			// Validate via SHOW INDEXES
			indexFilter := &sdk.ShowHybridTableIndexIn{
				Table: &tableId,
			}
			indexes, err := client.HybridTables.ShowIndexes(ctx, sdk.NewShowHybridTableIndexesRequest().WithIn(*indexFilter))
			require.NoError(t, err)
			require.NotEmpty(t, indexes)

			// Find our index (there will also be a system primary key index)
			var found *sdk.HybridTableIndex
			for i := range indexes {
				if indexes[i].Name == indexId.Name() {
					found = &indexes[i]
					break
				}
			}
			require.NotNil(t, found, "Created index not found in SHOW INDEXES")
			require.Equal(t, indexId.Name(), found.Name)
			require.Contains(t, found.Columns, "STATUS") // Snowflake returns uppercase
		})

		t.Run("CREATE INDEX with INCLUDE", func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			indexId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			// Create hybrid table
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "customer_id", Type: sdk.DataType("NUMBER(38,0)")},
					{Name: "order_date", Type: sdk.DataType("DATE")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(tableId, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(tableId).WithIfExists(true))
				require.NoError(t, err)
			})

			// Create index with INCLUDE
			indexReq := sdk.NewCreateHybridTableIndexRequest(indexId, tableId, []string{"customer_id"}).
				WithIncludeColumns([]string{"order_date"})
			err = client.HybridTables.CreateIndex(ctx, indexReq)
			require.NoError(t, err)

			// Validate via SHOW INDEXES
			indexFilter := &sdk.ShowHybridTableIndexIn{
				Table: &tableId,
			}
			indexes, err := client.HybridTables.ShowIndexes(ctx, sdk.NewShowHybridTableIndexesRequest().WithIn(*indexFilter))
			require.NoError(t, err)

			// Find our index
			var found *sdk.HybridTableIndex
			for i := range indexes {
				if indexes[i].Name == indexId.Name() {
					found = &indexes[i]
					break
				}
			}
			require.NotNil(t, found)
			require.Contains(t, found.Columns, "CUSTOMER_ID")
			require.NotEmpty(t, found.IncludedColumns)
			require.Contains(t, found.IncludedColumns, "ORDER_DATE")
		})

		t.Run("DROP INDEX", func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			indexId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			// Create hybrid table and index
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "status", Type: sdk.DataType("VARCHAR(50)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(tableId, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(tableId).WithIfExists(true))
				require.NoError(t, err)
			})

			indexReq := sdk.NewCreateHybridTableIndexRequest(indexId, tableId, []string{"status"})
			err = client.HybridTables.CreateIndex(ctx, indexReq)
			require.NoError(t, err)

			// Drop index using standalone DROP INDEX command
			dropReq := sdk.NewDropHybridTableIndexRequest(indexId)
			err = client.HybridTables.DropIndex(ctx, dropReq)
			require.NoError(t, err)

			// Verify index is gone
			indexFilter := &sdk.ShowHybridTableIndexIn{
				Table: &tableId,
			}
			indexes, err := client.HybridTables.ShowIndexes(ctx, sdk.NewShowHybridTableIndexesRequest().WithIn(*indexFilter))
			require.NoError(t, err)
			for _, idx := range indexes {
				require.NotEqual(t, indexId.Name(), idx.Name, "Index should be dropped")
			}
		})

		t.Run("SHOW INDEXES validates all columns", func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			indexId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			// Create hybrid table and index
			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					{Name: "status", Type: sdk.DataType("VARCHAR(50)")},
				},
			}

			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(tableId, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(tableId).WithIfExists(true))
				require.NoError(t, err)
			})

			indexReq := sdk.NewCreateHybridTableIndexRequest(indexId, tableId, []string{"status"})
			err = client.HybridTables.CreateIndex(ctx, indexReq)
			require.NoError(t, err)

			// Show indexes and validate ALL 10 columns
			indexFilter := &sdk.ShowHybridTableIndexIn{
				Table: &tableId,
			}
			indexes, err := client.HybridTables.ShowIndexes(ctx, sdk.NewShowHybridTableIndexesRequest().WithIn(*indexFilter))
			require.NoError(t, err)
			require.NotEmpty(t, indexes)

			// Find our index and validate all columns
			var found *sdk.HybridTableIndex
			for i := range indexes {
				if indexes[i].Name == indexId.Name() {
					found = &indexes[i]
					break
				}
			}
			require.NotNil(t, found)

			// Validate all 10 columns from SHOW INDEXES output
			require.NotZero(t, found.CreatedOn)
			require.Equal(t, indexId.Name(), found.Name)
			require.False(t, found.IsUnique) // Basic index is not unique
			require.NotEmpty(t, found.Columns)
			// IncludedColumns may be empty for basic index
			require.Equal(t, tableId.Name(), found.TableName)
			require.Equal(t, tableId.DatabaseName(), found.DatabaseName)
			require.Equal(t, tableId.SchemaName(), found.SchemaName)
			require.NotEmpty(t, found.Owner)
			require.NotEmpty(t, found.OwnerRoleType)
		})
	})

	t.Run("show filter operations", func(t *testing.T) {
		t.Run("SHOW with LIKE pattern", func(t *testing.T) {
			id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_like_abc")
			id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_like_xyz")

			// Create two tables with distinct prefixes
			for _, id := range []sdk.SchemaObjectIdentifier{id1, id2} {
				columns := sdk.HybridTableColumnsConstraintsAndIndexes{
					Columns: []sdk.HybridTableColumn{
						{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					},
				}
				err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
				require.NoError(t, err)

				t.Cleanup(func() {
					err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
					require.NoError(t, err)
				})
			}

			// SHOW with LIKE pattern matching first table
			likePattern := &sdk.Like{Pattern: sdk.String("TEST_LIKE_ABC%")}
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithLike(*likePattern))
			require.NoError(t, err)

			// Should find id1 but not id2
			var found1, found2 bool
			for _, t := range tables {
				if t.Name == id1.Name() {
					found1 = true
				}
				if t.Name == id2.Name() {
					found2 = true
				}
			}
			require.True(t, found1, "Expected to find first table")
			require.False(t, found2, "Should not find second table with LIKE filter")
		})

		t.Run("SHOW with IN DATABASE", func(t *testing.T) {
			id := createHybridTableBasic(t)

			// SHOW with IN DATABASE filter
			inClause := &sdk.In{Database: sdk.NewAccountObjectIdentifier(id.DatabaseName())}
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithIn(*inClause))
			require.NoError(t, err)

			// Should find our table
			var found bool
			for _, t := range tables {
				if t.Name == id.Name() {
					found = true
					break
				}
			}
			require.True(t, found, "Expected to find table in database")
		})

		t.Run("SHOW with IN SCHEMA", func(t *testing.T) {
			id := createHybridTableBasic(t)

			// SHOW with IN SCHEMA filter
			inClause := &sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())}
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithIn(*inClause))
			require.NoError(t, err)

			// Should find our table
			var found bool
			for _, t := range tables {
				if t.Name == id.Name() {
					found = true
					break
				}
			}
			require.True(t, found, "Expected to find table in schema")
		})

		t.Run("SHOW with STARTS WITH", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("SWTEST")

			columns := sdk.HybridTableColumnsConstraintsAndIndexes{
				Columns: []sdk.HybridTableColumn{
					{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
				},
			}
			err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
			require.NoError(t, err)

			t.Cleanup(func() {
				err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
				require.NoError(t, err)
			})

			// SHOW with STARTS WITH filter - use prefix that matches our table
			startsWithValue := "SWTEST"
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithStartsWith(startsWithValue))
			require.NoError(t, err)

			// Should find our table
			var found bool
			for _, t := range tables {
				if t.Name == id.Name() {
					found = true
					break
				}
			}
			require.True(t, found, "Expected to find table with STARTS WITH filter")
		})

		t.Run("SHOW with LIMIT", func(t *testing.T) {
			// Create multiple tables
			var ids []sdk.SchemaObjectIdentifier
			for i := 0; i < 3; i++ {
				id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
				ids = append(ids, id)

				columns := sdk.HybridTableColumnsConstraintsAndIndexes{
					Columns: []sdk.HybridTableColumn{
						{Name: "id", Type: sdk.DataType("NUMBER(38,0)"), InlineConstraint: &sdk.HybridTableColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}},
					},
				}
				err := client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(id, columns))
				require.NoError(t, err)

				t.Cleanup(func() {
					err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
					require.NoError(t, err)
				})
			}

			// SHOW with LIMIT 1
			limitClause := &sdk.LimitFrom{Rows: sdk.Int(1)}
			tables, err := client.HybridTables.Show(ctx, sdk.NewShowHybridTableRequest().WithLimit(*limitClause))
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(tables), 1, "Expected at least 1 table with LIMIT")
		})
	})

	t.Run("drop operations", func(t *testing.T) {
		t.Run("drop basic", func(t *testing.T) {
			id := createHybridTableBasic(t)

			// Drop the table (cleanup disabled for this test)
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id))
			require.NoError(t, err)

			// Verify it's gone
			_, err = client.HybridTables.ShowByID(ctx, id)
			require.Error(t, err)
		})

		t.Run("drop non-existent with IF EXISTS", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

			// Should not error with IF EXISTS
			err := client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
			require.NoError(t, err)
		})
	})
}
