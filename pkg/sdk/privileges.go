package sdk

type GlobalPrivilege string

const (
	GlobalPrivilegeCreateAccount             GlobalPrivilege = "CREATE ACCOUNT"
	GlobalPrivilegeCreateApplication         GlobalPrivilege = "CREATE APPLICATION"
	GlobalPrivilegeCreateApplicationPackage  GlobalPrivilege = "CREATE APPLICATION PACKAGE"
	GlobalPrivilegeCreateComputePool         GlobalPrivilege = "CREATE COMPUTE POOL"
	GlobalPrivilegeCreateDataExchangeListing GlobalPrivilege = "CREATE DATA EXCHANGE LISTING"
	GlobalPrivilegeCreateDatabase            GlobalPrivilege = "CREATE DATABASE"
	GlobalPrivilegeCreateFailoverGroup       GlobalPrivilege = "CREATE FAILOVER GROUP"
	GlobalPrivilegeCreateIntegration         GlobalPrivilege = "CREATE INTEGRATION"
	GlobalPrivilegeCreateNetworkPolicy       GlobalPrivilege = "CREATE NETWORK POLICY"
	GlobalPrivilegeCreateExternalVolume      GlobalPrivilege = "CREATE EXTERNAL VOLUME"
	GlobalPrivilegeCreateReplicationGroup    GlobalPrivilege = "CREATE REPLICATION GROUP"
	GlobalPrivilegeCreateRole                GlobalPrivilege = "CREATE ROLE"
	GlobalPrivilegeCreateShare               GlobalPrivilege = "CREATE SHARE"
	GlobalPrivilegeCreateUser                GlobalPrivilege = "CREATE USER"
	GlobalPrivilegeCreateWarehouse           GlobalPrivilege = "CREATE WAREHOUSE"

	GlobalPrivilegeAttachPolicy        GlobalPrivilege = "ATTACH POLICY"
	GlobalPrivilegeAudit               GlobalPrivilege = "AUDIT"
	GlobalPrivilegeBindServiceEndpoint GlobalPrivilege = "BIND SERVICE ENDPOINT"

	GlobalPrivilegeApplyAggregationPolicy    GlobalPrivilege = "APPLY AGGREGATION POLICY"
	GlobalPrivilegeApplyAuthenticationPolicy GlobalPrivilege = "APPLY AUTHENTICATION POLICY"
	GlobalPrivilegeApplyMaskingPolicy        GlobalPrivilege = "APPLY MASKING POLICY"
	GlobalPrivilegeApplyPackagesPolicy       GlobalPrivilege = "APPLY PACKAGES POLICY"
	GlobalPrivilegeApplyPasswordPolicy       GlobalPrivilege = "APPLY PASSWORD POLICY"
	GlobalPrivilegeApplyProjectionPolicy     GlobalPrivilege = "APPLY PROJECTION POLICY"
	GlobalPrivilegeApplyRowAccessPolicy      GlobalPrivilege = "APPLY ROW ACCESS POLICY"
	GlobalPrivilegeApplySessionPolicy        GlobalPrivilege = "APPLY SESSION POLICY"
	GlobalPrivilegeApplyTag                  GlobalPrivilege = "APPLY TAG"

	GlobalPrivilegeExecuteAlert              GlobalPrivilege = "EXECUTE ALERT"
	GlobalPrivilegeExecuteDataMetricFunction GlobalPrivilege = "EXECUTE DATA METRIC FUNCTION"
	GlobalPrivilegeExecuteDataManagedAlert   GlobalPrivilege = "EXECUTE MANAGED ALERT"
	GlobalPrivilegeExecuteDataManagedTask    GlobalPrivilege = "EXECUTE MANAGED TASK"
	GlobalPrivilegeExecuteTask               GlobalPrivilege = "EXECUTE TASK"

	GlobalPrivilegeImportShare GlobalPrivilege = "IMPORT SHARE"

	GlobalPrivilegeManageAccountSupportCases      GlobalPrivilege = "MANAGE ACCOUNT SUPPORT CASES"
	GlobalPrivilegeManageEventSharing             GlobalPrivilege = "MANAGE EVENT SHARING"
	GlobalPrivilegeManageGrants                   GlobalPrivilege = "MANAGE GRANTS"
	GlobalPrivilegeManageListingAutoFulfillment   GlobalPrivilege = "MANAGE LISTING AUTO FULFILLMENT"
	GlobalPrivilegeManageOrganizationSupportCases GlobalPrivilege = "MANAGE ORGANIZATION SUPPORT CASES"
	GlobalPrivilegeManagePolarisConnections       GlobalPrivilege = "MANAGE POLARIS CONNECTIONS"
	GlobalPrivilegeManageUserSupportCases         GlobalPrivilege = "MANAGE USER SUPPORT CASES"
	GlobalPrivilegeManageWarehouses               GlobalPrivilege = "MANAGE WAREHOUSES"
	GlobalPrivilegeManageShareTarget              GlobalPrivilege = "MANAGE SHARE TARGET"

	GlobalPrivilegeModifyLogLevel          GlobalPrivilege = "MODIFY LOG LEVEL"
	GlobalPrivilegeModifyTraceLevel        GlobalPrivilege = "MODIFY TRACE LEVEL"
	GlobalPrivilegeModifySessionLogLevel   GlobalPrivilege = "MODIFY SESSION LOG LEVEL"
	GlobalPrivilegeModifySessionTraceLevel GlobalPrivilege = "MODIFY SESSION TRACE LEVEL"

	GlobalPrivilegeMonitorExecution GlobalPrivilege = "MONITOR EXECUTION"
	GlobalPrivilegeMonitorSecurity  GlobalPrivilege = "MONITOR SECURITY"
	GlobalPrivilegeMonitorUsage     GlobalPrivilege = "MONITOR USAGE"

	GlobalPrivilegeOverrideShareRestrictions   GlobalPrivilege = "OVERRIDE SHARE RESTRICTIONS"
	GlobalPrivilegePurchaseDataExchangeListing GlobalPrivilege = "PURCHASE DATA EXCHANGE LISTING"
	GlobalPrivilegeResolveAll                  GlobalPrivilege = "RESOLVE ALL"
)

