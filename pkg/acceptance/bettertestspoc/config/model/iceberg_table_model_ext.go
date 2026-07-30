package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

// WithColumn satisfies the generated constructor's call for the complex list `column` attribute.
func (i *IcebergTableModel) WithColumn(column []sdk.TableColumnSignature) *IcebergTableModel {
	columns := make([]tfconfig.Variable, len(column))
	for idx, v := range column {
		columns[idx] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	i.Column = tfconfig.ListVariable(columns...)
	return i
}

// IcebergTableColumnRequest describes a single entry of the `column` attribute, covering the fields
// that go beyond name + type (not_null, comment, default expression, masking_policy, projection_policy).
type IcebergTableColumnRequest struct {
	Name               string
	Type               datatypes.DataType
	NotNull            *string
	Comment            string
	DefaultExpression  string
	MaskingPolicy      *sdk.SchemaObjectIdentifier
	MaskingPolicyUsing []string
	ProjectionPolicy   *sdk.SchemaObjectIdentifier
}

// WithColumns is like WithColumn, but supports setting the full set of per-column fields
// (not_null, comment, default expression, masking_policy, projection_policy).
func (i *IcebergTableModel) WithColumns(columns ...IcebergTableColumnRequest) *IcebergTableModel {
	vars := make([]tfconfig.Variable, len(columns))
	for idx, c := range columns {
		m := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(c.Name),
			"type": tfconfig.StringVariable(c.Type.ToSql()),
		}
		if c.NotNull != nil {
			m["not_null"] = tfconfig.StringVariable(*c.NotNull)
		}
		if c.Comment != "" {
			m["comment"] = tfconfig.StringVariable(c.Comment)
		}
		if c.DefaultExpression != "" {
			m["default"] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"expression": tfconfig.StringVariable(c.DefaultExpression),
			})
		}
		if c.MaskingPolicy != nil {
			maskingPolicy := map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(c.MaskingPolicy.FullyQualifiedName()),
			}
			if len(c.MaskingPolicyUsing) > 0 {
				maskingPolicy["using"] = tfconfig.ListVariable(collections.Map(c.MaskingPolicyUsing, func(s string) tfconfig.Variable {
					return tfconfig.StringVariable(s)
				})...)
			}
			m["masking_policy"] = tfconfig.ObjectVariable(maskingPolicy)
		}
		if c.ProjectionPolicy != nil {
			m["projection_policy"] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(c.ProjectionPolicy.FullyQualifiedName()),
			})
		}
		vars[idx] = tfconfig.ObjectVariable(m)
	}
	i.Column = tfconfig.ListVariable(vars...)
	return i
}

// uniquePKConstraintVariable builds the fields shared by primary_key_constraint and
// unique_constraint entries, reusing the SDK's TableOutOfLineUniquePKRequest for the field set.
func uniquePKConstraintVariable(c sdk.TableOutOfLineUniquePKRequest) tfconfig.Variable {
	m := map[string]tfconfig.Variable{
		"column": columnsVariable(c.Columns),
	}
	setStringIfNotNil(m, "name", c.Name)
	setStringIfNotNil(m, "comment", c.Comment)
	setBoolPairIfNotNil(m, "enforced", c.Enforced, c.NotEnforced)
	setBoolPairIfNotNil(m, "deferrable", c.Deferrable, c.NotDeferrable)
	setBoolPairIfNotNil(m, "initially_deferred", c.InitiallyDeferred, c.InitiallyImmediate)
	setBoolPairIfNotNil(m, "enable", c.Enable, c.Disable)
	setBoolPairIfNotNil(m, "validate", c.Validate, c.Novalidate)
	setBoolPairIfNotNil(m, "rely", c.Rely, c.Norely)
	return tfconfig.ObjectVariable(m)
}

// WithPrimaryKeyConstraints sets the `primary_key_constraint` attribute (table-level PRIMARY KEY
// constraints), reusing the SDK's TableOutOfLineUniquePKRequest for the field set.
func (i *IcebergTableModel) WithPrimaryKeyConstraints(constraints ...sdk.TableOutOfLineUniquePKRequest) *IcebergTableModel {
	vars := collections.Map(constraints, uniquePKConstraintVariable)
	return i.WithPrimaryKeyConstraintValue(tfconfig.ListVariable(vars...))
}

// WithUniqueConstraints sets the `unique_constraint` attribute (table-level UNIQUE constraints),
// reusing the SDK's TableOutOfLineUniquePKRequest for the field set.
func (i *IcebergTableModel) WithUniqueConstraints(constraints ...sdk.TableOutOfLineUniquePKRequest) *IcebergTableModel {
	vars := collections.Map(constraints, uniquePKConstraintVariable)
	return i.WithUniqueConstraintValue(tfconfig.ListVariable(vars...))
}

