package sdk

func (opts *AlterHybridTableOptions) additionalValidations() error {
	var errs []error
	for _, action := range opts.AlterColumnAction {
		if !exactlyOneValueSet(action.DropDefault, action.SetDefault, action.Type, action.Comment, action.UnsetComment) {
			errs = append(errs, errExactlyOneOf("AlterHybridTableOptions.AlterColumnAction", "DropDefault", "SetDefault", "Type", "Comment", "UnsetComment"))
		}
	}
	return JoinErrors(errs...)
}
