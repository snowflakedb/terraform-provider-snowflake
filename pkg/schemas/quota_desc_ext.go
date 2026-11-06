package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func quotaComputedValue(valueType schema.ValueType) *schema.Schema {
	return &schema.Schema{Type: valueType, Computed: true}
}

func quotaComputedStringList() *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}}
}

func quotaComputedNestedList(nested map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: nested}}
}

func (quotaDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scope": quotaComputedNestedList(map[string]*schema.Schema{
			"shared_resources": quotaComputedNestedList(map[string]*schema.Schema{
				"domain": quotaComputedValue(schema.TypeString),
				"id":     quotaComputedValue(schema.TypeInt),
				"name":   quotaComputedValue(schema.TypeString),
			}),
			"included_users": quotaComputedStringList(),
			"excluded_users": quotaComputedStringList(),
			// ALL_USERS when the users are not selected by tags; the tags are then the excluded ones
			"user_tags_operator": quotaComputedValue(schema.TypeString),
			"user_tags": quotaComputedNestedList(map[string]*schema.Schema{
				"tag":    quotaComputedValue(schema.TypeString),
				"values": quotaComputedStringList(),
			}),
		}),
		"config": quotaComputedNestedList(map[string]*schema.Schema{
			"per_user_limit":                       quotaComputedValue(schema.TypeInt),
			"per_user_limit_weekly":                quotaComputedValue(schema.TypeInt),
			"per_user_limit_daily":                 quotaComputedValue(schema.TypeInt),
			"block_enforcement_enabled":            quotaComputedValue(schema.TypeBool),
			"per_user_block_notifications_enabled": quotaComputedValue(schema.TypeBool),
			"admin_emails":                         quotaComputedStringList(),
			"refresh_tier":                         quotaComputedValue(schema.TypeString),
		}),
		"notification_thresholds": quotaComputedNestedList(map[string]*schema.Schema{
			"threshold":       quotaComputedValue(schema.TypeInt),
			"spend_strategy":  quotaComputedValue(schema.TypeString),
			"notify_user":     quotaComputedValue(schema.TypeBool),
			"cycle":           quotaComputedValue(schema.TypeString),
			"added_timestamp": quotaComputedValue(schema.TypeString),
		}),
		"notification_integrations": quotaComputedNestedList(map[string]*schema.Schema{
			"integration_name":       quotaComputedValue(schema.TypeString),
			"last_notification_time": quotaComputedValue(schema.TypeString),
			"added_timestamp":        quotaComputedValue(schema.TypeString),
		}),
		"custom_actions": quotaComputedNestedList(map[string]*schema.Schema{
			"action_id":                 quotaComputedValue(schema.TypeString),
			"procedure":                 quotaComputedValue(schema.TypeString),
			"procedure_arguments":       quotaComputedStringList(),
			"spend_strategy":            quotaComputedValue(schema.TypeString),
			"threshold":                 quotaComputedValue(schema.TypeInt),
			"last_trigger_attempt_time": quotaComputedValue(schema.TypeString),
			"added_timestamp":           quotaComputedValue(schema.TypeString),
		}),
		"cycle_start_action": quotaComputedNestedList(map[string]*schema.Schema{
			"action_id":                 quotaComputedValue(schema.TypeString),
			"procedure":                 quotaComputedValue(schema.TypeString),
			"procedure_arguments":       quotaComputedStringList(),
			"added_timestamp":           quotaComputedValue(schema.TypeString),
			"last_trigger_attempt_time": quotaComputedValue(schema.TypeString),
		}),
	}
}

