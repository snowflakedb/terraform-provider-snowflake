package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (d *DatabaseRolesModel) WithRows(rows int) *DatabaseRolesModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (d *DatabaseRolesModel) WithRowsAndFrom(rows int, from string) *DatabaseRolesModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
