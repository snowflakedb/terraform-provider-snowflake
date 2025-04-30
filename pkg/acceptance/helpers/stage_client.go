package helpers

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/terraform-provider-snowflake/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

const (
	nycWeatherDataURL = "s3://snowflake-workshop-lab/weather-nyc"
)

type StageClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStageClient(context *TestClientContext, idsGenerator *IdsGenerator) *StageClient {
	return &StageClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StageClient) client() sdk.Stages {
	return c.context.client.Stages
}

func (c *StageClient) CreateStageWithURL(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id).
		WithExternalStageParams(sdk.NewExternalS3StageParamsRequest(nycWeatherDataURL)))
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStageWithDirectory(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id).WithDirectoryTableOptions(sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(sdk.Bool(true))))
}

func (c *StageClient) CreateStage(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	return c.CreateStageInSchema(t, c.ids.SchemaId())
}

func (c *StageClient) CreateStageInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Stage, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	return c.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id))
}

func (c *StageClient) CreateStageWithRequest(t *testing.T, request *sdk.CreateInternalStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateInternal(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, request.ID())
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, request.ID())
}

func (c *StageClient) DropStageFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}

func (c *StageClient) PutOnStage(t *testing.T, id sdk.SchemaObjectIdentifier, filename string) {
	t.Helper()
	ctx := context.Background()

	path, err := filepath.Abs("./testdata/" + filename)
	require.NoError(t, err)
	absPath := "file://" + path

	_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT '%s' @%s AUTO_COMPRESS = FALSE`, absPath, id.FullyQualifiedName()))
	require.NoError(t, err)
}

func (c *StageClient) PutOnUserStageWithContent(t *testing.T, filename string, content string) string {
	t.Helper()
	ctx := context.Background()

	path := testhelpers.TestFile(t, filename, []byte(content))

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @~/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, path))
	require.NoError(t, err)

	t.Cleanup(c.RemoveFromUserStageFunc(t, path))

	return path
}

func (c *StageClient) RemoveFromUserStage(t *testing.T, pathOnStage string) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`REMOVE @~/%s`, pathOnStage))
	require.NoError(t, err)
}

func (c *StageClient) RemoveFromUserStageFunc(t *testing.T, pathOnStage string) func() {
	t.Helper()
	return func() {
		c.RemoveFromUserStage(t, pathOnStage)
	}
}

func (c *StageClient) RemoveFromStage(t *testing.T, stageLocation string, pathOnStage string) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`REMOVE %s/%s`, stageLocation, pathOnStage))
	require.NoError(t, err)
}

func (c *StageClient) RemoveFromStageFunc(t *testing.T, stageLocation string, pathOnStage string) func() {
	t.Helper()
	return func() {
		c.RemoveFromStage(t, stageLocation, pathOnStage)
	}
}

func (c *StageClient) PutOnStageWithContent(t *testing.T, id sdk.SchemaObjectIdentifier, filename string, content string) {
	t.Helper()
	ctx := context.Background()

	filePath := testhelpers.TestFile(t, filename, []byte(content))

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @%s AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, filePath, id.FullyQualifiedName()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`REMOVE @%s/%s`, id.FullyQualifiedName(), filename))
		require.NoError(t, err)
	})
}

func (c *StageClient) CopyIntoTableFromFile(t *testing.T, table, stage sdk.SchemaObjectIdentifier, filename string) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`COPY INTO %s
	FROM @%s/%s
	FILE_FORMAT = (type=json)
	MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE`, table.FullyQualifiedName(), stage.FullyQualifiedName(), filename))
	require.NoError(t, err)
}

func (c *StageClient) Rename(t *testing.T, id sdk.SchemaObjectIdentifier, newId sdk.SchemaObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterStageRequest(id).WithRenameTo(&newId))
	require.NoError(t, err)
}

func (c *StageClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) ([]sdk.StageProperty, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}
