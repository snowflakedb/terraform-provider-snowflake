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

// ExpectedColumn describes the fields of a single `column` block to assert on. MaskingPolicy and
// ProjectionPolicy are only checked when non-nil. MaskingPolicyUsing is only checked when MaskingPolicy
// is non-nil and MaskingPolicyUsing is non-empty.
type ExpectedColumn struct {
	Name               string
	Type               string
	NotNull            bool
	Comment            string
	DefaultExpression  string
	MaskingPolicy      *sdk.SchemaObjectIdentifier
	MaskingPolicyUsing []string
	ProjectionPolicy   *sdk.SchemaObjectIdentifier
}

func (i *IcebergTableResourceAssert) HasColumns(columns ...ExpectedColumn) *IcebergTableResourceAssert {
	i.CollectionLength("column", len(columns))
	for index, column := range columns {
		prefix := fmt.Sprintf("column.%d.", index)
		i.ValueSet(prefix+"name", column.Name)
		i.ValueSet(prefix+"type", column.Type)
		i.BoolValueSet(prefix+"not_null", column.NotNull)
		i.ValueSet(prefix+"comment", column.Comment)

		if column.DefaultExpression != "" {
			i.CollectionLength(prefix+"default", 1)
			i.ValueSet(prefix+"default.0.expression", column.DefaultExpression)
		} else {
			i.CollectionLength(prefix+"default", 0)
		}

		if column.MaskingPolicy != nil {
			i.CollectionLength(prefix+"masking_policy", 1)
			i.ValueSet(prefix+"masking_policy.0.policy_name", column.MaskingPolicy.FullyQualifiedName())
			i.CollectionLength(prefix+"masking_policy.0.using", len(column.MaskingPolicyUsing))
			for using, columnName := range column.MaskingPolicyUsing {
				i.ValueSet(fmt.Sprintf("%smasking_policy.0.using.%d", prefix, using), columnName)
			}
		} else {
			i.CollectionLength(prefix+"masking_policy", 0)
		}

		if column.ProjectionPolicy != nil {
			i.CollectionLength(prefix+"projection_policy", 1)
			i.ValueSet(prefix+"projection_policy.0.policy_name", column.ProjectionPolicy.FullyQualifiedName())
		} else {
			i.CollectionLength(prefix+"projection_policy", 0)
		}
	}
	return i
}

// boolPairExpectation renders the expected tri-state boolean string for a constraint enforcement
// field, given the pair of mutually exclusive SQL keyword flags (e.g. Enforced/NotEnforced). Mirrors
// model.boolPairVariable so assertions can reuse the same SDK request structs as the config builders.
func boolPairExpectation(positive, negative *bool) string {
	switch {
	case positive != nil && *positive:
		return "true"
	case negative != nil && *negative:
		return "false"
	default:
		return "default"
	}
}

func (i *IcebergTableResourceAssert) hasConstraintEnforcementFields(prefix string, enforced, notEnforced, deferrable, notDeferrable, initiallyDeferred, initiallyImmediate, enable, disable, validate, novalidate, rely, norely *bool) *IcebergTableResourceAssert {
	i.ValueSet(prefix+"enforced", boolPairExpectation(enforced, notEnforced))
	i.ValueSet(prefix+"deferrable", boolPairExpectation(deferrable, notDeferrable))
	i.ValueSet(prefix+"initially_deferred", boolPairExpectation(initiallyDeferred, initiallyImmediate))
	i.ValueSet(prefix+"enable", boolPairExpectation(enable, disable))
	i.ValueSet(prefix+"validate", boolPairExpectation(validate, novalidate))
	i.ValueSet(prefix+"rely", boolPairExpectation(rely, norely))
	return i
}

// hasUniquePKConstraintFields asserts the fields shared by primary_key_constraint and
// unique_constraint entries, reusing the SDK's TableOutOfLineUniquePKRequest for the expected
// field set.
func (i *IcebergTableResourceAssert) hasUniquePKConstraintFields(prefix string, c sdk.TableOutOfLineUniquePKRequest) *IcebergTableResourceAssert {
	i.OptionalStringValueSet(prefix+"name", c.Name)
	i.OptionalStringValueSet(prefix+"comment", c.Comment)
	i.CollectionLength(prefix+"column", len(c.Columns))
	for colIndex, column := range c.Columns {
		i.ValueSet(fmt.Sprintf("%scolumn.%d", prefix, colIndex), column.Value)
	}
	i.hasConstraintEnforcementFields(prefix, c.Enforced, c.NotEnforced, c.Deferrable, c.NotDeferrable, c.InitiallyDeferred, c.InitiallyImmediate, c.Enable, c.Disable, c.Validate, c.Novalidate, c.Rely, c.Norely)
	return i
}

