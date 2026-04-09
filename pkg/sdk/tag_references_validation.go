package sdk

import (
	"errors"
)

var _ validatable = new(getForEntityTagReferenceOptions)

func (opts *getForEntityTagReferenceOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.parameters) {
		errs = append(errs, errNotSet("getForEntityTagReferenceOptions", "parameters"))
	} else {
		if !valueSet(opts.parameters.arguments) {
			errs = append(errs, errNotSet("tagReferenceParameters", "arguments"))
		} else {
			if !valueSet(opts.parameters.arguments.objectName) {
				errs = append(errs, errNotSet("tagReferenceFunctionArguments", "objectName"))
			}
			if !valueSet(opts.parameters.arguments.objectDomain) {
				errs = append(errs, errNotSet("tagReferenceFunctionArguments", "objectDomain"))
			}
		}
	}
	return errors.Join(errs...)
}
