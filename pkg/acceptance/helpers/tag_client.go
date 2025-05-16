package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TagClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewTagClient(context *TestClientContext, idsGenerator *IdsGenerator) *TagClient {
	return &TagClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *TagClient) client() sdk.Tags {
	return c.context.client.Tags
}

func (c *TagClient) CreateTag(t *testing.T) (*sdk.Tag, func()) {
	t.Helper()
	return c.CreateTagInSchema(t, c.ids.SchemaId())
}

func (c *TagClient) CreateTagWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Tag, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateTagRequest(id))
}

func (c *TagClient) CreateTagInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Tag, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateTagRequest(c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)))
}

func (c *TagClient) CreateWithRequest(t *testing.T, req *sdk.CreateTagRequest) (*sdk.Tag, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	tag, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return tag, c.DropTagFunc(t, req.GetName())
}

func (c *TagClient) Unset(t *testing.T, objectType sdk.ObjectType, id sdk.ObjectIdentifier, unsetTags []sdk.ObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Unset(ctx, sdk.NewUnsetTagRequest(objectType, id).WithUnsetTags(unsetTags))
	require.NoError(t, err)
}

func (c *TagClient) Set(t *testing.T, objectType sdk.ObjectType, id sdk.ObjectIdentifier, setTags []sdk.TagAssociation) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Set(ctx, sdk.NewSetTagRequest(objectType, id).WithSetTags(setTags))
	require.NoError(t, err)
}

func (c *TagClient) Alter(t *testing.T, req *sdk.AlterTagRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *TagClient) DropTagFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropTagRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *TagClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Tag, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *TagClient) GetForObject(t *testing.T, tagId sdk.SchemaObjectIdentifier, objectId sdk.ObjectIdentifier, objectType sdk.ObjectType) (*string, error) {
	t.Helper()
	ctx := context.Background()
	client := c.context.client.SystemFunctions

	return client.GetTag(ctx, tagId, objectId, objectType)
}
