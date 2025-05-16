package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type FileFormatClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFileFormatClient(context *TestClientContext, idsGenerator *IdsGenerator) *FileFormatClient {
	return &FileFormatClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FileFormatClient) client() sdk.FileFormats {
	return c.context.client.FileFormats
}

func (c *FileFormatClient) CreateFileFormat(t *testing.T) (*sdk.FileFormat, func()) {
	t.Helper()
	return c.CreateFileFormatWithOptions(t, &sdk.CreateFileFormatOptions{
		Type: sdk.FileFormatTypeCSV,
	})
}

func (c *FileFormatClient) CreateFileFormatWithOptions(t *testing.T, opts *sdk.CreateFileFormatOptions) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) DropFileFormatFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropFileFormatOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
