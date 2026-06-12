package sdk

import "errors"

var (
	_ validatable = new(createDatabaseRoleOptions)
	_ validatable = new(alterDatabaseRoleOptions)
	_ validatable = new(dropDatabaseRoleOptions)
	_ validatable = new(showDatabaseRoleOptions)
	_ validatable = new(grantDatabaseRoleOptions)
	_ validatable = new(revokeDatabaseRoleOptions)
	_ validatable = new(grantDatabaseRoleToShareOptions)
	_ validatable = new(revokeDatabaseRoleFromShareOptions)
)

func (opts *createDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("createDatabaseRoleOptions", "OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *alterDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Rename, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("alterDatabaseRoleOptions", "Rename", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if opts.Rename != nil {
		if !ValidObjectIdentifier(opts.Rename) {
			errs = append(errs, errInvalidIdentifier("alterDatabaseRoleOptions", "Rename"))
		}
		if opts.name.DatabaseName() != opts.Rename.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !unset.Comment {
			errs = append(errs, errAtLeastOneOf("alterDatabaseRoleOptions.Unset", "Comment"))
		}
	}
	return errors.Join(errs...)
}

func (opts *dropDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *showDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.Database) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}

func (opts *grantDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.ParentRole.DatabaseRoleName, opts.ParentRole.AccountRoleName); !ok {
		errs = append(errs, errOneOf("DatabaseRoleName", "AccountRoleName"))
	}
	return errors.Join(errs...)
}

func (opts *revokeDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.ParentRole.DatabaseRoleName, opts.ParentRole.AccountRoleName); !ok {
		errs = append(errs, errOneOf("DatabaseRoleName", "AccountRoleName"))
	}
	return errors.Join(errs...)
}

func (opts *grantDatabaseRoleToShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.Share) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *revokeDatabaseRoleFromShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.Share) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
