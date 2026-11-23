package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// ExpectedTableColumn describes the fields of a single `column` block to assert on.
// An empty DefaultExpression means that the column has no `default` block.
type ExpectedTableColumn struct {
	Name              string
	Type              string
	DefaultExpression string
}

// HasColumns asserts all the columns, in the order of the configuration.
func (t *TableResourceAssert) HasColumns(columns ...ExpectedTableColumn) *TableResourceAssert {
	t.CollectionLength("column", len(columns))
	for index, column := range columns {
		prefix := fmt.Sprintf("column.%d.", index)
		t.ValueSet(prefix+"name", column.Name)
		t.ValueSet(prefix+"type", column.Type)

		if column.DefaultExpression != "" {
			t.CollectionLength(prefix+"default", 1)
			t.ValueSet(prefix+"default.0.expression", column.DefaultExpression)
		} else {
			t.CollectionLength(prefix+"default", 0)
		}
	}
	return t
}

func (t *TableResourceAssert) HasRowAccessPolicy(rowAccessPolicy sdk.SchemaObjectIdentifier, on ...string) *TableResourceAssert {
	t.ValueSet("row_access_policy.0.policy_name", rowAccessPolicy.FullyQualifiedName())
	t.CollectionLength("row_access_policy.0.on", len(on))
	for _, column := range on {
		t.SetContainsElem("row_access_policy.0.on", column)
	}
	return t
}

func (t *TableResourceAssert) HasAggregationPolicy(aggregationPolicy sdk.SchemaObjectIdentifier, entityKey ...string) *TableResourceAssert {
	t.ValueSet("aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())
	t.CollectionLength("aggregation_policy.0.entity_key", len(entityKey))
	for _, key := range entityKey {
		t.SetContainsElem("aggregation_policy.0.entity_key", key)
	}
	return t
}

func (t *TableResourceAssert) HasJoinPolicy(joinPolicy sdk.SchemaObjectIdentifier, allowedJoinKeys ...string) *TableResourceAssert {
	t.ValueSet("join_policy.0.policy_name", joinPolicy.FullyQualifiedName())
	t.CollectionLength("join_policy.0.allowed_join_keys", len(allowedJoinKeys))
	for _, key := range allowedJoinKeys {
		t.SetContainsElem("join_policy.0.allowed_join_keys", key)
	}
	return t
}

func (t *TableResourceAssert) HasNoRowAccessPolicy() *TableResourceAssert {
	t.ValueNotSet("row_access_policy.#")
	return t
}

func (t *TableResourceAssert) HasNoAggregationPolicy() *TableResourceAssert {
	t.ValueNotSet("aggregation_policy.#")
	return t
}

func (t *TableResourceAssert) HasNoJoinPolicy() *TableResourceAssert {
	t.ValueNotSet("join_policy.#")
	return t
}
