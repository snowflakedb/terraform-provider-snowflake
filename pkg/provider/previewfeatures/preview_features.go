package previewfeatures

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type feature string

const (
	AccountAuthenticationPolicyAttachmentResource  feature = "snowflake_account_authentication_policy_attachment_resource"
	AccountPasswordPolicyAttachmentResource        feature = "snowflake_account_password_policy_attachment_resource"
	AccountSessionPolicyAttachmentResource         feature = "snowflake_account_session_policy_attachment_resource"
	AlertResource                                  feature = "snowflake_alert_resource"
	AlertsDatasource                               feature = "snowflake_alerts_datasource"
	ApiIntegrationsDatasource                      feature = "snowflake_api_integrations_datasource"
	ApiIntegrationResource                         feature = "snowflake_api_integration_resource"
	ApiIntegrationAmazonApiGatewayResource         feature = "snowflake_api_integration_amazon_api_gateway_resource"
	ApiIntegrationAzureApiManagementResource       feature = "snowflake_api_integration_azure_api_management_resource"
	ApiIntegrationExternalMcpDynamicClientResource feature = "snowflake_api_integration_external_mcp_dynamic_client_resource"
	ApiIntegrationExternalMcpOAuth2Resource        feature = "snowflake_api_integration_external_mcp_oauth2_resource"
	ApiIntegrationGitRepositoryGithubAppResource   feature = "snowflake_api_integration_git_repository_github_app_resource"
	ApiIntegrationGitRepositoryOauth2Resource      feature = "snowflake_api_integration_git_repository_oauth2_resource"
	ApiIntegrationGitRepositoryPrivateLinkResource feature = "snowflake_api_integration_git_repository_private_link_resource"
	ApiIntegrationGitRepositoryTokenResource       feature = "snowflake_api_integration_git_repository_token_resource"
	ApiIntegrationGoogleCloudApiGatewayResource    feature = "snowflake_api_integration_google_cloud_api_gateway_resource"
	AuthenticationPolicyResource                   feature = "snowflake_authentication_policy_resource"
	AuthenticationPoliciesDatasource               feature = "snowflake_authentication_policies_datasource"
	CatalogIntegrationAwsGlueResource              feature = "snowflake_catalog_integration_aws_glue_resource"
	CatalogIntegrationObjectStorageResource        feature = "snowflake_catalog_integration_object_storage_resource"
	CatalogIntegrationOpenCatalogResource          feature = "snowflake_catalog_integration_open_catalog_resource"
	CatalogIntegrationIcebergRestResource          feature = "snowflake_catalog_integration_iceberg_rest_resource"
	CatalogIntegrationsDatasource                  feature = "snowflake_catalog_integrations_datasource"
	ComputePoolResource                            feature = "snowflake_compute_pool_resource"
	ComputePoolsDatasource                         feature = "snowflake_compute_pools_datasource"
	CortexAgentResource                            feature = "snowflake_cortex_agent_resource"
	CortexAgentsDatasource                         feature = "snowflake_cortex_agents_datasource"
	CortexSearchServiceResource                    feature = "snowflake_cortex_search_service_resource"
	CortexSearchServicesDatasource                 feature = "snowflake_cortex_search_services_datasource"
	CurrentAccountResource                         feature = "snowflake_current_account_resource"
	CurrentAccountDatasource                       feature = "snowflake_current_account_datasource"
	CurrentOrganizationAccountResource             feature = "snowflake_current_organization_account_resource"
	DatabaseDatasource                             feature = "snowflake_database_datasource"
	DatabaseRoleDatasource                         feature = "snowflake_database_role_datasource"
	DynamicTableResource                           feature = "snowflake_dynamic_table_resource"
	DynamicTablesDatasource                        feature = "snowflake_dynamic_tables_datasource"
	EmailNotificationIntegrationResource           feature = "snowflake_email_notification_integration_resource"
	ExternalAzureStageResource                     feature = "snowflake_stage_external_azure_resource"
	ExternalFunctionResource                       feature = "snowflake_external_function_resource"
	ExternalFunctionsDatasource                    feature = "snowflake_external_functions_datasource"
	ExternalGcsStageResource                       feature = "snowflake_stage_external_gcs_resource"
	ExternalS3StageResource                        feature = "snowflake_stage_external_s3_resource"
	ExternalS3CompatibleStageResource              feature = "snowflake_stage_external_s3_compatible_resource"
	ExternalTableResource                          feature = "snowflake_external_table_resource"
	ExternalTablesDatasource                       feature = "snowflake_external_tables_datasource"
	ExternalVolumeResource                         feature = "snowflake_external_volume_resource"
	ExternalVolumesDatasource                      feature = "snowflake_external_volumes_datasource"
	FailoverGroupResource                          feature = "snowflake_failover_group_resource"
	FailoverGroupsDatasource                       feature = "snowflake_failover_groups_datasource"
	FileFormatResource                             feature = "snowflake_file_format_resource"
	FileFormatsDatasource                          feature = "snowflake_file_formats_datasource"
	FunctionJavaResource                           feature = "snowflake_function_java_resource"
	FunctionJavascriptResource                     feature = "snowflake_function_javascript_resource"
	FunctionPythonResource                         feature = "snowflake_function_python_resource"
	FunctionScalaResource                          feature = "snowflake_function_scala_resource"
	FunctionSqlResource                            feature = "snowflake_function_sql_resource"
	FunctionsDatasource                            feature = "snowflake_functions_datasource"
	GitRepositoryResource                          feature = "snowflake_git_repository_resource"
	GitRepositoriesDatasource                      feature = "snowflake_git_repositories_datasource"
	// Currently, the hybrid table resource is added to the test provider.
	// HybridTableResource                           feature = "snowflake_hybrid_table_resource"
	IcebergTableFromDeltaFilesResource            feature = "snowflake_iceberg_table_from_delta_files_resource"
	IcebergTableFromFilesResource                 feature = "snowflake_iceberg_table_from_files_resource"
	ImageRepositoryResource                       feature = "snowflake_image_repository_resource"
	ImageRepositoriesDatasource                   feature = "snowflake_image_repositories_datasource"
	InternalStageResource                         feature = "snowflake_stage_internal_resource"
	JobServiceResource                            feature = "snowflake_job_service_resource"
	ListingResource                               feature = "snowflake_listing_resource"
	ListingsDatasource                            feature = "snowflake_listings_datasource"
	ManagedAccountResource                        feature = "snowflake_managed_account_resource"
	MaterializedViewResource                      feature = "snowflake_materialized_view_resource"
	MaterializedViewsDatasource                   feature = "snowflake_materialized_views_datasource"
	NetworkPolicyAttachmentResource               feature = "snowflake_network_policy_attachment_resource"
	NetworkRuleResource                           feature = "snowflake_network_rule_resource"
	NetworkRulesDatasource                        feature = "snowflake_network_rules_datasource"
	NotebookResource                              feature = "snowflake_notebook_resource"
	NotebooksDatasource                           feature = "snowflake_notebooks_datasource"
	NotificationIntegrationResource               feature = "snowflake_notification_integration_resource"
	ObjectParameterResource                       feature = "snowflake_object_parameter_resource"
	PasswordPoliciesDatasource                    feature = "snowflake_password_policies_datasource"
	PasswordPolicyResource                        feature = "snowflake_password_policy_resource"
	PipeResource                                  feature = "snowflake_pipe_resource"
	PipesDatasource                               feature = "snowflake_pipes_datasource"
	PostgresForkResource                          feature = "snowflake_postgres_fork_resource"
	PostgresInstanceResource                      feature = "snowflake_postgres_instance_resource"
	ProcedureJavaResource                         feature = "snowflake_procedure_java_resource"
	ProcedureJavascriptResource                   feature = "snowflake_procedure_javascript_resource"
	ProcedurePythonResource                       feature = "snowflake_procedure_python_resource"
	ProcedureScalaResource                        feature = "snowflake_procedure_scala_resource"
	ProcedureSqlResource                          feature = "snowflake_procedure_sql_resource"
	ProceduresDatasource                          feature = "snowflake_procedures_datasource"
	CurrentRoleDatasource                         feature = "snowflake_current_role_datasource"
	SemanticViewResource                          feature = "snowflake_semantic_view_resource"
	SemanticViewDatasource                        feature = "snowflake_semantic_views_datasource"
	SessionPoliciesDatasource                     feature = "snowflake_session_policies_datasource"
	SessionPolicyResource                         feature = "snowflake_session_policy_resource"
	ServiceResource                               feature = "snowflake_service_resource"
	ServicesDatasource                            feature = "snowflake_services_datasource"
	SequenceResource                              feature = "snowflake_sequence_resource"
	SequencesDatasource                           feature = "snowflake_sequences_datasource"
	ShareResource                                 feature = "snowflake_share_resource"
	SharesDatasource                              feature = "snowflake_shares_datasource"
	ParametersDatasource                          feature = "snowflake_parameters_datasource"
	StageResource                                 feature = "snowflake_stage_resource"
	StagesDatasource                              feature = "snowflake_stages_datasource"
	StorageIntegrationResource                    feature = "snowflake_storage_integration_resource"
	StorageIntegrationAwsResource                 feature = "snowflake_storage_integration_aws_resource"
	StorageIntegrationAzureResource               feature = "snowflake_storage_integration_azure_resource"
	StorageIntegrationGcsResource                 feature = "snowflake_storage_integration_gcs_resource"
	StorageIntegrationsDatasource                 feature = "snowflake_storage_integrations_datasource"
	StorageLifecyclePolicyResource                feature = "snowflake_storage_lifecycle_policy_resource"
	StorageLifecyclePoliciesDatasource            feature = "snowflake_storage_lifecycle_policies_datasource"
	SystemGenerateSCIMAccessTokenDatasource       feature = "snowflake_system_generate_scim_access_token_datasource"
	SystemGetAWSSNSIAMPolicyDatasource            feature = "snowflake_system_get_aws_sns_iam_policy_datasource"
	SystemGetPrivateLinkConfigDatasource          feature = "snowflake_system_get_privatelink_config_datasource"
	SystemGetSnowflakePlatformInfoDatasource      feature = "snowflake_system_get_snowflake_platform_info_datasource"
	TableResource                                 feature = "snowflake_table_resource"
	TablesDatasource                              feature = "snowflake_tables_datasource"
	TableColumnMaskingPolicyApplicationResource   feature = "snowflake_table_column_masking_policy_application_resource"
	TableConstraintResource                       feature = "snowflake_table_constraint_resource"
	TableStorageLifecyclePolicyAttachmentResource feature = "snowflake_table_storage_lifecycle_policy_attachment_resource"
	UserAuthenticationPolicyAttachmentResource    feature = "snowflake_user_authentication_policy_attachment_resource"
	UserPublicKeysResource                        feature = "snowflake_user_public_keys_resource"
	UserPasswordPolicyAttachmentResource          feature = "snowflake_user_password_policy_attachment_resource"
	UserProgrammaticAccessTokenResource           feature = "snowflake_user_programmatic_access_token_resource"
	UserSessionPolicyAttachmentResource           feature = "snowflake_user_session_policy_attachment_resource"
	UserProgrammaticAccessTokensDatasource        feature = "snowflake_user_programmatic_access_tokens_datasource"
	WarehouseAdaptiveResource                     feature = "snowflake_warehouse_adaptive_resource"
)

