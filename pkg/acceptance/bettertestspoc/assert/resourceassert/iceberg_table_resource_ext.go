package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *IcebergTableResourceAssert) HasRowAccessPolicy(rowAccessPolicy sdk.SchemaObjectIdentifier, on ...string) *IcebergTableResourceAssert {
	i.ValueSet("row_access_policy.0.policy_name", rowAccessPolicy.FullyQualifiedName())
	i.CollectionLength("row_access_policy.0.on", len(on))
	for _, column := range on {
		i.SetContainsElem("row_access_policy.0.on", column)
	}
	return i
}

func (i *IcebergTableResourceAssert) HasAggregationPolicy(aggregationPolicy sdk.SchemaObjectIdentifier, entityKey ...string) *IcebergTableResourceAssert {
	i.ValueSet("aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())
	i.CollectionLength("aggregation_policy.0.entity_key", len(entityKey))
	for _, key := range entityKey {
		i.SetContainsElem("aggregation_policy.0.entity_key", key)
	}
	return i
}

func (i *IcebergTableResourceAssert) HasNoRowAccessPolicy() *IcebergTableResourceAssert {
	i.ValueNotSet("row_access_policy.#")
	return i
}

func (i *IcebergTableResourceAssert) HasNoAggregationPolicy() *IcebergTableResourceAssert {
	i.ValueNotSet("aggregation_policy.#")
	return i
}
