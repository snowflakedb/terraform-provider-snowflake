package resourceshowoutputassert

func (a *ApiIntegrationAwsDescribeOutputAssert) HasApiKeyNotEmpty() *ApiIntegrationAwsDescribeOutputAssert {
	a.ValuePresent("api_key")
	return a
}

func (a *ApiIntegrationAwsDescribeOutputAssert) HasApiAwsIamUserArnNotEmpty() *ApiIntegrationAwsDescribeOutputAssert {
	a.ValuePresent("api_aws_iam_user_arn")
	return a
}

func (a *ApiIntegrationAwsDescribeOutputAssert) HasApiAwsExternalIdNotEmpty() *ApiIntegrationAwsDescribeOutputAssert {
	a.ValuePresent("api_aws_external_id")
	return a
}
