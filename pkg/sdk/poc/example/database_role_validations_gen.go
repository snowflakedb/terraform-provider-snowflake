package example

import "errors"

var (
	_ validatable = new(CreateDatabaseRoleOptions)
	_ validatable = new(AlterDatabaseRoleOptions)
)

func (opts *CreateDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Rename, opts.Set, opts.Unset) {
		errs = append(errs, errOneOf("Rename", "Set", "Unset"))
	}
	if valueSet(opts.Rename) {
		if !validObjectidentifier(opts.Rename.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if valueSet(opts.Set) {
		if valueSet(opts.Set.NestedThirdLevel) {
			if ok := anyValueSet(opts.Set.NestedThirdLevel.Field); !ok {
				errs = append(errs, errAtLeastOneOf("Field"))
			}
		}
	}
	if valueSet(opts.Unset) {
		if ok := anyValueSet(opts.Unset.Comment); !ok {
			errs = append(errs, errAtLeastOneOf("Comment"))
		}
		if valueSet(opts.Unset.NestedThirdLevel) {
			if ok := anyValueSet(opts.Unset.NestedThirdLevel.Field); !ok {
				errs = append(errs, errAtLeastOneOf("Field"))
			}
		}
	}
	return errors.Join(errs...)
}
