package sdk

import (
	"fmt"
	"slices"
	"strings"
)

func (s *CreateServiceRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func (s *ExecuteJobServiceRequest) GetName() SchemaObjectIdentifier {
	return s.Name
}

type ServiceType string

const (
	ServiceTypeService    ServiceType = "SERVICE"
	ServiceTypeJobService ServiceType = "JOB_SERVICE"
)

func (s Service) Type() ServiceType {
	if s.IsJob {
		return ServiceTypeJobService
	}
	return ServiceTypeService
}

var allServiceTypes = []ServiceType{
	ServiceTypeService,
	ServiceTypeJobService,
}

func ToServiceType(s string) (ServiceType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allServiceTypes, ServiceType(s)) {
		return "", fmt.Errorf("invalid service type: %s", s)
	}
	return ServiceType(s), nil
}