var allPreviewFeatures = []feature{
	AccountAuthenticationPolicyAttachmentResource,
	AccountPasswordPolicyAttachmentResource,
	AlertResource,
	AlertsDatasource,
	ApiIntegrationsDatasource,
	ApiIntegrationResource,
	ApiIntegrationAmazonApiGatewayResource,
	ApiIntegrationAzureApiManagementResource,
	ApiIntegrationExternalMcpDynamicClientResource,
	ApiIntegrationExternalMcpOAuth2Resource,
	ApiIntegrationGitRepositoryGithubAppResource,
	ApiIntegrationGitRepositoryOauth2Resource,
	ApiIntegrationGitRepositoryPrivateLinkResource,
	ApiIntegrationGitRepositoryTokenResource,
	ApiIntegrationGoogleCloudApiGatewayResource,
	CortexAgentResource,
	CortexAgentsDatasource,
	CortexSearchServiceResource,
	CortexSearchServicesDatasource,
	CurrentAccountDatasource,
	DatabaseDatasource,
	DatabaseRoleDatasource,
	DynamicTableResource,
	DynamicTablesDatasource,
	ExternalFunctionResource,
	ExternalFunctionsDatasource,
	ExternalTableResource,
	ExternalTablesDatasource,
	FailoverGroupResource,
	FailoverGroupsDatasource,
	FileFormatResource,
	FileFormatsDatasource,
	FunctionJavaResource,
	FunctionJavascriptResource,
	FunctionPythonResource,
	FunctionScalaResource,
	FunctionSqlResource,
	FunctionsDatasource,
	// Currently, the hybrid table resource is added to the test provider.
	// HybridTableResource,
	IcebergTableFromDeltaFilesResource,
	IcebergTableFromFilesResource,
	JobServiceResource,
	ListingsDatasource,
	ManagedAccountResource,
	MaterializedViewResource,
	MaterializedViewsDatasource,
	NetworkPolicyAttachmentResource,
	NetworkRulesDatasource,
	EmailNotificationIntegrationResource,
	NotificationIntegrationResource,
	ObjectParameterResource,
	PipeResource,
	PipesDatasource,
	// These resources are added to the test provider.
	// PostgresForkResource,
	// PostgresInstanceResource,
	CurrentRoleDatasource,
	SemanticViewResource,
	SemanticViewDatasource,
	SequenceResource,
	SequencesDatasource,
	ShareResource,
	SharesDatasource,
	ParametersDatasource,
	ProcedureJavaResource,
	ProcedureJavascriptResource,
	ProcedurePythonResource,
	ProcedureScalaResource,
	ProcedureSqlResource,
	ProceduresDatasource,
	StageResource,
	StagesDatasource,
	StorageIntegrationResource,
	StorageLifecyclePolicyResource,
	StorageLifecyclePoliciesDatasource,
	SystemGenerateSCIMAccessTokenDatasource,
	SystemGetAWSSNSIAMPolicyDatasource,
	SystemGetPrivateLinkConfigDatasource,
	SystemGetSnowflakePlatformInfoDatasource,
	TableColumnMaskingPolicyApplicationResource,
	TableConstraintResource,
	TableStorageLifecyclePolicyAttachmentResource,
	TableResource,
	TablesDatasource,
	UserAuthenticationPolicyAttachmentResource,
	UserPasswordPolicyAttachmentResource,
	UserPublicKeysResource,
	WarehouseAdaptiveResource,
}
var AllPreviewFeatures = sdk.AsStringList(allPreviewFeatures)

