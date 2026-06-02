package sdk

func (r streamlitsDetailRow) additionalConvert(result *StreamlitDetail) error {
	result.ExternalAccessIntegrations = ParseCommaSeparatedStringArray(r.ExternalAccessIntegrations, false)
	externalAccessIntegrations := make([]string, len(result.ExternalAccessIntegrations))
	for i, v := range result.ExternalAccessIntegrations {
		externalAccessIntegrations[i] = NewObjectIdentifierFromFullyQualifiedName(v).Name()
	}
	result.ExternalAccessIntegrations = externalAccessIntegrations
	return nil
}

func (opts *ShowStreamlitOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	return nil
}
