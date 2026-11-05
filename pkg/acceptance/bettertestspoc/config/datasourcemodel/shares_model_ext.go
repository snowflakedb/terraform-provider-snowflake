package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (s *SharesModel) WithLimit(rows int) *SharesModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (s *SharesModel) WithRowsAndFrom(rows int, from string) *SharesModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
