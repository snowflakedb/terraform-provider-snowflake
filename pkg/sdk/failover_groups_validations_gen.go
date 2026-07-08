package sdk

import (
	"errors"
	"slices"
)

var (
	_ validatable = new(CreateFailoverGroupOptions)
	_ validatable = new(CreateSecondaryReplicationGroupOptions)
	_ validatable = new(AlterSourceFailoverGroupOptions)
	_ validatable = new(AlterTargetFailoverGroupOptions)
	_ validatable = new(DropFailoverGroupOptions)
	_ validatable = new(ShowFailoverGroupOptions)
)

func (opts *CreateFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *CreateSecondaryReplicationGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.primaryFailoverGroup) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryReplicationGroupOptions", "primaryFailoverGroup"))
	}
	return errors.Join(errs...)
}

func (opts *AlterSourceFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.Add, opts.Move, opts.Remove, opts.NewName) {
		errs = append(errs, errExactlyOneOf("AlterSourceFailoverGroupOptions", "Set", "Unset", "Add", "Move", "Remove", "NewName"))
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
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Move) {
		if err := opts.Move.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *FailoverGroupSet) validate() error {
	if len(v.AllowedIntegrationTypes) > 0 {
		// INTEGRATIONS must be set in object types
		if !slices.Contains(v.ObjectTypes, PluralObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

func (v *FailoverGroupUnset) validate() error {
	if everyValueNil(v.ReplicationSchedule) {
		return errAtLeastOneOf("FailoverGroupUnset", "ReplicationSchedule")
	}
	return nil
}

func (v *FailoverGroupAdd) validate() error {
	return nil
}

func (v *FailoverGroupMove) validate() error {
	return nil
}

func (v *FailoverGroupRemove) validate() error {
	return nil
}

func (opts *AlterTargetFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Refresh, opts.Primary, opts.Suspend, opts.Resume) {
		errs = append(errs, errExactlyOneOf("AlterTargetFailoverGroupOptions", "Refresh", "Primary", "Suspend", "Resume"))
	}
	return errors.Join(errs...)
}

func (opts *DropFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowFailoverGroupOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}
