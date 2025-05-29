package previewfeatures

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type feature string

const (
	CurrentAccountDatasource                      feature = "snowflake_current_account_datasource"
	AccountAuthenticationPolicyAttachmentResource feature = "snowflake_account_authentication_policy_attachment_resource"
	AccountPasswordPolicyAttachmentResource       feature = "snowflake_account_password_policy_attachment_resource"
	AlertResource                                 feature = "snowflake_alert_resource"
	AlertsDatasource                              feature = "snowflake_alerts_datasource"
	ApiIntegrationResource                        feature = "snowflake_api_integration_resource"
	AuthenticationPolicyResource                  feature = "snowflake_authentication_policy_resource"
	ComputePoolResource                           feature = "snowflake_compute_pool_resource"
	CortexSearchServiceResource                   feature = "snowflake_cortex_search_service_resource"
	CortexSearchServicesDatasource                feature = "snowflake_cortex_search_services_datasource"
	DatabaseDatasource                            feature = "snowflake_database_datasource"
	DatabaseRoleDatasource                        feature = "snowflake_database_role_datasource"
	DynamicTableResource                          feature = "snowflake_dynamic_table_resource"
	DynamicTablesDatasource                       feature = "snowflake_dynamic_tables_datasource"
	ExternalFunctionResource                      feature = "snowflake_external_function_resource"
	ExternalFunctionsDatasource                   feature = "snowflake_external_functions_datasource"
	ExternalTableResource                         feature = "snowflake_external_table_resource"
	ExternalTablesDatasource                      feature = "snowflake_external_tables_datasource"
	ExternalVolumeResource                        feature = "snowflake_external_volume_resource"
	FailoverGroupResource                         feature = "snowflake_failover_group_resource"
	FailoverGroupsDatasource                      feature = "snowflake_failover_groups_datasource"
	FileFormatResource                            feature = "snowflake_file_format_resource"
	FileFormatsDatasource                         feature = "snowflake_file_formats_datasource"
	FunctionJavaResource                          feature = "snowflake_function_java_resource"
	FunctionJavascriptResource                    feature = "snowflake_function_javascript_resource"
	FunctionPythonResource                        feature = "snowflake_function_python_resource"
	FunctionScalaResource                         feature = "snowflake_function_scala_resource"
	FunctionSqlResource                           feature = "snowflake_function_sql_resource"
	FunctionsDatasource                           feature = "snowflake_functions_datasource"
	ImageRepositoryResource                       feature = "snowflake_image_repository_resource"
	ImageRepositoriesDatasource                   feature = "snowflake_image_repositories_datasource"
	ManagedAccountResource                        feature = "snowflake_managed_account_resource"
	MaterializedViewResource                      feature = "snowflake_materialized_view_resource"
	MaterializedViewsDatasource                   feature = "snowflake_materialized_views_datasource"
	NetworkPolicyAttachmentResource               feature = "snowflake_network_policy_attachment_resource"
	NetworkRuleResource                           feature = "snowflake_network_rule_resource"
	EmailNotificationIntegrationResource          feature = "snowflake_email_notification_integration_resource"
	NotificationIntegrationResource               feature = "snowflake_notification_integration_resource"
	ObjectParameterResource                       feature = "snowflake_object_parameter_resource"
	PasswordPolicyResource                        feature = "snowflake_password_policy_resource"
	PipeResource                                  feature = "snowflake_pipe_resource"
	PipesDatasource                               feature = "snowflake_pipes_datasource"
	ProcedureJavaResource                         feature = "snowflake_procedure_java_resource"
	ProcedureJavascriptResource                   feature = "snowflake_procedure_javascript_resource"
	ProcedurePythonResource                       feature = "snowflake_procedure_python_resource"
	ProcedureScalaResource                        feature = "snowflake_procedure_scala_resource"
	ProcedureSqlResource                          feature = "snowflake_procedure_sql_resource"
	ProceduresDatasource                          feature = "snowflake_procedures_datasource"
	CurrentRoleDatasource                         feature = "snowflake_current_role_datasource"
	SequenceResource                              feature = "snowflake_sequence_resource"
	SequencesDatasource                           feature = "snowflake_sequences_datasource"
	ShareResource                                 feature = "snowflake_share_resource"
	SharesDatasource                              feature = "snowflake_shares_datasource"
	ParametersDatasource                          feature = "snowflake_parameters_datasource"
	StageResource                                 feature = "snowflake_stage_resource"
	StagesDatasource                              feature = "snowflake_stages_datasource"
	StorageIntegrationResource                    feature = "snowflake_storage_integration_resource"
	StorageIntegrationsDatasource                 feature = "snowflake_storage_integrations_datasource"
	SystemGenerateSCIMAccessTokenDatasource       feature = "snowflake_system_generate_scim_access_token_datasource"
	SystemGetAWSSNSIAMPolicyDatasource            feature = "snowflake_system_get_aws_sns_iam_policy_datasource"
	SystemGetPrivateLinkConfigDatasource          feature = "snowflake_system_get_privatelink_config_datasource"
	SystemGetSnowflakePlatformInfoDatasource      feature = "snowflake_system_get_snowflake_platform_info_datasource"
	TableResource                                 feature = "snowflake_table_resource"
	TablesDatasource                              feature = "snowflake_tables_datasource"
	TableColumnMaskingPolicyApplicationResource   feature = "snowflake_table_column_masking_policy_application_resource"
	TableConstraintResource                       feature = "snowflake_table_constraint_resource"
	UserAuthenticationPolicyAttachmentResource    feature = "snowflake_user_authentication_policy_attachment_resource"
	UserPublicKeysResource                        feature = "snowflake_user_public_keys_resource"
	UserPasswordPolicyAttachmentResource          feature = "snowflake_user_password_policy_attachment_resource"
)