var promotedFeatures = []feature{
	AccountSessionPolicyAttachmentResource,
	AuthenticationPolicyResource,
	AuthenticationPoliciesDatasource,
	CatalogIntegrationAwsGlueResource,
	CatalogIntegrationObjectStorageResource,
	CatalogIntegrationOpenCatalogResource,
	CatalogIntegrationIcebergRestResource,
	CatalogIntegrationsDatasource,
	ComputePoolResource,
	ComputePoolsDatasource,
	CurrentAccountResource,
	CurrentOrganizationAccountResource,
	ExternalAzureStageResource,
	ExternalGcsStageResource,
	ExternalS3StageResource,
	ExternalS3CompatibleStageResource,
	ExternalVolumeResource,
	ExternalVolumesDatasource,
	GitRepositoryResource,
	GitRepositoriesDatasource,
	ImageRepositoryResource,
	ImageRepositoriesDatasource,
	InternalStageResource,
	ListingResource,
	NetworkRuleResource,
	NotebookResource,
	NotebooksDatasource,
	PasswordPoliciesDatasource,
	PasswordPolicyResource,
	ServiceResource,
	ServicesDatasource,
	SessionPoliciesDatasource,
	SessionPolicyResource,
	StorageIntegrationAwsResource,
	StorageIntegrationAzureResource,
	StorageIntegrationGcsResource,
	StorageIntegrationsDatasource,
	UserProgrammaticAccessTokenResource,
	UserProgrammaticAccessTokensDatasource,
	UserSessionPolicyAttachmentResource,
}
var PromotedFeatures = sdk.AsStringList(promotedFeatures)

