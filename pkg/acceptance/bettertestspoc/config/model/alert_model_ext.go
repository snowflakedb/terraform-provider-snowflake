package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AlertModel) WithAlertScheduleInterval(interval int) *AlertModel {
	return a.WithAlertScheduleValue(
		tfconfig.ListVariable(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"interval": tfconfig.IntegerVariable(interval),
			}),
		),
	)
}
