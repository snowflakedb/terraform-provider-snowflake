package sdk

func (opts *CreateApplicationOptions) additionalValidations() error {
	if valueSet(opts.DebugMode) && !valueSet(opts.Version) {
		return NewError("CreateApplicationOptions.DebugMode can be set only when CreateApplicationOptions.Version is set")
	}
	return nil
}

func (opts *AlterApplicationOptions) additionalValidations() error {
	if valueSet(opts.IfExists) {
		if !valueSet(opts.Set) && !valueSet(opts.Unset) {
			return NewError("AlterApplicationOptions.IfExists can be set only when AlterApplicationOptions.Set or AlterApplicationOptions.Unset is set")
		}
	}
	return nil
}
