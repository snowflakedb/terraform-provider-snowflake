package resourceshowoutputassert

func (e *ExternalAccessIntegrationDescribeOutputAssert) HasAllowedNetworkRules(expected ...string) *ExternalAccessIntegrationDescribeOutputAssert {
	e.SetContainsExactlyStringValues("allowed_network_rules", expected...)
	return e
}

func (e *ExternalAccessIntegrationDescribeOutputAssert) HasAllowedApiAuthenticationIntegrations(expected ...string) *ExternalAccessIntegrationDescribeOutputAssert {
	e.SetContainsExactlyStringValues("allowed_api_authentication_integrations", expected...)
	return e
}

func (e *ExternalAccessIntegrationDescribeOutputAssert) HasAllowedAuthenticationSecrets(expected ...string) *ExternalAccessIntegrationDescribeOutputAssert {
	e.SetContainsExactlyStringValues("allowed_authentication_secrets", expected...)
	return e
}
