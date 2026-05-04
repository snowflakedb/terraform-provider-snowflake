package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type OpenflowRuntimeNodeType string

const (
	OpenflowRuntimeNodeTypeSmall  OpenflowRuntimeNodeType = "SMALL"
	OpenflowRuntimeNodeTypeMedium OpenflowRuntimeNodeType = "MEDIUM"
	OpenflowRuntimeNodeTypeLarge  OpenflowRuntimeNodeType = "LARGE"
)

var AllOpenflowRuntimeNodeTypes = []OpenflowRuntimeNodeType{
	OpenflowRuntimeNodeTypeSmall,
	OpenflowRuntimeNodeTypeMedium,
	OpenflowRuntimeNodeTypeLarge,
}

func ToOpenflowRuntimeNodeType(s string) (OpenflowRuntimeNodeType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllOpenflowRuntimeNodeTypes, OpenflowRuntimeNodeType(s)) {
		return "", fmt.Errorf("invalid openflow runtime node type: %s", s)
	}
	return OpenflowRuntimeNodeType(s), nil
}

type OpenflowRuntimeStatus string

const (
	OpenflowRuntimeStatusCreating                   OpenflowRuntimeStatus = "CREATING"
	OpenflowRuntimeStatusCreateFailed               OpenflowRuntimeStatus = "CREATE_FAILED"
	OpenflowRuntimeStatusUpdating                   OpenflowRuntimeStatus = "UPDATING"
	OpenflowRuntimeStatusUpdateFailed               OpenflowRuntimeStatus = "UPDATE_FAILED"
	OpenflowRuntimeStatusSuspending                 OpenflowRuntimeStatus = "SUSPENDING"
	OpenflowRuntimeStatusSuspended                  OpenflowRuntimeStatus = "SUSPENDED"
	OpenflowRuntimeStatusSuspendFailed              OpenflowRuntimeStatus = "SUSPEND_FAILED"
	OpenflowRuntimeStatusActivating                 OpenflowRuntimeStatus = "ACTIVATING"
	OpenflowRuntimeStatusActive                     OpenflowRuntimeStatus = "ACTIVE"
	OpenflowRuntimeStatusActivateFailed             OpenflowRuntimeStatus = "ACTIVATE_FAILED"
	OpenflowRuntimeStatusDeleting                   OpenflowRuntimeStatus = "DELETING"
	OpenflowRuntimeStatusDeleted                    OpenflowRuntimeStatus = "DELETED"
	OpenflowRuntimeStatusDeleteFailed               OpenflowRuntimeStatus = "DELETE_FAILED"
	OpenflowRuntimeStatusCancelRequested            OpenflowRuntimeStatus = "CANCEL_REQUESTED"
	OpenflowRuntimeStatusRestarting                 OpenflowRuntimeStatus = "RESTARTING"
	OpenflowRuntimeStatusRestartFailed              OpenflowRuntimeStatus = "RESTART_FAILED"
	OpenflowRuntimeStatusUpgrading                  OpenflowRuntimeStatus = "UPGRADING"
	OpenflowRuntimeStatusUpgradeFailed              OpenflowRuntimeStatus = "UPGRADE_FAILED"
	OpenflowRuntimeStatusGeneratingDiagnosticBundle OpenflowRuntimeStatus = "GENERATING_DIAGNOSTIC_BUNDLE"
	OpenflowRuntimeStatusCleaningUp                 OpenflowRuntimeStatus = "CLEANING_UP"
	OpenflowRuntimeStatusInactive                   OpenflowRuntimeStatus = "INACTIVE"
)

var allOpenflowRuntimeStatuses = []OpenflowRuntimeStatus{
	OpenflowRuntimeStatusCreating,
	OpenflowRuntimeStatusCreateFailed,
	OpenflowRuntimeStatusUpdating,
	OpenflowRuntimeStatusUpdateFailed,
	OpenflowRuntimeStatusSuspending,
	OpenflowRuntimeStatusSuspended,
	OpenflowRuntimeStatusSuspendFailed,
	OpenflowRuntimeStatusActivating,
	OpenflowRuntimeStatusActive,
	OpenflowRuntimeStatusActivateFailed,
	OpenflowRuntimeStatusDeleting,
	OpenflowRuntimeStatusDeleted,
	OpenflowRuntimeStatusDeleteFailed,
	OpenflowRuntimeStatusCancelRequested,
	OpenflowRuntimeStatusRestarting,
	OpenflowRuntimeStatusRestartFailed,
	OpenflowRuntimeStatusUpgrading,
	OpenflowRuntimeStatusUpgradeFailed,
	OpenflowRuntimeStatusGeneratingDiagnosticBundle,
	OpenflowRuntimeStatusCleaningUp,
	OpenflowRuntimeStatusInactive,
}

func ToOpenflowRuntimeStatus(s string) (OpenflowRuntimeStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allOpenflowRuntimeStatuses, OpenflowRuntimeStatus(s)) {
		return "", fmt.Errorf("invalid openflow runtime status: %s", s)
	}
	return OpenflowRuntimeStatus(s), nil
}

func (r *CreateOpenflowRuntimeRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
