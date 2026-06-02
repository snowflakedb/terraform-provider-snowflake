//go:build non_account_level_tests

package testint

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_IcebergTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertPolicyReference := func(t *testing.T, policyRef sdk.PolicyReference,
		policyId sdk.SchemaObjectIdentifier,
		policyKind sdk.PolicyKind,
		tableId sdk.SchemaObjectIdentifier,
		refColumnName *string,
	) {
		t.Helper()
		assert.Equal(t, policyId.Name(), policyRef.PolicyName)
		assert.Equal(t, policyKind, policyRef.PolicyKind)
		assert.Equal(t, tableId.Name(), policyRef.RefEntityName)
		assert.Equal(t, "ICEBERG_TABLE", policyRef.RefEntityDomain)
		assert.Equal(t, "ACTIVE", *policyRef.PolicyStatus)
		if refColumnName != nil {
			assert.NotNil(t, policyRef.RefColumnName)
			assert.Equal(t, *refColumnName, *policyRef.RefColumnName)
		} else {
			assert.Nil(t, policyRef.RefColumnName)
		}
	}

	catalog := sdk.IcebergTableCatalogSnowflake
	snowflakeManagedExternalVolume := sdk.NewAccountObjectIdentifier("SNOWFLAKE_MANAGED")

	basicAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		assertThatObject(t, objectassert.IcebergTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner("ACCOUNTADMIN").
			HasExternalVolumeName(snowflakeManagedExternalVolume).
			HasCatalogName(sdk.NewAccountObjectIdentifier("SNOWFLAKE")).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasNoCatalogTableName().
			HasNoCatalogNamespace().
			HasBaseLocationPrefix(id).
			HasCanWriteMetadata(true).
			HasComment("").
			HasNoNameMapping().
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatus("").
			HasPartitionSpecsJson([]sdk.IcebergTablePartitionSpec{
				{
					SpecId: 0,
					Fields: []sdk.IcebergTablePartitionSpecField{},
				},
			}).
			HasCurrentPartitionSpecId(0).
			HasIcebergTableFormatVersion(2),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
			HasName("ID").
			HasType(testdatatypes.DataTypeNumber_38_0).
			HasSourceIcebergType(testdatatypes.DataTypeDecimal_38_0.ToSql()).
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(false).
			HasNoCheck().
			HasNoExpression().
			HasNoComment().
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		assertThatObject(t, objectparametersassert.IcebergTableParameters(t, id).
			HasAllDefaultsExplicit(),
		)
	}

	completeAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier, policyId sdk.SchemaObjectIdentifier) {
		t.Helper()

		assertThatObject(t, objectassert.IcebergTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner("ACCOUNTADMIN").
			HasExternalVolumeName(snowflakeManagedExternalVolume).
			HasCatalogName(sdk.NewAccountObjectIdentifier("SNOWFLAKE")).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasNoCatalogTableName().
			HasNoCatalogNamespace().
			HasBaseLocationPrefix(id).
			HasCanWriteMetadata(true).
			HasComment("integration test").
			HasNoNameMapping().
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatus("").
			HasCurrentPartitionSpecId(0).
			HasPartitionSpecsJson([]sdk.IcebergTablePartitionSpec{
				{
					SpecId: 0,
					Fields: []sdk.IcebergTablePartitionSpecField{
						{FieldId: 1000, Name: "REGION", SourceId: 4, Transform: "identity"},
						{FieldId: 1001, Name: "BUCKET_COL_bucket_4", SourceId: 5, Transform: "bucket[4]"},
						{FieldId: 1002, Name: "TRUNC_COL_trunc_10", SourceId: 6, Transform: "truncate[10]"},
						{FieldId: 1003, Name: "YEAR_COL_year", SourceId: 7, Transform: "year"},
						{FieldId: 1004, Name: "MONTH_COL_month", SourceId: 8, Transform: "month"},
						{FieldId: 1005, Name: "DAY_COL_day", SourceId: 9, Transform: "day"},
						{FieldId: 1006, Name: "HOUR_COL_hour", SourceId: 10, Transform: "hour"},
					},
				},
			}).
			HasIcebergTableFormatVersion(2),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 11)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
			HasName("ID").
			HasType(testdatatypes.DataTypeNumber_38_0).
			HasSourceIcebergType(testdatatypes.DataTypeDecimal_38_0.ToSql()).
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(true).
			HasUniqueKey(false).
			HasCheck("ID > 0").
			HasNoExpression().
			HasComment("id column").
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[1]).
			HasName("FK_ID").
			HasType(testdatatypes.DataTypeNumber_38_0).
			HasSourceIcebergType(testdatatypes.DataTypeDecimal_38_0.ToSql()).
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(false).
			HasNoCheck().
			HasNoExpression().
			HasNoComment().
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[2]).
			HasName("EVENT_TS").
			HasType(testdatatypes.DataTypeTimestampNTZ_6).
			HasSourceIcebergType("TIMESTAMP").
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(false).
			HasNoCheck().
			HasNoExpression().
			HasNoComment().
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[3]).
			HasName("REGION").
			HasType(testdatatypes.DataTypeVarcharIceberg).
			HasSourceIcebergType("STRING").
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(true).
			HasNoCheck().
			HasNoExpression().
			HasNoComment().
			HasPolicyName(policyId).
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		colDefs := []struct {
			name          string
			typ           datatypes.DataType
			sourceIceberg string
		}{
			{"BUCKET_COL", testdatatypes.DataTypeVarcharIceberg, "STRING"},
			{"TRUNC_COL", testdatatypes.DataTypeVarcharIceberg, "STRING"},
			{"YEAR_COL", testdatatypes.DataTypeTimestampNTZ_6, "TIMESTAMP"},
			{"MONTH_COL", testdatatypes.DataTypeTimestampNTZ_6, "TIMESTAMP"},
			{"DAY_COL", testdatatypes.DataTypeTimestampNTZ_6, "TIMESTAMP"},
			{"HOUR_COL", testdatatypes.DataTypeTimestampNTZ_6, "TIMESTAMP"},
		}
		for i, def := range colDefs {
			assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[4+i]).
				HasName(def.name).
				HasType(def.typ).
				HasSourceIcebergType(def.sourceIceberg).
				HasKind("COLUMN").
				HasIsNullable(false).
				HasNoDefault().
				HasPrimaryKey(false).
				HasUniqueKey(false).
				HasNoCheck().
				HasNoExpression().
				HasNoComment().
				HasNoPolicyName().
				HasNoPrivacyDomain().
				HasNoNameMapping().
				HasNoWriteDefault(),
			)
		}
		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[10]).
			HasName("STATUS").
			HasType(testdatatypes.DataTypeVarcharIceberg).
			HasSourceIcebergType("STRING").
			HasKind("COLUMN").
			HasIsNullable(false).
			HasDefault("'active'").
			HasPrimaryKey(false).
			HasUniqueKey(false).
			HasCheck("STATUS IN ('active', 'inactive')").
			HasNoExpression().
			HasNoComment().
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)
		// TODO (next PRs): add assertions for the out-of-line constraints.
	}

	t.Run("create Snowflake managed: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		basicAssertions(t, id)
	})

	t.Run("create Snowflake managed: all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)

		aggregationPolicyId, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)

		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarcharIceberg)
		t.Cleanup(maskingPolicyCleanup)

		projectionPolicyId, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)

		fkRefTable, fkRefCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForIcebergTable(t)
		t.Cleanup(fkRefCleanup)

		colDef := sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{
					Name:       "ID",
					ColumnType: testdatatypes.DataTypeNumber,
					NotNull:    new(true),
					InlineConstraint: &sdk.TableColumnInlineConstraintRequest{
						UniquePK: &sdk.TableColumnInlineUniquePKRequest{
							Name:       new("pk_id"),
							PrimaryKey: new(true),
						},
					},
					Comment: new("id column"),
				},
				{
					Name:       "FK_ID",
					ColumnType: testdatatypes.DataTypeNumber,
					InlineConstraint: &sdk.TableColumnInlineConstraintRequest{
						FK: &sdk.TableColumnInlineFKRequest{
							Name:       new("fk_ref"),
							References: fkRefTable.ID(),
							RefColumn:  []sdk.Column{{Value: "ID"}},
						},
					},
					// TODO (next PRs): it looks like masking policy and projection policy cannot be created at the same time.
					// Investigate this because according to the documentation, they can be created at the same time.
					ProjectionPolicy: sdk.NewTableColumnProjectionPolicyRequest(projectionPolicyId),
				},
				{Name: "EVENT_TS", ColumnType: testdatatypes.DataTypeTimestampNTZ_6},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: sdk.NewTableColumnMaskingPolicyRequest(maskingPolicy.ID()).WithUsing([]sdk.Column{{Value: "REGION"}})},
				{Name: "BUCKET_COL", ColumnType: testdatatypes.DataTypeVarcharIceberg},
				{Name: "TRUNC_COL", ColumnType: testdatatypes.DataTypeVarcharIceberg},
				{Name: "YEAR_COL", ColumnType: testdatatypes.DataTypeTimestampNTZ_6},
				{Name: "MONTH_COL", ColumnType: testdatatypes.DataTypeTimestampNTZ_6},
				{Name: "DAY_COL", ColumnType: testdatatypes.DataTypeTimestampNTZ_6},
				{Name: "HOUR_COL", ColumnType: testdatatypes.DataTypeTimestampNTZ_6},
				{
					Name:         "STATUS",
					ColumnType:   testdatatypes.DataTypeVarcharIceberg,
					DefaultValue: &sdk.ColumnDefaultValue{Expression: new("'active'")},
					InlineConstraint: &sdk.TableColumnInlineConstraintRequest{
						CH: &sdk.TableColumnInlineCHRequest{
							Name:       new("chk_status"),
							Expression: "STATUS IN ('active', 'inactive')",
						},
					},
				},
			},
			OutOfLineConstraint: []sdk.TableOutOfLineConstraintRequest{
				{
					UniquePK: &sdk.TableOutOfLineUniquePKRequest{
						Name:    new("uq_region"),
						Unique:  new(true),
						Columns: []sdk.Column{{Value: "REGION"}},
					},
				},
				{
					FK: &sdk.TableOutOfLineFKRequest{
						Name:       new("fk_out_ref"),
						Columns:    []sdk.Column{{Value: "FK_ID"}},
						References: fkRefTable.ID(),
						RefColumns: []sdk.Column{{Value: "ID"}},
					},
				},
				{
					CH: &sdk.TableOutOfLineCHRequest{
						Name:       new("chk_id_positive"),
						Expression: "ID > 0",
					},
				},
			},
		}

		req := sdk.NewCreateIcebergTableRequest(id, colDef).
			WithIfNotExists(true).
			WithCatalog(catalog).
			// TODO (next PRs): these are commented out for now because the current external volume is not writable.
			// WithExternalVolume(volumeId).
			// WithBaseLocation(baseLocation).
			// TODO (next PRs): handle CONTACT and CATALOG_SYNC.
			WithPartitionBy([]sdk.IcebergTablePartitionExpressionRequest{
				{Identity: new("REGION")},
				{Bucket: &sdk.IcebergTablePartitionBucketRequest{Args: sdk.IcebergTablePartitionBucketArgsRequest{NumBuckets: 4, Column: "BUCKET_COL"}}},
				{Truncate: &sdk.IcebergTablePartitionTruncateRequest{Args: sdk.IcebergTablePartitionTruncateArgsRequest{Width: 10, Column: "TRUNC_COL"}}},
				{Year: &sdk.IcebergTablePartitionYearRequest{Args: sdk.IcebergTablePartitionTimeArgsRequest{Column: "YEAR_COL"}}},
				{Month: &sdk.IcebergTablePartitionMonthRequest{Args: sdk.IcebergTablePartitionTimeArgsRequest{Column: "MONTH_COL"}}},
				{Day: &sdk.IcebergTablePartitionDayRequest{Args: sdk.IcebergTablePartitionTimeArgsRequest{Column: "DAY_COL"}}},
				{Hour: &sdk.IcebergTablePartitionHourRequest{Args: sdk.IcebergTablePartitionTimeArgsRequest{Column: "HOUR_COL"}}},
			}).
			WithPathLayout(sdk.IcebergTablePathLayoutHierarchical).
			WithTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
			WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			WithDataRetentionTimeInDays(1).
			WithMaxDataExtensionTimeInDays(8).
			WithChangeTracking(true).
			WithErrorLogging(true).
			WithComment("integration test").
			WithIcebergVersion(2).
			WithEnableIcebergMergeOnRead(true).
			WithEnableDataCompaction(true).
			WithRowAccessPolicy(sdk.IcebergTableRowAccessPolicyRequest{
				Name: rowAccessPolicy.ID(),
				On:   []sdk.Column{{Value: "ID"}},
			}).
			WithAggregationPolicy(sdk.IcebergTableAggregationPolicyRequest{
				AggregationPolicy: aggregationPolicyId,
			})

		err := client.IcebergTables.Create(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		completeAssertions(t, id, maskingPolicy.ID())

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 4)
		slices.SortFunc(references, func(x, y sdk.PolicyReference) int {
			return strings.Compare(string(x.PolicyKind), string(y.PolicyKind))
		})
		assertPolicyReference(t, references[0], aggregationPolicyId, sdk.PolicyKindAggregationPolicy, id, nil)
		assertPolicyReference(t, references[1], maskingPolicy.ID(), sdk.PolicyKindMaskingPolicy, id, new("REGION"))
		assertPolicyReference(t, references[2], projectionPolicyId, sdk.PolicyKindProjectionPolicy, id, new("FK_ID"))
		assertPolicyReference(t, references[3], rowAccessPolicy.ID(), sdk.PolicyKindRowAccessPolicy, id, nil)

		assertThatObject(t, objectparametersassert.IcebergTableParameters(t, id).
			HasAllowRowTimestamp(false).
			HasCatalog(testClientHelper().Database.TestDatabaseCatalog().Name()).
			HasCatalogSync("").
			HasDataMetricSchedule("60 MINUTES").
			HasDataRetentionTimeInDays(1).
			HasDefaultDdlCollation("").
			HasEnableDataCompaction(true).
			HasEnableIcebergMergeOnRead(true).
			HasExternalVolume("SNOWFLAKE_MANAGED").
			HasIcebergMergeOnReadBehavior("auto").
			HasLogEventLevel("OFF").
			HasMaxDataExtensionTimeInDays(8).
			HasOptimizeDataLayout(true).
			HasQuotedIdentifiersIgnoreCase(false).
			HasReplaceInvalidCharacters(false).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasTargetFileSize(sdk.IcebergTableTargetFileSize128mb),
		)
	})

	t.Run("create Snowflake managed: copy grants", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// CREATE OR REPLACE with COPY GRANTS copies the grants from the existing table of the same name.
		err = client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}).WithOrReplace(true).WithCopyGrants(true))
		require.NoError(t, err)

		_, err = client.IcebergTables.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("drop iceberg table: existing", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		err := client.IcebergTables.Drop(ctx, sdk.NewDropIcebergTableRequest(id))
		require.NoError(t, err)

		_, err = client.IcebergTables.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop iceberg table: non-existing", func(t *testing.T) {
		id := NonExistingSchemaObjectIdentifier

		err := client.IcebergTables.Drop(ctx, sdk.NewDropIcebergTableRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop iceberg table: non-existing with if exists option", func(t *testing.T) {
		id := NonExistingSchemaObjectIdentifier

		err := client.IcebergTables.Drop(ctx, sdk.NewDropIcebergTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("show iceberg tables: default", func(t *testing.T) {
		obj1, cleanup1 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup1)
		obj2, cleanup2 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup2)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest())
		require.NoError(t, err)

		assert.Contains(t, returned, *obj1)
		assert.Contains(t, returned, *obj2)
	})

	t.Run("show iceberg tables: with like option", func(t *testing.T) {
		obj1, cleanup1 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup1)
		obj2, cleanup2 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup2)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest().
			WithLike(sdk.Like{Pattern: new(obj1.Name)}))
		require.NoError(t, err)

		assert.Contains(t, returned, *obj1)
		assert.NotContains(t, returned, *obj2)
	})

	t.Run("show iceberg tables: with in schema option", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest().
			WithIn(sdk.In{Schema: obj.ID().SchemaId()}))
		require.NoError(t, err)

		assert.Contains(t, returned, *obj)
	})

	t.Run("describe iceberg table: existing", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)

		details, err := client.IcebergTables.Describe(ctx, obj.ID())
		require.NoError(t, err)

		require.Len(t, details, 1)
		assert.NotEmpty(t, details[0].Name)
	})

	t.Run("describe iceberg table: non-existing", func(t *testing.T) {
		_, err := client.IcebergTables.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
