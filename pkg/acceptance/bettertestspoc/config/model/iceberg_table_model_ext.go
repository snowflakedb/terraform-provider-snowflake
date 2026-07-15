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
