package testacc

const (
	// Aws constants
	awsRoleArn       = "arn:aws:iam::000000000001:role/test"
	awsAllowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	awsBlockedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/blocked/"

	// Azure constants
	azureTenantId        = "00000000-0000-0000-0000-000000000000"
	azureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	azureAllowedPrefix   = "https://apim-hello-world.azure-api.net/dev"
	azureBlockedPrefix   = "https://apim-hello-world.azure-api.net/blocked"

	// Google constants
	googleAudience      = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	googleAllowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	googleBlockedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod/blocked/"
)
