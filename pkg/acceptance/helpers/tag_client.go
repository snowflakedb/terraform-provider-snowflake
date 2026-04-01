package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

// SetupConflict creates two tables with conflicting tag values and a view that joins them,
// producing a tag propagation conflict on the view.
func (c *TagClient) SetupConflict(t *testing.T, tagId sdk.SchemaObjectIdentifier, value1, value2 string) *sdk.View {
	t.Helper()
	ctx := context.Background()

	tableClient := NewTableClient(c.context, c.ids)
	viewClient := NewViewClient(c.context, c.ids)

	table1, table1Cleanup := tableClient.Create(t)
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := tableClient.Create(t)
	t.Cleanup(table2Cleanup)

	err := c.client().Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table1.ID()).WithSetTags([]sdk.TagAssociation{
		{Name: tagId, Value: value1},
	}))
	require.NoError(t, err)

	err = c.client().Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table2.ID()).WithSetTags([]sdk.TagAssociation{
		{Name: tagId, Value: value2},
	}))
	require.NoError(t, err)

	view, viewCleanup := viewClient.CreateView(t, fmt.Sprintf(
		"SELECT t1.id AS id1, t2.id AS id2 FROM %s t1 JOIN %s t2 ON t1.id = t2.id",
		table1.ID().FullyQualifiedName(),
		table2.ID().FullyQualifiedName(),
	))
	t.Cleanup(viewCleanup)

	return view
}
