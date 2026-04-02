package sdk

import (
	"fmt"
	"strings"
)

var AllViewDataMetricScheduleMinutes = []int{5, 15, 30, 60, 720, 1440}

type ViewDataMetricScheduleStatusOperationOption string

const (
	ViewDataMetricScheduleStatusOperationResume  ViewDataMetricScheduleStatusOperationOption = "RESUME"
	ViewDataMetricScheduleStatusOperationSuspend ViewDataMetricScheduleStatusOperationOption = "SUSPEND"
)

var AllViewDataMetricScheduleStatusOperationOptions = []ViewDataMetricScheduleStatusOperationOption{
	ViewDataMetricScheduleStatusOperationResume,
	ViewDataMetricScheduleStatusOperationSuspend,
}

func ToViewDataMetricScheduleStatusOperationOption(s string) (ViewDataMetricScheduleStatusOperationOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ViewDataMetricScheduleStatusOperationResume):
		return ViewDataMetricScheduleStatusOperationResume, nil
	case string(ViewDataMetricScheduleStatusOperationSuspend):
		return ViewDataMetricScheduleStatusOperationSuspend, nil
	default:
		return "", fmt.Errorf("invalid ViewDataMetricScheduleStatusOperationOption: %s", s)
	}
}
