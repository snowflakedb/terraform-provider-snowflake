package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := quotasTestIdSchemaObjectIdentifier
	warehouseId := randomAccountObjectIdentifier()
	cortexAgentId := randomSchemaObjectIdentifier()
	snowflakeIntelligenceId := randomAccountObjectIdentifier()
	aiGatewayId := randomAccountObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifier()
	integrationId := randomAccountObjectIdentifier()
	procedureId := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR)
	tagReference := func(tagId SchemaObjectIdentifier) string {
		return "(SELECT SYSTEM$REFERENCE('TAG', '" + tagId.FullyQualifiedName() + "', 'SESSION', 'APPLYBUDGET'))"
	}

	quotasTests.Create.
		withExpectedSqlf(
			case_Quotas_sql_Create_basic,
			"CREATE SNOWFLAKE.CORE.QUOTA %s ()", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Create_all,
			func(opts *CreateQuotaOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("test comment")
			},
			"CREATE SNOWFLAKE.CORE.QUOTA IF NOT EXISTS %s () COMMENT = 'test comment'", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateQuotaOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE SNOWFLAKE.CORE.QUOTA %s ()", id.FullyQualifiedName(),
		)

	newId := randomSchemaObjectIdentifier()
	quotasTests.Alter.
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Alter_RenameTo,
			func(opts *AlterQuotaOptions) { opts.RenameTo = &newId },
			"ALTER SNOWFLAKE.CORE.QUOTA %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Alter_Set,
			func(opts *AlterQuotaOptions) { opts.Set = &QuotaSet{Comment: new("quota comment")} },
			"ALTER SNOWFLAKE.CORE.QUOTA %s SET COMMENT = 'quota comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Alter_Unset,
			func(opts *AlterQuotaOptions) { opts.Unset = &QuotaUnset{Comment: new(true)} },
			"ALTER SNOWFLAKE.CORE.QUOTA %s UNSET COMMENT", id.FullyQualifiedName(),
		)

	quotasTests.Drop.
		withExpectedSqlf(
			case_Quotas_sql_Drop_basic,
			"DROP SNOWFLAKE.CORE.QUOTA %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Drop_all,
			func(opts *DropQuotaOptions) { opts.IfExists = new(true) },
			"DROP SNOWFLAKE.CORE.QUOTA IF EXISTS %s", id.FullyQualifiedName(),
		)

	quotasTests.Show.
		withExpectedSql(case_Quotas_sql_Show_basic, "SHOW SNOWFLAKE.CORE.QUOTA").
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Show_all,
			func(opts *ShowQuotaOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &In{Schema: id.SchemaId()}
				opts.Limit = &LimitFrom{Rows: new(10), From: new("A")}
			},
			"SHOW SNOWFLAKE.CORE.QUOTA LIKE '%s' IN SCHEMA %s LIMIT 10 FROM 'A'", id.Name(), id.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Show_Like,
			func(opts *ShowQuotaOptions) { opts.Like = &Like{Pattern: new(id.Name())} },
			"SHOW SNOWFLAKE.CORE.QUOTA LIKE '%s'", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Show_In,
			func(opts *ShowQuotaOptions) { opts.In = &In{Account: new(true)} },
			"SHOW SNOWFLAKE.CORE.QUOTA IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_Quotas_sql_Show_Limit,
			func(opts *ShowQuotaOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW SNOWFLAKE.CORE.QUOTA LIMIT 10",
		)

	quotasTests.AddSharedResource.
		withDefaultOpts(func() *AddSharedResourceQuotaOptions {
			return &AddSharedResourceQuotaOptions{
				name: id,
				args: QuotaAddSharedResourceArgs{Domain: QuotaSharedResourceDomainAiFunction},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_AddSharedResource_basic,
			"CALL %s!ADD_SHARED_RESOURCE ('AI FUNCTION')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_aiFunction",
			func(opts *AddSharedResourceQuotaOptions) { opts.args.InstanceName = new("AI_CLASSIFY") },
			"CALL %s!ADD_SHARED_RESOURCE ('AI FUNCTION', 'AI_CLASSIFY')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_cortexCode",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainCortexCode
				opts.args.InstanceName = new("CORTEX_CODE_CLI")
			},
			"CALL %s!ADD_SHARED_RESOURCE ('CORTEX CODE', 'CORTEX_CODE_CLI')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_warehouse",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainWarehouse
				opts.args.Warehouse = &warehouseId
			},
			"CALL %s!ADD_SHARED_RESOURCE ('WAREHOUSE', (SELECT SYSTEM$REFERENCE('WAREHOUSE', '%s')))", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_cortexAgent",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainCortexAgent
				opts.args.CortexAgent = &cortexAgentId
			},
			"CALL %s!ADD_SHARED_RESOURCE ('CORTEX AGENT', (SELECT SYSTEM$REFERENCE('CORTEX AGENT', '%s')))", id.FullyQualifiedName(), cortexAgentId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_snowflakeIntelligence",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainSnowflakeIntelligence
				opts.args.SnowflakeIntelligence = &snowflakeIntelligenceId
			},
			"CALL %s!ADD_SHARED_RESOURCE ('SNOWFLAKE INTELLIGENCE', (SELECT SYSTEM$REFERENCE('SNOWFLAKE INTELLIGENCE', '%s')))", id.FullyQualifiedName(), snowflakeIntelligenceId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddSharedResource_aiGateway",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainAiGateway
				opts.args.AiGateway = &aiGatewayId
			},
			"CALL %s!ADD_SHARED_RESOURCE ('AI GATEWAY', (SELECT SYSTEM$REFERENCE('AI GATEWAY', '%s')))", id.FullyQualifiedName(), aiGatewayId.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_AddSharedResource_instanceNameForReferenceDomain",
			func(opts *AddSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainWarehouse
				opts.args.InstanceName = new("SNOWADHOC")
			},
			errors.New("InstanceName can only be set for the AI FUNCTION, CORTEX CODE domains, got WAREHOUSE"),
		).
		withAdditionalValidationCase(
			"validation_AddSharedResource_domainNotSet",
			func(opts *AddSharedResourceQuotaOptions) { opts.args.Domain = "" },
			errors.New("invalid quota shared resource domain"),
		).
		withAdditionalValidationCase(
			"validation_AddSharedResource_referenceForOtherDomain",
			func(opts *AddSharedResourceQuotaOptions) { opts.args.Warehouse = &warehouseId },
			errors.New("Warehouse can only be set for the WAREHOUSE domain, got AI FUNCTION"),
		)

	quotasTests.RemoveSharedResource.
		withDefaultOpts(func() *RemoveSharedResourceQuotaOptions {
			return &RemoveSharedResourceQuotaOptions{
				name: id,
				args: QuotaRemoveSharedResourceArgs{Domain: QuotaSharedResourceDomainCortexAgent},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_RemoveSharedResource_basic,
			"CALL %s!REMOVE_SHARED_RESOURCE ('CORTEX AGENT')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RemoveSharedResource_instance",
			func(opts *RemoveSharedResourceQuotaOptions) { opts.args.CortexAgent = &cortexAgentId },
			"CALL %s!REMOVE_SHARED_RESOURCE ('CORTEX AGENT', (SELECT SYSTEM$REFERENCE('CORTEX AGENT', '%s')))", id.FullyQualifiedName(), cortexAgentId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RemoveSharedResource_warehouse",
			func(opts *RemoveSharedResourceQuotaOptions) {
				opts.args.Domain = QuotaSharedResourceDomainWarehouse
				opts.args.Warehouse = &warehouseId
			},
			"CALL %s!REMOVE_SHARED_RESOURCE ('WAREHOUSE', (SELECT SYSTEM$REFERENCE('WAREHOUSE', '%s')))", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_RemoveSharedResource_referenceForOtherDomain",
			func(opts *RemoveSharedResourceQuotaOptions) { opts.args.AiGateway = &aiGatewayId },
			errors.New("AiGateway can only be set for the AI GATEWAY domain, got CORTEX AGENT"),
		)

	quotasTests.SetUserTags.
		withDefaultOpts(func() *SetUserTagsQuotaOptions {
			return &SetUserTagsQuotaOptions{
				name: id,
				args: QuotaSetUserTagsArgs{Tags: []QuotaUserTag{{Tag: tagId, Value: "finance"}}, Operator: QuotaUserTagsOperatorUnion},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_SetUserTags_basic,
			"CALL %s!SET_USER_TAGS ([[%s, 'finance']], 'UNION')", id.FullyQualifiedName(), tagReference(tagId),
		).
		withAdditionalSqlCasef(
			"sql_SetUserTags_all",
			func(opts *SetUserTagsQuotaOptions) {
				opts.args.Tags = append(opts.args.Tags, QuotaUserTag{Tag: tagId2, Value: "analytics"})
				opts.args.Operator = QuotaUserTagsOperatorIntersection
			},
			"CALL %s!SET_USER_TAGS ([[%s, 'finance'], [%s, 'analytics']], 'INTERSECTION')", id.FullyQualifiedName(), tagReference(tagId), tagReference(tagId2),
		).
		withAdditionalSqlCasef(
			"sql_SetUserTags_empty",
			func(opts *SetUserTagsQuotaOptions) { opts.args.Tags = []QuotaUserTag{} },
			"CALL %s!SET_USER_TAGS ([], 'UNION')", id.FullyQualifiedName(),
		)

	quotasTests.ExcludeUsers.
		withDefaultOpts(func() *ExcludeUsersQuotaOptions {
			return &ExcludeUsersQuotaOptions{
				name: id,
				args: QuotaExcludeUsersArgs{ByUser: &QuotaExcludeUsersByUser{Users: []string{"CAROL"}}},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_ExcludeUsers_basic,
			"CALL %s!EXCLUDE_USERS ('USER', ['CAROL'])", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_ExcludeUsers_noUsers",
			func(opts *ExcludeUsersQuotaOptions) { opts.args.ByUser.Users = nil },
			"CALL %s!EXCLUDE_USERS ('USER', [])", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_ExcludeUsers_byTag",
			func(opts *ExcludeUsersQuotaOptions) {
				opts.args.ByUser = nil
				opts.args.ByTag = &QuotaExcludeUsersByTag{Tags: []QuotaExcludedUserTag{{Tag: tagId, Value: "high"}}}
			},
			"CALL %s!EXCLUDE_USERS ('TAG', [[%s, 'high']])", id.FullyQualifiedName(), tagReference(tagId),
		).
		withAdditionalSqlCasef(
			"sql_ExcludeUsers_noTags",
			func(opts *ExcludeUsersQuotaOptions) {
				opts.args.ByUser = nil
				opts.args.ByTag = &QuotaExcludeUsersByTag{}
			},
			"CALL %s!EXCLUDE_USERS ('TAG', [])", id.FullyQualifiedName(),
		)

	quotasTests.IncludeUsers.
		withDefaultOpts(func() *IncludeUsersQuotaOptions {
			return &IncludeUsersQuotaOptions{
				name: id,
				args: QuotaIncludeUsersArgs{Users: []string{"ALICE", "BOB"}},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_IncludeUsers_basic,
			"CALL %s!INCLUDE_USERS (['ALICE', 'BOB'])", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_IncludeUsers_empty",
			func(opts *IncludeUsersQuotaOptions) { opts.args.Users = nil },
			"CALL %s!INCLUDE_USERS ([])", id.FullyQualifiedName(),
		)

	quotasTests.GetQuotaScope.
		withExpectedSqlf(
			case_Quotas_sql_GetQuotaScope_basic,
			"CALL %s!GET_QUOTA_SCOPE ()", id.FullyQualifiedName(),
		)

	quotasTests.SetPerUserLimit.
		withDefaultOpts(func() *SetPerUserLimitQuotaOptions {
			return &SetPerUserLimitQuotaOptions{
				name: id,
				args: QuotaSetPerUserLimitArgs{Limit: 120},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_SetPerUserLimit_basic,
			"CALL %s!SET_PER_USER_LIMIT (120)", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetPerUserLimit_cycle",
			func(opts *SetPerUserLimitQuotaOptions) { opts.args.Cycle = new(QuotaCycleWeekly) },
			"CALL %s!SET_PER_USER_LIMIT (120, 'WEEKLY')", id.FullyQualifiedName(),
		)

	quotasTests.UnsetPerUserLimit.
		withDefaultOpts(func() *UnsetPerUserLimitQuotaOptions {
			return &UnsetPerUserLimitQuotaOptions{
				name: id,
				args: QuotaUnsetPerUserLimitArgs{Cycle: QuotaCycleDaily},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_UnsetPerUserLimit_basic,
			"CALL %s!UNSET_PER_USER_LIMIT ('DAILY')", id.FullyQualifiedName(),
		)

	quotasTests.SetBlockEnforcementEnabled.
		withDefaultOpts(func() *SetBlockEnforcementEnabledQuotaOptions {
			return &SetBlockEnforcementEnabledQuotaOptions{
				name: id,
				args: QuotaSetBlockEnforcementEnabledArgs{Enabled: true},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_SetBlockEnforcementEnabled_basic,
			"CALL %s!SET_BLOCK_ENFORCEMENT_ENABLED (true)", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetBlockEnforcementEnabled_all",
			func(opts *SetBlockEnforcementEnabledQuotaOptions) { opts.args.NotifyUsers = new(true) },
			"CALL %s!SET_BLOCK_ENFORCEMENT_ENABLED (true, true)", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetBlockEnforcementEnabled_disabled",
			func(opts *SetBlockEnforcementEnabledQuotaOptions) {
				opts.args.Enabled = false
				opts.args.NotifyUsers = new(false)
			},
			"CALL %s!SET_BLOCK_ENFORCEMENT_ENABLED (false, false)", id.FullyQualifiedName(),
		)

	quotasTests.GetConfig.
		withExpectedSqlf(
			case_Quotas_sql_GetConfig_basic,
			"CALL %s!GET_CONFIG ()", id.FullyQualifiedName(),
		)

	quotasTests.AddNotificationThreshold.
		withDefaultOpts(func() *AddNotificationThresholdQuotaOptions {
			return &AddNotificationThresholdQuotaOptions{
				name: id,
				args: QuotaAddNotificationThresholdArgs{
					Threshold:     80,
					SpendStrategy: QuotaSpendStrategyProjected,
					NotifyUser:    true,
					Cycle:         QuotaCycleMonthly,
				},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_AddNotificationThreshold_basic,
			"CALL %s!ADD_NOTIFICATION_THRESHOLD (80, 'PROJECTED', true, 'MONTHLY')", id.FullyQualifiedName(),
		)

	quotasTests.RemoveNotificationThreshold.
		withDefaultOpts(func() *RemoveNotificationThresholdQuotaOptions {
			return &RemoveNotificationThresholdQuotaOptions{
				name: id,
				args: QuotaRemoveNotificationThresholdArgs{Threshold: 80, SpendStrategy: QuotaSpendStrategyProjected},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_RemoveNotificationThreshold_basic,
			"CALL %s!REMOVE_NOTIFICATION_THRESHOLD (80, 'PROJECTED')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RemoveNotificationThreshold_cycle",
			func(opts *RemoveNotificationThresholdQuotaOptions) { opts.args.Cycle = new(QuotaCycleDaily) },
			"CALL %s!REMOVE_NOTIFICATION_THRESHOLD (80, 'PROJECTED', 'DAILY')", id.FullyQualifiedName(),
		)

	quotasTests.GetNotificationThresholds.
		withExpectedSqlf(
			case_Quotas_sql_GetNotificationThresholds_basic,
			"CALL %s!GET_NOTIFICATION_THRESHOLDS ()", id.FullyQualifiedName(),
		)

	quotasTests.AddNotificationIntegration.
		withDefaultOpts(func() *AddNotificationIntegrationQuotaOptions {
			return &AddNotificationIntegrationQuotaOptions{
				name: id,
				args: QuotaAddNotificationIntegrationArgs{Integration: integrationId},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_AddNotificationIntegration_basic,
			`CALL %s!ADD_NOTIFICATION_INTEGRATION ('\"%s\"')`, id.FullyQualifiedName(), integrationId.Name(),
		)

	quotasTests.RemoveNotificationIntegration.
		withDefaultOpts(func() *RemoveNotificationIntegrationQuotaOptions {
			return &RemoveNotificationIntegrationQuotaOptions{
				name: id,
				args: QuotaRemoveNotificationIntegrationArgs{Integration: integrationId},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_RemoveNotificationIntegration_basic,
			`CALL %s!REMOVE_NOTIFICATION_INTEGRATION ('\"%s\"')`, id.FullyQualifiedName(), integrationId.Name(),
		)

	quotasTests.GetNotificationIntegrations.
		withExpectedSqlf(
			case_Quotas_sql_GetNotificationIntegrations_basic,
			"CALL %s!GET_NOTIFICATION_INTEGRATIONS ()", id.FullyQualifiedName(),
		)

	quotasTests.SetAdminEmails.
		withDefaultOpts(func() *SetAdminEmailsQuotaOptions {
			return &SetAdminEmailsQuotaOptions{
				name: id,
				args: QuotaSetAdminEmailsArgs{Emails: "admin1@example.com, admin2@example.com"},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_SetAdminEmails_basic,
			"CALL %s!SET_ADMIN_EMAILS ('admin1@example.com, admin2@example.com')", id.FullyQualifiedName(),
		)

	quotasTests.AddCustomAction.
		withDefaultOpts(func() *AddCustomActionQuotaOptions {
			return &AddCustomActionQuotaOptions{
				name: id,
				args: QuotaAddCustomActionArgs{
					Procedure:     procedureId,
					SpendStrategy: QuotaSpendStrategyActual,
					Threshold:     100,
				},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_AddCustomAction_basic,
			"CALL %s!ADD_CUSTOM_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), [], 'ACTUAL', 100)", id.FullyQualifiedName(), procedureId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AddCustomAction_arguments",
			func(opts *AddCustomActionQuotaOptions) { opts.args.Arguments = []string{"arg1", "arg2"} },
			"CALL %s!ADD_CUSTOM_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), ['arg1', 'arg2'], 'ACTUAL', 100)", id.FullyQualifiedName(), procedureId.FullyQualifiedName(),
		)

	quotasTests.RemoveCustomActions.
		withDefaultOpts(func() *RemoveCustomActionsQuotaOptions {
			return &RemoveCustomActionsQuotaOptions{name: id}
		}).
		withExpectedSqlf(
			case_Quotas_sql_RemoveCustomActions_basic,
			"CALL %s!REMOVE_CUSTOM_ACTIONS ()", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RemoveCustomActions_threshold",
			func(opts *RemoveCustomActionsQuotaOptions) { opts.args.Threshold = new(100) },
			"CALL %s!REMOVE_CUSTOM_ACTIONS (100)", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RemoveCustomActions_procedure",
			func(opts *RemoveCustomActionsQuotaOptions) {
				opts.args.Threshold = new(100)
				opts.args.ProcedureFqn = new("COST_MGMT_DB.ACTIONS.DISABLE_USERS_BY_ID(VARCHAR)")
			},
			"CALL %s!REMOVE_CUSTOM_ACTIONS (100, 'COST_MGMT_DB.ACTIONS.DISABLE_USERS_BY_ID(VARCHAR)')", id.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_RemoveCustomActions_procedureWithoutThreshold",
			func(opts *RemoveCustomActionsQuotaOptions) {
				opts.args.ProcedureFqn = new("COST_MGMT_DB.ACTIONS.DISABLE_USERS_BY_ID(VARCHAR)")
			},
			errors.New("ProcedureFqn can only be set together with Threshold"),
		)

	quotasTests.GetCustomActions.
		withExpectedSqlf(
			case_Quotas_sql_GetCustomActions_basic,
			"CALL %s!GET_CUSTOM_ACTIONS ()", id.FullyQualifiedName(),
		)

	quotasTests.SetCycleStartAction.
		withDefaultOpts(func() *SetCycleStartActionQuotaOptions {
			return &SetCycleStartActionQuotaOptions{
				name: id,
				args: QuotaSetCycleStartActionArgs{Procedure: procedureId},
			}
		}).
		withExpectedSqlf(
			case_Quotas_sql_SetCycleStartAction_basic,
			"CALL %s!SET_CYCLE_START_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), [])", id.FullyQualifiedName(), procedureId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetCycleStartAction_arguments",
			func(opts *SetCycleStartActionQuotaOptions) { opts.args.Arguments = []string{"arg1"} },
			"CALL %s!SET_CYCLE_START_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), ['arg1'])", id.FullyQualifiedName(), procedureId.FullyQualifiedName(),
		)

	quotasTests.RemoveCycleStartAction.
		withExpectedSqlf(
			case_Quotas_sql_RemoveCycleStartAction_basic,
			"CALL %s!REMOVE_CYCLE_START_ACTION ()", id.FullyQualifiedName(),
		)

	quotasTests.GetCycleStartAction.
		withExpectedSqlf(
			case_Quotas_sql_GetCycleStartAction_basic,
			"CALL %s!GET_CYCLE_START_ACTION ()", id.FullyQualifiedName(),
		)
}

func TestParseQuotaScope(t *testing.T) {
	t.Run("captured output with instances and selected users", func(t *testing.T) {
		scope, err := ParseQuotaScope(`{
  "excluded_users": [
    {
      "user_id": 22,
      "user_name": "BOB"
    }
  ],
  "included_users": [
    {
      "user_id": 41,
      "user_name": "ALICE"
    },
    {
      "user_id": 42,
      "user_name": "carol"
    }
  ],
  "shared_resources": [
    {
      "domain": "AI_FUNCTION",
      "id": -1,
      "name": "[ALL-CORTEX AI FUNCTIONS]"
    },
    {
      "domain": "CORTEX_CODE",
      "id": 1,
      "name": "CORTEX_CODE_CLI"
    }
  ],
  "user_tags": {
    "operator": "INTERSECTION",
    "tags": [
      {
        "tagDatabase": "DB",
        "tagId": 11,
        "tagName": "COST_CENTER",
        "tagSchema": "S1",
        "tagValues": [
          "Analytics",
          "finance"
        ]
      }
    ]
  }
}`)
		require.NoError(t, err)

		assert.Equal(t, []QuotaScopeUser{{UserId: "22", UserName: "BOB"}}, scope.ExcludedUsers)
		assert.Equal(t, []QuotaScopeUser{{UserId: "41", UserName: "ALICE"}, {UserId: "42", UserName: "carol"}}, scope.IncludedUsers)

		require.Len(t, scope.SharedResources, 2)
		assert.True(t, scope.SharedResources[0].IsWholeDomain())
		domain, err := scope.SharedResources[0].NormalizedDomain()
		require.NoError(t, err)
		assert.Equal(t, QuotaSharedResourceDomainAiFunction, domain)
		assert.False(t, scope.SharedResources[1].IsWholeDomain())
		assert.Equal(t, "CORTEX_CODE_CLI", scope.SharedResources[1].Name)

		require.NotNil(t, scope.UserTags)
		assert.False(t, scope.UserTags.AreExclusions())
		require.Len(t, scope.UserTags.Tags, 1)
		assert.Equal(t, NewSchemaObjectIdentifier("DB", "S1", "COST_CENTER"), scope.UserTags.Tags[0].TagId())
		assert.Equal(t, []string{"Analytics", "finance"}, scope.UserTags.Tags[0].TagValues)
	})

	t.Run("captured output of a new quota", func(t *testing.T) {
		scope, err := ParseQuotaScope(`{
  "excluded_users": [],
  "included_users": [],
  "shared_resources": [],
  "user_tags": {
    "operator": "ALL_USERS",
    "tags": []
  }
}`)
		require.NoError(t, err)
		assert.Empty(t, scope.SharedResources)
		require.NotNil(t, scope.UserTags)
		assert.True(t, scope.UserTags.AreExclusions())
		assert.Empty(t, scope.UserTags.Tags)
	})

	t.Run("captured output with excluded tags", func(t *testing.T) {
		scope, err := ParseQuotaScope(`{"excluded_users": [], "included_users": [], "shared_resources": [], "user_tags": {"operator": "ALL_USERS", "tags": [{"tagDatabase": "DB", "tagId": 11, "tagName": "COST_CENTER", "tagSchema": "S1", "tagValues": ["contractor"]}]}}`)
		require.NoError(t, err)
		assert.True(t, scope.UserTags.AreExclusions())
		assert.Equal(t, []string{"contractor"}, scope.UserTags.Tags[0].TagValues)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := ParseQuotaScope(`not json`)
		require.ErrorContains(t, err, "failed to parse quota scope")
	})
}

func TestQuotaConfig_PerUserLimitFor(t *testing.T) {
	config := QuotaConfig{PerUserLimit: new(1000), PerUserLimitDaily: new(60)}

	monthly, err := config.PerUserLimitFor(QuotaCycleMonthly)
	require.NoError(t, err)
	assert.Equal(t, new(1000), monthly)
	weekly, err := config.PerUserLimitFor(QuotaCycleWeekly)
	require.NoError(t, err)
	assert.Nil(t, weekly)
	daily, err := config.PerUserLimitFor(QuotaCycleDaily)
	require.NoError(t, err)
	assert.Equal(t, new(60), daily)
	_, err = config.PerUserLimitFor("YEARLY")
	require.ErrorContains(t, err, "unsupported quota cycle")
}

func TestParseQuotaAdminEmails(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{input: "", expected: []string{}},
		{input: "admin1@example.com", expected: []string{"admin1@example.com"}},
		// as returned by GET_CONFIG
		{input: "admin1@example.com,admin2@example.com", expected: []string{"admin1@example.com", "admin2@example.com"}},
		{input: "admin1@example.com, admin2@example.com", expected: []string{"admin1@example.com", "admin2@example.com"}},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			emails, err := ParseQuotaAdminEmails(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, emails)
		})
	}
}

func TestParseQuotaProcedureArguments(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{input: "", expected: []string{}},
		{input: "[]", expected: []string{}},
		{input: "null", expected: []string{}},
		{input: `["a", "b"]`, expected: []string{"a", "b"}},
		// as returned by GET_CUSTOM_ACTIONS
		{input: "[\n  \"other\"\n]", expected: []string{"other"}},
		{input: "[\n  \"it's\"\n]", expected: []string{"it's"}},
		{input: "[a, b]", expected: []string{"a", "b"}},
		{input: "['a', 'b']", expected: []string{"a", "b"}},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			arguments, err := ParseQuotaProcedureArguments(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, arguments)
		})
	}
}

func TestNewQuotaSharedResourceArgsRequestFromInstance(t *testing.T) {
	testCases := []struct {
		domain   QuotaSharedResourceDomain
		instance string
		expected QuotaAddSharedResourceArgsRequest
	}{
		{domain: QuotaSharedResourceDomainAiFunction, instance: "", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainAiFunction}},
		{domain: QuotaSharedResourceDomainAiFunction, instance: "AI_CLASSIFY", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainAiFunction, InstanceName: new("AI_CLASSIFY")}},
		{domain: QuotaSharedResourceDomainCortexCode, instance: "CORTEX_CODE_CLI", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainCortexCode, InstanceName: new("CORTEX_CODE_CLI")}},
		{domain: QuotaSharedResourceDomainWarehouse, instance: "SNOWADHOC", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainWarehouse, Warehouse: new(NewAccountObjectIdentifier("SNOWADHOC"))}},
		{domain: QuotaSharedResourceDomainWarehouse, instance: `"snow adhoc"`, expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainWarehouse, Warehouse: new(NewAccountObjectIdentifier("snow adhoc"))}},
		{domain: QuotaSharedResourceDomainCortexAgent, instance: `"DB"."SCH"."AGENT"`, expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainCortexAgent, CortexAgent: new(NewSchemaObjectIdentifier("DB", "SCH", "AGENT"))}},
		{domain: QuotaSharedResourceDomainSnowflakeIntelligence, instance: "MY_SI", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainSnowflakeIntelligence, SnowflakeIntelligence: new(NewAccountObjectIdentifier("MY_SI"))}},
		{domain: QuotaSharedResourceDomainAiGateway, instance: "SNOWFLAKE", expected: QuotaAddSharedResourceArgsRequest{Domain: QuotaSharedResourceDomainAiGateway, AiGateway: new(NewAccountObjectIdentifier("SNOWFLAKE"))}},
	}
	for _, tc := range testCases {
		t.Run(string(tc.domain)+"/"+tc.instance, func(t *testing.T) {
			addRequest, err := NewQuotaAddSharedResourceArgsRequestFromInstance(tc.domain, tc.instance)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, *addRequest)

			removeRequest, err := NewQuotaRemoveSharedResourceArgsRequestFromInstance(tc.domain, tc.instance)
			require.NoError(t, err)
			assert.Equal(t, QuotaRemoveSharedResourceArgsRequest(tc.expected), *removeRequest)
		})
	}

	t.Run("invalid cortex agent identifier", func(t *testing.T) {
		_, err := NewQuotaAddSharedResourceArgsRequestFromInstance(QuotaSharedResourceDomainCortexAgent, "AGENT")
		require.Error(t, err)
	})

	t.Run("unsupported domain", func(t *testing.T) {
		_, err := NewQuotaAddSharedResourceArgsRequestFromInstance("UNKNOWN", "X")
		require.ErrorContains(t, err, "unsupported quota shared resource domain")
	})
}

func TestNewQuotaSetAdminEmailsArgsRequestFromEmails(t *testing.T) {
	request := NewQuotaSetAdminEmailsArgsRequestFromEmails("admin1@example.com", "admin2@example.com")
	assert.Equal(t, "admin1@example.com, admin2@example.com", request.Emails)
}
