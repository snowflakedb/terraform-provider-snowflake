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
				AdminPassword: String(adminPassword),
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
				opts.AdminRsaPublicKey = String(adminRsaKey)
				opts.AdminUserType = Pointer(UserTypeService)
				opts.FirstName = String("Ad")
				opts.LastName = String("Min")
				opts.MustChangePassword = Bool(true)
				opts.RegionGroup = String("groupid")
				opts.Region = String("regionid")
				opts.Comment = String("Test account")
				opts.ConsumptionBillingEntity = String("be-name")
				opts.Polaris = Bool(true)
			},
			`CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_RSA_PUBLIC_KEY = '%s' ADMIN_USER_TYPE = SERVICE FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = true EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account' CONSUMPTION_BILLING_ENTITY = "be-name" POLARIS = true`,
			id.FullyQualifiedName(), adminRsaKey,
		).
		withAdditionalSqlCasef(
			"sql_Create_staticPassword",
			func(opts *CreateAccountOptions) {
				opts.FirstName = String("Ad")
				opts.LastName = String("Min")
				opts.MustChangePassword = Bool(false)
				opts.RegionGroup = String("groupid")
				opts.Region = String("regionid")
				opts.Comment = String("Test account")
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
					PasswordPolicy:            Bool(true),
					SessionPolicyUnset:        &AccountSessionPolicyUnset{SessionPolicy: new(true)},
					AuthenticationPolicyUnset: &AccountAuthenticationPolicyUnset{AuthenticationPolicy: new(true)},
				}
			},
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_ConsumptionBillingEntity_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{ConsumptionBillingEntity: String("be-name")}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_Unset_ConsumptionBillingEntity_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{ConsumptionBillingEntity: Bool(true)}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_Drop_requiresName",
			func(opts *AlterAccountOptions) {
				opts.Drop = &AccountDrop{OldUrl: Bool(true)}
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
					ConsumptionBillingEntity: String("my_consumption_billing_entity"),
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
				opts.Drop = &AccountDrop{OldUrl: Bool(true), OldOrganizationUrl: Bool(true)}
			},
			errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_Set,
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{ResourceMonitor: Pointer(NewAccountObjectIdentifier("mymonitor"))}
			},
			`ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_Unset,
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{ResourceMonitor: Bool(true)}
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
				opts.Drop = &AccountDrop{OldUrl: Bool(true)}
			},
			`ALTER ACCOUNT %s DROP OLD URL`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Alter_RenameTo,
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.RenameTo = &renameTarget
				opts.SaveOldURL = Bool(false)
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
							ClientEncryptionKeySize:       Int(128),
							PreventUnloadToInternalStages: Bool(true),
						},
						SessionParameters: &SessionParameters{
							JsonIndent: Int(16),
						},
						ObjectParameters: &ObjectParameters{
							MaxDataExtensionTimeInDays: Int(30),
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
				opts.Set = &AccountSet{PackagesPolicy: &policyId, Force: Bool(true)}
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
				opts.Set = &AccountSet{ConsumptionBillingEntity: String("my_consumption_billing_entity")}
			},
			`ALTER ACCOUNT %s SET CONSUMPTION_BILLING_ENTITY = "my_consumption_billing_entity"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_OrgAdmin",
			func(opts *AlterAccountOptions) {
				opts.Name = &id
				opts.Set = &AccountSet{OrgAdmin: Bool(true)}
			},
			`ALTER ACCOUNT %s SET IS_ORG_ADMIN = true`, id.FullyQualifiedName(),
		)

	accountsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_Parameters",
			func(opts *AlterAccountOptions) {
				opts.Set = &AccountSet{
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
				}
			},
			`ALTER ACCOUNT SET ABORT_DETACHED_QUERY = true, ACTIVE_PYTHON_PROFILER = "MEMORY", ALLOW_CLIENT_MFA_CACHING = true, ALLOW_ID_TOKEN = true, AUTOCOMMIT = false, BASE_LOCATION_PREFIX = "STORAGE_BASE_URL/", BINARY_INPUT_FORMAT = "BASE64", BINARY_OUTPUT_FORMAT = "BASE64", CATALOG = "SNOWFLAKE", CATALOG_SYNC = "CATALOG_SYNC", CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS = true, CLIENT_ENCRYPTION_KEY_SIZE = 256, CLIENT_MEMORY_LIMIT = 1540, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX = true, CLIENT_METADATA_USE_SESSION_DATABASE = true, CLIENT_PREFETCH_THREADS = 5, CLIENT_RESULT_CHUNK_SIZE = 159, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE = true, CLIENT_SESSION_KEEP_ALIVE = true, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY = 3599, CLIENT_TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_NTZ", CORTEX_ENABLED_CROSS_REGION = "ANY_REGION", CORTEX_MODELS_ALLOWLIST = "All", CSV_TIMESTAMP_FORMAT = "YYYY-MM-DD", DATA_RETENTION_TIME_IN_DAYS = 2, DATE_INPUT_FORMAT = "YYYY-MM-DD", DATE_OUTPUT_FORMAT = "YYYY-MM-DD", DEFAULT_DDL_COLLATION = "en-cs", DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = "CPU_X64_S", DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = "GPU_NV_S", DEFAULT_NULL_ORDERING = "FIRST", DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE = %[1]s, DISABLE_UI_DOWNLOAD_BUTTON = true, DISABLE_USER_PRIVILEGE_GRANTS = true, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG = false, ENABLE_EGRESS_COST_OPTIMIZER = false, ENABLE_IDENTIFIER_FIRST_LOGIN = false, ENABLE_INTERNAL_STAGES_PRIVATELINK = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY = true, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE = true, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING = false, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION = false, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR = true, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES = true, ERROR_ON_NONDETERMINISTIC_MERGE = false, ERROR_ON_NONDETERMINISTIC_UPDATE = true, EVENT_TABLE = %[4]s, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, EXTERNAL_VOLUME = %[3]s, GEOGRAPHY_OUTPUT_FORMAT = "WKT", GEOMETRY_OUTPUT_FORMAT = "WKT", HYBRID_TABLE_LOCK_TIMEOUT = 3599, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = 9.9, JDBC_TREAT_DECIMAL_AS_INT = false, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC = true, JDBC_USE_SESSION_TIMEZONE = false, JSON_INDENT = 4, JS_TREAT_INTEGER_AS_BIGINT = true, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE = "2 minutes", LOCK_TIMEOUT = 43201, LOG_LEVEL = "INFO", MAX_CONCURRENCY_LEVEL = 7, MAX_DATA_EXTENSION_TIME_IN_DAYS = 13, METRIC_LEVEL = "ALL", MIN_DATA_RETENTION_TIME_IN_DAYS = 1, MULTI_STATEMENT_COUNT = 0, NETWORK_POLICY = %[2]s, NOORDER_SEQUENCE_AS_DEFAULT = false, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST = false, ODBC_TREAT_DECIMAL_AS_INT = true, PERIODIC_DATA_REKEYING = false, PIPE_EXECUTION_PAUSED = true, PREVENT_UNLOAD_TO_INLINE_URL = true, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, PYTHON_PROFILER_MODULES = "module1, module2", PYTHON_PROFILER_TARGET_STAGE = %[5]s, QUERY_TAG = "test-query-tag", QUOTED_IDENTIFIERS_IGNORE_CASE = true, REPLACE_INVALID_CHARACTERS = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION = true, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION = true, ROWS_PER_RESULTSET = 1000, S3_STAGE_VPCE_DNS_NAME = "s3-vpce-dns-name", SEARCH_PATH = "$current, $public", SERVERLESS_TASK_MAX_STATEMENT_SIZE = "XLARGE", SERVERLESS_TASK_MIN_STATEMENT_SIZE = "SMALL", SIMULATED_DATA_SHARING_CONSUMER = "simulated-consumer", SSO_LOGIN_PAGE = true, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1, STATEMENT_TIMEOUT_IN_SECONDS = 1, STORAGE_SERIALIZATION_POLICY = "OPTIMIZED", STRICT_JSON_OUTPUT = true, SUSPEND_TASK_AFTER_NUM_FAILURES = 3, TASK_AUTO_RETRY_ATTEMPTS = 3, TIMESTAMP_DAY_IS_ALWAYS_24H = true, TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_LTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_NTZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_OUTPUT_FORMAT = "YYYY-MM-DD", TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_LTZ", TIMESTAMP_TZ_OUTPUT_FORMAT = "YYYY-MM-DD", TIMEZONE = "Europe/London", TIME_INPUT_FORMAT = "YYYY-MM-DD", TIME_OUTPUT_FORMAT = "YYYY-MM-DD", TRACE_LEVEL = "PROPAGATE", TRANSACTION_ABORT_ON_ERROR = true, TRANSACTION_DEFAULT_ISOLATION_LEVEL = "READ COMMITTED", TWO_DIGIT_CENTURY_START = 1971, UNSUPPORTED_DDL_ACTION = "FAIL", USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = "SMALL", USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10, USER_TASK_TIMEOUT_MS = 10, USE_CACHED_RESULT = false, WEEK_OF_YEAR_POLICY = 1, WEEK_START = 1`,
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
				opts.Drop = &AccountDrop{OldOrganizationUrl: Bool(true)}
			},
			`ALTER ACCOUNT %s DROP OLD ORGANIZATION URL`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_LegacyParameters",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{
					LegacyParameters: &AccountLevelParametersUnset{
						AccountParameters: &LegacyAccountParametersUnset{
							InitialReplicationSizeLimitInTB: Bool(true),
							SSOLoginPage:                    Bool(true),
						},
						SessionParameters: &SessionParametersUnset{
							SimulatedDataSharingConsumer: Bool(true),
							Timezone:                     Bool(true),
						},
						ObjectParameters: &ObjectParametersUnset{
							DefaultDDLCollation: Bool(true),
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
				opts.Unset = &AccountUnset{ConsumptionBillingEntity: Bool(true)}
			},
			`ALTER ACCOUNT %s UNSET CONSUMPTION_BILLING_ENTITY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_PackagesPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{PackagesPolicy: Bool(true)}
			},
			`ALTER ACCOUNT UNSET PACKAGES POLICY`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_FeaturePolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{FeaturePolicyUnset: &AccountFeaturePolicyUnset{FeaturePolicy: Bool(true)}}
			},
			`ALTER ACCOUNT UNSET FEATURE POLICY FOR ALL APPLICATIONS`,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_PasswordPolicy",
			func(opts *AlterAccountOptions) {
				opts.Unset = &AccountUnset{PasswordPolicy: Bool(true)}
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
				}
			},
			`ALTER ACCOUNT UNSET ABORT_DETACHED_QUERY, ACTIVE_PYTHON_PROFILER, ALLOW_CLIENT_MFA_CACHING, ALLOW_ID_TOKEN, AUTOCOMMIT, BASE_LOCATION_PREFIX, BINARY_INPUT_FORMAT, BINARY_OUTPUT_FORMAT, CATALOG, CATALOG_SYNC, CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS, CLIENT_ENCRYPTION_KEY_SIZE, CLIENT_MEMORY_LIMIT, CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX, CLIENT_METADATA_USE_SESSION_DATABASE, CLIENT_PREFETCH_THREADS, CLIENT_RESULT_CHUNK_SIZE, CLIENT_RESULT_COLUMN_CASE_INSENSITIVE, CLIENT_SESSION_KEEP_ALIVE, CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY, CLIENT_TIMESTAMP_TYPE_MAPPING, CORTEX_ENABLED_CROSS_REGION, CORTEX_MODELS_ALLOWLIST, CSV_TIMESTAMP_FORMAT, DATA_RETENTION_TIME_IN_DAYS, DATE_INPUT_FORMAT, DATE_OUTPUT_FORMAT, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, DEFAULT_NULL_ORDERING, DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE, DISABLE_UI_DOWNLOAD_BUTTON, DISABLE_USER_PRIVILEGE_GRANTS, ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG, ENABLE_EGRESS_COST_OPTIMIZER, ENABLE_IDENTIFIER_FIRST_LOGIN, ENABLE_INTERNAL_STAGES_PRIVATELINK, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY, ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE, ENABLE_UNHANDLED_EXCEPTIONS_REPORTING, ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION, ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, ENABLE_UNREDACTED_SECURE_OBJECT_ERROR, ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES, ERROR_ON_NONDETERMINISTIC_MERGE, ERROR_ON_NONDETERMINISTIC_UPDATE, EVENT_TABLE, EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, EXTERNAL_VOLUME, GEOGRAPHY_OUTPUT_FORMAT, GEOMETRY_OUTPUT_FORMAT, HYBRID_TABLE_LOCK_TIMEOUT, INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, JDBC_TREAT_DECIMAL_AS_INT, JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC, JDBC_USE_SESSION_TIMEZONE, JSON_INDENT, JS_TREAT_INTEGER_AS_BIGINT, LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE, LOCK_TIMEOUT, LOG_LEVEL, MAX_CONCURRENCY_LEVEL, MAX_DATA_EXTENSION_TIME_IN_DAYS, METRIC_LEVEL, MIN_DATA_RETENTION_TIME_IN_DAYS, MULTI_STATEMENT_COUNT, NETWORK_POLICY, NOORDER_SEQUENCE_AS_DEFAULT, OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST, ODBC_TREAT_DECIMAL_AS_INT, PERIODIC_DATA_REKEYING, PIPE_EXECUTION_PAUSED, PREVENT_UNLOAD_TO_INLINE_URL, PREVENT_UNLOAD_TO_INTERNAL_STAGES, PYTHON_PROFILER_MODULES, PYTHON_PROFILER_TARGET_STAGE, QUERY_TAG, QUOTED_IDENTIFIERS_IGNORE_CASE, REPLACE_INVALID_CHARACTERS, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION, REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION, ROWS_PER_RESULTSET, S3_STAGE_VPCE_DNS_NAME, SEARCH_PATH, SERVERLESS_TASK_MAX_STATEMENT_SIZE, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SIMULATED_DATA_SHARING_CONSUMER, SSO_LOGIN_PAGE, STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, STATEMENT_TIMEOUT_IN_SECONDS, STORAGE_SERIALIZATION_POLICY, STRICT_JSON_OUTPUT, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, TIMESTAMP_DAY_IS_ALWAYS_24H, TIMESTAMP_INPUT_FORMAT, TIMESTAMP_LTZ_OUTPUT_FORMAT, TIMESTAMP_NTZ_OUTPUT_FORMAT, TIMESTAMP_OUTPUT_FORMAT, TIMESTAMP_TYPE_MAPPING, TIMESTAMP_TZ_OUTPUT_FORMAT, TIMEZONE, TIME_INPUT_FORMAT, TIME_OUTPUT_FORMAT, TRACE_LEVEL, TRANSACTION_ABORT_ON_ERROR, TRANSACTION_DEFAULT_ISOLATION_LEVEL, TWO_DIGIT_CENTURY_START, UNSUPPORTED_DDL_ACTION, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, USER_TASK_TIMEOUT_MS, USE_CACHED_RESULT, WEEK_OF_YEAR_POLICY, WEEK_START`,
		)

	accountsTests.Drop.
		withDefaultOpts(func() *DropAccountOptions {
			return &DropAccountOptions{
				name:              id,
				GracePeriodInDays: Int(10),
			}
		}).
		withExpectedSqlf(
			case_Accounts_sql_Drop_basic,
			`DROP ACCOUNT %s GRACE_PERIOD_IN_DAYS = 10`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Drop_all,
			func(opts *DropAccountOptions) {
				opts.IfExists = Bool(true)
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
				opts.History = Bool(true)
				opts.Like = &Like{Pattern: String("myaccount")}
			},
			`SHOW ACCOUNTS HISTORY LIKE 'myaccount'`,
		).
		withModifyAndExpectedSqlf(
			case_Accounts_sql_Show_Like,
			func(opts *ShowAccountOptions) {
				opts.Like = &Like{Pattern: String("myaccount")}
			},
			`SHOW ACCOUNTS LIKE 'myaccount'`,
		)
}
