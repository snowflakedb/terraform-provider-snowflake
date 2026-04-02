package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (d *NotebooksModel) WithLimit(rows int) *NotebooksModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (d *NotebooksModel) WithRowsAndFrom(rows int, from string) *NotebooksModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
