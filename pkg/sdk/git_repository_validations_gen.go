package sdk

var (
	_ validatable = new(CreateGitRepositoryOptions)
	_ validatable = new(AlterGitRepositoryOptions)
	_ validatable = new(DropGitRepositoryOptions)
	_ validatable = new(DescribeGitRepositoryOptions)
	_ validatable = new(ShowGitRepositoryOptions)
)

func (opts *CreateGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateGitRepositoryOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.Fetch) {
		errs = append(errs, errExactlyOneOf("AlterGitRepositoryOptions", "Set", "Unset", "SetTags", "UnsetTags", "Fetch"))
	}
	return JoinErrors(errs...)
}

func (opts *DropGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
