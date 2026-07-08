package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = new(CreateShareOptions)
	_ validatable = new(AlterShareOptions)
	_ validatable = new(DropShareOptions)
	_ validatable = new(ShowShareOptions)
	_ validatable = new(describeShareOptions)
)

func (opts *CreateShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *DropShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *AlterShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Add, opts.Remove, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterShareOptions", "Add", "Remove", "Set", "Unset", "SetTag", "UnsetTag"))
	}
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *ShareAdd) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

func (v *ShareRemove) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

func (v *ShareSet) validate() error {
	if !anyValueSet(v.Accounts, v.Comment) {
		return errAtLeastOneOf("ShareSet", "Accounts", "Comment")
	}
	return nil
}

func (v *ShareUnset) validate() error {
	if !exactlyOneValueSet(v.Comment) {
		return errExactlyOneOf("ShareUnset", "Comment")
	}
	return nil
}

func (opts *ShowShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *describeShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}
