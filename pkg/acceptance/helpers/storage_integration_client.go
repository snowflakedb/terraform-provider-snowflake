package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type StorageIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStorageIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *StorageIntegrationClient {
	return &StorageIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StorageIntegrationClient) client() sdk.StorageIntegrations {
	return c.context.client.StorageIntegrations
}

func (c *StorageIntegrationClient) CreateS3(t *testing.T, awsBucketUrl, awsRoleArn string) (*sdk.StorageIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	allowedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/allowed-location",
			},
			{
				Path: prefix + "/allowed-location2",
			},
		}
	}
	s3AllowedLocations := allowedLocations(awsBucketUrl)

	blockedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/blocked-location",
			},
			{
				Path: prefix + "/blocked-location2",
			},
		}
	}
	s3BlockedLocations := blockedLocations(awsBucketUrl)

	id := c.ids.RandomAccountObjectIdentifier()
	req := sdk.NewCreateStorageIntegrationRequest(id, true, s3AllowedLocations).
		WithIfNotExists(true).
		WithS3StorageProviderParams(*sdk.NewS3StorageParamsRequest(sdk.RegularS3Protocol, awsRoleArn)).
		WithStorageBlockedLocations(s3BlockedLocations).
		WithComment("some comment")

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return integration, c.DropFunc(t, id)
}

func (c *StorageIntegrationClient) CreateAzure(t *testing.T, azureBucketUrl string, azureTenantId string) (*sdk.StorageIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	allowedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/allowed-location",
			},
			{
				Path: prefix + "/allowed-location2",
			},
		}
	}
	azureAllowedLocations := allowedLocations(azureBucketUrl)

	id := c.ids.RandomAccountObjectIdentifier()
	req := sdk.NewCreateStorageIntegrationRequest(id, true, azureAllowedLocations).
		WithAzureStorageProviderParams(*sdk.NewAzureStorageParamsRequest(azureTenantId))

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return integration, c.DropFunc(t, id)
}

func (c *StorageIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStorageIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *StorageIntegrationClient) Alter(t *testing.T, request *sdk.AlterStorageIntegrationRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *StorageIntegrationClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.StorageIntegration, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *StorageIntegrationClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.StorageIntegrationProperty, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}

func (c *StorageIntegrationClient) DescribeAws(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.StorageIntegrationAwsDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeAwsDetails(ctx, id)
}

func (c *StorageIntegrationClient) DescribeAzure(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.StorageIntegrationAzureDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeAzureDetails(ctx, id)
}

func (c *StorageIntegrationClient) DescribeGcs(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.StorageIntegrationGcsDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeGcsDetails(ctx, id)
}

func (c *StorageIntegrationClient) CreateWithoutEnabled(t *testing.T, id sdk.AccountObjectIdentifier, iamRole string, allowedLocation sdk.StorageLocation) error {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`CREATE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = '%s' STORAGE_ALLOWED_LOCATIONS = ('%s')`, id.FullyQualifiedName(), iamRole, allowedLocation.Path))
	t.Cleanup(c.DropFunc(t, id))

	return err
}
