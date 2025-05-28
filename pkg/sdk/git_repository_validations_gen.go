package sdk

var (
	_ validatable = new(CreateGitRepositoryOptions)
	_ validatable = new(AlterGitRepositoryOptions)
	_ validatable = new(DropGitRepositoryOptions)
	_ validatable = new(DescribeGitRepositoryOptions)
	_ validatable = new(ShowGitRepositoryOptions)
	_ validatable = new(ShowGitBranchesGitRepositoryOptions)
	_ validatable = new(ShowGitTagsGitRepositoryOptions)
)

func (opts *CreateGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.ApiIntegration) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.GitCredentials != nil && !ValidObjectIdentifier(opts.GitCredentials) {
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
	if valueSet(opts.Set) {
		if opts.Set.ApiIntegration != nil && !ValidObjectIdentifier(opts.Set.ApiIntegration) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if opts.Set.GitCredentials != nil && !ValidObjectIdentifier(opts.Set.GitCredentials) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
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

func (opts *ShowGitBranchesGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *ShowGitTagsGitRepositoryOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
