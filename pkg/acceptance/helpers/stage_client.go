package helpers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	externalS3Req := sdk.NewExternalS3StageParamsRequest(nycWeatherDataURL)

	err := c.client().CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id, *externalS3Req))
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStageWithDirectory(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id).WithDirectoryTableOptions(*sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(true)))
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

func (c *StageClient) CreateStageOnS3WithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateOnS3StageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnS3(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStageOnGCSWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateOnGCSStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnGCS(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStageOnAzureWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateOnAzureStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnAzure(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStageOnS3CompatibleWithRequest(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateOnS3CompatibleStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnS3Compatible(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnS3 creates an S3 stage with sane defaults using pre-created storage integration.
func (c *StageClient) CreateStageOnS3(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	id := c.ids.RandomSchemaObjectIdentifier()

	s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
		WithStorageIntegration(ids.PrecreatedS3StorageIntegration)
	request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

	err := c.client().CreateOnS3(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnS3WithCredentials creates an S3 stage using AWS credentials from env vars.
func (c *StageClient) CreateStageOnS3WithCredentials(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)
	id := c.ids.RandomSchemaObjectIdentifier()

	s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
		WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
			WithAwsKeyId(awsKeyId).
			WithAwsSecretKey(awsSecretKey))
	request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

	err := c.client().CreateOnS3(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnGCS creates a GCS stage with sane defaults using pre-created storage integration.
func (c *StageClient) CreateStageOnGCS(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	id := c.ids.RandomSchemaObjectIdentifier()

	gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
		WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)
	request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq)

	err := c.client().CreateOnGCS(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnAzure creates an Azure stage with sane defaults using pre-created storage integration.
func (c *StageClient) CreateStageOnAzure(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	id := c.ids.RandomSchemaObjectIdentifier()

	azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
		WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)
	request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

	err := c.client().CreateOnAzure(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnAzureWithCredentials creates an Azure stage using SAS token from env vars.
func (c *StageClient) CreateStageOnAzureWithCredentials(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureSasToken := testenvs.GetOrSkipTest(t, testenvs.AzureExternalSasToken)
	id := c.ids.RandomSchemaObjectIdentifier()

	azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
		WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))
	request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

	err := c.client().CreateOnAzure(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

// CreateStageOnS3Compatible creates an S3-compatible stage with sane defaults.
func (c *StageClient) CreateStageOnS3Compatible(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
	endpoint := "s3.us-west-2.amazonaws.com"
	id := c.ids.RandomSchemaObjectIdentifier()

	s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
		WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))
	request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq)

	err := c.client().CreateOnS3Compatible(ctx, request)
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) DropStageFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStageRequest(id).WithIfExists(true))
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

	path := testfiles.TestFile(t, filename, []byte(content))

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @~/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, path))
	require.NoError(t, err)

	t.Cleanup(c.RemoveFromUserStageFunc(t, path))

	return path
}

func (c *StageClient) PutOnStageWithPath(t *testing.T, id sdk.SchemaObjectIdentifier, stageLocation string, filename string) string {
	t.Helper()
	ctx := context.Background()

	filePath := filepath.Join("./testdata/", filename)
	absPath, err := filepath.Abs(filePath)
	require.NoError(t, err)

	stagePath := filepath.Join(stageLocation, filename)

	c.putInLocation(ctx, t, absPath, filePath, fmt.Sprintf("@%s/%s", id.FullyQualifiedName(), stagePath))

	return stagePath
}

func (c *StageClient) PutInLocationWithContent(t *testing.T, stageLocation string, filename string, content string) string {
	t.Helper()
	ctx := context.Background()

	filePath := testfiles.TestFile(t, filename, []byte(content))

	c.putInLocation(ctx, t, filePath, filename, stageLocation)

	return filePath
}

func (c *StageClient) putInLocation(ctx context.Context, t *testing.T, filePath string, filename string, location string) {
	t.Helper()
	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s %s AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, filePath, location))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`REMOVE %s/%s`, location, filename))
		// Only check the error if it's not related to the stage / file existence or access
		if !errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			require.NoError(t, err)
		}
	})
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

	filePath := testfiles.TestFile(t, filename, []byte(content))

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

	err := c.client().Alter(ctx, sdk.NewAlterStageRequest(id).WithRenameTo(newId))
	require.NoError(t, err)
}

func (c *StageClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) ([]sdk.StageProperty, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}

func (c *StageClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Stage, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
