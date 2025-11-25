package sdk

type ApiIntegrationAwsApiProviderType string

var (
	ApiIntegrationAwsApiGateway           ApiIntegrationAwsApiProviderType = "aws_api_gateway"
	ApiIntegrationAwsPrivateApiGateway    ApiIntegrationAwsApiProviderType = "aws_private_api_gateway"
	ApiIntegrationAwsGovApiGateway        ApiIntegrationAwsApiProviderType = "aws_gov_api_gateway"
	ApiIntegrationAwsGovPrivateApiGateway ApiIntegrationAwsApiProviderType = "aws_gov_private_api_gateway"
)