func (p GlobalPrivilege) String() string {
	return string(p)
}

type AccountObjectPrivilege string

const (
	// For COMPUTE POOL
	// AccountObjectPrivilegeOperate AccountObjectPrivilege = "OPERATE" (duplicate)
	// AccountObjectPrivilegeModify  AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	// AccountObjectPrivilegeUsage   AccountObjectPrivilege = "USAGE" (duplicate)

	// For DATABASE
	AccountObjectPrivilegeApplyBudget        AccountObjectPrivilege = "APPLYBUDGET"
	AccountObjectPrivilegeCreateDatabaseRole AccountObjectPrivilege = "CREATE DATABASE ROLE"
	AccountObjectPrivilegeCreateSchema       AccountObjectPrivilege = "CREATE SCHEMA"
	AccountObjectPrivilegeImportedPrivileges AccountObjectPrivilege = "IMPORTED PRIVILEGES"
	AccountObjectPrivilegeModify             AccountObjectPrivilege = "MODIFY"
	AccountObjectPrivilegeMonitor            AccountObjectPrivilege = "MONITOR"
	AccountObjectPrivilegeUsage              AccountObjectPrivilege = "USAGE"

	// For EXTERNAL VOLUME
	// AccountObjectPrivilegeUsage AccountObjectPrivilege = "USAGE" (duplicate)

	// For FAILOVER GROUP
	AccountObjectPrivilegeFailover AccountObjectPrivilege = "FAILOVER"
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	// AccountObjectPrivilegeReplicate AccountObjectPrivilege = "REPLICATE" (duplicate)

	// For INTEGRATION
	// AccountObjectPrivilegeUsage AccountObjectPrivilege = "USAGE" (duplicate)
	AccountObjectPrivilegeUseAnyRole AccountObjectPrivilege = "USE_ANY_ROLE"

	// For REPLICATION GROUP
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)
	AccountObjectPrivilegeReplicate AccountObjectPrivilege = "REPLICATE"

	// For RESOURCE MONITOR
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor AccountObjectPrivilege = "MONITOR" (duplicate)

	// For USER
	// AccountObjectPrivilegeModify AccountObjectPrivilege = "MODIFY" (duplicate)
	AccountObjectPrivilegeImpersonate AccountObjectPrivilege = "IMPERSONATE"

	// For WAREHOUSE
	// AccountObjectPrivilegeApplyBudget AccountObjectPrivilege = "APPLYBUDGET" (duplicate)
	// AccountObjectPrivilegeModify   	 AccountObjectPrivilege = "MODIFY" (duplicate)
	// AccountObjectPrivilegeMonitor  	 AccountObjectPrivilege = "MONITOR" (duplicate)
	// AccountObjectPrivilegeUsage    	 AccountObjectPrivilege = "USAGE" (duplicate)
	AccountObjectPrivilegeOperate AccountObjectPrivilege = "OPERATE"
)

