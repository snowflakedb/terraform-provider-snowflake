package sdk

var (
	_ validatable = new(CreateListingOptions)
	_ validatable = new(CreateFromStageListingOptions)
	_ validatable = new(AlterListingOptions)
	_ validatable = new(DropListingOptions)
	_ validatable = new(ShowListingOptions)
	_ validatable = new(DescribeListingOptions)
)

func (opts *CreateListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.With) {
		if !exactlyOneValueSet(opts.With.Share, opts.With.ApplicationPackage) {
			errs = append(errs, errExactlyOneOf("CreateListingOptions.With", "Share", "ApplicationPackage"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateFromStageListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.With) {
		if !exactlyOneValueSet(opts.With.Share, opts.With.ApplicationPackage) {
			errs = append(errs, errExactlyOneOf("CreateFromStageListingOptions.With", "Share", "ApplicationPackage"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfExists, opts.AddVersion) {
		errs = append(errs, errOneOf("AlterListingOptions", "IfExists", "AddVersion"))
	}
	if !exactlyOneValueSet(opts.Publish, opts.Unpublish, opts.Review, opts.AlterListingAs, opts.AddVersion, opts.RenameTo, opts.Set) {
		errs = append(errs, errExactlyOneOf("AlterListingOptions", "Publish", "Unpublish", "Review", "AlterListingAs", "AddVersion", "RenameTo", "Set"))
	}
	return JoinErrors(errs...)
}

func (opts *DropListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeListingOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