// WithForeignKeyConstraints sets the `foreign_key_constraint` attribute, reusing the SDK's
// TableOutOfLineFKRequest for the field set.
func (i *IcebergTableModel) WithForeignKeyConstraints(constraints ...sdk.TableOutOfLineFKRequest) *IcebergTableModel {
	vars := make([]tfconfig.Variable, len(constraints))
	for idx, c := range constraints {
		m := map[string]tfconfig.Variable{
			"column":     columnsVariable(c.Columns),
			"table_name": tfconfig.StringVariable(c.References.FullyQualifiedName()),
		}
		setStringIfNotNil(m, "name", c.Name)
		if len(c.RefColumns) > 0 {
			m["ref_column"] = columnsVariable(c.RefColumns)
		}
		if c.Match != nil {
			m["match"] = tfconfig.StringVariable(string(*c.Match))
		}
		if c.On != nil {
			if c.On.OnUpdate != nil {
				m["on_update"] = tfconfig.StringVariable(string(*c.On.OnUpdate))
			}
			if c.On.OnDelete != nil {
				m["on_delete"] = tfconfig.StringVariable(string(*c.On.OnDelete))
			}
		}
		setStringIfNotNil(m, "comment", c.Comment)
		setBoolPairIfNotNil(m, "enforced", c.Enforced, c.NotEnforced)
		setBoolPairIfNotNil(m, "deferrable", c.Deferrable, c.NotDeferrable)
		setBoolPairIfNotNil(m, "initially_deferred", c.InitiallyDeferred, c.InitiallyImmediate)
		setBoolPairIfNotNil(m, "enable", c.Enable, c.Disable)
		setBoolPairIfNotNil(m, "validate", c.Validate, c.Novalidate)
		setBoolPairIfNotNil(m, "rely", c.Rely, c.Norely)
		vars[idx] = tfconfig.ObjectVariable(m)
	}
	return i.WithForeignKeyConstraintValue(tfconfig.ListVariable(vars...))
}

// WithCheckConstraints sets the `check_constraint` attribute, reusing the SDK's TableOutOfLineCHRequest
// for the field set.
func (i *IcebergTableModel) WithCheckConstraints(constraints ...sdk.TableOutOfLineCHRequest) *IcebergTableModel {
	vars := make([]tfconfig.Variable, len(constraints))
	for idx, c := range constraints {
		m := map[string]tfconfig.Variable{
			"expression": tfconfig.StringVariable(c.Expression),
		}
		setStringIfNotNil(m, "name", c.Name)
		setBoolPairIfNotNil(m, "validate", c.EnableValidate, c.EnableNovalidate)
		vars[idx] = tfconfig.ObjectVariable(m)
	}
	return i.WithCheckConstraintValue(tfconfig.ListVariable(vars...))
}

func (i *IcebergTableModel) WithRowAccessPolicy(rap sdk.SchemaObjectIdentifier, on string) *IcebergTableModel {
	return i.WithRowAccessPolicyValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(rap.FullyQualifiedName()),
				"on":          tfconfig.ListVariable(tfconfig.StringVariable(on)),
			},
		),
	)
}

func (i *IcebergTableModel) WithAggregationPolicy(ap sdk.SchemaObjectIdentifier, entityKey ...string) *IcebergTableModel {
	m := map[string]tfconfig.Variable{
		"policy_name": tfconfig.StringVariable(ap.FullyQualifiedName()),
	}
	if len(entityKey) > 0 {
		m["entity_key"] = tfconfig.ListVariable(collections.Map(entityKey, func(s string) tfconfig.Variable {
			return tfconfig.StringVariable(s)
		})...)
	}
	return i.WithAggregationPolicyValue(
		tfconfig.ObjectVariable(
			m,
		),
	)
}

// WithClusterBy satisfies the generated constructor's call for the complex list `cluster_by` attribute.
func (i *IcebergTableModel) WithClusterBy(clusterBy ...string) *IcebergTableModel {
	return i.WithClusterByValue(tfconfig.ListVariable(collections.Map(clusterBy, func(s string) tfconfig.Variable {
		return tfconfig.StringVariable(s)
	})...))
}

// WithPartitionBy satisfies the generated constructor's call for the complex list `partition_by` attribute.
// Build each entry with IcebergTablePartitionByIdentity/Bucket/Truncate/Year/Month/Day/Hour.
func (i *IcebergTableModel) WithPartitionBy(entries ...tfconfig.Variable) *IcebergTableModel {
	return i.WithPartitionByValue(tfconfig.ListVariable(entries...))
}

func IcebergTablePartitionByIdentity(column string) tfconfig.Variable {
	return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"identity": tfconfig.StringVariable(column),
	})
}

func IcebergTablePartitionByBucket(numBuckets int, column string) tfconfig.Variable {
	return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"bucket": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"num_buckets": tfconfig.IntegerVariable(numBuckets),
			"column":      tfconfig.StringVariable(column),
		})),
	})
}

func IcebergTablePartitionByTruncate(width int, column string) tfconfig.Variable {
	return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"truncate": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"width":  tfconfig.IntegerVariable(width),
			"column": tfconfig.StringVariable(column),
		})),
	})
}

func IcebergTablePartitionByYear(column string) tfconfig.Variable {
	return icebergTablePartitionByTimeVariable("year", column)
}

func IcebergTablePartitionByMonth(column string) tfconfig.Variable {
	return icebergTablePartitionByTimeVariable("month", column)
}

func IcebergTablePartitionByDay(column string) tfconfig.Variable {
	return icebergTablePartitionByTimeVariable("day", column)
}

func IcebergTablePartitionByHour(column string) tfconfig.Variable {
	return icebergTablePartitionByTimeVariable("hour", column)
}

func icebergTablePartitionByTimeVariable(kind string, column string) tfconfig.Variable {
	return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		kind: tfconfig.StringVariable(column),
	})
}
