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

func (a *AlertModel) WithAlertScheduleCron(expression, timeZone string) *AlertModel {
	return a.WithAlertScheduleValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"cron": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{"expression": tfconfig.StringVariable(expression), "time_zone": tfconfig.StringVariable(timeZone)})),
	})))
}

// WithAlertScheduleIntervalAndCron builds an invalid configuration for conflict validation tests.
func (a *AlertModel) WithAlertScheduleIntervalAndCron(interval int, expression, timeZone string) *AlertModel {
	return a.WithAlertScheduleValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"interval": tfconfig.IntegerVariable(interval),
		"cron":     tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{"expression": tfconfig.StringVariable(expression), "time_zone": tfconfig.StringVariable(timeZone)})),
	})))
}
