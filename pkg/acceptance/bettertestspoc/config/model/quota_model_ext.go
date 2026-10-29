package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type QuotaSharedResourceModel struct {
	Domain   sdk.QuotaSharedResourceDomain
	Instance string
}

type QuotaTagValueModel struct {
	Tag   sdk.SchemaObjectIdentifier
	Value string
}

type QuotaNotificationThresholdModel struct {
	Threshold     int
	SpendStrategy sdk.QuotaSpendStrategy
	NotifyUser    bool
	Cycle         sdk.QuotaCycle
}

type QuotaCustomActionModel struct {
	Procedure     sdk.SchemaObjectIdentifierWithArguments
	Arguments     []string
	SpendStrategy sdk.QuotaSpendStrategy
	Threshold     int
}

func QuotaFromId(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
) *QuotaModel {
	return Quota(resourceName, id.DatabaseName(), id.SchemaName(), id.Name())
}

func quotaStringSet(values []string) tfconfig.Variable {
	if len(values) == 0 {
		return config.EmptyListVariable()
	}
	return tfconfig.SetVariable(collections.Map(values, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...)
}

func quotaStringList(values []string) tfconfig.Variable {
	if len(values) == 0 {
		return config.EmptyListVariable()
	}
	return tfconfig.ListVariable(collections.Map(values, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...)
}

func quotaTagValues(tags []QuotaTagValueModel) []tfconfig.Variable {
	return collections.Map(tags, func(tag QuotaTagValueModel) tfconfig.Variable {
		return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"tag":   tfconfig.StringVariable(tag.Tag.FullyQualifiedName()),
			"value": tfconfig.StringVariable(tag.Value),
		})
	})
}

func (q *QuotaModel) WithSharedResources(resources ...QuotaSharedResourceModel) *QuotaModel {
	return q.WithSharedResourceValue(tfconfig.SetVariable(collections.Map(resources, func(resource QuotaSharedResourceModel) tfconfig.Variable {
		values := map[string]tfconfig.Variable{"domain": tfconfig.StringVariable(string(resource.Domain))}
		if resource.Instance != "" {
			values["instance"] = tfconfig.StringVariable(resource.Instance)
		}
		return tfconfig.ObjectVariable(values)
	})...))
}

func (q *QuotaModel) WithUserTags(operator sdk.QuotaUserTagsOperator, tags ...QuotaTagValueModel) *QuotaModel {
	return q.WithUserTagsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"operator": tfconfig.StringVariable(string(operator)),
		"tag":      tfconfig.SetVariable(quotaTagValues(tags)...),
	})))
}

func (q *QuotaModel) WithExcludedUserTags(tags ...QuotaTagValueModel) *QuotaModel {
	return q.WithExcludedUserTagValue(tfconfig.SetVariable(quotaTagValues(tags)...))
}

func (q *QuotaModel) WithIncludedUsers(users ...string) *QuotaModel {
	return q.WithIncludedUsersValue(quotaStringSet(users))
}

func (q *QuotaModel) WithExcludedUsers(users ...string) *QuotaModel {
	return q.WithExcludedUsersValue(quotaStringSet(users))
}

func (q *QuotaModel) WithNotificationThresholds(thresholds ...QuotaNotificationThresholdModel) *QuotaModel {
	return q.WithNotificationThresholdValue(tfconfig.SetVariable(collections.Map(thresholds, func(threshold QuotaNotificationThresholdModel) tfconfig.Variable {
		return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"threshold":      tfconfig.IntegerVariable(threshold.Threshold),
			"spend_strategy": tfconfig.StringVariable(string(threshold.SpendStrategy)),
			"notify_user":    tfconfig.BoolVariable(threshold.NotifyUser),
			"cycle":          tfconfig.StringVariable(string(threshold.Cycle)),
		})
	})...))
}

func (q *QuotaModel) WithNotificationIntegrations(integrations ...sdk.AccountObjectIdentifier) *QuotaModel {
	return q.WithNotificationIntegrationsValue(quotaStringSet(collections.Map(integrations, sdk.AccountObjectIdentifier.Name)))
}

func (q *QuotaModel) WithAdminEmails(emails ...string) *QuotaModel {
	return q.WithAdminEmailsValue(quotaStringSet(emails))
}

func (q *QuotaModel) WithCustomActions(actions ...QuotaCustomActionModel) *QuotaModel {
	return q.WithCustomActionValue(tfconfig.SetVariable(collections.Map(actions, func(action QuotaCustomActionModel) tfconfig.Variable {
		return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"procedure":      tfconfig.StringVariable(action.Procedure.FullyQualifiedName()),
			"arguments":      quotaStringList(action.Arguments),
			"spend_strategy": tfconfig.StringVariable(string(action.SpendStrategy)),
			"threshold":      tfconfig.IntegerVariable(action.Threshold),
		})
	})...))
}

func (q *QuotaModel) WithCycleStartAction(procedure sdk.SchemaObjectIdentifierWithArguments, arguments ...string) *QuotaModel {
	return q.WithCycleStartActionValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"procedure": tfconfig.StringVariable(procedure.FullyQualifiedName()),
		"arguments": quotaStringList(arguments),
	})))
}
