package sdk

var (
	_ validatable = new(CreateOpenflowConnectorOptions)
	_ validatable = new(AlterOpenflowConnectorOptions)
	_ validatable = new(DropOpenflowConnectorOptions)
	_ validatable = new(ShowOpenflowConnectorOptions)
	_ validatable = new(DescribeOpenflowConnectorOptions)
)

func (opts *CreateOpenflowConnectorOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !anyValueSet(opts.FromDefinition, opts.FromStage) {
		errs = append(errs, errAtLeastOneOf("CreateOpenflowConnectorOptions", "FromDefinition", "FromStage"))
	}
	if valueSet(opts.FromDefinition) && valueSet(opts.FromStage) {
		errs = append(errs, errExactlyOneOf("CreateOpenflowConnectorOptions", "FromDefinition", "FromStage"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterOpenflowConnectorOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Start, opts.Stop, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterOpenflowConnectorOptions", "Start", "Stop", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.DisplayName, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterOpenflowConnectorOptions.Set", "DisplayName", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.DisplayName, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterOpenflowConnectorOptions.Unset", "DisplayName", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropOpenflowConnectorOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowOpenflowConnectorOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	return nil
}

func (opts *DescribeOpenflowConnectorOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
