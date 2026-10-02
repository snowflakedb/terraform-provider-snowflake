package sdk

func (opts *CreateNotebookOptions) additionalValidations() error {
	var errs []error
	if opts.IdleAutoShutdownTimeSeconds != nil && !validateIntGreaterThan(*opts.IdleAutoShutdownTimeSeconds, 0) {
		errs = append(errs, errIntValue("CreateNotebookOptions", "IdleAutoShutdownTimeSeconds", IntErrGreater, 0))
	}
	return JoinErrors(errs...)
}

func (s *NotebookSet) additionalValidations() error {
	var errs []error
	if s.IdleAutoShutdownTimeSeconds != nil && !validateIntGreaterThan(*s.IdleAutoShutdownTimeSeconds, 0) {
		errs = append(errs, errIntValue("AlterNotebookOptions", "IdleAutoShutdownTimeSeconds", IntErrGreater, 0))
	}
	return JoinErrors(errs...)
}

func (r *CreateNotebookRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r NotebookDetailsRow) additionalConvert(_ *NotebookDetails) error {
	return nil
}

func (d *NotebookDetails) ID() SchemaObjectIdentifier {
	return d.Id
}
