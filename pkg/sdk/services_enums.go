package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ServiceStatus string

const (
	ServiceStatusPending       ServiceStatus = "PENDING"
	ServiceStatusRunning       ServiceStatus = "RUNNING"
	ServiceStatusFailed        ServiceStatus = "FAILED"
	ServiceStatusDone          ServiceStatus = "DONE"
	ServiceStatusSuspending    ServiceStatus = "SUSPENDING"
	ServiceStatusSuspended     ServiceStatus = "SUSPENDED"
	ServiceStatusDeleting      ServiceStatus = "DELETING"
	ServiceStatusDeleted       ServiceStatus = "DELETED"
	ServiceStatusInternalError ServiceStatus = "INTERNAL_ERROR"
)

var allServiceStatuses = []ServiceStatus{
	ServiceStatusPending,
	ServiceStatusRunning,
	ServiceStatusFailed,
	ServiceStatusDone,
	ServiceStatusSuspending,
	ServiceStatusSuspended,
	ServiceStatusDeleting,
	ServiceStatusDeleted,
	ServiceStatusInternalError,
}

func ToServiceStatus(s string) (ServiceStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allServiceStatuses, ServiceStatus(s)) {
		return "", fmt.Errorf("invalid service status: %s", s)
	}
	return ServiceStatus(s), nil
}
