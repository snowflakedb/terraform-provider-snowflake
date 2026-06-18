package sdk

import "errors"

var (
	_ validatable = (*CreateRoleOptions)(nil)
	_ validatable = (*AlterRoleOptions)(nil)
	_ validatable = (*DropRoleOptions)(nil)
	_ validatable = (*ShowRoleOptions)(nil)
	_ validatable = (*GrantRoleOptions)(nil)
	_ validatable = (*RevokeRoleOptions)(nil)
)

func (opts *CreateRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateRoleOptions", "OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueNil(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errAtLeastOneOf("AlterRoleOptions", "RenameTo", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	}
	if anyValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) &&
		!exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterRoleOptions", "RenameTo", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	}
	return errors.Join(errs...)
}

func (opts *DropRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.InClass) && !ValidObjectIdentifier(opts.InClass.Class) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *GrantRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if (opts.Grant.Role != nil && opts.Grant.User != nil) || (opts.Grant.Role == nil && opts.Grant.User == nil) {
		errs = append(errs, errOneOf("GrantRoleOptions.Grant", "Role", "User"))
	}
	if opts.Grant.Role != nil && !ValidObjectIdentifier(opts.Grant.Role) {
		errs = append(errs, errInvalidIdentifier("GrantRoleOptions.Grant", "Role"))
	}
	if opts.Grant.User != nil && !ValidObjectIdentifier(opts.Grant.User) {
		errs = append(errs, errInvalidIdentifier("GrantRoleOptions.Grant", "User"))
	}
	return errors.Join(errs...)
}

func (opts *RevokeRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if (opts.Revoke.Role != nil && opts.Revoke.User != nil) || (opts.Revoke.Role == nil && opts.Revoke.User == nil) {
		errs = append(errs, errOneOf("RevokeRoleOptions.Revoke", "Role", "User"))
	}
	return errors.Join(errs...)
}
