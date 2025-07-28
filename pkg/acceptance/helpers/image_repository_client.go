package helpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ImageRepositoryClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewImageRepositoryClient(context *TestClientContext, idsGenerator *IdsGenerator) *ImageRepositoryClient {
	return &ImageRepositoryClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ImageRepositoryClient) client() sdk.ImageRepositories {
	return c.context.client.ImageRepositories
}

func (c *ImageRepositoryClient) Create(t *testing.T) (*sdk.ImageRepository, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(c.ids.RandomSchemaObjectIdentifier()))
}

// TODO(SNOW-2070746): Image repositories cannot be created in the default schema with lowercase letters and the schema must be provided for now.
func (c *ImageRepositoryClient) CreateInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.ImageRepository, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	return c.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(id))
}

func (c *ImageRepositoryClient) CreateWithRequest(t *testing.T, req *sdk.CreateImageRepositoryRequest) (*sdk.ImageRepository, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	imageRepository, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return imageRepository, c.DropImageRepositoryFunc(t, req.GetName())
}

func (c *ImageRepositoryClient) DropImageRepositoryFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropImageRepositoryRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ImageRepositoryClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.ImageRepository, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ImageRepositoryClient) Alter(t *testing.T, req *sdk.AlterImageRepositoryRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}
