package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (d *OpenflowDeploymentsModel) WithLimit(rows int) *OpenflowDeploymentsModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}
