package sdk

import (
	"fmt"
	"slices"
	"strings"
)

func (s *ServiceDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(s.DatabaseName, s.SchemaName, s.Name)
}

func (opts *CreateServiceOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.MinReadyInstances) {
		if !validateIntGreaterThan(*opts.MinReadyInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinReadyInstances", IntErrGreater, 0))
		}
		if valueSet(opts.MinInstances) && !validateIntGreaterThanOrEqual(*opts.MinInstances, *opts.MinReadyInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreaterOrEqual, *opts.MinReadyInstances))
		}
		if valueSet(opts.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.MaxInstances, *opts.MinReadyInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, *opts.MinReadyInstances))
		}
	}
	if valueSet(opts.MinInstances) {
		if !validateIntGreaterThan(*opts.MinInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreater, 0))
		}
		if valueSet(opts.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.MaxInstances, *opts.MinInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, *opts.MinInstances))
		}
	}
	if valueSet(opts.MaxInstances) {
		if !validateIntGreaterThan(*opts.MaxInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreater, 0))
		}
	}
	if valueSet(opts.AutoSuspendSecs) && !validateIntGreaterThanOrEqual(*opts.AutoSuspendSecs, 0) {
		errs = append(errs, errIntValue("CreateServiceOptions", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
	}
	return JoinErrors(errs...)
}

func (s *ServiceSet) additionalValidations() error {
	var errs []error
	if valueSet(s.MinReadyInstances) {
		if !validateIntGreaterThan(*s.MinReadyInstances, 0) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinReadyInstances", IntErrGreater, 0))
		}
		if valueSet(s.MinInstances) && !validateIntGreaterThanOrEqual(*s.MinInstances, *s.MinReadyInstances) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreaterOrEqual, *s.MinReadyInstances))
		}
		if valueSet(s.MaxInstances) && !validateIntGreaterThanOrEqual(*s.MaxInstances, *s.MinReadyInstances) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, *s.MinReadyInstances))
		}
	}
	if valueSet(s.MinInstances) {
		if !validateIntGreaterThan(*s.MinInstances, 0) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreater, 0))
		}
		if valueSet(s.MaxInstances) && !validateIntGreaterThanOrEqual(*s.MaxInstances, *s.MinInstances) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, *s.MinInstances))
		}
	}
	if valueSet(s.MaxInstances) {
		if !validateIntGreaterThan(*s.MaxInstances, 0) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreater, 0))
		}
	}
	if valueSet(s.AutoSuspendSecs) && !validateIntGreaterThanOrEqual(*s.AutoSuspendSecs, 0) {
		errs = append(errs, errIntValue("AlterServiceOptions.Set", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
	}
	return JoinErrors(errs...)
}

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

func (v *Service) Type() ServiceType {
	if v.IsJob {
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
