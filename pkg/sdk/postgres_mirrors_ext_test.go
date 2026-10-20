package sdk

import (
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	externalVolumeErr := NewError("CreatePostgresMirrorOptions.arguments field: ExternalVolume can only be set when TargetTableType is ICEBERG")

	postgresMirrorsTests.Create.
		withDefaultOpts(func() *CreatePostgresMirrorOptions {
			return &CreatePostgresMirrorOptions{
				arguments: CreatePostgresMirrorArguments{
					MirrorName:       "orders_mirror",
					PostgresInstance: "my-instance",
					PostgresDatabase: "postgres",
					TargetDatabase:   "ORDERS_DB",
					PostgresTables:   []StringListItemWrapper{{Value: "public.orders"}},
				},
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_arguments_PostgresTables_emptyList",
			func(opts *CreatePostgresMirrorOptions) { opts.arguments.PostgresTables = []StringListItemWrapper{} },
			errAtLeastOneOf("CreatePostgresMirrorOptions.arguments", "PostgresTables", "PostgresSchemas"),
		).
		withAdditionalValidationCase(
			"validation_Create_arguments_ExternalVolume_withoutTargetTableType",
			func(opts *CreatePostgresMirrorOptions) { opts.arguments.ExternalVolume = new("MY_VOLUME") },
			externalVolumeErr,
		).
		withAdditionalValidationCase(
			"validation_Create_arguments_ExternalVolume_withSnowflakeTargetTableType",
			func(opts *CreatePostgresMirrorOptions) {
				opts.arguments.ExternalVolume = new("MY_VOLUME")
				opts.arguments.TargetTableType = new(PostgresMirrorTargetTableTypeSnowflake)
			},
			externalVolumeErr,
		).
		withExpectedSql(
			case_PostgresMirrors_sql_Create_basic,
			`CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR (mirror_name => 'orders_mirror', postgres_instance => 'my-instance', postgres_database => 'postgres', target_database => 'ORDERS_DB', postgres_tables => ARRAY_CONSTRUCT ('public.orders'))`,
		).
		withModifyAndExpectedSqlf(
			case_PostgresMirrors_sql_Create_all,
			func(opts *CreatePostgresMirrorOptions) {
				opts.arguments.PostgresTables = []StringListItemWrapper{{Value: "public.orders"}, {Value: "public.items"}}
				opts.arguments.PostgresSchemas = []StringListItemWrapper{{Value: "sales"}}
				opts.arguments.RefreshInterval = new("1 minute")
				opts.arguments.TargetTableType = new(PostgresMirrorTargetTableTypeIceberg)
				opts.arguments.Warehouse = new("MY_WH")
				opts.arguments.ExternalVolume = new("MY_VOLUME")
			},
			`CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR (mirror_name => 'orders_mirror', postgres_instance => 'my-instance', postgres_database => 'postgres', target_database => 'ORDERS_DB', postgres_tables => ARRAY_CONSTRUCT ('public.orders', 'public.items'), postgres_schemas => ARRAY_CONSTRUCT ('sales'), refresh_interval => '1 minute', target_table_type => 'ICEBERG', warehouse => 'MY_WH', external_volume => 'MY_VOLUME')`,
		).
		withAdditionalSqlCasef(
			"sql_Create_snowflakeTargetTableType",
			func(opts *CreatePostgresMirrorOptions) {
				opts.arguments.TargetTableType = new(PostgresMirrorTargetTableTypeSnowflake)
			},
			`CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR (mirror_name => 'orders_mirror', postgres_instance => 'my-instance', postgres_database => 'postgres', target_database => 'ORDERS_DB', postgres_tables => ARRAY_CONSTRUCT ('public.orders'), target_table_type => 'SNOWFLAKE')`,
		).
		withAdditionalSqlCasef(
			"sql_Create_schemasOnly",
			func(opts *CreatePostgresMirrorOptions) {
				opts.arguments.PostgresTables = nil
				opts.arguments.PostgresSchemas = []StringListItemWrapper{{Value: "sales"}, {Value: "inventory"}}
			},
			`CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR (mirror_name => 'orders_mirror', postgres_instance => 'my-instance', postgres_database => 'postgres', target_database => 'ORDERS_DB', postgres_schemas => ARRAY_CONSTRUCT ('sales', 'inventory'))`,
		).
		withAdditionalSqlCasef(
			"sql_Create_quotesInValues",
			func(opts *CreatePostgresMirrorOptions) {
				opts.arguments.PostgresInstance = "my'instance"
				opts.arguments.PostgresTables = []StringListItemWrapper{{Value: `public."MixedCase"`}, {Value: "public.o'reilly"}}
			},
			`CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR (mirror_name => 'orders_mirror', postgres_instance => 'my\'instance', postgres_database => 'postgres', target_database => 'ORDERS_DB', postgres_tables => ARRAY_CONSTRUCT ('public.\"MixedCase\"', 'public.o\'reilly'))`,
		)

	postgresMirrorsTests.Drop.
		withDefaultOpts(func() *DropPostgresMirrorOptions {
			return &DropPostgresMirrorOptions{
				arguments: DropPostgresMirrorArguments{MirrorName: "orders_mirror"},
			}
		}).
		withExpectedSql(
			case_PostgresMirrors_sql_Drop_basic,
			`CALL SNOWFLAKE.POSTGRES.DROP_MIRROR (mirror_name => 'orders_mirror')`,
		).
		withModifyAndExpectedSqlf(
			case_PostgresMirrors_sql_Drop_all,
			func(opts *DropPostgresMirrorOptions) { opts.arguments.DropTargetDatabase = new(false) },
			`CALL SNOWFLAKE.POSTGRES.DROP_MIRROR (mirror_name => 'orders_mirror', drop_target_database => false)`,
		).
		withAdditionalSqlCasef(
			"sql_Drop_dropTargetDatabase",
			func(opts *DropPostgresMirrorOptions) { opts.arguments.DropTargetDatabase = new(true) },
			`CALL SNOWFLAKE.POSTGRES.DROP_MIRROR (mirror_name => 'orders_mirror', drop_target_database => true)`,
		)

	postgresMirrorsTests.Describe.
		withDefaultOpts(func() *DescribePostgresMirrorOptions {
			return &DescribePostgresMirrorOptions{
				arguments: DescribePostgresMirrorArguments{MirrorName: "orders_mirror"},
			}
		}).
		withExpectedSql(
			case_PostgresMirrors_sql_Describe_basic,
			`CALL SNOWFLAKE.POSTGRES.DESCRIBE_MIRROR (mirror_name => 'orders_mirror')`,
		)
}

func TestPostgresMirrorDetailsRow_convert(t *testing.T) {
	lastApplyTime := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	lastOperationTime := time.Date(2026, 10, 8, 12, 5, 0, 0, time.UTC)

	t.Run("full row", func(t *testing.T) {
		row := postgresMirrorDetailsRow{
			MirrorName:             "orders_mirror",
			PostgresInstance:       sql.NullString{String: "my-instance", Valid: true},
			PostgresDatabase:       sql.NullString{String: "postgres", Valid: true},
			TargetDatabase:         sql.NullString{String: "ORDERS_DB", Valid: true},
			Status:                 sql.NullString{String: "ACTIVE", Valid: true},
			RefreshInterval:        sql.NullString{String: "1 minute", Valid: true},
			TargetTableType:        sql.NullString{String: "SNOWFLAKE", Valid: true},
			Warehouse:              sql.NullString{String: "MY_WH", Valid: true},
			PostgresTables:         sql.NullString{String: "[\n  \"public.orders\",\n  \"public.\\\"MixedCase\\\"\"\n]", Valid: true},
			PostgresSchemas:        sql.NullString{String: "[]", Valid: true},
			FeedOnly:               sql.NullBool{Bool: false, Valid: true},
			LastOperationId:        sql.NullInt64{Int64: 42, Valid: true},
			LastApplyTime:          sql.NullTime{Time: lastApplyTime, Valid: true},
			LastOperationTime:      sql.NullTime{Time: lastOperationTime, Valid: true},
			ErrorMessage:           sql.NullString{String: "apply failed", Valid: true},
			CatalogIntegrationName: sql.NullString{String: "MY_CATALOG", Valid: true},
			ApplyLagSeconds:        sql.NullInt64{Int64: 0, Valid: true},
			ApplyLagOperations:     sql.NullInt64{Int64: 7, Valid: true},
			LastAppliedCommitLsn:   sql.NullString{String: "123456789", Valid: true},
			RecentQueryDurations:   sql.NullString{String: "[\n  1.5,\n  2\n]", Valid: true},
			PostgresStatus:         sql.NullString{String: "{\n  \"slot_active\": true\n}", Valid: true},
			LastAppliedCommitTime:  sql.NullTime{},
		}

		details, err := row.convert()
		require.NoError(t, err)
		assert.Equal(t, "orders_mirror", details.Name)
		assert.Equal(t, "my-instance", details.PostgresInstance)
		assert.Equal(t, "postgres", details.PostgresDatabase)
		assert.Equal(t, "ORDERS_DB", details.TargetDatabase)
		assert.Equal(t, PostgresMirrorStatusActive, details.Status)
		assert.Equal(t, "1 minute", details.RefreshInterval)
		assert.Equal(t, PostgresMirrorTargetTableTypeSnowflake, details.TargetTableType)
		assert.Nil(t, details.ExternalVolume)
		assert.Equal(t, String("MY_WH"), details.Warehouse)
		assert.Equal(t, []string{"public.orders", `public."MixedCase"`}, details.PostgresTables)
		assert.NotNil(t, details.PostgresSchemas)
		assert.Empty(t, details.PostgresSchemas)
		assert.True(t, details.IsMembershipKnown())
		assert.Equal(t, Bool(false), details.FeedOnly)
		assert.Equal(t, Pointer(int64(42)), details.LastOperationId)
		assert.Equal(t, &lastApplyTime, details.LastApplyTime)
		assert.Equal(t, &lastOperationTime, details.LastOperationTime)
		assert.Equal(t, String("apply failed"), details.ErrorMessage)
		assert.Equal(t, String("MY_CATALOG"), details.CatalogIntegrationName)
		assert.Equal(t, Int(0), details.ApplyLagSeconds)
		assert.Equal(t, Pointer(int64(7)), details.ApplyLagOperations)
		assert.Equal(t, String("123456789"), details.LastAppliedCommitLsn)
		assert.Nil(t, details.LastAppliedCommitTime)
		assert.Equal(t, String("[1.5,2]"), details.RecentQueryDurations)
		assert.Equal(t, String(`{"slot_active":true}`), details.PostgresStatus)
	})

	t.Run("counts above 32 bits", func(t *testing.T) {
		above := int64(math.MaxInt32) + 1
		details, err := postgresMirrorDetailsRow{
			LastOperationId:    sql.NullInt64{Int64: above, Valid: true},
			ApplyLagOperations: sql.NullInt64{Int64: above, Valid: true},
			ApplyLagBytes:      sql.NullInt64{Int64: above, Valid: true},
		}.convert()
		require.NoError(t, err)
		require.NotNil(t, details.LastOperationId)
		require.NotNil(t, details.ApplyLagOperations)
		require.NotNil(t, details.ApplyLagBytes)
		assert.Equal(t, above, *details.LastOperationId)
		assert.Equal(t, above, *details.ApplyLagOperations)
		assert.Equal(t, above, *details.ApplyLagBytes)
	})

	t.Run("membership unknown when the source cannot be read", func(t *testing.T) {
		details, err := postgresMirrorDetailsRow{MirrorName: "orders_mirror"}.convert()
		require.NoError(t, err)
		assert.Nil(t, details.PostgresTables)
		assert.Nil(t, details.PostgresSchemas)
		assert.False(t, details.IsMembershipKnown())
		assert.Nil(t, details.Warehouse)
	})

	t.Run("invalid membership JSON", func(t *testing.T) {
		_, err := postgresMirrorDetailsRow{
			MirrorName:     "orders_mirror",
			PostgresTables: sql.NullString{String: "not json", Valid: true},
		}.convert()
		require.ErrorContains(t, err, "POSTGRES_TABLES")
	})

	t.Run("lowercase target table type", func(t *testing.T) {
		details, err := postgresMirrorDetailsRow{
			TargetTableType: sql.NullString{String: "snowflake", Valid: true},
		}.convert()
		require.NoError(t, err)
		assert.Equal(t, PostgresMirrorTargetTableTypeSnowflake, details.TargetTableType)
	})

	t.Run("unknown target table type", func(t *testing.T) {
		_, err := postgresMirrorDetailsRow{
			TargetTableType: sql.NullString{String: "PARQUET", Valid: true},
		}.convert()
		require.ErrorContains(t, err, "TARGET_TABLE_TYPE")
	})
}
