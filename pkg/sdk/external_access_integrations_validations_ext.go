package sdk

func (opts *ExternalAccessIntegrationSet) additionalValidations() error {
	if opts.AllowedNetworkRules != nil && len(opts.AllowedNetworkRules) == 0 {
		return NewError("AllowedNetworkRules must not be empty when provided")
	}
	return nil
}