func (quotaDetailsToSchemaMapper) additionalToSchema(src *sdk.QuotaDetails, dst map[string]any) {
	dst["scope"] = []map[string]any{quotaScopeToSchema(src.Scope)}
	dst["config"] = []map[string]any{quotaConfigToSchema(src.Config)}

	thresholds := make([]map[string]any, len(src.NotificationThresholds))
	for i, threshold := range src.NotificationThresholds {
		thresholds[i] = map[string]any{
			"threshold":       threshold.Threshold,
			"spend_strategy":  string(threshold.SpendStrategy),
			"notify_user":     threshold.NotifyUser,
			"cycle":           string(threshold.Cycle),
			"added_timestamp": threshold.AddedTimestamp.String(),
		}
	}
	dst["notification_thresholds"] = thresholds

	integrations := make([]map[string]any, len(src.NotificationIntegrations))
	for i, integration := range src.NotificationIntegrations {
		integrationSchema := map[string]any{
			"integration_name": integration.IntegrationName,
			"added_timestamp":  integration.AddedTimestamp.String(),
		}
		if integration.LastNotificationTime != nil {
			integrationSchema["last_notification_time"] = *integration.LastNotificationTime
		}
		integrations[i] = integrationSchema
	}
	dst["notification_integrations"] = integrations

	actions := make([]map[string]any, len(src.CustomActions))
	for i, action := range src.CustomActions {
		actionSchema := map[string]any{
			"action_id":           action.ActionId,
			"procedure":           action.ProcedureFqn,
			"procedure_arguments": action.ProcedureArgs,
			"spend_strategy":      string(action.SpendStrategy),
			"threshold":           action.Threshold,
			"added_timestamp":     action.AddedTimestamp.String(),
		}
		if action.LastTriggerAttemptTime != nil {
			actionSchema["last_trigger_attempt_time"] = action.LastTriggerAttemptTime.String()
		}
		actions[i] = actionSchema
	}
	dst["custom_actions"] = actions

	cycleStartAction := make([]map[string]any, 0, 1)
	if action := src.CycleStartAction; action != nil {
		actionSchema := map[string]any{"procedure_arguments": action.ProcedureArgs}
		if action.ActionId != nil {
			actionSchema["action_id"] = *action.ActionId
		}
		if action.ProcedureFqn != nil {
			actionSchema["procedure"] = *action.ProcedureFqn
		}
		if action.AddedTimestamp != nil {
			actionSchema["added_timestamp"] = action.AddedTimestamp.String()
		}
		if action.LastTriggerAttemptTime != nil {
			actionSchema["last_trigger_attempt_time"] = action.LastTriggerAttemptTime.String()
		}
		cycleStartAction = append(cycleStartAction, actionSchema)
	}
	dst["cycle_start_action"] = cycleStartAction
}

func quotaScopeToSchema(scope sdk.QuotaScope) map[string]any {
	sharedResources := make([]map[string]any, len(scope.SharedResources))
	for i, resource := range scope.SharedResources {
		// GET_QUOTA_SCOPE returns domains with underscores (e.g. AI_FUNCTION); they're shown in the form used by the configuration.
		domain := resource.Domain
		if normalized, err := resource.NormalizedDomain(); err == nil {
			domain = string(normalized)
		}
		resourceSchema := map[string]any{
			"domain": domain,
			"name":   resource.Name,
		}
		if id, err := resource.Id.Int64(); err == nil {
			resourceSchema["id"] = int(id)
		}
		sharedResources[i] = resourceSchema
	}
	userNames := func(users []sdk.QuotaScopeUser) []string {
		names := make([]string, len(users))
		for i, user := range users {
			names[i] = user.UserName
		}
		return names
	}
	scopeSchema := map[string]any{
		"shared_resources": sharedResources,
		"included_users":   userNames(scope.IncludedUsers),
		"excluded_users":   userNames(scope.ExcludedUsers),
	}
	userTags := make([]map[string]any, 0)
	if scope.UserTags != nil {
		scopeSchema["user_tags_operator"] = scope.UserTags.Operator
		for _, tag := range scope.UserTags.Tags {
			userTags = append(userTags, map[string]any{
				"tag":    tag.TagId().FullyQualifiedName(),
				"values": tag.TagValues,
			})
		}
	}
	scopeSchema["user_tags"] = userTags
	return scopeSchema
}

func quotaConfigToSchema(config sdk.QuotaConfig) map[string]any {
	configSchema := make(map[string]any)
	for key, cycle := range map[string]sdk.QuotaCycle{
		"per_user_limit":        sdk.QuotaCycleMonthly,
		"per_user_limit_weekly": sdk.QuotaCycleWeekly,
		"per_user_limit_daily":  sdk.QuotaCycleDaily,
	} {
		if limit, err := config.PerUserLimitFor(cycle); err == nil && limit != nil {
			configSchema[key] = *limit
		}
	}
	configSchema["block_enforcement_enabled"] = config.BlockEnforcementEnabled
	configSchema["per_user_block_notifications_enabled"] = config.PerUserBlockNotificationsEnabled
	configSchema["admin_emails"] = config.AdminEmails
	if config.RefreshTier != nil {
		configSchema["refresh_tier"] = *config.RefreshTier
	}
	return configSchema
}
