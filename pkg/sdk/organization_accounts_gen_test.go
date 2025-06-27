package sdk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrganizationAccounts_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid CreateOrganizationAccountOptions
	defaultOpts := func() *CreateOrganizationAccountOptions {
		return &CreateOrganizationAccountOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOrganizationAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.AdminPassword opts.AdminRsaPublicKey] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("CreateOrganizationAccountOptions", "AdminPassword", "AdminRsaPublicKey"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.AdminName = "admin_name"
		opts.AdminPassword = String("pass")
		opts.Email = "example@email.com"
		opts.Edition = OrganizationAccountEditionEnterprise
		assertOptsValidAndSQLEquals(t, opts, "CREATE ORGANIZATION ACCOUNT %s ADMIN_NAME = admin_name ADMIN_PASSWORD = 'pass' EMAIL = 'example@email.com' EDITION = ENTERPRISE", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.AdminName = "admin_name"
		opts.AdminPassword = String("pass")
		opts.AdminRsaPublicKey = String("key")
		opts.FirstName = String("first_name")
		opts.LastName = String("last_name")
		opts.Email = "example@email.com"
		opts.MustChangePassword = Bool(false)
		opts.Edition = OrganizationAccountEditionEnterprise
		opts.RegionGroup = String("region_group")
		opts.Region = String("region")
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE ORGANIZATION ACCOUNT %s ADMIN_NAME = admin_name ADMIN_PASSWORD = 'pass' ADMIN_RSA_PUBLIC_KEY = 'key' FIRST_NAME = 'first_name' LAST_NAME = 'last_name' EMAIL = 'example@email.com' MUST_CHANGE_PASSWORD = false EDITION = ENTERPRISE REGION_GROUP = "region_group" REGION = "region" COMMENT = 'comment'`, id.FullyQualifiedName())
	})
}

func TestOrganizationAccounts_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid AlterOrganizationAccountOptions
	defaultOpts := func() *AlterOrganizationAccountOptions {
		return &AlterOrganizationAccountOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterOrganizationAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.Name] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.Name opts.Set]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &id
		opts.Set = new(OrganizationAccountSet)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterOrganizationAccountOptions", "Name", "Set"))
	})

	t.Run("validation: conflicting fields for [opts.Name opts.Unset]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &id
		opts.Unset = new(OrganizationAccountUnset)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterOrganizationAccountOptions", "Name", "Unset"))
	})

	t.Run("validation: conflicting fields for [opts.Name opts.SetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &id
		opts.SetTags = []TagAssociation{
			{
				Name:  randomSchemaObjectIdentifier(),
				Value: "tag-value",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterOrganizationAccountOptions", "Name", "SetTags"))
	})

	t.Run("validation: conflicting fields for [opts.Name opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &id
		opts.UnsetTags = []ObjectIdentifier{randomSchemaObjectIdentifier()}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterOrganizationAccountOptions", "Name", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Parameters opts.Set.ResourceMonitor opts.Set.PasswordPolicy opts.Set.SessionPolicy] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = new(OrganizationAccountSet)
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOrganizationAccountOptions.Set", "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Parameters opts.Unset.ResourceMonitor opts.Unset.PasswordPolicy opts.Unset.SessionPolicy] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = new(OrganizationAccountUnset)
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOrganizationAccountOptions.Unset", "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.RenameTo opts.DropOldUrl] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOrganizationAccountOptions", "Set", "Unset", "SetTags", "UnsetTags", "RenameTo", "DropOldUrl"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.RenameTo opts.DropOldUrl] should be present - two set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = new(OrganizationAccountSet)
		opts.Unset = new(OrganizationAccountUnset)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOrganizationAccountOptions", "Set", "Unset", "SetTags", "UnsetTags", "RenameTo", "DropOldUrl"))
	})

	t.Run("set params", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier()
		networkPolicyId := randomAccountObjectIdentifier()
		externalVolumeId := randomAccountObjectIdentifier()
		eventTableId := randomSchemaObjectIdentifier()
		stageId := randomSchemaObjectIdentifier()

		opts := &AlterOrganizationAccountOptions{
			Set: &OrganizationAccountSet{
				Parameters: &AccountParameters{
					AbortDetachedQuery:                               Bool(true),
					ActivePythonProfiler:                             Pointer(ActivePythonProfilerMemory),
					AllowClientMFACaching:                            Bool(true),
					AllowIDToken:                                     Bool(true),
					Autocommit:                                       Bool(false),
					BaseLocationPrefix:                               String("STORAGE_BASE_URL/"),
					BinaryInputFormat:                                Pointer(BinaryInputFormatBase64),
					BinaryOutputFormat:                               Pointer(BinaryOutputFormatBase64),
					Catalog:                                          String("SNOWFLAKE"),
					CatalogSync:                                      String("CATALOG_SYNC"),
					ClientEnableLogInfoStatementParameters:           Bool(true),
					ClientEncryptionKeySize:                          Int(256),
					ClientMemoryLimit:                                Int(1540),
					ClientMetadataRequestUseConnectionCtx:            Bool(true),
					ClientMetadataUseSessionDatabase:                 Bool(true),
					ClientPrefetchThreads:                            Int(5),
					ClientResultChunkSize:                            Int(159),
					ClientResultColumnCaseInsensitive:                Bool(true),
					ClientSessionKeepAlive:                           Bool(true),
					ClientSessionKeepAliveHeartbeatFrequency:         Int(3599),
					ClientTimestampTypeMapping:                       Pointer(ClientTimestampTypeMappingNtz),
					CortexEnabledCrossRegion:                         String("ANY_REGION"),
					CortexModelsAllowlist:                            String("All"),
					CsvTimestampFormat:                               String("YYYY-MM-DD"),
					DataRetentionTimeInDays:                          Int(2),
					DateInputFormat:                                  String("YYYY-MM-DD"),
					DateOutputFormat:                                 String("YYYY-MM-DD"),
					DefaultDDLCollation:                              String("en-cs"),
					DefaultNotebookComputePoolCpu:                    String("CPU_X64_S"),
					DefaultNotebookComputePoolGpu:                    String("GPU_NV_S"),
					DefaultNullOrdering:                              Pointer(DefaultNullOrderingFirst),
					DefaultStreamlitNotebookWarehouse:                Pointer(warehouseId),
					DisableUiDownloadButton:                          Bool(true),
					DisableUserPrivilegeGrants:                       Bool(true),
					EnableAutomaticSensitiveDataClassificationLog:    Bool(false),
					EnableEgressCostOptimizer:                        Bool(false),
					EnableIdentifierFirstLogin:                       Bool(false),
					EnableInternalStagesPrivatelink:                  Bool(true),
					EnableTriSecretAndRekeyOptOutForImageRepository:  Bool(true),
					EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: Bool(true),
					EnableUnhandledExceptionsReporting:               Bool(false),
					EnableUnloadPhysicalTypeOptimization:             Bool(false),
					EnableUnredactedQuerySyntaxError:                 Bool(true),
					EnableUnredactedSecureObjectError:                Bool(true),
					EnforceNetworkRulesForInternalStages:             Bool(true),
					ErrorOnNondeterministicMerge:                     Bool(false),
					ErrorOnNondeterministicUpdate:                    Bool(true),
					EventTable:                                       Pointer(eventTableId),
					ExternalOAuthAddPrivilegedRolesToBlockedList:     Bool(false),
					ExternalVolume:                                   Pointer(externalVolumeId),
					GeographyOutputFormat:                            Pointer(GeographyOutputFormatWKT),
					GeometryOutputFormat:                             Pointer(GeometryOutputFormatWKT),
					HybridTableLockTimeout:                           Int(3599),
					InitialReplicationSizeLimitInTB:                  String("9.9"),
					JdbcTreatDecimalAsInt:                            Bool(false),
					JdbcTreatTimestampNtzAsUtc:                       Bool(true),
					JdbcUseSessionTimezone:                           Bool(false),
					JsonIndent:                                       Int(4),
					JsTreatIntegerAsBigInt:                           Bool(true),
					ListingAutoFulfillmentReplicationRefreshSchedule: String("2 minutes"),
					LockTimeout:                                      Int(43201),
					LogLevel:                                         Pointer(LogLevelInfo),
					MaxConcurrencyLevel:                              Int(7),
					MaxDataExtensionTimeInDays:                       Int(13),
					MetricLevel:                                      Pointer(MetricLevelAll),
					MinDataRetentionTimeInDays:                       Int(1),
					MultiStatementCount:                              Int(0),
					NetworkPolicy:                                    Pointer(networkPolicyId),
					NoorderSequenceAsDefault:                         Bool(false),
					OAuthAddPrivilegedRolesToBlockedList:             Bool(false),
					OdbcTreatDecimalAsInt:                            Bool(true),
					PeriodicDataRekeying:                             Bool(false),
					PipeExecutionPaused:                              Bool(true),
					PreventUnloadToInlineURL:                         Bool(true),
					PreventUnloadToInternalStages:                    Bool(true),
					PythonProfilerModules:                            String("module1, module2"),
					PythonProfilerTargetStage:                        Pointer(stageId),
					QueryTag:                                         String("test-query-tag"),
					QuotedIdentifiersIgnoreCase:                      Bool(true),
					ReplaceInvalidCharacters:                         Bool(true),
					RequireStorageIntegrationForStageCreation:        Bool(true),
					RequireStorageIntegrationForStageOperation:       Bool(true),
					RowsPerResultset:                                 Int(1000),
					S3StageVpceDnsName:                               String("s3-vpce-dns-name"),
					SamlIdentityProvider:                             String("saml-idp"),
					SearchPath:                                       String("$current, $public"),
					ServerlessTaskMaxStatementSize:                   Pointer(WarehouseSizeXLarge),
					ServerlessTaskMinStatementSize:                   Pointer(WarehouseSizeSmall),
					SimulatedDataSharingConsumer:                     String("simulated-consumer"),
					SsoLoginPage:                                     Bool(true),
					StatementQueuedTimeoutInSeconds:                  Int(1),
					StatementTimeoutInSeconds:                        Int(1),
					StorageSerializationPolicy:                       Pointer(StorageSerializationPolicyOptimized),
					StrictJsonOutput:                                 Bool(true),
					SuspendTaskAfterNumFailures:                      Int(3),
					TaskAutoRetryAttempts:                            Int(3),
					TimestampDayIsAlways24h:                          Bool(true),
					TimestampInputFormat:                             String("YYYY-MM-DD"),
					TimestampLtzOutputFormat:                         String("YYYY-MM-DD"),
					TimestampNtzOutputFormat:                         String("YYYY-MM-DD"),
					TimestampOutputFormat:                            String("YYYY-MM-DD"),
					TimestampTypeMapping:                             Pointer(TimestampTypeMappingLtz),
					TimestampTzOutputFormat:                          String("YYYY-MM-DD"),
					Timezone:                                         String("Europe/London"),
					TimeInputFormat:                                  String("YYYY-MM-DD"),
					TimeOutputFormat:                                 String("YYYY-MM-DD"),
					TraceLevel:                                       Pointer(TraceLevelPropagate),
					TransactionAbortOnError:                          Bool(true),
					TransactionDefaultIsolationLevel:                 Pointer(TransactionDefaultIsolationLevelReadCommitted),
					TwoDigitCenturyStart:                             Int(1971),
					UnsupportedDdlAction:                             Pointer(UnsupportedDDLActionFail),
					UserTaskManagedInitialWarehouseSize:              Pointer(WarehouseSizeSmall),
					UserTaskMinimumTriggerIntervalInSeconds:          Int(10),
					UserTaskTimeoutMs:                                Int(10),
					UseCachedResult:                                  Bool(false),
					WeekOfYearPolicy:                                 Int(1),
					WeekStart:                                        Int(1),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ORGANIZATION ACCOUNT SET ABORT_DETACHED_QUERY = true, ACTIVE_PYTHON_PROFILER = "MEMORY", ALLOW_CLIENT_MFA_CACHING = true, ALLOW_ID_TOKEN = true, AUTOCOMMIT = false, BASE_LOCATION_PREFIX = "STORAGE_BASE_URL/", BINARY_INPUT_FORMAT = "BASE64", BINARY_OUTPUT_FORMAT = "BASE64", CATALOG = "SNOWFLAKE", CATALOG_SYNC = "CATALOG_SYNC", CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS = true, CLIENT_ENCRYPTION_KEY_SIZE = 256, CLIENT_MEMORY_LIMIT = 1540, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX = true, CLIENT_METADATA_USE_SESSION_DATABASE = true, CLIENT_PREFETCH_THREADS = 5, CLIENT_RESULT_CHUNK_SIZE = 159, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE = true, CLIENT_SESSION_KEEP_ALIVE = true, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY = 3599, CLIENT_TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_NTZ", CORTEX_ENABLED_CROSS_REGION = "ANY_REGION", CORTEX_MODELS_ALLOWLIST = "All", CSV_TIMESTAMP_FORMAT = "YYYY-MM-DD", DATA_RETENTION_TIME_IN_DAYS = 2, DATE_INPUT_FORMAT = "YYYY-MM-DD", DATE_OUTPUT_FORMAT = "YYYY-MM-DD", DEFAULT_DDL_COLLATION = "en-cs", DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = "CPU_X64_S", DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = "GPU_NV_S", DEFAULT_NULL_ORDERING = "FIRST", DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE = %[1]s, DISABLE_UI_DOWNLOAD_BUTTON = true, DISABLE_USER_PRIVILEGE_GRANTS = true, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG = false, ENABLE_EGRESS_COST_OPTIMIZER = false, ENABLE_IDENTIFIER_FIRST_LOGIN = false, ENABLE_INTERNAL_STAGES_PRIVATELINK = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE = true, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING = false, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION = false, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR = true, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES = true, ERROR_ON_NONDETERMINISTIC_MERGE = false, ERROR_ON_NONDETERMINISTIC_UPDATE = true, EVENT_TABLE = %[4]s, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, EXTERNAL_VOLUME = %[3]s, GEOGRAPHY_OUTPUT_FORMAT = "WKT", GEOMETRY_OUTPUT_FORMAT = "WKT", HYBRID_TABLE_LOCK_TIMEOUT = 3599, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = 9.9, JDBC_TREAT_DECIMAL_AS_INT = false, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC = true, JDBC_USE_SESSION_TIMEZONE = false, JSON_INDENT = 4, JS_TREAT_INTEGER_AS_BIGINT = true, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE = "2 minutes", LOCK_TIMEOUT = 43201, LOG_LEVEL = "INFO", MAX_CONCURRENCY_LEVEL = 7, MAX_DATA_EXTENSION_TIME_IN_DAYS = 13, METRIC_LEVEL = "ALL", MIN_DATA_RETENTION_TIME_IN_DAYS = 1, MULTI_STATEMENT_COUNT = 0, NETWORK_POLICY = %[2]s, NOORDER_SEQUENCE_AS_DEFAULT = false, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, ODBC_TREAT_DECIMAL_AS_INT = true, PERIODIC_DATA_REKEYING = false, PIPE_EXECUTION_PAUSED = true, PREVENT_UNLOAD_TO_INLINE_URL = true, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, PYTHON_PROFILER_MODULES = "module1, module2", PYTHON_PROFILER_TARGET_STAGE = %[5]s, QUERY_TAG = "test-query-tag", QUOTED_IDENTIFIERS_IGNORE_CASE = true, REPLACE_INVALID_CHARACTERS = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION = true, ROWS_PER_RESULTSET = 1000, S3_STAGE_VPCE_DNS_NAME = "s3-vpce-dns-name", SAML_IDENTITY_PROVIDER = "saml-idp", SEARCH_PATH = "$current, $public", SERVERLESS_TASK_MAX_STATEMENT_SIZE = "XLARGE", SERVERLESS_TASK_MIN_STATEMENT_SIZE = "SMALL", SIMULATED_DATA_SHARING_CONSUMER = "simulated-consumer", SSO_LOGIN_PAGE = true, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1, STATEMENT_TIMEOUT_IN_SECONDS = 1, STORAGE_SERIALIZATION_POLICY = "OPTIMIZED", STRICT_JSON_OUTPUT = true, SUSPEND_TASK_AFTER_NUM_FAILURES = 3, TASK_AUTO_RETRY_ATTEMPTS = 3, TIMESTAMP_DAY_IS_ALWAYS_24H = true, TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_LTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_NTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_LTZ", TIMESTAMP_TZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMEZONE = "Europe/London", TIME_INPUT_FORMAT = "YYYY-MM-DD", TIME_OUTPUT_FORMAT = "YYYY-MM-DD", TRACE_LEVEL = "PROPAGATE", TRANSACTION_ABORT_ON_ERROR = true, TRANSACTION_DEFAULT_ISOLATION_LEVEL = "READ COMMITTED", TWO_DIGIT_CENTURY_START = 1971, UNSUPPORTED_DDL_ACTION = "FAIL", USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = "SMALL", USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10, USER_TASK_TIMEOUT_MS = 10, USE_CACHED_RESULT = false, WEEK_OF_YEAR_POLICY = 1, WEEK_START = 1`,
			warehouseId.FullyQualifiedName(),
			networkPolicyId.FullyQualifiedName(),
			externalVolumeId.FullyQualifiedName(),
			eventTableId.FullyQualifiedName(),
			stageId.FullyQualifiedName(),
		)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterOrganizationAccountOptions{
			Unset: &OrganizationAccountUnset{
				Parameters: &AccountParametersUnset{
					AbortDetachedQuery:                               Bool(true),
					ActivePythonProfiler:                             Bool(true),
					AllowClientMFACaching:                            Bool(true),
					AllowIDToken:                                     Bool(true),
					Autocommit:                                       Bool(true),
					BaseLocationPrefix:                               Bool(true),
					BinaryInputFormat:                                Bool(true),
					BinaryOutputFormat:                               Bool(true),
					Catalog:                                          Bool(true),
					CatalogSync:                                      Bool(true),
					ClientEnableLogInfoStatementParameters:           Bool(true),
					ClientEncryptionKeySize:                          Bool(true),
					ClientMemoryLimit:                                Bool(true),
					ClientMetadataRequestUseConnectionCtx:            Bool(true),
					ClientMetadataUseSessionDatabase:                 Bool(true),
					ClientPrefetchThreads:                            Bool(true),
					ClientResultChunkSize:                            Bool(true),
					ClientResultColumnCaseInsensitive:                Bool(true),
					ClientSessionKeepAlive:                           Bool(true),
					ClientSessionKeepAliveHeartbeatFrequency:         Bool(true),
					ClientTimestampTypeMapping:                       Bool(true),
					CortexEnabledCrossRegion:                         Bool(true),
					CortexModelsAllowlist:                            Bool(true),
					CsvTimestampFormat:                               Bool(true),
					DataRetentionTimeInDays:                          Bool(true),
					DateInputFormat:                                  Bool(true),
					DateOutputFormat:                                 Bool(true),
					DefaultDDLCollation:                              Bool(true),
					DefaultNotebookComputePoolCpu:                    Bool(true),
					DefaultNotebookComputePoolGpu:                    Bool(true),
					DefaultNullOrdering:                              Bool(true),
					DefaultStreamlitNotebookWarehouse:                Bool(true),
					DisableUiDownloadButton:                          Bool(true),
					DisableUserPrivilegeGrants:                       Bool(true),
					EnableAutomaticSensitiveDataClassificationLog:    Bool(true),
					EnableEgressCostOptimizer:                        Bool(true),
					EnableIdentifierFirstLogin:                       Bool(true),
					EnableInternalStagesPrivatelink:                  Bool(true),
					EnableTriSecretAndRekeyOptOutForImageRepository:  Bool(true),
					EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: Bool(true),
					EnableUnhandledExceptionsReporting:               Bool(true),
					EnableUnloadPhysicalTypeOptimization:             Bool(true),
					EnableUnredactedQuerySyntaxError:                 Bool(true),
					EnableUnredactedSecureObjectError:                Bool(true),
					EnforceNetworkRulesForInternalStages:             Bool(true),
					ErrorOnNondeterministicMerge:                     Bool(true),
					ErrorOnNondeterministicUpdate:                    Bool(true),
					EventTable:                                       Bool(true),
					ExternalOAuthAddPrivilegedRolesToBlockedList:     Bool(true),
					ExternalVolume:                                   Bool(true),
					GeographyOutputFormat:                            Bool(true),
					GeometryOutputFormat:                             Bool(true),
					HybridTableLockTimeout:                           Bool(true),
					InitialReplicationSizeLimitInTB:                  Bool(true),
					JdbcTreatDecimalAsInt:                            Bool(true),
					JdbcTreatTimestampNtzAsUtc:                       Bool(true),
					JdbcUseSessionTimezone:                           Bool(true),
					JsonIndent:                                       Bool(true),
					JsTreatIntegerAsBigInt:                           Bool(true),
					ListingAutoFulfillmentReplicationRefreshSchedule: Bool(true),
					LockTimeout:                                      Bool(true),
					LogLevel:                                         Bool(true),
					MaxConcurrencyLevel:                              Bool(true),
					MaxDataExtensionTimeInDays:                       Bool(true),
					MetricLevel:                                      Bool(true),
					MinDataRetentionTimeInDays:                       Bool(true),
					MultiStatementCount:                              Bool(true),
					NetworkPolicy:                                    Bool(true),
					NoorderSequenceAsDefault:                         Bool(true),
					OAuthAddPrivilegedRolesToBlockedList:             Bool(true),
					OdbcTreatDecimalAsInt:                            Bool(true),
					PeriodicDataRekeying:                             Bool(true),
					PipeExecutionPaused:                              Bool(true),
					PreventUnloadToInlineURL:                         Bool(true),
					PreventUnloadToInternalStages:                    Bool(true),
					PythonProfilerModules:                            Bool(true),
					PythonProfilerTargetStage:                        Bool(true),
					QueryTag:                                         Bool(true),
					QuotedIdentifiersIgnoreCase:                      Bool(true),
					ReplaceInvalidCharacters:                         Bool(true),
					RequireStorageIntegrationForStageCreation:        Bool(true),
					RequireStorageIntegrationForStageOperation:       Bool(true),
					RowsPerResultset:                                 Bool(true),
					S3StageVpceDnsName:                               Bool(true),
					SamlIdentityProvider:                             Bool(true),
					SearchPath:                                       Bool(true),
					ServerlessTaskMaxStatementSize:                   Bool(true),
					ServerlessTaskMinStatementSize:                   Bool(true),
					SimulatedDataSharingConsumer:                     Bool(true),
					SsoLoginPage:                                     Bool(true),
					StatementQueuedTimeoutInSeconds:                  Bool(true),
					StatementTimeoutInSeconds:                        Bool(true),
					StorageSerializationPolicy:                       Bool(true),
					StrictJsonOutput:                                 Bool(true),
					SuspendTaskAfterNumFailures:                      Bool(true),
					TaskAutoRetryAttempts:                            Bool(true),
					TimestampDayIsAlways24h:                          Bool(true),
					TimestampInputFormat:                             Bool(true),
					TimestampLtzOutputFormat:                         Bool(true),
					TimestampNtzOutputFormat:                         Bool(true),
					TimestampOutputFormat:                            Bool(true),
					TimestampTypeMapping:                             Bool(true),
					TimestampTzOutputFormat:                          Bool(true),
					Timezone:                                         Bool(true),
					TimeInputFormat:                                  Bool(true),
					TimeOutputFormat:                                 Bool(true),
					TraceLevel:                                       Bool(true),
					TransactionAbortOnError:                          Bool(true),
					TransactionDefaultIsolationLevel:                 Bool(true),
					TwoDigitCenturyStart:                             Bool(true),
					UnsupportedDdlAction:                             Bool(true),
					UserTaskManagedInitialWarehouseSize:              Bool(true),
					UserTaskMinimumTriggerIntervalInSeconds:          Bool(true),
					UserTaskTimeoutMs:                                Bool(true),
					UseCachedResult:                                  Bool(true),
					WeekOfYearPolicy:                                 Bool(true),
					WeekStart:                                        Bool(true),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ORGANIZATION ACCOUNT UNSET ABORT_DETACHED_QUERY, ACTIVE_PYTHON_PROFILER, ALLOW_CLIENT_MFA_CACHING, ALLOW_ID_TOKEN, AUTOCOMMIT, BASE_LOCATION_PREFIX, BINARY_INPUT_FORMAT, BINARY_OUTPUT_FORMAT, CATALOG, CATALOG_SYNC, CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS, CLIENT_ENCRYPTION_KEY_SIZE, CLIENT_MEMORY_LIMIT, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX, CLIENT_METADATA_USE_SESSION_DATABASE, CLIENT_PREFETCH_THREADS, CLIENT_RESULT_CHUNK_SIZE, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE, CLIENT_SESSION_KEEP_ALIVE, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY, CLIENT_TIMESTAMP_TYPE_MAPPING, CORTEX_ENABLED_CROSS_REGION, CORTEX_MODELS_ALLOWLIST, CSV_TIMESTAMP_FORMAT, DATA_RETENTION_TIME_IN_DAYS, DATE_INPUT_FORMAT, DATE_OUTPUT_FORMAT, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, DEFAULT_NULL_ORDERING, DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE, DISABLE_UI_DOWNLOAD_BUTTON, DISABLE_USER_PRIVILEGE_GRANTS, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG, ENABLE_EGRESS_COST_OPTIMIZER, ENABLE_IDENTIFIER_FIRST_LOGIN, ENABLE_INTERNAL_STAGES_PRIVATELINK, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES, ERROR_ON_NONDETERMINISTIC_MERGE, ERROR_ON_NONDETERMINISTIC_UPDATE, EVENT_TABLE, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, EXTERNAL_VOLUME, GEOGRAPHY_OUTPUT_FORMAT, GEOMETRY_OUTPUT_FORMAT, HYBRID_TABLE_LOCK_TIMEOUT, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, JDBC_TREAT_DECIMAL_AS_INT, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC, JDBC_USE_SESSION_TIMEZONE, JSON_INDENT, JS_TREAT_INTEGER_AS_BIGINT, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE, LOCK_TIMEOUT, LOG_LEVEL, MAX_CONCURRENCY_LEVEL, MAX_DATA_EXTENSION_TIME_IN_DAYS, METRIC_LEVEL, MIN_DATA_RETENTION_TIME_IN_DAYS, MULTI_STATEMENT_COUNT, NETWORK_POLICY, NOORDER_SEQUENCE_AS_DEFAULT, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, ODBC_TREAT_DECIMAL_AS_INT, PERIODIC_DATA_REKEYING, PIPE_EXECUTION_PAUSED, PREVENT_UNLOAD_TO_INLINE_URL, PREVENT_UNLOAD_TO_INTERNAL_STAGES, PYTHON_PROFILER_MODULES, PYTHON_PROFILER_TARGET_STAGE, QUERY_TAG, QUOTED_IDENTIFIERS_IGNORE_CASE, REPLACE_INVALID_CHARACTERS, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION, ROWS_PER_RESULTSET, S3_STAGE_VPCE_DNS_NAME, SAML_IDENTITY_PROVIDER, SEARCH_PATH, SERVERLESS_TASK_MAX_STATEMENT_SIZE, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SIMULATED_DATA_SHARING_CONSUMER, SSO_LOGIN_PAGE, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, STATEMENT_TIMEOUT_IN_SECONDS, STORAGE_SERIALIZATION_POLICY, STRICT_JSON_OUTPUT, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, TIMESTAMP_DAY_IS_ALWAYS_24H, TIMESTAMP_INPUT_FORMAT, TIMESTAMP_LTZ_OUTPUT_FORMAT, TIMESTAMP_NTZ_OUTPUT_FORMAT, TIMESTAMP_OUTPUT_FORMAT, TIMESTAMP_TYPE_MAPPING, TIMESTAMP_TZ_OUTPUT_FORMAT, TIMEZONE, TIME_INPUT_FORMAT, TIME_OUTPUT_FORMAT, TRACE_LEVEL, TRANSACTION_ABORT_ON_ERROR, TRANSACTION_DEFAULT_ISOLATION_LEVEL, TWO_DIGIT_CENTURY_START, UNSUPPORTED_DDL_ACTION, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, USER_TASK_TIMEOUT_MS, USE_CACHED_RESULT, WEEK_OF_YEAR_POLICY, WEEK_START`)
	})

	t.Run("set resource monitor", func(t *testing.T) {
		resourceMonitorId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &OrganizationAccountSet{
			ResourceMonitor: &resourceMonitorId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT SET RESOURCE_MONITOR = %s", resourceMonitorId.FullyQualifiedName())
	})

	t.Run("unset resource monitor", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OrganizationAccountUnset{
			ResourceMonitor: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT UNSET RESOURCE_MONITOR")
	})

	t.Run("set password policy", func(t *testing.T) {
		passwordPolicyId := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &OrganizationAccountSet{
			PasswordPolicy: &passwordPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT SET PASSWORD POLICY %s", passwordPolicyId.FullyQualifiedName())
	})

	t.Run("unset password policy", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OrganizationAccountUnset{
			PasswordPolicy: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT UNSET PASSWORD POLICY")
	})

	t.Run("set session policy", func(t *testing.T) {
		sessionPolicyId := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &OrganizationAccountSet{
			SessionPolicy: &sessionPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT SET SESSION POLICY %s", sessionPolicyId.FullyQualifiedName())
	})

	t.Run("unset session policy", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OrganizationAccountUnset{
			SessionPolicy: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT UNSET SESSION POLICY")
	})

	t.Run("set tags", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  &tagId,
				Value: "tag-value",
			},
			{
				Name:  &tagId2,
				Value: "tag-value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT SET TAG %s = 'tag-value', %s = 'tag-value2'", tagId.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT UNSET TAG %s, %s", tagId.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Name = &id
		opts.RenameTo = &OrganizationAccountRename{
			NewName:    &newId,
			SaveOldUrl: Bool(false),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT %s RENAME TO %s SAVE_OLD_URL = false", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("drop old url", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = &id
		opts.DropOldUrl = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER ORGANIZATION ACCOUNT %s DROP OLD URL", id.FullyQualifiedName())
	})
}

func TestOrganizationAccounts_Show(t *testing.T) {
	// Minimal valid ShowOrganizationAccountOptions
	defaultOpts := func() *ShowOrganizationAccountOptions {
		return &ShowOrganizationAccountOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowOrganizationAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW ORGANIZATION ACCOUNTS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ORGANIZATION ACCOUNTS LIKE 'pattern'")
	})
}

func Test_Provider_ToOrganizationAccountEdition(t *testing.T) {
	type test struct {
		input string
		want  OrganizationAccountEdition
	}

	valid := []test{
		// Case insensitive.
		{input: "enterprise", want: OrganizationAccountEditionEnterprise},

		// Supported Values.
		{input: "ENTERPRISE", want: OrganizationAccountEditionEnterprise},
		{input: "BUSINESS_CRITICAL", want: OrganizationAccountEditionBusinessCritical},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToOrganizationAccountEdition(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToOrganizationAccountEdition(tc.input)
			require.Error(t, err)
		})
	}
}
