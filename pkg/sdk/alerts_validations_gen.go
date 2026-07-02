package sdk

import "errors"

var (
	_ validatable = new(CreateAlertOptions)
	_ validatable = new(AlterAlertOptions)
	_ validatable = new(DropAlertOptions)
	_ validatable = new(ShowAlertOptions)
)

func (opts *CreateAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *AlterAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Action, opts.Set, opts.Unset, opts.ModifyCondition, opts.ModifyAction) {
		errs = append(errs, errExactlyOneOf("AlterAlertOptions", "Action", "Set", "Unset", "ModifyCondition", "ModifyAction"))
	}
	return errors.Join(errs...)
}

func (opts *DropAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *describeAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}
