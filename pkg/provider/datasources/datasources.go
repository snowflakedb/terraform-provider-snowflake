package datasources

type datasource string

const (
	Accounts                       datasource = "snowflake_accounts"
	AccountRoles                   datasource = "snowflake_account_roles"
	Alerts                         datasource = "snowflake_alerts"
	Connections                    datasource = "snowflake_connections"
	CortexSearchServices           datasource = "snowflake_cortex_search_services"
	CurrentAccount                 datasource = "snowflake_current_account"
	CurrentRole                    datasource = "snowflake_current_role"
	Database                       datasource = "snowflake_database"
	DatabaseRole                   datasource = "snowflake_database_role"
	DatabaseRoles                  datasource = "snowflake_database_roles"
	Databases                      datasource = "snowflake_databases"
	DynamicTables                  datasource = "snowflake_dynamic_tables"
	ExternalFunctions              datasource = "snowflake_external_functions"
	ExternalTables                 datasource = "snowflake_external_tables"
	FailoverGroups                 datasource = "snowflake_failover_groups"
	FileFormats                    datasource = "snowflake_file_formats"
	Functions                      datasource = "snowflake_functions"
	Grants                         datasource = "snowflake_grants"
	ImageRepositories              datasource = "snowflake_image_repositories"
	MaskingPolicies                datasource = "snowflake_masking_policies"
	MaterializedViews              datasource = "snowflake_materialized_views"
	NetworkPolicies                datasource = "snowflake_network_policies"
	Parameters                     datasource = "snowflake_parameters"
	Pipes                          datasource = "snowflake_pipes"
	Procedures                     datasource = "snowflake_procedures"
	ResourceMonitors               datasource = "snowflake_resource_monitors"
	RowAccessPolicies              datasource = "snowflake_row_access_policies"
	Schemas                        datasource = "snowflake_schemas"
	Secrets                        datasource = "snowflake_secrets"
	SecurityIntegrations           datasource = "snowflake_security_integrations"
	Sequences                      datasource = "snowflake_sequences"
	Shares                         datasource = "snowflake_shares"
	Stages                         datasource = "snowflake_stages"
	StorageIntegrations            datasource = "snowflake_storage_integrations"
	Streams                        datasource = "snowflake_streams"
	Streamlits                     datasource = "snowflake_streamlits"
	SystemGenerateScimAccessToken  datasource = "snowflake_system_generate_scim_access_token"
	SystemGetAwsSnsIamPolicy       datasource = "snowflake_system_get_aws_sns_iam_policy"
	SystemGetPrivateLinkConfig     datasource = "snowflake_system_get_privatelink_config"
	SystemGetSnowflakePlatformInfo datasource = "snowflake_system_get_snowflake_platform_info"
	Tables                         datasource = "snowflake_tables"
	Tags                           datasource = "snowflake_tags"
	Tasks                          datasource = "snowflake_tasks"
	Users                          datasource = "snowflake_users"
	Views                          datasource = "snowflake_views"
	Warehouses                     datasource = "snowflake_warehouses"
)

type Datasource interface {
	xxxProtected()
	String() string
}

func (r datasource) xxxProtected() {}

func (r datasource) String() string {
	return string(r)
}
