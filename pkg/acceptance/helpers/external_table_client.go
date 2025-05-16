package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalTableClient {
	return &ExternalTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalTableClient) client() sdk.ExternalTables {
	return c.context.client.ExternalTables
}

func (c *ExternalTableClient) PublishDataToStage(t *testing.T, stageId sdk.SchemaObjectIdentifier, data []byte) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`copy into @%s/external_tables_test_data/test_data from (select parse_json('%s')) overwrite = true`, stageId.FullyQualifiedName(), string(data)))
	require.NoError(t, err)
}

func (c *ExternalTableClient) CreateWithLocation(t *testing.T, location string) (*sdk.ExternalTable, func()) {
	t.Helper()

	externalTableId := c.ids.RandomSchemaObjectIdentifier()
	req := sdk.NewCreateExternalTableRequest(externalTableId, location).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON)).WithColumns([]*sdk.ExternalTableColumnRequest{sdk.NewExternalTableColumnRequest("id", sdk.DataTypeNumber, "value:time::int")})

	return c.CreateWithRequest(t, req)
}

func (c *ExternalTableClient) CreateInSchemaWithLocation(t *testing.T, location string, schemaId sdk.DatabaseObjectIdentifier) (*sdk.ExternalTable, func()) {
	t.Helper()

	req := sdk.NewCreateExternalTableRequest(c.ids.RandomSchemaObjectIdentifierInSchema(schemaId), location).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON)).WithColumns([]*sdk.ExternalTableColumnRequest{sdk.NewExternalTableColumnRequest("id", sdk.DataTypeNumber, "value:time::int")})

	return c.CreateWithRequest(t, req)
}

func (c *ExternalTableClient) CreateWithRequest(t *testing.T, req *sdk.CreateExternalTableRequest) (*sdk.ExternalTable, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	stream, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return stream, c.DropFunc(t, req.GetName())
}

func (c *ExternalTableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropExternalTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
