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

	// Git constants
	gitAllowedPrefix = "https://github.com/my-org/"
	gitBlockedPrefix = "https://github.com/my-org/blocked/"

	// MCP constants
	mcpAllowedPrefix    = "https://mcp.example.com/api/"
	mcpBlockedPrefix    = "https://mcp.example.com/api/blocked/"
	mcpOauthResourceUrl = "https://mcp.atlassian.com/v1/mcp"

	// Git OAuth2 constants
	gitOauth2AuthorizationEndpoint         = "https://auth.example.com/authorize"
	gitOauth2TokenEndpoint                 = "https://auth.example.com/token" //nolint:gosec
	gitOauth2ClientId                      = "oauth-client-id-123"
	gitOauth2ClientSecret                  = "oauth-client-secret-456" //nolint:gosec
	gitOauth2ExternalAuthorizationEndpoint = "https://different.example.com/authorize"
	gitOauth2ExternalTokenEndpoint         = "https://different.example.com/token" //nolint:gosec
	gitOauth2ExternalClientId              = "different-client-id"
	gitOauth2UpdatedAuthorizationEndpoint  = "https://updated.example.com/authorize"
	gitOauth2UpdatedTokenEndpoint          = "https://updated.example.com/token" //nolint:gosec
	gitOauth2AccessTokenValidity           = 3600
	gitOauth2RefreshTokenValidity          = 86400
	gitOauth2Username                      = "test_user"
)
