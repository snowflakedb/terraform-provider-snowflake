package resourceassert

import (
	"fmt"
	"strconv"

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

func (i *IcebergTableResourceAssert) HasPartitionByLength(expected int) *IcebergTableResourceAssert {
	i.CollectionLength("partition_by", expected)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByIdentity(index int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.identity", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByBucket(index int, numBuckets int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.bucket.0.num_buckets", index), strconv.Itoa(numBuckets))
	i.ValueSet(fmt.Sprintf("partition_by.%d.bucket.0.column", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByTruncate(index int, width int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.truncate.0.width", index), strconv.Itoa(width))
	i.ValueSet(fmt.Sprintf("partition_by.%d.truncate.0.column", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByYear(index int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.year", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByMonth(index int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.month", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByDay(index int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.day", index), column)
	return i
}

func (i *IcebergTableResourceAssert) HasPartitionByHour(index int, column string) *IcebergTableResourceAssert {
	i.ValueSet(fmt.Sprintf("partition_by.%d.hour", index), column)
	return i
}
