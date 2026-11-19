//go:build non_account_level_tests

package testint

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
			assert.Equal(t, expectedColumn.Name, createdColumns[i].ColumnName)
			createdColumnDataType, err := datatypes.ParseDataType(createdColumns[i].DataType)
			assert.NoError(t, err)
			assert.Equal(t, expectedColumn.Type, sdk.LegacyDataTypeFrom(createdColumnDataType))
		}
	}

	assertTable := func(t *testing.T, table *sdk.Table, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, table.ID())
		assertThatObject(t, objectassert.TableFromObject(t, table).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasKind("TABLE").
			HasRows(0).
			HasOwner("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE").
			HasIsEvent(false).
			HasIsHybrid(false).
			HasIsIceberg(false).
			HasIsDynamic(false).
			HasIsImmutable(false).
			HasIsInteractive(false).
			HasRowTimestamp(false).
			HasErrorLogging(false))
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

	selectTableConstraints := func(t *testing.T, id sdk.SchemaObjectIdentifier) []sdk.TableConstraintDetails {
		t.Helper()
		constraints, err := client.Tables.SelectTableConstraints(ctx, sdk.NewSelectTableConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		slices.SortFunc(constraints, func(x, y sdk.TableConstraintDetails) int {
			return strings.Compare(x.ConstraintName, y.ConstraintName)
		})
		return constraints
	}

	addConstraint := func(t *testing.T, id sdk.SchemaObjectIdentifier, constraint sdk.TableOutOfLineConstraintRequest) {
		t.Helper()
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(*sdk.NewTableConstraintActionRequest().WithAdd(constraint)))
		require.NoError(t, err)
	}

	createTableWithColumns := func(t *testing.T, columnNames ...string) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := make([]sdk.TableColumnRequest, len(columnNames))
		for i, name := range columnNames {
			columns[i] = *sdk.NewTableColumnRequest(name, sdk.DataTypeNumber)
		}
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))
		return id
	}

	t.Run("create table: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("FIRST_COLUMN", sdk.DataTypeNumber).WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))),
			*sdk.NewTableColumnRequest("SECOND_COLUMN", sdk.DataTypeNumber).WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
	})

	t.Run("show: dynamic table is_dynamic", func(t *testing.T) {
		baseID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(baseID, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber),
		})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(baseID))

		dynamicID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(
			dynamicID,
			sdk.TargetLagRequest{MaximumDuration: sdk.String("1 hour")},
			testClientHelper().Ids.WarehouseId(),
			fmt.Sprintf("SELECT ID FROM %s", baseID.FullyQualifiedName()),
		))
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(dynamicID)))
		})

		table, err := client.Tables.ShowByID(ctx, dynamicID)
		require.NoError(t, err)
		assertThatObject(t, objectassert.TableFromObject(t, table).
			HasKind("TABLE").
			HasIsDynamic(true).
			HasIsHybrid(false).
			HasIsIceberg(false).
			HasIsEvent(false).
			HasIsImmutable(false).
			HasIsInteractive(false))
	})

	t.Run("create table: DECFLOAT", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("FIRST_COLUMN", sdk.DataTypeNumber).WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))),
			*sdk.NewTableColumnRequest("SECOND_COLUMN", datatypes.DecfloatLegacyDataType),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)

		returnedTableColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"FIRST_COLUMN", sdk.DataTypeNumber},
			{"SECOND_COLUMN", datatypes.DecfloatLegacyDataType},
		}
		assertColumns(t, expectedColumns, returnedTableColumns)
	})

	t.Run("create table: complete optionals", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)
		referencedId := createTableWithColumns(t, "id")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "id"}})))
		contactId, contactCleanup := testClientHelper().Contact.Create(t)
		t.Cleanup(contactCleanup)
		storageLifecyclePolicyId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().StorageLifecyclePolicy.CreateWithId(t, storageLifecyclePolicyId))
		aggregationPolicy, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)
		joinPolicy, joinPolicyCleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(joinPolicyCleanup)
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
			*sdk.NewTableColumnRequest("column_3", sdk.DataTypeVARCHAR).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("uq_Inline").
					WithUnique(true))),
			*sdk.NewTableColumnRequest("column_1", sdk.DataTypeVARCHAR).
				WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression("'default'")).
				WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(maskingPolicy.ID()).WithUsing([]sdk.Column{{Value: "column_1"}, {Value: "column_3"}})).
				WithTag(columnTags).
				WithNotNull(true),
			*sdk.NewTableColumnRequest("column_2", sdk.DataTypeNumber).
				WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("pk_Id").
					WithPrimaryKey(true).
					WithNotEnforced(true).
					WithRely(true))),
			*sdk.NewTableColumnRequest("ref_id", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("inline_ref", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithFK(*sdk.NewTableColumnInlineFKRequest(referencedId).
					WithName("fk_Inline").
					WithForeignKey(true).
					WithRefColumn([]sdk.Column{{Value: "id"}}).
					WithMatch(sdk.MatchTypeFull).
					WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction))))),
			*sdk.NewTableColumnRequest("amount", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithCH(*sdk.NewTableColumnInlineCHRequest(`"amount" > 0`).
					WithName("chk_Amount_Positive"))),
			*sdk.NewTableColumnRequest("DOUBLE_ID", sdk.DataTypeNumber).
				WithVirtual(*sdk.NewTableColumnVirtualRequest("(1 + 1)").WithGeneratedAlways(true).WithVirtual(true)),
		}
		constraints := []sdk.TableOutOfLineConstraintRequest{
			*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
				WithName("uq_Code").
				WithUnique(true).
				WithColumns([]sdk.Column{{Value: "column_1"}}).
				WithComment("unique comment")),
			*sdk.NewTableOutOfLineConstraintRequest().WithFK(*sdk.NewTableOutOfLineFKRequest(referencedId).
				WithName("fk_Ref").
				WithColumns([]sdk.Column{{Value: "ref_id"}}).
				WithRefColumns([]sdk.Column{{Value: "id"}}).
				WithMatch(sdk.MatchTypeFull).
				WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction)))),
			*sdk.NewTableOutOfLineConstraintRequest().WithCH(*sdk.NewTableOutOfLineCHRequest(`"column_2" > 0`).
				WithName("chk_Id_Positive")),
		}
		request := sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns).WithOutOfLineConstraint(constraints)).
			WithComment(comment).
			WithDataRetentionTimeInDays(30).
			WithMaxDataExtensionTimeInDays(30).
			WithErrorLogging(true).
			WithRowTimestamp(true).
			WithContact([]sdk.TableContact{{Purpose: "SUPPORT", Contact: contactId}}).
			WithStorageLifecyclePolicy(*sdk.NewCreateTableStorageLifecyclePolicyRequest(storageLifecyclePolicyId, []sdk.Column{{Value: "column_1"}})).
			WithAggregationPolicy(*sdk.NewTableAggregationPolicyRequest(aggregationPolicy).WithEntityKey([]sdk.Column{{Value: "column_1"}})).
			WithJoinPolicy(*sdk.NewTableJoinPolicyRequest(joinPolicy).WithAllowedJoinKeys([]sdk.Column{{Value: "column_1"}})).
			WithTag([]sdk.TagAssociation{{Name: tag1.ID(), Value: "v1"}})

		err := client.Tables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.TableFromObject(t, table).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasKind("TABLE").
			HasComment(comment).
			HasRetentionTime(30).
			HasOwner("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE").
			HasErrorLogging(true).
			HasRowTimestamp(true))

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assertThatObject(t, objectassert.PolicyReferenceFromObject(t, new(references[0])).HasRefEntityName(id.Name()))

		tableColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"column_3", sdk.DataTypeVARCHAR},
			{"column_1", sdk.DataTypeVARCHAR},
			{"column_2", sdk.DataTypeNumber},
			{"ref_id", sdk.DataTypeNumber},
			{"inline_ref", sdk.DataTypeNumber},
			{"amount", sdk.DataTypeNumber},
			{"DOUBLE_ID", sdk.DataTypeNumber},
		}
		assertColumns(t, expectedColumns, tableColumns)

		tableConstraints := selectTableConstraints(t, id)
		require.Len(t, tableConstraints, 5)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[0]).
			HasConstraintName("fk_Inline").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[1]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[2]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[3]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasComment("unique comment").
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[4]).
			HasConstraintName("uq_Inline").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		slices.SortFunc(checkConstraints, func(x, y sdk.TableCheckConstraintDetails) int {
			return strings.Compare(x.ConstraintName, y.ConstraintName)
		})
		require.Len(t, checkConstraints, 2)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Amount_Positive").
			HasCheckClause(`"amount" > 0`).
			HasConstraintTable(id.Name()))
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[1]).
			HasConstraintName("chk_Id_Positive").
			HasCheckClause(`"column_2" > 0`).
			HasConstraintTable(id.Name()))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		})).WithOrReplace(true).WithCopy([]sdk.TableCopyOption{sdk.TableCopyOptionComments, sdk.TableCopyOptionTags}))
		require.NoError(t, err)

		table, err = client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.TableFromObject(t, table).HasComment(comment))
	})

	t.Run("create table: out-of-line primary key, unique, foreign key, and check", func(t *testing.T) {
		referencedId := createTableWithColumns(t, "REF_ID")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "REF_ID"}})))

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("CODE", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("REF_ID", sdk.DataTypeNumber),
		}
		constraints := []sdk.TableOutOfLineConstraintRequest{
			*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
				WithName("pk_Id").
				WithPrimaryKey(true).
				WithColumns([]sdk.Column{{Value: "ID"}}).
				WithNotEnforced(true).
				WithRely(true).
				WithComment("primary key comment")),
			*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
				WithName("uq_Code").
				WithUnique(true).
				WithColumns([]sdk.Column{{Value: "CODE"}})),
			*sdk.NewTableOutOfLineConstraintRequest().WithFK(*sdk.NewTableOutOfLineFKRequest(referencedId).
				WithName("fk_Ref").
				WithColumns([]sdk.Column{{Value: "REF_ID"}}).
				WithRefColumns([]sdk.Column{{Value: "REF_ID"}}).
				WithMatch(sdk.MatchTypeFull).
				WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction)))),
			*sdk.NewTableOutOfLineConstraintRequest().WithCH(*sdk.NewTableOutOfLineCHRequest("ID > 0").
				WithName("chk_Id_Positive")),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns).WithOutOfLineConstraint(constraints)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableConstraints := selectTableConstraints(t, id)
		require.Len(t, tableConstraints, 3)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[0]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[1]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasComment("primary key comment").
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[2]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, checkConstraints, 1)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Id_Positive").
			HasCheckClause("ID > 0").
			HasConstraintTable(id.Name()))
	})

	t.Run("create table: inline primary key, unique, foreign key, and check", func(t *testing.T) {
		referencedId := createTableWithColumns(t, "REF_ID")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "REF_ID"}})))

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("pk_Id").
					WithPrimaryKey(true).
					WithNotEnforced(true).
					WithRely(true))),
			*sdk.NewTableColumnRequest("CODE", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("uq_Code").
					WithUnique(true))),
			*sdk.NewTableColumnRequest("REF_ID", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithFK(*sdk.NewTableColumnInlineFKRequest(referencedId).
					WithName("fk_Ref").
					WithForeignKey(true).
					WithRefColumn([]sdk.Column{{Value: "REF_ID"}}).
					WithMatch(sdk.MatchTypeFull).
					WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction))))),
			*sdk.NewTableColumnRequest("AMOUNT", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithCH(*sdk.NewTableColumnInlineCHRequest("AMOUNT > 0").
					WithName("chk_Amount_Positive"))),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableConstraints := selectTableConstraints(t, id)
		require.Len(t, tableConstraints, 3)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[0]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[1]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[2]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, checkConstraints, 1)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Amount_Positive").
			HasCheckClause("AMOUNT > 0").
			HasConstraintTable(id.Name()))
	})

	t.Run("create table: out-of-line primary key, unique, foreign key, and check", func(t *testing.T) {
		referencedId := createTableWithColumns(t, "REF_ID")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "REF_ID"}})))

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("CODE", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("REF_ID", sdk.DataTypeNumber),
		}
		constraints := []sdk.TableOutOfLineConstraintRequest{
			*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
				WithName("pk_Id").
				WithPrimaryKey(true).
				WithColumns([]sdk.Column{{Value: "ID"}}).
				WithNotEnforced(true).
				WithRely(true).
				WithComment("primary key comment")),
			*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
				WithName("uq_Code").
				WithUnique(true).
				WithColumns([]sdk.Column{{Value: "CODE"}})),
			*sdk.NewTableOutOfLineConstraintRequest().WithFK(*sdk.NewTableOutOfLineFKRequest(referencedId).
				WithName("fk_Ref").
				WithColumns([]sdk.Column{{Value: "REF_ID"}}).
				WithRefColumns([]sdk.Column{{Value: "REF_ID"}}).
				WithMatch(sdk.MatchTypeFull).
				WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction)))),
			*sdk.NewTableOutOfLineConstraintRequest().WithCH(*sdk.NewTableOutOfLineCHRequest("ID > 0").
				WithName("chk_Id_Positive")),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns).WithOutOfLineConstraint(constraints)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableConstraints := selectTableConstraints(t, id)
		require.Len(t, tableConstraints, 3)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[0]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[1]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasComment("primary key comment").
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[2]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, checkConstraints, 1)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Id_Positive").
			HasCheckClause("ID > 0").
			HasConstraintTable(id.Name()))
	})

	t.Run("create table: inline primary key, unique, foreign key, and check", func(t *testing.T) {
		referencedId := createTableWithColumns(t, "REF_ID")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "REF_ID"}})))

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("pk_Id").
					WithPrimaryKey(true).
					WithNotEnforced(true).
					WithRely(true))),
			*sdk.NewTableColumnRequest("CODE", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithUniquePK(*sdk.NewTableColumnInlineUniquePKRequest().
					WithName("uq_Code").
					WithUnique(true))),
			*sdk.NewTableColumnRequest("REF_ID", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithFK(*sdk.NewTableColumnInlineFKRequest(referencedId).
					WithName("fk_Ref").
					WithForeignKey(true).
					WithRefColumn([]sdk.Column{{Value: "REF_ID"}}).
					WithMatch(sdk.MatchTypeFull).
					WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction))))),
			*sdk.NewTableColumnRequest("AMOUNT", sdk.DataTypeNumber).
				WithInlineConstraint(*sdk.NewTableColumnInlineConstraintRequest().WithCH(*sdk.NewTableColumnInlineCHRequest("AMOUNT > 0").
					WithName("chk_Amount_Positive"))),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableConstraints := selectTableConstraints(t, id)
		require.Len(t, tableConstraints, 3)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[0]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[1]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &tableConstraints[2]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, checkConstraints, 1)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Amount_Positive").
			HasCheckClause("AMOUNT > 0").
			HasConstraintTable(id.Name()))
	})

	t.Run("create table as select", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
		t.Cleanup(maskingPolicyCleanup)
		columns := []sdk.TableAsSelectColumnRequest{
			*sdk.NewTableAsSelectColumnRequest("COLUMN_3").
				WithColumnType(sdk.DataTypeVARCHAR),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_1").
				WithColumnType(sdk.DataTypeVARCHAR),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_2").
				WithColumnType(sdk.DataTypeVARCHAR).
				WithMaskingPolicy(maskingPolicy.ID()),
		}

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := "SELECT 1, 2, 3"
		request := sdk.NewCreateAsSelectTableRequest(id, columns, query)

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
		fileFormat, fileFormatCleanup := testClientHelper().FileFormat.CreateCsv(t)
		t.Cleanup(fileFormatCleanup)
		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		filePath := testfiles.TestFile(t, "data.csv", []byte(` [{"name": "column1", "type" "INTEGER"},
									 {"name": "column2", "type" "INTEGER"} ]`))

		_, err := client.ExecForTests(ctx, fmt.Sprintf("PUT file://%s @%s", filePath, stage.ID().FullyQualifiedName()))
		require.NoError(t, err)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '@%s', FILE_FORMAT=>'%s', ignore_case => true))`, stage.ID().FullyQualifiedName(), fileFormat.ID().FullyQualifiedName())
		request := sdk.NewCreateUsingTemplateTableRequest(id).WithQuery([]string{query})

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
		request := sdk.NewCreateLikeTableRequest(id, sourceTable.ID()).WithCopy([]sdk.TableCopyOption{sdk.TableCopyOptionGrants})

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
		request := sdk.NewCreateCloneTableRequest(id, sourceTable.ID()).
			WithCopy([]sdk.TableCopyOption{sdk.TableCopyOptionGrants}).WithClonePoint(*sdk.NewTableClonePointRequest(sdk.CloneMomentAt, sdk.TimeTravel{Offset: sdk.Pointer(0)}))

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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)

		alterRequest := sdk.NewAlterTableRequest(id).WithRenameTo(newId)
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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(secondTableId, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(secondTableColumns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(secondTableId))

		alterRequest := sdk.NewAlterTableRequest(id).WithSwapWith(secondTableId)
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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		clusterByColumns := []string{"COLUMN_1", "COLUMN_2"}
		alterRequest := sdk.NewAlterTableRequest(id).WithClusteringAction(*sdk.NewTableClusteringActionRequest().WithClusterBy(clusterByColumns))

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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(*sdk.NewTableClusteringActionRequest().
				WithChangeReclusterState(*sdk.NewTableReclusterChangeStateRequest().WithState(sdk.ReclusterStateSuspend)))
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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(*sdk.NewTableClusteringActionRequest().
				WithDropClusteringKey(true))
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

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().
				WithAdd(*sdk.NewTableColumnAddActionRequest("COLUMN_3", sdk.DataTypeVARCHAR).WithComment("some comment")))
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
		assertThatObject(t, objectassert.TableFromObject(t, table).HasComment(""))

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAdd(
			*sdk.NewTableColumnAddActionRequest("DOUBLE_ID", sdk.DataTypeNumber).
				WithVirtual(*sdk.NewTableColumnVirtualRequest("(1 + 1)").WithGeneratedAlways(true).WithVirtual(true)),
		)))
		require.NoError(t, err)
	})

	// https://github.com/snowflakedb/terraform-provider-snowflake/issues/4730
	t.Run("alter table: add a column with a constant default", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().
				WithAdd(*sdk.NewTableColumnAddActionRequest("COLUMN_2", sdk.DataTypeBoolean).
					WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression("FALSE"))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeBoolean},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	t.Run("alter table: rename column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().
				WithRename(*sdk.NewTableRenameColumnActionRequest("COLUMN_1", "COLUMN_3")))
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
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(maskingPolicy.ID())),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableDetails, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeColumnsTableRequest(id))
		require.NoError(t, err)

		require.Len(t, tableDetails, 2)
		// TODO [SNOW-1348114]: make nicer during the table rework
		assert.Equal(t, maskingPolicy.ID().FullyQualifiedName(), sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(*tableDetails[0].PolicyName).FullyQualifiedName())

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().WithUnsetMaskingPolicy(*sdk.NewTableUnsetColumnMaskingPolicyRequest("COLUMN_1")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		tableDetails, err = client.Tables.DescribeColumns(ctx, sdk.NewDescribeColumnsTableRequest(id))
		require.NoError(t, err)

		require.Len(t, tableDetails, 2)
		assert.Empty(t, tableDetails[0].PolicyName)
	})

	t.Run("create table with projection policy on column", func(t *testing.T) {
		projectionPolicyId, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).WithProjectionPolicy(*sdk.NewTableColumnProjectionPolicyRequest(projectionPolicyId)),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		objectassert.PolicyReferenceFromObject(t, &references[0]).
			HasPolicyName(projectionPolicyId.Name()).
			HasPolicyKind(sdk.PolicyKindProjectionPolicy).
			HasRefEntityName(id.Name()).
			HasRefEntityDomain(string(sdk.PolicyEntityDomainTable)).
			HasRefColumnName("COLUMN_1")
	})

	t.Run("alter table: add, set, and unset projection policy on column", func(t *testing.T) {
		projectionPolicyId, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)
		projectionPolicy2Id, projectionPolicy2Cleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicy2Cleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		addRequest := sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAdd(
			*sdk.NewTableColumnAddActionRequest("COLUMN_2", sdk.DataTypeVARCHAR).WithProjectionPolicy(*sdk.NewTableColumnProjectionPolicyRequest(projectionPolicyId)),
		))
		err = client.Tables.Alter(ctx, addRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		objectassert.PolicyReferenceFromObject(t, &references[0]).
			HasPolicyName(projectionPolicyId.Name()).
			HasPolicyKind(sdk.PolicyKindProjectionPolicy).
			HasRefColumnName("COLUMN_2")

		setRequest := sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithSetProjectionPolicy(
			*sdk.NewTableSetColumnProjectionPolicyRequest("COLUMN_1", projectionPolicy2Id),
		))
		err = client.Tables.Alter(ctx, setRequest)
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 2)

		unsetRequest := sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithUnsetProjectionPolicy(
			*sdk.NewTableUnsetColumnProjectionPolicyRequest("COLUMN_1"),
		))
		err = client.Tables.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		objectassert.PolicyReferenceFromObject(t, &references[0]).
			HasPolicyName(projectionPolicyId.Name()).
			HasPolicyKind(sdk.PolicyKindProjectionPolicy).
			HasRefColumnName("COLUMN_2")
	})

	t.Run("alter table: add and drop storage lifecycle policy", func(t *testing.T) {
		storageLifecyclePolicyId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		storageLifecyclePolicyCleanup := testClientHelper().StorageLifecyclePolicy.CreateWithId(t, storageLifecyclePolicyId)
		t.Cleanup(storageLifecyclePolicyCleanup)

		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(tableId, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(tableId))

		addRequest := sdk.NewAlterTableRequest(tableId).WithAddStorageLifecyclePolicy(
			*sdk.NewTableStorageLifecyclePolicyRequest(storageLifecyclePolicyId, []sdk.Column{{Value: "COLUMN_1"}}),
		)
		err = client.Tables.Alter(ctx, addRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertThatObject(
			t, objectassert.PolicyReferenceFromObject(t, new(references[0])).
				HasPolicyDb(storageLifecyclePolicyId.DatabaseName()).
				HasPolicySchema(storageLifecyclePolicyId.SchemaName()).
				HasPolicyName(storageLifecyclePolicyId.Name()).
				HasPolicyKind(sdk.PolicyKindStorageLifecyclePolicy).
				HasRefDatabaseName(tableId.DatabaseName()).
				HasRefSchemaName(tableId.SchemaName()).
				HasRefEntityName(tableId.Name()).
				HasRefEntityDomain(string(sdk.PolicyEntityDomainTable)).
				HasRefArgColumnNames(`[ "COLUMN_1" ]`).
				HasPolicyStatus("ACTIVE"),
		)

		dropRequest := sdk.NewAlterTableRequest(tableId).WithDropStorageLifecyclePolicy(true)
		err = client.Tables.Alter(ctx, dropRequest)
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	// TODO [SNOW-4532966]: add an integration test for CREATE TABLE ... FROM BACKUP SET ... IDENTIFIER once the variant is added back to the SSOT definition.

	t.Run("alter table: set and unset aggregation policy", func(t *testing.T) {
		aggregationPolicy, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)
		aggregationPolicy2, aggregationPolicy2Cleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicy2Cleanup)

		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(tableId, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(tableId))

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithSetAggregationPolicy(*sdk.NewTableSetAggregationPolicyRequest(aggregationPolicy).
				WithEntityKey([]sdk.Column{{Value: "COLUMN_1"}})))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, aggregationPolicy.Name(), references[0].PolicyName)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithSetAggregationPolicy(*sdk.NewTableSetAggregationPolicyRequest(aggregationPolicy2).
				WithEntityKey([]sdk.Column{{Value: "COLUMN_1"}}).
				WithForce(true)))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, aggregationPolicy2.Name(), references[0].PolicyName)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithUnsetAggregationPolicy(*sdk.NewTableUnsetAggregationPolicyRequest()))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter table: set and unset join policy", func(t *testing.T) {
		joinPolicy1, joinPolicy1Cleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(joinPolicy1Cleanup)
		joinPolicy2, joinPolicy2Cleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(joinPolicy2Cleanup)

		tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(tableId, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(tableId))

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithSetJoinPolicy(*sdk.NewTableSetJoinPolicyRequest(joinPolicy1)))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, joinPolicy1.Name(), references[0].PolicyName)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithSetJoinPolicy(*sdk.NewTableSetJoinPolicyRequest(joinPolicy2).WithForce(true)))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, joinPolicy2.Name(), references[0].PolicyName)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableId).
			WithUnsetJoinPolicy(*sdk.NewTableUnsetJoinPolicyRequest()))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, tableId, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter table: drop columns", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().WithDropColumns(*sdk.NewTableDropColumnActionRequest([]sdk.Column{{Value: "COLUMN_1"}})))
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

	t.Run("alter constraint: add primary key, unique, foreign key, and check", func(t *testing.T) {
		referencedId := createTableWithColumns(t, "REF_ID")
		addConstraint(t, referencedId, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "REF_ID"}})))

		id := createTableWithColumns(t, "ID", "CODE", "REF_ID")

		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithName("pk_Id").
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "ID"}}).
			WithNotEnforced(true).
			WithRely(true).
			WithComment("primary key comment")))
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithName("uq_Code").
			WithUnique(true).
			WithColumns([]sdk.Column{{Value: "CODE"}})))
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithFK(*sdk.NewTableOutOfLineFKRequest(referencedId).
			WithName("fk_Ref").
			WithColumns([]sdk.Column{{Value: "REF_ID"}}).
			WithRefColumns([]sdk.Column{{Value: "REF_ID"}}).
			WithMatch(sdk.MatchTypeFull).
			WithOn(*sdk.NewForeignKeyOnAction().WithOnUpdate(new(sdk.ForeignKeyNoAction)).WithOnDelete(new(sdk.ForeignKeyNoAction)))))
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithCH(*sdk.NewTableOutOfLineCHRequest("ID > 0").
			WithName("chk_Id_Positive").
			WithEnableNovalidate(true)))

		constraints := selectTableConstraints(t, id)
		require.Len(t, constraints, 3)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &constraints[0]).
			HasConstraintName("fk_Ref").
			HasConstraintType(sdk.TableConstraintTypeForeignKey).
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &constraints[1]).
			HasConstraintName("pk_Id").
			HasConstraintType(sdk.TableConstraintTypePrimaryKey).
			HasEnforced(false).
			HasRely(true).
			HasComment("primary key comment").
			HasTableName(id.Name()))
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &constraints[2]).
			HasConstraintName("uq_Code").
			HasConstraintType(sdk.TableConstraintTypeUnique).
			HasTableName(id.Name()))

		checkConstraints, err := client.Tables.SelectCheckConstraints(ctx, sdk.NewSelectCheckConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, checkConstraints, 1)
		assertThatObject(t, objectassert.TableCheckConstraintDetailsFromObject(t, &checkConstraints[0]).
			HasConstraintName("chk_Id_Positive").
			HasCheckClause("ID > 0").
			HasConstraintTable(id.Name()))
	})

	t.Run("add constraint: not null", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{
				*sdk.NewTableColumnAlterActionRequest("COLUMN_1").
					WithNotNullConstraint(*sdk.NewTableColumnNotNullConstraintRequest().WithSet(true)),
			}))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter constraint: rename", func(t *testing.T) {
		id := createTableWithColumns(t, "COLUMN_1")
		oldConstraintName := "Out_Of_Line_Constraint"
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithName(oldConstraintName).
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "COLUMN_1"}})))

		newConstraintName := "New_Out_Of_Line_Constraint"
		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(*sdk.NewTableConstraintActionRequest().
				WithRename(*sdk.NewTableConstraintRenameActionRequest(oldConstraintName, newConstraintName)))
		err := client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		constraints := selectTableConstraints(t, id)
		require.Len(t, constraints, 1)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &constraints[0]).
			HasConstraintName(newConstraintName).
			HasConstraintType(sdk.TableConstraintTypePrimaryKey))
	})

	t.Run("alter constraint: alter", func(t *testing.T) {
		id := createTableWithColumns(t, "COLUMN_1")
		constraintName := "Out_Of_Line_Constraint"
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithName(constraintName).
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "COLUMN_1"}})))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(*sdk.NewTableConstraintActionRequest().WithAlter(*sdk.NewTableConstraintAlterActionRequest().WithConstraintName(constraintName).WithRely(true)))
		err := client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		constraints := selectTableConstraints(t, id)
		require.Len(t, constraints, 1)
		assertThatObject(t, objectassert.TableConstraintDetailsFromObject(t, &constraints[0]).
			HasConstraintName(constraintName).
			HasRely(true))
	})

	t.Run("alter constraint: drop constraint with name", func(t *testing.T) {
		id := createTableWithColumns(t, "COLUMN_1")
		constraintName := "Out_Of_Line_Constraint"
		addConstraint(t, id, *sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithName(constraintName).
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "COLUMN_1"}})))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(*sdk.NewTableConstraintActionRequest().WithDrop(*sdk.NewTableConstraintDropActionRequest().WithConstraintName(constraintName)))
		err := client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		require.Empty(t, selectTableConstraints(t, id))
	})

	t.Run("alter constraint: drop primary key without constraint name", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}
		outOfLineConstraint := sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*sdk.NewTableOutOfLineUniquePKRequest().
			WithPrimaryKey(true).
			WithColumns([]sdk.Column{{Value: "COLUMN_1"}}))

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns).WithOutOfLineConstraint([]sdk.TableOutOfLineConstraintRequest{*outOfLineConstraint})))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(*sdk.NewTableConstraintActionRequest().WithDrop(*sdk.NewTableConstraintDropActionRequest().WithPrimaryKey(true)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("external table: add column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(*sdk.NewTableExternalTableActionRequest().WithAdd(
				*sdk.NewTableExternalTableColumnAddActionRequest("COLUMN_3", sdk.DataTypeNumber).
					WithExpression([]string{"1 + 1"}).
					WithComment("some comment"),
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

	// TODO [SNOW-1007542]: check search optimization - after adding https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization
	t.Run("add search optimization", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithSearchOptimizationAction(*sdk.NewTableSearchOptimizationActionRequest().WithAdd(*sdk.NewTableAddSearchOptimizationRequest().WithOn([]sdk.TableSearchMethodWithTargetRequest{
				*sdk.NewTableSearchMethodWithTargetRequest(sdk.TableSearchMethodSubstring,
					*sdk.NewTableSearchMethodArgsRequest().WithTargets([]string{"*"})),
			})))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter table: set and unset complete options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		contactId, contactCleanup := testClientHelper().Contact.Create(t)
		t.Cleanup(contactCleanup)
		dmfId, dmfCleanup := testClientHelper().DataMetricFunctionClient.CreateWithArguments(t, []string{"ARG_T TABLE(ARG_C NUMBER)"}, []string{"TABLE(NUMBER)"})
		t.Cleanup(dmfCleanup)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("ID", sdk.DataTypeNumber),
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).
			WithSet(*sdk.NewTableSetRequest().
				WithEnableSchemaEvolution(true).
				WithDataRetentionTimeInDays(30).
				WithMaxDataExtensionTimeInDays(90).
				WithChangeTracking(false).
				WithDefaultDdlCollation(sdk.StringAllowEmpty{Value: "us"}).
				WithIcebergDefaultDdlCollation(sdk.StringAllowEmpty{Value: "en-ci"}).
				WithErrorLogging(true).
				WithContact([]sdk.TableContact{{Purpose: "SUPPORT", Contact: contactId}}).
				WithComment(comment).
				WithRowTimestamp(true)))
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.TableFromObject(t, table).
			HasComment(comment).
			HasRetentionTime(30).
			HasChangeTracking(false).
			HasEnableSchemaEvolution(true).
			HasErrorLogging(true).
			HasRowTimestamp(true))

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnset(*sdk.NewTableUnsetRequest().
			WithDataRetentionTimeInDays(true).
			WithMaxDataExtensionTimeInDays(true).
			WithChangeTracking(true).
			WithDefaultDdlCollation(true).
			WithIcebergDefaultDdlCollation(true).
			WithEnableSchemaEvolution(true).
			WithErrorLogging(true).
			WithComment(true).
			WithRowTimestamp(true)))
		require.NoError(t, err)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnset(*sdk.NewTableUnsetRequest().WithContact("SUPPORT")))
		require.NoError(t, err)

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSetDataMetricSchedule(*sdk.NewTableSetDataMetricScheduleRequest("TRIGGER_ON_CHANGES")))
		require.NoError(t, err)
		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithAddDataMetricFunction(*sdk.NewTableAddDataMetricFunctionRequest([]sdk.TableAddedFunctionsRequest{
			*sdk.NewTableAddedFunctionsRequest(dmfId.SchemaObjectId(), []sdk.Column{{Value: "ID"}}),
		})))
		require.NoError(t, err)
		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithModifyDataMetricFunction(*sdk.NewTableModifyDataMetricFunctionRequest([]sdk.TableModifiedFunctionsRequest{
			*sdk.NewTableModifiedFunctionsRequest(dmfId.SchemaObjectId(), []sdk.Column{{Value: "ID"}}).WithStatus(sdk.TableDataMetricScheduleStatusOperationSuspend),
		})))
		require.NoError(t, err)
		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithDropDataMetricFunction(*sdk.NewTableDropDataMetricFunctionRequest([]sdk.TableDroppedFunctionsRequest{
			*sdk.NewTableDroppedFunctionsRequest(dmfId.SchemaObjectId(), []sdk.Column{{Value: "ID"}}),
		})))
		require.NoError(t, err)
		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnsetDataMetricSchedule(true))
		require.NoError(t, err)
	})

	t.Run("drop table", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(table.ID()).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, table.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop table in non-existing schema", func(t *testing.T) {
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(nonExistingTableId).WithIfExists(true))
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop table in non-existing database", func(t *testing.T) {
		nonExistingDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		nonExistingSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(nonExistingDatabaseId)
		nonExistingTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(nonExistingSchemaId)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(nonExistingTableId).WithIfExists(true))
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
			*sdk.NewTableColumnRequest("c1", sdk.DataTypeNumber).WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))),
		}
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)))
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
			*sdk.NewTableColumnRequest("c1", sdk.DataTypeNumber).WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(1, 1))),
		}
		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(columns)).WithEnableSchemaEvolution(true))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(snowflakeroles.Accountadmin).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeEvolveSchema})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithSchemaObject(sdk.Object{ObjectType: sdk.ObjectTypeTable, Name: sdk.NewObjectIdentifierFromFullyQualifiedName(table.ID().FullyQualifiedName())}))))
		require.NoError(t, err)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		testClientHelper().Stage.PutOnStage(t, stage.ID(), "schema_evolution_record.json")

		testClientHelper().Stage.CopyIntoTableFromFile(t, table.ID(), stage.ID(), "schema_evolution_record.json")

		currentColumns := testClientHelper().Table.GetTableColumnsFor(t, table.ID())
		require.Len(t, currentColumns, 2)
		assert.NotEmpty(t, currentColumns[1].SchemaEvolutionRecord)

		descColumns, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeColumnsTableRequest(id))
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

	t.Run("show parameters", func(t *testing.T) {
		table, cleanup := testClientHelper().Table.Create(t)
		t.Cleanup(cleanup)

		parametersDirect, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{Table: table.ID()},
		})
		require.NoError(t, err)
		require.NotEmpty(t, parametersDirect)

		parameters, err := client.Tables.ShowParameters(ctx, table.ID())
		require.NoError(t, err)
		require.Equal(t, parametersDirect, parameters)
	})

	t.Run("show parameters of a dropped table", func(t *testing.T) {
		table, cleanup := testClientHelper().Table.Create(t)
		cleanup()

		_, err := client.Tables.ShowParameters(ctx, table.ID())
		require.Error(t, err)
	})
}
