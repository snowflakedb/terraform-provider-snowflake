//go:build !account_level_tests

package testint

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedColumn struct {
	Name string
	Type sdk.DataType
}

func TestInt_Table(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupTableProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(id))
			require.NoError(t, err)
		}
	}
	tag1, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	tag2, tagCleanup2 := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup2)

	assertColumns := func(t *testing.T, expectedColumns []expectedColumn, createdColumns []helpers.InformationSchemaColumns) {
		t.Helper()

		require.Len(t, createdColumns, len(expectedColumns))
		for i, expectedColumn := range expectedColumns {
			assert.Equal(t, strings.ToUpper(expectedColumn.Name), createdColumns[i].ColumnName)
			createdColumnDataType, err := datatypes.ParseDataType(createdColumns[i].DataType)
			assert.NoError(t, err)
			assert.Equal(t, expectedColumn.Type, sdk.LegacyDataTypeFrom(createdColumnDataType))
		}
	}

	assertTable := func(t *testing.T, table *sdk.Table, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, table.ID())
		assert.NotEmpty(t, table.CreatedOn)
		assert.Equal(t, id.Name(), table.Name)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), table.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), table.SchemaName)
		assert.Equal(t, "TABLE", table.Kind)
		assert.Equal(t, 0, table.Rows)
		assert.Equal(t, "ACCOUNTADMIN", table.Owner)
		assert.Equal(t, "ROLE", table.OwnerRoleType)
	}

	assertTableTerse := func(t *testing.T, table *sdk.Table, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, table.ID())
		assert.NotEmpty(t, table.CreatedOn)
		assert.Equal(t, id.Name(), table.Name)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), table.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), table.SchemaName)
		assert.Equal(t, "TABLE", table.Kind)
		assert.Empty(t, table.Rows)
		assert.Empty(t, table.Owner)
	}

	t.Run("create table: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("FIRST_COLUMN", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
			*sdk.NewTableColumnRequest("SECOND_COLUMN", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
	})

	t.Run("create table: complete optionals", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)
		table2, table2Cleanup := testClientHelper().Table.Create(t)
		t.Cleanup(table2Cleanup)
		comment := random.Comment()

		columnTags := []sdk.TagAssociation{
			{
				Name:  tag1.ID(),
				Value: "v1",
			},
			{
				Name:  tag2.ID(),
				Value: "v2",
			},
		}
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).
				WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithExpression(sdk.String("'default'"))).
				WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(maskingPolicy.ID()).WithUsing([]string{"COLUMN_1", "COLUMN_3"})).
				WithTags(columnTags).
				WithNotNull(sdk.Bool(true)),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypeForeignKey).
			WithName(sdk.String("OUT_OF_LINE_CONSTRAINT")).
			WithColumns([]string{"COLUMN_1"}).
			WithForeignKey(sdk.NewOutOfLineForeignKeyRequest(table2.ID(), []string{"id"}).
				WithMatch(sdk.Pointer(sdk.FullMatchType)).
				WithOn(sdk.NewForeignKeyOnAction().
					WithOnDelete(sdk.Pointer(sdk.ForeignKeySetNullAction)).WithOnUpdate(sdk.Pointer(sdk.ForeignKeyRestrictAction))))
		stageFileFormat := sdk.NewStageFileFormatRequest().
			WithType(sdk.Pointer(sdk.FileFormatTypeCSV)).
			WithOptions(sdk.NewFileFormatTypeOptionsRequest().WithCSVCompression(sdk.Pointer(sdk.CSVCompressionAuto)))
		stageCopyOptions := sdk.NewStageCopyOptionsRequest().WithOnError(sdk.NewStageCopyOnErrorOptionsRequest().WithSkipFile())
		request := sdk.NewCreateTableRequest(id, columns).
			WithOutOfLineConstraint(*outOfLineConstraint).
			WithStageFileFormat(*stageFileFormat).
			WithStageCopyOptions(*stageCopyOptions).
			WithComment(&comment).
			WithDataRetentionTimeInDays(sdk.Int(30)).
			WithMaxDataExtensionTimeInDays(sdk.Int(30))

		err := client.Tables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
		assert.Equal(t, comment, table.Comment)
		assert.Equal(t, 30, table.RetentionTime)

		param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterMaxDataExtensionTimeInDays, sdk.Object{ObjectType: sdk.ObjectTypeTable, Name: table.ID()})
		assert.NoError(t, err)
		assert.Equal(t, "30", param.Value)

		tableColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeNumber},
		}
		assertColumns(t, expectedColumns, tableColumns)
	})

	t.Run("create table as select", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
		t.Cleanup(maskingPolicyCleanup)
		columns := []sdk.TableAsSelectColumnRequest{
			*sdk.NewTableAsSelectColumnRequest("COLUMN_3").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_1").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_2").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)).WithMaskingPolicyName(sdk.Pointer(maskingPolicy.ID())),
		}

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := "SELECT 1, 2, 3"
		request := sdk.NewCreateTableAsSelectRequest(id, columns, query)

		err := client.Tables.CreateAsSelect(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		tableColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, tableColumns)
	})

	// TODO [SNOW-1007542]: fix this test, it should create two integer column but is creating 3 text ones instead
	t.Run("create table using template", func(t *testing.T) {
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateFileFormat(t)
		t.Cleanup(fileFormatCleanup)
		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		filePath := testhelpers.TestFile(t, "data.csv", []byte(` [{"name": "column1", "type" "INTEGER"},
									 {"name": "column2", "type" "INTEGER"} ]`))

		_, err := client.ExecForTests(ctx, fmt.Sprintf("PUT file://%s @%s", filePath, stage.ID().FullyQualifiedName()))
		require.NoError(t, err)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '@%s', FILE_FORMAT=>'%s', ignore_case => true))`, stage.ID().FullyQualifiedName(), fileFormat.ID().FullyQualifiedName())
		request := sdk.NewCreateTableUsingTemplateRequest(id, query)

		err = client.Tables.CreateUsingTemplate(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		returnedTableColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"C1", sdk.DataTypeVARCHAR},
			{"C2", sdk.DataTypeVARCHAR},
			{"C3", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, returnedTableColumns)
	})

	t.Run("create table like", func(t *testing.T) {
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", "NUMBER"),
			*sdk.NewTableColumnRequest("col2", "VARCHAR"),
			*sdk.NewTableColumnRequest("col3", "BOOLEAN"),
		}
		sourceTable, sourceTableCleanup := testClientHelper().Table.CreateWithColumns(t, columns)
		t.Cleanup(sourceTableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateTableLikeRequest(id, sourceTable.ID()).WithCopyGrants(sdk.Bool(true))

		err := client.Tables.CreateLike(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		sourceTableColumns := testClientHelper().Table.GetTableColumnsFor(t, sourceTable.ID())
		expectedColumns := []expectedColumn{
			{"id", sdk.DataTypeNumber},
			{"col2", sdk.DataTypeVARCHAR},
			{"col3", sdk.DataTypeBoolean},
		}
		assertColumns(t, expectedColumns, sourceTableColumns)

		likeTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		likeTableColumns := testClientHelper().Table.GetTableColumnsFor(t, likeTable.ID())
		assertColumns(t, expectedColumns, likeTableColumns)
	})

	t.Run("create table clone", func(t *testing.T) {
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", "NUMBER"),
			*sdk.NewTableColumnRequest("col2", "VARCHAR"),
			*sdk.NewTableColumnRequest("col3", "BOOLEAN"),
		}
		sourceTable, sourceTableCleanup := testClientHelper().Table.CreateWithColumns(t, columns)
		t.Cleanup(sourceTableCleanup)

		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		// ensure that time travel is allowed (and revert if needed after the test)
		testClientHelper().Schema.UpdateDataRetentionTime(t, schema.ID(), 1)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateTableCloneRequest(id, sourceTable.ID()).
			WithCopyGrants(sdk.Bool(true)).WithClonePoint(sdk.NewClonePointRequest().
			WithAt(*sdk.NewTimeTravelRequest().WithOffset(sdk.Pointer(0))).
			WithMoment(sdk.CloneMomentAt))

		err := client.Tables.CreateClone(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		sourceTableColumns := testClientHelper().Table.GetTableColumnsFor(t, sourceTable.ID())
		expectedColumns := []expectedColumn{
			{"id", sdk.DataTypeNumber},
			{"col2", sdk.DataTypeVARCHAR},
			{"col3", sdk.DataTypeBoolean},
		}
		assertColumns(t, expectedColumns, sourceTableColumns)

		cloneTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		cloneTableColumns := testClientHelper().Table.GetTableColumnsFor(t, cloneTable.ID())
		assertColumns(t, expectedColumns, cloneTableColumns)
	})

	t.Run("alter table: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)

		alterRequest := sdk.NewAlterTableRequest(id).WithNewName(&newId)
		err = client.Tables.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupTableProvider(id))
		} else {
			t.Cleanup(cleanupTableProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		table, err := client.Tables.ShowByID(ctx, newId)
		require.NoError(t, err)
		assertTable(t, table, newId)
	})

	t.Run("alter table: swap with", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		secondTableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secondTableColumns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(secondTableId))

		alterRequest := sdk.NewAlterTableRequest(id).WithSwapWith(&secondTableId)
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, secondTableId)
		require.NoError(t, err)

		assertTable(t, table, secondTableId)

		secondTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, secondTable, id)
	})

	t.Run("alter table: cluster by", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		clusterByColumns := []string{"COLUMN_1", "COLUMN_2"}
		alterRequest := sdk.NewAlterTableRequest(id).WithClusteringAction(sdk.NewTableClusteringActionRequest().WithClusterBy(clusterByColumns))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
		assert.Equal(t, "", table.Comment)
		clusterByString := "LINEAR(" + strings.Join(clusterByColumns, ", ") + ")"
		assert.Equal(t, clusterByString, table.ClusterBy)
	})

	t.Run("alter table: resume recluster", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(sdk.NewTableClusteringActionRequest().
				WithChangeReclusterState(sdk.Pointer(sdk.ReclusterStateSuspend)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		clusterByString := "LINEAR(" + strings.Join(clusterBy, ", ") + ")"
		assert.Equal(t, clusterByString, table.ClusterBy)
	})

	t.Run("alter table: drop clustering key", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(sdk.NewTableClusteringActionRequest().
				WithDropClusteringKey(sdk.Bool(true)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", table.ClusterBy)
	})

	t.Run("alter table: add a column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().
				WithAdd(sdk.NewTableColumnAddActionRequest("COLUMN_3", sdk.DataTypeVARCHAR).WithComment(sdk.String("some comment"))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
			{"COLUMN_3", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, "", table.Comment)
	})

	t.Run("alter table: rename column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().
				WithRename(sdk.NewTableColumnRenameActionRequest("COLUMN_1", "COLUMN_3")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, "", table.Comment)
	})

	t.Run("alter table: unset masking policy", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
		t.Cleanup(maskingPolicyCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(maskingPolicy.ID())),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableDetails, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
		require.NoError(t, err)

		require.Len(t, tableDetails, 2)
		// TODO [SNOW-1348114]: make nicer during the table rework
		assert.Equal(t, maskingPolicy.ID().FullyQualifiedName(), sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(*tableDetails[0].PolicyName).FullyQualifiedName())

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().WithUnsetMaskingPolicy(sdk.NewTableColumnAlterUnsetMaskingPolicyActionRequest("COLUMN_1")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		tableDetails, err = client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
		require.NoError(t, err)

		require.Len(t, tableDetails, 2)
		assert.Empty(t, tableDetails[0].PolicyName)
	})

	t.Run("alter table: drop columns", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().WithDropColumns([]string{"COLUMN_1"}))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, "", table.Comment)
	})

	// TODO [SNOW-1007542]: check added constraints
	// Add method similar to getTableColumnsFor based on https://docs.snowflake.com/en/sql-reference/info-schema/table_constraints.
	t.Run("alter constraint: add", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		secondTableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secondTableColumns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(secondTableId))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().
				WithAdd(sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypeForeignKey).WithName(sdk.String("OUT_OF_LINE_CONSTRAINT")).WithColumns([]string{"COLUMN_1"}).
					WithForeignKey(sdk.NewOutOfLineForeignKeyRequest(secondTableId, []string{"COLUMN_3"}))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("add constraint: not null", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{
				*sdk.NewTableColumnAlterActionRequest("COLUMN_1").
					WithNotNullConstraint(sdk.NewTableColumnNotNullConstraintRequest().WithSet(sdk.Bool(true))),
			}))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO [SNOW-1007542]: check renamed constraint
	t.Run("alter constraint: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		oldConstraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithName(sdk.String(oldConstraintName)).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		newConstraintName := "NEW_OUT_OF_LINE_CONSTRAINT_NAME"
		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().
				WithRename(sdk.NewTableConstraintRenameActionRequest().
					WithOldName(oldConstraintName).
					WithNewName(newConstraintName)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO [SNOW-1007542]: check altered constraint
	t.Run("alter constraint: alter", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithName(sdk.String(constraintName)).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().WithAlter(sdk.NewTableConstraintAlterActionRequest().WithConstraintName(sdk.String(constraintName)).WithEnforced(sdk.Bool(true))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO [SNOW-1007542]: check dropped constraint
	t.Run("alter constraint: drop constraint with name", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithName(sdk.String(constraintName)).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().WithDrop(sdk.NewTableConstraintDropActionRequest().WithConstraintName(sdk.String(constraintName))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter constraint: drop primary key without constraint name", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().WithDrop(sdk.NewTableConstraintDropActionRequest().WithPrimaryKey(sdk.Bool(true))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("external table: add column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithAdd(sdk.NewTableExternalTableColumnAddActionRequest().
				WithName("COLUMN_3").
				WithType(sdk.DataTypeNumber).
				WithExpression("1 + 1").
				WithComment(sdk.String("some comment")),
			))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
			{"COLUMN_3", sdk.DataTypeNumber},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	t.Run("external table: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithRename(sdk.NewTableExternalTableColumnRenameActionRequest().WithOldName("COLUMN_1").WithNewName("COLUMN_3")))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", table.Comment)
		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	t.Run("external table: drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithDrop(sdk.NewTableExternalTableColumnDropActionRequest([]string{"COLUMN_2"})))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	// TODO [SNOW-1007542]: check search optimization - after adding https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization
	t.Run("add search optimization", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithSearchOptimizationAction(sdk.NewTableSearchOptimizationActionRequest().WithAddSearchOptimizationOn([]string{"SUBSTRING(*)", "GEO(*)"}))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO [SNOW-1007542]: try to check more sets (ddl collation, max data extension time in days, etc.)
	t.Run("set: with complete options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		stageFileFormats := sdk.StageFileFormatRequest{
			Type: sdk.Pointer(sdk.FileFormatTypeCSV),
		}
		stageCopyOptions := sdk.StageCopyOptionsRequest{
			OnError: sdk.NewStageCopyOnErrorOptionsRequest().WithSkipFile(),
		}
		alterRequest := sdk.NewAlterTableRequest(id).
			WithSet(sdk.NewTableSetRequest().
				WithEnableSchemaEvolution(sdk.Bool(true)).
				WithStageFileFormat(stageFileFormats).
				WithStageCopyOptions(stageCopyOptions).
				WithDataRetentionTimeInDays(sdk.Int(30)).
				WithMaxDataExtensionTimeInDays(sdk.Int(90)).
				WithChangeTracking(sdk.Bool(false)).
				WithDefaultDDLCollation(sdk.String("us")).
				WithComment(&comment))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, comment, table.Comment)
		assert.Equal(t, 30, table.RetentionTime)
		assert.False(t, table.ChangeTracking)
		assert.True(t, table.EnableSchemaEvolution)
	})

	t.Run("drop table", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(table.ID()).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, table.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop table in non-existing schema", func(t *testing.T) {
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(nonExistingTableId).WithIfExists(sdk.Bool(true)))
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop table in non-existing database", func(t *testing.T) {
		nonExistingDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(nonExistingDatabaseId)
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(nonExistingTableId).WithIfExists(sdk.Bool(true)))
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop safely table", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		err := client.Tables.DropSafely(ctx, table.ID())
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, table.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop safely table in non-existing schema", func(t *testing.T) {
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.DropSafely(ctx, nonExistingTableId)
		require.NoError(t, err)
	})

	t.Run("drop safely table in non-existing database", func(t *testing.T) {
		nonExistingDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(nonExistingDatabaseId)
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.DropSafely(ctx, nonExistingTableId)
		require.NoError(t, err)
	})

	t.Run("show tables", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		table2, table2Cleanup := testClientHelper().Table.Create(t)
		t.Cleanup(table2Cleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest())
		require.NoError(t, err)

		t1, err := collections.FindFirst(tables, func(t sdk.Table) bool { return t.ID().FullyQualifiedName() == table.ID().FullyQualifiedName() })
		require.NoError(t, err)
		t2, err := collections.FindFirst(tables, func(t sdk.Table) bool { return t.ID().FullyQualifiedName() == table2.ID().FullyQualifiedName() })
		require.NoError(t, err)

		assertTable(t, t1, table.ID())
		assertTable(t, t2, table2.ID())
	})

	t.Run("with terse", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithTerse(true).WithLike(sdk.Like{
			Pattern: sdk.String(table.Name),
		}))
		require.NoError(t, err)
		assert.Len(t, tables, 1)

		assertTableTerse(t, &tables[0], table.ID())
	})

	t.Run("with starts with", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithStartsWith(table.Name))
		require.NoError(t, err)
		assert.Len(t, tables, 1)

		assertTable(t, &tables[0], table.ID())
	})

	t.Run("when searching a non-existent table", func(t *testing.T) {
		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithLike(sdk.Like{
			Pattern: sdk.String("non-existent"),
		}))
		require.NoError(t, err)
		assert.Empty(t, tables)
	})
}

func TestInt_TablesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupTableHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("c1", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createTableHandle(t, id1)
		createTableHandle(t, id2)

		e1, err := client.Tables.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Tables.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: check schema evolution record", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("c1", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithEnableSchemaEvolution(sdk.Pointer(true)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.Grants.GrantPrivilegesToAccountRole(ctx,
			&sdk.AccountRoleGrantPrivileges{SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeEvolveSchema}},
			&sdk.AccountRoleGrantOn{SchemaObject: &sdk.GrantOnSchemaObject{SchemaObject: &sdk.Object{ObjectType: sdk.ObjectTypeTable, Name: sdk.NewObjectIdentifierFromFullyQualifiedName(table.ID().FullyQualifiedName())}}},
			snowflakeroles.Accountadmin,
			nil)
		require.NoError(t, err)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		testClientHelper().Stage.PutOnStage(t, stage.ID(), "schema_evolution_record.json")

		testClientHelper().Stage.CopyIntoTableFromFile(t, table.ID(), stage.ID(), "schema_evolution_record.json")

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		require.Len(t, currentColumns, 2)
		assert.NotEmpty(t, currentColumns[1].SchemaEvolutionRecord)

		descColumns, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
		require.NoError(t, err)
		require.Len(t, descColumns, 2)
		assert.NotEmpty(t, descColumns[1].SchemaEvolutionRecord)
	})

	t.Run("show by id: missing database", func(t *testing.T) {
		databaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
		_, err := client.Tables.ShowByID(ctx, tableId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("show by id: missing schema", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
		_, err := client.Tables.ShowByID(ctx, tableId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("show by id safely", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createTableHandle(t, id)
		table, err := client.Tables.ShowByIDSafely(ctx, id)
		assert.NotNil(t, table)
		assert.NoError(t, err)
	})

	t.Run("show by id safely: missing database", func(t *testing.T) {
		databaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
		_, err := client.Tables.ShowByIDSafely(ctx, tableId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("show by id safely: missing schema", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
		_, err := client.Tables.ShowByIDSafely(ctx, tableId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("show by id safely: missing table", func(t *testing.T) {
		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, err := client.Tables.ShowByIDSafely(ctx, tableId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})
}
