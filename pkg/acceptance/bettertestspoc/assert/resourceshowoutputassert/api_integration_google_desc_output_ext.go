package resourceshowoutputassert

func (a *ApiIntegrationGoogleDescribeOutputAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationGoogleDescribeOutputAssert {
	a.ValuePresent("google_api_service_account")
	return a
}
