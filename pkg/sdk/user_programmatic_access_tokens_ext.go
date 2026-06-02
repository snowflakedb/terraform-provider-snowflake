package sdk

import "context"

func (opts *AddUserProgrammaticAccessTokenOptions) additionalValidations() error {
	var errs []error
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("AddUserProgrammaticAccessTokenOptions", "UserName"))
	}
	if valueSet(opts.DaysToExpiry) {
		if !validateIntGreaterThanOrEqual(*opts.DaysToExpiry, 1) {
			errs = append(errs, errIntValue("AddUserProgrammaticAccessTokenOptions", "DaysToExpiry", IntErrGreaterOrEqual, 1))
		}
	}
	if valueSet(opts.MinsToBypassNetworkPolicyRequirement) {
		if !validateIntGreaterThanOrEqual(*opts.MinsToBypassNetworkPolicyRequirement, 1) {
			errs = append(errs, errIntValue("AddUserProgrammaticAccessTokenOptions", "MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ModifyUserProgrammaticAccessTokenOptions) additionalValidations() error {
	var errs []error
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("ModifyUserProgrammaticAccessTokenOptions", "UserName"))
	}
	if valueSet(opts.Set) && valueSet(opts.Set.MinsToBypassNetworkPolicyRequirement) {
		if !validateIntGreaterThanOrEqual(*opts.Set.MinsToBypassNetworkPolicyRequirement, 1) {
			errs = append(errs, errIntValue("ModifyUserProgrammaticAccessTokenOptions", "Set.MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
		}
	}
	return JoinErrors(errs...)
}

func (opts *RotateUserProgrammaticAccessTokenOptions) additionalValidations() error {
	var errs []error
	if !ValidObjectIdentifier(opts.UserName) {
		errs = append(errs, errInvalidIdentifier("RotateUserProgrammaticAccessTokenOptions", "UserName"))
	}
	if valueSet(opts.ExpireRotatedTokenAfterHours) {
		if !validateIntGreaterThanOrEqual(*opts.ExpireRotatedTokenAfterHours, 0) {
			errs = append(errs, errIntValue("RotateUserProgrammaticAccessTokenOptions", "ExpireRotatedTokenAfterHours", IntErrGreaterOrEqual, 0))
		}
	}
	return JoinErrors(errs...)
}

func (opts *RemoveUserProgrammaticAccessTokenOptions) additionalValidations() error {
	if !ValidObjectIdentifier(opts.UserName) {
		return errInvalidIdentifier("RemoveUserProgrammaticAccessTokenOptions", "UserName")
	}
	return nil
}

func (r *AddProgrammaticAccessTokenResult) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(r.TokenName)
}

func (v *ProgrammaticAccessToken) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *userProgrammaticAccessTokens) RemoveByIDSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return SafeRemoveProgrammaticAccessToken(v.client, ctx, request)
}
