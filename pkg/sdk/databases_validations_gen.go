package sdk

import (
	"errors"
	"fmt"
)

func (opts *CreateDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Clone) {
		if err := opts.Clone.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateDatabaseOptions", "Catalog"))
	}
	return errors.Join(errs...)
}

func (opts *CreateSharedDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSharedDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "Catalog"))
	}
	if !ValidObjectIdentifier(opts.fromShare) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "fromShare"))
	}
	return errors.Join(errs...)
}

func (opts *CreateSecondaryDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.primaryDatabase) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "primaryDatabase"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSecondaryDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "Catalog"))
	}
	return errors.Join(errs...)
}

func (opts *CreateDatabaseFromListingOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.fromListing == "" {
		errs = append(errs, fmt.Errorf("CreateDatabaseFromListingOptions: listing global name must not be empty"))
	}
	return errors.Join(errs...)
}

func (opts *AlterDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.NewName != nil && !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AlterDatabaseOptions", "NewName"))
	}
	if opts.SwapWith != nil && !ValidObjectIdentifier(opts.SwapWith) {
		errs = append(errs, errInvalidIdentifier("AlterDatabaseOptions", "SwapWith"))
	}
	if !exactlyOneValueSet(opts.NewName, opts.Set, opts.Unset, opts.SwapWith, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
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

func (v *DatabaseSet) validate() error {
	var errs []error
	if v.ExternalVolume != nil && !ValidObjectIdentifier(v.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("DatabaseSet", "ExternalVolume"))
	}
	if v.Catalog != nil && !ValidObjectIdentifier(v.Catalog) {
		errs = append(errs, errInvalidIdentifier("DatabaseSet", "Catalog"))
	}
	if !anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ExternalVolume,
		v.Catalog,
		v.ReplaceInvalidCharacters,
		v.DefaultDDLCollation,
		v.StorageSerializationPolicy,
		v.LogLevel,
		v.LogEventLevel,
		v.TraceLevel,
		v.SuspendTaskAfterNumFailures,
		v.TaskAutoRetryAttempts,
		v.UserTaskManagedInitialWarehouseSize,
		v.UserTaskTimeoutMs,
		v.UserTaskMinimumTriggerIntervalInSeconds,
		v.QuotedIdentifiersIgnoreCase,
		v.EnableConsoleOutput,
		v.Comment,
	) {
		errs = append(errs, errAtLeastOneOf(
			"DatabaseSet",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"LogEventLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"Comment",
		))
	}
	return errors.Join(errs...)
}

func (v *DatabaseUnset) validate() error {
	var errs []error
	if !anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ExternalVolume,
		v.Catalog,
		v.ReplaceInvalidCharacters,
		v.DefaultDDLCollation,
		v.StorageSerializationPolicy,
		v.LogLevel,
		v.LogEventLevel,
		v.TraceLevel,
		v.SuspendTaskAfterNumFailures,
		v.TaskAutoRetryAttempts,
		v.UserTaskManagedInitialWarehouseSize,
		v.UserTaskTimeoutMs,
		v.UserTaskMinimumTriggerIntervalInSeconds,
		v.QuotedIdentifiersIgnoreCase,
		v.EnableConsoleOutput,
		v.Comment,
	) {
		errs = append(errs, errAtLeastOneOf(
			"DatabaseUnset",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"LogEventLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"Comment",
		))
	}
	return errors.Join(errs...)
}

func (opts *AlterDatabaseReplicationOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.EnableReplication, opts.DisableReplication, opts.Refresh) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	}
	return errors.Join(errs...)
}

func (opts *AlterDatabaseFailoverOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.EnableFailover, opts.DisableFailover, opts.Primary) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
	}
	return errors.Join(errs...)
}

func (opts *DropDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Cascade, opts.Restrict) {
		errs = append(errs, errOneOf("DropDatabaseOptions", "Cascade", "Restrict"))
	}
	return JoinErrors(errs...)
}

func (opts *undropDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowDatabasesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *describeDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}
