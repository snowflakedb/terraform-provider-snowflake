package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowResultSchemaDef is the input definition for show/describe schema generation.
type ShowResultSchemaDef struct {
	ObjectStruct any
	// IsDescribe generates DescribeXSchema and {snake}_desc_gen.go (trims a trailing _details).
	IsDescribe bool
	// SkipFields omits snake_case schema keys from both the schema map and mapping.
	SkipFields []string
	// TypeOverrides maps snake_case schema keys to a Terraform schema type.
	// Currently, no logic is implemented.
	TypeOverrides map[string]schema.ValueType
	// UsedAsListEntry generates NameSchema (no Show/Describe prefix) for a property-row list Elem.
	UsedAsListEntry bool
	// AdditionalMapping generates a mapper type that must implement additionalSchemaMapper[T] in *_ext.go. Not inferred from SkipFields.
	AdditionalMapping bool
}

// ShowResultSchemaDetails is the extracted generator input (struct fields + definition metadata).
type ShowResultSchemaDetails struct {
	IsDescribe        bool
	SkipFields        []string
	TypeOverrides     map[string]schema.ValueType
	UsedAsListEntry   bool
	AdditionalMapping bool
	genhelpers.StructDetails
}

var SdkShowResultStructs = []ShowResultSchemaDef{
	{ObjectStruct: sdk.Account{}},
	{ObjectStruct: sdk.Alert{}},
	{ObjectStruct: sdk.ApiIntegration{}},
	{ObjectStruct: sdk.ApplicationPackage{}},
	{ObjectStruct: sdk.ApplicationRole{}},
	{ObjectStruct: sdk.Application{}},
	{ObjectStruct: sdk.AuthenticationPolicy{}, SkipFields: []string{"target_scopes"}}, // TODO [next PRs]: un-skip target_scopes (stale public schema).
	{ObjectStruct: sdk.CatalogIntegration{}},
	{ObjectStruct: sdk.ComputePool{}, SkipFields: []string{"backup_instance_families"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Connection{}, SkipFields: []string{"failover_allowed_to_accounts"}, AdditionalMapping: true},
	{ObjectStruct: sdk.CortexAgent{}, SkipFields: []string{"profile"}, AdditionalMapping: true},
	{ObjectStruct: sdk.DatabaseRole{}},
	{ObjectStruct: sdk.Database{}},
	{ObjectStruct: sdk.DynamicTable{}},
	{ObjectStruct: sdk.EventTable{}},
	{ObjectStruct: sdk.ExternalAccessIntegration{}},
	{ObjectStruct: sdk.ExternalFunction{}},
	{ObjectStruct: sdk.ExternalTable{}},
	{ObjectStruct: sdk.ExternalVolume{}},
	{ObjectStruct: sdk.FailoverGroup{}},
	{ObjectStruct: sdk.FileFormat{}},
	{ObjectStruct: sdk.FileFormatLegacy{}},
	{ObjectStruct: sdk.Function{}, SkipFields: []string{"arguments_old", "return_type_old"}},
	{ObjectStruct: sdk.GitRepository{}},
	{ObjectStruct: sdk.Grant{}, SkipFields: []string{"grant_on", "grant_to"}},
	{ObjectStruct: sdk.HybridTable{}},
	{ObjectStruct: sdk.HybridTableConstraint{}, SkipFields: []string{"columns", "referenced_table", "referenced_columns", "delete_rule", "update_rule"}, AdditionalMapping: true},
	{ObjectStruct: sdk.HybridTableIndex{}},
	{ObjectStruct: sdk.IcebergTable{}, SkipFields: []string{"auto_refresh_status", "partition_specs"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ImageRepository{}},
	{ObjectStruct: sdk.Listing{}},
	{ObjectStruct: sdk.ManagedAccount{}},
	{ObjectStruct: sdk.MaskingPolicy{}, SkipFields: []string{"options"}},
	{ObjectStruct: sdk.MaterializedView{}},
	{ObjectStruct: sdk.McpServer{}},
	{ObjectStruct: sdk.NetworkPolicy{}},
	{ObjectStruct: sdk.NetworkRule{}},
	{ObjectStruct: sdk.Notebook{}},
	{ObjectStruct: sdk.NotificationIntegration{}},
	{ObjectStruct: sdk.OpenflowConnectorDefinition{}, SkipFields: []string{"categories"}, AdditionalMapping: true},
	{ObjectStruct: sdk.OpenflowConnector{}},
	{ObjectStruct: sdk.OpenflowDeployment{}},
	{ObjectStruct: sdk.OpenflowRuntime{}, SkipFields: []string{"external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.OrganizationAccount{}},
	{ObjectStruct: sdk.Parameter{}},
	{ObjectStruct: sdk.PasswordPolicy{}},
	{ObjectStruct: sdk.Pipe{}},
	{ObjectStruct: sdk.PolicyReference{}},
	{ObjectStruct: sdk.PostgresInstance{}, SkipFields: []string{"is_highly_available"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Procedure{}, SkipFields: []string{"arguments_old", "return_type_old"}},
	{ObjectStruct: sdk.ReplicationAccount{}, SkipFields: []string{"comment"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ReplicationDatabase{}},
	{ObjectStruct: sdk.Region{}},
	{ObjectStruct: sdk.ResourceMonitor{}, SkipFields: []string{"notify_at", "notify_users", "suspend_immediately_at"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Role{}},
	{ObjectStruct: sdk.RowAccessPolicy{}},
	{ObjectStruct: sdk.Schema{}},
	{ObjectStruct: sdk.Secret{}, SkipFields: []string{"oauth_scopes"}, AdditionalMapping: true},
	{ObjectStruct: sdk.SecurityIntegration{}},
	{ObjectStruct: sdk.SemanticView{}},
	{ObjectStruct: sdk.Service{}, SkipFields: []string{"external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Sequence{}},
	{ObjectStruct: sdk.SessionPolicy{}, SkipFields: []string{"target_scopes"}}, // TODO [next PRs]: un-skip target_scopes (stale public schema).
	// SkipFields `name` remapped in ext; `owner_account` is stale public schema (add later).
	// TODO [next PRs]: un-skip owner_account.
	{ObjectStruct: sdk.Share{}, SkipFields: []string{"owner_account", "name"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Stage{}},
	{ObjectStruct: sdk.StorageIntegration{}},
	{ObjectStruct: sdk.StorageLifecyclePolicy{}},
	{ObjectStruct: sdk.Streamlit{}},
	{ObjectStruct: sdk.Stream{}, SkipFields: []string{"base_tables"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Table{}},
	{ObjectStruct: sdk.Tag{}, SkipFields: []string{"allowed_values"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Task{}, SkipFields: []string{"predecessors", "task_relations", "target_completion_interval"}, AdditionalMapping: true},
	{ObjectStruct: sdk.User{}},
	{ObjectStruct: sdk.ProgrammaticAccessToken{}},
	{ObjectStruct: sdk.View{}},
	{ObjectStruct: sdk.Warehouse{}},
	{ObjectStruct: sdk.AlertDetails{}, IsDescribe: true},
	// SkipFields uses generator snake_case (OAuth → o_auth, SigV4 → sig_v4); ext re-adds public oauth_* / sigv4_* keys.
	{ObjectStruct: sdk.CatalogIntegrationAllDetails{}, IsDescribe: true, SkipFields: []string{"rest_config", "o_auth_rest_authentication", "bearer_rest_authentication", "sig_v4_rest_authentication"}, AdditionalMapping: true},
	{ObjectStruct: sdk.CatalogIntegrationAwsGlueDetails{}, IsDescribe: true},
	// SkipFields uses generator snake_case (OAuth → o_auth, SigV4 → sig_v4); ext re-adds public oauth_* / sigv4_* keys.
	{ObjectStruct: sdk.CatalogIntegrationIcebergRestDetails{}, IsDescribe: true, SkipFields: []string{"rest_config", "o_auth_rest_authentication", "bearer_rest_authentication", "sig_v4_rest_authentication"}, AdditionalMapping: true},
	{ObjectStruct: sdk.CatalogIntegrationObjectStorageDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.CatalogIntegrationOpenCatalogDetails{}, IsDescribe: true, SkipFields: []string{"rest_config", "rest_authentication"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ComputePoolDetails{}, IsDescribe: true, SkipFields: []string{"backup_instance_families"}, AdditionalMapping: true},
	{ObjectStruct: sdk.CortexAgentDetails{}, IsDescribe: true, SkipFields: []string{"profile"}, AdditionalMapping: true},
	// SkipFields `attribute_columns` / `columns` re-added in ext.
	// TODO [next PRs]: un-skip attribute_columns / columns once MapToSchemaField maps []string.
	// TODO [next PRs]: un-skip serving_state / primary_key_columns / scoring_profile_count / full_index_build_interval_days (stale public schema).
	{ObjectStruct: sdk.CortexSearchServiceDetails{}, IsDescribe: true, SkipFields: []string{
		"attribute_columns", "columns",
		"serving_state", "primary_key_columns", "scoring_profile_count", "full_index_build_interval_days",
	}, AdditionalMapping: true},
	{ObjectStruct: sdk.DynamicTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.EventTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.ExternalAccessIntegrationDetails{}, IsDescribe: true, SkipFields: []string{"allowed_network_rules", "allowed_api_authentication_integrations", "allowed_authentication_secrets"}, AdditionalMapping: true},
	// SkipFields `id` is the SDK identifier (public schema has no `id`); nested storage_locations in ext.
	{ObjectStruct: sdk.ExternalVolumeDetails{}, IsDescribe: true, SkipFields: []string{"id", "storage_locations"}, AdditionalMapping: true},
	// SkipFields nested type-specific structs (public schema is a flat union of per-type keys; nested mapping in ext).
	{ObjectStruct: sdk.FileFormatAllDetails{}, IsDescribe: true, SkipFields: []string{"csv", "json", "avro", "orc", "parquet", "xml"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip null_if once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatAvro{}, IsDescribe: true, SkipFields: []string{"null_if"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip null_if once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatCsv{}, IsDescribe: true, SkipFields: []string{"null_if"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip null_if once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatJson{}, IsDescribe: true, SkipFields: []string{"null_if"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip null_if once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatOrc{}, IsDescribe: true, SkipFields: []string{"null_if"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip null_if once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatParquet{}, IsDescribe: true, SkipFields: []string{"null_if"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip disable_snowflake_data (stale public schema; still a DESCRIBE column).
	{ObjectStruct: sdk.FileFormatXml{}, IsDescribe: true, SkipFields: []string{"disable_snowflake_data"}},
	// SkipFields omits SDK-only identifier/normalized fields (not Snowflake DESCRIBE properties; schema is not wired yet).
	// TODO [next PRs]: un-skip return_data_type once MapToSchemaField maps datatypes.DataType via ToSql().
	{ObjectStruct: sdk.FunctionDetails{}, IsDescribe: true, SkipFields: []string{
		"id", "normalized_imports", "normalized_target_path", "return_data_type",
		"normalized_arguments", "normalized_external_access_integrations", "normalized_secrets", "normalized_packages",
	}, AdditionalMapping: true},
	{ObjectStruct: sdk.HybridTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.IcebergTableDetails{}, IsDescribe: true, SkipFields: []string{"type", "data_type_raw"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ListingDetails{}, IsDescribe: true},
	// TODO [next PRs]: un-skip return_type once MapToSchemaField maps datatypes.DataType via ToSql(); signature stays ext (slice of structs).
	{ObjectStruct: sdk.MaskingPolicyDetails{}, IsDescribe: true, SkipFields: []string{"signature", "return_type"}, AdditionalMapping: true},
	{ObjectStruct: sdk.MaterializedViewDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.McpServerDetails{}, IsDescribe: true},
	// Property-row list entry (same pattern as SecurityIntegrationProperty). Flattened DescribeNetworkPolicySchema stays the name-keyed consumer.
	{ObjectStruct: sdk.NetworkPolicyProperty{}, UsedAsListEntry: true},
	{ObjectStruct: sdk.NetworkRuleDetails{}, IsDescribe: true, SkipFields: []string{"value_list"}, AdditionalMapping: true},
	// SkipFields `id` is the SDK identifier (public schema uses `name`).
	{ObjectStruct: sdk.NotebookDetails{}, IsDescribe: true, SkipFields: []string{"id"}},
	// SkipFields `id` is the SDK identifier (public schema uses `name`).
	{ObjectStruct: sdk.OpenflowConnectorDetails{}, IsDescribe: true, SkipFields: []string{"id"}},
	{ObjectStruct: sdk.OpenflowDeploymentDetails{}, IsDescribe: true},
	// SkipFields `id` is the SDK identifier (public schema uses `name`); slice is re-added in ext.
	{ObjectStruct: sdk.OpenflowRuntimeDetails{}, IsDescribe: true, SkipFields: []string{"id", "external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.PasswordPolicyDetails{}, IsDescribe: true},
	// SkipFields omits SDK parse helpers (not DESCRIBE columns; keep omitted).
	{ObjectStruct: sdk.PostgresInstanceDetails{}, IsDescribe: true, SkipFields: []string{"has_any_running_operations", "operation_errors"}},
	// SkipFields omits SDK-only identifier/normalized fields (not Snowflake DESCRIBE properties; schema is not wired yet).
	// TODO [next PRs]: un-skip return_data_type once MapToSchemaField maps datatypes.DataType via ToSql().
	{ObjectStruct: sdk.ProcedureDetails{}, IsDescribe: true, SkipFields: []string{
		"id", "normalized_imports", "normalized_target_path", "return_data_type",
		"normalized_arguments", "normalized_external_access_integrations", "normalized_secrets", "normalized_packages",
		"snowpark_version",
	}, AdditionalMapping: true},
	{ObjectStruct: sdk.SchemaDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.SecretDetails{}, IsDescribe: true, SkipFields: []string{"oauth_scopes"}, AdditionalMapping: true},
	{ObjectStruct: sdk.SecurityIntegrationProperty{}, UsedAsListEntry: true},
	{ObjectStruct: sdk.ServiceDetails{}, IsDescribe: true, SkipFields: []string{"external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.SessionPolicyDetails{}, IsDescribe: true, SkipFields: []string{"allowed_secondary_roles", "blocked_secondary_roles"}, AdditionalMapping: true},
	{ObjectStruct: sdk.StorageIntegrationAllDetails{}, IsDescribe: true, SkipFields: []string{"allowed_locations", "blocked_locations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.StorageIntegrationAwsDetails{}, IsDescribe: true, SkipFields: []string{"allowed_locations", "blocked_locations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.StorageIntegrationAzureDetails{}, IsDescribe: true, SkipFields: []string{"allowed_locations", "blocked_locations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.StorageIntegrationGcsDetails{}, IsDescribe: true, SkipFields: []string{"allowed_locations", "blocked_locations"}, AdditionalMapping: true},
	// TODO [next PRs]: un-skip return_type once MapToSchemaField maps datatypes.DataType via ToSql(); signature stays ext (slice of structs).
	{ObjectStruct: sdk.StorageLifecyclePolicyDetails{}, IsDescribe: true, SkipFields: []string{"signature", "return_type"}, AdditionalMapping: true},
	// SkipFields `root_location` (ParseRootLocation rewrite) and []string TypeSets.
	{ObjectStruct: sdk.StreamlitDetail{}, IsDescribe: true, SkipFields: []string{"root_location", "user_packages", "import_urls", "external_access_integrations"}, AdditionalMapping: true},
	// SkipFields `check`: SDK is *bool; public describe_output is TypeString.
	{ObjectStruct: sdk.TableColumnDetails{}, IsDescribe: true, SkipFields: []string{"check"}, AdditionalMapping: true},
	// TODO [next PRs]: UserDetails SkipFields+ext is temporary (no XxxProperty mapping).
	// P2: flatten the SDK struct or add dedicated property handling, then un-skip.
	// Keep omitted: password (secret). Add later: rsa_public_key_last_set_time / rsa_public_key2_last_set_time (stale public schema).
	{ObjectStruct: sdk.UserDetails{}, IsDescribe: true, SkipFields: []string{
		"name", "comment", "display_name", "type", "login_name", "first_name", "middle_name", "last_name", "email", "password",
		"must_change_password", "disabled", "snowflake_lock", "snowflake_support", "days_to_expiry", "mins_to_unlock",
		"default_warehouse", "default_namespace", "default_role", "default_secondary_roles", "ext_authn_duo", "ext_authn_uid",
		"mins_to_bypass_mfa", "mins_to_bypass_network_policy", "rsa_public_key", "rsa_public_key_fp", "rsa_public_key_last_set_time",
		"rsa_public_key2", "rsa_public_key2_fp", "rsa_public_key2_last_set_time", "password_last_set_time",
		"custom_landing_page_url", "custom_landing_page_url_flush_next_ui_load", "has_mfa", "has_workload_identity",
	}, AdditionalMapping: true},
	// SkipFields `check`: SDK is *bool; public describe_output is TypeString.
	{ObjectStruct: sdk.ViewDetails{}, IsDescribe: true, SkipFields: []string{"check"}, AdditionalMapping: true},
	{ObjectStruct: sdk.WarehouseDetails{}, IsDescribe: true},
}

func GetShowResultSchemaDetails() []ShowResultSchemaDetails {
	allDetails := make([]ShowResultSchemaDetails, len(SdkShowResultStructs))
	for idx, d := range SdkShowResultStructs {
		allDetails[idx] = ShowResultSchemaDetails{
			IsDescribe:        d.IsDescribe,
			SkipFields:        d.SkipFields,
			TypeOverrides:     d.TypeOverrides,
			UsedAsListEntry:   d.UsedAsListEntry,
			AdditionalMapping: d.AdditionalMapping,
			StructDetails:     genhelpers.ExtractStructDetails(d.ObjectStruct),
		}
	}
	return allDetails
}