func (p AccountObjectPrivilege) String() string {
	return string(p)
}

type SchemaPrivilege string

const (
	SchemaPrivilegeAddSearchOptimization             SchemaPrivilege = "ADD SEARCH OPTIMIZATION"
	SchemaPrivilegeApplyBudget                       SchemaPrivilege = "APPLYBUDGET"
	SchemaPrivilegeCreateAlert                       SchemaPrivilege = "CREATE ALERT"
	SchemaPrivilegeCreateCortexSearchService         SchemaPrivilege = "CREATE CORTEX SEARCH SERVICE"
	SchemaPrivilegeCreateDataset                     SchemaPrivilege = "CREATE DATASET"
	SchemaPrivilegeCreateEventTable                  SchemaPrivilege = "CREATE EVENT TABLE"
	SchemaPrivilegeCreateFileFormat                  SchemaPrivilege = "CREATE FILE FORMAT"
	SchemaPrivilegeCreateFunction                    SchemaPrivilege = "CREATE FUNCTION"
	SchemaPrivilegeCreateGitRepository               SchemaPrivilege = "CREATE GIT REPOSITORY"
	SchemaPrivilegeCreateImageRepository             SchemaPrivilege = "CREATE IMAGE REPOSITORY"
	SchemaPrivilegeCreateModel                       SchemaPrivilege = "CREATE MODEL"
	SchemaPrivilegeCreateNetworkRule                 SchemaPrivilege = "CREATE NETWORK RULE"
	SchemaPrivilegeCreateNotebook                    SchemaPrivilege = "CREATE NOTEBOOK"
	SchemaPrivilegeCreatePipe                        SchemaPrivilege = "CREATE PIPE"
	SchemaPrivilegeCreateProcedure                   SchemaPrivilege = "CREATE PROCEDURE"
	SchemaPrivilegeCreateAggregationPolicy           SchemaPrivilege = "CREATE AGGREGATION POLICY"
	SchemaPrivilegeCreateAuthenticationPolicy        SchemaPrivilege = "CREATE AUTHENTICATION POLICY"
	SchemaPrivilegeCreateMaskingPolicy               SchemaPrivilege = "CREATE MASKING POLICY"
	SchemaPrivilegeCreatePackagesPolicy              SchemaPrivilege = "CREATE PACKAGES POLICY"
	SchemaPrivilegeCreatePasswordPolicy              SchemaPrivilege = "CREATE PASSWORD POLICY"
	SchemaPrivilegeCreateProjectionPolicy            SchemaPrivilege = "CREATE PROJECTION POLICY"
	SchemaPrivilegeCreateRowAccessPolicy             SchemaPrivilege = "CREATE ROW ACCESS POLICY"
	SchemaPrivilegeCreateSessionPolicy               SchemaPrivilege = "CREATE SESSION POLICY"
	SchemaPrivilegeCreateSecret                      SchemaPrivilege = "CREATE SECRET"
	SchemaPrivilegeCreateSequence                    SchemaPrivilege = "CREATE SEQUENCE"
	SchemaPrivilegeCreateService                     SchemaPrivilege = "CREATE SERVICE"
	SchemaPrivilegeCreateSnapshot                    SchemaPrivilege = "CREATE SNAPSHOT"
	SchemaPrivilegeCreateStage                       SchemaPrivilege = "CREATE STAGE"
	SchemaPrivilegeCreateStream                      SchemaPrivilege = "CREATE STREAM"
	SchemaPrivilegeCreateStreamlit                   SchemaPrivilege = "CREATE STREAMLIT"
	SchemaPrivilegeCreateSnowflakeCoreBudget         SchemaPrivilege = "CREATE SNOWFLAKE.CORE.BUDGET"
	SchemaPrivilegeCreateSnowflakeMlAnomalyDetection SchemaPrivilege = "CREATE SNOWFLAKE.ML.ANOMALY_DETECTION"
	SchemaPrivilegeCreateSnowflakeMlForecast         SchemaPrivilege = "CREATE SNOWFLAKE.ML.FORECAST"
	SchemaPrivilegeCreateDynamicTable                SchemaPrivilege = "CREATE DYNAMIC TABLE"
	SchemaPrivilegeCreateExternalTable               SchemaPrivilege = "CREATE EXTERNAL TABLE"
	SchemaPrivilegeCreateHybridTable                 SchemaPrivilege = "CREATE HYBRID TABLE"
	SchemaPrivilegeCreateIcebergTable                SchemaPrivilege = "CREATE ICEBERG TABLE"
	SchemaPrivilegeCreateTable                       SchemaPrivilege = "CREATE TABLE"
	SchemaPrivilegeCreateTag                         SchemaPrivilege = "CREATE TAG"
	SchemaPrivilegeCreateTask                        SchemaPrivilege = "CREATE TASK"
	SchemaPrivilegeCreateMaterializedView            SchemaPrivilege = "CREATE MATERIALIZED VIEW"
	SchemaPrivilegeCreateView                        SchemaPrivilege = "CREATE VIEW"
	SchemaPrivilegeModify                            SchemaPrivilege = "MODIFY"
	SchemaPrivilegeMonitor                           SchemaPrivilege = "MONITOR"
	SchemaPrivilegeUsage                             SchemaPrivilege = "USAGE"
)

