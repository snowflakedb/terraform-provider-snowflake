package sdk

func init() {
	id := organizationAccountsTestIdAccountObjectIdentifier

	organizationAccountsTests.Create.
		withDefaultOpts(func() *CreateOrganizationAccountOptions {
			return &CreateOrganizationAccountOptions{
				name:          id,
				AdminName:     "admin_name",
				AdminPassword: new("pass"),
				Email:         "example@email.com",
				Edition:       OrganizationAccountEditionEnterprise,
			}
		}).
		withExpectedSqlf(
			case_OrganizationAccounts_sql_Create_basic,
			`CREATE ORGANIZATION ACCOUNT %s ADMIN_NAME = admin_name ADMIN_PASSWORD = 'pass' EMAIL = 'example@email.com' EDITION = ENTERPRISE`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Create_all,
			func(opts *CreateOrganizationAccountOptions) {
				opts.AdminRsaPublicKey = new("key")
				opts.FirstName = new("first_name")
				opts.LastName = new("last_name")
				opts.MustChangePassword = new(false)
				opts.RegionGroup = new("region_group")
				opts.Region = new("region")
				opts.Comment = new("comment")
			},
			`CREATE ORGANIZATION ACCOUNT %s ADMIN_NAME = admin_name ADMIN_PASSWORD = 'pass' ADMIN_RSA_PUBLIC_KEY = 'key' FIRST_NAME = 'first_name' LAST_NAME = 'last_name' EMAIL = 'example@email.com' MUST_CHANGE_PASSWORD = false EDITION = ENTERPRISE REGION_GROUP = "region_group" REGION = "region" COMMENT = 'comment'`,
			id.FullyQualifiedName(),
		)

	resourceMonitorId := randomAccountObjectIdentifier()
	passwordPolicyId := randomSchemaObjectIdentifier()
	sessionPolicyId := randomSchemaObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifier()
	renameToId := randomAccountObjectIdentifier()
	warehouseId := randomAccountObjectIdentifier()
	networkPolicyId := randomAccountObjectIdentifier()
	externalVolumeId := randomAccountObjectIdentifier()
	eventTableId := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()

	organizationAccountsTests.Alter.
		withModify(
			case_OrganizationAccounts_validation_Alter_opts_ConflictingFields_Name_SetTags,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = new(randomAccountObjectIdentifier())
				opts.SetTags = []TagAssociation{{Name: randomSchemaObjectIdentifier(), Value: "tag-value"}}
			},
		).
		withModify(
			case_OrganizationAccounts_validation_Alter_opts_ConflictingFields_Name_UnsetTags,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = new(randomAccountObjectIdentifier())
				opts.UnsetTags = []ObjectIdentifier{randomSchemaObjectIdentifier()}
			},
		).
		withModify(
			case_OrganizationAccounts_validation_Alter_opts_Set_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Set = &OrganizationAccountSet{
					ResourceMonitor: new(randomAccountObjectIdentifier()),
					Comment:         new("comment"),
				}
			},
		).
		withModify(
			case_OrganizationAccounts_validation_Alter_opts_Unset_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Unset = &OrganizationAccountUnset{
					ResourceMonitor: new(true),
					PasswordPolicy:  new(true),
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_Set,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Set = &OrganizationAccountSet{
					Parameters: &AccountParameters{
						AbortDetachedQuery:                               new(true),
						ActivePythonProfiler:                             new(ActivePythonProfilerMemory),
						AllowClientMFACaching:                            new(true),
						AllowIDToken:                                     new(true),
						Autocommit:                                       new(false),
						BaseLocationPrefix:                               new("STORAGE_BASE_URL/"),
						BinaryInputFormat:                                new(BinaryInputFormatBase64),
						BinaryOutputFormat:                               new(BinaryOutputFormatBase64),
						Catalog:                                          new("SNOWFLAKE"),
						CatalogSync:                                      new("CATALOG_SYNC"),
						ClientEnableLogInfoStatementParameters:           new(true),
						ClientEncryptionKeySize:                          new(256),
						ClientMemoryLimit:                                new(1540),
						ClientMetadataRequestUseConnectionCtx:            new(true),
						ClientMetadataUseSessionDatabase:                 new(true),
						ClientPrefetchThreads:                            new(5),
						ClientResultChunkSize:                            new(159),
						ClientResultColumnCaseInsensitive:                new(true),
						ClientSessionKeepAlive:                           new(true),
						ClientSessionKeepAliveHeartbeatFrequency:         new(3599),
						ClientTimestampTypeMapping:                       new(ClientTimestampTypeMappingNtz),
						CortexEnabledCrossRegion:                         new("ANY_REGION"),
						CortexModelsAllowlist:                            new("All"),
						CsvTimestampFormat:                               new("YYYY-MM-DD"),
						DataRetentionTimeInDays:                          new(2),
						DateInputFormat:                                  new("YYYY-MM-DD"),
						DateOutputFormat:                                 new("YYYY-MM-DD"),
						DefaultDDLCollation:                              new("en-cs"),
						DefaultNotebookComputePoolCpu:                    new("CPU_X64_S"),
						DefaultNotebookComputePoolGpu:                    new("GPU_NV_S"),
						DefaultNullOrdering:                              new(DefaultNullOrderingFirst),
						DefaultStreamlitNotebookWarehouse:                new(warehouseId),
						DisableUiDownloadButton:                          new(true),
						DisableUserPrivilegeGrants:                       new(true),
						EnableAutomaticSensitiveDataClassificationLog:    new(false),
						EnableEgressCostOptimizer:                        new(false),
						EnableIdentifierFirstLogin:                       new(false),
						EnableInternalStagesPrivatelink:                  new(true),
						EnableTriSecretAndRekeyOptOutForImageRepository:  new(true),
						EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: new(true),
						EnableUnhandledExceptionsReporting:               new(false),
						EnableUnloadPhysicalTypeOptimization:             new(false),
						EnableUnredactedQuerySyntaxError:                 new(true),
						EnableUnredactedSecureObjectError:                new(true),
						EnforceNetworkRulesForInternalStages:             new(true),
						ErrorOnNondeterministicMerge:                     new(false),
						ErrorOnNondeterministicUpdate:                    new(true),
						EventTable:                                       new(eventTableId),
						ExternalOAuthAddPrivilegedRolesToBlockedList:     new(false),
						ExternalVolume:                                   new(externalVolumeId),
						GeographyOutputFormat:                            new(GeographyOutputFormatWKT),
						GeometryOutputFormat:                             new(GeometryOutputFormatWKT),
						HybridTableLockTimeout:                           new(3599),
						InitialReplicationSizeLimitInTB:                  new("9.9"),
						JdbcTreatDecimalAsInt:                            new(false),
						JdbcTreatTimestampNtzAsUtc:                       new(true),
						JdbcUseSessionTimezone:                           new(false),
						JsonIndent:                                       new(4),
						JsTreatIntegerAsBigInt:                           new(true),
						ListingAutoFulfillmentReplicationRefreshSchedule: new("2 minutes"),
						LockTimeout:                                      new(43201),
						LogLevel:                                         new(LogLevelInfo),
						MaxConcurrencyLevel:                              new(7),
						MaxDataExtensionTimeInDays:                       new(13),
						MetricLevel:                                      new(MetricLevelAll),
						MinDataRetentionTimeInDays:                       new(1),
						MultiStatementCount:                              new(0),
						NetworkPolicy:                                    new(networkPolicyId),
						NoorderSequenceAsDefault:                         new(false),
						OAuthAddPrivilegedRolesToBlockedList:             new(false),
						OdbcTreatDecimalAsInt:                            new(true),
						PeriodicDataRekeying:                             new(false),
						PipeExecutionPaused:                              new(true),
						PreventUnloadToInlineURL:                         new(true),
						PreventUnloadToInternalStages:                    new(true),
						PythonProfilerModules:                            new("module1, module2"),
						PythonProfilerTargetStage:                        new(stageId),
						QueryTag:                                         new("test-query-tag"),
						QuotedIdentifiersIgnoreCase:                      new(true),
						ReplaceInvalidCharacters:                         new(true),
						RequireStorageIntegrationForStageCreation:        new(true),
						RequireStorageIntegrationForStageOperation:       new(true),
						RowsPerResultset:                                 new(1000),
						S3StageVpceDnsName:                               new("s3-vpce-dns-name"),
						SearchPath:                                       new("$current, $public"),
						ServerlessTaskMaxStatementSize:                   new(WarehouseSizeXLarge),
						ServerlessTaskMinStatementSize:                   new(WarehouseSizeSmall),
						SimulatedDataSharingConsumer:                     new("simulated-consumer"),
						SsoLoginPage:                                     new(true),
						StatementQueuedTimeoutInSeconds:                  new(1),
						StatementTimeoutInSeconds:                        new(1),
						StorageSerializationPolicy:                       new(StorageSerializationPolicyOptimized),
						StrictJsonOutput:                                 new(true),
						SuspendTaskAfterNumFailures:                      new(3),
						TaskAutoRetryAttempts:                            new(3),
						TimestampDayIsAlways24h:                          new(true),
						TimestampInputFormat:                             new("YYYY-MM-DD"),
						TimestampLtzOutputFormat:                         new("YYYY-MM-DD"),
						TimestampNtzOutputFormat:                         new("YYYY-MM-DD"),
						TimestampOutputFormat:                            new("YYYY-MM-DD"),
						TimestampTypeMapping:                             new(TimestampTypeMappingLtz),
						TimestampTzOutputFormat:                          new("YYYY-MM-DD"),
						Timezone:                                         new("Europe/London"),
						TimeInputFormat:                                  new("YYYY-MM-DD"),
						TimeOutputFormat:                                 new("YYYY-MM-DD"),
						TraceLevel:                                       new(TraceLevelPropagate),
						TransactionAbortOnError:                          new(true),
						TransactionDefaultIsolationLevel:                 new(TransactionDefaultIsolationLevelReadCommitted),
						TwoDigitCenturyStart:                             new(1971),
						UnsupportedDdlAction:                             new(UnsupportedDDLActionFail),
						UserTaskManagedInitialWarehouseSize:              new(WarehouseSizeSmall),
						UserTaskMinimumTriggerIntervalInSeconds:          new(10),
						UserTaskTimeoutMs:                                new(10),
						UseCachedResult:                                  new(false),
						WeekOfYearPolicy:                                 new(1),
						WeekStart:                                        new(1),
					},
				}
			},
			`ALTER ORGANIZATION ACCOUNT SET ABORT_DETACHED_QUERY = true, ACTIVE_PYTHON_PROFILER = "MEMORY", ALLOW_CLIENT_MFA_CACHING = true, ALLOW_ID_TOKEN = true, AUTOCOMMIT = false, BASE_LOCATION_PREFIX = "STORAGE_BASE_URL/", BINARY_INPUT_FORMAT = "BASE64", BINARY_OUTPUT_FORMAT = "BASE64", CATALOG = "SNOWFLAKE", CATALOG_SYNC = "CATALOG_SYNC", CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS = true, CLIENT_ENCRYPTION_KEY_SIZE = 256, CLIENT_MEMORY_LIMIT = 1540, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX = true, CLIENT_METADATA_USE_SESSION_DATABASE = true, CLIENT_PREFETCH_THREADS = 5, CLIENT_RESULT_CHUNK_SIZE = 159, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE = true, CLIENT_SESSION_KEEP_ALIVE = true, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY = 3599, CLIENT_TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_NTZ", CORTEX_ENABLED_CROSS_REGION = "ANY_REGION", CORTEX_MODELS_ALLOWLIST = "All", CSV_TIMESTAMP_FORMAT = "YYYY-MM-DD", DATA_RETENTION_TIME_IN_DAYS = 2, DATE_INPUT_FORMAT = "YYYY-MM-DD", DATE_OUTPUT_FORMAT = "YYYY-MM-DD", DEFAULT_DDL_COLLATION = "en-cs", DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = "CPU_X64_S", DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = "GPU_NV_S", DEFAULT_NULL_ORDERING = "FIRST", DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE = %[1]s, DISABLE_UI_DOWNLOAD_BUTTON = true, DISABLE_USER_PRIVILEGE_GRANTS = true, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG = false, ENABLE_EGRESS_COST_OPTIMIZER = false, ENABLE_IDENTIFIER_FIRST_LOGIN = false, ENABLE_INTERNAL_STAGES_PRIVATELINK = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE = true, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING = false, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION = false, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR = true, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES = true, ERROR_ON_NONDETERMINISTIC_MERGE = false, ERROR_ON_NONDETERMINISTIC_UPDATE = true, EVENT_TABLE = %[4]s, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, EXTERNAL_VOLUME = %[3]s, GEOGRAPHY_OUTPUT_FORMAT = "WKT", GEOMETRY_OUTPUT_FORMAT = "WKT", HYBRID_TABLE_LOCK_TIMEOUT = 3599, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = 9.9, JDBC_TREAT_DECIMAL_AS_INT = false, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC = true, JDBC_USE_SESSION_TIMEZONE = false, JSON_INDENT = 4, JS_TREAT_INTEGER_AS_BIGINT = true, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE = "2 minutes", LOCK_TIMEOUT = 43201, LOG_LEVEL = "INFO", MAX_CONCURRENCY_LEVEL = 7, MAX_DATA_EXTENSION_TIME_IN_DAYS = 13, METRIC_LEVEL = "ALL", MIN_DATA_RETENTION_TIME_IN_DAYS = 1, MULTI_STATEMENT_COUNT = 0, NETWORK_POLICY = %[2]s, NOORDER_SEQUENCE_AS_DEFAULT = false, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, ODBC_TREAT_DECIMAL_AS_INT = true, PERIODIC_DATA_REKEYING = false, PIPE_EXECUTION_PAUSED = true, PREVENT_UNLOAD_TO_INLINE_URL = true, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, PYTHON_PROFILER_MODULES = "module1, module2", PYTHON_PROFILER_TARGET_STAGE = %[5]s, QUERY_TAG = "test-query-tag", QUOTED_IDENTIFIERS_IGNORE_CASE = true, REPLACE_INVALID_CHARACTERS = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION = true, ROWS_PER_RESULTSET = 1000, S3_STAGE_VPCE_DNS_NAME = "s3-vpce-dns-name", SEARCH_PATH = "$current, $public", SERVERLESS_TASK_MAX_STATEMENT_SIZE = "XLARGE", SERVERLESS_TASK_MIN_STATEMENT_SIZE = "SMALL", SIMULATED_DATA_SHARING_CONSUMER = "simulated-consumer", SSO_LOGIN_PAGE = true, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1, STATEMENT_TIMEOUT_IN_SECONDS = 1, STORAGE_SERIALIZATION_POLICY = "OPTIMIZED", STRICT_JSON_OUTPUT = true, SUSPEND_TASK_AFTER_NUM_FAILURES = 3, TASK_AUTO_RETRY_ATTEMPTS = 3, TIMESTAMP_DAY_IS_ALWAYS_24H = true, TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_LTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_NTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_LTZ", TIMESTAMP_TZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMEZONE = "Europe/London", TIME_INPUT_FORMAT = "YYYY-MM-DD", TIME_OUTPUT_FORMAT = "YYYY-MM-DD", TRACE_LEVEL = "PROPAGATE", TRANSACTION_ABORT_ON_ERROR = true, TRANSACTION_DEFAULT_ISOLATION_LEVEL = "READ COMMITTED", TWO_DIGIT_CENTURY_START = 1971, UNSUPPORTED_DDL_ACTION = "FAIL", USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = "SMALL", USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10, USER_TASK_TIMEOUT_MS = 10, USE_CACHED_RESULT = false, WEEK_OF_YEAR_POLICY = 1, WEEK_START = 1`,
			warehouseId.FullyQualifiedName(),
			networkPolicyId.FullyQualifiedName(),
			externalVolumeId.FullyQualifiedName(),
			eventTableId.FullyQualifiedName(),
			stageId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_resourceMonitor",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Set = &OrganizationAccountSet{ResourceMonitor: new(resourceMonitorId)}
			},
			`ALTER ORGANIZATION ACCOUNT SET RESOURCE_MONITOR = %s`, resourceMonitorId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_passwordPolicy",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Set = &OrganizationAccountSet{PasswordPolicy: new(passwordPolicyId)}
			},
			`ALTER ORGANIZATION ACCOUNT SET PASSWORD POLICY %s`, passwordPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_sessionPolicy",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Set = &OrganizationAccountSet{SessionPolicy: new(sessionPolicyId)}
			},
			`ALTER ORGANIZATION ACCOUNT SET SESSION POLICY %s`, sessionPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_Unset,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Unset = &OrganizationAccountUnset{
					Parameters: &AccountParametersUnset{
						AbortDetachedQuery:                               new(true),
						ActivePythonProfiler:                             new(true),
						AllowClientMFACaching:                            new(true),
						AllowIDToken:                                     new(true),
						Autocommit:                                       new(true),
						BaseLocationPrefix:                               new(true),
						BinaryInputFormat:                                new(true),
						BinaryOutputFormat:                               new(true),
						Catalog:                                          new(true),
						CatalogSync:                                      new(true),
						ClientEnableLogInfoStatementParameters:           new(true),
						ClientEncryptionKeySize:                          new(true),
						ClientMemoryLimit:                                new(true),
						ClientMetadataRequestUseConnectionCtx:            new(true),
						ClientMetadataUseSessionDatabase:                 new(true),
						ClientPrefetchThreads:                            new(true),
						ClientResultChunkSize:                            new(true),
						ClientResultColumnCaseInsensitive:                new(true),
						ClientSessionKeepAlive:                           new(true),
						ClientSessionKeepAliveHeartbeatFrequency:         new(true),
						ClientTimestampTypeMapping:                       new(true),
						CortexEnabledCrossRegion:                         new(true),
						CortexModelsAllowlist:                            new(true),
						CsvTimestampFormat:                               new(true),
						DataRetentionTimeInDays:                          new(true),
						DateInputFormat:                                  new(true),
						DateOutputFormat:                                 new(true),
						DefaultDDLCollation:                              new(true),
						DefaultNotebookComputePoolCpu:                    new(true),
						DefaultNotebookComputePoolGpu:                    new(true),
						DefaultNullOrdering:                              new(true),
						DefaultStreamlitNotebookWarehouse:                new(true),
						DisableUiDownloadButton:                          new(true),
						DisableUserPrivilegeGrants:                       new(true),
						EnableAutomaticSensitiveDataClassificationLog:    new(true),
						EnableEgressCostOptimizer:                        new(true),
						EnableIdentifierFirstLogin:                       new(true),
						EnableInternalStagesPrivatelink:                  new(true),
						EnableTriSecretAndRekeyOptOutForImageRepository:  new(true),
						EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: new(true),
						EnableUnhandledExceptionsReporting:               new(true),
						EnableUnloadPhysicalTypeOptimization:             new(true),
						EnableUnredactedQuerySyntaxError:                 new(true),
						EnableUnredactedSecureObjectError:                new(true),
						EnforceNetworkRulesForInternalStages:             new(true),
						ErrorOnNondeterministicMerge:                     new(true),
						ErrorOnNondeterministicUpdate:                    new(true),
						EventTable:                                       new(true),
						ExternalOAuthAddPrivilegedRolesToBlockedList:     new(true),
						ExternalVolume:                                   new(true),
						GeographyOutputFormat:                            new(true),
						GeometryOutputFormat:                             new(true),
						HybridTableLockTimeout:                           new(true),
						InitialReplicationSizeLimitInTB:                  new(true),
						JdbcTreatDecimalAsInt:                            new(true),
						JdbcTreatTimestampNtzAsUtc:                       new(true),
						JdbcUseSessionTimezone:                           new(true),
						JsonIndent:                                       new(true),
						JsTreatIntegerAsBigInt:                           new(true),
						ListingAutoFulfillmentReplicationRefreshSchedule: new(true),
						LockTimeout:                                      new(true),
						LogLevel:                                         new(true),
						MaxConcurrencyLevel:                              new(true),
						MaxDataExtensionTimeInDays:                       new(true),
						MetricLevel:                                      new(true),
						MinDataRetentionTimeInDays:                       new(true),
						MultiStatementCount:                              new(true),
						NetworkPolicy:                                    new(true),
						NoorderSequenceAsDefault:                         new(true),
						OAuthAddPrivilegedRolesToBlockedList:             new(true),
						OdbcTreatDecimalAsInt:                            new(true),
						PeriodicDataRekeying:                             new(true),
						PipeExecutionPaused:                              new(true),
						PreventUnloadToInlineURL:                         new(true),
						PreventUnloadToInternalStages:                    new(true),
						PythonProfilerModules:                            new(true),
						PythonProfilerTargetStage:                        new(true),
						QueryTag:                                         new(true),
						QuotedIdentifiersIgnoreCase:                      new(true),
						ReplaceInvalidCharacters:                         new(true),
						RequireStorageIntegrationForStageCreation:        new(true),
						RequireStorageIntegrationForStageOperation:       new(true),
						RowsPerResultset:                                 new(true),
						S3StageVpceDnsName:                               new(true),
						SearchPath:                                       new(true),
						ServerlessTaskMaxStatementSize:                   new(true),
						ServerlessTaskMinStatementSize:                   new(true),
						SimulatedDataSharingConsumer:                     new(true),
						SsoLoginPage:                                     new(true),
						StatementQueuedTimeoutInSeconds:                  new(true),
						StatementTimeoutInSeconds:                        new(true),
						StorageSerializationPolicy:                       new(true),
						StrictJsonOutput:                                 new(true),
						SuspendTaskAfterNumFailures:                      new(true),
						TaskAutoRetryAttempts:                            new(true),
						TimestampDayIsAlways24h:                          new(true),
						TimestampInputFormat:                             new(true),
						TimestampLtzOutputFormat:                         new(true),
						TimestampNtzOutputFormat:                         new(true),
						TimestampOutputFormat:                            new(true),
						TimestampTypeMapping:                             new(true),
						TimestampTzOutputFormat:                          new(true),
						Timezone:                                         new(true),
						TimeInputFormat:                                  new(true),
						TimeOutputFormat:                                 new(true),
						TraceLevel:                                       new(true),
						TransactionAbortOnError:                          new(true),
						TransactionDefaultIsolationLevel:                 new(true),
						TwoDigitCenturyStart:                             new(true),
						UnsupportedDdlAction:                             new(true),
						UserTaskManagedInitialWarehouseSize:              new(true),
						UserTaskMinimumTriggerIntervalInSeconds:          new(true),
						UserTaskTimeoutMs:                                new(true),
						UseCachedResult:                                  new(true),
						WeekOfYearPolicy:                                 new(true),
						WeekStart:                                        new(true),
					},
				}
			},
			`ALTER ORGANIZATION ACCOUNT UNSET ABORT_DETACHED_QUERY, ACTIVE_PYTHON_PROFILER, ALLOW_CLIENT_MFA_CACHING, ALLOW_ID_TOKEN, AUTOCOMMIT, BASE_LOCATION_PREFIX, BINARY_INPUT_FORMAT, BINARY_OUTPUT_FORMAT, CATALOG, CATALOG_SYNC, CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS, CLIENT_ENCRYPTION_KEY_SIZE, CLIENT_MEMORY_LIMIT, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX, CLIENT_METADATA_USE_SESSION_DATABASE, CLIENT_PREFETCH_THREADS, CLIENT_RESULT_CHUNK_SIZE, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE, CLIENT_SESSION_KEEP_ALIVE, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY, CLIENT_TIMESTAMP_TYPE_MAPPING, CORTEX_ENABLED_CROSS_REGION, CORTEX_MODELS_ALLOWLIST, CSV_TIMESTAMP_FORMAT, DATA_RETENTION_TIME_IN_DAYS, DATE_INPUT_FORMAT, DATE_OUTPUT_FORMAT, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, DEFAULT_NULL_ORDERING, DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE, DISABLE_UI_DOWNLOAD_BUTTON, DISABLE_USER_PRIVILEGE_GRANTS, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG, ENABLE_EGRESS_COST_OPTIMIZER, ENABLE_IDENTIFIER_FIRST_LOGIN, ENABLE_INTERNAL_STAGES_PRIVATELINK, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES, ERROR_ON_NONDETERMINISTIC_MERGE, ERROR_ON_NONDETERMINISTIC_UPDATE, EVENT_TABLE, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, EXTERNAL_VOLUME, GEOGRAPHY_OUTPUT_FORMAT, GEOMETRY_OUTPUT_FORMAT, HYBRID_TABLE_LOCK_TIMEOUT, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, JDBC_TREAT_DECIMAL_AS_INT, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC, JDBC_USE_SESSION_TIMEZONE, JSON_INDENT, JS_TREAT_INTEGER_AS_BIGINT, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE, LOCK_TIMEOUT, LOG_LEVEL, MAX_CONCURRENCY_LEVEL, MAX_DATA_EXTENSION_TIME_IN_DAYS, METRIC_LEVEL, MIN_DATA_RETENTION_TIME_IN_DAYS, MULTI_STATEMENT_COUNT, NETWORK_POLICY, NOORDER_SEQUENCE_AS_DEFAULT, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, ODBC_TREAT_DECIMAL_AS_INT, PERIODIC_DATA_REKEYING, PIPE_EXECUTION_PAUSED, PREVENT_UNLOAD_TO_INLINE_URL, PREVENT_UNLOAD_TO_INTERNAL_STAGES, PYTHON_PROFILER_MODULES, PYTHON_PROFILER_TARGET_STAGE, QUERY_TAG, QUOTED_IDENTIFIERS_IGNORE_CASE, REPLACE_INVALID_CHARACTERS, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION, ROWS_PER_RESULTSET, S3_STAGE_VPCE_DNS_NAME, SEARCH_PATH, SERVERLESS_TASK_MAX_STATEMENT_SIZE, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SIMULATED_DATA_SHARING_CONSUMER, SSO_LOGIN_PAGE, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, STATEMENT_TIMEOUT_IN_SECONDS, STORAGE_SERIALIZATION_POLICY, STRICT_JSON_OUTPUT, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, TIMESTAMP_DAY_IS_ALWAYS_24H, TIMESTAMP_INPUT_FORMAT, TIMESTAMP_LTZ_OUTPUT_FORMAT, TIMESTAMP_NTZ_OUTPUT_FORMAT, TIMESTAMP_OUTPUT_FORMAT, TIMESTAMP_TYPE_MAPPING, TIMESTAMP_TZ_OUTPUT_FORMAT, TIMEZONE, TIME_INPUT_FORMAT, TIME_OUTPUT_FORMAT, TRACE_LEVEL, TRANSACTION_ABORT_ON_ERROR, TRANSACTION_DEFAULT_ISOLATION_LEVEL, TWO_DIGIT_CENTURY_START, UNSUPPORTED_DDL_ACTION, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, USER_TASK_TIMEOUT_MS, USE_CACHED_RESULT, WEEK_OF_YEAR_POLICY, WEEK_START`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_resourceMonitor",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Unset = &OrganizationAccountUnset{ResourceMonitor: new(true)}
			},
			`ALTER ORGANIZATION ACCOUNT UNSET RESOURCE_MONITOR`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_passwordPolicy",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Unset = &OrganizationAccountUnset{PasswordPolicy: new(true)}
			},
			`ALTER ORGANIZATION ACCOUNT UNSET PASSWORD POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_sessionPolicy",
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.Unset = &OrganizationAccountUnset{SessionPolicy: new(true)}
			},
			`ALTER ORGANIZATION ACCOUNT UNSET SESSION POLICY`,
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_SetTags,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.SetTags = []TagAssociation{
					{Name: new(tagId), Value: "tag-value"},
					{Name: new(tagId2), Value: "tag-value2"},
				}
			},
			`ALTER ORGANIZATION ACCOUNT SET TAG %s = 'tag-value', %s = 'tag-value2'`,
			tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_UnsetTags,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = nil
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER ORGANIZATION ACCOUNT UNSET TAG %s, %s`,
			tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_RenameTo,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = new(id)
				opts.RenameTo = &OrganizationAccountRename{
					RenameTo:   new(renameToId),
					SaveOldUrl: new(false),
				}
			},
			`ALTER ORGANIZATION ACCOUNT %s RENAME TO %s SAVE_OLD_URL = false`,
			id.FullyQualifiedName(), renameToId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Alter_DropOldUrl,
			func(opts *AlterOrganizationAccountOptions) {
				opts.Name = new(id)
				opts.DropOldUrl = new(true)
			},
			`ALTER ORGANIZATION ACCOUNT %s DROP OLD URL`,
			id.FullyQualifiedName(),
		)

	organizationAccountsTests.Show.
		withExpectedSqlf(
			case_OrganizationAccounts_sql_Show_basic,
			`SHOW ORGANIZATION ACCOUNTS`,
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Show_all,
			func(opts *ShowOrganizationAccountOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
			},
			`SHOW ORGANIZATION ACCOUNTS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_OrganizationAccounts_sql_Show_Like,
			func(opts *ShowOrganizationAccountOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
			},
			`SHOW ORGANIZATION ACCOUNTS LIKE 'pattern'`,
		)
}
