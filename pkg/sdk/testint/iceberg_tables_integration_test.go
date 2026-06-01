//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_IcebergTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	catalog := sdk.IcebergTableCatalogSnowflake
	snowflakeManagedExternalVolume := sdk.NewAccountObjectIdentifier("SNOWFLAKE_MANAGED")

	// createExternalVolume creates a writable S3-backed external volume and registers
	// a cleanup. It skips t if the required AWS env vars are absent — call it only in
	// sub-tests that actually need the volume.
	createExternalVolume := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleArn := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
		storageLocations := []sdk.ExternalVolumeStorageLocationItem{
			{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
				Name: "iceberg_table_test_location",
				S3StorageLocationParams: &sdk.S3StorageLocationParams{
					StorageProvider:   sdk.S3StorageProviderS3,
					StorageAwsRoleArn: awsRoleArn,
					StorageBaseUrl:    awsBucketUrl,
				},
			}},
		}
		volumeId, volumeCleanup := testClientHelper().ExternalVolume.CreateWithRequest(t,
			sdk.NewCreateExternalVolumeRequest(testClientHelper().Ids.RandomAccountObjectIdentifier(), storageLocations),
		)
		t.Cleanup(volumeCleanup)
		return volumeId
	}

	basicAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier, baseLocation string) {
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
			HasBaseLocation(baseLocation).
			HasCanWriteMetadata(true).
			HasNoComment().
			HasNoNameMapping().
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatus("").
			HasPartitionSpecs("[]").
			HasCurrentPartitionSpecId(0).
			HasIcebergTableFormatVersion(2),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 1)

		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
			HasName("ID").
			HasType("NUMBER(38,0)").
			HasSourceIcebergType("long").
			HasKind("COLUMN").
			HasIsNullable(true).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(false).
			HasNoCheck().
			HasNoExpression().
			HasComment("").
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)
	}

	completeAssertions := func(t *testing.T, id sdk.SchemaObjectIdentifier, baseLocation string, volumeId sdk.AccountObjectIdentifier) {
		t.Helper()

		assertThatObject(t, objectassert.IcebergTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner("ACCOUNTADMIN").
			HasExternalVolumeName(volumeId).
			HasCatalogName(sdk.NewAccountObjectIdentifier("SNOWFLAKE")).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasNoCatalogTableName().
			HasNoCatalogNamespace().
			HasBaseLocation(baseLocation).
			HasCanWriteMetadata(true).
			HasComment("integration test").
			HasNoNameMapping().
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatus("").
			HasCurrentPartitionSpecId(0).
			HasIcebergTableFormatVersion(2),
		)

		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		require.Len(t, details, 10)

		// ID — NOT NULL, PK=true, comment="id column"
		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[0]).
			HasName("ID").
			HasType("NUMBER(38,0)").
			HasSourceIcebergType("long").
			HasKind("COLUMN").
			HasIsNullable(false).
			HasNoDefault().
			HasPrimaryKey(true).
			HasUniqueKey(false).
			HasNoCheck().
			HasNoExpression().
			HasComment("id column").
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		// FK_ID — nullable, no pk/uk
		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[1]).
			HasName("FK_ID").
			HasType("NUMBER(38,0)").
			HasSourceIcebergType("long").
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

		// EVENT_TS
		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[2]).
			HasName("EVENT_TS").
			HasType("TIMESTAMP_NTZ(9)").
			HasSourceIcebergType("timestamptz").
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

		// REGION
		assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[3]).
			HasName("REGION").
			HasType("VARCHAR(16777216)").
			HasSourceIcebergType("string").
			HasKind("COLUMN").
			HasIsNullable(true).
			HasNoDefault().
			HasPrimaryKey(false).
			HasUniqueKey(true).
			HasNoCheck().
			HasNoExpression().
			HasNoComment().
			HasNoPolicyName().
			HasNoPrivacyDomain().
			HasNoNameMapping().
			HasNoWriteDefault(),
		)

		// BUCKET_COL, TRUNC_COL, YEAR_COL, MONTH_COL, DAY_COL, HOUR_COL
		colDefs := []struct {
			name          string
			typ           string
			sourceIceberg string
		}{
			{"BUCKET_COL", "VARCHAR(16777216)", "string"},
			{"TRUNC_COL", "VARCHAR(16777216)", "string"},
			{"YEAR_COL", "TIMESTAMP_NTZ(9)", "timestamptz"},
			{"MONTH_COL", "TIMESTAMP_NTZ(9)", "timestamptz"},
			{"DAY_COL", "TIMESTAMP_NTZ(9)", "timestamptz"},
			{"HOUR_COL", "TIMESTAMP_NTZ(9)", "timestamptz"},
		}
		for i, def := range colDefs {
			assertThatObject(t, objectassert.IcebergTableDetailsFromObject(t, &details[4+i]).
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
	}

	// print all SHOW and DESCRIBE fields for discovery
	logIcebergTableFields := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		obj, err := client.IcebergTables.ShowByID(ctx, id)
		require.NoError(t, err)
		t.Logf("SHOW fields:\n"+
			"  Name=%q DatabaseName=%q SchemaName=%q Owner=%v\n"+
			"  ExternalVolumeName=%v CatalogName=%v IcebergTableType=%q\n"+
			"  CatalogTableName=%v CatalogNamespace=%v BaseLocation=%q\n"+
			"  CanWriteMetadata=%v Comment=%v NameMapping=%v\n"+
			"  OwnerRoleType=%q CatalogSyncName=%q AutoRefreshStatus=%q\n"+
			"  PartitionSpecs=%q CurrentPartitionSpecId=%d IcebergTableFormatVersion=%d",
			obj.Name, obj.DatabaseName, obj.SchemaName, obj.Owner,
			obj.ExternalVolumeName, obj.CatalogName, obj.IcebergTableType,
			obj.CatalogTableName, obj.CatalogNamespace, obj.BaseLocation,
			obj.CanWriteMetadata, obj.Comment, obj.NameMapping,
			obj.OwnerRoleType, obj.CatalogSyncName, obj.AutoRefreshStatus,
			obj.PartitionSpecs, obj.CurrentPartitionSpecId, obj.IcebergTableFormatVersion,
		)
		details, err := client.IcebergTables.Describe(ctx, id)
		require.NoError(t, err)
		for i, d := range details {
			t.Logf("DESCRIBE col[%d]:\n"+
				"  Name=%q Type=%q SourceIcebergType=%q Kind=%q\n"+
				"  IsNullable=%v Default=%v PrimaryKey=%v UniqueKey=%v\n"+
				"  Check=%v Expression=%v Comment=%v PolicyName=%v\n"+
				"  PrivacyDomain=%v NameMapping=%v WriteDefault=%v",
				i, d.Name, d.Type, d.SourceIcebergType, d.Kind,
				d.IsNullable, d.Default, d.PrimaryKey, d.UniqueKey,
				d.Check, d.Expression, d.Comment, d.PolicyName,
				d.PrivacyDomain, d.NameMapping, d.WriteDefault,
			)
		}
	}

	t.Run("create Snowflake managed: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		baseLocation := random.AlphaLowerN(8)

		err := client.IcebergTables.Create(ctx, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{
				{Name: "ID", ColumnType: testdatatypes.DataTypeNumber},
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().IcebergTable.DropFunc(t, id))

		logIcebergTableFields(t, id)
		basicAssertions(t, id, baseLocation)
	})

	t.Run("create Snowflake managed: all options", func(t *testing.T) {
		volumeId := createExternalVolume(t)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		baseLocation := random.AlphaLowerN(8)

		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)

		aggregationPolicyId, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)

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
							RefColumn:  new("ID"),
						},
					},
				},
				{Name: "EVENT_TS", ColumnType: testdatatypes.DataTypeTimestampNTZIceberg},
				{Name: "REGION", ColumnType: testdatatypes.DataTypeVarcharIceberg},
				{Name: "BUCKET_COL", ColumnType: testdatatypes.DataTypeVarcharIceberg},
				{Name: "TRUNC_COL", ColumnType: testdatatypes.DataTypeVarcharIceberg},
				{Name: "YEAR_COL", ColumnType: testdatatypes.DataTypeTimestampNTZIceberg},
				{Name: "MONTH_COL", ColumnType: testdatatypes.DataTypeTimestampNTZIceberg},
				{Name: "DAY_COL", ColumnType: testdatatypes.DataTypeTimestampNTZIceberg},
				{Name: "HOUR_COL", ColumnType: testdatatypes.DataTypeTimestampNTZIceberg},
			},
			OutOfLineConstraint: []sdk.TableOutOfLineConstraintRequest{
				{
					UniquePK: &sdk.TableOutOfLineUniquePKRequest{
						Name:    new("uq_region"),
						Unique:  new(true),
						Columns: []sdk.Column{{Value: "REGION"}},
					},
				},
			},
		}

		req := sdk.NewCreateIcebergTableRequest(id, colDef).
			WithIfNotExists(true).
			WithCatalog(catalog).
			// TODO: these are commmented out for now because the current external volume is not writable.
			// WithExternalVolume(volumeId).
			// WithBaseLocation(baseLocation).
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

		logIcebergTableFields(t, id)
		completeAssertions(t, id, baseLocation, volumeId)
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
		volumeId := createExternalVolume(t)
		obj, _ := testClientHelper().IcebergTable.Create(t, volumeId)
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
		volumeId := createExternalVolume(t)
		obj1, _ := testClientHelper().IcebergTable.Create(t, volumeId)
		obj2, _ := testClientHelper().IcebergTable.Create(t, volumeId)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest())
		require.NoError(t, err)

		assert.Contains(t, returned, *obj1)
		assert.Contains(t, returned, *obj2)
	})

	t.Run("show iceberg tables: with like option", func(t *testing.T) {
		volumeId := createExternalVolume(t)
		obj1, _ := testClientHelper().IcebergTable.Create(t, volumeId)
		obj2, _ := testClientHelper().IcebergTable.Create(t, volumeId)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest().
			WithLike(sdk.Like{Pattern: new(obj1.Name)}))
		require.NoError(t, err)

		assert.Contains(t, returned, *obj1)
		assert.NotContains(t, returned, *obj2)
	})

	t.Run("show iceberg tables: with in schema option", func(t *testing.T) {
		volumeId := createExternalVolume(t)
		obj, _ := testClientHelper().IcebergTable.Create(t, volumeId)

		returned, err := client.IcebergTables.Show(ctx, sdk.NewShowIcebergTableRequest().
			WithIn(sdk.In{Schema: obj.ID().SchemaId()}))
		require.NoError(t, err)

		assert.Contains(t, returned, *obj)
	})

	t.Run("describe iceberg table: existing", func(t *testing.T) {
		volumeId := createExternalVolume(t)
		obj, _ := testClientHelper().IcebergTable.Create(t, volumeId)

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
