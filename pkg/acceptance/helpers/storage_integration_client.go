package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

func (c *StorageIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStorageIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
