package sdk

import (
	"fmt"
	"strings"
)

type DataMetricFunctionRefEntityDomainOption string

const (
	DataMetricFunctionRefEntityDomainView DataMetricFunctionRefEntityDomainOption = "VIEW"
)

type DataMetricScheduleStatusOption string

const (
	DataMetricScheduleStatusStarted                                                   DataMetricScheduleStatusOption = "STARTED"
	DataMetricScheduleStatusStartedAndPendingScheduleUpdate                           DataMetricScheduleStatusOption = "STARTED_AND_PENDING_SCHEDULE_UPDATE"
	DataMetricScheduleStatusSuspended                                                 DataMetricScheduleStatusOption = "SUSPENDED"
	DataMetricScheduleStatusSuspendedTableDoesNotExistOrNotAuthorized                 DataMetricScheduleStatusOption = "SUSPENDED_TABLE_DOES_NOT_EXIST_OR_NOT_AUTHORIZED"
	DataMetricScheduleStatusSuspendedDataMetricFunctionDoesNotExistOrNotAuthorized    DataMetricScheduleStatusOption = "SUSPENDED_DATA_METRIC_FUNCTION_DOES_NOT_EXIST_OR_NOT_AUTHORIZED"
	DataMetricScheduleStatusSuspendedTableColumnDoesNotExistOrNotAuthorized           DataMetricScheduleStatusOption = "SUSPENDED_TABLE_COLUMN_DOES_NOT_EXIST_OR_NOT_AUTHORIZED"
	DataMetricScheduleStatusSuspendedInsufficientPrivilegeToExecuteDataMetricFunction DataMetricScheduleStatusOption = "SUSPENDED_INSUFFICIENT_PRIVILEGE_TO_EXECUTE_DATA_METRIC_FUNCTION"
	DataMetricScheduleStatusSuspendedActiveEventTableDoesNotExistOrNotAuthorized      DataMetricScheduleStatusOption = "SUSPENDED_ACTIVE_EVENT_TABLE_DOES_NOT_EXIST_OR_NOT_AUTHORIZED"
	DataMetricScheduleStatusSuspendedByUserAction                                     DataMetricScheduleStatusOption = "SUSPENDED_BY_USER_ACTION"
)

var AllAllowedDataMetricScheduleStatusOptions = []DataMetricScheduleStatusOption{
	DataMetricScheduleStatusStarted,
	DataMetricScheduleStatusSuspended,
}

var AllDataMetricScheduleStatusStartedOptions = []DataMetricScheduleStatusOption{
	DataMetricScheduleStatusStarted,
	DataMetricScheduleStatusStartedAndPendingScheduleUpdate,
}

var AllDataMetricScheduleStatusSuspendedOptions = []DataMetricScheduleStatusOption{
	DataMetricScheduleStatusSuspended,
	DataMetricScheduleStatusSuspendedTableDoesNotExistOrNotAuthorized,
	DataMetricScheduleStatusSuspendedDataMetricFunctionDoesNotExistOrNotAuthorized,
	DataMetricScheduleStatusSuspendedTableColumnDoesNotExistOrNotAuthorized,
	DataMetricScheduleStatusSuspendedInsufficientPrivilegeToExecuteDataMetricFunction,
	DataMetricScheduleStatusSuspendedActiveEventTableDoesNotExistOrNotAuthorized,
	DataMetricScheduleStatusSuspendedByUserAction,
}

func ToAllowedDataMetricScheduleStatusOption(s string) (DataMetricScheduleStatusOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(DataMetricScheduleStatusStarted):
		return DataMetricScheduleStatusStarted, nil
	case string(DataMetricScheduleStatusSuspended):
		return DataMetricScheduleStatusSuspended, nil
	default:
		return "", fmt.Errorf("invalid DataMetricScheduleStatusOption: %s", s)
	}
}

func ToDataMetricScheduleStatusOption(s string) (DataMetricScheduleStatusOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(DataMetricScheduleStatusStarted):
		return DataMetricScheduleStatusStarted, nil
	case string(DataMetricScheduleStatusStartedAndPendingScheduleUpdate):
		return DataMetricScheduleStatusStartedAndPendingScheduleUpdate, nil
	case string(DataMetricScheduleStatusSuspended):
		return DataMetricScheduleStatusSuspended, nil
	case string(DataMetricScheduleStatusSuspendedTableDoesNotExistOrNotAuthorized):
		return DataMetricScheduleStatusSuspendedTableDoesNotExistOrNotAuthorized, nil
	case string(DataMetricScheduleStatusSuspendedDataMetricFunctionDoesNotExistOrNotAuthorized):
		return DataMetricScheduleStatusSuspendedDataMetricFunctionDoesNotExistOrNotAuthorized, nil
	case string(DataMetricScheduleStatusSuspendedTableColumnDoesNotExistOrNotAuthorized):
		return DataMetricScheduleStatusSuspendedTableColumnDoesNotExistOrNotAuthorized, nil
	case string(DataMetricScheduleStatusSuspendedInsufficientPrivilegeToExecuteDataMetricFunction):
		return DataMetricScheduleStatusSuspendedInsufficientPrivilegeToExecuteDataMetricFunction, nil
	case string(DataMetricScheduleStatusSuspendedActiveEventTableDoesNotExistOrNotAuthorized):
		return DataMetricScheduleStatusSuspendedActiveEventTableDoesNotExistOrNotAuthorized, nil
	case string(DataMetricScheduleStatusSuspendedByUserAction):
		return DataMetricScheduleStatusSuspendedByUserAction, nil
	default:
		return "", fmt.Errorf("invalid DataMetricScheduleStatusOption: %s", s)
	}
}
