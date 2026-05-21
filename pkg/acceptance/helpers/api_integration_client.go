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

func (c *ApiIntegrationClient) CreateApiIntegration(t *testing.T) (*sdk.ApiIntegration, func()) {
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

func (c *ApiIntegrationClient) CreateApiIntegrationForGitRepository(t *testing.T, origin string) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	apiAllowedPrefixes := []sdk.ApiIntegrationEndpointPrefix{{Path: origin}}
	gitParams := sdk.NewGitHttpsApiGithubAppParamsRequest().
		WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))
	req := sdk.NewCreateApiIntegrationRequest(id, apiAllowedPrefixes, true).
		WithGitHttpsApiGithubAppProviderParams(*gitParams)

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	return id, c.DropApiIntegrationFunc(t, id)
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
