package sdk

import "errors"

var (
	_ validatable = new(AlterSessionOptions)
)

func (opts *AlterSessionOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueNil(opts.Set, opts.Unset) {
		errs = append(errs, errOneOf("AlterSessionOptions", "Set", "Unset"))
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

func (v *SessionSet) validate() error {
	if err := v.SessionParameters.validate(); err != nil {
		return err
	}
	return nil
}

func (v *SessionUnset) validate() error {
	if err := v.SessionParametersUnset.validate(); err != nil {
		return err
	}
	return nil
}
