package previewfeatures

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StringToFeature(t *testing.T) {
	type test struct {
		input string
		want  feature
	}

	valid := []test{
		// Case insensitive.
		{input: "SNOWFLAKE_CURRENT_ACCOUNT_DATASOURCE", want: CurrentAccountDatasource},

		// Supported Values.
		{input: "snowflake_account_password_policy_attachment_resource", want: AccountPasswordPolicyAttachmentResource},
		{input: "snowflake_alert_resource", want: AlertResource},
		{input: "snowflake_alerts_datasource", want: AlertsDatasource},
		{input: "snowflake_api_integration_resource", want: ApiIntegrationResource},
		{input: "snowflake_compute_pool_resource", want: ComputePoolResource},
		{input: "snowflake_compute_pools_datasource", want: ComputePoolsDatasource},
		{input: "snowflake_cortex_search_service_resource", want: CortexSearchServiceResource},
		{input: "snowflake_cortex_search_services_datasource", want: CortexSearchServicesDatasource},
		{input: "snowflake_current_account_resource", want: CurrentAccountResource},
		{input: "snowflake_current_account_datasource", want: CurrentAccountDatasource},
		{input: "snowflake_current_organization_account_resource", want: CurrentOrganizationAccountResource},
		{input: "snowflake_database_datasource", want: DatabaseDatasource},
		{input: "snowflake_database_role_datasource", want: DatabaseRoleDatasource},
		{input: "snowflake_dynamic_table_resource", want: DynamicTableResource},
		{input: "snowflake_dynamic_tables_datasource", want: DynamicTablesDatasource},
		{input: "snowflake_external_function_resource", want: ExternalFunctionResource},
		{input: "snowflake_external_functions_datasource", want: ExternalFunctionsDatasource},
		{input: "snowflake_external_table_resource", want: ExternalTableResource},
		{input: "snowflake_external_tables_datasource", want: ExternalTablesDatasource},
		{input: "snowflake_external_volume_resource", want: ExternalVolumeResource},
		{input: "snowflake_failover_group_resource", want: FailoverGroupResource},
		{input: "snowflake_failover_groups_datasource", want: FailoverGroupsDatasource},
		{input: "snowflake_file_format_resource", want: FileFormatResource},
		{input: "snowflake_file_formats_datasource", want: FileFormatsDatasource},
		{input: "snowflake_git_repository_resource", want: GitRepositoryResource},
		{input: "snowflake_git_repositories_datasource", want: GitRepositoriesDatasource},
		{input: "snowflake_image_repository_resource", want: ImageRepositoryResource},
		{input: "snowflake_image_repositories_datasource", want: ImageRepositoriesDatasource},
		{input: "snowflake_job_service_resource", want: JobServiceResource},
		{input: "snowflake_listing_resource", want: ListingResource},
		{input: "snowflake_managed_account_resource", want: ManagedAccountResource},
		{input: "snowflake_materialized_view_resource", want: MaterializedViewResource},
		{input: "snowflake_materialized_views_datasource", want: MaterializedViewsDatasource},
		{input: "snowflake_network_policy_attachment_resource", want: NetworkPolicyAttachmentResource},
		{input: "snowflake_network_rule_resource", want: NetworkRuleResource},
		{input: "snowflake_email_notification_integration_resource", want: EmailNotificationIntegrationResource},
		{input: "snowflake_notification_integration_resource", want: NotificationIntegrationResource},
		{input: "snowflake_object_parameter_resource", want: ObjectParameterResource},
		{input: "snowflake_password_policy_resource", want: PasswordPolicyResource},
		{input: "snowflake_pipe_resource", want: PipeResource},
		{input: "snowflake_pipes_datasource", want: PipesDatasource},
		{input: "snowflake_current_role_datasource", want: CurrentRoleDatasource},
		{input: "snowflake_service_resource", want: ServiceResource},
		{input: "snowflake_services_datasource", want: ServicesDatasource},
		{input: "snowflake_sequence_resource", want: SequenceResource},
		{input: "snowflake_sequences_datasource", want: SequencesDatasource},
		{input: "snowflake_share_resource", want: ShareResource},
		{input: "snowflake_shares_datasource", want: SharesDatasource},
		{input: "snowflake_parameters_datasource", want: ParametersDatasource},
		{input: "snowflake_stage_resource", want: StageResource},
		{input: "snowflake_stages_datasource", want: StagesDatasource},
		{input: "snowflake_storage_integration_resource", want: StorageIntegrationResource},
		{input: "snowflake_storage_integrations_datasource", want: StorageIntegrationsDatasource},
		{input: "snowflake_system_generate_scim_access_token_datasource", want: SystemGenerateSCIMAccessTokenDatasource},
		{input: "snowflake_system_get_aws_sns_iam_policy_datasource", want: SystemGetAWSSNSIAMPolicyDatasource},
		{input: "snowflake_system_get_privatelink_config_datasource", want: SystemGetPrivateLinkConfigDatasource},
		{input: "snowflake_system_get_snowflake_platform_info_datasource", want: SystemGetSnowflakePlatformInfoDatasource},
		{input: "snowflake_table_column_masking_policy_application_resource", want: TableColumnMaskingPolicyApplicationResource},
		{input: "snowflake_table_constraint_resource", want: TableConstraintResource},
		{input: "snowflake_user_authentication_policy_attachment_resource", want: UserAuthenticationPolicyAttachmentResource},
		{input: "snowflake_user_public_keys_resource", want: UserPublicKeysResource},
		{input: "snowflake_user_password_policy_attachment_resource", want: UserPasswordPolicyAttachmentResource},
		{input: "snowflake_user_programmatic_access_token_resource", want: UserProgrammaticAccessTokenResource},
		{input: "snowflake_user_programmatic_access_tokens_datasource", want: UserProgrammaticAccessTokensDatasource},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := StringToFeature(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := StringToFeature(tc.input)
			require.ErrorContains(t, err, "invalid feature")
		})
	}
}

