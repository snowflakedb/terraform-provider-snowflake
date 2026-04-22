package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type OpenflowDeploymentType string

const (
	OpenflowDeploymentTypeSnowflake OpenflowDeploymentType = "SNOWFLAKE"
	OpenflowDeploymentTypeByoc      OpenflowDeploymentType = "BYOC"
)

var AllOpenflowDeploymentTypes = []OpenflowDeploymentType{
	OpenflowDeploymentTypeSnowflake,
	OpenflowDeploymentTypeByoc,
}

func ToOpenflowDeploymentType(s string) (OpenflowDeploymentType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllOpenflowDeploymentTypes, OpenflowDeploymentType(s)) {
		return "", fmt.Errorf("invalid openflow deployment type: %s", s)
	}
	return OpenflowDeploymentType(s), nil
}

type OpenflowVpcType string

const (
	OpenflowVpcTypeManaged  OpenflowVpcType = "MANAGED"
	OpenflowVpcTypeProvided OpenflowVpcType = "PROVIDED"
)

var AllOpenflowVpcTypes = []OpenflowVpcType{
	OpenflowVpcTypeManaged,
	OpenflowVpcTypeProvided,
}

func ToOpenflowVpcType(s string) (OpenflowVpcType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllOpenflowVpcTypes, OpenflowVpcType(s)) {
		return "", fmt.Errorf("invalid openflow vpc type: %s", s)
	}
	return OpenflowVpcType(s), nil
}

type OpenflowDeploymentStatus string

const (
	OpenflowDeploymentStatusCreating             OpenflowDeploymentStatus = "CREATING"
	OpenflowDeploymentStatusActive               OpenflowDeploymentStatus = "ACTIVE"
	OpenflowDeploymentStatusInactive             OpenflowDeploymentStatus = "INACTIVE"
	OpenflowDeploymentStatusProvisioning         OpenflowDeploymentStatus = "PROVISIONING"
	OpenflowDeploymentStatusNotReporting         OpenflowDeploymentStatus = "NOT_REPORTING"
	OpenflowDeploymentStatusNotHealthy           OpenflowDeploymentStatus = "NOT_HEALTHY"
	OpenflowDeploymentStatusUpgrading            OpenflowDeploymentStatus = "UPGRADING"
	OpenflowDeploymentStatusUpgradeFailed        OpenflowDeploymentStatus = "UPGRADE_FAILED"
	OpenflowDeploymentStatusDeactivationRequired OpenflowDeploymentStatus = "DEACTIVATION_REQUIRED"
	OpenflowDeploymentStatusDeleting             OpenflowDeploymentStatus = "DELETING"
	OpenflowDeploymentStatusDeleted              OpenflowDeploymentStatus = "DELETED"
	OpenflowDeploymentStatusCreateFailed         OpenflowDeploymentStatus = "CREATE_FAILED"
	OpenflowDeploymentStatusDeleteFailed         OpenflowDeploymentStatus = "DELETE_FAILED"
)

var allOpenflowDeploymentStatuses = []OpenflowDeploymentStatus{
	OpenflowDeploymentStatusCreating,
	OpenflowDeploymentStatusActive,
	OpenflowDeploymentStatusInactive,
	OpenflowDeploymentStatusProvisioning,
	OpenflowDeploymentStatusNotReporting,
	OpenflowDeploymentStatusNotHealthy,
	OpenflowDeploymentStatusUpgrading,
	OpenflowDeploymentStatusUpgradeFailed,
	OpenflowDeploymentStatusDeactivationRequired,
	OpenflowDeploymentStatusDeleting,
	OpenflowDeploymentStatusDeleted,
	OpenflowDeploymentStatusCreateFailed,
	OpenflowDeploymentStatusDeleteFailed,
}

func ToOpenflowDeploymentStatus(s string) (OpenflowDeploymentStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allOpenflowDeploymentStatuses, OpenflowDeploymentStatus(s)) {
		return "", fmt.Errorf("invalid openflow deployment status: %s", s)
	}
	return OpenflowDeploymentStatus(s), nil
}

func (r *CreateOpenflowDeploymentRequest) GetName() AccountObjectIdentifier {
	return r.name
}
