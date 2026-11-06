package resourceassert

import "strconv"

func (b *BackupPolicyResourceAssert) HasScheduleHours(hours int) *BackupPolicyResourceAssert {
	b.ValueSet("schedule.#", "1")
	b.ValueSet("schedule.0.hours", strconv.Itoa(hours))
	b.ValueSet("schedule.0.minutes", "0")
	b.ValueSet("schedule.0.using_cron", "")
	return b
}

func (b *BackupPolicyResourceAssert) HasScheduleMinutes(minutes int) *BackupPolicyResourceAssert {
	b.ValueSet("schedule.#", "1")
	b.ValueSet("schedule.0.minutes", strconv.Itoa(minutes))
	b.ValueSet("schedule.0.hours", "0")
	b.ValueSet("schedule.0.using_cron", "")
	return b
}

func (b *BackupPolicyResourceAssert) HasScheduleCron(cron string) *BackupPolicyResourceAssert {
	b.ValueSet("schedule.#", "1")
	b.ValueSet("schedule.0.using_cron", cron)
	b.ValueSet("schedule.0.hours", "0")
	b.ValueSet("schedule.0.minutes", "0")
	return b
}

func (b *BackupPolicyResourceAssert) HasNoSchedule() *BackupPolicyResourceAssert {
	b.ValueSet("schedule.#", "0")
	return b
}
