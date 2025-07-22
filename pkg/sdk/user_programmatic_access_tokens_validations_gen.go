package sdk

var (
	_ validatable = new(AddUserProgrammaticAccessTokenOptions)
	_ validatable = new(ModifyUserProgrammaticAccessTokenOptions)
	_ validatable = new(RotateUserProgrammaticAccessTokenOptions)
	_ validatable = new(RemoveUserProgrammaticAccessTokenOptions)
	_ validatable = new(ShowUserProgrammaticAccessTokenOptions)
)

func (opts *AddUserProgrammaticAccessTokenOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// adjusted manually
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("AddUserProgrammaticAccessTokenOptions", "UserName"))
	}
	// adjusted manually
	if valueSet(opts.DaysToExpiry) {
		if !validateIntGreaterThanOrEqual(*opts.DaysToExpiry, 1) {
			errs = append(errs, errIntValue("AddUserProgrammaticAccessTokenOptions", "DaysToExpiry", IntErrGreaterOrEqual, 1))
		}
	}
	// adjusted manually
	if valueSet(opts.MinsToBypassNetworkPolicyRequirement) {
		if !validateIntGreaterThanOrEqual(*opts.MinsToBypassNetworkPolicyRequirement, 1) {
			errs = append(errs, errIntValue("AddUserProgrammaticAccessTokenOptions", "MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
		}
	}
	if opts.RoleRestriction != nil && !ValidObjectIdentifier(opts.RoleRestriction) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ModifyUserProgrammaticAccessTokenOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// adjusted manually
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("ModifyUserProgrammaticAccessTokenOptions", "UserName"))
	}
	// adjusted manually
	if valueSet(opts.Set) && valueSet(opts.Set.MinsToBypassNetworkPolicyRequirement) {
		if !validateIntGreaterThanOrEqual(*opts.Set.MinsToBypassNetworkPolicyRequirement, 1) {
			errs = append(errs, errIntValue("ModifyUserProgrammaticAccessTokenOptions", "Set.MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
		}
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.RenameTo) {
		errs = append(errs, errExactlyOneOf("ModifyUserProgrammaticAccessTokenOptions", "Set", "Unset", "RenameTo"))
	}
	return JoinErrors(errs...)
}

func (opts *RotateUserProgrammaticAccessTokenOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// adjusted manually
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("RotateUserProgrammaticAccessTokenOptions", "UserName"))
	}
	// adjusted manually
	if valueSet(opts.ExpireRotatedTokenAfterHours) {
		if !validateIntGreaterThanOrEqual(*opts.ExpireRotatedTokenAfterHours, 0) {
			errs = append(errs, errIntValue("RotateUserProgrammaticAccessTokenOptions", "ExpireRotatedTokenAfterHours", IntErrGreaterOrEqual, 0))
		}
	}
	return JoinErrors(errs...)
}

func (opts *RemoveUserProgrammaticAccessTokenOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// adjusted manually
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("RemoveUserProgrammaticAccessTokenOptions", "UserName"))
	}
	return JoinErrors(errs...)
}

func (opts *ShowUserProgrammaticAccessTokenOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
