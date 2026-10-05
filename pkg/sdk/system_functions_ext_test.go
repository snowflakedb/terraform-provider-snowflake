package sdk

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	tagId := randomSchemaObjectIdentifier()
	objectId := randomSchemaObjectIdentifier()
	schemaObjectId := randomSchemaObjectIdentifier()
	databaseId := randomAccountObjectIdentifier()
	bundle := "2025_01"

	systemFunctionsTests.GetTag.
		withDefaultOpts(func() *GetTagOptions {
			return &GetTagOptions{
				Arguments: GetTagArguments{
					TagId:      tagId,
					ObjectId:   objectId,
					ObjectType: ObjectTypeTable,
				},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_GetTag_basic,
			`SELECT SYSTEM$GET_TAG ('%s', '%s', 'TABLE') AS "TAG"`,
			singleQuotedIdentifier(tagId), singleQuotedIdentifier(objectId),
		).
		withAdditionalValidationCase(
			"validation_GetTag_invalidObjectType",
			func(opts *GetTagOptions) { opts.Arguments.ObjectType = "SEQUENCE;" },
			NewError("invalid object type"),
		)

	systemFunctionsTests.PipeStatus.
		withDefaultOpts(func() *PipeStatusOptions {
			return &PipeStatusOptions{
				Arguments: PipeStatusArguments{Name: schemaObjectId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_PipeStatus_basic,
			`SELECT SYSTEM$PIPE_STATUS ('%s') AS "PIPE_STATUS"`,
			singleQuotedIdentifier(schemaObjectId),
		)

	systemFunctionsTests.PipeForceResume.
		withDefaultOpts(func() *PipeForceResumeOptions {
			return &PipeForceResumeOptions{
				Arguments: PipeForceResumeArguments{Name: schemaObjectId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_PipeForceResume_basic,
			`SELECT SYSTEM$PIPE_FORCE_RESUME ('%s')`,
			singleQuotedIdentifier(schemaObjectId),
		).
		withAdditionalSqlCasef(
			"sql_PipeForceResume_stalenessCheckOverride",
			func(opts *PipeForceResumeOptions) {
				opts.Arguments.Options = []ForceResumePipeOption{ForceResumePipeOptionStalenessCheckOverride}
			},
			`SELECT SYSTEM$PIPE_FORCE_RESUME ('%s', 'STALENESS_CHECK_OVERRIDE')`,
			singleQuotedIdentifier(schemaObjectId),
		).
		withAdditionalSqlCasef(
			"sql_PipeForceResume_bothOptions",
			func(opts *PipeForceResumeOptions) {
				opts.Arguments.Options = []ForceResumePipeOption{
					ForceResumePipeOptionStalenessCheckOverride,
					ForceResumePipeOptionOwnershipTransferCheckOverride,
				}
			},
			`SELECT SYSTEM$PIPE_FORCE_RESUME ('%s', 'STALENESS_CHECK_OVERRIDE, OWNERSHIP_TRANSFER_CHECK_OVERRIDE')`,
			singleQuotedIdentifier(schemaObjectId),
		)

	systemFunctionsTests.EnableBehaviorChangeBundle.
		withDefaultOpts(func() *EnableBehaviorChangeBundleOptions {
			return &EnableBehaviorChangeBundleOptions{
				Arguments: BehaviorChangeBundleArguments{Bundle: bundle},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_EnableBehaviorChangeBundle_basic,
			`SELECT SYSTEM$ENABLE_BEHAVIOR_CHANGE_BUNDLE ('%s')`,
			bundle,
		)

	systemFunctionsTests.DisableBehaviorChangeBundle.
		withDefaultOpts(func() *DisableBehaviorChangeBundleOptions {
			return &DisableBehaviorChangeBundleOptions{
				Arguments: BehaviorChangeBundleArguments{Bundle: bundle},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_DisableBehaviorChangeBundle_basic,
			`SELECT SYSTEM$DISABLE_BEHAVIOR_CHANGE_BUNDLE ('%s')`,
			bundle,
		)

	systemFunctionsTests.ShowActiveBehaviorChangeBundles.
		withExpectedSql(
			case_SystemFunctions_sql_ShowActiveBehaviorChangeBundles_basic,
			`SELECT SYSTEM$SHOW_ACTIVE_BEHAVIOR_CHANGE_BUNDLES() AS "BUNDLES"`,
		).
		withModifyAndExpectedSqlf(
			case_SystemFunctions_sql_ShowActiveBehaviorChangeBundles_all,
			func(opts *ShowActiveBehaviorChangeBundlesOptions) {},
			`SELECT SYSTEM$SHOW_ACTIVE_BEHAVIOR_CHANGE_BUNDLES() AS "BUNDLES"`,
		)

	systemFunctionsTests.BehaviorChangeBundleStatus.
		withDefaultOpts(func() *BehaviorChangeBundleStatusOptions {
			return &BehaviorChangeBundleStatusOptions{
				Arguments: BehaviorChangeBundleArguments{Bundle: bundle},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_BehaviorChangeBundleStatus_basic,
			`SELECT SYSTEM$BEHAVIOR_CHANGE_BUNDLE_STATUS ('%s') AS "STATUS"`,
			bundle,
		)

	systemFunctionsTests.GetIcebergTableInformation.
		withDefaultOpts(func() *GetIcebergTableInformationOptions {
			return &GetIcebergTableInformationOptions{
				Arguments: GetIcebergTableInformationArguments{Name: schemaObjectId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_GetIcebergTableInformation_basic,
			`SELECT SYSTEM$GET_ICEBERG_TABLE_INFORMATION ('%s') AS "ICEBERG_TABLE_INFORMATION"`,
			singleQuotedIdentifier(schemaObjectId),
		)

	systemFunctionsTests.GetClusteringInformation.
		withDefaultOpts(func() *GetClusteringInformationOptions {
			return &GetClusteringInformationOptions{
				Arguments: GetClusteringInformationArguments{Name: schemaObjectId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_GetClusteringInformation_basic,
			`SELECT SYSTEM$CLUSTERING_INFORMATION ('%s') AS "CLUSTERING_INFORMATION"`,
			singleQuotedIdentifier(schemaObjectId),
		).
		withAdditionalSqlCasef(
			"sql_GetClusteringInformation_columns",
			func(opts *GetClusteringInformationOptions) {
				opts.Arguments.Columns = []ClusteringInformationColumn{{Name: "REGION"}, {Name: "id"}}
			},
			`SELECT SYSTEM$CLUSTERING_INFORMATION ('%s', '(\"REGION\", \"id\")') AS "CLUSTERING_INFORMATION"`,
			singleQuotedIdentifier(schemaObjectId),
		)

	systemFunctionsTests.GetCatalogLinkedDatabaseConfig.
		withDefaultOpts(func() *GetCatalogLinkedDatabaseConfigOptions {
			return &GetCatalogLinkedDatabaseConfigOptions{
				Arguments: GetCatalogLinkedDatabaseConfigArguments{Name: databaseId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_GetCatalogLinkedDatabaseConfig_basic,
			`SELECT SYSTEM$GET_CATALOG_LINKED_DATABASE_CONFIG ('%s') AS "CONFIG"`,
			singleQuotedIdentifier(databaseId),
		)

	systemFunctionsTests.GetCatalogLinkStatus.
		withDefaultOpts(func() *GetCatalogLinkStatusOptions {
			return &GetCatalogLinkStatusOptions{
				Arguments: GetCatalogLinkStatusArguments{Name: databaseId},
			}
		}).
		withExpectedSqlf(
			case_SystemFunctions_sql_GetCatalogLinkStatus_basic,
			`SELECT SYSTEM$CATALOG_LINK_STATUS ('%s') AS "STATUS"`,
			singleQuotedIdentifier(databaseId),
		)
}

func singleQuotedIdentifier(id ObjectIdentifier) string {
	return strings.ReplaceAll(id.FullyQualifiedName(), `"`, `\"`)
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
		require.Equal(t, NewAccountObjectIdentifier("MY_CATALOG_INT"), config.CatalogIntegration)
		require.Nil(t, config.CatalogName)
		require.NotNil(t, config.ExternalVolume)
		require.Equal(t, NewAccountObjectIdentifier("MY_EXTERNAL_VOL"), *config.ExternalVolume)
		require.NotNil(t, config.SyncIntervalSeconds)
		require.Equal(t, 60, *config.SyncIntervalSeconds)
		require.NotNil(t, config.NamespaceMode)
		require.Equal(t, CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace, *config.NamespaceMode)
		require.NotNil(t, config.NamespaceFlattenDelimiter)
		require.Equal(t, "-", *config.NamespaceFlattenDelimiter)
		require.NotNil(t, config.AllowedWriteOperations)
		require.Equal(t, CatalogLinkedDatabaseAllowedWriteOperationsAll, *config.AllowedWriteOperations)
		require.NotNil(t, config.CatalogCaseSensitivity)
		require.Equal(t, DatabaseCatalogCaseSensitivityCaseInsensitive, *config.CatalogCaseSensitivity)
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

func Test_parseCatalogLinkStatus(t *testing.T) {
	expectedLastLinkAttemptStartTime := time.Date(2026, 8, 17, 10, 27, 15, 945000000, time.UTC)

	t.Run("valid running output", func(t *testing.T) {
		// Output shape based on https://docs.snowflake.com/en/sql-reference/functions/system_catalog_link_status.
		raw := `{
  "executionState" : "RUNNING",
  "lastLinkAttemptStartTime" : "2026-08-17T10:27:15.945Z"
}`
		status, err := parseCatalogLinkStatus(raw)
		require.NoError(t, err)
		require.NotNil(t, status)
		require.Equal(t, "RUNNING", status.ExecutionState)
		require.Nil(t, status.FailedExecutionStateReason)
		require.Nil(t, status.FailedExecutionStateErrorCode)
		require.NotNil(t, status.LastLinkAttemptStartTime)
		require.Equal(t, expectedLastLinkAttemptStartTime, *status.LastLinkAttemptStartTime)
		require.Empty(t, status.FailureDetails)
	})

	t.Run("valid failed output", func(t *testing.T) {
		raw := `{
  "executionState" : "FAILED",
  "failedExecutionStateReason" : "some reason",
  "failedExecutionStateErrorCode" : "391408",
  "lastLinkAttemptStartTime" : "2026-08-17T10:27:15.945Z",
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
		require.NotNil(t, status.LastLinkAttemptStartTime)
		require.Equal(t, expectedLastLinkAttemptStartTime, *status.LastLinkAttemptStartTime)
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
