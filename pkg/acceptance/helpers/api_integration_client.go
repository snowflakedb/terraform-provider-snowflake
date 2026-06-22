package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApiIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewApiIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *ApiIntegrationClient {
	return &ApiIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ApiIntegrationClient) client() sdk.ApiIntegrations {
	return c.context.client.ApiIntegrations
}

func (c *ApiIntegrationClient) Create(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	apiAllowedPrefixes := []sdk.ApiIntegrationEndpointPrefix{{Path: "https://xyz.execute-api.us-west-2.amazonaws.com/production"}}
	req := sdk.NewCreateApiIntegrationRequest(id, apiAllowedPrefixes, true)
	req.WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, "arn:aws:iam::123456789012:role/hello_cloud_account_role"))

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	apiIntegration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return apiIntegration, c.DropApiIntegrationFunc(t, id)
}

func (c *ApiIntegrationClient) CreateWithRequest(t *testing.T, request *sdk.CreateApiIntegrationRequest) (*sdk.ApiIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return integration, c.DropApiIntegrationFunc(t, request.GetName())
}

func (c *ApiIntegrationClient) CreateAws(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://123456.execute-api.us-west-2.amazonaws.com/dev/"}}, true).
		WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, "arn:aws:iam::000000000001:/role/test")),
	)
}

func (c *ApiIntegrationClient) CreateAzure(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()

	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://apim-hello-world.azure-api.net/dev"}}, true).
		WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest("00000000-0000-0000-0000-000000000000", "11111111-1111-1111-1111-111111111111")),
	)
}

func (c *ApiIntegrationClient) CreateGoogle(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://gateway-id-123456.uc.gateway.dev/prod"}}, true).
		WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest("api-gateway-id-123456.apigateway.gcp-project.cloud.goog")),
	)
}

func (c *ApiIntegrationClient) CreateGitToken(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/my-org/"}}, true).
		WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
			WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))),
	)
}

func (c *ApiIntegrationClient) CreateGitTokenWithAllowedOrigin(t *testing.T, origin string) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	integration, cleanup := c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: origin}}, true).
		WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
			WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))),
	)

	return integration.ID(), cleanup
}

func (c *ApiIntegrationClient) CreateGitGithubApp(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/my-org/"}}, true).
		WithGitHttpsApiGithubAppProviderParams(*sdk.NewGitHttpsApiGithubAppParamsRequest()),
	)
}

func (c *ApiIntegrationClient) CreateGitOAuth2(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	auth := sdk.NewOAuth2GitUserAuthenticationRequest("https://auth.example.com/authorize", "https://auth.example.com/token", "oauth-client-id-123", "oauth-client-secret-456")
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/my-org/"}}, true).
		WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().WithApiUserAuthentication(*auth)),
	)
}

func (c *ApiIntegrationClient) CreateGitPrivateLink(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://github.com/my-org/"}}, true).
		WithGitHttpsApiPrivateLinkProviderParams(*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true).
			WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))),
	)
}

func (c *ApiIntegrationClient) CreateMcpOAuth2(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	auth := sdk.NewOAuth2McpUserAuthenticationRequest("oauth-client-id-123", "oauth-client-secret-456", "https://auth.example.com/token", "https://auth.example.com/authorize")
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://mcp.example.com/api/"}}, true).
		WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*auth)),
	)
}

func (c *ApiIntegrationClient) CreateMcpDynamicClient(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	auth := sdk.NewDynamicClientMcpUserAuthenticationRequest("https://mcp.atlassian.com/v1/mcp")
	return c.CreateWithRequest(t, sdk.NewCreateApiIntegrationRequest(id,
		[]sdk.ApiIntegrationEndpointPrefix{{Path: "https://mcp.atlassian.com/v1/mcp"}}, true).
		WithExternalMcpDynamicClientProviderParams(*sdk.NewExternalMcpDynamicClientParamsRequest().WithApiUserAuthentication(*auth)),
	)
}

func (c *ApiIntegrationClient) DropApiIntegrationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApiIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ApiIntegrationClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegration, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByIDSafely(ctx, id)
}

func (c *ApiIntegrationClient) Alter(t *testing.T, request *sdk.AlterApiIntegrationRequest) {
	t.Helper()
	err := c.client().Alter(context.Background(), request)
	require.NoError(t, err)
}

func (c *ApiIntegrationClient) DescribeAws(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationAwsDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeAwsDetails(context.Background(), id)
}

func (c *ApiIntegrationClient) DescribeAzure(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationAzureDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeAzureDetails(context.Background(), id)
}

func (c *ApiIntegrationClient) DescribeGoogle(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationGoogleDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeGoogleDetails(context.Background(), id)
}

func (c *ApiIntegrationClient) DescribeGitHttpsApi(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationGitHttpsApiDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeGitHttpsApiDetails(context.Background(), id)
}

func (c *ApiIntegrationClient) DescribeExternalMcp(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationExternalMcpDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeExternalMcpDetails(context.Background(), id)
}
