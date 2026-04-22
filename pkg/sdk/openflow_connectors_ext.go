package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type OpenflowConnectorStatus string

const (
	OpenflowConnectorStatusCreating     OpenflowConnectorStatus = "CREATING"
	OpenflowConnectorStatusCreateFailed OpenflowConnectorStatus = "CREATE_FAILED"
	OpenflowConnectorStatusStarting     OpenflowConnectorStatus = "STARTING"
	OpenflowConnectorStatusStartFailed  OpenflowConnectorStatus = "START_FAILED"
	OpenflowConnectorStatusRunning      OpenflowConnectorStatus = "RUNNING"
	OpenflowConnectorStatusStopping     OpenflowConnectorStatus = "STOPPING"
	OpenflowConnectorStatusStopFailed   OpenflowConnectorStatus = "STOP_FAILED"
	OpenflowConnectorStatusStopped      OpenflowConnectorStatus = "STOPPED"
	OpenflowConnectorStatusUpdating     OpenflowConnectorStatus = "UPDATING"
	OpenflowConnectorStatusUpdateFailed OpenflowConnectorStatus = "UPDATE_FAILED"
	OpenflowConnectorStatusDeleting     OpenflowConnectorStatus = "DELETING"
	OpenflowConnectorStatusDeleteFailed OpenflowConnectorStatus = "DELETE_FAILED"
	OpenflowConnectorStatusDeleted      OpenflowConnectorStatus = "DELETED"
)

var allOpenflowConnectorStatuses = []OpenflowConnectorStatus{
	OpenflowConnectorStatusCreating,
	OpenflowConnectorStatusCreateFailed,
	OpenflowConnectorStatusStarting,
	OpenflowConnectorStatusStartFailed,
	OpenflowConnectorStatusRunning,
	OpenflowConnectorStatusStopping,
	OpenflowConnectorStatusStopFailed,
	OpenflowConnectorStatusStopped,
	OpenflowConnectorStatusUpdating,
	OpenflowConnectorStatusUpdateFailed,
	OpenflowConnectorStatusDeleting,
	OpenflowConnectorStatusDeleteFailed,
	OpenflowConnectorStatusDeleted,
}

func ToOpenflowConnectorStatus(s string) (OpenflowConnectorStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allOpenflowConnectorStatuses, OpenflowConnectorStatus(s)) {
		return "", fmt.Errorf("invalid openflow connector status: %s", s)
	}
	return OpenflowConnectorStatus(s), nil
}

func (r *CreateOpenflowConnectorRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
