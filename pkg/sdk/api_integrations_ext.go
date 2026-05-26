package sdk

func (s *ApiIntegrationSet) additionalValidations() error {
	if moreThanOneValueSet(s.AwsParams, s.AzureParams, s.GoogleParams) {
		return errMoreThanOneOf("AlterApiIntegrationOptions.Set", "AwsParams", "AzureParams", "GoogleParams")
	}
	return nil
}

func (r *CreateApiIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}
