package sdk

var (
	_ validatable = new(CreateTaskOptions)
	_ validatable = new(CreateOrAlterTaskOptions)
	_ validatable = new(CloneTaskOptions)
	_ validatable = new(AlterTaskOptions)
	_ validatable = new(DropTaskOptions)
	_ validatable = new(ShowTaskOptions)
	_ validatable = new(DescribeTaskOptions)
	_ validatable = new(ExecuteTaskOptions)
)

func (opts *CreateTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Warehouse) {
		if !exactlyOneValueSet(opts.Warehouse.Warehouse, opts.Warehouse.UserTaskManagedInitialWarehouseSize) {
			errs = append(errs, errExactlyOneOf("CreateTaskOptions.Warehouse", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
		}
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateTaskOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ErrorIntegration != nil && !ValidObjectIdentifier(opts.ErrorIntegration) {
		errs = append(errs, errInvalidIdentifier("CreateTaskOptions", "ErrorIntegration"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOrAlterTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Warehouse) {
		if !exactlyOneValueSet(opts.Warehouse.Warehouse, opts.Warehouse.UserTaskManagedInitialWarehouseSize) {
			errs = append(errs, errExactlyOneOf("CreateOrAlterTaskOptions.CreateTaskWarehouse", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
		}
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.ErrorIntegration != nil && !ValidObjectIdentifier(opts.ErrorIntegration) {
		errs = append(errs, errInvalidIdentifier("CreateOrAlterTaskOptions", "ErrorIntegration"))
	}
	return JoinErrors(errs...)
}

func (opts *CloneTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.sourceTask) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Resume, opts.Suspend, opts.RemoveAfter, opts.AddAfter, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.SetFinalize, opts.UnsetFinalize, opts.ModifyAs, opts.ModifyWhen, opts.RemoveWhen) {
		errs = append(errs, errExactlyOneOf("AlterTaskOptions", "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "SetFinalize", "UnsetFinalize", "ModifyAs", "ModifyWhen", "RemoveWhen"))
	}
	if valueSet(opts.Set) {
		if valueSet(opts.Set.SessionParameters) {
			if err := opts.Set.SessionParameters.validate(); err != nil {
				errs = append(errs, err)
			}
		}
		if !anyValueSet(opts.Set.Warehouse, opts.Set.UserTaskManagedInitialWarehouseSize, opts.Set.Schedule, opts.Set.Config, opts.Set.AllowOverlappingExecution, opts.Set.UserTaskTimeoutMs, opts.Set.SuspendTaskAfterNumFailures, opts.Set.ErrorIntegration, opts.Set.Comment, opts.Set.SessionParameters, opts.Set.TaskAutoRetryAttempts, opts.Set.UserTaskMinimumTriggerIntervalInSeconds, opts.Set.TargetCompletionInterval, opts.Set.ServerlessTaskMinStatementSize, opts.Set.ServerlessTaskMaxStatementSize) {
			errs = append(errs, errAtLeastOneOf("AlterTaskOptions.Set", "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParameters", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize"))
		}
		if everyValueSet(opts.Set.Warehouse, opts.Set.UserTaskManagedInitialWarehouseSize) {
			errs = append(errs, errOneOf("AlterTaskOptions.Set", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
		}
		if opts.Set.ErrorIntegration != nil && !ValidObjectIdentifier(opts.Set.ErrorIntegration) {
			errs = append(errs, errInvalidIdentifier("AlterTaskOptions.Set", "ErrorIntegration"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Warehouse, opts.Unset.UserTaskManagedInitialWarehouseSize, opts.Unset.Schedule, opts.Unset.Config, opts.Unset.AllowOverlappingExecution, opts.Unset.UserTaskTimeoutMs, opts.Unset.SuspendTaskAfterNumFailures, opts.Unset.ErrorIntegration, opts.Unset.Comment, opts.Unset.SessionParametersUnset, opts.Unset.TaskAutoRetryAttempts, opts.Unset.UserTaskMinimumTriggerIntervalInSeconds, opts.Unset.TargetCompletionInterval, opts.Unset.ServerlessTaskMinStatementSize, opts.Unset.ServerlessTaskMaxStatementSize) {
			errs = append(errs, errAtLeastOneOf("AlterTaskOptions.Unset", "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParametersUnset", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize"))
		}
		if valueSet(opts.Unset.SessionParametersUnset) {
			if err := opts.Unset.SessionParametersUnset.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ExecuteTaskOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
