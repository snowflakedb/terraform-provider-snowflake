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
