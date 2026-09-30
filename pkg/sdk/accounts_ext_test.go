package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	id := accountsTestIdAccountObjectIdentifier
	adminPassword := random.Password()
	adminRsaKey := random.Password()

	policyId := randomSchemaObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()

	warehouseId := randomAccountObjectIdentifier()
	networkPolicyId := randomAccountObjectIdentifier()
	externalVolumeId := randomAccountObjectIdentifier()
	eventTableId := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()

	tagId1 := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
	unsetTagId := randomSchemaObjectIdentifier()

	accountsTests.Create.
		withDefaultOpts(func() *CreateAccountOptions {
			return &CreateAccountOptions{
				name:          id,
				AdminName:     "someadmin",
				AdminPassword: new(adminPassword),
				Email:         "admin@example.com",
				Edition:       AccountEditionBusinessCritical,
			}
		}).
		withExpectedSqlf(
			case_Accounts_sql_Create_basic,
			`CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = '%s' EMAIL = 'admin@example.com' EDITION = BUSINESS_CRITICAL`,
			id.FullyQualifiedName(), adminPassword,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Create_all,
			func(opts *CreateAccountOptions) {
				opts.AdminPassword = nil
				opts.AdminRsaPublicKey = new(adminRsaKey)
				opts.AdminUserType = new(UserTypeService)
				opts.FirstName = new("Ad")
				opts.LastName = new("Min")
				opts.MustChangePassword = new(true)
				opts.RegionGroup = new("groupid")
				opts.Region = new("regionid")
				opts.Comment = new("Test account")
				opts.ConsumptionBillingEntity = new("be-name")
				opts.Polaris = new(true)
			},
			`CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_RSA_PUBLIC_KEY = '%s' ADMIN_USER_TYPE = SERVICE FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = true EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account' CONSUMPTION_BILLING_ENTITY = "be-name" POLARIS = true`,
			id.FullyQualifiedName(), adminRsaKey,
		).
		withAdditionalSqlCasef(
			"sql_Create_staticPassword",
			func(opts *CreateAccountOptions) {
				opts.FirstName = new("Ad")
				opts.LastName = new("Min")
				opts.MustChangePassword = new(false)
				opts.RegionGroup = new("groupid")
				opts.Region = new("regionid")
				opts.Comment = new("Test account")
			},
			`CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = '%s' FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = false EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account'`,
			id.FullyQualifiedName(), adminPassword,
		)

	accountsTests.Alter.
		withDefaultOpts(func() *AlterAccountOptions {
			return &AlterAccountOptions{}
		}).
		withModify(
			case_Accounts_validation_Alter_opts_Set_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{
					PasswordPolicy:          new(randomSchemaObjectIdentifier()),
					SessionPolicySet:        &AccountSessionPolicySet{SessionPolicy: new(randomSchemaObjectIdentifier())},
					AuthenticationPolicySet: &AccountAuthenticationPolicySet{AuthenticationPolicy: new(randomSchemaObjectIdentifier())},
				}
			},
		).
		withModify(
			case_Accounts_validation_Alter_opts_Unset_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{
					PasswordPolicy:            new(true),
					SessionPolicyUnset:        &AccountSessionPolicyUnset{SessionPolicy: new(true)},
					AuthenticationPolicyUnset: &AccountAuthenticationPolicyUnset{AuthenticationPolicy: new(true)},
				}
			},
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_ConsumptionBillingEntity_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{ConsumptionBillingEntity: new("be-name")}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_Unset_ConsumptionBillingEntity_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{ConsumptionBillingEntity: new(true)}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_Drop_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Drop = &AccountDrop{OldUrl: new(true)}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_RenameTo_requiresName",
			func(opts *AlterAccountOptions) {
				opts.RenameTo = &renameTarget
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_Force_requiresPolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{
					ConsumptionBillingEntity: new("my_consumption_billing_entity"),
					Force:                    new(true),
				}
			},
			NewError("force can only be set with PackagesPolicy, PasswordPolicy, SessionPolicy, AuthenticationPolicy, or FeaturePolicy"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Drop_ExactlyOneValueSet_NoneSet",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Drop = &AccountDrop{}
			},
			errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Drop_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Drop = &AccountDrop{OldUrl: new(true), OldOrganizationUrl: new(true)}
			},
			errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_Set,
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{ResourceMonitor: new(NewAccountObjectIdentifier("mymonitor"))}
			},
			`ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_Unset,
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{ResourceMonitor: new(true)}
			},
			`ALTER ACCOUNT UNSET RESOURCE_MONITOR`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_SetTag,
			func(opts *AlterAccountOptions) {
				opts.SetTag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`ALTER ACCOUNT SET TAG %s = 'v1', %s = 'v2'`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_UnsetTag,
			func(opts *AlterAccountOptions) {
				opts.UnsetTag = []ObjectIdentifier{unsetTagId}
			},
			`ALTER ACCOUNT UNSET TAG %s`, unsetTagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_Drop,
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Drop = &AccountDrop{OldUrl: new(true)}
			},
			`ALTER ACCOUNT %s DROP OLD URL`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_RenameTo,
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.RenameTo = &renameTarget
				opts.SaveOldURL = new(false)
			},
			`ALTER ACCOUNT %s RENAME TO %s SAVE_OLD_URL = false`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		)

	accountsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_LegacyParameters",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{
					LegacyParameters: &AccountLevelParameters{
						AccountParameters: &LegacyAccountParameters{
							ClientEncryptionKeySize:       new(128),
							PreventUnloadToInternalStages: new(true),
						},
						SessionParameters: &SessionParameters{
							JsonIndent: new(16),
						},
						ObjectParameters: &ObjectParameters{
							MaxDataExtensionTimeInDays: new(30),
						},
					},
				}
			},
			`ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 128, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, JSON_INDENT = 16, MAX_DATA_EXTENSION_TIME_IN_DAYS = 30`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_FeaturePolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{FeaturePolicySet: &AccountFeaturePolicySet{FeaturePolicy: &policyId}}
			},
			`ALTER ACCOUNT SET FEATURE POLICY %s FOR ALL APPLICATIONS`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_FeaturePolicy_Force",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{FeaturePolicySet: &AccountFeaturePolicySet{FeaturePolicy: &policyId}, Force: new(true)}
			},
			`ALTER ACCOUNT SET FEATURE POLICY %s FOR ALL APPLICATIONS FORCE`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PackagesPolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{PackagesPolicy: &policyId}
			},
			`ALTER ACCOUNT SET PACKAGES POLICY %s`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PackagesPolicy_Force",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{PackagesPolicy: &policyId, Force: new(true)}
			},
			`ALTER ACCOUNT SET PACKAGES POLICY %s FORCE`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PasswordPolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{PasswordPolicy: &policyId}
			},
			`ALTER ACCOUNT SET PASSWORD POLICY %s`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PasswordPolicy_Force",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{PasswordPolicy: &policyId, Force: new(true)}
			},
			`ALTER ACCOUNT SET PASSWORD POLICY %s FORCE`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{SessionPolicySet: &AccountSessionPolicySet{SessionPolicy: &policyId}}
			},
			`ALTER ACCOUNT SET SESSION POLICY %s`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy_Force",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{SessionPolicySet: &AccountSessionPolicySet{SessionPolicy: &policyId}, Force: new(true)}
			},
			`ALTER ACCOUNT SET SESSION POLICY %s FORCE`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy_ForAllPersonUsers",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{SessionPolicySet: &AccountSessionPolicySet{SessionPolicy: &policyId, ForAllPersonUsers: new(true)}}
			},
			`ALTER ACCOUNT SET SESSION POLICY %s FOR ALL PERSON USERS`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy_ForAllServiceUsers",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{SessionPolicySet: &AccountSessionPolicySet{SessionPolicy: &policyId, ForAllServiceUsers: new(true)}}
			},
			`ALTER ACCOUNT SET SESSION POLICY %s FOR ALL SERVICE USERS`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{AuthenticationPolicySet: &AccountAuthenticationPolicySet{AuthenticationPolicy: &policyId}}
			},
			`ALTER ACCOUNT SET AUTHENTICATION POLICY %s`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy_Force",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{AuthenticationPolicySet: &AccountAuthenticationPolicySet{AuthenticationPolicy: &policyId}, Force: new(true)}
			},
			`ALTER ACCOUNT SET AUTHENTICATION POLICY %s FORCE`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy_ForAllPersonUsers",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{AuthenticationPolicySet: &AccountAuthenticationPolicySet{AuthenticationPolicy: &policyId, ForAllPersonUsers: new(true)}}
			},
			`ALTER ACCOUNT SET AUTHENTICATION POLICY %s FOR ALL PERSON USERS`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy_ForAllServiceUsers",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{AuthenticationPolicySet: &AccountAuthenticationPolicySet{AuthenticationPolicy: &policyId, ForAllServiceUsers: new(true)}}
			},
			`ALTER ACCOUNT SET AUTHENTICATION POLICY %s FOR ALL SERVICE USERS`, policyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_ConsumptionBillingEntity",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Set = &AccountSet{ConsumptionBillingEntity: new("my_consumption_billing_entity")}
			},
			`ALTER ACCOUNT %s SET CONSUMPTION_BILLING_ENTITY = "my_consumption_billing_entity"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_OrgAdmin",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Set = &AccountSet{OrgAdmin: new(true)}
			},
			`ALTER ACCOUNT %s SET IS_ORG_ADMIN = true`, id.FullyQualifiedName(),
		)

	accountsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_Parameters",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{
					Parameters: &AccountParameters{
						AbortDetachedQuery:                               new(true),
						ActivePythonProfiler:                             new(ActivePythonProfilerMemory),
						AllowClientMfaCaching:                            new(true),
						AllowIdToken:                                     new(true),
						Autocommit:                                       new(false),
						BaseLocationPrefix:                               new("STORAGE_BASE_URL/"),
						BinaryInputFormat:                                new(BinaryInputFormatBase64),
						BinaryOutputFormat:                               new(BinaryOutputFormatBase64),
						Catalog:                                          new(NewAccountObjectIdentifier("SNOWFLAKE")),
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
						DefaultDdlCollation:                              new(StringAllowEmpty{Value: "en-cs"}),
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
						ExternalOauthAddPrivilegedRolesToBlockedList:     new(false),
						ExternalVolume:                                   new(externalVolumeId),
						GeographyOutputFormat:                            new(GeographyOutputFormatWKT),
						GeometryOutputFormat:                             new(GeometryOutputFormatWKT),
						HybridTableLockTimeout:                           new(3599),
						InitialReplicationSizeLimitInTb:                  new("9.9"),
						JdbcTreatDecimalAsInt:                            new(false),
						JdbcTreatTimestampNtzAsUtc:                       new(true),
						JdbcUseSessionTimezone:                           new(false),
						JsonIndent:                                       new(4),
						JsTreatIntegerAsBigint:                           new(true),
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
						OauthAddPrivilegedRolesToBlockedList:             new(false),
						OdbcTreatDecimalAsInt:                            new(true),
						PeriodicDataRekeying:                             new(false),
						PipeExecutionPaused:                              new(true),
						PreventUnloadToInlineUrl:                         new(true),
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
						TimestampDayIsAlways24H:                          new(true),
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
			`ALTER ACCOUNT SET ABORT_DETACHED_QUERY = true, ACTIVE_PYTHON_PROFILER = 'MEMORY', ALLOW_CLIENT_MFA_CACHING = true, ALLOW_ID_TOKEN = true, AUTOCOMMIT = false, BASE_LOCATION_PREFIX = 'STORAGE_BASE_URL/', BINARY_INPUT_FORMAT = 'BASE64', BINARY_OUTPUT_FORMAT = 'BASE64', CATALOG = "SNOWFLAKE", CATALOG_SYNC = 'CATALOG_SYNC', CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS = true, CLIENT_ENCRYPTION_KEY_SIZE = 256, CLIENT_MEMORY_LIMIT = 1540, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX = true, CLIENT_METADATA_USE_SESSION_DATABASE = true, CLIENT_PREFETCH_THREADS = 5, CLIENT_RESULT_CHUNK_SIZE = 159, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE = true, CLIENT_SESSION_KEEP_ALIVE = true, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY = 3599, CLIENT_TIMESTAMP_TYPE_MAPPING = 'TIMESTAMP_NTZ', CORTEX_ENABLED_CROSS_REGION = 'ANY_REGION', CORTEX_MODELS_ALLOWLIST = 'All', CSV_TIMESTAMP_FORMAT = 'YYYY-MM-DD', DATA_RETENTION_TIME_IN_DAYS = 2, DATE_INPUT_FORMAT = 'YYYY-MM-DD', DATE_OUTPUT_FORMAT = 'YYYY-MM-DD', DEFAULT_DDL_COLLATION = 'en-cs', DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S', DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S', DEFAULT_NULL_ORDERING = 'FIRST', DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE = %[1]s, DISABLE_UI_DOWNLOAD_BUTTON = true, DISABLE_USER_PRIVILEGE_GRANTS = true, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG = false, ENABLE_EGRESS_COST_OPTIMIZER = false, ENABLE_IDENTIFIER_FIRST_LOGIN = false, ENABLE_INTERNAL_STAGES_PRIVATELINK = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE = true, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING = false, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION = false, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR = true, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES = true, ERROR_ON_NONDETERMINISTIC_MERGE = false, ERROR_ON_NONDETERMINISTIC_UPDATE = true, EVENT_TABLE = %[4]s, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, EXTERNAL_VOLUME = %[3]s, GEOGRAPHY_OUTPUT_FORMAT = 'WKT', GEOMETRY_OUTPUT_FORMAT = 'WKT', HYBRID_TABLE_LOCK_TIMEOUT = 3599, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = '9.9', JDBC_TREAT_DECIMAL_AS_INT = false, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC = true, JDBC_USE_SESSION_TIMEZONE = false, JS_TREAT_INTEGER_AS_BIGINT = true, JSON_INDENT = 4, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE = '2 minutes', LOCK_TIMEOUT = 43201, LOG_LEVEL = 'INFO', MAX_CONCURRENCY_LEVEL = 7, MAX_DATA_EXTENSION_TIME_IN_DAYS = 13, METRIC_LEVEL = 'ALL', MIN_DATA_RETENTION_TIME_IN_DAYS = 1, MULTI_STATEMENT_COUNT = 0, NETWORK_POLICY = %[2]s, NOORDER_SEQUENCE_AS_DEFAULT = false, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, ODBC_TREAT_DECIMAL_AS_INT = true, PERIODIC_DATA_REKEYING = false, PIPE_EXECUTION_PAUSED = true, PREVENT_UNLOAD_TO_INLINE_URL = true, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, PYTHON_PROFILER_MODULES = 'module1, module2', PYTHON_PROFILER_TARGET_STAGE = %[5]s, QUERY_TAG = 'test-query-tag', QUOTED_IDENTIFIERS_IGNORE_CASE = true, REPLACE_INVALID_CHARACTERS = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION = true, ROWS_PER_RESULTSET = 1000, S3_STAGE_VPCE_DNS_NAME = 's3-vpce-dns-name', SEARCH_PATH = '$current, $public', SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'XLARGE', SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'SMALL', SIMULATED_DATA_SHARING_CONSUMER = 'simulated-consumer', SSO_LOGIN_PAGE = true, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1, STATEMENT_TIMEOUT_IN_SECONDS = 1, STORAGE_SERIALIZATION_POLICY = 'OPTIMIZED', STRICT_JSON_OUTPUT = true, SUSPEND_TASK_AFTER_NUM_FAILURES = 3, TASK_AUTO_RETRY_ATTEMPTS = 3, TIME_INPUT_FORMAT = 'YYYY-MM-DD', TIME_OUTPUT_FORMAT = 'YYYY-MM-DD', TIMESTAMP_DAY_IS_ALWAYS_24H = true, TIMESTAMP_INPUT_FORMAT = 'YYYY-MM-DD', TIMESTAMP_LTZ_OUTPUT_FORMAT = 'YYYY-MM-DD', TIMESTAMP_NTZ_OUTPUT_FORMAT = 'YYYY-MM-DD', TIMESTAMP_OUTPUT_FORMAT = 'YYYY-MM-DD', TIMESTAMP_TYPE_MAPPING = 'TIMESTAMP_LTZ', TIMESTAMP_TZ_OUTPUT_FORMAT = 'YYYY-MM-DD', TIMEZONE = 'Europe/London', TRACE_LEVEL = 'PROPAGATE', TRANSACTION_ABORT_ON_ERROR = true, TRANSACTION_DEFAULT_ISOLATION_LEVEL = 'READ COMMITTED', TWO_DIGIT_CENTURY_START = 1971, UNSUPPORTED_DDL_ACTION = 'FAIL', USE_CACHED_RESULT = false, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'SMALL', USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10, USER_TASK_TIMEOUT_MS = 10, WEEK_OF_YEAR_POLICY = 1, WEEK_START = 1`,
			warehouseId.FullyQualifiedName(),
			networkPolicyId.FullyQualifiedName(),
			externalVolumeId.FullyQualifiedName(),
			eventTableId.FullyQualifiedName(),
			stageId.FullyQualifiedName(),
		)

	accountsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Drop_OldOrganizationUrl",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Drop = &AccountDrop{OldOrganizationUrl: new(true)}
			},
			`ALTER ACCOUNT %s DROP OLD ORGANIZATION URL`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_LegacyParameters",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{
					LegacyParameters: &AccountLevelParametersUnset{
						AccountParameters: &LegacyAccountParametersUnset{
							InitialReplicationSizeLimitInTB: new(true),
							SSOLoginPage:                    new(true),
						},
						SessionParameters: &SessionParametersUnset{
							SimulatedDataSharingConsumer: new(true),
							Timezone:                     new(true),
						},
						ObjectParameters: &ObjectParametersUnset{
							DefaultDDLCollation: new(true),
						},
					},
				}
			},
			`ALTER ACCOUNT UNSET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, SSO_LOGIN_PAGE, SIMULATED_DATA_SHARING_CONSUMER, TIMEZONE, DEFAULT_DDL_COLLATION`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_ConsumptionBillingEntity",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Unset = &AccountUnset{ConsumptionBillingEntity: new(true)}
			},
			`ALTER ACCOUNT %s UNSET CONSUMPTION_BILLING_ENTITY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_PackagesPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{PackagesPolicy: new(true)}
			},
			`ALTER ACCOUNT UNSET PACKAGES POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_FeaturePolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{FeaturePolicyUnset: &AccountFeaturePolicyUnset{FeaturePolicy: new(true)}}
			},
			`ALTER ACCOUNT UNSET FEATURE POLICY FOR ALL APPLICATIONS`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_PasswordPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{PasswordPolicy: new(true)}
			},
			`ALTER ACCOUNT UNSET PASSWORD POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_SessionPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{SessionPolicyUnset: &AccountSessionPolicyUnset{SessionPolicy: new(true)}}
			},
			`ALTER ACCOUNT UNSET SESSION POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_SessionPolicy_ForAllPersonUsers",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{SessionPolicyUnset: &AccountSessionPolicyUnset{SessionPolicy: new(true), ForAllPersonUsers: new(true)}}
			},
			`ALTER ACCOUNT UNSET SESSION POLICY FOR ALL PERSON USERS`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_SessionPolicy_ForAllServiceUsers",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{SessionPolicyUnset: &AccountSessionPolicyUnset{SessionPolicy: new(true), ForAllServiceUsers: new(true)}}
			},
			`ALTER ACCOUNT UNSET SESSION POLICY FOR ALL SERVICE USERS`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_AuthenticationPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{AuthenticationPolicyUnset: &AccountAuthenticationPolicyUnset{AuthenticationPolicy: new(true)}}
			},
			`ALTER ACCOUNT UNSET AUTHENTICATION POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_AuthenticationPolicy_ForAllPersonUsers",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{AuthenticationPolicyUnset: &AccountAuthenticationPolicyUnset{AuthenticationPolicy: new(true), ForAllPersonUsers: new(true)}}
			},
			`ALTER ACCOUNT UNSET AUTHENTICATION POLICY FOR ALL PERSON USERS`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_AuthenticationPolicy_ForAllServiceUsers",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{AuthenticationPolicyUnset: &AccountAuthenticationPolicyUnset{AuthenticationPolicy: new(true), ForAllServiceUsers: new(true)}}
			},
			`ALTER ACCOUNT UNSET AUTHENTICATION POLICY FOR ALL SERVICE USERS`,
		)

	accountsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Unset_Parameters",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{
					Parameters: &AccountParametersUnset{
						AbortDetachedQuery:                               new(true),
						ActivePythonProfiler:                             new(true),
						AllowClientMfaCaching:                            new(true),
						AllowIdToken:                                     new(true),
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
						DefaultDdlCollation:                              new(true),
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
						ExternalOauthAddPrivilegedRolesToBlockedList:     new(true),
						ExternalVolume:                                   new(true),
						GeographyOutputFormat:                            new(true),
						GeometryOutputFormat:                             new(true),
						HybridTableLockTimeout:                           new(true),
						InitialReplicationSizeLimitInTb:                  new(true),
						JdbcTreatDecimalAsInt:                            new(true),
						JdbcTreatTimestampNtzAsUtc:                       new(true),
						JdbcUseSessionTimezone:                           new(true),
						JsonIndent:                                       new(true),
						JsTreatIntegerAsBigint:                           new(true),
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
						OauthAddPrivilegedRolesToBlockedList:             new(true),
						OdbcTreatDecimalAsInt:                            new(true),
						PeriodicDataRekeying:                             new(true),
						PipeExecutionPaused:                              new(true),
						PreventUnloadToInlineUrl:                         new(true),
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
						TimestampDayIsAlways24H:                          new(true),
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
			`ALTER ACCOUNT UNSET ABORT_DETACHED_QUERY, ACTIVE_PYTHON_PROFILER, ALLOW_CLIENT_MFA_CACHING, ALLOW_ID_TOKEN, AUTOCOMMIT, BASE_LOCATION_PREFIX, BINARY_INPUT_FORMAT, BINARY_OUTPUT_FORMAT, CATALOG, CATALOG_SYNC, CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS, CLIENT_ENCRYPTION_KEY_SIZE, CLIENT_MEMORY_LIMIT, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX, CLIENT_METADATA_USE_SESSION_DATABASE, CLIENT_PREFETCH_THREADS, CLIENT_RESULT_CHUNK_SIZE, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE, CLIENT_SESSION_KEEP_ALIVE, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY, CLIENT_TIMESTAMP_TYPE_MAPPING, CORTEX_ENABLED_CROSS_REGION, CORTEX_MODELS_ALLOWLIST, CSV_TIMESTAMP_FORMAT, DATA_RETENTION_TIME_IN_DAYS, DATE_INPUT_FORMAT, DATE_OUTPUT_FORMAT, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, DEFAULT_NULL_ORDERING, DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE, DISABLE_UI_DOWNLOAD_BUTTON, DISABLE_USER_PRIVILEGE_GRANTS, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG, ENABLE_EGRESS_COST_OPTIMIZER, ENABLE_IDENTIFIER_FIRST_LOGIN, ENABLE_INTERNAL_STAGES_PRIVATELINK, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES, ERROR_ON_NONDETERMINISTIC_MERGE, ERROR_ON_NONDETERMINISTIC_UPDATE, EVENT_TABLE, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, EXTERNAL_VOLUME, GEOGRAPHY_OUTPUT_FORMAT, GEOMETRY_OUTPUT_FORMAT, HYBRID_TABLE_LOCK_TIMEOUT, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, JDBC_TREAT_DECIMAL_AS_INT, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC, JDBC_USE_SESSION_TIMEZONE, JS_TREAT_INTEGER_AS_BIGINT, JSON_INDENT, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE, LOCK_TIMEOUT, LOG_LEVEL, MAX_CONCURRENCY_LEVEL, MAX_DATA_EXTENSION_TIME_IN_DAYS, METRIC_LEVEL, MIN_DATA_RETENTION_TIME_IN_DAYS, MULTI_STATEMENT_COUNT, NETWORK_POLICY, NOORDER_SEQUENCE_AS_DEFAULT, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, ODBC_TREAT_DECIMAL_AS_INT, PERIODIC_DATA_REKEYING, PIPE_EXECUTION_PAUSED, PREVENT_UNLOAD_TO_INLINE_URL, PREVENT_UNLOAD_TO_INTERNAL_STAGES, PYTHON_PROFILER_MODULES, PYTHON_PROFILER_TARGET_STAGE, QUERY_TAG, QUOTED_IDENTIFIERS_IGNORE_CASE, REPLACE_INVALID_CHARACTERS, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION, ROWS_PER_RESULTSET, S3_STAGE_VPCE_DNS_NAME, SEARCH_PATH, SERVERLESS_TASK_MAX_STATEMENT_SIZE, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SIMULATED_DATA_SHARING_CONSUMER, SSO_LOGIN_PAGE, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, STATEMENT_TIMEOUT_IN_SECONDS, STORAGE_SERIALIZATION_POLICY, STRICT_JSON_OUTPUT, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, TIME_INPUT_FORMAT, TIME_OUTPUT_FORMAT, TIMESTAMP_DAY_IS_ALWAYS_24H, TIMESTAMP_INPUT_FORMAT, TIMESTAMP_LTZ_OUTPUT_FORMAT, TIMESTAMP_NTZ_OUTPUT_FORMAT, TIMESTAMP_OUTPUT_FORMAT, TIMESTAMP_TYPE_MAPPING, TIMESTAMP_TZ_OUTPUT_FORMAT, TIMEZONE, TRACE_LEVEL, TRANSACTION_ABORT_ON_ERROR, TRANSACTION_DEFAULT_ISOLATION_LEVEL, TWO_DIGIT_CENTURY_START, UNSUPPORTED_DDL_ACTION, USE_CACHED_RESULT, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, USER_TASK_TIMEOUT_MS, WEEK_OF_YEAR_POLICY, WEEK_START`,
		)

	accountsTests.Drop.
		withDefaultOpts(func() *DropAccountOptions {
			return &DropAccountOptions{
				name:              id,
				GracePeriodInDays: new(10),
			}
		}).
		withExpectedSqlf(
			case_Accounts_sql_Drop_basic,
			`DROP ACCOUNT %s GRACE_PERIOD_IN_DAYS = 10`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Drop_all,
			func(opts *DropAccountOptions) {
				opts.IfExists = new(true)
			},
			`DROP ACCOUNT IF EXISTS %s GRACE_PERIOD_IN_DAYS = 10`, id.FullyQualifiedName(),
		)

	accountsTests.Undrop.
		withExpectedSqlf(
			case_Accounts_sql_Undrop_basic,
			`UNDROP ACCOUNT %s`, id.FullyQualifiedName(),
		)

	accountsTests.Show.
		withExpectedSqlf(
			case_Accounts_sql_Show_basic,
			`SHOW ACCOUNTS`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Show_all,
			func(opts *ShowAccountOptions) {
				opts.History = new(true)
				opts.Like = &Like{Pattern: new("myaccount")}
			},
			`SHOW ACCOUNTS HISTORY LIKE 'myaccount'`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Show_Like,
			func(opts *ShowAccountOptions) {
				opts.Like = &Like{Pattern: new("myaccount")}
			},
			`SHOW ACCOUNTS LIKE 'myaccount'`,
		)
}
