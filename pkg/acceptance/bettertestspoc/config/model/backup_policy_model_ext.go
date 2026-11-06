package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (b *BackupPolicyModel) WithScheduleHours(hours int) *BackupPolicyModel {
	b.Schedule = tfconfig.MapVariable(map[string]tfconfig.Variable{
		"hours": tfconfig.IntegerVariable(hours),
	})
	return b
}

func (b *BackupPolicyModel) WithScheduleMinutes(minutes int) *BackupPolicyModel {
	b.Schedule = tfconfig.MapVariable(map[string]tfconfig.Variable{
		"minutes": tfconfig.IntegerVariable(minutes),
	})
	return b
}

func (b *BackupPolicyModel) WithScheduleCron(cron string) *BackupPolicyModel {
	b.Schedule = tfconfig.MapVariable(map[string]tfconfig.Variable{
		"using_cron": tfconfig.StringVariable(cron),
	})
	return b
}
