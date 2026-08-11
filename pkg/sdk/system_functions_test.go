package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/stretchr/testify/require"
)

func Test_ToBehaviorChangeBundleStatus(t *testing.T) {
	type test struct {
		input string
		want  BehaviorChangeBundleStatus
	}

	valid := []test{
		{input: "enabled", want: BehaviorChangeBundleStatusEnabled},
		{input: "ENABLED", want: BehaviorChangeBundleStatusEnabled},
		{input: "DISABLED", want: BehaviorChangeBundleStatusDisabled},
		{input: "RELEASED", want: BehaviorChangeBundleStatusReleased},
	}

	invalid := []string{
		"",
		"foo",
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBehaviorChangeBundleStatus(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := ToBehaviorChangeBundleStatus(in)
			require.Error(t, err)
		})
	}
}

func Test_getClusteringInformationOptions_SQL(t *testing.T) {
	id := NewSchemaObjectIdentifier("db", "schema", "table")

	buildColumns := func(columns ...string) []clusteringInformationColumn {
		return collections.Map(columns, func(col string) clusteringInformationColumn {
			return clusteringInformationColumn{Name: col}
		})
	}

	t.Run("without columns", func(t *testing.T) {
		opts := &getClusteringInformationOptions{
			arguments: &clusteringInformationArgs{Name: id},
		}
		got, err := structToSQL(opts)
		require.NoError(t, err)
		require.Equal(t, `SELECT SYSTEM$CLUSTERING_INFORMATION ('\"db\".\"schema\".\"table\"') AS "CLUSTERING_INFORMATION"`, got)
	})

	t.Run("with columns", func(t *testing.T) {
		opts := &getClusteringInformationOptions{
			arguments: &clusteringInformationArgs{
				Name:    id,
				Columns: buildColumns("REGION", "id"),
			},
		}
		got, err := structToSQL(opts)
		require.NoError(t, err)
		require.Equal(t, `SELECT SYSTEM$CLUSTERING_INFORMATION ('\"db\".\"schema\".\"table\"', '(\"REGION\", \"id\")') AS "CLUSTERING_INFORMATION"`, got)
	})

	t.Run("column with double quote is escaped", func(t *testing.T) {
		opts := &getClusteringInformationOptions{
			arguments: &clusteringInformationArgs{
				Name:    id,
				Columns: buildColumns(`we"ird`),
			},
		}
		got, err := structToSQL(opts)
		require.NoError(t, err)
		require.Equal(t, `SELECT SYSTEM$CLUSTERING_INFORMATION ('\"db\".\"schema\".\"table\"', '(\"we\"\"ird\")') AS "CLUSTERING_INFORMATION"`, got)
	})
}

func Test_parseClusteringInformation(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		// Output captured from SYSTEM$CLUSTERING_INFORMATION on a clustered table.
		raw := `{
  "cluster_by_keys" : "LINEAR(REGION)",
  "total_partition_count" : 1,
  "total_constant_partition_count" : 0,
  "average_overlaps" : 0.0,
  "average_depth" : 1.0,
  "partition_depth_histogram" : {
    "00000" : 0,
    "00001" : 1,
    "00002" : 0,
    "00016" : 0
  },
  "clustering_errors" : [ { "timestamp" : "2024-01-01 00:00:00", "error" : "some error" } ]
}`
		info, err := parseClusteringInformation(raw)
		require.NoError(t, err)
		require.NotNil(t, info)
		require.Equal(t, "LINEAR(REGION)", info.ClusterByKeys)
		require.Equal(t, 1, info.TotalPartitionCount)
		require.Equal(t, 0, info.TotalConstantPartitionCount)
		require.InDelta(t, 0.0, info.AverageOverlaps, 0.0001)
		require.InDelta(t, 1.0, info.AverageDepth, 0.0001)
		require.Equal(t, 1, info.PartitionDepthHistogram["00001"])
		require.Len(t, info.ClusteringErrors, 1)
		require.Equal(t, "2024-01-01 00:00:00", info.ClusteringErrors[0].Timestamp)
		require.Equal(t, "some error", info.ClusteringErrors[0].Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseClusteringInformation("not json")
		require.Error(t, err)
	})
}

func Test_getCatalogLinkedDatabaseConfigOptions_SQL(t *testing.T) {
	id := NewAccountObjectIdentifier("my_db")

	opts := &getCatalogLinkedDatabaseConfigOptions{
		arguments: &catalogLinkedDatabaseArgs{Name: id},
	}
	got, err := structToSQL(opts)
	require.NoError(t, err)
	require.Equal(t, `SELECT SYSTEM$GET_CATALOG_LINKED_DATABASE_CONFIG ('\"my_db\"') AS "CONFIG"`, got)
}

