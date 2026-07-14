package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

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
