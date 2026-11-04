package sdk

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	alertCronSchedulePrefix     = "USING CRON "
	alertIntervalScheduleSuffix = " MINUTE"
)

// AlertSchedule is a parsed alert SCHEDULE value.
// A minute interval sets Minutes. A cron schedule sets CronExpression and TimeZone.
type AlertSchedule struct {
	Minutes        int
	CronExpression string
	TimeZone       string
}

// FormatAlertCronSchedule builds the SCHEDULE value for a cron expression and time zone.
func FormatAlertCronSchedule(expression, timeZone string) string {
	return fmt.Sprintf("%s%s %s", alertCronSchedulePrefix, expression, timeZone)
}

// FormatAlertIntervalSchedule builds the SCHEDULE value for a minute interval.
func FormatAlertIntervalSchedule(minutes int) string {
	return fmt.Sprintf("%d%s", minutes, alertIntervalScheduleSuffix)
}

// ParseAlertSchedule parses a SCHEDULE value from SHOW ALERTS.
// Cron values use the USING CRON prefix. The last token is the time zone.
// Interval values use the MINUTE suffix.
func ParseAlertSchedule(schedule string) (AlertSchedule, error) {
	if rest, ok := strings.CutPrefix(schedule, alertCronSchedulePrefix); ok {
		// A cron expression is five fields, so the time zone is the last token.
		idx := strings.LastIndex(rest, " ")
		if idx <= 0 || idx == len(rest)-1 {
			return AlertSchedule{}, fmt.Errorf("invalid alert cron schedule %q", schedule)
		}
		return AlertSchedule{CronExpression: rest[:idx], TimeZone: rest[idx+1:]}, nil
	}
	if before, ok := strings.CutSuffix(schedule, alertIntervalScheduleSuffix); ok {
		minutes, err := strconv.Atoi(before)
		if err != nil {
			return AlertSchedule{}, fmt.Errorf("invalid alert interval schedule %q: %w", schedule, err)
		}
		return AlertSchedule{Minutes: minutes}, nil
	}
	return AlertSchedule{}, fmt.Errorf("invalid alert schedule %q", schedule)
}

// NewAlertConditionFromString builds the IF (EXISTS (...)) block from a SQL condition string.
func NewAlertConditionFromString(condition string) AlertConditionRequest {
	return *NewAlertConditionRequest().WithCondition([]string{condition})
}

func (r *CreateAlertRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
