package sdk

import "errors"

var (
	_ validatable = new(CreateMaskingPolicyOptions)
	_ validatable = new(AlterMaskingPolicyOptions)
	_ validatable = new(DropMaskingPolicyOptions)
	_ validatable = new(ShowMaskingPolicyOptions)
	_ validatable = new(describeMaskingPolicyOptions)
)

func (opts *CreateMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("CreateMaskingPolicyOptions", "OrReplace", "IfNotExists"))
	}
	if !valueSet(opts.signature) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "signature"))
	}
	if !valueSet(opts.returns) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "returns"))
	}
	if !valueSet(opts.body) {
		errs = append(errs, errNotSet("CreateMaskingPolicyOptions", "body"))
	}
	return errors.Join(errs...)
}

func (opts *AlterMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.NewName != nil {
		if !ValidObjectIdentifier(opts.NewName) {
			errs = append(errs, errInvalidIdentifier("AlterMaskingPolicyOptions", "NewName"))
		}
		if opts.name.DatabaseName() != opts.NewName.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
		if opts.name.SchemaName() != opts.NewName.SchemaName() {
			errs = append(errs, ErrDifferentSchema)
		}
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.NewName) {
		errs = append(errs, errExactlyOneOf("AlterMaskingPolicyOptions", "Set", "Unset", "SetTag", "UnsetTag", "NewName"))
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

func (v *MaskingPolicySet) validate() error {
	if !exactlyOneValueSet(v.Body, v.Comment) {
		return errExactlyOneOf("MaskingPolicySet", "Body", "Comment")
	}
	return nil
}

func (v *MaskingPolicyUnset) validate() error {
	if !exactlyOneValueSet(v.Comment) {
		return errExactlyOneOf("MaskingPolicyUnset", "Comment")
	}
	return nil
}

func (opts *DropMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *describeMaskingPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}
