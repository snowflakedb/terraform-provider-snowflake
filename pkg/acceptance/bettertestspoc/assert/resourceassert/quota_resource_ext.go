package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (q *QuotaResourceAssert) HasSharedResourcesCount(expected int) *QuotaResourceAssert {
	q.AddAssertion(assert.ValueSet("shared_resource.#", strconv.Itoa(expected)))
	return q
}

func (q *QuotaResourceAssert) HasSharedResource(domain sdk.QuotaSharedResourceDomain, instance string) *QuotaResourceAssert {
	q.AddAssertion(assert.SetElemNested("shared_resource", map[string]string{"domain": string(domain), "instance": instance}))
	return q
}

func (q *QuotaResourceAssert) HasUserTags(operator sdk.QuotaUserTagsOperator, tagCount int) *QuotaResourceAssert {
	q.AddAssertion(assert.ValueSet("user_tags.#", "1"))
	q.AddAssertion(assert.ValueSet("user_tags.0.operator", string(operator)))
	q.AddAssertion(assert.ValueSet("user_tags.0.tag.#", strconv.Itoa(tagCount)))
	return q
}

func (q *QuotaResourceAssert) HasUserTag(tag sdk.SchemaObjectIdentifier, value string) *QuotaResourceAssert {
	q.AddAssertion(assert.SetElemNested("user_tags.0.tag", map[string]string{"tag": tag.FullyQualifiedName(), "value": value}))
	return q
}

func (q *QuotaResourceAssert) HasExcludedUserTag(tag sdk.SchemaObjectIdentifier, value string) *QuotaResourceAssert {
	q.AddAssertion(assert.SetElemNested("excluded_user_tag", map[string]string{"tag": tag.FullyQualifiedName(), "value": value}))
	return q
}

func (q *QuotaResourceAssert) HasNotificationThresholdsCount(expected int) *QuotaResourceAssert {
	q.AddAssertion(assert.ValueSet("notification_threshold.#", strconv.Itoa(expected)))
	return q
}

func (q *QuotaResourceAssert) HasNotificationThreshold(threshold int, spendStrategy sdk.QuotaSpendStrategy, notifyUser bool, cycle sdk.QuotaCycle) *QuotaResourceAssert {
	q.AddAssertion(assert.SetElemNested("notification_threshold", map[string]string{
		"threshold":      strconv.Itoa(threshold),
		"spend_strategy": string(spendStrategy),
		"notify_user":    strconv.FormatBool(notifyUser),
		"cycle":          string(cycle),
	}))
	return q
}

func (q *QuotaResourceAssert) HasCustomActionsCount(expected int) *QuotaResourceAssert {
	q.AddAssertion(assert.ValueSet("custom_action.#", strconv.Itoa(expected)))
	return q
}

func (q *QuotaResourceAssert) HasCustomAction(procedure sdk.SchemaObjectIdentifierWithArguments, spendStrategy sdk.QuotaSpendStrategy, threshold int) *QuotaResourceAssert {
	q.AddAssertion(assert.SetElemNested("custom_action", map[string]string{
		"procedure":      procedure.FullyQualifiedName(),
		"spend_strategy": string(spendStrategy),
		"threshold":      strconv.Itoa(threshold),
	}))
	return q
}

func (q *QuotaResourceAssert) HasCycleStartAction(procedure sdk.SchemaObjectIdentifierWithArguments) *QuotaResourceAssert {
	q.AddAssertion(assert.ValueSet("cycle_start_action.#", "1"))
	q.AddAssertion(assert.ValueSet("cycle_start_action.0.procedure", procedure.FullyQualifiedName()))
	return q
}
