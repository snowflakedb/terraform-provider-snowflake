package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (c *FileFormatClient) CreateCsv(t *testing.T) (*sdk.FileFormat, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateCsvWithRequest(t, id, sdk.NewCreateCsvFileFormatRequest(id))
}

func (c *FileFormatClient) CreateCsvWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateCsvFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateCsv(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) CreateJsonWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateJsonFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateJson(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) CreateAvroWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateAvroFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateAvro(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) CreateOrcWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateOrcFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOrc(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) CreateParquetWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateParquetFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateParquet(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) CreateXmlWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateXmlFileFormatRequest) (*sdk.FileFormat, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateXml(ctx, request)
	require.NoError(t, err)

	fileFormat, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return fileFormat, c.DropFileFormatFunc(t, id)
}

func (c *FileFormatClient) AlterCsv(t *testing.T, request *sdk.AlterCsvFileFormatRequest) {
	t.Helper()
	err := c.client().AlterCsv(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) AlterJson(t *testing.T, request *sdk.AlterJsonFileFormatRequest) {
	t.Helper()
	err := c.client().AlterJson(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) AlterAvro(t *testing.T, request *sdk.AlterAvroFileFormatRequest) {
	t.Helper()
	err := c.client().AlterAvro(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) AlterOrc(t *testing.T, request *sdk.AlterOrcFileFormatRequest) {
	t.Helper()
	err := c.client().AlterOrc(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) AlterParquet(t *testing.T, request *sdk.AlterParquetFileFormatRequest) {
	t.Helper()
	err := c.client().AlterParquet(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) AlterXml(t *testing.T, request *sdk.AlterXmlFileFormatRequest) {
	t.Helper()
	err := c.client().AlterXml(context.Background(), request)
	require.NoError(t, err)
}

func (c *FileFormatClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormat, error) {
	t.Helper()
	return c.client().ShowByIDSafely(context.Background(), id)
}

func (c *FileFormatClient) DescribeCsvDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatCsv, error) {
	t.Helper()
	return c.client().DescribeCsvDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeJsonDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatJson, error) {
	t.Helper()
	return c.client().DescribeJsonDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeAvroDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatAvro, error) {
	t.Helper()
	return c.client().DescribeAvroDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeOrcDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatOrc, error) {
	t.Helper()
	return c.client().DescribeOrcDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeParquetDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatParquet, error) {
	t.Helper()
	return c.client().DescribeParquetDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeXmlDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatXml, error) {
	t.Helper()
	return c.client().DescribeXmlDetails(context.Background(), id)
}

func (c *FileFormatClient) DescribeAllDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.FileFormatAllDetails, error) {
	t.Helper()
	return c.client().DescribeAllDetails(context.Background(), id)
}

func (c *FileFormatClient) DropFileFormatFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropFileFormatRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
