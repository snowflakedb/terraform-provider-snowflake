package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SecretClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSecretClient(context *TestClientContext, idsGenerator *IdsGenerator) *SecretClient {
	return &SecretClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SecretClient) client() sdk.Secrets {
	return c.context.client.Secrets
}

func (c *SecretClient) CreateWithOAuthClientCredentialsFlow(t *testing.T, id sdk.SchemaObjectIdentifier, apiIntegration sdk.AccountObjectIdentifier, oauthScopes []sdk.ApiIntegrationScope) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, apiIntegration).
		WithOauthScopes(sdk.OauthScopesListRequest{OauthScopesList: oauthScopes})

	err := c.client().CreateWithOAuthClientCredentialsFlow(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.DropFunc(t, id)
}

func (c *SecretClient) CreateWithOAuthAuthorizationCodeFlow(t *testing.T, id sdk.SchemaObjectIdentifier, apiIntegration sdk.AccountObjectIdentifier, refreshToken, refreshTokenExpiryTime string) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, apiIntegration)

	err := c.client().CreateWithOAuthAuthorizationCodeFlow(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.DropFunc(t, id)
}

func (c *SecretClient) CreateWithBasicAuthenticationFlow(t *testing.T, id sdk.SchemaObjectIdentifier, username, password string) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)

	err := c.client().CreateWithBasicAuthentication(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.DropFunc(t, id)
}

func (c *SecretClient) CreateWithGenericString(t *testing.T, id sdk.SchemaObjectIdentifier, secretString string) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithGenericStringSecretRequest(id, secretString)

	err := c.client().CreateWithGenericString(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.DropFunc(t, id)
}

func (c *SecretClient) Alter(t *testing.T, req *sdk.AlterSecretRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *SecretClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	}
}

func (c *SecretClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Secret, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
