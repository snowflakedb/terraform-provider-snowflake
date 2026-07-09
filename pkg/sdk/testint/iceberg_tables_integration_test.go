//go:build non_account_level_tests

package testint

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_IcebergTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocationPrefix := "iceberg_test_table"
	metadataFilePath := baseLocationPrefix + "/metadata/v1.metadata.json"

	externalVolumeId, externalVolumeCleanup := testClientHelper().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogForIcebergFilesId, catalogForIcebergFilesCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogForIcebergFilesCleanup)

	dbForIcebergFilesId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	dbForIcebergFiles, dbForIcebergFilesCleanup := testClientHelper().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(dbForIcebergFilesId).WithCatalog(catalogForIcebergFilesId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbForIcebergFilesCleanup)
	schemaIdForIcebergFiles := sdk.NewDatabaseObjectIdentifier(dbForIcebergFiles.ID().Name(), "PUBLIC")

	deltaBaseLocation := "delta_lake_test_table"

	catalogForDeltaLakeId, catalogForDeltaLakeCleanup := testClientHelper().CatalogIntegration.CreateFunc(
		t,
		sdk.NewCreateCatalogIntegrationRequest(testClientHelper().Ids.RandomAccountObjectIdentifier(), true).
			WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatDelta)),
	)
	t.Cleanup(catalogForDeltaLakeCleanup)

	dbForDeltaLakeId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	dbForDeltaLake, dbForDeltaLakeCleanup := testClientHelper().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(dbForDeltaLakeId).WithCatalog(catalogForDeltaLakeId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbForDeltaLakeCleanup)
	schemaIdForDeltaLake := sdk.NewDatabaseObjectIdentifier(dbForDeltaLake.ID().Name(), "PUBLIC")

	contactId, contactCleanup := testClientHelper().Contact.Create(t)
	t.Cleanup(contactCleanup)

	rowAccessPolicy1, rowAccessPolicy1Cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
	t.Cleanup(rowAccessPolicy1Cleanup)

	assertPolicyReference := func(t *testing.T, policyRef sdk.PolicyReference,
		policyId sdk.SchemaObjectIdentifier,
		policyKind sdk.PolicyKind,
		tableId sdk.SchemaObjectIdentifier,
		refColumnName *string,
	) {
		t.Helper()
		ass := objectassert.PolicyReferenceFromObject(t, &policyRef).
			HasPolicyDb(policyId.DatabaseName()).
			HasPolicySchema(policyId.SchemaName()).
			HasPolicyName(policyId.Name()).
			HasPolicyKind(policyKind).
			HasRefDatabaseName(tableId.DatabaseName()).
			HasRefSchemaName(tableId.SchemaName()).
			HasRefEntityName(tableId.Name()).
			HasRefEntityDomain(string(sdk.PolicyEntityDomainIcebergTable)).
			HasPolicyStatus("ACTIVE")
		if refColumnName != nil {
			ass.HasRefColumnName(*refColumnName)
		} else {
			ass.HasNoRefColumnName()
		}
		assertThatObject(t, ass)
	}

	assertConstraint := func(t *testing.T, expected sdk.TableConstraintDetails, actual sdk.TableConstraintDetails) {
		t.Helper()
		assert.Equal(t, expected.ConstraintName, actual.ConstraintName)
		assert.Equal(t, expected.ConstraintType, actual.ConstraintType)
		assert.Equal(t, expected.Enforced, actual.Enforced)
		assert.Equal(t, expected.Rely, actual.Rely)
		assert.Equal(t, expected.IsDeferrable, actual.IsDeferrable)
		assert.Equal(t, expected.InitiallyDeferred, actual.InitiallyDeferred)
		assert.Equal(t, expected.Comment, actual.Comment)
		assert.Equal(t, expected.ConstraintCatalog, actual.ConstraintCatalog)
		assert.Equal(t, expected.ConstraintSchema, actual.ConstraintSchema)
		assert.Equal(t, expected.TableCatalog, actual.TableCatalog)
		assert.Equal(t, expected.TableSchema, actual.TableSchema)
		assert.Equal(t, expected.TableName, actual.TableName)
	}

	snowflakeCatalog := sdk.IcebergTableCatalogSnowflake
	snowflakeManagedExternalVolume := sdk.NewAccountObjectIdentifier("SNOWFLAKE_MANAGED")

	basicAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(snowflakeManagedExternalVolume).
				HasCatalogName(sdk.NewAccountObjectIdentifier("SNOWFLAKE")).
				HasIcebergTableType(sdk.IcebergTableTypeManaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasBaseLocationIdPrefix(id).
				HasCanWriteMetadata(true).
				HasComment("").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus().
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

		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
				HasName("ID").
				HasType(testdatatypes.DataTypeNumber_38_0).
				HasSourceIcebergType(testdatatypes.DataTypeDecimal_38_0.ToSql()).
				HasKind("COLUMN").
				HasIsNullable(true).
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

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasAllDefaultsExplicit(),
		)
	}

	completeAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier, policyId sdk.SchemaObjectIdentifier) {
		t.Helper()

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(sdk.NewAccountObjectIdentifier("SNOWFLAKE")).
				HasIcebergTableType(sdk.IcebergTableTypeManaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasBaseLocationPrefix(baseLocationPrefix).
				HasCanWriteMetadata(true).
				HasComment("integration test").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus().
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

		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
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

		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[1]).
				HasName("FK_ID").
				HasType(testdatatypes.DataTypeNumber_38_0).
				HasSourceIcebergType(testdatatypes.DataTypeDecimal_38_0.ToSql()).
				HasKind("COLUMN").
				HasIsNullable(true).
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

		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[2]).
				HasName("EVENT_TS").
				HasType(testdatatypes.DataTypeTimestampNTZ_6).
				HasSourceIcebergType("TIMESTAMP").
				HasKind("COLUMN").
				HasIsNullable(true).
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

		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[3]).
				HasName("REGION").
				HasType(testdatatypes.DataTypeVarcharIceberg).
				HasSourceIcebergType("STRING").
				HasKind("COLUMN").
				HasIsNullable(true).
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
			assertThatObject(
				t, objectassert.IcebergTableDetailsFromObject(t, &details[4+i]).
					HasName(def.name).
					HasType(def.typ).
					HasSourceIcebergType(def.sourceIceberg).
					HasKind("COLUMN").
					HasIsNullable(true).
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
		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[10]).
				HasName("STATUS").
				HasType(testdatatypes.DataTypeVarcharIceberg).
				HasSourceIcebergType("STRING").
				HasKind("COLUMN").
				HasIsNullable(true).
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

		info := testClientHelper().IcebergTable.GetIcebergTableInformation(t, id)
		assert.NotEmpty(t, info.MetadataLocation)
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
			WithCatalog(snowflakeCatalog).
			WithExternalVolume(externalVolumeId).
			WithBaseLocation(baseLocationPrefix).
			// TODO (next PRs): handle CATALOG_SYNC
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
			}).
			WithContact([]sdk.TableContact{
				{Purpose: "SUPPORT", Contact: contactId},
			})

		err := client.IcebergTables.Create(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		completeAssertions(t, id, maskingPolicy.ID())

		info := testClientHelper().IcebergTable.GetIcebergTableInformation(t, id)
		assert.NotEmpty(t, info.MetadataLocation)

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

		constraints, err := client.Tables.SelectTableConstraints(ctx, sdk.NewSelectTableConstraintsTableRequest(id.DatabaseId(), id.SchemaName(), id.Name()))
		require.NoError(t, err)
		require.Len(t, constraints, 4)
		// Sort the constraints because the order is not guaranteed.
		slices.SortFunc(constraints, func(x, y sdk.TableConstraintDetails) int {
			return strings.Compare(x.ConstraintName, y.ConstraintName)
		})
		assertConstraint(t, sdk.TableConstraintDetails{
			ConstraintName:    "fk_out_ref",
			ConstraintType:    sdk.TableConstraintTypeForeignKey,
			Enforced:          false,
			Rely:              false,
			IsDeferrable:      false,
			InitiallyDeferred: true,
			Comment:           nil,
			ConstraintCatalog: id.DatabaseName(),
			ConstraintSchema:  id.SchemaName(),
			TableCatalog:      id.DatabaseName(),
			TableSchema:       id.SchemaName(),
			TableName:         id.Name(),
		}, constraints[0])
		assertConstraint(t, sdk.TableConstraintDetails{
			ConstraintName:    "fk_ref",
			ConstraintType:    sdk.TableConstraintTypeForeignKey,
			Enforced:          false,
			Rely:              false,
			IsDeferrable:      false,
			InitiallyDeferred: true,
			Comment:           nil,
			ConstraintCatalog: id.DatabaseName(),
			ConstraintSchema:  id.SchemaName(),
			TableCatalog:      id.DatabaseName(),
			TableSchema:       id.SchemaName(),
			TableName:         id.Name(),
		}, constraints[1])
		assertConstraint(t, sdk.TableConstraintDetails{
			ConstraintName:    "pk_id",
			ConstraintType:    sdk.TableConstraintTypePrimaryKey,
			Enforced:          false,
			Rely:              false,
			IsDeferrable:      false,
			InitiallyDeferred: true,
			Comment:           nil,
			ConstraintCatalog: id.DatabaseName(),
			ConstraintSchema:  id.SchemaName(),
			TableCatalog:      id.DatabaseName(),
			TableSchema:       id.SchemaName(),
			TableName:         id.Name(),
		}, constraints[2])
		assertConstraint(t, sdk.TableConstraintDetails{
			ConstraintName:    "uq_region",
			ConstraintType:    sdk.TableConstraintTypeUnique,
			Enforced:          false,
			Rely:              false,
			IsDeferrable:      false,
			InitiallyDeferred: true,
			Comment:           nil,
			ConstraintCatalog: id.DatabaseName(),
			ConstraintSchema:  id.SchemaName(),
			TableCatalog:      id.DatabaseName(),
			TableSchema:       id.SchemaName(),
			TableName:         id.Name(),
		}, constraints[3])

		// TODO (next PRs): add assertions for CHECK constraints
		// like SELECT * FROM "A" . INFORMATION_SCHEMA.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'B' AND CONSTRAINT_TABLE = 'C'

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasAllowRowTimestamp(false).
				HasCatalog(testClientHelper().Database.TestDatabaseCatalog().Name()).
				HasCatalogSync("").
				HasDataMetricSchedule("60 MINUTES").
				HasDataRetentionTimeInDays(1).
				HasDefaultDdlCollation("").
				HasEnableDataCompaction(true).
				HasEnableIcebergMergeOnRead(true).
				HasExternalVolume(externalVolumeId.FullyQualifiedName()).
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

	t.Run("create Snowflake managed: cluster by", func(t *testing.T) {
		// PARTITION BY and CLUSTER BY are mutually exclusive for Iceberg tables (err 099207), so clustering is
		// tested in a separate table from the partitioned "all options" one above.
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}).WithClusterBy([]string{"ID", "REGION"}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// The clustering key is not returned from SHOW/DESCRIBE, so verify it with SYSTEM$CLUSTERING_INFORMATION.
		clusteringInfo, err := client.SystemFunctions.GetClusteringInformation(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "LINEAR(ID, REGION)", clusteringInfo.ClusterByKeys)
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

		info := testClientHelper().IcebergTable.GetIcebergTableInformation(t, id)
		assert.NotEmpty(t, info.MetadataLocation)
	})
	t.Run("create from iceberg files: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(catalogForIcebergFilesId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasCanWriteMetadata(true).
				HasComment("").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasAllowRowTimestamp(false).
				HasCatalog(catalogForIcebergFilesId.Name()).
				HasCatalogSync("").
				HasDataMetricSchedule("60 MINUTES").
				HasDataRetentionTimeInDays(1).
				HasDefaultDdlCollation("").
				HasEnableDataCompaction(true).
				HasEnableIcebergMergeOnRead(true).
				HasExternalVolume(externalVolumeId.Name()).
				HasIcebergMergeOnReadBehavior("auto").
				HasLogEventLevel("OFF").
				HasMaxDataExtensionTimeInDays(14).
				HasOptimizeDataLayout(true).
				HasQuotedIdentifiersIgnoreCase(false).
				HasReplaceInvalidCharacters(false).
				HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
				HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 2)

		// With this type, we cannot manage table columns. These assertions are basic and assert the values of the precreated file in the external volume.
		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
				HasKind("COLUMN").
				HasName("ID").
				HasType(testdatatypes.DataTypeNumber_19_0),
		)
		assertThatObject(
			t, objectassert.IcebergTableDetailsFromObject(t, &details[1]).
				HasKind("COLUMN").
				HasName("NAME").
				HasType(testdatatypes.DataTypeVarcharIceberg),
		)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Empty(t, references)

		info := testClientHelper().IcebergTable.GetIcebergTableInformation(t, id)
		assert.Equal(t, s3CompatBaseUrl+metadataFilePath, info.MetadataLocation)
	})

	t.Run("create from iceberg files: all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath).
			WithOrReplace(true).
			WithExternalVolume(externalVolumeId).
			WithCatalog(catalogForIcebergFilesId).
			WithReplaceInvalidCharacters(true).
			WithComment("integration test").
			WithContact([]sdk.TableContact{
				{Purpose: "SUPPORT", Contact: contactId},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(catalogForIcebergFilesId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasCanWriteMetadata(true).
				HasComment("integration test").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasAllowRowTimestamp(false).
				HasCatalog(catalogForIcebergFilesId.FullyQualifiedName()).
				HasCatalogSync("").
				HasDataMetricSchedule("60 MINUTES").
				HasDataRetentionTimeInDays(1).
				HasDefaultDdlCollation("").
				HasEnableDataCompaction(true).
				HasEnableIcebergMergeOnRead(true).
				HasExternalVolume(externalVolumeId.FullyQualifiedName()).
				HasIcebergMergeOnReadBehavior("auto").
				HasLogEventLevel("OFF").
				HasMaxDataExtensionTimeInDays(1).
				HasOptimizeDataLayout(true).
				HasQuotedIdentifiersIgnoreCase(false).
				HasReplaceInvalidCharacters(true).
				HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
				HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, details)
		for _, col := range details {
			assertThatObject(
				t, objectassert.IcebergTableDetailsFromObject(t, &col).
					HasKind("COLUMN").
					HasNoPolicyName().
					HasNoPrivacyDomain().
					HasNoWriteDefault(),
			)
		}

		info := testClientHelper().IcebergTable.GetIcebergTableInformation(t, id)
		assert.Equal(t, s3CompatBaseUrl+metadataFilePath, info.MetadataLocation)
	})

	t.Run("create from iceberg files: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// IF NOT EXISTS should not error when the table already exists.
		err = client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath).
			WithIfNotExists(true))
		require.NoError(t, err)
	})

	t.Run("create from delta lake: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForDeltaLake)

		err := client.IcebergTables.CreateFromDeltaLake(ctx, sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, deltaBaseLocation))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(catalogForDeltaLakeId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasCanWriteMetadata(true).
				HasComment("").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(catalogForDeltaLakeId.Name()).
				HasExternalVolume(externalVolumeId.Name()).
				HasReplaceInvalidCharacters(false),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, details)
		for _, col := range details {
			assertThatObject(
				t, objectassert.IcebergTableDetailsFromObject(t, &col).
					HasKind("COLUMN").
					HasNoPolicyName().
					HasNoPrivacyDomain().
					HasNoWriteDefault(),
			)
		}
	})

	t.Run("create from delta lake: all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.CreateFromDeltaLake(ctx, sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, deltaBaseLocation).
			WithOrReplace(true).
			WithExternalVolume(externalVolumeId).
			WithCatalog(catalogForDeltaLakeId).
			WithReplaceInvalidCharacters(true).
			WithAutoRefresh(true).
			WithComment("integration test").
			WithContact([]sdk.TableContact{
				{Purpose: "SUPPORT", Contact: contactId},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(catalogForDeltaLakeId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasNoCatalogTableName().
				HasNoCatalogNamespace().
				HasCanWriteMetadata(true).
				HasComment("integration test").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasAutoRefreshStatus(sdk.IcebergTableAutoRefreshStatus{
					CurrentSnapshotId:    0,
					PendingSnapshotCount: 0,
					ExecutionState:       "RUNNING",
				}),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(catalogForDeltaLakeId.FullyQualifiedName()).
				HasExternalVolume(externalVolumeId.FullyQualifiedName()).
				HasReplaceInvalidCharacters(true),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, details)
		for _, col := range details {
			assertThatObject(
				t, objectassert.IcebergTableDetailsFromObject(t, &col).
					HasKind("COLUMN").
					HasNoPolicyName().
					HasNoPrivacyDomain().
					HasNoWriteDefault(),
			)
		}
	})

	t.Run("create from delta lake: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForDeltaLake)

		err := client.IcebergTables.CreateFromDeltaLake(ctx, sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, deltaBaseLocation))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// IF NOT EXISTS should not error when the table already exists.
		err = client.IcebergTables.CreateFromDeltaLake(ctx, sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, deltaBaseLocation).
			WithIfNotExists(true))
		require.NoError(t, err)
	})

	t.Run("alter: add column", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("STATUS", testdatatypes.DataTypeVarcharIceberg)))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 2)
		assert.Equal(t, "STATUS", details[1].Name)
	})

	t.Run("alter: add column if not exists", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("STATUS", testdatatypes.DataTypeVarcharIceberg).WithIfNotExists(true)))
		require.NoError(t, err)

		// Adding the same column again with IF NOT EXISTS should not error.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("STATUS", testdatatypes.DataTypeVarcharIceberg).WithIfNotExists(true)))
		require.NoError(t, err)
	})

	t.Run("alter: drop column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "STATUS", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithDropColumnAction(*sdk.NewTableDropColumnActionRequest([]sdk.Column{{Value: "STATUS"}})))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		assert.Equal(t, "ID", details[0].Name)
	})

	t.Run("alter: rename column", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "OLD_NAME", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithRenameColumnAction(*sdk.NewTableRenameColumnActionRequest("OLD_NAME").WithNewName("NEW_NAME")))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		names := make([]string, len(details))
		for i, d := range details {
			names[i] = d.Name
		}
		assert.Contains(t, names, "NEW_NAME")
		assert.NotContains(t, names, "OLD_NAME")
	})

	t.Run("alter: alter column set/drop not null", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "STATUS", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("STATUS").WithSetNotNull(true),
			}))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		statusIdx := slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "STATUS" })
		require.NotEqual(t, -1, statusIdx)
		assert.False(t, details[statusIdx].IsNullable)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("STATUS").WithDropNotNull(true),
			}))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		statusIdx = slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "STATUS" })
		require.NotEqual(t, -1, statusIdx)
		assert.True(t, details[statusIdx].IsNullable)
	})

	t.Run("alter: alter column set/unset comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("ID").WithComment("my comment"),
			}))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		require.NotNil(t, details[0].Comment)
		assert.Equal(t, "my comment", *details[0].Comment)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("ID").WithUnsetComment(true),
			}))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		assert.Nil(t, details[0].Comment)
	})

	t.Run("alter: alter column set/drop write default", func(t *testing.T) {
		// WRITE DEFAULT requires an Iceberg v3 table; it is not supported on the default v2 tables
		// used in the other alter tests (err 093695: WRITE DEFAULT feature disabled).
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "STATUS", ColumnType: testdatatypes.DataTypeVarcharIceberg, DefaultValue: &sdk.ColumnDefaultValue{Expression: new("'active'")}},
			},
		}).WithIcebergVersion(3))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("STATUS").WithSetWriteDefault("'active'"),
			}))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		statusIdx := slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "STATUS" })
		require.NotEqual(t, -1, statusIdx)
		require.NotNil(t, details[statusIdx].WriteDefault)
		assert.Equal(t, "'active'", *details[statusIdx].WriteDefault)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("STATUS").WithDropWriteDefault(true),
			}))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		statusIdx = slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "STATUS" })
		require.NotEqual(t, -1, statusIdx)
		assert.Nil(t, details[statusIdx].WriteDefault)
	})

	t.Run("alter: set and unset masking policy on column", func(t *testing.T) {
		// Create the policies before the table so that on cleanup the table (which references a
		// policy) is dropped first, releasing it before the policies are dropped.
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeNumber_38_0)
		t.Cleanup(maskingPolicyCleanup)
		maskingPolicy2, maskingPolicy2Cleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeNumber_38_0)
		t.Cleanup(maskingPolicy2Cleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		// set
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetMaskingPolicyOnColumn(*sdk.NewTableSetColumnMaskingPolicyRequest("ID", maskingPolicy.ID())))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		require.NotNil(t, details[0].PolicyName)
		assert.Equal(t, maskingPolicy.ID().Name(), details[0].PolicyName.Name())

		// set with FORCE atomically replaces the existing masking policy on the column
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetMaskingPolicyOnColumn(*sdk.NewTableSetColumnMaskingPolicyRequest("ID", maskingPolicy2.ID()).WithForce(true)))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		require.NotNil(t, details[0].PolicyName)
		assert.Equal(t, maskingPolicy2.ID().Name(), details[0].PolicyName.Name())

		// set with USING explicitly listing the column the policy is applied to
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetMaskingPolicyOnColumn(*sdk.NewTableSetColumnMaskingPolicyRequest("ID", maskingPolicy.ID()).
				WithUsing([]sdk.Column{{Value: "ID"}}).WithForce(true)))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		require.NotNil(t, details[0].PolicyName)
		assert.Equal(t, maskingPolicy.ID().Name(), details[0].PolicyName.Name())

		// unset
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnsetMaskingPolicyOnColumn(*sdk.NewTableUnsetColumnMaskingPolicyRequest("ID")))
		require.NoError(t, err)

		details, err = client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)
		assert.Nil(t, details[0].PolicyName)
	})

	t.Run("alter: set projection policy on column", func(t *testing.T) {
		projectionPolicy, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)
		projectionPolicy2, projectionPolicy2Cleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicy2Cleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetProjectionPolicyOnColumn(*sdk.NewTableSetColumnProjectionPolicyRequest("ID", projectionPolicy)))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], projectionPolicy, sdk.PolicyKindProjectionPolicy, id, new("ID"))

		// set with FORCE atomically replaces the existing projection policy on the column
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetProjectionPolicyOnColumn(*sdk.NewTableSetColumnProjectionPolicyRequest("ID", projectionPolicy2).WithForce(true)))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], projectionPolicy2, sdk.PolicyKindProjectionPolicy, id, new("ID"))

		// unset
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnsetProjectionPolicyOnColumn(*sdk.NewTableUnsetColumnProjectionPolicyRequest("ID")))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assert.Empty(t, references)
	})

	t.Run("alter: clustering action - cluster by", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().WithClusterBy([]string{"REGION"})))
		require.NoError(t, err)

		// The clustering key is not returned from SHOW/DESCRIBE, so verify it with SYSTEM$CLUSTERING_INFORMATION.
		info, err := client.SystemFunctions.GetClusteringInformation(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "LINEAR(REGION)", info.ClusterByKeys)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().WithDropClusteringKey(true)))
		require.NoError(t, err)

		// After dropping the clustering key, the table is no longer clustered.
		_, err = client.SystemFunctions.GetClusteringInformation(ctx, id)
		require.ErrorIs(t, err, sdk.ErrTableNotClustered)
	})

	t.Run("alter: set and unset properties", func(t *testing.T) {
		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSet(
				*sdk.NewIcebergTableSetPropertiesRequest().
					// REPLACE_INVALID_CHARACTERS and LOG_EVENT_LEVEL are only supported for Iceberg tables using an external catalog - tested in other test
					WithComment("new comment").
					WithDataRetentionTimeInDays(2).
					// TODO (next PRs): handle CATALOG_SYNC
					WithMaxDataExtensionTimeInDays(5).
					WithEnableDataCompaction(false).
					WithTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
					WithEnableIcebergMergeOnRead(false).
					WithErrorLogging(false),
			))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasComment("new comment"),
		)
		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasDataRetentionTimeInDays(2).
				HasMaxDataExtensionTimeInDays(5).
				HasEnableDataCompaction(false).
				HasTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
				HasEnableIcebergMergeOnRead(false),
		)
		// Error logging is not returned from Snowflake.

		// unset
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnset(
				*sdk.NewIcebergTableUnsetPropertiesRequest().
					// REPLACE_INVALID_CHARACTERS and LOG_EVENT_LEVEL are only supported for Iceberg tables using an external catalog - tested in other test
					WithCatalogSync(true).
					WithDataRetentionTimeInDays(true).
					WithMaxDataExtensionTimeInDays(true).
					WithTargetFileSize(true).
					WithErrorLogging(true).
					WithEnableDataCompaction(true).
					WithEnableIcebergMergeOnRead(true).
					WithComment(true),
			))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasComment(""),
		)
		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasAllDefaultsExplicit(),
		)
	})
	t.Run("alter: add plus drop row access policy", func(t *testing.T) {
		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		// add
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddRowAccessPolicy(sdk.ViewAddRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy.ID(),
				On:              []sdk.Column{{Value: "ID"}},
			}))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], rowAccessPolicy.ID(), sdk.PolicyKindRowAccessPolicy, id, nil)

		// drop
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithDropRowAccessPolicy(sdk.ViewDropRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy.ID(),
			}))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assert.Empty(t, references)
	})

	t.Run("alter: drop and add row access policy", func(t *testing.T) {
		rowAccessPolicy2, rowAccessPolicy2Cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicy2Cleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}).WithRowAccessPolicy(sdk.IcebergTableRowAccessPolicyRequest{
			Name: rowAccessPolicy1.ID(),
			On:   []sdk.Column{{Value: "ID"}},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithDropAndAddRowAccessPolicy(sdk.ViewDropAndAddRowAccessPolicy{
				Drop: sdk.ViewDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicy1.ID()},
				Add:  sdk.ViewAddRowAccessPolicy{RowAccessPolicy: rowAccessPolicy2.ID(), On: []sdk.Column{{Value: "ID"}}},
			}))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], rowAccessPolicy2.ID(), sdk.PolicyKindRowAccessPolicy, id, nil)
	})

	t.Run("alter: drop all row access policies", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}).WithRowAccessPolicy(sdk.IcebergTableRowAccessPolicyRequest{
			Name: rowAccessPolicy1.ID(),
			On:   []sdk.Column{{Value: "ID"}},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithDropAllRowAccessPolicies(true))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assert.Empty(t, references)
	})

	t.Run("alter: set and unset aggregation policy", func(t *testing.T) {
		aggregationPolicy, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)
		aggregationPolicy2, aggregationPolicy2Cleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicy2Cleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		// set with an explicit entity key
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetAggregationPolicy(*sdk.NewTableSetAggregationPolicyRequest(aggregationPolicy).
				WithEntityKey([]sdk.Column{{Value: "ID"}})))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], aggregationPolicy, sdk.PolicyKindAggregationPolicy, id, nil)

		// set with FORCE atomically replaces the existing aggregation policy
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetAggregationPolicy(*sdk.NewTableSetAggregationPolicyRequest(aggregationPolicy2).
				WithEntityKey([]sdk.Column{{Value: "ID"}}).
				WithForce(true)))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assertPolicyReference(t, references[0], aggregationPolicy2, sdk.PolicyKindAggregationPolicy, id, nil)

		// unset
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnsetAggregationPolicy(*sdk.NewTableUnsetAggregationPolicyRequest()))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assert.Empty(t, references)
	})

	t.Run("alter: alter column set data type", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "AMOUNT", ColumnType: testdatatypes.DataTypeNumber_2_0},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// Iceberg tables only support widening a NUMBER's precision (with the scale unchanged).
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAlterColumnAction([]sdk.IcebergTableAlterColumnActionRequest{
				*sdk.NewIcebergTableAlterColumnActionRequest("AMOUNT").WithDataType(testdatatypes.DataTypeNumber_19_0),
			}))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		amountIdx := slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "AMOUNT" })
		require.NotEqual(t, -1, amountIdx)
		assert.Equal(t, testdatatypes.DataTypeNumber_19_0.ToSql(), details[amountIdx].Type.ToSql())
	})

	t.Run("alter: clustering action - suspend and resume recluster", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}).WithClusterBy([]string{"REGION"}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().
				WithChangeReclusterState(*sdk.NewIcebergTableReclusterChangeStateRequest().WithState(sdk.ReclusterStateSuspend))))
		require.NoError(t, err)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().
				WithChangeReclusterState(*sdk.NewIcebergTableReclusterChangeStateRequest().WithState(sdk.ReclusterStateResume))))
		require.NoError(t, err)

		// Recluster (automatic clustering) state is only available in SHOW TABLES... It is not returned by DESCRIBE nor modeled on the sdk.IcebergTable SHOW output.
	})

	t.Run("alter: search optimization action - add and drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// Add EQUALITY and FULL_TEXT (with an ANALYZER argument - DEFAULT_ANALYZER, UNICODE_ANALYZER
		// or NO_OP_ANALYZER) in a single statement.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSearchOptimizationAction(*sdk.NewTableSearchOptimizationActionRequest().
				WithAdd(*sdk.NewTableAddSearchOptimizationRequest().
					WithOn([]sdk.TableSearchMethodWithTargetRequest{
						*sdk.NewTableSearchMethodWithTargetRequest(sdk.TableSearchMethodEquality).
							WithArgs(*sdk.NewTableSearchMethodArgsRequest().WithTargets([]string{"REGION"})),
						*sdk.NewTableSearchMethodWithTargetRequest(sdk.TableSearchMethodFullText).
							WithArgs(*sdk.NewTableSearchMethodArgsRequest().
								WithTargets([]string{"REGION"}).
								WithAnalyzer("DEFAULT_ANALYZER")),
					}))))
		require.NoError(t, err)

		details, err := client.Tables.DescribeSearchOptimization(ctx, sdk.NewDescribeSearchOptimizationTableRequest(id))
		require.NoError(t, err)
		require.Len(t, details, 2)
		objectassert.TableSearchOptimizationDetailsFromObject(t, &details[0]).
			HasExpressionId(1).
			HasActive(true).
			HasMethod(string(sdk.TableSearchMethodEquality)).
			HasTarget("REGION").
			HasTargetDataType(testdatatypes.DataTypeVarcharIceberg)
		objectassert.TableSearchOptimizationDetailsFromObject(t, &details[1]).
			HasExpressionId(2).
			HasActive(true).
			HasMethod(string(sdk.TableSearchMethodFullText)).
			HasTarget("REGION").
			HasTargetDataType(testdatatypes.DataTypeVarcharIceberg)

		// Drop by method/target: removes only the EQUALITY entry. The analyzer is not part of the matcher.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSearchOptimizationAction(*sdk.NewTableSearchOptimizationActionRequest().
				WithDrop(*sdk.NewTableDropSearchOptimizationRequest().
					WithOn([]sdk.TableDropSearchOptimizationOnRequest{
						*sdk.NewTableDropSearchOptimizationOnRequest().WithSearchMethodWithTarget(*sdk.NewTableSearchMethodWithTargetRequest(sdk.TableSearchMethodEquality).
							WithArgs(*sdk.NewTableSearchMethodArgsRequest().WithTargets([]string{"REGION"}))),
					}))))
		require.NoError(t, err)

		details, err = client.Tables.DescribeSearchOptimization(ctx, sdk.NewDescribeSearchOptimizationTableRequest(id))
		require.NoError(t, err)
		require.Len(t, details, 1)

		// Drop by column name: removes the remaining (FULL_TEXT) search optimization on the column.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSearchOptimizationAction(*sdk.NewTableSearchOptimizationActionRequest().
				WithDrop(*sdk.NewTableDropSearchOptimizationRequest().
					WithOn([]sdk.TableDropSearchOptimizationOnRequest{
						*sdk.NewTableDropSearchOptimizationOnRequest().WithColumnName("REGION"),
					}))))
		require.NoError(t, err)

		details, err = client.Tables.DescribeSearchOptimization(ctx, sdk.NewDescribeSearchOptimizationTableRequest(id))
		require.NoError(t, err)
		require.Empty(t, details)
	})

	t.Run("alter: set and unset join policy", func(t *testing.T) {
		// Create the policies before the table so that on cleanup the table (which references a
		// policy) is dropped first, releasing it before the policies are dropped.
		joinPolicy1, joinPolicy1Cleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(joinPolicy1Cleanup)
		joinPolicy2, joinPolicy2Cleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(joinPolicy2Cleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		// set
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetJoinPolicy(*sdk.NewTableSetJoinPolicyRequest(joinPolicy1)))
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, joinPolicy1.Name(), references[0].PolicyName)

		// set with FORCE atomically replaces the existing join policy
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSetJoinPolicy(*sdk.NewTableSetJoinPolicyRequest(joinPolicy2).WithForce(true)))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 1)
		assert.Equal(t, joinPolicy2.Name(), references[0].PolicyName)

		// unset
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnsetJoinPolicy(*sdk.NewTableUnsetJoinPolicyRequest()))
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		assert.Empty(t, references)
	})

	t.Run("alter: if exists on non-existing table", func(t *testing.T) {
		// ALTER ICEBERG TABLE ... IF EXISTS on a non-existing table should not error.
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true).
			WithSet(*sdk.NewIcebergTableSetPropertiesRequest().WithComment("noop")))
		require.NoError(t, err)
	})

	t.Run("alter: add column with options", func(t *testing.T) {
		// Create the policies before the table so that on cleanup the table (which references the
		// policies) is dropped first, releasing them before they are dropped.
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarcharIceberg)
		t.Cleanup(maskingPolicyCleanup)

		projectionPolicyId, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		obj, cleanup := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup)
		id := obj.ID()

		// ADD COLUMN with a masking policy and a tag.
		err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("DESCRIPTION", testdatatypes.DataTypeVarcharIceberg).
				WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(maskingPolicy.ID()).WithUsing([]sdk.Column{{Value: "DESCRIPTION"}})).
				WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "tag-value"}})))
		require.NoError(t, err)

		// ADD COLUMN with a projection policy.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("CODE", testdatatypes.DataTypeVarcharIceberg).
				WithProjectionPolicy(*sdk.NewTableColumnProjectionPolicyRequest(projectionPolicyId))))
		require.NoError(t, err)

		// ADD COLUMN with a default value and an inline check constraint.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithAddColumnAction(*sdk.NewIcebergTableAddColumnActionRequest("STATUS", testdatatypes.DataTypeVarcharIceberg).
				WithDefaultValue(sdk.ColumnDefaultValue{Expression: new("'active'")}).
				WithInlineConstraint(sdk.TableColumnInlineConstraintRequest{
					CH: &sdk.TableColumnInlineCHRequest{
						Name:       new("chk_status_added"),
						Expression: "STATUS IN ('active', 'inactive')",
					},
				})))
		require.NoError(t, err)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		// ID + DESCRIPTION + CODE + STATUS
		require.Len(t, details, 4)

		descriptionIdx := slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "DESCRIPTION" })
		require.NotEqual(t, -1, descriptionIdx)
		require.NotNil(t, details[descriptionIdx].PolicyName)
		assert.Equal(t, maskingPolicy.ID().Name(), details[descriptionIdx].PolicyName.Name())

		// the tag set on the DESCRIPTION column at ADD COLUMN time
		descriptionColumnId := sdk.NewTableColumnIdentifier(id.DatabaseName(), id.SchemaName(), id.Name(), "DESCRIPTION")
		tagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), descriptionColumnId, sdk.ObjectTypeColumn)
		require.NoError(t, err)
		require.NotNil(t, tagValue)
		assert.Equal(t, "tag-value", *tagValue)

		statusIdx := slices.IndexFunc(details, func(d sdk.IcebergTableDetails) bool { return d.Name == "STATUS" })
		require.NotEqual(t, -1, statusIdx)
		require.NotNil(t, details[statusIdx].Default)
		assert.Equal(t, "'active'", *details[statusIdx].Default)
		require.NotNil(t, details[statusIdx].Check)
		assert.Equal(t, "STATUS IN ('active', 'inactive')", *details[statusIdx].Check)

		// The masking policy (DESCRIPTION) and the projection policy (CODE) are both referenced.
		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, id, sdk.PolicyEntityDomainTable)
		require.NoError(t, err)
		require.Len(t, references, 2)
	})

	t.Run("alter: set and unset log event level and replace invalid characters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSet(
				*sdk.NewIcebergTableSetPropertiesRequest().
					WithLogEventLevel(sdk.IcebergTableLogEventLevelDebug).
					WithReplaceInvalidCharacters(true),
			))
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasLogEventLevel(string(sdk.IcebergTableLogEventLevelDebug)).
				HasReplaceInvalidCharacters(true),
		)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithUnset(
				*sdk.NewIcebergTableUnsetPropertiesRequest().
					WithLogEventLevel(true).
					WithReplaceInvalidCharacters(true),
			))
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasDefaultLogEventLevelValueExplicit().
				HasDefaultReplaceInvalidCharactersValueExplicit(),
		)
	})

	t.Run("alter: set auto refresh on table from iceberg files", func(t *testing.T) {
		// AUTO_REFRESH is only valid for unmanaged tables (created from iceberg files / Delta Lake),
		// not the Snowflake-managed tables used in the other alter tests.
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSet(*sdk.NewIcebergTableSetPropertiesRequest().WithAutoRefresh(true)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasAutoRefreshStatusNotEmpty(),
		)

		// Disable auto refresh again. UNSET is not supported.
		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).
			WithSet(*sdk.NewIcebergTableSetPropertiesRequest().WithAutoRefresh(false)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasNoAutoRefreshStatus(),
		)
	})

	t.Run("alter iceberg table from files: set and unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)

		err := client.IcebergTables.CreateFromIcebergFiles(ctx, sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithSet(
			*sdk.NewIcebergTableSetPropertiesRequest().
				WithComment("integration test comment").
				WithReplaceInvalidCharacters(true),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasComment("integration test comment"),
		)
		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasReplaceInvalidCharacters(true),
		)

		err = client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithUnset(
			*sdk.NewIcebergTableUnsetPropertiesRequest().
				WithComment(true).
				WithReplaceInvalidCharacters(true),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasComment(""),
		)
		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasReplaceInvalidCharacters(false),
		)
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

	t.Run("show iceberg tables: in account", func(t *testing.T) {
		obj1, cleanup1 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup1)
		obj2, cleanup2 := testClientHelper().IcebergTable.Create(t)
		t.Cleanup(cleanup2)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest().WithIn(sdk.In{Account: new(true)}))
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
			WithLike(sdk.Like{Pattern: new(obj1.Name)}).
			// The test setup creates a new database which changes the connection context, so we need to specify the schema explicitly.
			WithIn(sdk.In{Schema: obj1.ID().SchemaId()}))
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

