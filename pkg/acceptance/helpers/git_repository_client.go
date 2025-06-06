package helpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type GitRepositoryClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewGitRepositoryClient(context *TestClientContext, idsGenerator *IdsGenerator) *GitRepositoryClient {
	return &GitRepositoryClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *GitRepositoryClient) client() sdk.GitRepositories {
	return c.context.client.GitRepositories
}

func (c *GitRepositoryClient) Create(t *testing.T, name sdk.SchemaObjectIdentifier, origin string, apiIntegration sdk.AccountObjectIdentifier) (*sdk.GitRepository, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateGitRepositoryRequest(name, origin, apiIntegration))
}

func (c *GitRepositoryClient) CreateWithRequest(t *testing.T, req *sdk.CreateGitRepositoryRequest) (*sdk.GitRepository, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	gitRepository, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return gitRepository, c.DropFunc(t, req.GetName())
}

func (c *GitRepositoryClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropGitRepositoryRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *GitRepositoryClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.GitRepository, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
