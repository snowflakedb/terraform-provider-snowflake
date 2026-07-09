package resourceshowoutputassert

import "fmt"

func (a *ApiIntegrationAllDescribeOutputAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationAllDescribeOutputAssert {
	a.StringValueSet("allowed_prefixes.#", fmt.Sprintf("%d", len(expected)))
	for i, v := range expected {
		a.StringValueSet(fmt.Sprintf("allowed_prefixes.%d", i), v)
	}
	return a
}

func (a *ApiIntegrationAllDescribeOutputAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationAllDescribeOutputAssert {
	a.StringValueSet("blocked_prefixes.#", fmt.Sprintf("%d", len(expected)))
	for i, v := range expected {
		a.StringValueSet(fmt.Sprintf("blocked_prefixes.%d", i), v)
	}
	return a
}

func (a *ApiIntegrationAllDescribeOutputAssert) HasApiAwsIamUserArnNotEmpty() *ApiIntegrationAllDescribeOutputAssert {
	a.ValuePresent("api_aws_iam_user_arn")
	return a
}

func (a *ApiIntegrationAllDescribeOutputAssert) HasApiAwsExternalIdNotEmpty() *ApiIntegrationAllDescribeOutputAssert {
	a.ValuePresent("api_aws_external_id")
	return a
}

func (a *ApiIntegrationAllDescribeOutputAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationAllDescribeOutputAssert {
	a.ValuePresent("google_api_service_account")
	return a
}