func (p SchemaPrivilege) String() string {
	return string(p)
}

type SchemaObjectPrivilege string

const (
	SchemaObjectOwnership SchemaObjectPrivilege = "OWNERSHIP"

	// For ALERT
	// SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "MONITOR" (duplicate)
	SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE"

	// For APPLICATION PACKAGE
	SchemaObjectPrivilegeAttachListing  SchemaObjectPrivilege = "ATTACH LISTING"
	SchemaObjectPrivilegeDevelop        SchemaObjectPrivilege = "DEVELOP"
	SchemaObjectPrivilegeInstall        SchemaObjectPrivilege = "INSTALL"
	SchemaObjectPrivilegeManageVersions SchemaObjectPrivilege = "MANAGE VERSIONS"
	SchemaObjectPrivilegeManageReleases SchemaObjectPrivilege = "MANAGE RELEASES"

	// For DYNAMIC TABLE
	// SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "MONITOR" (duplicate)
	// SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE" (duplicate)
	// SchemaObjectPrivilegeSelect  SchemaObjectPrivilege = "SELECT" (duplicate)

	// For CORTEX SEARCH SERVICE
	// SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "USAGE" (duplicate)

	// For EVENT TABLE
	SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT"
	SchemaObjectPrivilegeInsert SchemaObjectPrivilege = "INSERT"

	// For FILE FORMAT, FUNCTION (UDF or external function), PROCEDURE, SECRET, or SEQUENCE
	SchemaObjectPrivilegeUsage SchemaObjectPrivilege = "USAGE"

	// For HYBRID TABLE
	// SchemaObjectPrivilegeInsert SchemaObjectPrivilege = "INSERT" (duplicate)
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeUpdate SchemaObjectPrivilege = "UPDATE" (duplicate)

	// For IMAGE REPOSITORY
	// SchemaObjectPrivilegeRead  SchemaObjectPrivilege = "READ" (duplicate)
	// SchemaObjectPrivilegeWrite SchemaObjectPrivilege = "WRITE" (duplicate)

	// For ICEBERG TABLE
	SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET"
	// SchemaObjectPrivilegeDelete     SchemaObjectPrivilege = "DELETE" (duplicate)
	// SchemaObjectPrivilegeInsert     SchemaObjectPrivilege = "INSERT" (duplicate)
	// SchemaObjectPrivilegeReferences SchemaObjectPrivilege = "REFERENCES" (duplicate)
	// SchemaObjectPrivilegeSelect     SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeTruncate   SchemaObjectPrivilege = "TRUNCATE" (duplicate)
	// SchemaObjectPrivilegeUpdate     SchemaObjectPrivilege = "UPDATE" (duplicate)

	// For PIPE
	// SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET" (duplicate)
	SchemaObjectPrivilegeMonitor SchemaObjectPrivilege = "MONITOR"
	// SchemaObjectPrivilegeOperate SchemaObjectPrivilege = "OPERATE" (duplicate)

	// For { MASKING | PASSWORD | ROW ACCESS | SESSION } POLICY or TAG
	SchemaObjectPrivilegeApply SchemaObjectPrivilege = "APPLY"

	// For external STAGE
	// SchemaObjectPrivilegeUsage SchemaObjectPrivilege = "USAGE" (duplicate)

	// For internal STAGE
	SchemaObjectPrivilegeRead  SchemaObjectPrivilege = "READ"
	SchemaObjectPrivilegeWrite SchemaObjectPrivilege = "WRITE"

	// For STREAM
	// SchemaObjectPrivilegeSelect SchemaObjectPrivilege = "SELECT" (duplicate)

	// For STREAMLIT
	// SchemaObjectPrivilegeUsage SchemaObjectPrivilege = "USAGE" (duplicate)

	// For TABLE
	// SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET" (duplicate)
	// SchemaObjectPrivilegeSelect 		SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeInsert 		SchemaObjectPrivilege = "INSERT" (duplicate)
	SchemaObjectPrivilegeEvolveSchema SchemaObjectPrivilege = "EVOLVE SCHEMA"
	SchemaObjectPrivilegeUpdate       SchemaObjectPrivilege = "UPDATE"
	SchemaObjectPrivilegeDelete       SchemaObjectPrivilege = "DELETE"
	SchemaObjectPrivilegeTruncate     SchemaObjectPrivilege = "TRUNCATE"
	SchemaObjectPrivilegeReferences   SchemaObjectPrivilege = "REFERENCES"

	// For Tag
	// SchemaObjectPrivilegeRead SchemaObjectPrivilege = "READ" (duplicate)

	// For TASK
	// SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET" (duplicate)
	// SchemaObjectPrivilegeMonitor 	SchemaObjectPrivilege = "MONITOR" (duplicate)
	// SchemaObjectPrivilegeOperate 	SchemaObjectPrivilege = "OPERATE" (duplicate)

	// For VIEW
	// SchemaObjectPrivilegeSelect		SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeReferences  SchemaObjectPrivilege = "REFERENCES" (duplicate)

	// For MATERIALIZED VIEW
	// SchemaObjectPrivilegeApplyBudget SchemaObjectPrivilege = "APPLYBUDGET" (duplicate)
	// SchemaObjectPrivilegeSelect 		SchemaObjectPrivilege = "SELECT" (duplicate)
	// SchemaObjectPrivilegeReferences 	SchemaObjectPrivilege = "REFERENCES" (duplicate)
)

func (p SchemaObjectPrivilege) String() string {
	return string(p)
}

type ObjectPrivilege string

const (
	ObjectPrivilegeReferenceUsage ObjectPrivilege = "REFERENCE_USAGE"
	ObjectPrivilegeUsage          ObjectPrivilege = "USAGE"
	ObjectPrivilegeSelect         ObjectPrivilege = "SELECT"
	ObjectPrivilegeRead           ObjectPrivilege = "READ"
)

func (p ObjectPrivilege) String() string {
	return string(p)
}
