package sdk

import (
	"context"
	"errors"
)

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
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema, opts.In.Application, opts.In.ApplicationPackage) {
		errs = append(errs, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema", "Application", "ApplicationPackage"))
	}
	return errors.Join(errs...)
}

func (v *pipes) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Pipe: id,
		},
	})
}
