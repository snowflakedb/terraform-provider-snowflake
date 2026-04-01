package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = new(createTagOptions)
	_ validatable = new(alterTagOptions)
	_ validatable = new(showTagOptions)
	_ validatable = new(dropTagOptions)
	_ validatable = new(undropTagOptions)
	_ validatable = new(AllowedValues)
	_ validatable = new(TagPropagate)
	_ validatable = new(TagOnConflict)
	_ validatable = new(TagSet)
	_ validatable = new(TagUnset)
)

func (opts *createTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("createTagOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.AllowedValues) {
		if err := opts.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Propagate) {
		if err := opts.Propagate.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *AllowedValues) validate() error {
	if !validateIntInRangeInclusive(len(v.Values), 1, 300) {
		return errIntBetween("AllowedValues", "Values", 1, 300)
	}
	return nil
}

func (v *TagPropagate) validate() error {
	if v == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(v.PropagationMethod) {
		errs = append(errs, errNotSet("TagPropagate", "PropagationMethod"))
	}
	if valueSet(v.OnConflict) {
		if err := v.OnConflict.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *TagOnConflict) validate() error {
	if v == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !exactlyOneValueSet(v.CustomValue, v.AllowedValuesSequence) {
		errs = append(errs, errExactlyOneOf("TagOnConflict", "CustomValue", "AllowedValuesSequence"))
	}
	return errors.Join(errs...)
}

func (v *TagSet) validate() error {
	var errs []error
	if valueSet(v.MaskingPolicies) && anyValueSet(v.AllowedValues, v.Propagate, v.Comment) {
		errs = append(errs, errOneOf("TagSet", "MaskingPolicies", "AllowedValues", "Propagate", "Comment"))
	}
	if !anyValueSet(v.MaskingPolicies, v.AllowedValues, v.Propagate, v.Comment) {
		errs = append(errs, errAtLeastOneOf("TagSet", "MaskingPolicies", "AllowedValues", "Propagate", "Comment"))
	}
	if valueSet(v.MaskingPolicies) {
		if !validateIntGreaterThan(len(v.MaskingPolicies.MaskingPolicies), 0) {
			errs = append(errs, errIntValue("TagSet.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
		}
	}
	if valueSet(v.AllowedValues) {
		if err := v.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Propagate) {
		if err := v.Propagate.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *TagUnset) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.MaskingPolicies, v.AllowedValues, v.Propagate, v.OnConflict, v.Comment) {
		errs = append(errs, errExactlyOneOf("TagUnset", "MaskingPolicies", "AllowedValues", "Propagate", "OnConflict", "Comment"))
	}
	if valueSet(v.MaskingPolicies) {
		if !validateIntGreaterThan(len(v.MaskingPolicies.MaskingPolicies), 0) {
			errs = append(errs, errIntValue("TagUnset.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
		}
	}
	return errors.Join(errs...)
}

func (opts *alterTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Add, opts.Drop, opts.Set, opts.Unset, opts.Rename) {
		errs = append(errs, errExactlyOneOf("alterTagOptions", "Add", "Drop", "Set", "Unset", "Rename"))
	}
	if valueSet(opts.Add) && valueSet(opts.Add.AllowedValues) {
		if err := opts.Add.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Drop) && valueSet(opts.Drop.AllowedValues) {
		if err := opts.Drop.AllowedValues.validate(); err != nil {
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
	if valueSet(opts.Rename) {
		if !ValidObjectIdentifier(opts.Rename.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (opts *showTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("showTagOptions.In", "Account", "Database", "Schema"))
	}
	return errors.Join(errs...)
}

func (opts *dropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *undropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *setTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.objectName) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !canBeAssociatedWithTag(opts.objectType) {
		return fmt.Errorf("tagging for object type %s is not supported", opts.objectType)
	}
	return errors.Join(errs...)
}

func (opts *unsetTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.objectName) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !canBeAssociatedWithTag(opts.objectType) {
		return fmt.Errorf("tagging for object type %s is not supported", opts.objectType)
	}
	return errors.Join(errs...)
}
