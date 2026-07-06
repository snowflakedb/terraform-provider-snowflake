package resourceshowoutputassert

func (a *ApiIntegrationAzureDescribeOutputAssert) HasApiKeyNotEmpty() *ApiIntegrationAzureDescribeOutputAssert {
	a.ValuePresent("api_key")
	return a
}

func (a *ApiIntegrationAzureDescribeOutputAssert) HasAzureMultiTenantAppNameNotEmpty() *ApiIntegrationAzureDescribeOutputAssert {
	a.ValuePresent("azure_multi_tenant_app_name")
	return a
}

func (a *ApiIntegrationAzureDescribeOutputAssert) HasAzureConsentUrlNotEmpty() *ApiIntegrationAzureDescribeOutputAssert {
	a.ValuePresent("azure_consent_url")
	return a
}