func Test_parseCatalogLinkedDatabaseConfig(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		// Output shape based on https://docs.snowflake.com/en/sql-reference/functions/system_get_catalog_linked_database_config.
		raw := `{
  "catalog_integration" : "MY_CATALOG_INT",
  "catalog_name" : null,
  "external_volume" : "MY_EXTERNAL_VOL",
  "sync_interval_seconds" : 60,
  "namespace_mode" : "FLATTEN_NESTED_NAMESPACE",
  "namespace_flatten_delimiter" : "-",
  "allowed_write_operations" : "ALL",
  "catalog_case_sensitivity" : "CASE_INSENSITIVE",
  "is_suspended" : false,
  "allowed_namespaces" : [ "ns1", "ns2" ],
  "blocked_namespaces" : [ "ns3" ]
}`
		config, err := parseCatalogLinkedDatabaseConfig(raw)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, "MY_CATALOG_INT", config.CatalogIntegration)
		require.Nil(t, config.CatalogName)
		require.NotNil(t, config.ExternalVolume)
		require.Equal(t, "MY_EXTERNAL_VOL", *config.ExternalVolume)
		require.NotNil(t, config.SyncIntervalSeconds)
		require.Equal(t, 60, *config.SyncIntervalSeconds)
		require.NotNil(t, config.NamespaceMode)
		require.Equal(t, "FLATTEN_NESTED_NAMESPACE", *config.NamespaceMode)
		require.NotNil(t, config.NamespaceFlattenDelimiter)
		require.Equal(t, "-", *config.NamespaceFlattenDelimiter)
		require.NotNil(t, config.AllowedWriteOperations)
		require.Equal(t, "ALL", *config.AllowedWriteOperations)
		require.NotNil(t, config.CatalogCaseSensitivity)
		require.Equal(t, "CASE_INSENSITIVE", *config.CatalogCaseSensitivity)
		require.NotNil(t, config.IsSuspended)
		require.False(t, *config.IsSuspended)
		require.Equal(t, []string{"ns1", "ns2"}, config.AllowedNamespaces)
		require.Equal(t, []string{"ns3"}, config.BlockedNamespaces)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseCatalogLinkedDatabaseConfig("not json")
		require.Error(t, err)
	})
}

func Test_getCatalogLinkStatusOptions_SQL(t *testing.T) {
	id := NewAccountObjectIdentifier("my_db")

	opts := &getCatalogLinkStatusOptions{
		arguments: &catalogLinkedDatabaseArgs{Name: id},
	}
	got, err := structToSQL(opts)
	require.NoError(t, err)
	require.Equal(t, `SELECT SYSTEM$CATALOG_LINK_STATUS ('\"my_db\"') AS "STATUS"`, got)
}

func Test_parseCatalogLinkStatus(t *testing.T) {
	t.Run("valid running output", func(t *testing.T) {
		// Output shape based on https://docs.snowflake.com/en/sql-reference/functions/system_catalog_link_status.
		raw := `{
  "executionState" : "RUNNING",
  "lastLinkAttemptStartTime" : "2024-01-01 00:00:00"
}`
		status, err := parseCatalogLinkStatus(raw)
		require.NoError(t, err)
		require.NotNil(t, status)
		require.Equal(t, "RUNNING", status.ExecutionState)
		require.Nil(t, status.FailedExecutionStateReason)
		require.Nil(t, status.FailedExecutionStateErrorCode)
		require.NotNil(t, status.LastLinkAttemptStartTime)
		require.Equal(t, "2024-01-01 00:00:00", *status.LastLinkAttemptStartTime)
		require.Empty(t, status.FailureDetails)
	})

	t.Run("valid failed output", func(t *testing.T) {
		raw := `{
  "executionState" : "FAILED",
  "failedExecutionStateReason" : "some reason",
  "failedExecutionStateErrorCode" : "391408",
  "lastLinkAttemptStartTime" : "2024-01-01 00:00:00",
  "failureDetails" : [ {
    "qualifiedEntityName" : "ns1.table1",
    "entityDomain" : "TABLE",
    "operation" : "CREATE",
    "errorCode" : "391408",
    "errorMessage" : "some error"
  } ]
}`
		status, err := parseCatalogLinkStatus(raw)
		require.NoError(t, err)
		require.NotNil(t, status)
		require.Equal(t, "FAILED", status.ExecutionState)
		require.NotNil(t, status.FailedExecutionStateReason)
		require.Equal(t, "some reason", *status.FailedExecutionStateReason)
		require.NotNil(t, status.FailedExecutionStateErrorCode)
		require.Equal(t, "391408", *status.FailedExecutionStateErrorCode)
		require.Len(t, status.FailureDetails, 1)
		require.Equal(t, "ns1.table1", status.FailureDetails[0].QualifiedEntityName)
		require.Equal(t, "TABLE", status.FailureDetails[0].EntityDomain)
		require.Equal(t, "CREATE", status.FailureDetails[0].Operation)
		require.Equal(t, "391408", status.FailureDetails[0].ErrorCode)
		require.Equal(t, "some error", status.FailureDetails[0].ErrorMessage)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseCatalogLinkStatus("not json")
		require.Error(t, err)
	})
}
