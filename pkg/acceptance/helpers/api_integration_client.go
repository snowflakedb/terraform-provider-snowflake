package helpers

import (
	"context"
	"fmt"
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

type ApiIntegrationAllClient struct {
	context *TestClientContext
}

func NewApiIntegrationAllClient(context *TestClientContext) *ApiIntegrationAllClient {
	return &ApiIntegrationAllClient{context: context}
}

func (c *ApiIntegrationAllClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationAllDetails, error) {
	t.Helper()
	return c.context.client.ApiIntegrations.DescribeDetails(context.Background(), id)
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

// TODO(SNOW-1348334): change raw sqls to proper client
func (c *ApiIntegrationClient) CreateApiIntegrationForGitRepository(t *testing.T, origin string) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`CREATE OR REPLACE API INTEGRATION %s
	  API_PROVIDER = GIT_HTTPS_API
	  API_ALLOWED_PREFIXES = ('%s')
	  ALLOWED_AUTHENTICATION_SECRETS = ALL
	  ENABLED = TRUE;`, id.FullyQualifiedName(), origin))
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
