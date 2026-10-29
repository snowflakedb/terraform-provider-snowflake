package sdk

import "errors"

// additionalValidations for CreatePipeOptions — copyStatement must be non-empty.
// The generator cannot express "required string must be non-empty" natively.
func (opts *CreatePipeOptions) additionalValidations() error {
	if opts.copyStatement == "" {
		return errors.Join(errNotSet("CreatePipeOptions", "copyStatement"))
	}
	return nil
}

// additionalValidations for ShowPipeOptions — Like requires Pattern; In requires exactly one scope.
func (opts *ShowPipeOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"))
	}
	return errors.Join(errs...)
}
