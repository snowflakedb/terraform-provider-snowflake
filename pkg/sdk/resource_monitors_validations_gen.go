package sdk

import (
	"errors"
	"fmt"
)

func (opts *CreateResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateResourceMonitorOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, errors.Join(ErrInvalidObjectIdentifier))
	}
	if valueSet(opts.With) && everyValueNil(opts.With.CreditQuota, opts.With.Frequency, opts.With.StartTimestamp, opts.With.EndTimestamp, opts.With.NotifyUsers) && valueSet(opts.With.Triggers) {
		errs = append(errs, fmt.Errorf("due to Snowflake limitations you cannot create Resource Monitor with only triggers set"))
	}
	return errors.Join(errs...)
}

func (opts *AlterResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueNil(opts.Set, opts.Unset, opts.Triggers) {
		errs = append(errs, errAtLeastOneOf("AlterResourceMonitorOptions", "Set", "Unset", "Triggers"))
	}
	if set := opts.Set; valueSet(set) {
		if everyValueNil(set.CreditQuota, set.Frequency, set.StartTimestamp, set.EndTimestamp, set.NotifyUsers) {
			errs = append(errs, errAtLeastOneOf("ResourceMonitorSet", "CreditQuota", "Frequency", "StartTimestamp", "EndTimestamp", "NotifyUsers"))
		}
		if (set.Frequency != nil && set.StartTimestamp == nil) || (set.Frequency == nil && set.StartTimestamp != nil) {
			errs = append(errs, errors.New("must specify frequency and start time together"))
		}
	}
	return errors.Join(errs...)
}

func (opts *DropResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}