var ValidPreviewFeatures = append(AllPreviewFeatures, PromotedFeatures...)

func EnsurePreviewFeatureEnabled(feat feature, enabledFeatures []string) error {
	if !slices.ContainsFunc(enabledFeatures, func(s string) bool {
		return s == string(feat)
	}) {
		return fmt.Errorf("%[1]s is currently a preview feature, and must be enabled by adding %[1]s to `preview_features_enabled` in Terraform configuration.", feat)
	}
	return nil
}

func StringToFeature(featRaw string) (feature, error) {
	feat := feature(strings.ToLower(featRaw))
	if !slices.Contains(ValidPreviewFeatures, string(feat)) {
		return "", fmt.Errorf("invalid feature: %s", featRaw)
	}
	return feat, nil
}

func GetPromotedFeatures(enabledFeatures []string) []string {
	containedPromotedFeatures := make([]string, 0)
	if enabledFeatures == nil {
		return containedPromotedFeatures
	}
	for _, enabledFeature := range enabledFeatures {
		if IsPromotedFeature(enabledFeature) {
			containedPromotedFeatures = append(containedPromotedFeatures, enabledFeature)
		}
	}
	return containedPromotedFeatures
}

func IsPromotedFeature(rawFeature string) bool {
	return slices.ContainsFunc(PromotedFeatures, func(s string) bool {
		return strings.EqualFold(rawFeature, s)
	})
}

type PreviewFeature interface {
	xxxProtected()
	String() string
}

func (f feature) xxxProtected() {}

func (f feature) String() string {
	return string(f)
}
