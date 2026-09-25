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
	// SkipFields omits snake_case schema keys that should not become public (SDK-internal,
	// secret, stale, or add-later). A comment is left in place; no ext mapping is expected.
	SkipFields []string
	// ManualFields omits snake_case schema keys that ext must add
	// (unsupported conversion, rename, nested/slice mapping). A comment is left in place
	// and the additional-mapping hook is generated.
	ManualFields []string
	// TypeOverrides maps snake_case schema keys to a Terraform schema type.
	// Currently, no logic is implemented.
	TypeOverrides map[string]schema.ValueType
	// UsedAsListEntry generates NameSchema (no Show/Describe prefix) for a property-row list Elem.
	UsedAsListEntry bool
}

// ShowResultSchemaDetails is the extracted generator input (struct fields + definition metadata).
type ShowResultSchemaDetails struct {
	IsDescribe      bool
	SkipFields      []string
	ManualFields    []string
	TypeOverrides   map[string]schema.ValueType
	UsedAsListEntry bool
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
	{ObjectStruct: sdk.ComputePool{}, ManualFields: []string{"backup_instance_families"}},
	{ObjectStruct: sdk.Connection{}, ManualFields: []string{"failover_allowed_to_accounts"}},
	{ObjectStruct: sdk.CortexAgent{}, ManualFields: []string{"profile"}},
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
	{ObjectStruct: sdk.HybridTableConstraint{}, ManualFields: []string{"columns", "referenced_table", "referenced_columns", "delete_rule", "update_rule"}},
	{ObjectStruct: sdk.HybridTableIndex{}},
	{ObjectStruct: sdk.IcebergTable{}, ManualFields: []string{"auto_refresh_status", "partition_specs"}},
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
	{ObjectStruct: sdk.OpenflowConnectorDefinition{}, ManualFields: []string{"categories"}},
	{ObjectStruct: sdk.OpenflowConnector{}},
	{ObjectStruct: sdk.OpenflowDeployment{}},
	{ObjectStruct: sdk.OpenflowRuntime{}, ManualFields: []string{"external_access_integrations"}},
	{ObjectStruct: sdk.OrganizationAccount{}},
	{ObjectStruct: sdk.Parameter{}},
	{ObjectStruct: sdk.PasswordPolicy{}},
	{ObjectStruct: sdk.Pipe{}},
	{ObjectStruct: sdk.PolicyReference{}},
	{ObjectStruct: sdk.PostgresInstance{}, ManualFields: []string{"is_highly_available"}},
	{ObjectStruct: sdk.Procedure{}, SkipFields: []string{"arguments_old", "return_type_old"}},
	{ObjectStruct: sdk.ReplicationAccount{}, ManualFields: []string{"comment"}},
	{ObjectStruct: sdk.ReplicationDatabase{}},
	{ObjectStruct: sdk.Region{}},
	{ObjectStruct: sdk.ResourceMonitor{}, SkipFields: []string{"notify_at", "notify_users"}, ManualFields: []string{"suspend_immediately_at"}},
	{ObjectStruct: sdk.Role{}},
	{ObjectStruct: sdk.RowAccessPolicy{}},
	{ObjectStruct: sdk.Schema{}},
	{ObjectStruct: sdk.Secret{}, ManualFields: []string{"oauth_scopes"}},
	{ObjectStruct: sdk.SecurityIntegration{}},
	{ObjectStruct: sdk.SemanticView{}},
	{ObjectStruct: sdk.Service{}, ManualFields: []string{"external_access_integrations"}},
	{ObjectStruct: sdk.Sequence{}},
	{ObjectStruct: sdk.SessionPolicy{}, SkipFields: []string{"target_scopes"}}, // TODO [next PRs]: un-skip target_scopes (stale public schema).
	// ManualFields `name` remapped in ext; SkipFields `owner_account` is stale public schema (add later).
	// TODO [next PRs]: un-skip owner_account.
	{ObjectStruct: sdk.Share{}, SkipFields: []string{"owner_account"}, ManualFields: []string{"name"}},
	{ObjectStruct: sdk.Stage{}},
	{ObjectStruct: sdk.StorageIntegration{}},
	{ObjectStruct: sdk.StorageLifecyclePolicy{}},
	{ObjectStruct: sdk.Streamlit{}},
	{ObjectStruct: sdk.Stream{}, ManualFields: []string{"base_tables"}},
	{ObjectStruct: sdk.Table{}},
	{ObjectStruct: sdk.Tag{}, ManualFields: []string{"allowed_values"}},
	{ObjectStruct: sdk.Task{}, ManualFields: []string{"predecessors", "task_relations", "target_completion_interval"}},
	{ObjectStruct: sdk.User{}},
	{ObjectStruct: sdk.ProgrammaticAccessToken{}},
	{ObjectStruct: sdk.View{}},
	// Union of all SHOW WAREHOUSES columns (AllDetails analog). ManualFields `tables` re-added in warehouse_ext.go.
	// Keep omitted: actives / pendings / failed / suspended / uuid (Snowflake: internal use, will be removed).
	// TODO [next PRs]: drop tables from ManualFields once MapToSchemaField maps []SchemaObjectIdentifier.
	{ObjectStruct: sdk.Warehouse{}, SkipFields: []string{"actives", "pendings", "failed", "suspended", "uuid"}, ManualFields: []string{"tables"}},
	{ObjectStruct: sdk.WarehouseAdaptive{}},
	{ObjectStruct: sdk.WarehouseInteractive{}, ManualFields: []string{"tables"}},
	{ObjectStruct: sdk.WarehouseRegular{}},
	{ObjectStruct: sdk.AlertDetails{}, IsDescribe: true},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeApiAuthDetails parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.ApiAuthenticationSecurityIntegrationDetails{}, IsDescribe: true, ManualFields: []string{
		"enabled", "oauth_access_token_validity", "oauth_refresh_token_validity", "oauth_client_auth_method",
		"oauth_authorization_endpoint", "oauth_token_endpoint", "oauth_allowed_scopes", "oauth_grant",
		"parent_integration", "auth_type", "comment",
	}},
	// SkipFields `id` is the SDK identifier (public schema has no `id`).
	// ManualFields: api_key is Sensitive; api_provider needs strings.ToLower — neither can be generated.
	// TODO [next PRs]: drop prefixes/scopes/certs from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.ApiIntegrationAllDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{
		"api_key", "api_provider", "allowed_prefixes", "blocked_prefixes", "oauth_allowed_scopes", "tls_trusted_certificates",
	}},
	{ObjectStruct: sdk.ApiIntegrationAwsDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{
		"api_key", "api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	{ObjectStruct: sdk.ApiIntegrationAzureDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{
		"api_key", "api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	{ObjectStruct: sdk.ApiIntegrationExternalMcpDynamicClient{}, IsDescribe: true, ManualFields: []string{
		"api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	{ObjectStruct: sdk.ApiIntegrationExternalMcpOauth2{}, IsDescribe: true, ManualFields: []string{
		"api_provider", "allowed_prefixes", "blocked_prefixes", "oauth_allowed_scopes",
	}},
	{ObjectStruct: sdk.ApiIntegrationGitRepositoryGithubApp{}, IsDescribe: true, ManualFields: []string{
		"api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	{ObjectStruct: sdk.ApiIntegrationGitRepositoryOauth2{}, IsDescribe: true, ManualFields: []string{
		"allowed_prefixes", "blocked_prefixes", "oauth_allowed_scopes",
	}},
	{ObjectStruct: sdk.ApiIntegrationGitRepositoryPrivateLink{}, IsDescribe: true, ManualFields: []string{
		"api_provider", "allowed_prefixes", "blocked_prefixes", "tls_trusted_certificates",
	}},
	{ObjectStruct: sdk.ApiIntegrationGitRepositoryToken{}, IsDescribe: true, ManualFields: []string{
		"api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	{ObjectStruct: sdk.ApiIntegrationGoogleDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{
		"api_key", "api_provider", "allowed_prefixes", "blocked_prefixes",
	}},
	// TypeString of each property .Value (today’s public describe_output).
	{ObjectStruct: sdk.AuthenticationPolicyDescribeDetails{}, IsDescribe: true},
	// Property-row list entry (same pattern as SecurityIntegrationProperty / NetworkPolicyProperty).
	{ObjectStruct: sdk.AuthenticationPolicyDescription{}, UsedAsListEntry: true},
	// ManualFields uses generator snake_case (OAuth → o_auth, SigV4 → sig_v4); ext re-adds public oauth_* / sigv4_* keys.
	{ObjectStruct: sdk.CatalogIntegrationAllDetails{}, IsDescribe: true, ManualFields: []string{"rest_config", "o_auth_rest_authentication", "bearer_rest_authentication", "sig_v4_rest_authentication"}},
	{ObjectStruct: sdk.CatalogIntegrationAwsGlueDetails{}, IsDescribe: true},
	// ManualFields uses generator snake_case (OAuth → o_auth, SigV4 → sig_v4); ext re-adds public oauth_* / sigv4_* keys.
	{ObjectStruct: sdk.CatalogIntegrationIcebergRestDetails{}, IsDescribe: true, ManualFields: []string{"rest_config", "o_auth_rest_authentication", "bearer_rest_authentication", "sig_v4_rest_authentication"}},
	{ObjectStruct: sdk.CatalogIntegrationObjectStorageDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.CatalogIntegrationOpenCatalogDetails{}, IsDescribe: true, ManualFields: []string{"rest_config", "rest_authentication"}},
	{ObjectStruct: sdk.ComputePoolDetails{}, IsDescribe: true, ManualFields: []string{"backup_instance_families"}},
	{ObjectStruct: sdk.CortexAgentDetails{}, IsDescribe: true, ManualFields: []string{"profile"}},
	// ManualFields `attribute_columns` / `columns` re-added in ext.
	// TODO [next PRs]: un-skip attribute_columns / columns once MapToSchemaField maps []string.
	// TODO [next PRs]: un-skip serving_state / primary_key_columns / scoring_profile_count / full_index_build_interval_days (stale public schema).
	{ObjectStruct: sdk.CortexSearchServiceDetails{}, IsDescribe: true, ManualFields: []string{"attribute_columns", "columns"}, SkipFields: []string{
		"serving_state", "primary_key_columns", "scoring_profile_count", "full_index_build_interval_days",
	}},
	// Public describe_output Elem is the row (created_on / name / kind). DatabaseDetails is a Rows wrapper; list helper in ext.
	{ObjectStruct: sdk.DatabaseDetailsRow{}, IsDescribe: true},
	{ObjectStruct: sdk.DynamicTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.EventTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.ExternalAccessIntegrationDetails{}, IsDescribe: true, ManualFields: []string{"allowed_network_rules", "allowed_api_authentication_integrations", "allowed_authentication_secrets"}},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// ManualFields `external_oauth_rsa_public_key2` is ToSnakeCase of ExternalOauthRsaPublicKey2;
	// public key is external_oauth_rsa_public_key_2.
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeExternalOauthDetails parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.ExternalOauthSecurityIntegrationDetails{}, IsDescribe: true, ManualFields: []string{
		"enabled", "external_oauth_issuer", "external_oauth_jws_keys_url", "external_oauth_any_role_mode",
		"external_oauth_rsa_public_key", "external_oauth_rsa_public_key2", "external_oauth_blocked_roles_list",
		"external_oauth_allowed_roles_list", "external_oauth_audience_list", "external_oauth_token_user_mapping_claim",
		"external_oauth_snowflake_user_mapping_attribute", "external_oauth_scope_delimiter", "comment",
	}},
	// SkipFields `id` is the SDK identifier (public schema has no `id`); nested storage_locations is ManualFields.
	{ObjectStruct: sdk.ExternalVolumeDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{"storage_locations"}},
	// ManualFields nested type-specific structs (public schema is a flat union of per-type keys; nested mapping in ext).
	{ObjectStruct: sdk.FileFormatAllDetails{}, IsDescribe: true, ManualFields: []string{"csv", "json", "avro", "orc", "parquet", "xml"}},
	// TODO [next PRs]: drop null_if from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatAvro{}, IsDescribe: true, ManualFields: []string{"null_if"}},
	// TODO [next PRs]: drop null_if from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatCsv{}, IsDescribe: true, ManualFields: []string{"null_if"}},
	// TODO [next PRs]: drop null_if from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatJson{}, IsDescribe: true, ManualFields: []string{"null_if"}},
	// TODO [next PRs]: drop null_if from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatOrc{}, IsDescribe: true, ManualFields: []string{"null_if"}},
	// TODO [next PRs]: drop null_if from ManualFields once MapToSchemaField maps []string.
	{ObjectStruct: sdk.FileFormatParquet{}, IsDescribe: true, ManualFields: []string{"null_if"}},
	// TODO [next PRs]: un-skip disable_snowflake_data (stale public schema; still a DESCRIBE column).
	{ObjectStruct: sdk.FileFormatXml{}, IsDescribe: true, SkipFields: []string{"disable_snowflake_data"}},
	// SkipFields omits SDK-only identifier/normalized fields (not Snowflake DESCRIBE properties; schema is not wired yet).
	// TODO [next PRs]: drop return_data_type from ManualFields once MapToSchemaField maps datatypes.DataType via ToSql().
	{ObjectStruct: sdk.FunctionDetails{}, IsDescribe: true, SkipFields: []string{
		"id", "normalized_imports", "normalized_target_path",
		"normalized_arguments", "normalized_external_access_integrations", "normalized_secrets", "normalized_packages",
	}, ManualFields: []string{"return_data_type"}},
	{ObjectStruct: sdk.HybridTableDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.IcebergTableDetails{}, IsDescribe: true, SkipFields: []string{"data_type_raw"}, ManualFields: []string{"type"}},
	{ObjectStruct: sdk.ListingDetails{}, IsDescribe: true},
	// TODO [next PRs]: drop return_type from ManualFields once MapToSchemaField maps datatypes.DataType via ToSql(); signature stays ext (slice of structs).
	{ObjectStruct: sdk.MaskingPolicyDetails{}, IsDescribe: true, ManualFields: []string{"signature", "return_type"}},
	{ObjectStruct: sdk.MaterializedViewDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.McpServerDetails{}, IsDescribe: true},
	// TypeString of each property .Value (today’s public describe_output).
	{ObjectStruct: sdk.NetworkPolicyDetails{}, IsDescribe: true},
	// Property-row list entry (same pattern as SecurityIntegrationProperty). Flattened DescribeNetworkPolicyDetailsSchema stays the name-keyed consumer.
	{ObjectStruct: sdk.NetworkPolicyProperty{}, UsedAsListEntry: true},
	{ObjectStruct: sdk.NetworkRuleDetails{}, IsDescribe: true, ManualFields: []string{"value_list"}},
	// SkipFields `id` is the SDK identifier (public schema uses `name`).
	{ObjectStruct: sdk.NotebookDetails{}, IsDescribe: true, SkipFields: []string{"id"}},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// ManualFields `oauth_client_rsa_public_key2_fp` is ToSnakeCase of OauthClientRsaPublicKey2Fp;
	// public key is oauth_client_rsa_public_key_2_fp.
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeOauthCustomDetails parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.OauthIntegrationForCustomClientsDetails{}, IsDescribe: true, ManualFields: []string{
		"oauth_client_type", "enabled", "oauth_allow_non_tls_redirect_uri", "oauth_enforce_pkce",
		"oauth_use_secondary_roles", "pre_authorized_roles_list", "allowed_roles_list", "blocked_roles_list",
		"oauth_issue_refresh_tokens", "oauth_refresh_token_validity", "network_policy",
		"oauth_client_rsa_public_key_fp", "oauth_client_rsa_public_key2_fp", "comment",
		"oauth_authorization_endpoint", "oauth_token_endpoint",
		"oauth_allowed_authorization_endpoints", "oauth_allowed_token_endpoints",
	}},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// ManualFields `oauth_client_rsa_public_key2_fp` is ToSnakeCase of OauthClientRsaPublicKey2Fp;
	// public key is oauth_client_rsa_public_key_2_fp.
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeOauthPartnerDetails parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.OauthIntegrationForPartnerApplicationsDetails{}, IsDescribe: true, ManualFields: []string{
		"oauth_client_type", "enabled", "oauth_allow_non_tls_redirect_uri", "oauth_enforce_pkce",
		"oauth_use_secondary_roles", "pre_authorized_roles_list", "allowed_roles_list", "blocked_roles_list",
		"oauth_issue_refresh_tokens", "oauth_refresh_token_validity", "network_policy",
		"oauth_client_rsa_public_key_fp", "oauth_client_rsa_public_key2_fp", "comment",
		"oauth_authorization_endpoint", "oauth_token_endpoint",
		"oauth_allowed_authorization_endpoints", "oauth_allowed_token_endpoints",
	}},
	// SkipFields `id` is the SDK identifier (public schema uses `name`).
	{ObjectStruct: sdk.OpenflowConnectorDetails{}, IsDescribe: true, SkipFields: []string{"id"}},
	{ObjectStruct: sdk.OpenflowDeploymentDetails{}, IsDescribe: true},
	// SkipFields `id` is the SDK identifier (public schema uses `name`); slice is ManualFields.
	{ObjectStruct: sdk.OpenflowRuntimeDetails{}, IsDescribe: true, SkipFields: []string{"id"}, ManualFields: []string{"external_access_integrations"}},
	{ObjectStruct: sdk.PasswordPolicyDetails{}, IsDescribe: true},
	// SkipFields omits SDK parse helpers (not DESCRIBE columns; keep omitted).
	{ObjectStruct: sdk.PostgresInstanceDetails{}, IsDescribe: true, SkipFields: []string{"has_any_running_operations", "operation_errors"}},
	// SkipFields omits SDK-only identifier/normalized fields (not Snowflake DESCRIBE properties; schema is not wired yet).
	// TODO [next PRs]: drop return_data_type from ManualFields once MapToSchemaField maps datatypes.DataType via ToSql().
	{ObjectStruct: sdk.ProcedureDetails{}, IsDescribe: true, SkipFields: []string{
		"id", "normalized_imports", "normalized_target_path",
		"normalized_arguments", "normalized_external_access_integrations", "normalized_secrets", "normalized_packages",
		"snowpark_version",
	}, ManualFields: []string{"return_data_type"}},
	// ManualFields `signature` stays ext (slice of structs). `return_type` is string (not datatypes.DataType) so it generates natively.
	{ObjectStruct: sdk.RowAccessPolicyDescription{}, IsDescribe: true, ManualFields: []string{"signature"}},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeSaml2Details parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.Saml2SecurityIntegrationDetails{}, IsDescribe: true, ManualFields: []string{
		"saml2_issuer", "saml2_sso_url", "saml2_provider", "saml2_sp_initiated_login_page_label",
		"saml2_enable_sp_initiated", "saml2_sign_request", "saml2_requested_nameid_format",
		"saml2_post_logout_redirect_url", "saml2_force_authn", "saml2_snowflake_issuer_url",
		"saml2_snowflake_acs_url", "saml2_snowflake_metadata", "saml2_digest_methods_used",
		"saml2_signature_methods_used", "allowed_user_domains", "allowed_email_patterns", "comment",
	}},
	{ObjectStruct: sdk.SchemaDetails{}, IsDescribe: true},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output).
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via a details
	// struct + DescribeScimDetails parse (API/storage-integration analog); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.ScimSecurityIntegrationDetails{}, IsDescribe: true, ManualFields: []string{
		"enabled", "network_policy", "run_as_role", "sync_password", "comment",
	}},
	{ObjectStruct: sdk.SecretDetails{}, IsDescribe: true, ManualFields: []string{"oauth_scopes"}},
	{ObjectStruct: sdk.SecurityIntegrationProperty{}, UsedAsListEntry: true},
	{ObjectStruct: sdk.ServiceDetails{}, IsDescribe: true, ManualFields: []string{"external_access_integrations"}},
	{ObjectStruct: sdk.SessionPolicyDetails{}, IsDescribe: true, ManualFields: []string{"allowed_secondary_roles", "blocked_secondary_roles"}},
	// Nested directory_table / file_format / location / privatelink in ext (public key privatelink, not private_link).
	// Keep omitted: id (SDK identifier), credentials (secret).
	// TODO [next PRs]: first move only — generated schemas are empty (every public key is nested/slice).
	// Native slices/nested structs should drop from ManualFields and shrink the ext. Same pattern for API variant DESCRIBE.
	// TODO [next PRs]: DirectoryTable directory_notification_channel / aws_sns_topic (stale public schema).
	{ObjectStruct: sdk.StageAws{}, IsDescribe: true, ManualFields: []string{
		"file_format_name", "file_format_csv", "file_format_json", "file_format_avro", "file_format_orc", "file_format_parquet", "file_format_xml",
		"directory_table", "private_link", "location",
	}},
	{ObjectStruct: sdk.StageAwsCompatible{}, IsDescribe: true, ManualFields: []string{
		"file_format_name", "file_format_csv", "file_format_json", "file_format_avro", "file_format_orc", "file_format_parquet", "file_format_xml",
		"directory_table", "location",
	}},
	{ObjectStruct: sdk.StageCommon{}, IsDescribe: true, ManualFields: []string{
		"file_format_name", "file_format_csv", "file_format_json", "file_format_avro", "file_format_orc", "file_format_parquet", "file_format_xml",
		"directory_table",
	}},
	{ObjectStruct: sdk.StageDetails{}, IsDescribe: true, SkipFields: []string{"id", "credentials"}, ManualFields: []string{
		"file_format_name", "file_format_csv", "file_format_json", "file_format_avro", "file_format_orc", "file_format_parquet", "file_format_xml",
		"directory_table", "private_link", "location",
	}},
	{ObjectStruct: sdk.StorageIntegrationAllDetails{}, IsDescribe: true, ManualFields: []string{"allowed_locations", "blocked_locations"}},
	{ObjectStruct: sdk.StorageIntegrationAwsDetails{}, IsDescribe: true, ManualFields: []string{"allowed_locations", "blocked_locations"}},
	{ObjectStruct: sdk.StorageIntegrationAzureDetails{}, IsDescribe: true, ManualFields: []string{"allowed_locations", "blocked_locations"}},
	// Nested DescribePropertyListSchema in ext (today’s public describe_output on the legacy resource).
	// TODO [v3]: flatten to typed scalars (enabled.0.value → enabled) via StorageIntegrationAllDetails
	// / DescribeDetails (already exist); native TypeList of the property row is not the end state.
	{ObjectStruct: sdk.StorageIntegrationDetails{}, IsDescribe: true, ManualFields: []string{
		"enabled", "storage_provider", "storage_allowed_locations", "storage_blocked_locations",
		"storage_aws_iam_user_arn", "storage_aws_object_acl", "storage_aws_role_arn", "storage_aws_external_id",
		"storage_gcp_service_account", "azure_consent_url", "azure_multi_tenant_app_name",
		"use_privatelink_endpoint", "comment",
	}},
	{ObjectStruct: sdk.StorageIntegrationGcsDetails{}, IsDescribe: true, ManualFields: []string{"allowed_locations", "blocked_locations"}},
	// Property-row list entry (same pattern as SecurityIntegrationProperty). Nested DescribeStorageIntegrationDetailsSchema stays the name-keyed consumer.
	{ObjectStruct: sdk.StorageIntegrationProperty{}, UsedAsListEntry: true},
	// TODO [next PRs]: drop return_type from ManualFields once MapToSchemaField maps datatypes.DataType via ToSql(); signature stays ext (slice of structs).
	{ObjectStruct: sdk.StorageLifecyclePolicyDetails{}, IsDescribe: true, ManualFields: []string{"signature", "return_type"}},
	// ManualFields `root_location` (ParseRootLocation rewrite) and []string TypeSets.
	{ObjectStruct: sdk.StreamlitDetail{}, IsDescribe: true, ManualFields: []string{"root_location", "user_packages", "import_urls", "external_access_integrations"}},
	// ManualFields `check`: SDK is *bool; public describe_output is TypeString.
	{ObjectStruct: sdk.TableColumnDetails{}, IsDescribe: true, ManualFields: []string{"check"}},
	// Keep omitted: password (secret). Add later: rsa_public_key_last_set_time / rsa_public_key2_last_set_time (stale public schema).
	{ObjectStruct: sdk.UserDetails{}, IsDescribe: true, SkipFields: []string{
		"password", "rsa_public_key_last_set_time", "rsa_public_key2_last_set_time",
	}},
	// ManualFields `check`: SDK is *bool; public describe_output is TypeString.
	{ObjectStruct: sdk.ViewDetails{}, IsDescribe: true, ManualFields: []string{"check"}},
	{ObjectStruct: sdk.WarehouseDetails{}, IsDescribe: true},
}

func GetShowResultSchemaDetails() []ShowResultSchemaDetails {
	allDetails := make([]ShowResultSchemaDetails, len(SdkShowResultStructs))
	for idx, d := range SdkShowResultStructs {
		allDetails[idx] = ShowResultSchemaDetails{
			IsDescribe:      d.IsDescribe,
			SkipFields:      d.SkipFields,
			ManualFields:    d.ManualFields,
			TypeOverrides:   d.TypeOverrides,
			UsedAsListEntry: d.UsedAsListEntry,
			StructDetails:   genhelpers.ExtractStructDetails(d.ObjectStruct),
		}
	}
	return allDetails
}
