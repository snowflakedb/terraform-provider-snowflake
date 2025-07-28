package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (d *ComputePoolsModel) WithLimit(rows int) *ComputePoolsModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (d *ComputePoolsModel) WithRowsAndFrom(rows int, from string) *ComputePoolsModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
