package sdk

import "context"

func (v *computePools) dropSafelyHook(ctx context.Context, id AccountObjectIdentifier) error {
	return v.client.ComputePools.Alter(ctx, NewAlterComputePoolRequest(id).WithIfExists(true).WithStopAll(true))
}

func (opts *CreateComputePoolOptions) additionalValidations() error {
	var errs []error
	if !validateIntGreaterThan(opts.MinNodes, 0) {
		errs = append(errs, errIntValue("CreateComputePoolOptions", "MinNodes", IntErrGreater, 0))
	}
	if !validateIntGreaterThanOrEqual(opts.MaxNodes, opts.MinNodes) {
		errs = append(errs, errIntValue("CreateComputePoolOptions", "MaxNodes", IntErrGreaterOrEqual, opts.MinNodes))
	}
	return JoinErrors(errs...)
}

func (s *ComputePoolSet) additionalValidations() error {
	var errs []error
	if valueSet(s.MinNodes) && !validateIntGreaterThan(*s.MinNodes, 0) {
		errs = append(errs, errIntValue("AlterComputePoolOptions", "Set.MinNodes", IntErrGreater, 0))
	}
	if valueSet(s.MaxNodes) && !validateIntGreaterThan(*s.MaxNodes, 0) {
		errs = append(errs, errIntValue("AlterComputePoolOptions", "Set.MaxNodes", IntErrGreater, 0))
	}
	if valueSet(s.MinNodes) && valueSet(s.MaxNodes) && !validateIntGreaterThanOrEqual(*s.MaxNodes, *s.MinNodes) {
		errs = append(errs, errIntValue("AlterComputePoolOptions", "Set.MaxNodes", IntErrGreaterOrEqual, *s.MinNodes))
	}
	return JoinErrors(errs...)
}

func (r *CreateComputePoolRequest) GetName() AccountObjectIdentifier {
	return r.name
}