func Test_GetPromotedFeatures(t *testing.T) {
	type test struct {
		enabledList []string
		expected    []string
	}

	promotedFeature1 := string(GitRepositoryResource)
	promotedFeature2 := string(GitRepositoriesDatasource)
	previewFeature1 := string(TableResource)
	previewFeature2 := string(TablesDatasource)

	previewOnly := []string{previewFeature1, previewFeature2}
	promotedOnly := []string{promotedFeature1, promotedFeature2}
	mixed1 := []string{promotedFeature1, previewFeature2, promotedFeature2}
	mixed2 := []string{previewFeature1, promotedFeature1, previewFeature2}
	var empty []string

	valid := []test{
		{enabledList: nil, expected: empty},
		{enabledList: empty, expected: empty},
		{enabledList: previewOnly, expected: empty},
		{enabledList: promotedOnly, expected: promotedOnly},
		{enabledList: mixed1, expected: promotedOnly},
		{enabledList: mixed2, expected: []string{promotedFeature1}},
	}

	for _, tc := range valid {
		t.Run(fmt.Sprintf("Enabled list: %v, expected promoted list: %v", tc.enabledList, tc.expected), func(t *testing.T) {
			promoted := GetPromotedFeatures(tc.enabledList)
			require.ElementsMatch(t, tc.expected, promoted)
		})
	}
}

func Test_IsPromotedFeature(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	promotedFeature1 := string(GitRepositoryResource)
	previewFeature1 := string(TableResource)
	unknownFeature := string(TableResource)

	valid := []test{
		{input: promotedFeature1, expected: true},
		{input: previewFeature1, expected: false},
		{input: unknownFeature, expected: false},
		{input: "", expected: false},
	}

	for _, tc := range valid {
		t.Run(fmt.Sprintf("Feature: %s, expected: %t", tc.input, tc.expected), func(t *testing.T) {
			got := IsPromotedFeature(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}
