package resourceshowoutputassert

func (c *CatalogIntegrationAwsGlueDescribeOutputAssert) HasGlueAwsIamUserArnNotEmpty() *CatalogIntegrationAwsGlueDescribeOutputAssert {
	c.ValuePresent("glue_aws_iam_user_arn")
	return c
}

func (c *CatalogIntegrationAwsGlueDescribeOutputAssert) HasGlueAwsExternalIdNotEmpty() *CatalogIntegrationAwsGlueDescribeOutputAssert {
	c.ValuePresent("glue_aws_external_id")
	return c
}