var allPreviewFeatures = []feature{
	CurrentAccountDatasource,
	AccountAuthenticationPolicyAttachmentResource,
	AccountPasswordPolicyAttachmentResource,
	AlertResource,
	AlertsDatasource,
	ApiIntegrationResource,
	AuthenticationPolicyResource,
	ComputePoolResource,
	CortexSearchServiceResource,
	CortexSearchServicesDatasource,
	DatabaseDatasource,
	DatabaseRoleDatasource,
	DynamicTableResource,
	DynamicTablesDatasource,
	ExternalFunctionResource,
	ExternalFunctionsDatasource,
	ExternalTableResource,
	ExternalTablesDatasource,
	ExternalVolumeResource,
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
	ImageRepositoryResource,
	ImageRepositoriesDatasource,
	ManagedAccountResource,
	MaterializedViewResource,
	MaterializedViewsDatasource,
	NetworkPolicyAttachmentResource,
	NetworkRuleResource,
	EmailNotificationIntegrationResource,
	NotificationIntegrationResource,
	ObjectParameterResource,
	PasswordPolicyResource,
	PipeResource,
	PipesDatasource,
	CurrentRoleDatasource,
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
	StorageIntegrationsDatasource,
	SystemGenerateSCIMAccessTokenDatasource,
	SystemGetAWSSNSIAMPolicyDatasource,
	SystemGetPrivateLinkConfigDatasource,
	SystemGetSnowflakePlatformInfoDatasource,
	TableColumnMaskingPolicyApplicationResource,
	TableConstraintResource,
	TableResource,
	TablesDatasource,
	UserAuthenticationPolicyAttachmentResource,
	UserPublicKeysResource,
	UserPasswordPolicyAttachmentResource,
}
var AllPreviewFeatures = sdk.AsStringList(allPreviewFeatures)

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
	if !slices.Contains(allPreviewFeatures, feat) {
		return "", fmt.Errorf("invalid feature: %s", featRaw)
	}
	return feat, nil
}
