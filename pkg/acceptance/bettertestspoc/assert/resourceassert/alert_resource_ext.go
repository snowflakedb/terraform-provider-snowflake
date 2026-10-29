package resourceassert

import "strconv"

func (a *AlertResourceAssert) HasAlertScheduleInterval(interval int) *AlertResourceAssert {
	a.ValueSet("alert_schedule.#", "1")
	a.ValueSet("alert_schedule.0.interval", strconv.Itoa(interval))
	a.ValueSet("alert_schedule.0.cron.#", "0")
	return a
}

func (a *AlertResourceAssert) HasAlertScheduleCron(expression, timeZone string) *AlertResourceAssert {
	a.ValueSet("alert_schedule.#", "1")
	a.ValueSet("alert_schedule.0.cron.#", "1")
	a.ValueSet("alert_schedule.0.cron.0.expression", expression)
	a.ValueSet("alert_schedule.0.cron.0.time_zone", timeZone)
	return a
}
