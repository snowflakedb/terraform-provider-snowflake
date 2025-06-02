package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalVolumeClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalVolumeClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalVolumeClient {
	return &ExternalVolumeClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalVolumeClient) client() sdk.ExternalVolumes {
	return c.context.client.ExternalVolumes
}

// TODO(SNOW-999142): Switch to returning *sdk.ExternalVolume. Need to update existing acceptance tests for this.
func (c *ExternalVolumeClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	kmsKeyId := "1234abcd-12ab-34cd-56ef-1234567890ab"
	storageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:              "my-s3-us-west-2",
				StorageProvider:   "S3",
				StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
				StorageBaseUrl:    "s3://my-example-bucket/",
				Encryption: &sdk.ExternalVolumeS3Encryption{
					Type:     "AWS_SSE_KMS",
					KmsKeyId: &kmsKeyId,
				},
			},
		},
	}

	req := sdk.NewCreateExternalVolumeRequest(id, storageLocations)
	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	_, showErr := c.client().ShowByID(ctx, id)
	require.NoError(t, showErr)

	return id, c.DropFunc(t, id)
}

func (c *ExternalVolumeClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolume, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ExternalVolumeClient) Alter(t *testing.T, req *sdk.AlterExternalVolumeRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ExternalVolumeClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropExternalVolumeRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
