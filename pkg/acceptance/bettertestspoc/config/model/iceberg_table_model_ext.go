package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// WithColumn satisfies the generated constructor's call for the complex list `column` attribute.
func (i *IcebergTableModel) WithColumn(column []sdk.TableColumnSignature) *IcebergTableModel {
	columns := make([]tfconfig.Variable, len(column))
	for idx, v := range column {
		columns[idx] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	i.Column = tfconfig.ListVariable(columns...)
	return i
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
