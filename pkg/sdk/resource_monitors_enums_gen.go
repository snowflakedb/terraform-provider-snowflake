package sdk

import (
	"fmt"
	"strings"
)

type ResourceMonitorLevel string

const (
	ResourceMonitorLevelAccount   ResourceMonitorLevel = "ACCOUNT"
	ResourceMonitorLevelWarehouse ResourceMonitorLevel = "WAREHOUSE"
)

func ToResourceMonitorLevel(s string) (ResourceMonitorLevel, error) {
	switch level := ResourceMonitorLevel(strings.ToUpper(s)); level {
	case ResourceMonitorLevelAccount,
		ResourceMonitorLevelWarehouse:
		return level, nil
	default:
		return "", fmt.Errorf("invalid resource monitor level: %s", s)
	}
}

type TriggerAction string

const (
	TriggerActionSuspend          TriggerAction = "SUSPEND"
	TriggerActionSuspendImmediate TriggerAction = "SUSPEND_IMMEDIATE"
	TriggerActionNotify           TriggerAction = "NOTIFY"
)

func ToResourceMonitorTriggerAction(s string) (*TriggerAction, error) {
	switch action := TriggerAction(strings.ToUpper(s)); action {
	case TriggerActionSuspend,
		TriggerActionSuspendImmediate,
		TriggerActionNotify:
		return &action, nil
	default:
		return nil, fmt.Errorf("invalid trigger action type: %s", s)
	}
}

type Frequency string

const (
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyYearly  Frequency = "YEARLY"
	FrequencyNever   Frequency = "NEVER"
)

var AllFrequencyValues = []Frequency{
	FrequencyMonthly,
	FrequencyDaily,
	FrequencyWeekly,
	FrequencyYearly,
	FrequencyNever,
}

func ToResourceMonitorFrequency(s string) (*Frequency, error) {
	switch frequency := Frequency(strings.ToUpper(s)); frequency {
	case FrequencyDaily,
		FrequencyWeekly,
		FrequencyMonthly,
		FrequencyYearly,
		FrequencyNever:
		return &frequency, nil
	default:
		return nil, fmt.Errorf("invalid frequency type: %s", s)
	}
}