// HasPrimaryKeyConstraints asserts the `primary_key_constraint` attribute, reusing the SDK's
// TableOutOfLineUniquePKRequest for the expected field set.
func (i *IcebergTableResourceAssert) HasPrimaryKeyConstraints(constraints ...sdk.TableOutOfLineUniquePKRequest) *IcebergTableResourceAssert {
	i.CollectionLength("primary_key_constraint", len(constraints))
	for index, c := range constraints {
		i.hasUniquePKConstraintFields(fmt.Sprintf("primary_key_constraint.%d.", index), c)
	}
	return i
}

// HasUniqueConstraints asserts the `unique_constraint` attribute, reusing the SDK's
// TableOutOfLineUniquePKRequest for the expected field set.
func (i *IcebergTableResourceAssert) HasUniqueConstraints(constraints ...sdk.TableOutOfLineUniquePKRequest) *IcebergTableResourceAssert {
	i.CollectionLength("unique_constraint", len(constraints))
	for index, c := range constraints {
		i.hasUniquePKConstraintFields(fmt.Sprintf("unique_constraint.%d.", index), c)
	}
	return i
}

// HasForeignKeyConstraints asserts the `foreign_key_constraint` attribute, reusing the SDK's
// TableOutOfLineFKRequest for the expected field set.
func (i *IcebergTableResourceAssert) HasForeignKeyConstraints(constraints ...sdk.TableOutOfLineFKRequest) *IcebergTableResourceAssert {
	i.CollectionLength("foreign_key_constraint", len(constraints))
	for index, c := range constraints {
		prefix := fmt.Sprintf("foreign_key_constraint.%d.", index)
		i.OptionalStringValueSet(prefix+"name", c.Name)
		i.ValueSet(prefix+"table_name", c.References.FullyQualifiedName())
		i.OptionalStringValueSet(prefix+"comment", c.Comment)
		i.CollectionLength(prefix+"column", len(c.Columns))
		for colIndex, column := range c.Columns {
			i.ValueSet(fmt.Sprintf("%scolumn.%d", prefix, colIndex), column.Value)
		}
		i.CollectionLength(prefix+"ref_column", len(c.RefColumns))
		for colIndex, column := range c.RefColumns {
			i.ValueSet(fmt.Sprintf("%sref_column.%d", prefix, colIndex), column.Value)
		}
		if c.Match != nil {
			i.ValueSet(prefix+"match", string(*c.Match))
		}
		if c.On != nil && c.On.OnUpdate != nil {
			i.ValueSet(prefix+"on_update", string(*c.On.OnUpdate))
		}
		if c.On != nil && c.On.OnDelete != nil {
			i.ValueSet(prefix+"on_delete", string(*c.On.OnDelete))
		}
		i.hasConstraintEnforcementFields(prefix, c.Enforced, c.NotEnforced, c.Deferrable, c.NotDeferrable, c.InitiallyDeferred, c.InitiallyImmediate, c.Enable, c.Disable, c.Validate, c.Novalidate, c.Rely, c.Norely)
	}
	return i
}

// HasCheckConstraints asserts the `check_constraint` attribute, reusing the SDK's
// TableOutOfLineCHRequest for the expected field set.
func (i *IcebergTableResourceAssert) HasCheckConstraints(constraints ...sdk.TableOutOfLineCHRequest) *IcebergTableResourceAssert {
	i.CollectionLength("check_constraint", len(constraints))
	for index, c := range constraints {
		prefix := fmt.Sprintf("check_constraint.%d.", index)
		i.OptionalStringValueSet(prefix+"name", c.Name)
		i.ValueSet(prefix+"expression", c.Expression)
		i.ValueSet(prefix+"validate", boolPairExpectation(c.EnableValidate, c.EnableNovalidate))
	}
	return i
}