// TestInt_IcebergTables_WithAdditionalDependencies covers creating Iceberg tables from an external
// catalog (AWS Glue and Iceberg REST). Unlike TestInt_IcebergTables, these require additional
// preconfigured dependencies we can't provide dynamically for now.
// TODO(SNOW-3725859): Provide the external volume and catalog integrations dynamically. Unskip and move these tests to the main test suite.
func TestInt_IcebergTables_WithAdditionalDependencies(t *testing.T) {
	t.Skip("Iceberg REST and AWS Glue tests require preconfigured external catalog integrations and are not run by default")

	client := testClient(t)
	ctx := testContext(t)

	// These tests reuse preexisting, manually configured dependencies instead of creating their own
	// external volume and catalog integrations.
	externalVolumeId := sdk.NewAccountObjectIdentifier("GLUE_EXTERNAL_VOLUME")
	awsGlueCatalogId := sdk.NewAccountObjectIdentifier("GLUE_CATALOG_INTEGRATION")
	restCatalogId := sdk.NewAccountObjectIdentifier("REST_CATALOG_INTEGRATION")

	catalogTableName := "TEST"
	catalogNamespace := "glue_iceberg_schema"

	dbForAwsGlueId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	dbForAwsGlue, dbForAwsGlueCleanup := testClientHelper().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(dbForAwsGlueId).WithCatalog(awsGlueCatalogId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbForAwsGlueCleanup)
	schemaIdForAwsGlue := sdk.NewDatabaseObjectIdentifier(dbForAwsGlue.ID().Name(), "PUBLIC")

	// Separate database wired to the Iceberg REST catalog integration and the external volume.
	dbForRestId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	dbForRest, dbForRestCleanup := testClientHelper().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(dbForRestId).WithCatalog(restCatalogId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbForRestCleanup)
	schemaIdForRest := sdk.NewDatabaseObjectIdentifier(dbForRest.ID().Name(), "PUBLIC")

	contactId, contactCleanup := testClientHelper().Contact.Create(t)
	t.Cleanup(contactCleanup)

	assertUnmanagedColumns := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, details)
		// With an external catalog, we cannot manage table columns; assert only the generic properties.
		for _, col := range details {
			assertThatObject(
				t, objectassert.IcebergTableDetailsFromObject(t, &col).
					HasKind("COLUMN").
					HasNoPolicyName().
					HasNoPrivacyDomain().
					HasNoWriteDefault(),
			)
		}
	}

	t.Run("create from aws glue: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForAwsGlue)

		err := client.IcebergTables.CreateFromAwsGlue(ctx, sdk.NewCreateFromAwsGlueIcebergTableRequest(id, catalogTableName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(awsGlueCatalogId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasCatalogTableName(catalogTableName).
				HasCatalogNamespace(catalogNamespace).
				HasCanWriteMetadata(true).
				HasComment("").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(awsGlueCatalogId.Name()).
				HasExternalVolume(externalVolumeId.Name()).
				HasReplaceInvalidCharacters(false),
		)

		assertUnmanagedColumns(t, id)
	})

	t.Run("create from aws glue: all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.CreateFromAwsGlue(ctx, sdk.NewCreateFromAwsGlueIcebergTableRequest(id, catalogTableName).
			WithOrReplace(true).
			WithExternalVolume(externalVolumeId).
			WithCatalog(awsGlueCatalogId).
			WithCatalogNamespace(catalogNamespace).
			WithReplaceInvalidCharacters(true).
			WithAutoRefresh(true).
			WithComment("integration test").
			WithContact([]sdk.TableContact{
				{Purpose: "SUPPORT", Contact: contactId},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(awsGlueCatalogId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasCatalogTableName(catalogTableName).
				HasCatalogNamespace(catalogNamespace).
				HasCanWriteMetadata(true).
				HasComment("integration test").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasAutoRefreshStatusNotEmpty(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(awsGlueCatalogId.FullyQualifiedName()).
				HasExternalVolume(externalVolumeId.FullyQualifiedName()).
				HasReplaceInvalidCharacters(true),
		)

		assertUnmanagedColumns(t, id)
	})

	t.Run("create from aws glue: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForAwsGlue)

		err := client.IcebergTables.CreateFromAwsGlue(ctx, sdk.NewCreateFromAwsGlueIcebergTableRequest(id, catalogTableName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// IF NOT EXISTS should not error when the table already exists.
		err = client.IcebergTables.CreateFromAwsGlue(ctx, sdk.NewCreateFromAwsGlueIcebergTableRequest(id, catalogTableName).
			WithIfNotExists(true))
		require.NoError(t, err)
	})

	t.Run("create from iceberg rest: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForRest)

		err := client.IcebergTables.CreateFromIcebergRest(ctx, sdk.NewCreateFromIcebergRestIcebergTableRequest(id, catalogTableName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(restCatalogId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasCatalogTableName(catalogTableName).
				HasCatalogNamespace(catalogNamespace).
				HasCanWriteMetadata(true).
				HasComment("").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasNoAutoRefreshStatus(),
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(restCatalogId.Name()).
				HasExternalVolume(externalVolumeId.Name()).
				HasReplaceInvalidCharacters(false).
				HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
				HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto).
				HasEnableIcebergMergeOnRead(true).
				HasIcebergMergeOnReadBehavior("auto"),
		)

		assertUnmanagedColumns(t, id)
	})

	t.Run("create from iceberg rest: all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.IcebergTables.CreateFromIcebergRest(ctx, sdk.NewCreateFromIcebergRestIcebergTableRequest(id, catalogTableName).
			WithOrReplace(true).
			WithExternalVolume(externalVolumeId).
			WithCatalog(restCatalogId).
			WithCatalogNamespace(catalogNamespace).
			WithPathLayout(sdk.IcebergTablePathLayoutHierarchical).
			WithTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
			WithReplaceInvalidCharacters(true).
			WithAutoRefresh(true).
			WithComment("integration test").
			WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			WithIcebergMergeOnReadBehavior(sdk.IcebergTableIcebergMergeOnReadBehaviorEnabled).
			WithEnableIcebergMergeOnRead(true).
			WithContact([]sdk.TableContact{
				{Purpose: "SUPPORT", Contact: contactId},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		assertThatObject(
			t, objectassert.IcebergTable(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner("ACCOUNTADMIN").
				HasExternalVolumeName(externalVolumeId).
				HasCatalogName(restCatalogId).
				HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
				HasCatalogTableName(catalogTableName).
				HasCatalogNamespace(catalogNamespace).
				HasCanWriteMetadata(true).
				HasComment("integration test").
				HasNoNameMapping().
				HasOwnerRoleType("ROLE").
				HasCatalogSyncName("").
				HasAutoRefreshStatusNotEmpty(),
			// Path layout is not returned from Snowflake.
		)

		assertThatObject(
			t, objectparametersassert.IcebergTableParameters(t, id).
				HasCatalog(restCatalogId.FullyQualifiedName()).
				HasExternalVolume(externalVolumeId.FullyQualifiedName()).
				HasReplaceInvalidCharacters(true).
				HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
				HasTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
				HasEnableIcebergMergeOnRead(true).
				HasIcebergMergeOnReadBehavior(string(sdk.IcebergTableIcebergMergeOnReadBehaviorEnabled)),
		)

		assertUnmanagedColumns(t, id)
	})

	t.Run("create from iceberg rest: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForRest)

		err := client.IcebergTables.CreateFromIcebergRest(ctx, sdk.NewCreateFromIcebergRestIcebergTableRequest(id, catalogTableName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		// IF NOT EXISTS should not error when the table already exists.
		err = client.IcebergTables.CreateFromIcebergRest(ctx, sdk.NewCreateFromIcebergRestIcebergTableRequest(id, catalogTableName).
			WithIfNotExists(true))
		require.NoError(t, err)
	})
}

// TODO: alter set/unset CATALOG_SYNC - requires a configured catalog-sync integration (see also the create test TODO).
// TODO (next PRs): alter set CONTACT - SET CONTACT (...) on ALTER is unverified; contact is currently only covered at CREATE time. This should be ultimately handled similarly to tags.
